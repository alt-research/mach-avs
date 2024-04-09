package operator

import (
	"context"
	"fmt"
	"net/rpc"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/alt-research/avs/api/grpc/aggregator"
	"github.com/alt-research/avs/core/alert"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
	"github.com/alt-research/avs/metrics"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AggregatorGRpcClient struct {
	rpcClient                  *rpc.Client
	metrics                    metrics.Metrics
	logger                     logging.Logger
	config                     config.NodeConfig
	operatorId                 sdktypes.OperatorId
	operatorAddr               common.Address
	OperatorStateRetrieverAddr common.Address
	RegistryCoordinatorAddr    common.Address
	gRPCAggregatorIpPortAddr   string
	timeout                    time.Duration
}

func NewAggregatorGRpcClient(config config.NodeConfig, operatorId sdktypes.OperatorId, operatorAddr common.Address, logger logging.Logger, metrics metrics.Metrics) (*AggregatorGRpcClient, error) {
	return &AggregatorGRpcClient{
		// set to nil so that we can create an rpc client even if the aggregator is not running
		rpcClient:                  nil,
		metrics:                    metrics,
		logger:                     logger,
		config:                     config,
		operatorId:                 operatorId,
		operatorAddr:               operatorAddr,
		OperatorStateRetrieverAddr: common.HexToAddress(config.OperatorStateRetrieverAddress),
		RegistryCoordinatorAddr:    common.HexToAddress(config.AVSRegistryCoordinatorAddress),
		gRPCAggregatorIpPortAddr:   config.AggregatorGRPCServerIpPortAddress,
		timeout:                    1 * time.Second,
	}, nil
}

// CreateAlertTaskToAggregator create a new alert task, if had existing, just return current alert task.
func (c *AggregatorGRpcClient) InitOperatorToAggregator() error {
	conn, err := grpc.Dial(
		c.gRPCAggregatorIpPortAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	n := aggregator.NewAggregatorClient(conn)
	nodeCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	request := &aggregator.InitOperatorRequest{
		Layer1ChainId:              c.config.Layer1ChainId,
		ChainId:                    c.config.Layer2ChainId,
		OperatorId:                 c.operatorId[:],
		OperatorAddress:            c.operatorAddr.Hex(),
		OperatorStateRetrieverAddr: c.config.OperatorStateRetrieverAddress,
		RegistryCoordinatorAddr:    c.config.AVSRegistryCoordinatorAddress,
	}

	c.logger.Info("Init operator to aggregator", "req", fmt.Sprintf("%#v", request))

	reply, err := n.InitOperator(nodeCtx, request)
	if err != nil {
		return fmt.Errorf("call initOperatorToAggregator failed: %v", err.Error())
	}

	if !reply.GetOk() {
		return fmt.Errorf("initOperatorToAggregator failed: %v", reply.GetReason())
	}

	return nil

}

// CreateAlertTaskToAggregator create a new alert task, if had existing, just return current alert task.
func (c *AggregatorGRpcClient) CreateAlertTaskToAggregator(alertHash [32]byte) (*message.AlertTaskInfo, error) {
	conn, err := grpc.Dial(
		c.gRPCAggregatorIpPortAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	n := aggregator.NewAggregatorClient(conn)
	nodeCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	request := &aggregator.CreateTaskRequest{
		AlertHash: alertHash[:],
	}

	c.logger.Info("CreateAlertTask to aggregator", "req", fmt.Sprintf("%#v", request))

	reply, err := n.CreateTask(nodeCtx, request)
	if err != nil {
		return nil, fmt.Errorf("call CreateAlertTask failed: %v", err.Error())
	}

	info, err := message.NewAlertTaskInfo(reply.GetInfo())
	if err != nil {
		return nil, fmt.Errorf("call CreateAlertTask failed by decode alert info: %v", err.Error())
	}

	return info, nil
}

// SendSignedTaskResponseToAggregator sends a signed task response to the aggregator.
// it is meant to be ran inside a go thread, so doesn't return anything.
// this is because sending the signed task response to the aggregator is time sensitive,
// so there is no point in retrying if it fails for a few times.
// Currently hardcoded to retry sending the signed task response 5 times, waiting 2 seconds in between each attempt.
func (c *AggregatorGRpcClient) SendSignedTaskResponseToAggregator(signedTaskResponse *message.SignedTaskRespRequest, resChan chan alert.AlertResponse) {
	conn, err := grpc.Dial(
		c.gRPCAggregatorIpPortAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		resChan <- alert.AlertResponse{
			Err: err,
			Msg: "dial initOperatorToAggregator connection failed",
		}
		return
	}

	n := aggregator.NewAggregatorClient(conn)
	nodeCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	request := &aggregator.SignedTaskRespRequest{
		Alert:                    signedTaskResponse.Alert.ToPbType(),
		OperatorRequestSignature: signedTaskResponse.BlsSignature.Serialize(),
		OperatorId:               signedTaskResponse.OperatorId[:],
	}

	c.logger.Info("CreateAlertTask to aggregator", "req", fmt.Sprintf("%#v", request))

	response, err := n.ProcessSignedTaskResponse(nodeCtx, request)
	if err != nil {
		resChan <- alert.AlertResponse{
			Err: err,
			Msg: fmt.Sprintf("call CreateAlertTask failed by %v", err),
		}
		return
	}

	c.logger.Info("Signed task response header accepted by aggregator.", "response", fmt.Sprintf("%#v", response))

	res := alert.AlertResponse{
		Code:      0,
		TaskIndex: signedTaskResponse.Alert.TaskIndex,
	}
	copy(res.TxHash[:], response.GetTxHash()[:32])

	c.logger.Info("Signed task resp", "response", res)
	c.metrics.IncNumTasksAcceptedByAggregator()

	resChan <- res
}
