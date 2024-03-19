// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alt-research/avs/core/chainio (interfaces: AvsWriterer)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/avs_writer.go -package=mocks github.com/alt-research/avs/core/chainio AvsWriterer
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	ecdsa "crypto/ecdsa"
	big "math/big"
	reflect "reflect"

	contractRegistryCoordinator "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	bls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
	contractMachServiceManager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "go.uber.org/mock/gomock"
)

// MockAvsWriterer is a mock of AvsWriterer interface.
type MockAvsWriterer struct {
	ctrl     *gomock.Controller
	recorder *MockAvsWritererMockRecorder
}

// MockAvsWritererMockRecorder is the mock recorder for MockAvsWriterer.
type MockAvsWritererMockRecorder struct {
	mock *MockAvsWriterer
}

// NewMockAvsWriterer creates a new mock instance.
func NewMockAvsWriterer(ctrl *gomock.Controller) *MockAvsWriterer {
	mock := &MockAvsWriterer{ctrl: ctrl}
	mock.recorder = &MockAvsWritererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAvsWriterer) EXPECT() *MockAvsWritererMockRecorder {
	return m.recorder
}

// DeregisterOperator mocks base method.
func (m *MockAvsWriterer) DeregisterOperator(arg0 context.Context, arg1 []byte, arg2 contractRegistryCoordinator.BN254G1Point) (*types.Receipt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeregisterOperator", arg0, arg1, arg2)
	ret0, _ := ret[0].(*types.Receipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeregisterOperator indicates an expected call of DeregisterOperator.
func (mr *MockAvsWritererMockRecorder) DeregisterOperator(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeregisterOperator", reflect.TypeOf((*MockAvsWriterer)(nil).DeregisterOperator), arg0, arg1, arg2)
}

// RegisterOperatorInQuorumWithAVSRegistryCoordinator mocks base method.
func (m *MockAvsWriterer) RegisterOperatorInQuorumWithAVSRegistryCoordinator(arg0 context.Context, arg1 *ecdsa.PrivateKey, arg2 [32]byte, arg3 *big.Int, arg4 *bls.KeyPair, arg5 []byte, arg6 string) (*types.Receipt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterOperatorInQuorumWithAVSRegistryCoordinator", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(*types.Receipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterOperatorInQuorumWithAVSRegistryCoordinator indicates an expected call of RegisterOperatorInQuorumWithAVSRegistryCoordinator.
func (mr *MockAvsWritererMockRecorder) RegisterOperatorInQuorumWithAVSRegistryCoordinator(arg0, arg1, arg2, arg3, arg4, arg5, arg6 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterOperatorInQuorumWithAVSRegistryCoordinator", reflect.TypeOf((*MockAvsWriterer)(nil).RegisterOperatorInQuorumWithAVSRegistryCoordinator), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// SendConfirmAlert mocks base method.
func (m *MockAvsWriterer) SendConfirmAlert(arg0 context.Context, arg1 contractMachServiceManager.IMachServiceManagerAlertHeader, arg2 contractMachServiceManager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Receipt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendConfirmAlert", arg0, arg1, arg2)
	ret0, _ := ret[0].(*types.Receipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendConfirmAlert indicates an expected call of SendConfirmAlert.
func (mr *MockAvsWritererMockRecorder) SendConfirmAlert(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendConfirmAlert", reflect.TypeOf((*MockAvsWriterer)(nil).SendConfirmAlert), arg0, arg1, arg2)
}

// UpdateStakesOfEntireOperatorSetForQuorums mocks base method.
func (m *MockAvsWriterer) UpdateStakesOfEntireOperatorSetForQuorums(arg0 context.Context, arg1 [][]common.Address, arg2 []byte) (*types.Receipt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStakesOfEntireOperatorSetForQuorums", arg0, arg1, arg2)
	ret0, _ := ret[0].(*types.Receipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStakesOfEntireOperatorSetForQuorums indicates an expected call of UpdateStakesOfEntireOperatorSetForQuorums.
func (mr *MockAvsWritererMockRecorder) UpdateStakesOfEntireOperatorSetForQuorums(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStakesOfEntireOperatorSetForQuorums", reflect.TypeOf((*MockAvsWriterer)(nil).UpdateStakesOfEntireOperatorSetForQuorums), arg0, arg1, arg2)
}

// UpdateStakesOfOperatorSubsetForAllQuorums mocks base method.
func (m *MockAvsWriterer) UpdateStakesOfOperatorSubsetForAllQuorums(arg0 context.Context, arg1 []common.Address) (*types.Receipt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStakesOfOperatorSubsetForAllQuorums", arg0, arg1)
	ret0, _ := ret[0].(*types.Receipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStakesOfOperatorSubsetForAllQuorums indicates an expected call of UpdateStakesOfOperatorSubsetForAllQuorums.
func (mr *MockAvsWritererMockRecorder) UpdateStakesOfOperatorSubsetForAllQuorums(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStakesOfOperatorSubsetForAllQuorums", reflect.TypeOf((*MockAvsWriterer)(nil).UpdateStakesOfOperatorSubsetForAllQuorums), arg0, arg1)
}
