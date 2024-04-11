package generic

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	sdkclients "github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/services/avsregistry"
	blsagg "github.com/Layr-Labs/eigensdk-go/services/bls_aggregation"
	oppubkeysserv "github.com/Layr-Labs/eigensdk-go/services/operatorpubkeys"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"

	"github.com/alt-research/avs/aggregator/types"
	"github.com/alt-research/avs/core"
	"github.com/alt-research/avs/core/chainio"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
	"github.com/ethereum/go-ethereum/common"

	// TODO: use a generic bind for dependencies contracts from eigenlayer
	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
)

const (
	// number of blocks after which a task is considered expired
	// this hardcoded here because it's also hardcoded in the contracts, but should
	// ideally be fetched from the contracts
	taskChallengeWindowBlock = 100
	blockTimeDuration        = 12 * time.Second
)

type OperatorStatus struct {
	LastTime   int64               `json:"lastTime"`
	OperatorId sdktypes.OperatorId `json:"operatorId"`
}

type GenericFinishedTaskStatus struct {
	Info             *message.GenericTaskData
	TxHash           common.Hash
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

type AVSGenericService struct {
	logger logging.Logger
	cfg    *config.Config

	avsReader chainio.AvsReaderer
	avsWriter chainio.AvsWriterer
	ethClient eth.Client

	avsConfig             message.GenericAVSConfig
	tasks                 *AVSGenericTasks
	blsAggregationService blsagg.BlsAggregationService

	wg sync.WaitGroup
}

func NewAVSGenericTasksAggregatorService(c *config.Config, avsConfig message.GenericAVSConfig) (*AVSGenericService, error) {
	avsWriter, err := chainio.BuildAvsWriter(
		c.TxMgr,
		avsConfig.AVSRegistryCoordinatorAddress,
		avsConfig.OperatorStateRetrieverAddress,
		c.EthHttpClient, c.Logger, &avsConfig,
	)
	if err != nil {
		c.Logger.Errorf("Cannot create avsWriter", "err", err)
		return nil, err
	}

	avsReader, err := chainio.BuildAvsReader(
		avsConfig.AVSRegistryCoordinatorAddress,
		avsConfig.OperatorStateRetrieverAddress,
		c.EthHttpClient,
		c.Logger,
	)
	if err != nil {
		c.Logger.Error("Cannot create avsReader", "err", err)
		return nil, err
	}

	chainioConfig := sdkclients.BuildAllConfig{
		EthHttpUrl:                 c.EthHttpRpcUrl,
		EthWsUrl:                   c.EthWsRpcUrl,
		RegistryCoordinatorAddr:    avsConfig.AVSRegistryCoordinatorAddress.Hex(),
		OperatorStateRetrieverAddr: avsConfig.OperatorStateRetrieverAddress.Hex(),
		AvsName:                    avsConfig.AVSName,
		// TODO: split metrics from chainio config, for multiple avs in one aggregator, we should use one metrics
		PromMetricsIpPortAddress: ":9090",
	}
	clients, err := sdkclients.BuildAll(chainioConfig, c.PrivateKey, c.Logger)
	if err != nil {
		c.Logger.Errorf("Cannot create sdk clients", "err", err)
		return nil, err
	}

	operatorPubkeysService := oppubkeysserv.NewOperatorPubkeysServiceInMemory(context.Background(), clients.AvsRegistryChainSubscriber, clients.AvsRegistryChainReader, c.Logger)
	avsRegistryService := avsregistry.NewAvsRegistryServiceChainCaller(avsReader, operatorPubkeysService, c.Logger)
	blsAggregationService := blsagg.NewBlsAggregatorService(avsRegistryService, c.Logger)

	return &AVSGenericService{
		logger:    c.Logger,
		cfg:       c,
		avsReader: avsReader,
		avsWriter: avsWriter,
		ethClient: clients.EthHttpClient,

		avsConfig:             avsConfig,
		tasks:                 NewAVSGenericTasks(c.Logger),
		blsAggregationService: blsAggregationService,
	}, nil

}

func (t *AVSGenericService) Start(ctx context.Context) error {
	t.wg.Add(1)
	defer func() {
		t.logger.Info("Stop AVSGenericTasks aggregator service", "name", t.avsConfig.AVSName)
		t.wg.Done()
	}()

	t.logger.Info("Starting AVSGenericTasks aggregator service", "name", t.avsConfig.AVSName)
	// t.logger.Debug("AVSGenericTasks aggregator details", "config", fmt.Sprintf("%#v", t.avsConfig))

	for {
		select {
		case <-ctx.Done():
			return nil
		case blsAggServiceResp := <-t.blsAggregationService.GetResponseChannel():
			t.logger.Info("Received response from blsAggregationService", "blsAggServiceResp", blsAggServiceResp)
			t.sendAggregatedResponseToContract(blsAggServiceResp)
		}
	}
}

func (t *AVSGenericService) Wait() {
	t.wg.Wait()
}

func (t *AVSGenericService) sendAggregatedResponseToContract(blsAggServiceResp blsagg.BlsAggregationServiceResponse) {
	// TODO: check if blsAggServiceResp contains an err
	if blsAggServiceResp.Err != nil {
		t.logger.Error("BlsAggregationServiceResponse contains an error", "err", blsAggServiceResp.Err)
		// panicing to help with debugging (fail fast), but we shouldn't panic if we run this in production
		panic(blsAggServiceResp.Err)
	}
	nonSignerPubkeys := []csservicemanager.BN254G1Point{}
	for _, nonSignerPubkey := range blsAggServiceResp.NonSignersPubkeysG1 {
		nonSignerPubkeys = append(nonSignerPubkeys, core.ConvertToBN254G1Point(nonSignerPubkey))
	}
	quorumApks := []csservicemanager.BN254G1Point{}
	for _, quorumApk := range blsAggServiceResp.QuorumApksG1 {
		quorumApks = append(quorumApks, core.ConvertToBN254G1Point(quorumApk))
	}
	nonSignerStakesAndSignature := csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerPubkeys:             nonSignerPubkeys,
		QuorumApks:                   quorumApks,
		ApkG2:                        core.ConvertToBN254G2Point(blsAggServiceResp.SignersApkG2),
		Sigma:                        core.ConvertToBN254G1Point(blsAggServiceResp.SignersAggSigG1.G1Point),
		NonSignerQuorumBitmapIndices: blsAggServiceResp.NonSignerQuorumBitmapIndices,
		QuorumApkIndices:             blsAggServiceResp.QuorumApkIndices,
		TotalStakeIndices:            blsAggServiceResp.TotalStakeIndices,
		NonSignerStakeIndices:        blsAggServiceResp.NonSignerStakeIndices,
	}

	t.logger.Info("Threshold reached. Sending aggregated response onchain.",
		"taskIndex", blsAggServiceResp.TaskIndex,
	)

	task := t.tasks.GetTaskByIndex(blsAggServiceResp.TaskIndex)

	err := t.sendToContract(context.Background(), task, nonSignerStakesAndSignature)
	if err != nil {
		t.logger.Error("Aggregator failed to respond to task", "err", err)
	}
}

func (t *AVSGenericService) sendToContract(
	ctx context.Context,
	task *message.GenericTaskData,
	nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) error {
	// aggregator will collect the bls sig, when reach the quorum commit to
	// contract by given method and call params:
	// ```
	//  function methodName(
	//      CallParams calldata callParams,
	//      NonSignerStakesAndSignature memory nonSignerStakesAndSignature
	//  )
	// ```
	//
	// callParams 's bytes is given by operator in create bls sig task,`methodName`
	// is a config for avs, then the avs 's contract need check if callParams is valid,
	// this mode can support our mach AVS and eigenda, and some avs forked from eigenda
	res, err := t.avsWriter.SendGenericConfirm(ctx, task, nonSignerStakesAndSignature)
	if err != nil {
		t.logger.Error("Aggregator failed to respond to task", "err", err)
	}

	if res != nil {
		t.tasks.SetFinishedTask(task.TaskSigHash, &GenericFinishedTaskStatus{
			Info:             task,
			TxHash:           res.TxHash,
			BlockHash:        res.BlockHash,
			BlockNumber:      res.BlockNumber,
			TransactionIndex: res.TransactionIndex,
		})
	} else {
		t.logger.Error("the send confirm alert res is failed by nil return", "hash", task.TaskSigHash)
	}

	return err
}

// rpc endpoint which is called by operator
// will init operator, just for keep config valid
func (agg *AVSGenericService) InitOperator(req *message.InitOperatorDatas) (*message.InitOperatorResponse, error) {
	agg.logger.Infof("Received InitOperator: %#v", req)

	reply := &message.InitOperatorResponse{
		Ok: false,
	}

	if agg.avsConfig.OperatorStateRetrieverAddress != req.OperatorStateRetrieverAddr {
		reply.Res = fmt.Sprintf("OperatorStateRetrieverAddr invaild, expect %s", agg.avsConfig.OperatorStateRetrieverAddress.Hex())
		return reply, nil
	}

	if agg.avsConfig.AVSRegistryCoordinatorAddress != req.RegistryCoordinatorAddr {
		reply.Res = fmt.Sprintf("RegistryCoordinatorAddr invaild, expect %s", agg.avsConfig.OperatorStateRetrieverAddress.Hex())
		return reply, nil
	}

	if agg.cfg.Layer1ChainId != req.Layer1ChainId {
		reply.Res = fmt.Sprintf("Layer1ChainId invaild, expect %d", agg.cfg.Layer1ChainId)
		return reply, nil
	}

	agg.tasks.operatorStatusMu.Lock()
	defer agg.tasks.operatorStatusMu.Unlock()

	agg.tasks.operatorStatus[req.OperatorAddress] = &OperatorStatus{
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
func (agg *AVSGenericService) CreateTask(hash message.Bytes32, method string, params []interface{}) (*message.GenericTaskData, error) {
	agg.logger.Info("Received CreateTask", "sigHash", hash)

	finished := agg.tasks.GetFinishedTaskByAlertHash(hash)
	if finished != nil {
		return nil, fmt.Errorf("the task 0x%x already finished: 0x%x", hash, finished.TxHash)
	}

	task := agg.tasks.GetTaskByHash(hash)
	if task == nil {
		agg.logger.Info("create new task", "alert", hash)
		taskIndex := agg.tasks.newIndex()

		var err error
		task, err = agg.sendNewTask(hash, taskIndex, method, params)

		if err != nil {
			agg.logger.Error("send new task failed", "err", err)
			return nil, err
		}
	}

	return task, nil
}

// rpc endpoint which is called by operator
// reply doesn't need to be checked. If there are no errors, the task response is accepted
// rpc framework forces a reply type to exist, so we put bool as a placeholder
func (agg *AVSGenericService) ProcessSignedTaskResponse(taskData *message.GenericTaskData, blsSignature bls.Signature, operatorId sdktypes.OperatorId) (*message.Bytes32, error) {
	agg.logger.Info(
		"Received signed task response",
		"task", fmt.Sprintf("%#v", taskData),
		"operatorId", hex.EncodeToString(operatorId[:]),
	)

	taskIndex := taskData.TaskIndex

	if task := agg.tasks.GetTaskByIndex(taskIndex); task == nil {
		agg.logger.Error("ProcessNewSignature error by no task exist", "taskIndex", taskIndex)
		return nil, fmt.Errorf("task not found")
	}

	agg.logger.Info("ProcessNewSignature", "index", taskData.TaskIndex, "hash", taskData.TaskSigHash, "method", taskData.CallMethod)
	err := agg.blsAggregationService.ProcessNewSignature(
		context.Background(), taskIndex, taskData.TaskSigHash.UnderlyingType(),
		&blsSignature, operatorId,
	)

	if err != nil {
		agg.logger.Error("ProcessNewSignature error", "err", err)
	}

	return &message.Bytes32{}, err
}

// sendNewTask sends a new task to the task manager contract, and updates the Task dict struct
// with the information of operators opted into quorum 0 at the block of task creation.
func (agg *AVSGenericService) sendNewTask(hash message.Bytes32, taskIndex types.TaskIndex, method string, params []interface{}) (*message.GenericTaskData, error) {
	agg.logger.Info("Aggregator sending new task", "hash", hash, "task", taskIndex)

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
	agg.logger.Info("get quorumNumbers from layer1", "quorumNumbers", fmt.Sprintf("%v", quorumNumbers))

	if len(quorumNumbers) < len(agg.avsConfig.QuorumNumbers) {
		agg.logger.Error("the cfg quorum numbers is larger to the layer1, it will commit failed")
		return nil, fmt.Errorf("the quorum numbers is larger to the layer1 %v, expected %v", agg.avsConfig.QuorumNumbers, quorumNumbers)
	}

	// just use config value
	quorumNumbers = core.ConvertQuorumNumbersFromBytes(agg.avsConfig.QuorumNumbers)

	quorumThresholdPercentages, err := agg.avsReader.GetQuorumThresholdPercentages(context.Background(), uint32(referenceBlockNumber), quorumNumbers)
	if err != nil {
		agg.logger.Error("GetQuorumThresholdPercentages failed", "err", err)
		return nil, err
	}

	agg.logger.Info(
		"quorum datas",
		"numbers", fmt.Sprintf("%v", quorumNumbers.UnderlyingType()),
		"thresholdPercentages", fmt.Sprintf("%v", quorumThresholdPercentages.UnderlyingType()),
	)

	newAlertTask := &message.GenericTaskData{
		TaskIndex:                  taskIndex,
		TaskSigHash:                hash,
		QuorumNumbers:              quorumNumbers,
		QuorumThresholdPercentages: quorumThresholdPercentages,
		CallMethod:                 method,
		CallParams:                 params,
		ReferenceBlockNumber:       referenceBlockNumber,
	}

	agg.tasks.SetNewTask(newAlertTask)

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
