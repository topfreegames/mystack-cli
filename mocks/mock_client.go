// Automatically generated by MockGen. DO NOT EDIT!
// Source: client_interface.go

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of ClientInterface interface
type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockClientInterfaceRecorder
}

// Recorder for MockClientInterface (not exported)
type _MockClientInterfaceRecorder struct {
	mock *MockClientInterface
}

func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &_MockClientInterfaceRecorder{mock}
	return mock
}

func (_m *MockClientInterface) EXPECT() *_MockClientInterfaceRecorder {
	return _m.recorder
}

func (_m *MockClientInterface) Get(url string, host string) ([]byte, int, error) {
	ret := _m.ctrl.Call(_m, "Get", url, host)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockClientInterfaceRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0, arg1)
}

func (_m *MockClientInterface) GetToStdOut(url string, host string) error {
	ret := _m.ctrl.Call(_m, "GetToStdOut", url, host)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientInterfaceRecorder) GetToStdOut(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetToStdOut", arg0, arg1)
}

func (_m *MockClientInterface) Put(url string, body map[string]interface{}) ([]byte, int, error) {
	ret := _m.ctrl.Call(_m, "Put", url, body)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockClientInterfaceRecorder) Put(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Put", arg0, arg1)
}

func (_m *MockClientInterface) Delete(url string) ([]byte, int, error) {
	ret := _m.ctrl.Call(_m, "Delete", url)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockClientInterfaceRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}
