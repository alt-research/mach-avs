package operator

import (
	"fmt"
	"net/rpc"
	"strings"
	"time"

	"github.com/alt-research/avs/core/alert"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
	"github.com/alt-research/avs/metrics"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
)

type AggregatorRpcClienter interface {
	InitOperatorToAggregator() error
	CreateAlertTaskToAggregator(alertHash [32]byte) (*message.AlertTaskInfo, error)
	SendSignedTaskResponseToAggregator(signedTaskResponse *message.SignedTaskRespRequest, resChan chan alert.AlertResponse)
}
type AggregatorRpcClient struct {
	rpcClient                  *rpc.Client
	metrics                    metrics.Metrics
	logger                     logging.Logger
	config                     config.NodeConfig
	operatorId                 sdktypes.OperatorId
	operatorAddr               common.Address
	OperatorStateRetrieverAddr common.Address
	RegistryCoordinatorAddr    common.Address
	aggregatorIpPortAddr       string
}

func NewAggregatorRpcClient(config config.NodeConfig, operatorId sdktypes.OperatorId, operatorAddr common.Address, logger logging.Logger, metrics metrics.Metrics) (*AggregatorRpcClient, error) {
	return &AggregatorRpcClient{
		// set to nil so that we can create an rpc client even if the aggregator is not running
		rpcClient:                  nil,
		metrics:                    metrics,
		logger:                     logger,
		config:                     config,
		operatorId:                 operatorId,
		operatorAddr:               operatorAddr,
		OperatorStateRetrieverAddr: common.HexToAddress(config.OperatorStateRetrieverAddress),
		RegistryCoordinatorAddr:    common.HexToAddress(config.AVSRegistryCoordinatorAddress),
		aggregatorIpPortAddr:       config.AggregatorServerIpPortAddress,
	}, nil
}

func (c *AggregatorRpcClient) dialAggregatorRpcClient() error {
	client, err := rpc.DialHTTP("tcp", c.aggregatorIpPortAddr)
	if err != nil {
		return err
	}
	c.rpcClient = client
	return nil
}

// CreateAlertTaskToAggregator create a new alert task, if had existing, just return current alert task.
func (c *AggregatorRpcClient) InitOperatorToAggregator() error {
	if c.rpcClient == nil {
		c.logger.Info("rpc client is nil. Dialing aggregator rpc client")
		err := c.dialAggregatorRpcClient()
		if err != nil {
			c.logger.Error("Could not dial aggregator rpc client. Not sending signed task response header to aggregator. Is aggregator running?", "err", err)
			return err
		}
	}
	// we don't check this bool. It's just needed because rpc.Call requires rpc methods to have a return value
	var reply message.InitOperatorResponse
	req := message.InitOperatorRequest{
		Layer1ChainId:              c.config.Layer1ChainId,
		ChainId:                    c.config.Layer2ChainId,
		OperatorId:                 c.operatorId,
		OperatorAddress:            c.operatorAddr,
		OperatorStateRetrieverAddr: c.OperatorStateRetrieverAddr,
		RegistryCoordinatorAddr:    c.RegistryCoordinatorAddr,
	}

	c.logger.Info("Create task header to aggregator", "req", fmt.Sprintf("%#v", req))

	for i := 0; i < 5; i++ {
		err := c.rpcClient.Call("Aggregator.InitOperator", req, &reply)
		if err != nil {
			c.logger.Info("Received error from aggregator", "err", err)
			if strings.Contains(err.Error(), "already finished") {
				return err
			}
		} else {
			c.logger.Info("init operator accepted by aggregator.", "reply", reply)
			c.metrics.IncNumTasksAcceptedByAggregator()

			if !reply.Ok {
				if reply.Res != "" {
					return fmt.Errorf("init operator failed by %s", reply.Res)
				} else {
					return fmt.Errorf("init operator failed by unknown")
				}
			}

			return nil
		}
		c.logger.Infof("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
	c.logger.Errorf("Could not send init operator to aggregator. Tried 5 times.")

	return fmt.Errorf("Could not send init operator to aggregator")
}

// CreateAlertTaskToAggregator create a new alert task, if had existing, just return current alert task.
func (c *AggregatorRpcClient) CreateAlertTaskToAggregator(alertHash [32]byte) (*message.AlertTaskInfo, error) {
	if c.rpcClient == nil {
		c.logger.Info("rpc client is nil. Dialing aggregator rpc client")
		err := c.dialAggregatorRpcClient()
		if err != nil {
			c.logger.Error("Could not dial aggregator rpc client. Not sending signed task response header to aggregator. Is aggregator running?", "err", err)
			return nil, err
		}
	}
	// we don't check this bool. It's just needed because rpc.Call requires rpc methods to have a return value
	var reply message.CreateTaskResponse
	req := message.CreateTaskRequest{
		AlertHash: alertHash,
	}

	c.logger.Info("Create task header to aggregator", "req", fmt.Sprintf("%#v", req))

	for i := 0; i < 5; i++ {
		err := c.rpcClient.Call("Aggregator.CreateTask", req, &reply)
		if err != nil {
			c.logger.Info("Received error from aggregator", "err", err)
			if strings.Contains(err.Error(), "already finished") {
				return nil, err
			}
		} else {
			c.logger.Info("create task response header accepted by aggregator.", "reply", reply)
			c.metrics.IncNumTasksAcceptedByAggregator()
			return &reply.Info, nil
		}
		c.logger.Infof("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
	c.logger.Errorf("Could not send signed task response to aggregator. Tried 5 times.")

	return nil, fmt.Errorf("Could not send signed task response to aggregator")
}

// SendSignedTaskResponseToAggregator sends a signed task response to the aggregator.
// it is meant to be ran inside a go thread, so doesn't return anything.
// this is because sending the signed task response to the aggregator is time sensitive,
// so there is no point in retrying if it fails for a few times.
// Currently hardcoded to retry sending the signed task response 5 times, waiting 2 seconds in between each attempt.
func (c *AggregatorRpcClient) SendSignedTaskResponseToAggregator(signedTaskResponse *message.SignedTaskRespRequest, resChan chan alert.AlertResponse) {
	if c.rpcClient == nil {
		c.logger.Info("rpc client is nil. Dialing aggregator rpc client")
		err := c.dialAggregatorRpcClient()
		if err != nil {
			c.logger.Error("Could not dial aggregator rpc client. Not sending signed task response header to aggregator. Is aggregator running?", "err", err)
			resChan <- alert.AlertResponse{
				Err: err,
				Msg: "Could not dial aggregator rpc client",
			}
			return
		}
	}
	// we don't check this bool. It's just needed because rpc.Call requires rpc methods to have a return value
	var response message.SignedTaskRespResponse
	// We try to send the response 5 times to the aggregator, waiting 2 times in between each attempt.
	// This is mostly only necessary for local testing, since the aggregator sometimes is not ready to process task responses
	// before the operator gets the new task created log from anvil (because blocks are mined instantly)
	// the aggregator needs to read some onchain data related to quorums before it can accept operator signed task responses.
	c.logger.Info("Sending signed task response header to aggregator", "signedTaskResponse", fmt.Sprintf("%#v", signedTaskResponse))
	var err error
	for i := 0; i < 5; i++ {
		err = c.rpcClient.Call("Aggregator.ProcessSignedTaskResponse", signedTaskResponse, &response)
		if err != nil {
			c.logger.Info("Received error from aggregator", "err", err)
		} else {
			c.logger.Info("Signed task response header accepted by aggregator.", "response", response)
			c.metrics.IncNumTasksAcceptedByAggregator()

			resChan <- alert.AlertResponse{
				Code:      0,
				TxHash:    response.TxHash,
				TaskIndex: signedTaskResponse.Alert.TaskIndex,
			}

			return
		}
		c.logger.Infof("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
	c.logger.Errorf("Could not send signed task response to aggregator. Tried 5 times.")

	resChan <- alert.AlertResponse{
		Err: fmt.Errorf("Could not send signed task response to aggregator by %v.", err),
	}
}
