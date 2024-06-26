package actions

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	sdkecdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	sdkutils "github.com/Layr-Labs/eigensdk-go/utils"
	"github.com/alt-research/avs/legacy/core/config"
	"github.com/alt-research/avs/legacy/operator"
	"github.com/pkg/errors"
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

	quorumNumbers := []byte{0}

	quorumNumbersEnvStr, ok := os.LookupEnv("REG_QUORUM_NUMBERS")
	if ok && quorumNumbersEnvStr != "" {
		log.Printf("use REG_QUORUM_NUMBERS %v", quorumNumbersEnvStr)
		quorumNumbers, err = parseQuorumNumbers(quorumNumbersEnvStr)
		if err != nil {
			return errors.Wrap(err, "failed to parse REG_QUORUM_NUMBERS into quorumNumbers")
		}
	}

	quorumNumbersStr := ctx.String("quorum-numbers")
	if quorumNumbersStr != "" {
		log.Printf("use --quorum-numbers %v", quorumNumbersEnvStr)
		quorumNumbers, err = parseQuorumNumbers(quorumNumbersEnvStr)
		if err != nil {
			return errors.Wrap(err, "failed to parse --quorum-numbers into quorumNumbers")
		}
	}

	err = operator.RegisterOperatorWithAvs(operatorEcdsaPrivKey, quorumNumbers)
	if err != nil {
		return err
	}

	return nil
}

func parseQuorumNumbers(str string) ([]byte, error) {
	strArray := strings.Split(str, ",")

	res := make([]byte, 0, len(strArray))
	for i, s := range strArray {
		si, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to atoi %s in pos %v", s, i)
		}

		if si < 0 || si >= 192 {
			return nil, errors.Errorf("failed by %v not in [0, 192)", si)
		}

		res = append(res, byte(si))
	}

	return res, nil
}
