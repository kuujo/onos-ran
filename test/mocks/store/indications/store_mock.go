// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/onosproject/onos-ric/pkg/store/indications (interfaces: Store)

// Package mock_indications_store is a generated GoMock package.
package mock_indications_store

import (
	gomock "github.com/golang/mock/gomock"
	indications "github.com/onosproject/onos-ric/pkg/store/indications"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Clear mocks base method
func (m *MockStore) Clear() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear")
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear
func (mr *MockStoreMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockStore)(nil).Clear))
}

// Close mocks base method
func (m *MockStore) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockStoreMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStore)(nil).Close))
}

// Discard mocks base method
func (m *MockStore) Discard(arg0 indications.ID, arg1 ...indications.DiscardOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Discard", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Discard indicates an expected call of Discard
func (mr *MockStoreMockRecorder) Discard(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discard", reflect.TypeOf((*MockStore)(nil).Discard), varargs...)
}

// List mocks base method
func (m *MockStore) List(arg0 chan<- indications.Indication) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// List indicates an expected call of List
func (mr *MockStoreMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockStore)(nil).List), arg0)
}

// Lookup mocks base method
func (m *MockStore) Lookup(arg0 indications.ID, arg1 ...indications.LookupOption) (*indications.Indication, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Lookup", varargs...)
	ret0, _ := ret[0].(*indications.Indication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lookup indicates an expected call of Lookup
func (mr *MockStoreMockRecorder) Lookup(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lookup", reflect.TypeOf((*MockStore)(nil).Lookup), varargs...)
}

// Record mocks base method
func (m *MockStore) Record(arg0 *indications.Indication, arg1 ...indications.RecordOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Record", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Record indicates an expected call of Record
func (mr *MockStoreMockRecorder) Record(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Record", reflect.TypeOf((*MockStore)(nil).Record), varargs...)
}

// Watch mocks base method
func (m *MockStore) Watch(arg0 chan<- indications.Event, arg1 ...indications.WatchOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Watch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Watch indicates an expected call of Watch
func (mr *MockStoreMockRecorder) Watch(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockStore)(nil).Watch), varargs...)
}
