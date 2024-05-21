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
	"github.com/alt-research/avs/core/message"
	"github.com/sourcegraph/jsonrpc2"
)

type RpcResponse struct {
	TaskIndex uint64                  `json:"task_index"`
	TxHash    alert.HexEncodedBytes32 `json:"tx_hash"`
	AlertHash alert.HexEncodedBytes32 `json:"alert_hash"`
}

type RpcServer struct {
	logger             logging.Logger
	server             *http.Server
	serverIpPortAddr   string
	newTaskCreatedChan chan alert.AlertRequest
}

// NewRpcServer creates a new rpc server then init the server.
func NewRpcServer(logger logging.Logger, addr string, taskChain chan alert.AlertRequest) RpcServer {
	res := RpcServer{
		logger: logger,
		server: &http.Server{
			Addr: addr,
		},
		serverIpPortAddr:   addr,
		newTaskCreatedChan: taskChain,
	}
	res.SetHandler(res.setupHandlers())

	return res
}

func (s *RpcServer) SetHandler(handler http.Handler) {
	s.server.Handler = handler
}

func (s *RpcServer) StartServer(ctx context.Context) error {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Errorf("failed to start HTTP RPC server: %v", err.Error())
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
	mux.HandleFunc("/", s.HttpRPCHandler)

	return mux
}

func (s *RpcServer) HttpRPCHandler(w http.ResponseWriter, r *http.Request) {
	rpcRequest := jsonrpc2.Request{}
	err := json.NewDecoder(r.Body).Decode(&rpcRequest)
	if err != nil {
		s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, 1, err)
		return
	}

	if rpcRequest.Params == nil {
		err := errors.New("failed to unmarshal request.Params for mevBundle from mev-builder, error: EOF")
		s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, 1, err)
		return
	}

	s.HttpRPCHandlerRequest(w, rpcRequest)
}

func (s *RpcServer) HttpRPCHandlerRequest(w http.ResponseWriter, rpcRequest jsonrpc2.Request) {
	s.HttpRPCHandlerRequestByAVS("", w, rpcRequest)
}

func (s *RpcServer) HttpRPCHandlerRequestByAVS(avsName string, w http.ResponseWriter, rpcRequest jsonrpc2.Request) {
	switch rpcRequest.Method {
	case "health_check":
		{
			var msg message.BlockWorkProof
			if err := json.Unmarshal(*rpcRequest.Params, &msg); err != nil {
				s.logger.Error("the unmarshal", "err", err)
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, 3, fmt.Errorf("failed to unmarshal alert bundle params: %s", err.Error()))
				return
			}

			s.logger.Info("health_check", "avs_name", avsName, "msg", msg)

			s.writeJSON(w, rpcRequest.ID, http.StatusOK, true)
		}
	case "alert_blockMismatch":
		{
			var alert alert.AlertBlockMismatch
			if err := json.Unmarshal(*rpcRequest.Params, &alert); err != nil {
				s.logger.Error("the unmarshal", "err", err)
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, 3, fmt.Errorf("failed to unmarshal alert bundle params: %s", err.Error()))
				return
			}
			if alert.AVSName == "" && avsName != "" {
				alert.AVSName = avsName
			}

			res := s.AlertBlockMismatch(&alert)
			if res.Err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, res.Code, fmt.Errorf("failed to call alert: %s", res.Err.Error()))
				return
			}

			response := RpcResponse{
				TaskIndex: uint64(res.TaskIndex),
				TxHash:    res.TxHash,
				AlertHash: alert.MessageHash(),
			}

			s.writeJSON(w, rpcRequest.ID, http.StatusOK, response)
		}
	case "alert_blockOutputOracleMismatch":
		{
			var alert alert.AlertBlockOutputOracleMismatch
			if err := json.Unmarshal(*rpcRequest.Params, &alert); err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, 3, fmt.Errorf("failed to unmarshal alert bundle params: %s", err.Error()))
				return
			}
			if alert.AVSName == "" && avsName != "" {
				alert.AVSName = avsName
			}

			res := s.AlertBlockOutputOracleMismatch(&alert)
			if res.Err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, res.Code, fmt.Errorf("failed to call alert output oracle: %s", res.Err.Error()))
				return
			}

			response := RpcResponse{
				TaskIndex: uint64(res.TaskIndex),
				TxHash:    res.TxHash,
				AlertHash: alert.MessageHash(),
			}

			s.writeJSON(w, rpcRequest.ID, http.StatusOK, response)
		}
	case "alert_blockHash":
		{
			var alert alert.AlertBlockHashMismatch
			if err := json.Unmarshal(*rpcRequest.Params, &alert); err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, 3, fmt.Errorf("failed to unmarshal alert bundle params: %s", err.Error()))
				return
			}
			if alert.AVSName == "" && avsName != "" {
				alert.AVSName = avsName
			}

			res := s.AlertBlockHashMismatch(&alert)
			if res.Err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, res.Code, fmt.Errorf("failed to call alert block hash: %s", res.Err.Error()))
				return
			}

			response := RpcResponse{
				TaskIndex: uint64(res.TaskIndex),
				TxHash:    res.TxHash,
				AlertHash: alert.MessageHash(),
			}

			s.writeJSON(w, rpcRequest.ID, http.StatusOK, response)
		}
	default:
		err := fmt.Errorf("got unsupported method name: %v", rpcRequest.Method)
		s.writeErrorJSON(w, rpcRequest.ID, http.StatusNotFound, 1, err)
	}
}

func (s *RpcServer) writeErrorJSON(w http.ResponseWriter, id jsonrpc2.ID, statusCode int, code uint32, err error) {
	WriteErrorJSON(s.logger, w, id, statusCode, code, err)
}

func (s *RpcServer) writeJSON(w http.ResponseWriter, id jsonrpc2.ID, resultHTTPCode int, jsonAnswer interface{}) {
	WriteJSON(s.logger, w, id, resultHTTPCode, jsonAnswer)
}

func (s *RpcServer) AlertBlockMismatch(alertReq *alert.AlertBlockMismatch) alert.AlertResponse {
	s.logger.Info("AlertBlockMismatch", "alert", alertReq)

	responseChan := make(chan alert.AlertResponse)

	s.newTaskCreatedChan <- alert.AlertRequest{
		Alert:   alertReq,
		ResChan: responseChan,
	}

	response := <-responseChan

	if response.Msg != "" {
		s.logger.Error("AlertBlockMismatch failed", "msg", response.Msg)
	}

	return response
}

func (s *RpcServer) AlertBlockOutputOracleMismatch(alertReq *alert.AlertBlockOutputOracleMismatch) alert.AlertResponse {
	s.logger.Info("AlertBlockOutputOracleMismatch", "alert", alertReq)

	responseChan := make(chan alert.AlertResponse)

	s.newTaskCreatedChan <- alert.AlertRequest{
		Alert:   alertReq,
		ResChan: responseChan,
	}

	response := <-responseChan

	if response.Msg != "" {
		s.logger.Error("AlertBlockOutputOracleMismatch failed", "msg", response.Msg)
	}

	return response
}

func (s *RpcServer) AlertBlockHashMismatch(alertReq *alert.AlertBlockHashMismatch) alert.AlertResponse {
	s.logger.Info("AlertBlockHashMismatch", "alert", alertReq)

	responseChan := make(chan alert.AlertResponse)

	s.newTaskCreatedChan <- alert.AlertRequest{
		Alert:   alertReq,
		ResChan: responseChan,
	}

	response := <-responseChan

	if response.Msg != "" {
		s.logger.Error("AlertBlockHashMismatch failed", "msg", response.Msg)
	}

	if response.Err != nil {
		s.logger.Error("AlertBlockHashMismatch failed by err", "err", response.Err)
	}

	return response
}

func WriteErrorJSON(logger logging.Logger, w http.ResponseWriter, id jsonrpc2.ID, statusCode int, code uint32, err error) {
	logger.Info("writeErrorJSON", "id", id, "err", err)

	jsonrpcErr := jsonrpc2.Error{
		Code:    int64(code),
		Message: err.Error(),
	}

	resp := jsonrpc2.Response{
		ID:    id,
		Error: &jsonrpcErr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func WriteJSON(logger logging.Logger, w http.ResponseWriter, id jsonrpc2.ID, resultHTTPCode int, jsonAnswer interface{}) {
	resp := &jsonrpc2.Response{
		ID: id,
	}
	if err := resp.SetResult(jsonAnswer); err != nil {
		logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resultHTTPCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
