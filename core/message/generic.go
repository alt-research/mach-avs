package message

import (
	"fmt"

	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/alt-research/avs/aggregator/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

type CreateGenericTaskResponse struct {
	// The hash of alert
	TaskSigHash hexutil.Bytes `json:"sig_hash"`
	// QuorumNumbers of task
	QuorumNumbers []uint8 `json:"quorum_numbers"`
	// QuorumThresholdPercentages of task
	QuorumThresholdPercentages []uint8 `json:"quorum_threshold_percentages"`
	// TaskIndex
	TaskIndex uint32 `json:"task_index"`
	// ReferenceBlockNumber
	ReferenceBlockNumber uint64 `json:"reference_block_number"`
}

func (c *CreateGenericTaskResponse) SigHash() ([32]byte, error) {
	var res [32]byte
	if len(c.TaskSigHash) != 32 {
		return res, fmt.Errorf("task sig hash len not 32")
	}

	copy(res[:], c.TaskSigHash[:32])

	return res, nil
}
