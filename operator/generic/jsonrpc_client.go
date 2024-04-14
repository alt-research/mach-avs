package generic_operator

import (
	"context"
	"fmt"
	"net/rpc"
	"time"

	"github.com/pkg/errors"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	aggGeneric "github.com/alt-research/avs/aggregator/generic"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
	"github.com/alt-research/avs/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

type GenericAggregatorClient struct {
	rpcClient                   *rpc.Client
	metrics                     metrics.Metrics
	logger                      logging.Logger
	config                      config.NodeConfig
	avsCfg                      config.GenericAVSConfig
	operatorId                  sdktypes.OperatorId
	operatorAddr                common.Address
	jsonRPCAggregatorIpPortAddr string
	timeout                     time.Duration
}

func NewGenericAggregatorClient(config config.NodeConfig, avsCfg config.GenericAVSConfig, operatorId sdktypes.OperatorId, operatorAddr common.Address, logger logging.Logger, metrics metrics.Metrics) (*GenericAggregatorClient, error) {
	return &GenericAggregatorClient{
		// set to nil so that we can create an rpc client even if the aggregator is not running
		rpcClient:    nil,
		metrics:      metrics,
		logger:       logger,
		config:       config,
		avsCfg:       avsCfg,
		operatorId:   operatorId,
		operatorAddr: operatorAddr,

		jsonRPCAggregatorIpPortAddr: config.AggregatorJSONRPCServerIpPortAddr,
		timeout:                     1 * time.Second,
	}, nil
}

// CreateAlertTaskToAggregator create a new alert task, if had existing, just return current alert task.
func (c *GenericAggregatorClient) InitOperatorToAggregator() error {
	client, err := gethrpc.DialContext(context.Background(), c.jsonRPCAggregatorIpPortAddr)
	if err != nil {
		return fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	var res aggGeneric.InitOperatorResponse

	err = client.CallContext(
		context.Background(), &res, "genericAggregator_initOperator",
		c.avsCfg.AVSName,
		c.config.Layer1ChainId,
		hexutil.Bytes(c.operatorId[:]),
		c.operatorAddr.Hex(),
		c.avsCfg.OperatorStateRetrieverAddress.Hex(),
		c.avsCfg.AVSRegistryCoordinatorAddress.Hex(),
	)
	if err != nil {
		return fmt.Errorf("call initOperatorToAggregator failed: %v", err.Error())
	}

	if !res.Ok {
		return fmt.Errorf("initOperatorToAggregator failed: %v", res.Reason)
	}

	return nil

}

// CreateAlertTaskToAggregator create a new alert task, if had existing, just return current alert task.
func (c *GenericAggregatorClient) CreateTaskToAggregator(sigHash [32]byte, method string, params []interface{}) (*message.CreateGenericTaskResponse, error) {
	client, err := gethrpc.DialContext(context.Background(), c.jsonRPCAggregatorIpPortAddr)
	if err != nil {
		return nil, fmt.Errorf("dial CreateAlertTask connection failed: %v", err.Error())
	}

	paramsRaw, err := packCallParams(
		c.avsCfg.AVSName,
		&c.avsCfg.Abi,
		method, params,
	)
	if err != nil {
		return nil, errors.Wrap(err, "pack call params failed")
	}

	var res message.CreateGenericTaskResponse
	err = client.CallContext(
		context.Background(), &res, "genericAggregator_createTask",
		c.avsCfg.AVSName,
		hexutil.Bytes(sigHash[:]),
		method,
		hexutil.Bytes(paramsRaw),
	)

	if err != nil {
		return nil, fmt.Errorf("call CreateAlertTask failed: %v", err.Error())
	}

	return &res, nil
}

// SendSignedTaskResponseToAggregator sends a signed task response to the aggregator.
// it is meant to be ran inside a go thread, so doesn't return anything.
// this is because sending the signed task response to the aggregator is time sensitive,
// so there is no point in retrying if it fails for a few times.
// Currently hardcoded to retry sending the signed task response 5 times, waiting 2 seconds in between each attempt.
func (c *GenericAggregatorClient) SendSignedTaskResponseToAggregator(
	method string,
	params []interface{},
	operatorRequestSignature bls.Signature,
	operatorId sdktypes.OperatorId,
	taskInfo message.CreateGenericTaskResponse,
	resChan chan GenericResponse,
) {
	client, err := gethrpc.DialContext(context.Background(), c.jsonRPCAggregatorIpPortAddr)
	if err != nil {
		resChan <- GenericResponse{
			Err: err,
			Msg: "dial initOperatorToAggregator connection failed",
		}
		return
	}

	qperatorRequestSignature := operatorRequestSignature.Serialize()
	paramsRaw, err := packCallParams(
		c.avsCfg.AVSName,
		&c.avsCfg.Abi,
		method, params,
	)
	if err != nil {
		resChan <- GenericResponse{
			Err: err,
			Msg: "pack call params failed",
		}
		return
	}

	c.logger.Info("CreateAlertTask to aggregator", "alert", fmt.Sprintf("%#v", taskInfo))

	var resp aggGeneric.SignedTaskRespResponse
	err = client.CallContext(
		context.Background(), &resp, "genericAggregator_processSignedTaskResponse",
		c.avsCfg.AVSName,
		taskInfo,
		method,
		hexutil.Bytes(paramsRaw),
		hexutil.Bytes(qperatorRequestSignature),
		hexutil.Bytes(operatorId[:]),
	)
	if err != nil {
		resChan <- GenericResponse{
			Err: err,
			Msg: "call CreateAlertTask failed",
		}
		return
	}

	c.logger.Info("Signed task response header accepted by aggregator.", "response", fmt.Sprintf("%#v", resp))

	res := GenericResponse{
		Code:      0,
		TaskIndex: taskInfo.TaskIndex,
	}
	copy(res.TxHash[:], resp.TxHash[:32])

	c.logger.Info("Signed task resp", "response", res)
	c.metrics.IncNumTasksAcceptedByAggregator()

	resChan <- res
}
