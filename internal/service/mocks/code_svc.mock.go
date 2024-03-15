// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/code/code.go
//
// Generated by this command:
//
//	mockgen -source=./internal/service/code/code.go -destination=./internal/service/mocks/code_svc.mock.go -package=svcmocks
//

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockService) Send(ctx context.Context, biz, phone string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, biz, phone)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockServiceMockRecorder) Send(ctx, biz, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockService)(nil).Send), ctx, biz, phone)
}

// Verify mocks base method.
func (m *MockService) Verify(ctx context.Context, biz, phone, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", ctx, biz, phone, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockServiceMockRecorder) Verify(ctx, biz, phone, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockService)(nil).Verify), ctx, biz, phone, code)
}