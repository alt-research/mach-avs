package message

import (
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/alt-research/avs/aggregator/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type GenericTaskType int

const (
	GenericTaskTypeByHash GenericTaskType = iota
)

// Register a generic task
type RegisterGenericTaskData struct {
	QuorumNumbers                 []uint8
	QuorumThresholdPercentages    []uint8
	AVSRegistryCoordinatorAddress string
	OperatorStateRetrieverAddress string
	AVSContractAddress            string
}

// Register a generic task
type GenericAVSConfig struct {
	AVSName                       string
	QuorumNumbers                 []uint8
	QuorumThresholdPercentages    []uint8
	AVSRegistryCoordinatorAddress string
	OperatorStateRetrieverAddress string
	AVSContractAddress            common.Address
	Abi                           abi.ABI
}

// Register a generic task
type GenericTaskData struct {
	TaskIndex                  types.TaskIndex
	TaskSigHash                Bytes32
	QuorumNumbers              sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
	CallMethod                 string
	CallParams                 []interface{}
	ReferenceBlockNumber       uint64
}

type InitOperatorDatas struct {
	AVSName                    string              `json:"avs_name"`
	Layer1ChainId              uint32              `json:"layer1_chain_id"`
	OperatorId                 sdktypes.OperatorId `json:"operator_id"`
	OperatorAddress            common.Address      `json:"operator_address"`
	OperatorStateRetrieverAddr common.Address      `json:"operator_state_retriever_addr"`
	RegistryCoordinatorAddr    common.Address      `json:"registry_coordinator_addr"`
}
