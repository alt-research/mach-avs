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
