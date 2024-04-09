package generic

import (
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/aggregator/types"
	"github.com/alt-research/avs/core/message"
	"github.com/ethereum/go-ethereum/common"
)

type AVSGenericTasks struct {
	logger logging.Logger

	tasks         map[types.TaskIndex]*message.GenericTaskData
	finishedTasks map[message.Bytes32]*GenericFinishedTaskStatus
	nextTaskIndex types.TaskIndex
	tasksMu       sync.RWMutex

	operatorStatus   map[common.Address]*OperatorStatus
	operatorStatusMu sync.RWMutex
}

func NewAVSGenericTasks(logger logging.Logger) *AVSGenericTasks {
	return &AVSGenericTasks{
		logger: logger,

		tasks:          make(map[types.TaskIndex]*message.GenericTaskData),
		finishedTasks:  make(map[message.Bytes32]*GenericFinishedTaskStatus),
		operatorStatus: make(map[common.Address]*OperatorStatus),
	}
}

func (agg *AVSGenericTasks) GetTaskByHash(sigHash [32]byte) *message.GenericTaskData {
	agg.tasksMu.RLock()
	defer agg.tasksMu.RUnlock()

	for _, task := range agg.tasks {
		if task.TaskSigHash == sigHash {
			return task
		}
	}

	return nil
}

func (agg *AVSGenericTasks) SetNewTask(task *message.GenericTaskData) {
	agg.tasksMu.Lock()
	defer agg.tasksMu.Unlock()

	agg.tasks[task.TaskIndex] = task
}

func (agg *AVSGenericTasks) GetTaskByIndex(taskIndex types.TaskIndex) *message.GenericTaskData {
	agg.tasksMu.RLock()
	defer agg.tasksMu.RUnlock()

	res := agg.tasks[taskIndex]

	return res
}

func (agg *AVSGenericTasks) newIndex() types.TaskIndex {
	agg.tasksMu.Lock()
	defer agg.tasksMu.Unlock()

	res := agg.nextTaskIndex
	agg.nextTaskIndex += 1

	return res
}

func (agg *AVSGenericTasks) GetFinishedTaskByAlertHash(sigHash message.Bytes32) *GenericFinishedTaskStatus {
	agg.tasksMu.RLock()
	defer agg.tasksMu.RUnlock()

	return agg.finishedTasks[sigHash]
}

func (agg *AVSGenericTasks) SetFinishedTask(sigHash message.Bytes32, finished *GenericFinishedTaskStatus) {
	agg.tasksMu.Lock()
	defer agg.tasksMu.Unlock()

	agg.finishedTasks[sigHash] = finished
}
