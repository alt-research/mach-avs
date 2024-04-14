package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// Register a generic task
type GenericAVSConfig struct {
	AVSName                       string
	QuorumNumbers                 []uint8
	AVSRegistryCoordinatorAddress common.Address
	OperatorStateRetrieverAddress common.Address
	AVSContractAddress            common.Address
	Abi                           abi.ABI
}

func NewAVSConfigs(ctx *cli.Context) ([]GenericAVSConfig, error) {
	raws, err := newAVSConfigRaws(ctx)
	if err != nil {
		return nil, err
	}

	return NewAVSConfigsFromRaw(raws)
}

func NewAVSConfigsFromRaw(configRaws []avsConfigRaw) ([]GenericAVSConfig, error) {
	res := make([]GenericAVSConfig, 0, len(configRaws))

	for _, raw := range configRaws {
		cfg, err := NewAVSConfig(raw)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create config by avs raw")
		}

		res = append(res, cfg)
	}

	return res, nil
}

func NewAVSConfig(configRaw avsConfigRaw) (GenericAVSConfig, error) {
	if !common.IsHexAddress(configRaw.AVSRegistryCoordinatorAddress) {
		return GenericAVSConfig{}, fmt.Errorf("avs config %s 's avs_registry_coordinator_address not a hex address", configRaw.AVSName)
	}
	AVSRegistryCoordinatorAddress := common.HexToAddress(configRaw.AVSRegistryCoordinatorAddress)

	if !common.IsHexAddress(configRaw.OperatorStateRetrieverAddress) {
		return GenericAVSConfig{}, fmt.Errorf("avs config %s 's operator_state_retriever_address not a hex address", configRaw.AVSName)
	}
	OperatorStateRetrieverAddress := common.HexToAddress(configRaw.OperatorStateRetrieverAddress)

	if !common.IsHexAddress(configRaw.AVSContractAddress) {
		return GenericAVSConfig{}, fmt.Errorf("avs config %s 's contract address not a hex address", configRaw.AVSName)
	}

	AVSContractAddress := common.HexToAddress(configRaw.AVSContractAddress)

	var avsAbi abi.ABI
	err := avsAbi.UnmarshalJSON(configRaw.Abi)
	if err != nil {
		return GenericAVSConfig{}, errors.Wrapf(err, "avs config %s unmarshal failed", configRaw.AVSName)
	}

	return GenericAVSConfig{
		AVSName:                       configRaw.AVSName,
		QuorumNumbers:                 configRaw.QuorumNumbers,
		AVSRegistryCoordinatorAddress: AVSRegistryCoordinatorAddress,
		OperatorStateRetrieverAddress: OperatorStateRetrieverAddress,
		AVSContractAddress:            AVSContractAddress,
		Abi:                           avsAbi,
	}, nil

}
