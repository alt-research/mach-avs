package generic

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	aggregatorCmdInitOperator = iota + 1
	aggregatorCmdCreateTask
	aggregatorCmdGetSignedTask
	aggregatorCmdProcessSignedTaskResponse
)

type aggregatorCmd struct {
	cmdType           int
	avsName           string
	initOperatorDatas *message.InitOperatorDatas
	taskData          *message.GenericTaskData
	blsSignature      bls.Signature
	operatorId        sdktypes.OperatorId
	resp              chan<- aggregatorCmdRes
}

type aggregatorCmdRes struct {
	err      error
	ok       bool
	reason   string
	taskData *message.GenericTaskData
	txHash   common.Hash
}

type GenericAggregator struct {
	logger     logging.Logger
	cfg        *config.Config
	avsConfigs []*message.GenericAVSConfig

	jsonRpcServerIpPortAddr string

	serviceCtx context.Context
	services   map[string]*AVSGenericService
	wg         sync.WaitGroup
	cmdChan    chan aggregatorCmd
}

func NewAVSGenericService(c *config.Config, avsConfigs []*message.GenericAVSConfig) (*GenericAggregator, error) {
	return &GenericAggregator{
		logger:                  c.Logger,
		cfg:                     c,
		avsConfigs:              avsConfigs,
		jsonRpcServerIpPortAddr: c.AggregatorJSONRPCServerIpPortAddr,
		services:                make(map[string]*AVSGenericService, 32),
		cmdChan:                 make(chan aggregatorCmd, 128),
	}, nil
}

func (agg *GenericAggregator) getService(name string) *AVSGenericService {
	res := agg.services[name]

	return res
}

func (agg *GenericAggregator) newAVS(avsConfig message.GenericAVSConfig) error {
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

	agg.services[name] = service

	return nil
}

func (agg *GenericAggregator) initAVSServices() error {
	agg.logger.Info("init AVS services", "len", len(agg.avsConfigs))

	for _, cfg := range agg.avsConfigs {
		agg.logger.Info(
			"init AVS service",
			"name", cfg.AVSName,
			"address", cfg.AVSContractAddress,
			"registryCoordinator", cfg.AVSRegistryCoordinatorAddress,
			"operatorStateRetriever", cfg.OperatorStateRetrieverAddress,
		)

		err := agg.newAVS(*cfg)
		if err != nil {
			return fmt.Errorf("init avs service %s failed: %v", cfg.AVSName, err)
		}
	}

	return nil
}

func (agg *GenericAggregator) Start(ctx context.Context) error {
	agg.logger.Info("Start generic aggregator")
	agg.serviceCtx = ctx

	err := agg.initAVSServices()
	if err != nil {
		return errors.Wrap(err, "initAVSServices")
	}

	agg.wg.Add(1)
	defer func() {
		agg.logger.Info("Stop GenericAggregator")
		agg.wg.Done()
	}()

	agg.logger.Info("Start generic aggregator handler cmds")

	for {
		select {
		case <-ctx.Done():
			return nil
		case cmd := <-agg.cmdChan:
			agg.handlerCmd(cmd)
		}
	}
}

func (agg *GenericAggregator) handlerCmd(cmd aggregatorCmd) {
	name := cmd.avsName

	agg.logger.Debug("handler cmd", "name", name, "type", cmd.cmdType)

	switch cmd.cmdType {
	case aggregatorCmdInitOperator:
		{
			resp, err := agg.initOperator(name, cmd.initOperatorDatas)
			if err != nil {
				cmd.resp <- aggregatorCmdRes{
					err: err,
				}
			}
			cmd.resp <- aggregatorCmdRes{
				ok:     resp.Ok,
				reason: resp.Res,
			}
		}
	case aggregatorCmdCreateTask:
		{
			resp, err := agg.createTask(name, cmd.taskData.TaskSigHash, cmd.taskData.CallMethod, cmd.taskData.CallParams)
			if err != nil {
				cmd.resp <- aggregatorCmdRes{
					err: err,
				}
			}
			cmd.resp <- aggregatorCmdRes{
				taskData: resp,
			}
		}
	case aggregatorCmdGetSignedTask:
		{
			_, err := agg.getSignedTask(name, cmd.taskData.TaskSigHash, cmd.taskData.CallMethod)
			if err != nil {
				cmd.resp <- aggregatorCmdRes{
					err: err,
				}
			}
			cmd.resp <- aggregatorCmdRes{}
		}
	case aggregatorCmdProcessSignedTaskResponse:
		{
			resp, err := agg.processSignedTaskResponse(name, cmd.taskData, cmd.blsSignature, cmd.operatorId)
			if err != nil {
				cmd.resp <- aggregatorCmdRes{
					err: err,
				}
			}
			cmd.resp <- aggregatorCmdRes{
				txHash: resp.UnderlyingType(),
			}
		}
	}
}

func (agg *GenericAggregator) Wait() {
	for _, service := range agg.services {
		service.Wait()
	}

	agg.wg.Wait()
}

func (agg *GenericAggregator) initOperator(avsName string, req *message.InitOperatorDatas) (*message.InitOperatorResponse, error) {
	service := agg.getService(avsName)
	if service == nil {
		agg.logger.Error("not found service for init operator", "name", avsName)
	}

	agg.logger.Info("init operator service", "name", avsName, "operator", req.OperatorAddress.Hex())

	rsp, err := service.InitOperator(req)
	if err != nil {
		agg.logger.Error("init operator service failed", "name", avsName, "err", err.Error())
		return nil, err
	}

	return rsp, nil
}

func (agg *GenericAggregator) createTask(
	avsName string,
	hash message.Bytes32,
	method string,
	params []interface{}) (*message.GenericTaskData, error) {
	service := agg.getService(avsName)
	if service == nil {
		agg.logger.Error("not found service for create task", "name", avsName)
	}

	agg.logger.Info("create task service", "name", avsName, "operator", hash.String(), "method", method)

	rsp, err := service.CreateTask(hash, method, params)
	if err != nil {
		agg.logger.Error("create task service failed", "name", avsName, "err", err.Error())
		return nil, err
	}

	return rsp, nil
}

func (agg *GenericAggregator) getSignedTask(
	avsName string,
	hash message.Bytes32,
	method string) (*message.CreateTaskResponse, error) {
	return nil, nil
}

func (agg *GenericAggregator) processSignedTaskResponse(
	avsName string,
	taskData *message.GenericTaskData,
	blsSignature bls.Signature,
	operatorId sdktypes.OperatorId) (*message.Bytes32, error) {
	service := agg.getService(avsName)
	if service == nil {
		agg.logger.Error("not found service for process signed task", "name", avsName)
	}

	agg.logger.Info("process signed task service",
		"name", avsName,
		"hash", taskData.TaskSigHash.String(),
		"index", taskData.TaskIndex,
		"method", taskData.CallMethod)

	rsp, err := service.ProcessSignedTaskResponse(taskData, blsSignature, operatorId)
	if err != nil {
		agg.logger.Error("process signed task service failed", "name", avsName, "err", err.Error())
		return nil, err
	}

	return rsp, nil
}

func (agg *GenericAggregator) InitOperator(avsName string, req *message.InitOperatorDatas) (*message.InitOperatorResponse, error) {
	resChan := make(chan aggregatorCmdRes, 1)

	agg.cmdChan <- aggregatorCmd{
		cmdType:           aggregatorCmdInitOperator,
		avsName:           avsName,
		initOperatorDatas: req,
	}

	res := <-resChan
	if res.err != nil {
		return nil, res.err
	}

	return &message.InitOperatorResponse{
		Ok:  res.ok,
		Res: res.reason,
	}, nil
}

func (agg *GenericAggregator) CreateTask(
	avsName string,
	hash message.Bytes32,
	method string,
	params []interface{}) (*message.GenericTaskData, error) {
	resChan := make(chan aggregatorCmdRes, 1)

	agg.cmdChan <- aggregatorCmd{
		cmdType: aggregatorCmdInitOperator,
		avsName: avsName,
		taskData: &message.GenericTaskData{
			TaskSigHash: hash,
			CallMethod:  method,
			CallParams:  params,
		},
	}

	res := <-resChan
	if res.err != nil {
		return nil, res.err
	}

	return res.taskData, nil
}

func (agg *GenericAggregator) GetSignedTask(
	avsName string,
	hash message.Bytes32,
	method string) (*message.CreateTaskResponse, error) {
	return nil, nil
}

func (agg *GenericAggregator) ProcessSignedTaskResponse(
	avsName string,
	taskData *message.GenericTaskData,
	blsSignature bls.Signature,
	operatorId sdktypes.OperatorId) (*common.Hash, error) {
	resChan := make(chan aggregatorCmdRes, 1)

	agg.cmdChan <- aggregatorCmd{
		cmdType:      aggregatorCmdInitOperator,
		avsName:      avsName,
		taskData:     taskData,
		blsSignature: blsSignature,
		operatorId:   operatorId,
	}

	res := <-resChan
	if res.err != nil {
		return nil, res.err
	}

	return &res.txHash, nil
}
