package operator

import (
	"context"
	"net/http"
	"net/rpc"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/core/alert"
)

type RpcServer struct {
	logger             logging.Logger
	serverIpPortAddr   string
	newTaskCreatedChan chan alert.Alert
}

func (s *RpcServer) startServer(ctx context.Context) error {
	err := rpc.Register(s)
	if err != nil {
		s.logger.Fatal("Format of service TaskManager isn't correct. ", "err", err)
	}
	rpc.HandleHTTP()
	err = http.ListenAndServe(s.serverIpPortAddr, nil)
	if err != nil {
		s.logger.Fatal("ListenAndServe", "err", err)
	}

	return nil
}

func (s *RpcServer) AlertBlockMismatch(alert *alert.AlertBlockMismatch, reply *bool) error {
	s.logger.Info("AlertBlockMismatch", "alert", alert)

	s.newTaskCreatedChan <- alert

	return nil
}

func (s *RpcServer) AlertBlockOutputOracleMismatch(alert *alert.AlertBlockOutputOracleMismatch, reply *bool) error {
	s.logger.Info("AlertBlockOutputOracleMismatch", "alert", alert)

	s.newTaskCreatedChan <- alert

	return nil
}
