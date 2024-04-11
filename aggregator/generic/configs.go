package generic

import (
	"fmt"

	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func NewAVSConfigs(configRaws []config.AVSConfigRaw) ([]message.GenericAVSConfig, error) {
	res := make([]message.GenericAVSConfig, 0, len(configRaws))

	for _, raw := range configRaws {
		cfg, err := NewAVSConfig(raw)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create config by avs raw")
		}

		res = append(res, cfg)
	}

	return res, nil
}

func NewAVSConfig(configRaw config.AVSConfigRaw) (message.GenericAVSConfig, error) {
	if !common.IsHexAddress(configRaw.AVSRegistryCoordinatorAddress) {
		return message.GenericAVSConfig{}, fmt.Errorf("avs config %s 's avs_registry_coordinator_address not a hex address", configRaw.AVSName)
	}
	AVSRegistryCoordinatorAddress := common.HexToAddress(configRaw.AVSRegistryCoordinatorAddress)

	if !common.IsHexAddress(configRaw.OperatorStateRetrieverAddress) {
		return message.GenericAVSConfig{}, fmt.Errorf("avs config %s 's operator_state_retriever_address not a hex address", configRaw.AVSName)
	}
	OperatorStateRetrieverAddress := common.HexToAddress(configRaw.OperatorStateRetrieverAddress)

	if !common.IsHexAddress(configRaw.AVSContractAddress) {
		return message.GenericAVSConfig{}, fmt.Errorf("avs config %s 's contract address not a hex address", configRaw.AVSName)
	}

	AVSContractAddress := common.HexToAddress(configRaw.AVSContractAddress)

	var avsAbi abi.ABI
	err := avsAbi.UnmarshalJSON(configRaw.Abi)
	if err != nil {
		return message.GenericAVSConfig{}, errors.Wrapf(err, "avs config %s unmarshal failed", configRaw.AVSName)
	}

	return message.GenericAVSConfig{
		AVSName:                       configRaw.AVSName,
		QuorumNumbers:                 configRaw.QuorumNumbers,
		AVSRegistryCoordinatorAddress: AVSRegistryCoordinatorAddress,
		OperatorStateRetrieverAddress: OperatorStateRetrieverAddress,
		AVSContractAddress:            AVSContractAddress,
		Abi:                           avsAbi,
	}, nil

}
