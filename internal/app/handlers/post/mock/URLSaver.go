// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/handlers/post/post.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/avGenie/url-shortener/internal/app/entity"
	model "github.com/avGenie/url-shortener/internal/app/storage/api/model"
	gomock "github.com/golang/mock/gomock"
)

// MockURLSaver is a mock of URLSaver interface.
type MockURLSaver struct {
	ctrl     *gomock.Controller
	recorder *MockURLSaverMockRecorder
}

// MockURLSaverMockRecorder is the mock recorder for MockURLSaver.
type MockURLSaverMockRecorder struct {
	mock *MockURLSaver
}

// NewMockURLSaver creates a new mock instance.
func NewMockURLSaver(ctrl *gomock.Controller) *MockURLSaver {
	mock := &MockURLSaver{ctrl: ctrl}
	mock.recorder = &MockURLSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLSaver) EXPECT() *MockURLSaverMockRecorder {
	return m.recorder
}

// SaveURL mocks base method.
func (m *MockURLSaver) SaveURL(ctx context.Context, key, value entity.URL) entity.URLResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURL", ctx, key, value)
	ret0, _ := ret[0].(entity.URLResponse)
	return ret0
}

// SaveURL indicates an expected call of SaveURL.
func (mr *MockURLSaverMockRecorder) SaveURL(ctx, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURL", reflect.TypeOf((*MockURLSaver)(nil).SaveURL), ctx, key, value)
}

// MockURLBatchSaver is a mock of URLBatchSaver interface.
type MockURLBatchSaver struct {
	ctrl     *gomock.Controller
	recorder *MockURLBatchSaverMockRecorder
}

// MockURLBatchSaverMockRecorder is the mock recorder for MockURLBatchSaver.
type MockURLBatchSaverMockRecorder struct {
	mock *MockURLBatchSaver
}

// NewMockURLBatchSaver creates a new mock instance.
func NewMockURLBatchSaver(ctrl *gomock.Controller) *MockURLBatchSaver {
	mock := &MockURLBatchSaver{ctrl: ctrl}
	mock.recorder = &MockURLBatchSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLBatchSaver) EXPECT() *MockURLBatchSaverMockRecorder {
	return m.recorder
}

// SaveBatchURL mocks base method.
func (m *MockURLBatchSaver) SaveBatchURL(ctx context.Context, batch model.Batch) model.BatchResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBatchURL", ctx, batch)
	ret0, _ := ret[0].(model.BatchResponse)
	return ret0
}

// SaveBatchURL indicates an expected call of SaveBatchURL.
func (mr *MockURLBatchSaverMockRecorder) SaveBatchURL(ctx, batch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBatchURL", reflect.TypeOf((*MockURLBatchSaver)(nil).SaveBatchURL), ctx, batch)
}
