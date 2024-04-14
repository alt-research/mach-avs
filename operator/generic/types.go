package generic_operator

import "github.com/alt-research/avs/core/message"

type GenericRequest struct {
	AVSName      string
	Method       string
	Params       []interface{}
	SigHash      message.Bytes32
	ResponseChan chan GenericResponse
}

func (g *GenericRequest) SendRespose(code uint32, err error, msg string) {
	g.ResponseChan <- GenericResponse{
		Code: code,
		Err:  err,
		Msg:  "ProcessNewTaskCreatedLog failed",
	}
}

type GenericResponse struct {
	Code      uint32
	TxHash    [32]byte
	TaskIndex uint32
	Err       error
	Msg       string
}
