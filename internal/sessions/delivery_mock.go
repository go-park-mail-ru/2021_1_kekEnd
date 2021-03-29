package sessions
// Code generated by MockGen. DO NOT EDIT.
// Source: delivery.go

// Package mock_sessions is a generated GoMock package.

import (
reflect "reflect"
time "time"

gomock "github.com/golang/mock/gomock"
)

// MockDelivery is a mock of Delivery interface.
type MockDelivery struct {
	ctrl     *gomock.Controller
	recorder *MockDeliveryMockRecorder
}

// MockDeliveryMockRecorder is the mock recorder for MockDelivery.
type MockDeliveryMockRecorder struct {
	mock *MockDelivery
}

// NewMockDelivery creates a new mock instance.
func NewMockDelivery(ctrl *gomock.Controller) *MockDelivery {
	mock := &MockDelivery{ctrl: ctrl}
	mock.recorder = &MockDeliveryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDelivery) EXPECT() *MockDeliveryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockDelivery) Create(userID string, expires time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID, expires)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockDeliveryMockRecorder) Create(userID, expires interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDelivery)(nil).Create), userID, expires)
}

// Delete mocks base method.
func (m *MockDelivery) Delete(sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeliveryMockRecorder) Delete(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDelivery)(nil).Delete), sessionID)
}

// GetUser mocks base method.
func (m *MockDelivery) GetUser(sessionID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", sessionID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockDeliveryMockRecorder) GetUser(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockDelivery)(nil).GetUser), sessionID)
}