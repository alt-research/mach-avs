package genericproxy

import proxyUtils "github.com/alt-research/avs-generic-aggregator/proxy/utils"

type MachProxyConfig struct {
	proxyUtils.ProxyConfig

	ChainIds map[string]uint32 `yaml:"chain_ids"`
}
