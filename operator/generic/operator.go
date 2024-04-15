package generic_operator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	sdkelcontracts "github.com/Layr-Labs/eigensdk-go/chainio/clients/elcontracts"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	sdkmetrics "github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/Layr-Labs/eigensdk-go/metrics/collectors/economic"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	"github.com/Layr-Labs/eigensdk-go/nodeapi"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"

	"github.com/alt-research/avs/core"
	"github.com/alt-research/avs/core/chainio"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/metrics"
)

const SEM_VER = "0.0.1"

func GenericOperatorMain(cliCtx *cli.Context, ctx context.Context, nodeConfig config.NodeConfig) error {
	nodeConfig.WithEnv()

	avsConfigs, err := config.NewAVSConfigs(cliCtx)
	if err != nil {
		return err
	}

	// FIXME: currently just suppot one avs for a operator node.
	if len(avsConfigs) != 1 {
		return fmt.Errorf("currently the operator node just support one avs")
	}

	avsConfig := avsConfigs[0]

	// FIXME: split node config into two parts, avs config and network/node config
	nodeConfig.AVSName = avsConfig.AVSName

	operator, err := NewOperatorFromConfig(nodeConfig, avsConfig)
	if err != nil {
		return err
	}

	err = operator.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

type Operator struct {
	config           config.NodeConfig
	avsCfg           config.GenericAVSConfig
	logger           sdklogging.Logger
	ethClient        eth.Client
	metricsReg       *prometheus.Registry
	metrics          metrics.Metrics
	nodeApi          *nodeapi.NodeApi
	avsWriter        *chainio.AvsWriter
	avsReader        chainio.AvsReaderer
	eigenlayerReader sdkelcontracts.ELReader
	eigenlayerWriter sdkelcontracts.ELWriter
	blsKeypair       *bls.KeyPair
	operatorId       sdktypes.OperatorId
	operatorAddr     common.Address
	// receive new tasks in this chan (typically from mach service)
	newTaskCreatedChan chan GenericRequest
	// ip address of aggregator
	aggregatorServerIpPortAddr string
	// rpc client to send signed task responses to aggregator
	aggregatorRpcClient *GenericAggregatorClient
	// rpc server for other service to create bls sig task
	rpcServer *GenericRpcServer
}

func NewOperatorFromConfig(cfg config.NodeConfig, avsCfg config.GenericAVSConfig) (*Operator, error) {
	logger, err := core.NewLogger(cfg.Production)
	if err != nil {
		return nil, errors.Wrap(err, "New logger")
	}

	reg := prometheus.NewRegistry()
	eigenMetrics := sdkmetrics.NewEigenMetrics(cfg.AVSName, cfg.EigenMetricsIpPortAddress, reg, logger)
	avsAndEigenMetrics := metrics.NewAvsAndEigenMetrics(cfg.AVSName, eigenMetrics, reg)

	// Setup Node Api
	nodeApi := nodeapi.NewNodeApi(cfg.AVSName, SEM_VER, cfg.NodeApiIpPortAddress, logger)

	var ethRpcClient eth.Client
	if cfg.EnableMetrics {
		rpcCallsCollector := rpccalls.NewCollector(cfg.AVSName, reg)
		ethRpcClient, err = eth.NewInstrumentedClient(cfg.EthRpcUrl, rpcCallsCollector)
		if err != nil {
			logger.Errorf("Cannot create http ethclient", "err", err)
			return nil, err
		}
	} else {
		ethRpcClient, err = eth.NewClient(cfg.EthRpcUrl)
		if err != nil {
			logger.Errorf("Cannot create http ethclient", "err", err)
			return nil, err
		}
	}

	blsKeyPassword, ok := os.LookupEnv("OPERATOR_BLS_KEY_PASSWORD")
	if !ok {
		logger.Warnf("OPERATOR_BLS_KEY_PASSWORD env var not set. using empty string")
	}
	blsKeyPair, err := bls.ReadPrivateKeyFromFile(cfg.BlsPrivateKeyStorePath, blsKeyPassword)
	if err != nil {
		logger.Errorf("Cannot parse bls private key", "err", err)
		return nil, err
	}

	var operatorAddress common.Address
	var avsWriter *chainio.AvsWriter
	var privateKey *ecdsa.PrivateKey

	// just use a default value, it just to avoid build client panic
	// will never use this value to generate tx.
	ecdsaPrivateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		panic(err)
	}
	privateKey = ecdsaPrivateKey

	if cfg.EcdsaPrivateKeyStorePath != "" {
		operatorAddress, err = sdkEcdsa.GetAddressFromKeyStoreFile(cfg.EcdsaPrivateKeyStorePath)
		if err != nil {
			panic(err)
		}
	} else {
		if cfg.OperatorEcdsaAddress == "" {
			return nil, fmt.Errorf("If not use EcdsaPrivateKeyStorePath, must use operator_ecdsa_address or `OPERATOR_ECDSA_ADDRESS` env to select ecdsa address!")
		}

		if !common.IsHexAddress(cfg.OperatorEcdsaAddress) {
			return nil, fmt.Errorf("the operator_ecdsa_address format is not hex address!")
		}

		operatorAddress = common.HexToAddress(cfg.OperatorEcdsaAddress)
	}

	chainioConfig := clients.BuildAllConfig{
		EthHttpUrl:                 cfg.EthRpcUrl,
		EthWsUrl:                   cfg.EthWsUrl,
		RegistryCoordinatorAddr:    avsCfg.AVSRegistryCoordinatorAddress.String(),
		OperatorStateRetrieverAddr: avsCfg.OperatorStateRetrieverAddress.String(),
		AvsName:                    avsCfg.AVSName,
		PromMetricsIpPortAddress:   cfg.EigenMetricsIpPortAddress,
	}
	sdkClients, err := clients.BuildAll(chainioConfig, privateKey, logger)
	if err != nil {
		panic(err)
	}

	avsReader, err := chainio.BuildAvsReader(
		avsCfg.AVSRegistryCoordinatorAddress,
		avsCfg.OperatorStateRetrieverAddress,
		ethRpcClient, logger)
	if err != nil {
		logger.Error("Cannot create AvsReader", "err", err)
		return nil, err
	}

	// We must register the economic metrics separately because they are exported metrics (from jsonrpc or subgraph calls)
	// and not instrumented metrics: see https://prometheus.io/docs/instrumenting/writing_clientlibs/#overall-structure
	quorumNames := map[sdktypes.QuorumNum]string{
		0: "quorum0",
	}
	economicMetricsCollector := economic.NewCollector(
		sdkClients.ElChainReader, sdkClients.AvsRegistryChainReader,
		cfg.AVSName, logger, operatorAddress, quorumNames)
	reg.MustRegister(economicMetricsCollector)

	// OperatorId is set in contract during registration so we get it after registering operator.
	operatorId, err := sdkClients.AvsRegistryChainReader.GetOperatorId(&bind.CallOpts{}, operatorAddress)
	if err != nil {
		logger.Error("Cannot get operator id", "err", err)
		return nil, err
	}

	aggregatorRpcClient, err := NewGenericAggregatorClient(cfg, avsCfg, operatorId, operatorAddress, logger, avsAndEigenMetrics)
	if err != nil {
		logger.Error("buildAggregatorClient falied", "err", err)
		return nil, err
	}

	newTaskCreatedChan := make(chan GenericRequest, 32)

	rpcServer := &GenericRpcServer{
		logger:             logger,
		serverIpPortAddr:   cfg.OperatorServerIpPortAddr,
		newTaskCreatedChan: newTaskCreatedChan,
	}

	operator := &Operator{
		config:                     cfg,
		avsCfg:                     avsCfg,
		logger:                     logger,
		metricsReg:                 reg,
		metrics:                    avsAndEigenMetrics,
		nodeApi:                    nodeApi,
		ethClient:                  ethRpcClient,
		avsWriter:                  avsWriter,
		avsReader:                  avsReader,
		eigenlayerReader:           sdkClients.ElChainReader,
		eigenlayerWriter:           sdkClients.ElChainWriter,
		blsKeypair:                 blsKeyPair,
		operatorAddr:               operatorAddress,
		aggregatorServerIpPortAddr: cfg.AggregatorServerIpPortAddress,
		aggregatorRpcClient:        aggregatorRpcClient,
		newTaskCreatedChan:         newTaskCreatedChan,
		operatorId:                 operatorId,
		rpcServer:                  rpcServer,
	}

	logger.Info("Operator info",
		"operatorId", operatorId,
		"operatorAddr", operatorAddress,
		"operatorG1Pubkey", operator.blsKeypair.GetPubKeyG1(),
		"operatorG2Pubkey", operator.blsKeypair.GetPubKeyG2(),
	)

	return operator, nil

}

func (o *Operator) Start(ctx context.Context) error {
	o.logger.Info("Start operator", "address", o.operatorAddr)
	operatorIsRegistered, err := o.avsReader.IsOperatorRegistered(&bind.CallOpts{}, o.operatorAddr)
	if err != nil {
		o.logger.Error("Error checking if operator is registered", "err", err)
		return err
	}
	if !operatorIsRegistered {
		// We bubble the error all the way up instead of using logger.Fatal because logger.Fatal prints a huge stack trace
		// that hides the actual error message. This error msg is more explicit and doesn't require showing a stack trace to the user.
		return fmt.Errorf("operator is not registered. Registering operator using the operator-cli before starting operator")
	}

	o.logger.Infof("Starting operator, Init operator to aggregator.")
	err = o.aggregatorRpcClient.InitOperatorToAggregator()
	if err != nil {
		o.logger.Errorf("Init operator to aggregator failed: %v", err)
		return err
	}
	o.logger.Infof("Init operator to aggregator succeeded.")

	nodeApiErrChan := o.StartNodeApi()
	metricsErrChan := o.StartMetrics(ctx)

	o.logger.Info("start rpc server for got alert")
	if err = o.rpcServer.StartServer(ctx); err != nil {
		o.logger.Error("Error start Rpc server", "err", err)
		return err
	}
	defer func() {
		err := o.rpcServer.Stop()
		if err != nil {
			o.logger.Error("Stop Rpc server failed", "err", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-nodeApiErrChan:
			o.logger.Fatal("Error in node api server", "err", err)
		case err := <-metricsErrChan:
			o.logger.Fatal("Error in metrics server", "err", err)
		case newTaskCreatedLog := <-o.newTaskCreatedChan:
			o.logger.Info(
				"newTaskCreatedLog",
				"avs", newTaskCreatedLog.AVSName,
				"hash", newTaskCreatedLog.SigHash.String(),
				"method", newTaskCreatedLog.Method,
			)
			o.metrics.IncNumTasksReceived()
			err := o.CreateTaskThenCommitToAggregator(newTaskCreatedLog)
			if err != nil {
				o.logger.Errorf("create task then commit to aggregator failed: %v", err)
			}
		}
	}
}

func (o *Operator) CreateTaskThenCommitToAggregator(request GenericRequest) error {
	o.logger.Info(
		"new task to sign",
		"avs", request.AVSName,
		"method", request.Method,
		"hash", request.SigHash.String(),
	)

	o.metrics.IncNumTasksReceived()

	taskCreateResponse, err := o.aggregatorRpcClient.CreateTaskToAggregator(request.SigHash, request.Method, request.ParamsRaw)
	if err != nil {
		o.logger.Error("process task failed by create task", "err", err)
		var code uint32
		if strings.Contains(err.Error(), "already finished") {
			code = 2
		} else {
			code = 1
		}

		request.SendRespose(code, err, "CreateTaskToAggregator failed")
		return errors.Wrap(err, "create task to aggregator failed")
	}

	hash, err := taskCreateResponse.SigHash()
	if err != nil {
		return errors.Wrap(err, "sign task response failed")
	}

	blsSignature := o.blsKeypair.SignMessage(hash)

	go o.aggregatorRpcClient.SendSignedTaskResponseToAggregator(
		request.Method,
		request.ParamsRaw,
		*blsSignature,
		o.operatorId,
		*taskCreateResponse,
		request.ResponseChan,
	)

	return nil
}
