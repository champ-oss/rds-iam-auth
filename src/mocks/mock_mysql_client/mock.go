// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/mysql_client/mysql_client.go

// Package mock_mysql_client is a generated GoMock package.
package mock_mysql_client

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMysqlClientInterface is a mock of MysqlClientInterface interface.
type MockMysqlClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMysqlClientInterfaceMockRecorder
}

// MockMysqlClientInterfaceMockRecorder is the mock recorder for MockMysqlClientInterface.
type MockMysqlClientInterfaceMockRecorder struct {
	mock *MockMysqlClientInterface
}

// NewMockMysqlClientInterface creates a new mock instance.
func NewMockMysqlClientInterface(ctrl *gomock.Controller) *MockMysqlClientInterface {
	mock := &MockMysqlClientInterface{ctrl: ctrl}
	mock.recorder = &MockMysqlClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMysqlClientInterface) EXPECT() *MockMysqlClientInterfaceMockRecorder {
	return m.recorder
}

// CloseDb mocks base method.
func (m *MockMysqlClientInterface) CloseDb() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseDb")
}

// CloseDb indicates an expected call of CloseDb.
func (mr *MockMysqlClientInterfaceMockRecorder) CloseDb() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseDb", reflect.TypeOf((*MockMysqlClientInterface)(nil).CloseDb))
}

// Query mocks base method.
func (m *MockMysqlClientInterface) Query(sql string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", sql)
	ret0, _ := ret[0].(error)
	return ret0
}

// Query indicates an expected call of Query.
func (mr *MockMysqlClientInterfaceMockRecorder) Query(sql interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockMysqlClientInterface)(nil).Query), sql)
}
