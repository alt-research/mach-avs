package rpc

import (
	"errors"

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

type AggregatorRpcHandler interface {
	InitOperator(req *message.InitOperatorRequest) (*message.InitOperatorResponse, error)
	CreateTask(req *message.CreateTaskRequest) (*message.CreateTaskResponse, error)
	ProcessSignedTaskResponse(signedTaskResponse *message.SignedTaskRespRequest) (*message.SignedTaskRespResponse, error)
}
