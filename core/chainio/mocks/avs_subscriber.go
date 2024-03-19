// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alt-research/avs/core/chainio (interfaces: AvsSubscriberer)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/avs_subscriber.go -package=mocks github.com/alt-research/avs/core/chainio AvsSubscriberer
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	contractMachServiceManager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
	event "github.com/ethereum/go-ethereum/event"
	gomock "go.uber.org/mock/gomock"
)

// MockAvsSubscriberer is a mock of AvsSubscriberer interface.
type MockAvsSubscriberer struct {
	ctrl     *gomock.Controller
	recorder *MockAvsSubscribererMockRecorder
}

// MockAvsSubscribererMockRecorder is the mock recorder for MockAvsSubscriberer.
type MockAvsSubscribererMockRecorder struct {
	mock *MockAvsSubscriberer
}

// NewMockAvsSubscriberer creates a new mock instance.
func NewMockAvsSubscriberer(ctrl *gomock.Controller) *MockAvsSubscriberer {
	mock := &MockAvsSubscriberer{ctrl: ctrl}
	mock.recorder = &MockAvsSubscribererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAvsSubscriberer) EXPECT() *MockAvsSubscribererMockRecorder {
	return m.recorder
}

// SubscribeToAlertConfirmed mocks base method.
func (m *MockAvsSubscriberer) SubscribeToAlertConfirmed(arg0 chan *contractMachServiceManager.ContractMachServiceManagerAlertConfirmed) event.Subscription {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToAlertConfirmed", arg0)
	ret0, _ := ret[0].(event.Subscription)
	return ret0
}

// SubscribeToAlertConfirmed indicates an expected call of SubscribeToAlertConfirmed.
func (mr *MockAvsSubscribererMockRecorder) SubscribeToAlertConfirmed(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToAlertConfirmed", reflect.TypeOf((*MockAvsSubscriberer)(nil).SubscribeToAlertConfirmed), arg0)
}
