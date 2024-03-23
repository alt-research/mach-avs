package message

import (
	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"

	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/alt-research/avs/aggregator/types"
)

// The Alert task Information
type AlertTaskInfo struct {
	AlertHash                  [32]byte
	QuorumNumbers              sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
	TaskIndex                  types.TaskIndex
	ReferenceBlockNumber       uint64
}

func (a *AlertTaskInfo) EncodeSigHash() ([]byte, error) {
	// The order here has to match the field ordering of ReducedBatchHeader defined in IEigenDAServiceManager.sol
	// ref: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L43
	AlertType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "messageHash",
			Type: "bytes32",
		},
		{
			Name: "referenceBlockNumber",
			Type: "uint32",
		},
	})
	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{
		{
			Type: AlertType,
		},
	}

	s := struct {
		MessageHash          [32]byte
		ReferenceBlockNumber uint32
	}{
		MessageHash:          a.AlertHash,
		ReferenceBlockNumber: uint32(a.ReferenceBlockNumber),
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (a AlertTaskInfo) SignHash() ([32]byte, error) {
	alertBytes, err := a.EncodeSigHash()
	if err != nil {
		return [32]byte{}, err
	}

	var hash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(alertBytes)
	copy(hash[:], hasher.Sum(nil)[:32])

	return hash, nil
}

func (a AlertTaskInfo) ToIMachServiceManagerAlertHeader() csservicemanager.IMachServiceManagerAlertHeader {
	quorumNumbers := make([]byte, len(a.QuorumNumbers))
	quorumThresholdPercentages := make([]byte, len(a.QuorumThresholdPercentages))

	for i, _ := range a.QuorumNumbers {
		quorumNumbers[i] = byte(a.QuorumNumbers[i])
	}

	for i, _ := range a.QuorumThresholdPercentages {
		quorumThresholdPercentages[i] = byte(a.QuorumThresholdPercentages[i])
	}

	return csservicemanager.IMachServiceManagerAlertHeader{
		MessageHash:                a.AlertHash,
		QuorumNumbers:              quorumNumbers,
		QuorumThresholdPercentages: quorumThresholdPercentages,
		ReferenceBlockNumber:       uint32(a.ReferenceBlockNumber),
	}
}

// The Alert task create request
type CreateTaskRequest struct {
	AlertHash [32]byte
}

// The Alert task create response
type CreateTaskResponse struct {
	Info AlertTaskInfo
}
