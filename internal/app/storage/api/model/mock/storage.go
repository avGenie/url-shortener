// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/storage/api/model/storage.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/avGenie/url-shortener/internal/app/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() entity.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(entity.Response)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// GetURL mocks base method.
func (m *MockStorage) GetURL(ctx context.Context, key entity.URL) entity.URLResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", ctx, key)
	ret0, _ := ret[0].(entity.URLResponse)
	return ret0
}

// GetURL indicates an expected call of GetURL.
func (mr *MockStorageMockRecorder) GetURL(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockStorage)(nil).GetURL), ctx, key)
}

// PingServer mocks base method.
func (m *MockStorage) PingServer(ctx context.Context) entity.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingServer", ctx)
	ret0, _ := ret[0].(entity.Response)
	return ret0
}

// PingServer indicates an expected call of PingServer.
func (mr *MockStorageMockRecorder) PingServer(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingServer", reflect.TypeOf((*MockStorage)(nil).PingServer), ctx)
}

// SaveURL mocks base method.
func (m *MockStorage) SaveURL(ctx context.Context, key, value entity.URL) entity.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURL", ctx, key, value)
	ret0, _ := ret[0].(entity.Response)
	return ret0
}

// SaveURL indicates an expected call of SaveURL.
func (mr *MockStorageMockRecorder) SaveURL(ctx, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURL", reflect.TypeOf((*MockStorage)(nil).SaveURL), ctx, key, value)
}