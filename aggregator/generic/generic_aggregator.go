package generic

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/aggregator/rpc"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
)

type GenericAggregator struct {
	logger logging.Logger
	cfg    *config.Config

	jsonRpcServerIpPortAddr string

	serviceCtx    context.Context
	services      map[string]*AVSGenericService
	serviceMu     sync.RWMutex
	wg            sync.WaitGroup
	jsonrpcServer *rpc.JsonRpcServer
}

func (agg *GenericAggregator) getService(name string) *AVSGenericService {
	agg.serviceMu.RLock()
	defer agg.serviceMu.RUnlock()

	res := agg.services[name]

	return res
}

func (agg *GenericAggregator) newAVS(avsConfig *message.GenericAVSConfig) error {
	name := avsConfig.AVSName

	if agg.getService(name) != nil {
		agg.logger.Error("the avs service had already been created", "name", name)
		return fmt.Errorf("the avs service had already been created")
	}

	agg.logger.Info(
		"new avs aggregator service",
		"name", name,
		"contract", avsConfig.AVSContractAddress.Hex(),
	)

	service, err := NewAVSGenericTasksAggregatorService(agg.cfg, avsConfig)
	if err != nil {
		return fmt.Errorf("create avs generic tasks aggregator service failed: %v", err)
	}

	agg.wg.Add(1)
	go func() {
		defer agg.wg.Done()

		err := service.Start(agg.serviceCtx)
		if err != nil {
			agg.logger.Error("service run failed", "name", name, "err", err.Error())
		}

		agg.logger.Info("service stopped", "name", name)
	}()

	agg.serviceMu.Lock()
	defer agg.serviceMu.Unlock()

	agg.services[name] = service

	return nil
}

func (agg *GenericAggregator) Wait() {
	agg.serviceMu.Lock()
	defer agg.serviceMu.Unlock()

	for _, service := range agg.services {
		service.Wait()
	}

	agg.wg.Wait()
}

func (agg *GenericAggregator) Start() error {
	return nil
}

// rpc endpoint which is called by operator
// will init operator, just for keep config valid
func (agg *GenericAggregator) InitOperator(avsName string, req *message.InitOperatorRequest) (*message.InitOperatorResponse, error) {
	return nil, nil
}

// rpc endpoint which is called by operator
// will try to init the task, if currently had a same task for the alert,
// it will return the existing task.
func (agg *GenericAggregator) CreateTask(avsName string, req *message.CreateTaskRequest) (*message.CreateTaskResponse, error) {
	return nil, nil
}

func (agg *GenericAggregator) GetSignedTask(avsName string, req *message.CreateTaskRequest) (*message.CreateTaskResponse, error) {
	return nil, nil
}

// rpc endpoint which is called by operator
// reply doesn't need to be checked. If there are no errors, the task response is accepted
// rpc framework forces a reply type to exist, so we put bool as a placeholder
func (agg *GenericAggregator) ProcessSignedTaskResponse(avsName string, signedTaskResponse *message.SignedTaskRespRequest) (*message.SignedTaskRespResponse, error) {
	return nil, nil
}
