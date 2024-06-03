package rpc

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/legacy/api/grpc/aggregator"
	"github.com/alt-research/avs/legacy/core/message"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/node"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

type JsonRpcServer struct {
	logger  logging.Logger
	handler JsonRpcHandler
	vhosts  []string
	cors    []string
	wg      *sync.WaitGroup
}

func NewJsonRpcServer(logger logging.Logger, aggreagtor AggregatorRpcHandler, vhosts []string, cors []string) *JsonRpcServer {
	return &JsonRpcServer{
		logger: logger,
		handler: JsonRpcHandler{
			logger:     logger,
			aggreagtor: aggreagtor,
		},
		vhosts: vhosts,
		cors:   cors,
		wg:     &sync.WaitGroup{},
	}
}

func (s *JsonRpcServer) GetAPI() gethrpc.API {
	return gethrpc.API{
		Namespace: "aggregator",
		Service:   &s.handler,
	}
}

func (s *JsonRpcServer) StartServer(ctx context.Context, serverIpPortAddr string) {
	s.logger.Info("Start JSON RPC Server", "addr", serverIpPortAddr)

	rpcAPI := []gethrpc.API{s.GetAPI()}

	srv := gethrpc.NewServer()
	err := node.RegisterApis(rpcAPI, []string{"aggregator"}, srv)
	if err != nil {
		s.logger.Fatalf("Could not register API: %w", err)
	}
	handler := node.NewHTTPHandlerStack(srv, s.cors, s.vhosts, nil)

	httpServer, addr, err := node.StartHTTPEndpoint(serverIpPortAddr, gethrpc.DefaultHTTPTimeouts, handler)
	if err != nil {
		s.logger.Fatalf("Could not start RPC api: %v", err)
	}

	extapiURL := fmt.Sprintf("http://%v/", addr)
	s.logger.Info("HTTP endpoint opened", "url", extapiURL)

	serverErr := make(chan error, 1)

	s.wg.Add(1)
	defer s.wg.Done()

	select {
	case <-ctx.Done():
		s.logger.Info("Stop JSON RPC Server by Done")
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			s.logger.Errorf("Stop JSON RPC Server by error: %v", err.Error())
		}
	case err = <-serverErr:
	}

	if err != nil {
		s.logger.Error("JSON RPC Server serve stopped by error", "err", err)
	} else {
		s.logger.Info("JSON RPC Server serve stopped")
	}
}

func (s *JsonRpcServer) Wait() {
	s.wg.Wait()
}

type JsonRpcHandler struct {
	logger     logging.Logger
	aggreagtor AggregatorRpcHandler
}

type InitOperatorResponse struct {
	Ok     bool   `json:"ok"`
	Reason string `json:"reason"`
}

func (h *JsonRpcHandler) InitOperator(
	ctx context.Context,
	layer1ChainId uint32,
	chainId uint32,
	operatorId hexutil.Bytes,
	operatorAddress string,
	operatorStateRetrieverAddr string,
	registryCoordinatorAddr string,
) (InitOperatorResponse, error) {
	req, err := message.NewInitOperatorRequest(&aggregator.InitOperatorRequest{
		Layer1ChainId:              layer1ChainId,
		ChainId:                    chainId,
		OperatorId:                 operatorId,
		OperatorAddress:            operatorAddress,
		OperatorStateRetrieverAddr: operatorStateRetrieverAddr,
		RegistryCoordinatorAddr:    registryCoordinatorAddr,
	})
	if err != nil {
		return InitOperatorResponse{}, fmt.Errorf("initOperator parse request falied: %v", err)
	}

	res, err := h.aggreagtor.InitOperator(req)
	if err != nil {
		return InitOperatorResponse{}, fmt.Errorf("initOperator process request falied: %v", err)
	}

	resp := InitOperatorResponse{
		Ok:     res.Ok,
		Reason: res.Res,
	}

	return resp, nil
}

type AlertTaskInfo struct {
	// The hash of alert
	AlertHash hexutil.Bytes `json:"alert_hash"`
	// QuorumNumbers of task
	QuorumNumbers []uint8 `json:"quorum_numbers"`
	// QuorumThresholdPercentages of task
	QuorumThresholdPercentages []uint8 `json:"quorum_threshold_percentages"`
	// TaskIndex
	TaskIndex uint32 `json:"task_index"`
	// ReferenceBlockNumber
	ReferenceBlockNumber uint64 `json:"reference_block_number"`
}

func (h *JsonRpcHandler) CreateTask(
	ctx context.Context,
	alertHash hexutil.Bytes,
) (AlertTaskInfo, error) {
	req, err := message.NewCreateTaskRequest(&aggregator.CreateTaskRequest{
		AlertHash: alertHash,
	})
	if err != nil {
		return AlertTaskInfo{}, fmt.Errorf("createTask parse request falied: %v", err)
	}

	res, err := h.aggreagtor.CreateTask(req)
	if err != nil {
		return AlertTaskInfo{}, fmt.Errorf("createTask process request falied: %v", err)
	}

	info := res.Info.ToPbType()

	resp := AlertTaskInfo{
		AlertHash:                  info.AlertHash,
		QuorumNumbers:              info.QuorumNumbers,
		QuorumThresholdPercentages: info.QuorumThresholdPercentages,
		TaskIndex:                  info.TaskIndex,
		ReferenceBlockNumber:       info.ReferenceBlockNumber,
	}

	return resp, nil
}

type SignedTaskRespResponse struct {
	// If need reply
	Reply bool `json:"reply"`
	// The tx hash of send
	TxHash []byte `json:"tx_hash"`
}

func (h *JsonRpcHandler) ProcessSignedTaskResponse(
	ctx context.Context,
	alertInfo AlertTaskInfo,
	operatorRequestSignature hexutil.Bytes,
	operatorId hexutil.Bytes,
) (SignedTaskRespResponse, error) {
	req, err := message.NewSignedTaskRespRequest(&aggregator.SignedTaskRespRequest{
		Alert: &aggregator.AlertTaskInfo{
			AlertHash:                  alertInfo.AlertHash,
			QuorumNumbers:              alertInfo.QuorumNumbers,
			QuorumThresholdPercentages: alertInfo.QuorumThresholdPercentages,
			TaskIndex:                  alertInfo.TaskIndex,
			ReferenceBlockNumber:       alertInfo.ReferenceBlockNumber,
		},
		OperatorRequestSignature: operatorRequestSignature,
		OperatorId:               operatorId,
	})
	if err != nil {
		return SignedTaskRespResponse{}, fmt.Errorf("processSignedTaskResponse parse request falied: %v", err)
	}

	res, err := h.aggreagtor.ProcessSignedTaskResponse(req)
	if err != nil {
		return SignedTaskRespResponse{}, fmt.Errorf("processSignedTaskResponse process request falied: %v", err)
	}

	resp := SignedTaskRespResponse{
		Reply:  res.Reply,
		TxHash: res.TxHash[:],
	}

	return resp, nil
}
