package message

import (
	"encoding/hex"
	"fmt"
	"log/slog"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/alt-research/avs/legacy/aggregator/types"
	"github.com/alt-research/avs/legacy/api/grpc/aggregator"
	"github.com/alt-research/avs/legacy/core"

	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// make TaskResponseDigests print as hex encoded strings instead of a sequence of bytes
type Bytes32 [32]byte

func (m Bytes32) LogValue() slog.Value {
	return slog.StringValue(m.String())
}

func (m Bytes32) String() string {
	return hex.EncodeToString(m[:])
}

func (m *Bytes32) UnderlyingType() [32]byte {
	return *m
}

// The Alert task Information
type AlertTaskInfo struct {
	AlertHash                  Bytes32
	QuorumNumbers              sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
	TaskIndex                  types.TaskIndex
	ReferenceBlockNumber       uint64
}

func NewAlertTaskInfo(req *aggregator.AlertTaskInfo) (*AlertTaskInfo, error) {
	alertHash := req.GetAlertHash()
	if len(alertHash) != 32 {
		return nil, fmt.Errorf("alertHash len should be 32")
	}

	res := &AlertTaskInfo{
		QuorumNumbers:              core.ConvertQuorumNumbersFromBytes(req.GetQuorumNumbers()),
		QuorumThresholdPercentages: core.ConvertQuorumThresholdPercentagesFromBytes(req.GetQuorumThresholdPercentages()),
		TaskIndex:                  req.GetTaskIndex(),
		ReferenceBlockNumber:       req.GetReferenceBlockNumber(),
	}

	copy(res.AlertHash[:], alertHash[:32])

	return res, nil
}

func (r AlertTaskInfo) ToPbType() *aggregator.AlertTaskInfo {
	return &aggregator.AlertTaskInfo{
		AlertHash:                  r.AlertHash[:],
		QuorumNumbers:              r.QuorumNumbers.UnderlyingType(),
		QuorumThresholdPercentages: r.QuorumThresholdPercentages.UnderlyingType(),
		TaskIndex:                  r.TaskIndex,
		ReferenceBlockNumber:       r.ReferenceBlockNumber,
	}
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
	return csservicemanager.IMachServiceManagerAlertHeader{
		MessageHash:                a.AlertHash,
		QuorumNumbers:              a.QuorumNumbers.UnderlyingType(),
		QuorumThresholdPercentages: a.QuorumThresholdPercentages.UnderlyingType(),
		ReferenceBlockNumber:       uint32(a.ReferenceBlockNumber),
	}
}

// The init operator request
type InitOperatorRequest struct {
	Layer1ChainId              uint32
	ChainId                    uint32
	OperatorId                 sdktypes.OperatorId
	OperatorAddress            common.Address
	OperatorStateRetrieverAddr common.Address
	RegistryCoordinatorAddr    common.Address
}

func NewInitOperatorRequest(req *aggregator.InitOperatorRequest) (*InitOperatorRequest, error) {
	operatorId := req.GetOperatorId()
	if len(operatorId) != 32 {
		return nil, fmt.Errorf("operator ID len should be 32, got %d", len(operatorId))
	}

	if !common.IsHexAddress(req.GetOperatorAddress()) {
		return nil, fmt.Errorf("operatorAddress not a hex address")
	}
	operatorAddress := common.HexToAddress(req.OperatorAddress)

	if !common.IsHexAddress(req.GetOperatorStateRetrieverAddr()) {
		return nil, fmt.Errorf("operatorStateRetrieverAddr not a hex address")
	}
	operatorStateRetrieverAddr := common.HexToAddress(req.GetOperatorStateRetrieverAddr())

	if !common.IsHexAddress(req.GetRegistryCoordinatorAddr()) {
		return nil, fmt.Errorf("registryCoordinatorAddr not a hex address")
	}
	registryCoordinatorAddr := common.HexToAddress(req.GetRegistryCoordinatorAddr())

	res := &InitOperatorRequest{
		Layer1ChainId:              req.GetLayer1ChainId(),
		ChainId:                    req.GetChainId(),
		OperatorAddress:            operatorAddress,
		OperatorStateRetrieverAddr: operatorStateRetrieverAddr,
		RegistryCoordinatorAddr:    registryCoordinatorAddr,
	}

	copy(res.OperatorId[:], operatorId[:32])

	return res, nil
}

// The init operator response
type InitOperatorResponse struct {
	Ok  bool
	Res string
}

func (r InitOperatorResponse) ToPbType() *aggregator.InitOperatorResponse {
	return &aggregator.InitOperatorResponse{
		Ok:     r.Ok,
		Reason: r.Res,
	}
}

// The Alert task create request
type CreateTaskRequest struct {
	AlertHash Bytes32
}

func NewCreateTaskRequest(req *aggregator.CreateTaskRequest) (*CreateTaskRequest, error) {
	alertHash := req.GetAlertHash()
	if len(alertHash) != 32 {
		return nil, fmt.Errorf("alertHash len should be 32")
	}

	res := &CreateTaskRequest{}

	copy(res.AlertHash[:], alertHash[:32])

	return res, nil
}

// The Alert task create response
type CreateTaskResponse struct {
	Info AlertTaskInfo
}

func (r CreateTaskResponse) ToPbType() *aggregator.CreateTaskResponse {
	return &aggregator.CreateTaskResponse{
		Info: r.Info.ToPbType(),
	}
}

type SignedTaskRespRequest struct {
	Alert        AlertTaskInfo
	BlsSignature bls.Signature
	OperatorId   sdktypes.OperatorId
}

func NewSignedTaskRespRequest(req *aggregator.SignedTaskRespRequest) (*SignedTaskRespRequest, error) {
	operatorId := req.GetOperatorId()
	if len(operatorId) != 32 {
		return nil, fmt.Errorf("operator ID len should be 32, got %d", len(operatorId))
	}

	alert, err := NewAlertTaskInfo(req.GetAlert())
	if err != nil {
		return nil, fmt.Errorf("new alert task info failed: %v", err.Error())
	}

	signRaw := req.GetOperatorRequestSignature()
	if len(signRaw) != 64 {
		return nil, fmt.Errorf("operatorRequestSignature len should be 64")
	}

	g1Point := bls.NewZeroG1Point().Deserialize(signRaw)
	sign := bls.Signature{G1Point: g1Point}

	res := &SignedTaskRespRequest{
		Alert:        *alert,
		BlsSignature: sign,
	}

	copy(res.OperatorId[:], operatorId[:32])

	return res, nil
}

type SignedTaskRespResponse struct {
	Reply  bool
	TxHash [32]byte
}

func (r SignedTaskRespResponse) ToPbType() *aggregator.SignedTaskRespResponse {
	return &aggregator.SignedTaskRespResponse{
		Reply:  r.Reply,
		TxHash: r.TxHash[:],
	}
}
