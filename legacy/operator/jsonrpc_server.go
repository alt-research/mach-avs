package operator

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/alt-research/avs/legacy/core/alert"
	"github.com/alt-research/avs/legacy/core/message"
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
	newWorkProofChan   chan message.HealthCheckMsg
}

// NewRpcServer creates a new rpc server then init the server.
func NewRpcServer(
	logger logging.Logger,
	addr string,
	taskChain chan alert.AlertRequest,
	newWorkProofChan chan message.HealthCheckMsg) RpcServer {
	res := RpcServer{
		logger: logger,
		server: &http.Server{
			Addr: addr,
		},
		serverIpPortAddr:   addr,
		newTaskCreatedChan: taskChain,
		newWorkProofChan:   newWorkProofChan,
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
		WriteErrorJSON(s.logger, w, rpcRequest.ID, http.StatusBadRequest, 1, err)
		return
	}

	if rpcRequest.Params == nil {
		err := errors.New("failed to unmarshal request.Params for mevBundle from mev-builder, error: EOF")
		WriteErrorJSON(s.logger, w, rpcRequest.ID, http.StatusBadRequest, 1, err)
		return
	}

	s.HttpRPCHandlerRequest(w, rpcRequest)
}

func (s *RpcServer) HttpRPCHandlerRequest(w http.ResponseWriter, rpcRequest jsonrpc2.Request) {
	s.HttpRPCHandlerRequestByAVS("", w, rpcRequest)
}

func unmarshalParams(logger logging.Logger, rpcRequest jsonrpc2.Request, val any) error {
	var params []json.RawMessage
	logger.Debug("params", "raw", *rpcRequest.Params)
	if err := json.Unmarshal(*rpcRequest.Params, &params); err != nil {
		logger.Error("the unmarshal", "err", err)
		return errors.Wrap(err, "failed to unmarshal alert bundle params")
	}

	if len(params) == 0 {
		logger.Error("failed to unmarshal params by no msg")
		return errors.New("failed to unmarshal health check params by no msg")
	}

	if err := json.Unmarshal(params[0], val); err != nil {
		logger.Error("the unmarshal", "err", err)
		return errors.Wrap(err, "failed to unmarshal alert bundle params")
	}

	return nil
}

func (s *RpcServer) HttpRPCHandlerRequestByAVS(avsName string, w http.ResponseWriter, rpcRequest jsonrpc2.Request) {
	logger := s.logger.With(
		"avsName", avsName,
		"rpcMethod", rpcRequest.Method,
	)

	switch rpcRequest.Method {
	case "health_check":
		{
			var msg message.BlockWorkProof
			if err := unmarshalParams(logger, rpcRequest, &msg); err != nil {
				WriteErrorJSON(logger, w, rpcRequest.ID, http.StatusBadRequest, 3, err)
				return
			}

			s.logger.Info("on health check", "msg", msg)

			s.newWorkProofChan <- message.HealthCheckMsg{
				AvsName: avsName,
				Proof:   msg,
			}

			WriteJSON(logger, w, rpcRequest.ID, true)
		}
	case "alert_blockMismatch":
		{
			var alert alert.AlertBlockMismatch
			if err := unmarshalParams(logger, rpcRequest, &alert); err != nil {
				WriteErrorJSON(logger, w, rpcRequest.ID, http.StatusBadRequest, 3, err)
				return
			}

			if alert.AVSName == "" && avsName != "" {
				alert.AVSName = avsName
			}

			res := s.SendAlert(logger.With("alertType", "AlertBlockMismatch"), &alert)
			if res.Err != nil {
				WriteErrorJSON(
					logger,
					w, rpcRequest.ID,
					http.StatusBadRequest, res.Code,
					errors.Wrap(res.Err, "failed to call alert"),
				)
				return
			}

			response := RpcResponse{
				TaskIndex: uint64(res.TaskIndex),
				TxHash:    res.TxHash,
				AlertHash: alert.MessageHash(),
			}

			WriteJSON(logger, w, rpcRequest.ID, response)
		}
	case "alert_blockOutputOracleMismatch":
		{
			var alert alert.AlertBlockOutputOracleMismatch
			if err := unmarshalParams(logger, rpcRequest, &alert); err != nil {
				WriteErrorJSON(logger, w, rpcRequest.ID, http.StatusBadRequest, 3, err)
				return
			}

			if alert.AVSName == "" && avsName != "" {
				alert.AVSName = avsName
			}

			res := s.SendAlert(logger.With("alertType", "AlertBlockOutputOracleMismatch"), &alert)
			if res.Err != nil {
				WriteErrorJSON(
					logger,
					w, rpcRequest.ID, http.StatusBadRequest, res.Code,
					errors.Wrap(res.Err, "failed to call alert output oracle"),
				)
				return
			}

			response := RpcResponse{
				TaskIndex: uint64(res.TaskIndex),
				TxHash:    res.TxHash,
				AlertHash: alert.MessageHash(),
			}

			WriteJSON(logger, w, rpcRequest.ID, response)
		}
	case "alert_blockHash":
		{
			var alert alert.AlertBlockHashMismatch
			if err := unmarshalParams(logger, rpcRequest, &alert); err != nil {
				WriteErrorJSON(logger, w, rpcRequest.ID, http.StatusBadRequest, 3, err)
				return
			}

			if alert.AVSName == "" && avsName != "" {
				alert.AVSName = avsName
			}

			res := s.SendAlert(logger.With("alertType", "AlertBlockHashMismatch"), &alert)
			if res.Err != nil {
				WriteErrorJSON(
					logger,
					w, rpcRequest.ID, http.StatusBadRequest,
					res.Code, errors.Wrap(res.Err, "failed to call alert block hash"),
				)
				return
			}

			response := RpcResponse{
				TaskIndex: uint64(res.TaskIndex),
				TxHash:    res.TxHash,
				AlertHash: alert.MessageHash(),
			}

			WriteJSON(logger, w, rpcRequest.ID, response)
		}
	default:
		err := errors.Errorf("got unsupported method name: %v", rpcRequest.Method)
		WriteErrorJSON(logger, w, rpcRequest.ID, http.StatusNotFound, 1, err)
	}
}

func (s *RpcServer) SendAlert(logger logging.Logger, alertReq alert.Alert) alert.AlertResponse {
	logger.Info("Send alert", "alert", alertReq)

	responseChan := make(chan alert.AlertResponse)

	s.newTaskCreatedChan <- alert.AlertRequest{
		Alert:   alertReq,
		ResChan: responseChan,
	}

	response := <-responseChan

	if response.Msg != "" {
		logger.Error("Send alert failed", "msg", response.Msg)
	}

	if response.Err != nil {
		logger.Error("Send alert failed by err", "err", response.Err)
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

func WriteJSON(logger logging.Logger, w http.ResponseWriter, id jsonrpc2.ID, jsonAnswer interface{}) {
	resp := &jsonrpc2.Response{
		ID: id,
	}
	if err := resp.SetResult(jsonAnswer); err != nil {
		logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
