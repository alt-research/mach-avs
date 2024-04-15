package generic_operator

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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sourcegraph/jsonrpc2"
)

type CreateSigTaskReq struct {
	AvsName   string                  `json:"avs_name"`
	Method    string                  `json:"method"`
	ParamsRaw hexutil.Bytes           `json:"params_raw"`
	SigHash   alert.HexEncodedBytes32 `json:"sig_hash"`
}

type CreateSigTaskResp struct {
	TaskIndex uint64                  `json:"task_index"`
	TxHash    alert.HexEncodedBytes32 `json:"tx_hash"`
	SigHash   alert.HexEncodedBytes32 `json:"sig_hash"`
}

type GenericRpcServer struct {
	logger             logging.Logger
	server             *http.Server
	serverIpPortAddr   string
	newTaskCreatedChan chan<- GenericRequest
}

func (s *GenericRpcServer) StartServer(ctx context.Context) error {
	s.server = &http.Server{
		Addr: s.serverIpPortAddr,
	}

	s.server.Handler = s.setupHandlers()

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Errorf("failed to start HTTP RPC server: %v", err.Error())
		}
	}()

	return nil
}

func (s *GenericRpcServer) Stop() error {
	s.logger.Info("stop rpc server")
	if s.server == nil {
		s.logger.Warnf("stopping http server that was not initialized")
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}

func (s *GenericRpcServer) setupHandlers() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpRPCHandler)

	return mux
}

func (s *GenericRpcServer) httpRPCHandler(w http.ResponseWriter, r *http.Request) {
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

	switch rpcRequest.Method {
	case "operator_createSigTask":
		{
			var req CreateSigTaskReq
			if err = json.Unmarshal(*rpcRequest.Params, &req); err != nil {
				s.logger.Error("the unmarshal", "err", err)
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, 3, fmt.Errorf("failed to unmarshal req bundle params: %s", err.Error()))
				return
			}

			res := s.CreateSigTask(&req)
			if res.Err != nil {
				s.writeErrorJSON(w, rpcRequest.ID, http.StatusBadRequest, res.Code, fmt.Errorf("failed to call req: %s", res.Err.Error()))
				return
			}

			response := CreateSigTaskResp{
				TaskIndex: uint64(res.TaskIndex),
				TxHash:    res.TxHash,
				SigHash:   req.SigHash,
			}

			s.writeJSON(w, rpcRequest.ID, http.StatusOK, response)
		}
	default:
		err := fmt.Errorf("got unsupported method name: %v", rpcRequest.Method)
		s.writeErrorJSON(w, rpcRequest.ID, http.StatusNotFound, 1, err)
	}
}

func (s *GenericRpcServer) writeErrorJSON(w http.ResponseWriter, id jsonrpc2.ID, statusCode int, code uint32, err error) {
	s.logger.Info("writeErrorJSON", "id", id, "err", err)

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
		s.logger.Errorf("error: failed to marshal json to render an error, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *GenericRpcServer) writeJSON(w http.ResponseWriter, id jsonrpc2.ID, resultHTTPCode int, jsonAnswer interface{}) {
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

func (s *GenericRpcServer) CreateSigTask(req *CreateSigTaskReq) GenericResponse {
	s.logger.Info("CreateSigTask", "req", req)

	responseChan := make(chan GenericResponse, 1)

	s.newTaskCreatedChan <- GenericRequest{
		AVSName:      req.AvsName,
		Method:       req.Method,
		ParamsRaw:    req.ParamsRaw,
		SigHash:      message.Bytes32(req.SigHash),
		ResponseChan: responseChan,
	}

	response := <-responseChan

	if response.Msg != "" {
		s.logger.Error("AlertBlockMismatch failed", "msg", response.Msg)
	}

	return response
}
