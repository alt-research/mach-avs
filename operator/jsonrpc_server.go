package operator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/core/alert"
	"github.com/sourcegraph/jsonrpc2"
)

type RpcResponse struct {
	result bool
}

type RpcServer struct {
	logger             logging.Logger
	server             *http.Server
	serverIpPortAddr   string
	newTaskCreatedChan chan alert.Alert
}

func (s *RpcServer) StartServer(ctx context.Context) error {
	s.server = &http.Server{
		Addr: s.serverIpPortAddr,
	}

	s.server.Handler = s.setupHandlers()

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Errorf("failed to start HTTP RPC server: %v", err)
		}
	}()

	return nil
}

func (s *RpcServer) Stop() error {
	s.logger.Info("stop rpc server")
	if s.server == nil {
		s.logger.Warnf("stopping http server that was not initialized")
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}

func (s *RpcServer) setupHandlers() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpRPCHandler)

	return mux
}

func (s *RpcServer) httpRPCHandler(w http.ResponseWriter, r *http.Request) {
	rpcRequest := jsonrpc2.Request{}
	err := json.NewDecoder(r.Body).Decode(&rpcRequest)
	if err != nil {
		s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, err)
		return
	}

	if rpcRequest.Params == nil {
		err := errors.New("failed to unmarshal request.Params for mevBundle from mev-builder, error: EOF")
		s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, err)
		return
	}

	switch rpcRequest.Method {
	case "alert_blockMismatch":
		{
			var alert alert.AlertBlockMismatch
			if err = json.Unmarshal(*rpcRequest.Params, &alert); err != nil {
				s.logger.Error("the unmarshal", "err", err)
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, fmt.Errorf("failed to unmarshal alert bundle params: %v", err))
				return
			}

			if err = s.AlertBlockMismatch(&alert); err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, fmt.Errorf("failed to call alert: %v", err))
				return
			}

			s.writeJSON(w, rpcRequest.ID, http.StatusOK, RpcResponse{
				result: true,
			})
		}
	case "alert_blockOutputOracleMismatch":
		{
			var alert alert.AlertBlockOutputOracleMismatch
			if err = json.Unmarshal(*rpcRequest.Params, &alert); err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, fmt.Errorf("failed to unmarshal alert bundle params: %v", err))
				return
			}

			if err = s.AlertBlockOutputOracleMismatch(&alert); err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, fmt.Errorf("failed to call alert output oracle: %v", err))
				return
			}

			s.writeJSON(w, rpcRequest.ID, http.StatusOK, RpcResponse{
				result: true,
			})
		}
	case "alert_blockHash":
		{
			var alert alert.AlertBlockHashMismatch
			if err = json.Unmarshal(*rpcRequest.Params, &alert); err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, fmt.Errorf("failed to unmarshal alert bundle params: %v", err))
				return
			}

			if err = s.AlertBlockHashMismatch(&alert); err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, fmt.Errorf("failed to call alert output oracle: %v", err))
				return
			}

			s.writeJSON(w, rpcRequest.ID, http.StatusOK, RpcResponse{
				result: true,
			})
		}
	default:
		err := fmt.Errorf("got unsupported method name: %v", rpcRequest.Method)
		s.writeErrorJSON(w, rpcRequest.ID, http.StatusNotFound, err)
	}
}

func (s *RpcServer) writeErrorJSON(w http.ResponseWriter, id jsonrpc2.ID, statusCode int, err error) {
	s.logger.Info("writeErrorJSON", "id", id, "err", err)

	jsonrpcErr := jsonrpc2.Error{}
	jsonrpcErr.SetError(err)

	resp := jsonrpc2.Response{
		ID:    id,
		Error: &jsonrpcErr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *RpcServer) writeJSON(w http.ResponseWriter, id jsonrpc2.ID, resultHTTPCode int, jsonAnswer interface{}) {
	resp := &jsonrpc2.Response{
		ID: id,
	}
	if err := resp.SetResult(jsonAnswer); err != nil {
		s.logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resultHTTPCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *RpcServer) AlertBlockMismatch(alert *alert.AlertBlockMismatch) error {
	s.logger.Info("AlertBlockMismatch", "alert", alert)

	s.newTaskCreatedChan <- alert

	return nil
}

func (s *RpcServer) AlertBlockOutputOracleMismatch(alert *alert.AlertBlockOutputOracleMismatch) error {
	s.logger.Info("AlertBlockOutputOracleMismatch", "alert", alert)

	s.newTaskCreatedChan <- alert

	return nil
}

func (s *RpcServer) AlertBlockHashMismatch(alert *alert.AlertBlockHashMismatch) error {
	s.logger.Info("AlertBlockHashMismatch", "alert", alert)

	s.newTaskCreatedChan <- alert

	return nil
}
