package aggregator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"

	"github.com/alt-research/avs/aggregator/types"
	"github.com/alt-research/avs/core/message"
)

var (
	TaskNotFoundError400                     = errors.New("400. Task not found")
	OperatorNotPartOfTaskQuorum400           = errors.New("400. Operator not part of quorum")
	TaskResponseDigestNotFoundError500       = errors.New("500. Failed to get task response digest")
	UnknownErrorWhileVerifyingSignature400   = errors.New("400. Failed to verify signature")
	SignatureVerificationFailed400           = errors.New("400. Signature verification failed")
	CallToGetCheckSignaturesIndicesFailed500 = errors.New("500. Failed to get check signatures indices")
)

func (agg *Aggregator) startServer(ctx context.Context) error {

	err := rpc.Register(agg)
	if err != nil {
		agg.logger.Fatal("Format of service TaskManager isn't correct. ", "err", err)
	}
	rpc.HandleHTTP()
	err = http.ListenAndServe(agg.serverIpPortAddr, nil)
	if err != nil {
		agg.logger.Fatal("ListenAndServe", "err", err)
	}

	return nil
}

func (agg *Aggregator) GetTaskByAlertHash(alertHash [32]byte) *message.AlertTaskInfo {
	agg.tasksMu.RLock()
	defer agg.tasksMu.RUnlock()

	for _, task := range agg.tasks {
		if task.AlertHash == alertHash {
			return task
		}
	}

	return nil
}

func (agg *Aggregator) GetTaskByIndex(taskIndex types.TaskIndex) *message.AlertTaskInfo {
	agg.tasksMu.RLock()
	defer agg.tasksMu.RUnlock()

	res, _ := agg.tasks[taskIndex]

	return res
}

func (agg *Aggregator) newIndex() types.TaskIndex {
	agg.tasksMu.Lock()
	defer agg.tasksMu.Unlock()

	res := agg.nextTaskIndex
	agg.nextTaskIndex += 1

	return res
}

func (agg *Aggregator) GetFinishedTaskByAlertHash(alertHash [32]byte) *FinishedTaskStatus {
	agg.finishedTasksMu.RLock()
	defer agg.finishedTasksMu.RUnlock()

	return agg.finishedTasks[alertHash]
}

// rpc endpoint which is called by operator
// will try to init the task, if currently had a same task for the alert,
// it will return the existing task.
func (agg *Aggregator) CreateTask(req *message.CreateTaskRequest, reply *message.CreateTaskResponse) error {
	agg.logger.Infof("Received signed task response: %#v", req)

	finished := agg.GetFinishedTaskByAlertHash(req.AlertHash)
	if finished != nil {
		return fmt.Errorf("The task %v already finished: %#v", req.AlertHash, finished.TxHash)
	}

	task := agg.GetTaskByAlertHash(req.AlertHash)
	if task == nil {
		agg.logger.Info("create new task", "alert", req.AlertHash)
		taskIndex := agg.newIndex()

		var err error
		task, err = agg.sendNewTask(req.AlertHash, taskIndex)

		if err != nil {
			agg.logger.Error("send new task failed", "err", err)
			return err
		}
	}

	reply.Info = *task

	return nil
}

// rpc endpoint which is called by operator
// reply doesn't need to be checked. If there are no errors, the task response is accepted
// rpc framework forces a reply type to exist, so we put bool as a placeholder
func (agg *Aggregator) ProcessSignedTaskResponse(signedTaskResponse *message.SignedTaskRespRequest, reply *message.SignedTaskRespResponse) error {
	agg.logger.Infof("Received signed task response: %#v", signedTaskResponse)
	taskIndex := signedTaskResponse.Alert.TaskIndex
	taskResponseDigest, err := signedTaskResponse.Alert.SignHash()
	if err != nil {
		return err
	}

	if task := agg.GetTaskByIndex(taskIndex); task == nil {
		agg.logger.Error("ProcessNewSignature error by no task exist", "taskIndex", taskIndex)
		return fmt.Errorf("task not found")
	}

	agg.logger.Infof("ProcessNewSignature: %#v", signedTaskResponse.Alert.TaskIndex)
	err = agg.blsAggregationService.ProcessNewSignature(
		context.Background(), taskIndex, taskResponseDigest,
		&signedTaskResponse.BlsSignature, signedTaskResponse.OperatorId,
	)

	if err != nil {
		agg.logger.Error("ProcessNewSignature error", "err", err)
	}

	return err
}
