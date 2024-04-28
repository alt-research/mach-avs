package genericproxy

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"

	"github.com/alt-research/avs-generic-aggregator/core/config"
	"github.com/alt-research/avs-generic-aggregator/core/types"
	proxyUtils "github.com/alt-research/avs-generic-aggregator/proxy/utils"
	"github.com/alt-research/avs/core/alert"
	"github.com/alt-research/avs/operator"
)

// Hash32HeaderParam is the 1st parameter for the Mach AVS contract 's confirmAlert
type AlertHeaderParam struct {
	MessageHash                [32]byte `abi:"messageHash"`
	QuorumNumbers              []byte   `abi:"quorumNumbers"`
	QuorumThresholdPercentages []byte   `abi:"quorumThresholdPercentages"`
	ReferenceBlockNumber       uint32   `abi:"referenceBlockNumber"`
}

type CreateSigTaskForByte32Hash struct {
	Hash types.Bytes32 `json:"hash"`
}

type ProxyHashRpcServer struct {
	*proxyUtils.ProxyRpcServerBase
	logger logging.Logger
	// We use this rpc server as the handler for jsonrpc
	// to make sure same as legacy operator 's api
	rpcServer operator.RpcServer
	// receive new tasks in this chan (typically from mach service)
	newTaskCreatedChan chan alert.AlertRequest
}

func NewAlertProxyRpcServer(
	logger logging.Logger,
	ethClient eth.Client,
	avsCfg config.GenericAVSConfig,
	method string,
	genericOperatorAddr string,
	jsonrpcCfg config.JsonRpcServerConfig,
) *ProxyHashRpcServer {
	base := proxyUtils.NewProxyRpcServerBase(
		logger,
		ethClient,
		avsCfg,
		method,
		genericOperatorAddr,
		jsonrpcCfg,
	)

	newTaskCreatedChan := make(chan alert.AlertRequest, 32)
	rpcServer := operator.NewRpcServer(logger, jsonrpcCfg.Addr, newTaskCreatedChan)

	server := &ProxyHashRpcServer{
		ProxyRpcServerBase: base,
		logger:             logger,
		newTaskCreatedChan: newTaskCreatedChan,
		rpcServer:          rpcServer,
	}

	// not need do this because we not need use this server impl
	// server.SetHandler(server)

	return server
}

func (s *ProxyHashRpcServer) Start(ctx context.Context) error {
	s.logger.Info("start rpc server for got alert")
	if err := s.rpcServer.StartServer(ctx); err != nil {
		s.logger.Error("Error start Rpc server", "err", err)
		return err
	}

	defer func() {
		err := s.rpcServer.Stop()
		if err != nil {
			s.logger.Error("Stop Rpc server failed", "err", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case newTaskCreatedLog := <-s.newTaskCreatedChan:
			s.logger.Info("newTaskCreatedLog", "new", newTaskCreatedLog.Alert)
			resp, err := s.createSigTaskH32(ctx, newTaskCreatedLog.Alert.MessageHash())
			if err != nil {
				newTaskCreatedLog.ResChan <- alert.AlertResponse{
					Err: err,
					Msg: "createSigTaskH32 failed",
				}
				continue
			}

			newTaskCreatedLog.ResChan <- alert.AlertResponse{
				TxHash:    resp.TxHash,
				TaskIndex: uint32(resp.TaskIndex),
			}
		}
	}
}

func (s *ProxyHashRpcServer) createSigTaskH32(ctx context.Context, hash [32]byte) (proxyUtils.CreateSigTaskResp, error) {
	res, err := s.CreateGenericSigTask(
		func(
			referenceBlockNumber uint64,
			quorumNumbers sdktypes.QuorumNums,
			quorumThresholdPercentages sdktypes.QuorumThresholdPercentages,
		) []interface{} {
			return []interface{}{
				AlertHeaderParam{
					MessageHash:                hash,
					QuorumNumbers:              quorumNumbers.UnderlyingType(),
					QuorumThresholdPercentages: quorumThresholdPercentages.UnderlyingType(),
					ReferenceBlockNumber:       uint32(referenceBlockNumber),
				}}
		})
	if err != nil {
		return proxyUtils.CreateSigTaskResp{}, errors.Wrap(err, "create sig task failed")
	}

	response := proxyUtils.CreateSigTaskResp{
		TaskIndex: uint64(res.TaskIndex),
		TxHash:    res.TxHash,
		SigHash:   res.SigHash,
	}

	return response, nil
}
