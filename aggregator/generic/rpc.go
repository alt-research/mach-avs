package generic

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/api/grpc/aggregator"
	"github.com/alt-research/avs/core"
	"github.com/alt-research/avs/core/message"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/node"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type JsonRpcServer struct {
	logger  logging.Logger
	handler JsonRpcHandler
	vhosts  []string
	cors    []string
	wg      *sync.WaitGroup
}

func NewJsonRpcServer(logger logging.Logger, aggreagtor *AVSGenericServices, vhosts []string, cors []string) *JsonRpcServer {
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
		Namespace: "genericAggregator",
		Service:   &s.handler,
	}
}

func (s *JsonRpcServer) StartServer(ctx context.Context, serverIpPortAddr string) {
	s.logger.Info("Start JSON RPC Server", "addr", serverIpPortAddr)

	rpcAPI := []gethrpc.API{s.GetAPI()}

	srv := gethrpc.NewServer()
	err := node.RegisterApis(rpcAPI, []string{"genericAggregator"}, srv)
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
	aggreagtor *AVSGenericServices
}

type InitOperatorResponse struct {
	Ok     bool   `json:"ok"`
	Reason string `json:"reason"`
}

func (h *JsonRpcHandler) InitOperator(
	ctx context.Context,
	avsName string,
	layer1ChainId uint32,
	operatorId hexutil.Bytes,
	operatorAddress string,
	operatorStateRetrieverAddr string,
	registryCoordinatorAddr string,
) (InitOperatorResponse, error) {
	h.logger.Debug("handle init operator", "name", avsName)

	req, err := message.NewInitOperatorRequest(&aggregator.InitOperatorRequest{
		Layer1ChainId:              layer1ChainId,
		OperatorId:                 operatorId,
		OperatorAddress:            operatorAddress,
		OperatorStateRetrieverAddr: operatorStateRetrieverAddr,
		RegistryCoordinatorAddr:    registryCoordinatorAddr,
	})
	if err != nil {
		return InitOperatorResponse{}, fmt.Errorf("initOperator parse request falied: %v", err)
	}

	res, err := h.aggreagtor.InitOperator(avsName, &message.InitOperatorDatas{
		AVSName:                    avsName,
		Layer1ChainId:              layer1ChainId,
		OperatorId:                 req.OperatorId,
		OperatorAddress:            req.OperatorAddress,
		OperatorStateRetrieverAddr: req.OperatorStateRetrieverAddr,
		RegistryCoordinatorAddr:    req.RegistryCoordinatorAddr,
	})

	resp := InitOperatorResponse{
		Ok:     res.Ok,
		Reason: res.Res,
	}

	return resp, nil
}

type GenericTaskInfo struct {
	// The hash of alert
	TaskSigHash hexutil.Bytes `json:"sig_hash"`
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
	avsName string,
	hash hexutil.Bytes,
	method string,
	params hexutil.Bytes,
) (GenericTaskInfo, error) {
	hashBytes32, err := message.NewBytes32(hash)
	if err != nil {
		return GenericTaskInfo{}, fmt.Errorf("createTask parse request falied: %v", err)
	}

	res, err := h.aggreagtor.CreateTask(avsName, hashBytes32, method, params)
	if err != nil {
		return GenericTaskInfo{}, fmt.Errorf("createTask process request falied: %v", err)
	}

	resp := GenericTaskInfo{
		TaskSigHash:                res.TaskSigHash[:],
		QuorumNumbers:              res.QuorumNumbers.UnderlyingType(),
		QuorumThresholdPercentages: res.QuorumThresholdPercentages.UnderlyingType(),
		TaskIndex:                  res.TaskIndex,
		ReferenceBlockNumber:       res.ReferenceBlockNumber,
	}

	return resp, nil
}

type SignedTaskRespResponse struct {
	// If need reply
	Reply bool `json:"reply"`
	// The tx hash of send
	TxHash hexutil.Bytes `json:"tx_hash"`
}

func (h *JsonRpcHandler) ProcessSignedTaskResponse(
	ctx context.Context,
	avsName string,
	taskInfo GenericTaskInfo,
	method string,
	params hexutil.Bytes,
	operatorRequestSignature hexutil.Bytes,
	operatorId hexutil.Bytes,
) (SignedTaskRespResponse, error) {
	callParams, err := h.aggreagtor.parseParams(avsName, method, params)
	if err != nil {
		return SignedTaskRespResponse{}, errors.Wrapf(err, "parse %v:%v params failed", avsName, method)
	}

	taskSigHash, err := message.NewBytes32(taskInfo.TaskSigHash)
	if err != nil {
		return SignedTaskRespResponse{}, errors.Wrapf(err, "parse %v:%v taskSigHash failed", avsName, method)
	}

	operatorIdBytes32, err := message.NewBytes32(operatorId)
	if err != nil {
		return SignedTaskRespResponse{}, errors.Wrapf(err, "parse %v:%v operatorId failed", avsName, method)
	}

	g1Point := bls.NewZeroG1Point().Deserialize(operatorRequestSignature)
	sign := bls.Signature{G1Point: g1Point}

	res, err := h.aggreagtor.ProcessSignedTaskResponse(avsName, &message.GenericTaskData{
		TaskIndex:                  taskInfo.TaskIndex,
		TaskSigHash:                taskSigHash,
		QuorumNumbers:              core.ConvertQuorumNumbersFromBytes(taskInfo.QuorumNumbers),
		QuorumThresholdPercentages: core.ConvertQuorumThresholdPercentagesFromBytes(taskInfo.QuorumThresholdPercentages),
		CallMethod:                 method,
		CallParams:                 callParams,
		ReferenceBlockNumber:       taskInfo.ReferenceBlockNumber,
	}, sign, operatorIdBytes32.UnderlyingType())

	if err != nil {
		return SignedTaskRespResponse{}, errors.Wrapf(err, "call %v:%v process signed task response failed", avsName, method)
	}

	resp := SignedTaskRespResponse{
		Reply:  false,
		TxHash: res.Bytes(),
	}

	return resp, nil
}
