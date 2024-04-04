package rpc

import (
	"context"
	"net/http"
	"net/rpc"
	"sync"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/alt-research/avs/core/message"
)

type LegacyRpcHandler struct {
	logger     logging.Logger
	aggreagtor AggregatorRpcHandler
	wg         *sync.WaitGroup
}

func NewLegacyRpcHandler(logger logging.Logger, aggreagtor AggregatorRpcHandler) *LegacyRpcHandler {
	return &LegacyRpcHandler{
		logger:     logger,
		aggreagtor: aggreagtor,
		wg:         &sync.WaitGroup{},
	}
}

func (agg *LegacyRpcHandler) Wait() {
	agg.wg.Wait()
}

func (agg *LegacyRpcHandler) StartServer(ctx context.Context, shutdownTimeout time.Duration, serverIpPortAddr string) {
	agg.logger.Info("Start LegacyRpcServer", "addr", serverIpPortAddr)

	err := rpc.RegisterName("Aggregator", agg)
	if err != nil {
		agg.logger.Fatal("Format of service TaskManager isn't correct. ", "err", err)
	}

	rpc.HandleHTTP()

	serverErr := make(chan error, 1)
	server := &http.Server{
		Addr:    serverIpPortAddr,
		Handler: nil,
	}

	agg.wg.Add(1)
	go func() {
		defer agg.wg.Done()
		serverErr <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		agg.logger.Info("Stop LegacyRpcServer by Done")

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err = server.Shutdown(ctx)
	case err = <-serverErr:
	}

	if err != nil && err != http.ErrAbortHandler {
		agg.logger.Error("LegacyRpcServer serve stopped by error", "err", err)
	} else {
		agg.logger.Info("LegacyRpcServer serve stopped")
	}
}

// rpc endpoint which is called by operator
// will init operator, just for keep config valid
func (agg *LegacyRpcHandler) InitOperator(req *message.InitOperatorRequest, reply *message.InitOperatorResponse) error {
	res, err := agg.aggreagtor.InitOperator(req)
	if err != nil {
		return err
	}

	*reply = *res
	return nil
}

// rpc endpoint which is called by operator
// will try to init the task, if currently had a same task for the alert,
// it will return the existing task.
func (agg *LegacyRpcHandler) CreateTask(req *message.CreateTaskRequest, reply *message.CreateTaskResponse) error {
	res, err := agg.aggreagtor.CreateTask(req)
	if err != nil {
		return err
	}

	*reply = *res
	return nil
}

// rpc endpoint which is called by operator
// reply doesn't need to be checked. If there are no errors, the task response is accepted
// rpc framework forces a reply type to exist, so we put bool as a placeholder
func (agg *LegacyRpcHandler) ProcessSignedTaskResponse(req *message.SignedTaskRespRequest, reply *message.SignedTaskRespResponse) error {
	res, err := agg.aggreagtor.ProcessSignedTaskResponse(req)
	if err != nil {
		return err
	}

	*reply = *res
	return nil
}
