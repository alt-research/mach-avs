package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"

	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/operator"
	generic_operator "github.com/alt-research/avs/operator/generic"

	sdkutils "github.com/Layr-Labs/eigensdk-go/utils"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{config.ConfigFileFlag, config.AVSConfigFlag}
	app.Name = "mach-operator"
	app.Usage = "Mach Operator"
	app.Description = "Service that handle the alert, and sends them to the aggregator."

	app.Action = operatorMain
	err := app.Run(os.Args)
	if err != nil {
		log.Panicln("Application failed. Message:", err)
	}
}

func operatorMain(ctx *cli.Context) error {

	log.Println("Initializing Operator")
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
		log.Println("Config from file:", string(configJson))
	}

	avsConfig := ctx.GlobalString(config.AVSConfigFlag.Name)
	if avsConfig != "" {
		log.Println("init generic avs config:", avsConfig)
		mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		return generic_operator.GenericOperatorMain(ctx, mainCtx, nodeConfig)

	} else {
		log.Println("initializing operator")
		operator, err := operator.NewOperatorFromConfig(nodeConfig, false)
		if err != nil {
			return err
		}
		log.Println("initialized operator")

		log.Println("starting operator")
		err = operator.Start(context.Background())
		if err != nil {
			return err
		}
		log.Println("started operator")
	}

	return nil

}
