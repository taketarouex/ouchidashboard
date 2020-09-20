// Code generated by MockGen. DO NOT EDIT.
// Source: ouchi.go

// Package ouchi is a generated GoMock package.
package ouchi

import (
	gomock "github.com/golang/mock/gomock"
	collector "github.com/tktkc72/ouchidashboard/collector"
	enum "github.com/tktkc72/ouchidashboard/enum"
	reflect "reflect"
	time "time"
)

// MockIOuchi is a mock of IOuchi interface
type MockIOuchi struct {
	ctrl     *gomock.Controller
	recorder *MockIOuchiMockRecorder
}

// MockIOuchiMockRecorder is the mock recorder for MockIOuchi
type MockIOuchiMockRecorder struct {
	mock *MockIOuchi
}

// NewMockIOuchi creates a new mock instance
func NewMockIOuchi(ctrl *gomock.Controller) *MockIOuchi {
	mock := &MockIOuchi{ctrl: ctrl}
	mock.recorder = &MockIOuchiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIOuchi) EXPECT() *MockIOuchiMockRecorder {
	return m.recorder
}

// GetLogs mocks base method
func (m *MockIOuchi) GetLogs(logType enum.LogType, start, end time.Time, opts ...GetOption) ([]Log, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{logType, start, end}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLogs", varargs...)
	ret0, _ := ret[0].([]Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs
func (mr *MockIOuchiMockRecorder) GetLogs(logType, start, end interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{logType, start, end}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockIOuchi)(nil).GetLogs), varargs...)
}

// MocknoRoom is a mock of noRoom interface
type MocknoRoom struct {
	ctrl     *gomock.Controller
	recorder *MocknoRoomMockRecorder
}

// MocknoRoomMockRecorder is the mock recorder for MocknoRoom
type MocknoRoomMockRecorder struct {
	mock *MocknoRoom
}

// NewMocknoRoom creates a new mock instance
func NewMocknoRoom(ctrl *gomock.Controller) *MocknoRoom {
	mock := &MocknoRoom{ctrl: ctrl}
	mock.recorder = &MocknoRoomMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MocknoRoom) EXPECT() *MocknoRoomMockRecorder {
	return m.recorder
}

// noRoom mocks base method
func (m *MocknoRoom) noRoom() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "noRoom")
	ret0, _ := ret[0].(bool)
	return ret0
}

// noRoom indicates an expected call of noRoom
func (mr *MocknoRoomMockRecorder) noRoom() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "noRoom", reflect.TypeOf((*MocknoRoom)(nil).noRoom))
}

// MockIRepository is a mock of IRepository interface
type MockIRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRepositoryMockRecorder
}

// MockIRepositoryMockRecorder is the mock recorder for MockIRepository
type MockIRepositoryMockRecorder struct {
	mock *MockIRepository
}

// NewMockIRepository creates a new mock instance
func NewMockIRepository(ctrl *gomock.Controller) *MockIRepository {
	mock := &MockIRepository{ctrl: ctrl}
	mock.recorder = &MockIRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIRepository) EXPECT() *MockIRepositoryMockRecorder {
	return m.recorder
}

// SourceID mocks base method
func (m *MockIRepository) SourceID() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SourceID")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SourceID indicates an expected call of SourceID
func (mr *MockIRepositoryMockRecorder) SourceID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SourceID", reflect.TypeOf((*MockIRepository)(nil).SourceID))
}

// Add mocks base method
func (m *MockIRepository) Add(arg0 []collector.CollectLog) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockIRepositoryMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockIRepository)(nil).Add), arg0)
}

// Fetch mocks base method
func (m *MockIRepository) Fetch(logType enum.LogType, start, end time.Time, limit int, order enum.Order) ([]Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", logType, start, end, limit, order)
	ret0, _ := ret[0].([]Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch
func (mr *MockIRepositoryMockRecorder) Fetch(logType, start, end, limit, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockIRepository)(nil).Fetch), logType, start, end, limit, order)
}
