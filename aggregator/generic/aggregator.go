package generic

import (
	"context"
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/core/config"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type GenericAggregator struct {
	logger                  logging.Logger
	jsonRpcServerIpPortAddr string

	services *AVSGenericServices
	rpc      *JsonRpcServer
	wg       sync.WaitGroup
}

func GenericAggregatorMain(cliCtx *cli.Context, ctx context.Context, c *config.Config) (*GenericAggregator, error) {
	avsConfigs, err := config.NewAVSConfigs(cliCtx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new avs configs from raws")
	}

	if len(avsConfigs) == 0 {
		c.Logger.Info("not start generic aggregator by not use avs configs")
		return nil, nil
	}

	aggregator, err := NewGenericAggregator(c, avsConfigs)
	if err != nil {
		return nil, errors.Wrap(err, "new generic aggregator failed")
	}

	err = aggregator.Start(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "start aggregator failed")
	}

	return aggregator, nil
}

func NewGenericAggregator(c *config.Config, awsConfig []config.GenericAVSConfig) (*GenericAggregator, error) {
	services, err := NewAVSGenericServices(c, awsConfig)
	if err != nil {
		return nil, errors.Wrap(err, "NewAVSGenericServices")
	}

	rpc := NewJsonRpcServer(c.Logger, services, c.RpcVhosts, c.RpcCors)

	return &GenericAggregator{
		logger:                  c.Logger,
		jsonRpcServerIpPortAddr: c.AggregatorJSONRPCServerIpPortAddr,

		services: services,
		rpc:      rpc,
	}, nil
}

func (g *GenericAggregator) Start(ctx context.Context) error {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		err := g.services.Start(ctx)
		if err != nil {
			g.logger.Error("start services failed", "err", err.Error())
		}
	}()

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		g.rpc.StartServer(ctx, g.jsonRpcServerIpPortAddr)
	}()

	g.logger.Info("GenericAggregator started")

	return nil
}

func (g *GenericAggregator) Wait() {
	g.wg.Wait()
}
