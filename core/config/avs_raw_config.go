package config

import (
	"encoding/json"

	sdkutils "github.com/Layr-Labs/eigensdk-go/utils"
	"github.com/urfave/cli"
)

type avsConfigRaw struct {
	AVSName                       string          `json:"avs_name"`
	QuorumNumbers                 []uint8         `json:"quorum_numbers"`
	AVSRegistryCoordinatorAddress string          `json:"avs_registry_coordinator_address"`
	OperatorStateRetrieverAddress string          `json:"operator_state_retriever_address"`
	AVSContractAddress            string          `json:"avs_contract_address"`
	Abi                           json.RawMessage `json:"abi"`
}

func newAVSConfigRaws(ctx *cli.Context) ([]avsConfigRaw, error) {
	var configRaw []avsConfigRaw

	configFilePath := ctx.GlobalString(AVSConfigFlag.Name)
	if configFilePath != "" {
		err := sdkutils.ReadJsonConfig(configFilePath, &configRaw)
		if err != nil {
			return nil, err
		}
	}

	return configRaw, nil
}
