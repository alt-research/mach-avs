package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type NodeConfig struct {
	// used to set the logger level (true = info, false = debug)
	Production                        bool   `yaml:"production"`
	AVSName                           string `yaml:"avs_name"`
	OperatorStateRetrieverAddress     string `yaml:"operator_state_retriever_address"`
	AVSRegistryCoordinatorAddress     string `yaml:"avs_registry_coordinator_address"`
	EthRpcUrl                         string `yaml:"eth_rpc_url"`
	EthWsUrl                          string `yaml:"eth_ws_url"`
	BlsPrivateKeyStorePath            string `yaml:"bls_private_key_store_path"`
	EcdsaPrivateKeyStorePath          string `yaml:"ecdsa_private_key_store_path"`
	OperatorEcdsaAddress              string `yaml:"operator_ecdsa_address"`
	AggregatorServerIpPortAddress     string `yaml:"aggregator_server_ip_port_address"`
	AggregatorGRPCServerIpPortAddress string `yaml:"aggregator_grpc_server_ip_port_address"`
	AggregatorJSONRPCServerIpPortAddr string `yaml:"aggregator_jsonrpc_server_ip_port_address"`
	EigenMetricsIpPortAddress         string `yaml:"eigen_metrics_ip_port_address"`
	EnableMetrics                     bool   `yaml:"enable_metrics"`
	NodeApiIpPortAddress              string `yaml:"node_api_ip_port_address"`
	EnableNodeApi                     bool   `yaml:"enable_node_api"`
	OperatorServerIpPortAddr          string `yaml:"operator_server_ip_port_addr"`
	MetadataURI                       string `yaml:"metadata_uri"`
	OperatorSocket                    string `yaml:"operator_socket"`
	Layer1ChainId                     uint32 `yaml:"layer1_chain_id"`
	Layer2ChainId                     uint32 `yaml:"layer2_chain_id"`
}

// use the env config first for some keys
func (c *NodeConfig) WithEnv() {
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

	avsName, ok := os.LookupEnv("AVS_NAME")
	if ok && avsName != "" {
		c.AVSName = avsName
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
}

type AvsConfig struct {
	AVSName                       string
	QuorumNumbers                 []uint8
	AVSRegistryCoordinatorAddress common.Address
	OperatorStateRetrieverAddress common.Address
	Abi                           abi.ABI
}
