package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"

	"github.com/alt-research/avs-generic-aggregator/core"
	"github.com/alt-research/avs-generic-aggregator/core/config"
	proxyUtils "github.com/alt-research/avs-generic-aggregator/proxy/utils"

	genericproxy "github.com/alt-research/avs/generic-operator-proxy"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{config.ConfigFileFlag, config.AVSConfigFlag}
	app.Name = "mach-operator-alerter-proxy"
	app.Usage = "Mach Operator"
	app.Description = "Service that handle the alert, and sends them to a generic operator."

	app.Action = operatorProxyMain
	err := app.Run(os.Args)
	if err != nil {
		log.Panicln("Application failed. Message:", err)
	}
}

func operatorProxyMain(ctx *cli.Context) error {
	log.Println("Initializing Operator Alerter Proxy")

	proxyCfg := proxyUtils.ProxyConfig{}
	if err := config.ReadConfig(ctx, &proxyCfg); err != nil {
		panic(err)
	}
	proxyCfg.WithEnv()

	nodeConfig := genericproxy.MachProxyConfig{}
	if err := config.ReadConfig(ctx, &nodeConfig); err != nil {
		panic(err)
	}
	nodeConfig.ProxyConfig = proxyCfg

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, err := core.NewLogger(nodeConfig.Production)
	if err != nil {
		return errors.Wrap(err, "New logger")
	}

	avsConfigs, err := config.NewAVSConfigs(ctx)
	if err != nil {
		return err
	}

	ethRpcClient, err := eth.NewClient(nodeConfig.EthRpcUrl)
	if err != nil {
		logger.Errorf("Cannot create http ethclient", "err", err)
		return err
	}

	rpcServer := genericproxy.NewAlertProxyRpcServer(
		logger,
		ethRpcClient,
		avsConfigs,
		nodeConfig.AVSName,
		nodeConfig.Method,
		nodeConfig.GenericOperatorAddr,
		nodeConfig.RpcCfg,
		nodeConfig.ChainIds,
	)

	return rpcServer.Start(mainCtx)
}
