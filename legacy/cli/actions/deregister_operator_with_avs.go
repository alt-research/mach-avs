package actions

import (
	"encoding/json"
	"log"

	"github.com/alt-research/avs/legacy/core/config"
	"github.com/alt-research/avs/legacy/operator"
	"github.com/urfave/cli"
)

func DeregisterOperatorWithAvs(ctx *cli.Context) error {
	configPath := ctx.GlobalString(config.ConfigFileFlag.Name)
	nodeConfig := config.NodeConfig{}

	if configPath != "" {
		err := config.ReadYamlConfig(configPath, &nodeConfig)
		if err != nil {
			return err
		}
		configJson, err := json.MarshalIndent(nodeConfig, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.Println("Config:", string(configJson))
	}

	operator, err := operator.NewOperatorFromConfig(nodeConfig, true)
	if err != nil {
		return err
	}

	err = operator.DeregisterOperatorWithAvs()
	if err != nil {
		return err
	}

	return nil
}
