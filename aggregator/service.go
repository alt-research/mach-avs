package aggregator

import (
	"context"
	"fmt"
	"sync"
	"time"

	sdkclients "github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/services/avsregistry"
	blsagg "github.com/Layr-Labs/eigensdk-go/services/bls_aggregation"
	oppubkeysserv "github.com/Layr-Labs/eigensdk-go/services/operatorpubkeys"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"

	"github.com/alt-research/avs/aggregator/rpc"
	"github.com/alt-research/avs/aggregator/types"
	"github.com/alt-research/avs/core/chainio"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
	"github.com/ethereum/go-ethereum/common"
)

type AggregatorService struct {
	logger logging.Logger
	cfg    *config.Config

	avsReader chainio.AvsReaderer
	ethClient eth.Client

	blsAggregationService blsagg.BlsAggregationService
	tasks                 map[types.TaskIndex]*message.AlertTaskInfo
	tasksMu               sync.RWMutex
	finishedTasks         map[[32]byte]*FinishedTaskStatus
	finishedTasksMu       sync.RWMutex
	nextTaskIndex         types.TaskIndex
	nextTaskIndexMu       sync.RWMutex
	operatorStatus        map[common.Address]*OperatorStatus
	operatorStatusMu      sync.RWMutex
}

// NewAggregator creates a new Aggregator with the provided config.
func NewAggregatorService(c *config.Config) (*AggregatorService, error) {
	avsReader, err := chainio.BuildAvsReaderFromConfig(c)
	if err != nil {
		c.Logger.Error("Cannot create avsReader", "err", err)
		return nil, err
	}

	chainioConfig := sdkclients.BuildAllConfig{
		EthHttpUrl:                 c.EthHttpRpcUrl,
		EthWsUrl:                   c.EthWsRpcUrl,
		RegistryCoordinatorAddr:    c.RegistryCoordinatorAddr.String(),
		OperatorStateRetrieverAddr: c.OperatorStateRetrieverAddr.String(),
		AvsName:                    avsName,
		PromMetricsIpPortAddress:   ":9090",
	}
	clients, err := sdkclients.BuildAll(chainioConfig, c.PrivateKey, c.Logger)
	if err != nil {
		c.Logger.Errorf("Cannot create sdk clients", "err", err)
		return nil, err
	}

	operatorPubkeysService := oppubkeysserv.NewOperatorPubkeysServiceInMemory(context.Background(), clients.AvsRegistryChainSubscriber, clients.AvsRegistryChainReader, c.Logger)
	avsRegistryService := avsregistry.NewAvsRegistryServiceChainCaller(avsReader, operatorPubkeysService, c.Logger)
	blsAggregationService := blsagg.NewBlsAggregatorService(avsRegistryService, c.Logger)

	return &AggregatorService{
		logger:                c.Logger,
		avsReader:             avsReader,
		ethClient:             clients.EthHttpClient,
		blsAggregationService: blsAggregationService,
		tasks:                 make(map[types.TaskIndex]*message.AlertTaskInfo),
		finishedTasks:         make(map[[32]byte]*FinishedTaskStatus),
		operatorStatus:        make(map[common.Address]*OperatorStatus),
		cfg:                   c,
	}, nil
}

var _ rpc.AggregatorRpcHandler = (*AggregatorService)(nil)

func (agg *AggregatorService) GetTaskByAlertHash(alertHash [32]byte) *message.AlertTaskInfo {
	agg.tasksMu.RLock()
	defer agg.tasksMu.RUnlock()

	for _, task := range agg.tasks {
		if task.AlertHash == alertHash {
			return task
		}
	}

	return nil
}

func (agg *AggregatorService) GetTaskByIndex(taskIndex types.TaskIndex) *message.AlertTaskInfo {
	agg.tasksMu.RLock()
	defer agg.tasksMu.RUnlock()

	res := agg.tasks[taskIndex]

	return res
}

func (agg *AggregatorService) newIndex() types.TaskIndex {
	agg.nextTaskIndexMu.Lock()
	defer agg.nextTaskIndexMu.Unlock()

	res := agg.nextTaskIndex
	agg.nextTaskIndex += 1

	return res
}

func (agg *AggregatorService) GetFinishedTaskByAlertHash(alertHash [32]byte) *FinishedTaskStatus {
	agg.finishedTasksMu.RLock()
	defer agg.finishedTasksMu.RUnlock()

	return agg.finishedTasks[alertHash]
}

func (agg *AggregatorService) SetFinishedTask(alertHash [32]byte, finished *FinishedTaskStatus) {
	agg.finishedTasksMu.Lock()
	defer agg.finishedTasksMu.Unlock()

	agg.finishedTasks[alertHash] = finished
}

// rpc endpoint which is called by operator
// will init operator, just for keep config valid
func (agg *AggregatorService) InitOperator(req *message.InitOperatorRequest) (*message.InitOperatorResponse, error) {
	agg.logger.Infof("Received InitOperator: %#v", req)

	reply := &message.InitOperatorResponse{
		Ok: false,
	}

	if agg.cfg.OperatorStateRetrieverAddr != req.OperatorStateRetrieverAddr {
		reply.Res = fmt.Sprintf("OperatorStateRetrieverAddr invaild, expect %s", agg.cfg.OperatorStateRetrieverAddr.Hex())
		return reply, nil
	}

	if agg.cfg.RegistryCoordinatorAddr != req.RegistryCoordinatorAddr {
		reply.Res = fmt.Sprintf("RegistryCoordinatorAddr invaild, expect %s", agg.cfg.RegistryCoordinatorAddr.Hex())
		return reply, nil
	}

	if agg.cfg.Layer1ChainId != req.Layer1ChainId {
		reply.Res = fmt.Sprintf("Layer1ChainId invaild, expect %d", agg.cfg.Layer1ChainId)
		return reply, nil
	}

	if agg.cfg.Layer2ChainId != req.ChainId {
		reply.Res = fmt.Sprintf("Layer2ChainId invaild, expect %d", agg.cfg.Layer2ChainId)
		return reply, nil
	}

	agg.operatorStatusMu.Lock()
	defer agg.operatorStatusMu.Unlock()

	agg.operatorStatus[req.OperatorAddress] = &OperatorStatus{
		LastTime:   time.Now().Unix(),
		OperatorId: req.OperatorId,
	}

	reply.Ok = true

	agg.logger.Infof("new operator status: %s", req.OperatorAddress.Hex())

	return reply, nil
}

// rpc endpoint which is called by operator
// will try to init the task, if currently had a same task for the alert,
// it will return the existing task.
func (agg *AggregatorService) CreateTask(req *message.CreateTaskRequest) (*message.CreateTaskResponse, error) {
	agg.logger.Infof("Received CreateTask: %#v", req)

	finished := agg.GetFinishedTaskByAlertHash(req.AlertHash)
	if finished != nil {
		return nil, fmt.Errorf("the task 0x%x already finished: 0x%x", req.AlertHash, finished.TxHash)
	}

	task := agg.GetTaskByAlertHash(req.AlertHash)
	if task == nil {
		agg.logger.Info("create new task", "alert", req.AlertHash)
		taskIndex := agg.newIndex()

		var err error
		task, err = agg.sendNewTask(req.AlertHash, taskIndex)

		if err != nil {
			agg.logger.Error("send new task failed", "err", err)
			return nil, err
		}
	}

	return &message.CreateTaskResponse{Info: *task}, nil
}

// rpc endpoint which is called by operator
// reply doesn't need to be checked. If there are no errors, the task response is accepted
// rpc framework forces a reply type to exist, so we put bool as a placeholder
func (agg *AggregatorService) ProcessSignedTaskResponse(signedTaskResponse *message.SignedTaskRespRequest) (*message.SignedTaskRespResponse, error) {
	agg.logger.Infof("Received signed task response: %#v", signedTaskResponse)
	taskIndex := signedTaskResponse.Alert.TaskIndex
	taskResponseDigest, err := signedTaskResponse.Alert.SignHash()
	if err != nil {
		return nil, err
	}

	agg.logger.Infof("ProcessNewSignature11: %#v", signedTaskResponse.Alert.TaskIndex)

	if task := agg.GetTaskByIndex(taskIndex); task == nil {
		agg.logger.Error("ProcessNewSignature error by no task exist", "taskIndex", taskIndex)
		return nil, fmt.Errorf("task not found")
	}

	agg.logger.Infof("ProcessNewSignature: %#v", signedTaskResponse.Alert.TaskIndex)
	err = agg.blsAggregationService.ProcessNewSignature(
		context.Background(), taskIndex, taskResponseDigest,
		&signedTaskResponse.BlsSignature, signedTaskResponse.OperatorId,
	)

	if err != nil {
		agg.logger.Error("ProcessNewSignature error", "err", err)
	}

	return &message.SignedTaskRespResponse{}, err
}

// GetResponseChannel returns the single channel that meant to be used as the response channel
// Any task that is completed (see the completion criterion in the comment above InitializeNewTask)
// will be sent on this channel along with all the necessary information to call BLSSignatureChecker onchain
func (agg *AggregatorService) GetResponseChannel() <-chan blsagg.BlsAggregationServiceResponse {
	return agg.blsAggregationService.GetResponseChannel()
}

// sendNewTask sends a new task to the task manager contract, and updates the Task dict struct
// with the information of operators opted into quorum 0 at the block of task creation.
func (agg *AggregatorService) sendNewTask(alertHash [32]byte, taskIndex types.TaskIndex) (*message.AlertTaskInfo, error) {
	agg.logger.Info("Aggregator sending new task", "alert", alertHash, "task", taskIndex)

	// TODO: use cfg
	quorumNumbersValue := []sdktypes.QuorumNum{0}
	quorumThresholdPercentagesValue := []sdktypes.QuorumThresholdPercentage{100}

	var err error

	var referenceBlockNumber uint64
	if referenceBlockNumber, err = agg.ethClient.BlockNumber(context.Background()); err != nil {
		return nil, err
	}

	// the reference block number must < the current block number.
	referenceBlockNumber -= 1

	agg.logger.Info("get from layer1", "referenceBlockNumber", referenceBlockNumber)

	quorumNumbers, err := agg.avsReader.GetQuorumsByBlockNumber(context.Background(), uint32(referenceBlockNumber))
	if err != nil {
		agg.logger.Error("GetQuorumCountByBlockNumber failed", "err", err)
		return nil, err
	}

	quorumThresholdPercentages, err := agg.avsReader.GetQuorumThresholdPercentages(context.Background(), uint32(referenceBlockNumber), quorumNumbers)
	if err != nil {
		agg.logger.Error("GetQuorumThresholdPercentages failed", "err", err)
		return nil, err
	}

	agg.logger.Infof("quorum %v %v", quorumNumbers, quorumThresholdPercentages)

	newAlertTask := &message.AlertTaskInfo{
		AlertHash:                  alertHash,
		QuorumNumbers:              quorumNumbersValue,
		QuorumThresholdPercentages: quorumThresholdPercentagesValue,
		TaskIndex:                  taskIndex,
		ReferenceBlockNumber:       referenceBlockNumber,
	}

	agg.tasksMu.Lock()
	agg.tasks[taskIndex] = newAlertTask
	agg.tasksMu.Unlock()

	// TODO(samlaf): we use seconds for now, but we should ideally pass a blocknumber to the blsAggregationService
	// and it should monitor the chain and only expire the task aggregation once the chain has reached that block number.
	taskTimeToExpiry := taskChallengeWindowBlock * blockTimeDuration

	agg.logger.Infof("InitializeNewTask %v %v", taskIndex, taskTimeToExpiry)
	err = agg.blsAggregationService.InitializeNewTask(
		taskIndex,
		uint32(newAlertTask.ReferenceBlockNumber),
		newAlertTask.QuorumNumbers,
		newAlertTask.QuorumThresholdPercentages,
		taskTimeToExpiry,
	)
	if err != nil {
		agg.logger.Error("InitializeNewTask failed", "err", err)
		return nil, err
	}
	return newAlertTask, nil
}
