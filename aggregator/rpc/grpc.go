package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/api/grpc/aggregator"
	"github.com/alt-research/avs/core/message"
	"google.golang.org/grpc"
)

type GRpcHandler struct {
	aggregator.UnimplementedAggregatorServer
	logger     logging.Logger
	aggreagtor AggregatorRpcHandler
	wg         *sync.WaitGroup
}

func NewGRpcHandler(logger logging.Logger, aggreagtor AggregatorRpcHandler) *GRpcHandler {
	return &GRpcHandler{
		logger:     logger,
		aggreagtor: aggreagtor,
		wg:         &sync.WaitGroup{},
	}
}

func (s *GRpcHandler) StartServer(ctx context.Context, serverIpPortAddr string) {
	s.logger.Info("Start GRpcServer", "addr", serverIpPortAddr)

	lis, err := net.Listen("tcp", serverIpPortAddr)
	if err != nil {
		s.logger.Fatalf("GRpcServer failed to listen: %v", err)
	}

	server := grpc.NewServer()
	aggregator.RegisterAggregatorServer(server, s)

	serverErr := make(chan error, 1)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		serverErr <- server.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Stop GRpcServer by Done")
		server.Stop()
	case err = <-serverErr:
	}

	if err != nil {
		s.logger.Error("GRpcServer serve stopped by error", "err", err)
	} else {
		s.logger.Info("GRpcServer serve stopped")
	}
}

func (s *GRpcHandler) Wait() {
	s.wg.Wait()
}

// Send Init operator to aggregator from operator, will check if the config is matching
func (s *GRpcHandler) InitOperator(ctx context.Context, req *aggregator.InitOperatorRequest) (*aggregator.InitOperatorResponse, error) {
	msg, err := message.NewInitOperatorRequest(req)
	if err != nil {
		return nil, fmt.Errorf("initOperator message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.InitOperator(msg)
	if err != nil {
		return nil, fmt.Errorf("initOperator handler error: %v", err.Error())
	}

	return resp.ToPbType(), nil
}

// Create a alert task
func (s *GRpcHandler) CreateTask(ctx context.Context, req *aggregator.CreateTaskRequest) (*aggregator.CreateTaskResponse, error) {
	msg, err := message.NewCreateTaskRequest(req)
	if err != nil {
		return nil, fmt.Errorf("createTask message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.CreateTask(msg)
	if err != nil {
		return nil, fmt.Errorf("createTask handler error: %v", err.Error())
	}

	return resp.ToPbType(), nil
}

// Send signed task for alert
func (s *GRpcHandler) ProcessSignedTaskResponse(ctx context.Context, req *aggregator.SignedTaskRespRequest) (*aggregator.SignedTaskRespResponse, error) {
	msg, err := message.NewSignedTaskRespRequest(req)
	if err != nil {
		return nil, fmt.Errorf("processSignedTaskResponse message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.ProcessSignedTaskResponse(msg)
	if err != nil {
		return nil, fmt.Errorf("processSignedTaskResponse handler error: %v", err.Error())
	}

	return resp.ToPbType(), nil
}
