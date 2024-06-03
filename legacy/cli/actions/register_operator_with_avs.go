package actions

import (
	"encoding/json"
	"log"
	"os"

	sdkecdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	sdkutils "github.com/Layr-Labs/eigensdk-go/utils"
	"github.com/alt-research/avs/legacy/core/config"
	"github.com/alt-research/avs/legacy/operator"
	"github.com/urfave/cli"
)

func RegisterOperatorWithAvs(ctx *cli.Context) error {

	configPath := ctx.GlobalString(config.ConfigFileFlag.Name)
	nodeConfig := config.NodeConfig{}

	if configPath != "" {
		err := sdkutils.ReadYamlConfig(configPath, &nodeConfig)
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

	ecdsaKeyPassword, ok := os.LookupEnv("OPERATOR_ECDSA_KEY_PASSWORD")
	if !ok {
		log.Printf("OPERATOR_ECDSA_KEY_PASSWORD env var not set. using empty string")
	}
	operatorEcdsaPrivKey, err := sdkecdsa.ReadKey(
		operator.Config().EcdsaPrivateKeyStorePath,
		ecdsaKeyPassword,
	)
	if err != nil {
		return err
	}

	err = operator.RegisterOperatorWithAvs(operatorEcdsaPrivKey)
	if err != nil {
		return err
	}

	return nil
}
