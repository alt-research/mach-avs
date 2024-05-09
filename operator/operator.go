package operator

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/alt-research/avs/core"
	"github.com/alt-research/avs/core/alert"
	"github.com/alt-research/avs/core/chainio"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
	"github.com/alt-research/avs/metrics"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	sdkelcontracts "github.com/Layr-Labs/eigensdk-go/chainio/clients/elcontracts"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	sdkmetrics "github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/Layr-Labs/eigensdk-go/metrics/collectors/economic"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	"github.com/Layr-Labs/eigensdk-go/nodeapi"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
)

const AVS_NAME = "mach"
const SEM_VER = "0.0.1"

type Operator struct {
	config           config.NodeConfig
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
	metadataURI      string
	rpcServer        RpcServer
	// receive new tasks in this chan (typically from mach service)
	newTaskCreatedChan chan alert.AlertRequest
	// ip address of aggregator
	aggregatorServerIpPortAddr string
	// rpc client to send signed task responses to aggregator
	aggregatorRpcClient AggregatorRpcClienter
	// needed when opting in to avs (allow this service manager contract to slash operator)
	serviceManagerAddr common.Address
}

// use the env config first for some keys
func withEnvConfig(c config.NodeConfig) config.NodeConfig {
	// This keys can use the environment:
	//
	// - `ETH_RPC_URL` : eth_rpc_url
	// - `ETH_WS_URL` : eth_ws_url
	// - `ECDSA_PRIVATE_KEY_PATH` : ecdsa_private_key_store_path
	// - `BLS_PRIVATE_KEY_PATH` : bls_private_key_store_path
	// - `AGGREGATOR_SERVER_URL` : eth_rpc_url
	// - `EIGEN_METRICS_URL` : eigen_metrics_ip_port_address
	// - `NODE_API_URL` : node_api_ip_port_address
	// - `ENABLE_METRICS` : enable_metrics
	// - `ENABLE_NODE_API` : enable_node_api
	// - `AVS_REGISTRY_COORDINATOR_ADDRESS` : avs_registry_coordinator_address
	// - `OPERATOR_STATE_RETRIEVER_ADDRESS` : operator_state_retriever_address
	// - `OPERATOR_SERVER_URL` : operator_server_ip_port_addr
	// - `METADATA_URI` : metadata_uri

	Production, ok := os.LookupEnv("OPERATOR_PRODUCTION")
	if ok && Production != "" {
		c.Production = Production == "true"
	}

	ethRpcUrl, ok := os.LookupEnv("ETH_RPC_URL")
	if ok && ethRpcUrl != "" {
		c.EthRpcUrl = ethRpcUrl
	}

	EthWsUrl, ok := os.LookupEnv("ETH_WS_URL")
	if ok && EthWsUrl != "" {
		c.EthWsUrl = EthWsUrl
	}

	ecdsaPrivateKeyStorePath, ok := os.LookupEnv("ECDSA_PRIVATE_KEY_PATH")
	if ok && ecdsaPrivateKeyStorePath != "" {
		c.EcdsaPrivateKeyStorePath = ecdsaPrivateKeyStorePath
	}

	blsPrivateKeyStorePath, ok := os.LookupEnv("BLS_PRIVATE_KEY_PATH")
	if ok && blsPrivateKeyStorePath != "" {
		c.BlsPrivateKeyStorePath = blsPrivateKeyStorePath
	}

	aggregatorServerIpPortAddress, ok := os.LookupEnv("AGGREGATOR_SERVER_URL")
	if ok && aggregatorServerIpPortAddress != "" {
		c.AggregatorServerIpPortAddress = aggregatorServerIpPortAddress
	}

	grpcAggregatorServerIpPortAddress, ok := os.LookupEnv("AGGREGATOR_GRPC_SERVER_URL")
	if ok && grpcAggregatorServerIpPortAddress != "" {
		c.AggregatorGRPCServerIpPortAddress = grpcAggregatorServerIpPortAddress
	}

	jsonRPCAggregatorServerIpPortAddress, ok := os.LookupEnv("AGGREGATOR_JSONRPC_SERVER_URL")
	if ok && jsonRPCAggregatorServerIpPortAddress != "" {
		c.AggregatorJSONRPCServerIpPortAddr = jsonRPCAggregatorServerIpPortAddress
	}

	eigenMetricsIpPortAddress, ok := os.LookupEnv("EIGEN_METRICS_URL")
	if ok && eigenMetricsIpPortAddress != "" {
		c.EigenMetricsIpPortAddress = eigenMetricsIpPortAddress
	}

	nodeApiIpPortAddress, ok := os.LookupEnv("NODE_API_URL")
	if ok && nodeApiIpPortAddress != "" {
		c.NodeApiIpPortAddress = nodeApiIpPortAddress
	}

	enableMetrics, ok := os.LookupEnv("ENABLE_METRICS")
	if ok && enableMetrics != "" {
		c.EnableMetrics = enableMetrics == "true"
	}

	enableNodeApi, ok := os.LookupEnv("ENABLE_NODE_API")
	if ok && enableNodeApi != "" {
		c.EnableNodeApi = enableNodeApi == "true"
	}

	aVSRegistryCoordinatorAddress, ok := os.LookupEnv("AVS_REGISTRY_COORDINATOR_ADDRESS")
	if ok && aVSRegistryCoordinatorAddress != "" {
		c.AVSRegistryCoordinatorAddress = aVSRegistryCoordinatorAddress
	}

	operatorStateRetrieverAddress, ok := os.LookupEnv("OPERATOR_STATE_RETRIEVER_ADDRESS")
	if ok && operatorStateRetrieverAddress != "" {
		c.OperatorStateRetrieverAddress = operatorStateRetrieverAddress
	}

	operatorServerIpPortAddr, ok := os.LookupEnv("OPERATOR_SERVER_URL")
	if ok && operatorServerIpPortAddr != "" {
		c.OperatorServerIpPortAddr = operatorServerIpPortAddr
	}

	metadataURI, ok := os.LookupEnv("METADATA_URI")
	if ok && metadataURI != "" {
		c.MetadataURI = metadataURI
	}

	operatorSocket, ok := os.LookupEnv("OPERATOR_SOCKET")
	if ok && operatorSocket != "" {
		c.OperatorSocket = operatorSocket
	}

	layer1ChainId, ok := os.LookupEnv("LAYER1_CHAIN_ID")
	if ok && layer1ChainId != "" {
		layer1ChainId, err := strconv.Atoi(layer1ChainId)
		if err != nil {
			panic(fmt.Sprintf("layer1_chain_id parse error: %v", err))
		}

		c.Layer1ChainId = uint32(layer1ChainId)
	}

	layer2ChainId, ok := os.LookupEnv("LAYER2_CHAIN_ID")
	if ok && layer2ChainId != "" {
		layer2ChainId, err := strconv.Atoi(layer2ChainId)
		if err != nil {
			panic(fmt.Sprintf("layer2_chain_id parse error: %v", err))
		}

		c.Layer2ChainId = uint32(layer2ChainId)
	}

	operatorEcdsaAddress, ok := os.LookupEnv("OPERATOR_ECDSA_ADDRESS")
	if ok && operatorEcdsaAddress != "" {
		c.OperatorEcdsaAddress = operatorEcdsaAddress
	}

	configJson, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}

	log.Println("Config Env:", string(configJson))

	return c
}

// TODO(samlaf): config is a mess right now, since the chainio client constructors
//
//	take the config in core (which is shared with aggregator and challenger)
func NewOperatorFromConfig(cfg config.NodeConfig, isUseEcdsaKey bool) (*Operator, error) {
	var logLevel sdklogging.LogLevel
	if cfg.Production {
		logLevel = sdklogging.Production
	} else {
		logLevel = sdklogging.Development
	}
	logger, err := core.NewZapLogger(logLevel)
	if err != nil {
		return nil, err
	}

	c := withEnvConfig(cfg)

	reg := prometheus.NewRegistry()
	eigenMetrics := sdkmetrics.NewEigenMetrics(AVS_NAME, c.EigenMetricsIpPortAddress, reg, logger)
	avsAndEigenMetrics := metrics.NewAvsAndEigenMetrics(AVS_NAME, eigenMetrics, reg)

	// Setup Node Api
	nodeApi := nodeapi.NewNodeApi(AVS_NAME, SEM_VER, c.NodeApiIpPortAddress, logger)

	var ethRpcClient eth.Client
	if c.EnableMetrics {
		rpcCallsCollector := rpccalls.NewCollector(AVS_NAME, reg)
		ethRpcClient, err = eth.NewInstrumentedClient(c.EthRpcUrl, rpcCallsCollector)
		if err != nil {
			logger.Errorf("Cannot create http ethclient", "err", err)
			return nil, err
		}
	} else {
		ethRpcClient, err = eth.NewClient(c.EthRpcUrl)
		if err != nil {
			logger.Errorf("Cannot create http ethclient", "err", err)
			return nil, err
		}
	}

	blsKeyPassword, ok := os.LookupEnv("OPERATOR_BLS_KEY_PASSWORD")
	if !ok {
		logger.Warnf("OPERATOR_BLS_KEY_PASSWORD env var not set. using empty string")
	}
	blsKeyPair, err := bls.ReadPrivateKeyFromFile(c.BlsPrivateKeyStorePath, blsKeyPassword)
	if err != nil {
		logger.Errorf("Cannot parse bls private key", "err", err)
		return nil, err
	}
	// TODO(samlaf): should we add the chainId to the config instead?
	// this way we can prevent creating a signer that signs on mainnet by mistake
	// if the config says chainId=5, then we can only create a goerli signer
	chainId, err := ethRpcClient.ChainID(context.Background())
	if err != nil {
		logger.Error("Cannot get chainId", "err", err)
		return nil, err
	}

	var operatorAddress common.Address
	var avsWriter *chainio.AvsWriter
	var privateKey *ecdsa.PrivateKey
	if isUseEcdsaKey {
		var err error
		operatorAddress, err = sdkEcdsa.GetAddressFromKeyStoreFile(c.EcdsaPrivateKeyStorePath)
		if err != nil {
			panic(err)
		}

		ecdsaKeyPassword, ok := os.LookupEnv("OPERATOR_ECDSA_KEY_PASSWORD")
		if !ok {
			logger.Warnf("OPERATOR_ECDSA_KEY_PASSWORD env var not set. using empty string")
		}

		signerConfig := signerv2.Config{
			KeystorePath: c.EcdsaPrivateKeyStorePath,
			Password:     ecdsaKeyPassword,
		}
		signerV2, _, err := signerv2.SignerFromConfig(signerConfig, chainId)
		if err != nil {
			panic(err)
		}

		privateKey, err = sdkEcdsa.ReadKey(signerConfig.KeystorePath, signerConfig.Password)
		if err != nil {
			return nil, err
		}

		txSender, err := wallet.NewPrivateKeyWallet(ethRpcClient, signerV2, operatorAddress, logger)
		if err != nil {
			return nil, err
		}
		txMgr := txmgr.NewSimpleTxManager(txSender, ethRpcClient, logger, operatorAddress)

		avsWriter, err = chainio.BuildAvsWriter(
			txMgr, common.HexToAddress(c.AVSRegistryCoordinatorAddress),
			common.HexToAddress(c.OperatorStateRetrieverAddress), ethRpcClient, logger,
		)
		if err != nil {
			logger.Error("Cannot create AvsWriter", "err", err)
			return nil, err
		}
	} else {
		// just use a default value, it just to avoid build client panic
		// will never use this value to generate tx.
		ecdsaPrivateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
		if err != nil {
			panic(err)
		}
		privateKey = ecdsaPrivateKey

		if c.EcdsaPrivateKeyStorePath != "" {
			operatorAddress, err = sdkEcdsa.GetAddressFromKeyStoreFile(c.EcdsaPrivateKeyStorePath)
			if err != nil {
				panic(err)
			}
		} else {
			if c.OperatorEcdsaAddress == "" {
				return nil, fmt.Errorf("If not use EcdsaPrivateKeyStorePath, must use operator_ecdsa_address or `OPERATOR_ECDSA_ADDRESS` env to select ecdsa address!")
			}

			if !common.IsHexAddress(c.OperatorEcdsaAddress) {
				return nil, fmt.Errorf("the operator_ecdsa_address format is not hex address!")
			}

			operatorAddress = common.HexToAddress(c.OperatorEcdsaAddress)
		}
	}

	chainioConfig := clients.BuildAllConfig{
		EthHttpUrl:                 c.EthRpcUrl,
		EthWsUrl:                   c.EthWsUrl,
		RegistryCoordinatorAddr:    c.AVSRegistryCoordinatorAddress,
		OperatorStateRetrieverAddr: c.OperatorStateRetrieverAddress,
		AvsName:                    AVS_NAME,
		PromMetricsIpPortAddress:   c.EigenMetricsIpPortAddress,
	}
	sdkClients, err := clients.BuildAll(chainioConfig, privateKey, logger)
	if err != nil {
		panic(err)
	}

	avsReader, err := chainio.BuildAvsReader(
		common.HexToAddress(c.AVSRegistryCoordinatorAddress),
		common.HexToAddress(c.OperatorStateRetrieverAddress),
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
		AVS_NAME, logger, operatorAddress, quorumNames)
	reg.MustRegister(economicMetricsCollector)

	// OperatorId is set in contract during registration so we get it after registering operator.
	operatorId, err := sdkClients.AvsRegistryChainReader.GetOperatorId(&bind.CallOpts{}, operatorAddress)
	if err != nil {
		logger.Error("Cannot get operator id", "err", err)
		return nil, err
	}

	aggregatorRpcClient, err := buildAggregatorClient(c, operatorId, operatorAddress, logger, avsAndEigenMetrics)
	if err != nil {
		logger.Error("buildAggregatorClient falied", "err", err)
		return nil, err
	}

	newTaskCreatedChan := make(chan alert.AlertRequest, 32)
	rpcServer := NewRpcServer(logger, c.OperatorServerIpPortAddr, newTaskCreatedChan)

	operator := &Operator{
		config:                     c,
		logger:                     logger,
		metricsReg:                 reg,
		metrics:                    avsAndEigenMetrics,
		nodeApi:                    nodeApi,
		ethClient:                  ethRpcClient,
		avsWriter:                  avsWriter,
		avsReader:                  avsReader,
		eigenlayerReader:           sdkClients.ElChainReader,
		eigenlayerWriter:           sdkClients.ElChainWriter,
		rpcServer:                  rpcServer,
		blsKeypair:                 blsKeyPair,
		operatorAddr:               operatorAddress,
		aggregatorServerIpPortAddr: c.AggregatorServerIpPortAddress,
		aggregatorRpcClient:        aggregatorRpcClient,
		newTaskCreatedChan:         newTaskCreatedChan,
		serviceManagerAddr:         common.HexToAddress(c.AVSRegistryCoordinatorAddress),
		metadataURI:                c.MetadataURI,
		operatorId:                 operatorId,
	}

	logger.Info("Operator info",
		"operatorId", operatorId,
		"operatorAddr", operatorAddress,
		"operatorG1Pubkey", operator.blsKeypair.GetPubKeyG1(),
		"operatorG2Pubkey", operator.blsKeypair.GetPubKeyG2(),
	)

	return operator, nil

}

func buildAggregatorClient(c config.NodeConfig, operatorId sdktypes.OperatorId, operatorAddr common.Address, logger sdklogging.Logger, metrics metrics.Metrics) (AggregatorRpcClienter, error) {
	if c.AggregatorJSONRPCServerIpPortAddr != "" {
		logger.Info("Use json rpc server to connect to the aggregator", "address", c.AggregatorJSONRPCServerIpPortAddr)
		cli, err := NewAggregatorJsonRpcClient(c, operatorId, operatorAddr, logger, metrics)
		if err != nil {
			logger.Error("Cannot create AggregatorGRpcClient. Is aggregator running?", "err", err)
			return nil, err
		}

		return cli, nil
	}

	if c.AggregatorGRPCServerIpPortAddress != "" {
		logger.Info("Use grpc server to connect to the aggregator", "address", c.AggregatorGRPCServerIpPortAddress)

		cli, err := NewAggregatorGRpcClient(c, operatorId, operatorAddr, logger, metrics)
		if err != nil {
			logger.Error("Cannot create AggregatorGRpcClient. Is aggregator running?", "err", err)
			return nil, err
		}

		return cli, nil
	} else {
		logger.Info("Use legacy rpc server to connect to the aggregator", "address", c.AggregatorServerIpPortAddress)

		cli, err := NewAggregatorRpcClient(c, operatorId, operatorAddr, logger, metrics)
		if err != nil {
			logger.Error("Cannot create AggregatorRpcClient. Is aggregator running?", "err", err)
			return nil, err
		}

		return cli, nil
	}
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

	o.logger.Infof("Starting operator.")

	if o.config.EnableNodeApi {
		o.nodeApi.RegisterNewService(
			ServiceOperator,
			ServiceOperator,
			"operator to commit alert to aggregator",
			nodeapi.ServiceStatusInitializing,
		)

		o.nodeApi.RegisterNewService(
			ServiceOperatorAggregator,
			ServiceOperatorAggregator,
			"operator aggregator work",
			nodeapi.ServiceStatusInitializing,
		)

		o.nodeApi.RegisterNewService(
			ServiceOperatorVerifier,
			ServiceOperatorVerifier,
			"operator verifier work",
			nodeapi.ServiceStatusInitializing,
		)

		o.nodeApi.UpdateHealth(nodeapi.Healthy)
		o.nodeApi.Start()
	}

	o.logger.Infof("Init operator to aggregator.")
	err = o.aggregatorRpcClient.InitOperatorToAggregator()
	if err != nil {
		o.logger.Errorf("Init operator to aggregator failed: %v", err)
		return err
	}
	o.logger.Infof("Init operator to aggregator succeeded.")

	if o.config.EnableNodeApi {
		o.nodeApi.UpdateServiceStatus(
			ServiceOperator,
			nodeapi.ServiceStatusUp,
		)

		o.nodeApi.UpdateServiceStatus(
			ServiceOperatorAggregator,
			nodeapi.ServiceStatusUp,
		)

		// TODO: check the verifier working
		o.nodeApi.UpdateServiceStatus(
			ServiceOperatorVerifier,
			nodeapi.ServiceStatusUp,
		)
	}

	var metricsErrChan <-chan error
	if o.config.EnableMetrics {
		metricsErrChan = o.metrics.Start(ctx, o.metricsReg)
	} else {
		metricsErrChan = make(chan error, 1)
	}

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
		case err := <-metricsErrChan:
			// TODO(samlaf); we should also register the service as unhealthy in the node api
			// https://eigen.nethermind.io/docs/spec/api/
			o.logger.Fatal("Error in metrics server", "err", err)
		case newTaskCreatedLog := <-o.newTaskCreatedChan:
			o.logger.Info("newTaskCreatedLog", "new", newTaskCreatedLog.Alert)
			o.metrics.IncNumTasksReceived()
			taskResponse, err := o.ProcessNewTaskCreatedLog(newTaskCreatedLog.Alert)
			if err != nil {
				o.logger.Error("newTaskCreatedLog failed by new", "err", err)
				var code uint32
				if strings.Contains(err.Error(), "already finished") {
					code = 2
				}
				newTaskCreatedLog.ResChan <- alert.AlertResponse{
					Code: code,
					Err:  err,
					Msg:  "ProcessNewTaskCreatedLog failed",
				}
				continue
			}

			signedTaskResponse, err := o.SignTaskResponse(taskResponse)
			if err != nil {
				o.logger.Error("newTaskCreatedLog failed by sign task", "err", err)
				newTaskCreatedLog.ResChan <- alert.AlertResponse{
					Err: err,
					Msg: "SignTaskResponse failed",
				}
				continue
			}
			go o.aggregatorRpcClient.SendSignedTaskResponseToAggregator(signedTaskResponse, newTaskCreatedLog.ResChan)
		}
	}
}

// Takes a NewTaskCreatedLog struct as input and returns a TaskResponseHeader struct.
// The TaskResponseHeader struct is the struct that is signed and sent to the contract as a task response.
func (o *Operator) ProcessNewTaskCreatedLog(newAlert alert.Alert) (*message.AlertTaskInfo, error) {
	alertHash := newAlert.MessageHash()

	o.logger.Debug("Received new task", "task", newAlert)
	o.logger.Info("Received new task",
		"alert", alertHash,
	)

	return o.aggregatorRpcClient.CreateAlertTaskToAggregator(alertHash)
}

func (o *Operator) SignTaskResponse(taskResponse *message.AlertTaskInfo) (*message.SignedTaskRespRequest, error) {
	hash, err := taskResponse.SignHash()
	if err != nil {
		return nil, err
	}

	blsSignature := o.blsKeypair.SignMessage(hash)
	signedTaskResponse := &message.SignedTaskRespRequest{
		Alert:        *taskResponse,
		BlsSignature: *blsSignature,
		OperatorId:   o.operatorId,
	}
	o.logger.Debug("Signed task response", "signedTaskResponse", signedTaskResponse)
	return signedTaskResponse, nil
}

func (o Operator) Config() config.NodeConfig {
	return o.config
}
