package aggregator

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"

	blsagg "github.com/Layr-Labs/eigensdk-go/services/bls_aggregation"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"

	"github.com/alt-research/avs/aggregator/rpc"
	"github.com/alt-research/avs/core"
	"github.com/alt-research/avs/core/chainio"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"

	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
)

const (
	// number of blocks after which a task is considered expired
	// this hardcoded here because it's also hardcoded in the contracts, but should
	// ideally be fetched from the contracts
	taskChallengeWindowBlock = 100
	blockTimeDuration        = 12 * time.Second
	avsName                  = "mach"
)

type FinishedTaskStatus struct {
	Message          *message.AlertTaskInfo
	TxHash           common.Hash
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

type OperatorStatus struct {
	LastTime   int64               `json:"lastTime"`
	OperatorId sdktypes.OperatorId `json:"operatorId"`
}

// Aggregator sends tasks (numbers to square) onchain, then listens for operator signed TaskResponses.
// It aggregates responses signatures, and if any of the TaskResponses reaches the QuorumThresholdPercentage for each quorum
// (currently we only use a single quorum of the ERC20Mock token), it sends the aggregated TaskResponse and signature onchain.
//
// The signature is checked in the BLSSignatureChecker.sol contract, which expects a
//
//	struct NonSignerStakesAndSignature {
//		uint32[] nonSignerQuorumBitmapIndices;
//		BN254.G1Point[] nonSignerPubkeys;
//		BN254.G1Point[] quorumApks;
//		BN254.G2Point apkG2;
//		BN254.G1Point sigma;
//		uint32[] quorumApkIndices;
//		uint32[] totalStakeIndices;
//		uint32[][] nonSignerStakeIndices; // nonSignerStakeIndices[quorumNumberIndex][nonSignerIndex]
//	}
//
// A task can only be responded onchain by having enough operators sign on it such that their stake in each quorum reaches the QuorumThresholdPercentage.
// In order to verify this onchain, the Registry contracts store the history of stakes and aggregate pubkeys (apks) for each operators and each quorum. These are
// updated everytime an operator registers or deregisters with the BLSRegistryCoordinatorWithIndices.sol contract, or calls UpdateStakes() on the StakeRegistry.sol contract,
// after having received new delegated shares or having delegated shares removed by stakers queuing withdrawals. Each of these pushes to their respective datatype array a new entry.
//
// This is true for quorumBitmaps (represent the quorums each operator is opted into), quorumApks (apks per quorum), totalStakes (total stake per quorum), and nonSignerStakes (stake per quorum per operator).
// The 4 "indices" in NonSignerStakesAndSignature basically represent the index at which to fetch their respective data, given a blockNumber at which the task was created.
// Note that different data types might have different indices, since for eg QuorumBitmaps are updated for operators registering/deregistering, but not for UpdateStakes.
// Thankfully, we have deployed a helper contract BLSOperatorStateRetriever.sol whose function getCheckSignaturesIndices() can be used to fetch the indices given a block number.
//
// The 4 other fields nonSignerPubkeys, quorumApks, apkG2, and sigma, however, must be computed individually.
// apkG2 and sigma are just the aggregated signature and pubkeys of the operators who signed the task response (aggregated over all quorums, so individual signatures might be duplicated).
// quorumApks are the G1 aggregated pubkeys of the operators who signed the task response, but one per quorum, as opposed to apkG2 which is summed over all quorums.
// nonSignerPubkeys are the G1 pubkeys of the operators who did not sign the task response, but were opted into the quorum at the blocknumber at which the task was created.
// Upon sending a task onchain (or receiving a NewTaskCreated Event if the tasks were sent by an external task generator), the aggregator can get the list of all operators opted into each quorum at that
// block number by calling the getOperatorState() function of the BLSOperatorStateRetriever.sol contract.
type Aggregator struct {
	logger logging.Logger

	serverIpPortAddr        string
	grpcServerIpPortAddr    string
	jsonRpcServerIpPortAddr string

	avsWriter chainio.AvsWriterer

	service       *AggregatorService
	legacyRpc     *rpc.LegacyRpcHandler
	gRpc          *rpc.GRpcHandler
	jsonrpcServer *rpc.JsonRpcServer
}

// NewAggregator creates a new Aggregator with the provided config.
func NewAggregator(c *config.Config) (*Aggregator, error) {
	if c.MachAVSCfg == nil {
		return nil, fmt.Errorf(
			"If not use env `AVS_REGISTRY_COORDINATOR_ADDRESS` and `OPERATOR_STATE_RETRIEVER_ADDRESS`, should use --avs-deployment to use config for avs contract addresses!",
		)
	}

	avsWriter, err := chainio.BuildAvsWriter(c.TxMgr, c.MachAVSCfg.RegistryCoordinatorAddr, c.MachAVSCfg.OperatorStateRetrieverAddr, c.EthHttpClient, c.Logger, nil)
	if err != nil {
		c.Logger.Errorf("Cannot create avsWriter", "err", err)
		return nil, err
	}

	service, err := NewAggregatorService(c)
	if err != nil {
		c.Logger.Errorf("Cannot create NewAggregatorService", "err", err)
		return nil, err
	}

	legacyRpc := rpc.NewLegacyRpcHandler(c.Logger, service)

	var grpcServer *rpc.GRpcHandler
	if c.AggregatorGRPCServerIpPortAddr != "" {
		c.Logger.Infof("Create grpc server in %s", c.AggregatorGRPCServerIpPortAddr)
		grpcServer = rpc.NewGRpcHandler(c.Logger, service)
	}

	var jsonrpcServer *rpc.JsonRpcServer
	if c.AggregatorJSONRPCServerIpPortAddr != "" {
		c.Logger.Infof("Create json rpc server in %s", c.AggregatorJSONRPCServerIpPortAddr)
		jsonrpcServer = rpc.NewJsonRpcServer(c.Logger, service, c.RpcVhosts, c.RpcCors)
	}

	return &Aggregator{
		logger:                  c.Logger,
		serverIpPortAddr:        c.AggregatorServerIpPortAddr,
		grpcServerIpPortAddr:    c.AggregatorGRPCServerIpPortAddr,
		jsonRpcServerIpPortAddr: c.AggregatorJSONRPCServerIpPortAddr,
		avsWriter:               avsWriter,
		service:                 service,
		legacyRpc:               legacyRpc,
		gRpc:                    grpcServer,
		jsonrpcServer:           jsonrpcServer,
	}, nil
}

func (agg *Aggregator) Start(ctx context.Context, wg *sync.WaitGroup) error {
	defer func() {
		agg.wait()
		wg.Done()
	}()

	agg.logger.Infof("Starting aggregator.")
	agg.logger.Infof("Starting aggregator rpc server.")

	agg.startRpcServer(ctx)

	agg.logger.Info("Aggregator Rpc Server Started.")
	for {
		select {
		case <-ctx.Done():
			return nil
		case blsAggServiceResp := <-agg.service.GetResponseChannel():
			agg.logger.Info("Received response from blsAggregationService", "blsAggServiceResp", blsAggServiceResp)
			agg.sendAggregatedResponseToContract(blsAggServiceResp)
		}
	}
}

func (agg *Aggregator) startRpcServer(ctx context.Context) {
	go agg.legacyRpc.StartServer(ctx, 1*time.Second, agg.serverIpPortAddr)

	if agg.gRpc != nil {
		go agg.gRpc.StartServer(ctx, agg.grpcServerIpPortAddr)
	}

	if agg.jsonrpcServer != nil {
		go agg.jsonrpcServer.StartServer(ctx, agg.jsonRpcServerIpPortAddr)
	}
}

func (agg *Aggregator) wait() {
	agg.logger.Info("The aggregator is wait to exit")

	if agg.legacyRpc != nil {
		agg.legacyRpc.Wait()
	}

	if agg.gRpc != nil {
		agg.gRpc.Wait()
	}

	if agg.jsonrpcServer != nil {
		agg.jsonrpcServer.Wait()
	}

	agg.logger.Info("The aggregator is exited")
}

func (agg *Aggregator) sendAggregatedResponseToContract(blsAggServiceResp blsagg.BlsAggregationServiceResponse) {
	// TODO: check if blsAggServiceResp contains an err
	if blsAggServiceResp.Err != nil {
		agg.logger.Error("BlsAggregationServiceResponse contains an error", "err", blsAggServiceResp.Err)
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

	agg.logger.Info("Threshold reached. Sending aggregated response onchain.",
		"taskIndex", blsAggServiceResp.TaskIndex,
	)

	task := agg.service.GetTaskByIndex(blsAggServiceResp.TaskIndex)

	res, err := agg.avsWriter.SendConfirmAlert(context.Background(), task, nonSignerStakesAndSignature)
	if err != nil {
		agg.logger.Error("Aggregator failed to respond to task", "err", err)
	}

	if res != nil {
		agg.service.SetFinishedTask(task.AlertHash, &FinishedTaskStatus{
			Message:          task,
			TxHash:           res.TxHash,
			BlockHash:        res.BlockHash,
			BlockNumber:      res.BlockNumber,
			TransactionIndex: res.TransactionIndex,
		})
	} else {
		agg.logger.Error("the send confirm alert res is failed by nil return", "hash", task.AlertHash)
	}
}
