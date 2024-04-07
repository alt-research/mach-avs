package config

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"

	"github.com/alt-research/avs/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	"github.com/Layr-Labs/eigensdk-go/types"

	sdkutils "github.com/Layr-Labs/eigensdk-go/utils"
)

// Config contains all of the configuration information for a mach aggregators and challengers.
// Operators use a separate config. (see config-files/operator.anvil.yaml)
type Config struct {
	BlsPrivateKey             *bls.PrivateKey
	Logger                    sdklogging.Logger
	EigenMetricsIpPortAddress string
	// we need the url for the eigensdk currently... eventually standardize api so as to
	// only take an ethclient or an rpcUrl (and build the ethclient at each constructor site)
	EthHttpRpcUrl                     string
	EthWsRpcUrl                       string
	EthHttpClient                     eth.Client
	EthWsClient                       eth.Client
	OperatorStateRetrieverAddr        common.Address
	RegistryCoordinatorAddr           common.Address
	AggregatorServerIpPortAddr        string
	AggregatorGRPCServerIpPortAddr    string
	AggregatorJSONRPCServerIpPortAddr string
	Layer1ChainId                     uint32
	Layer2ChainId                     uint32
	QuorumNums                        types.QuorumNums
	// json:"-" skips this field when marshaling (only used for logging to stdout), since SignerFn doesnt implement marshalJson
	SignerFn          signerv2.SignerFn `json:"-"`
	PrivateKey        *ecdsa.PrivateKey `json:"-"`
	TxMgr             txmgr.TxManager
	AggregatorAddress common.Address
}

// These are read from ConfigFileFlag
type ConfigRaw struct {
	Environment                       sdklogging.LogLevel `yaml:"environment"`
	EthRpcUrl                         string              `yaml:"eth_rpc_url"`
	EthWsUrl                          string              `yaml:"eth_ws_url"`
	AggregatorServerIpPortAddr        string              `yaml:"aggregator_server_ip_port_address"`
	AggregatorGRPCServerIpPortAddr    string              `yaml:"aggregator_grpc_server_ip_port_address"`
	AggregatorJSONRPCServerIpPortAddr string              `yaml:"aggregator_jsonrpc_server_ip_port_address"`
	Layer1ChainId                     uint32              `yaml:"layer1_chain_id"`
	Layer2ChainId                     uint32              `yaml:"layer2_chain_id"`
	QuorumNums                        []uint8             `yaml:"quorum_nums"`
}

// These are read from DeploymentFileFlag
type MachAvsDeploymentRaw struct {
	RegistryCoordinatorAddr    string `json:"registryCoordinator"`
	OperatorStateRetrieverAddr string `json:"operatorStateRetriever"`
}

// NewConfig parses config file to read from from flags or environment variables
// Note: This config is shared by challenger and aggregator and so we put in the core.
// Operator has a different config and is meant to be used by the operator CLI.
func NewConfig(ctx *cli.Context) (*Config, error) {

	var configRaw ConfigRaw
	configFilePath := ctx.GlobalString(ConfigFileFlag.Name)
	if configFilePath != "" {
		err := sdkutils.ReadYamlConfig(configFilePath, &configRaw)
		if err != nil {
			return nil, err
		}
	}

	ethRpcUrl, ok := os.LookupEnv("ETH_RPC_URL")
	if ok && ethRpcUrl != "" {
		configRaw.EthRpcUrl = ethRpcUrl
	}

	EthWsUrl, ok := os.LookupEnv("ETH_WS_URL")
	if ok && EthWsUrl != "" {
		configRaw.EthWsUrl = EthWsUrl
	}

	aggregatorServerIpPortAddress, ok := os.LookupEnv("AGGREGATOR_SERVER_URL")
	if ok && aggregatorServerIpPortAddress != "" {
		configRaw.AggregatorServerIpPortAddr = aggregatorServerIpPortAddress
	}

	aggregatorGRPCServerIpPortAddress, ok := os.LookupEnv("AGGREGATOR_GRPC_SERVER_URL")
	if ok && aggregatorGRPCServerIpPortAddress != "" {
		configRaw.AggregatorGRPCServerIpPortAddr = aggregatorGRPCServerIpPortAddress
	}

	aggregatorJSONRPCServerIpPortAddr, ok := os.LookupEnv("AGGREGATOR_JSONRPC_SERVER_URL")
	if ok && aggregatorJSONRPCServerIpPortAddr != "" {
		configRaw.AggregatorJSONRPCServerIpPortAddr = aggregatorJSONRPCServerIpPortAddr
	}

	var deploymentRaw MachAvsDeploymentRaw

	avsRegistryCoordinatorAddress, rcOk := os.LookupEnv("AVS_REGISTRY_COORDINATOR_ADDRESS")
	operatorStateRetrieverAddress, osOk := os.LookupEnv("OPERATOR_STATE_RETRIEVER_ADDRESS")

	if rcOk && osOk && avsRegistryCoordinatorAddress != "" && operatorStateRetrieverAddress != "" {
		deploymentRaw.OperatorStateRetrieverAddr = operatorStateRetrieverAddress
		deploymentRaw.RegistryCoordinatorAddr = avsRegistryCoordinatorAddress
	} else {
		deploymentFilePath := ctx.GlobalString(DeploymentFileFlag.Name)
		if deploymentFilePath == "" {
			panic("If not use env `AVS_REGISTRY_COORDINATOR_ADDRESS` and `OPERATOR_STATE_RETRIEVER_ADDRESS`, should use --avs-deployment to use config for avs contract addresses!")
		}

		if _, err := os.Stat(deploymentFilePath); errors.Is(err, os.ErrNotExist) {
			panic("Path " + deploymentFilePath + " does not exist")
		}
		if err := sdkutils.ReadJsonConfig(deploymentFilePath, &deploymentRaw); err != nil {
			panic(err)
		}
	}

	logger, err := core.NewZapLogger(configRaw.Environment)
	if err != nil {
		return nil, err
	}

	ethRpcClient, err := eth.NewClient(configRaw.EthRpcUrl)
	if err != nil {
		logger.Errorf("Cannot create http ethclient", "err", err)
		return nil, err
	}

	layer1ChainIdFromRpc, err := ethRpcClient.ChainID(context.Background())
	if err != nil {
		logger.Errorf("Cannot got chain id from eth rpc client", "err", err)
		return nil, err
	}

	if layer1ChainIdFromRpc.Uint64() != uint64(configRaw.Layer1ChainId) {
		logger.Errorf("The layer1 chain id not expect", "layer1 rpc", layer1ChainIdFromRpc, "config", configRaw.Layer1ChainId)
		return nil, fmt.Errorf("layer1 chain id not expect")
	}

	ethWsClient, err := eth.NewClient(configRaw.EthWsUrl)
	if err != nil {
		logger.Errorf("Cannot create ws ethclient", "err", err)
		return nil, err
	}

	ecdsaPrivateKeyString := ctx.GlobalString(EcdsaPrivateKeyFlag.Name)
	if ecdsaPrivateKeyString[:2] == "0x" {
		ecdsaPrivateKeyString = ecdsaPrivateKeyString[2:]
	}
	ecdsaPrivateKey, err := crypto.HexToECDSA(ecdsaPrivateKeyString)
	if err != nil {
		logger.Errorf("Cannot parse ecdsa private key", "err", err)
		return nil, err
	}

	aggregatorAddr, err := sdkutils.EcdsaPrivateKeyToAddress(ecdsaPrivateKey)
	if err != nil {
		logger.Error("Cannot get operator address", "err", err)
		return nil, err
	}

	chainId, err := ethRpcClient.ChainID(context.Background())
	if err != nil {
		logger.Error("Cannot get chainId", "err", err)
		return nil, err
	}

	signerV2, _, err := signerv2.SignerFromConfig(signerv2.Config{PrivateKey: ecdsaPrivateKey}, chainId)
	if err != nil {
		panic(err)
	}

	txSender, err := wallet.NewPrivateKeyWallet(ethRpcClient, signerV2, aggregatorAddr, logger)
	if err != nil {
		return nil, types.WrapError(errors.New("failed to create transaction sender"), err)
	}
	txMgr := txmgr.NewSimpleTxManager(txSender, ethRpcClient, logger, aggregatorAddr)

	quorumNums := make([]types.QuorumNum, len(configRaw.QuorumNums))
	for i, quorumNum := range configRaw.QuorumNums {
		quorumNums[i] = types.QuorumNum(quorumNum)
	}

	if len(quorumNums) == 0 {
		// default use zero
		logger.Warn("not quorumNums, just use [0]")
		quorumNums = []types.QuorumNum{types.QuorumNum(0)}
	}
	logger.Info(
		"the quorumNums",
		"quorumNums", fmt.Sprintf("%#v", quorumNums),
		"raw", fmt.Sprintf("%#v", configRaw.QuorumNums),
	)

	config := &Config{
		Logger:                            logger,
		EthWsRpcUrl:                       configRaw.EthWsUrl,
		EthHttpRpcUrl:                     configRaw.EthRpcUrl,
		EthHttpClient:                     ethRpcClient,
		EthWsClient:                       ethWsClient,
		OperatorStateRetrieverAddr:        common.HexToAddress(deploymentRaw.OperatorStateRetrieverAddr),
		RegistryCoordinatorAddr:           common.HexToAddress(deploymentRaw.RegistryCoordinatorAddr),
		AggregatorServerIpPortAddr:        configRaw.AggregatorServerIpPortAddr,
		AggregatorGRPCServerIpPortAddr:    configRaw.AggregatorGRPCServerIpPortAddr,
		AggregatorJSONRPCServerIpPortAddr: configRaw.AggregatorJSONRPCServerIpPortAddr,
		SignerFn:                          signerV2,
		PrivateKey:                        ecdsaPrivateKey,
		TxMgr:                             txMgr,
		AggregatorAddress:                 aggregatorAddr,
		Layer1ChainId:                     configRaw.Layer1ChainId,
		Layer2ChainId:                     configRaw.Layer2ChainId,
		QuorumNums:                        quorumNums,
	}
	config.validate()
	return config, nil
}

func (c *Config) validate() {
	// TODO: make sure every pointer is non-nil
	if c.OperatorStateRetrieverAddr == common.HexToAddress("") {
		panic("Config: BLSOperatorStateRetrieverAddr is required")
	}
	if c.RegistryCoordinatorAddr == common.HexToAddress("") {
		panic("Config: RegistryCoordinatorAddr is required")
	}
}

var (
	/* Required Flags */
	ConfigFileFlag = cli.StringFlag{
		Name:     "config",
		Required: false,
		Usage:    "Load configuration from `FILE`",
	}
	DeploymentFileFlag = cli.StringFlag{
		Name:     "avs-deployment",
		Required: false,
		Usage:    "Load avs contract addresses from `FILE`",
	}
	EcdsaPrivateKeyFlag = cli.StringFlag{
		Name:     "ecdsa-private-key",
		Usage:    "Ethereum private key",
		Required: true,
		EnvVar:   "ECDSA_PRIVATE_KEY",
	}
	/* Optional Flags */
)

var requiredFlags = []cli.Flag{
	ConfigFileFlag,
	DeploymentFileFlag,
	EcdsaPrivateKeyFlag,
}

var optionalFlags = []cli.Flag{}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag
