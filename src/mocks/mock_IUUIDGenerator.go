// Code generated by MockGen. DO NOT EDIT.
// Source: microservice_gokit_base/src/domain/utils (interfaces: IUUIDGenerator)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIUUIDGenerator is a mock of IUUIDGenerator interface
type MockIUUIDGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockIUUIDGeneratorMockRecorder
}

// MockIUUIDGeneratorMockRecorder is the mock recorder for MockIUUIDGenerator
type MockIUUIDGeneratorMockRecorder struct {
	mock *MockIUUIDGenerator
}

// NewMockIUUIDGenerator creates a new mock instance
func NewMockIUUIDGenerator(ctrl *gomock.Controller) *MockIUUIDGenerator {
	mock := &MockIUUIDGenerator{ctrl: ctrl}
	mock.recorder = &MockIUUIDGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIUUIDGenerator) EXPECT() *MockIUUIDGeneratorMockRecorder {
	return m.recorder
}

// GenerateID mocks base method
func (m *MockIUUIDGenerator) GenerateID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GenerateID indicates an expected call of GenerateID
func (mr *MockIUUIDGeneratorMockRecorder) GenerateID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateID", reflect.TypeOf((*MockIUUIDGenerator)(nil).GenerateID))
}
