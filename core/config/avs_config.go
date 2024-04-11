package config

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
