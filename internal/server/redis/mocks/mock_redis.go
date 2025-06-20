// Code generated by MockGen. DO NOT EDIT.
// Source: redis.go

// Package redis_mocks is a generated GoMock package.
package redis_mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/kliuchnikovv/word-of-yoda/domain"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStore) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStoreMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStore)(nil).Close))
}

// DeleteChallenge mocks base method.
func (m *MockStore) DeleteChallenge(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChallenge", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChallenge indicates an expected call of DeleteChallenge.
func (mr *MockStoreMockRecorder) DeleteChallenge(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChallenge", reflect.TypeOf((*MockStore)(nil).DeleteChallenge), ctx, id)
}

// Exists mocks base method.
func (m *MockStore) Exists(ctx context.Context, key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockStoreMockRecorder) Exists(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockStore)(nil).Exists), ctx, key)
}

// GetChallenge mocks base method.
func (m *MockStore) GetChallenge(ctx context.Context, id string) (*domain.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChallenge", ctx, id)
	ret0, _ := ret[0].(*domain.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChallenge indicates an expected call of GetChallenge.
func (mr *MockStoreMockRecorder) GetChallenge(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChallenge", reflect.TypeOf((*MockStore)(nil).GetChallenge), ctx, id)
}

// GetTTL mocks base method.
func (m *MockStore) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTTL", ctx, key)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTTL indicates an expected call of GetTTL.
func (mr *MockStoreMockRecorder) GetTTL(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTTL", reflect.TypeOf((*MockStore)(nil).GetTTL), ctx, key)
}

// ListChallenges mocks base method.
func (m *MockStore) ListChallenges(ctx context.Context, pattern string, limit int) ([]*domain.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChallenges", ctx, pattern, limit)
	ret0, _ := ret[0].([]*domain.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChallenges indicates an expected call of ListChallenges.
func (mr *MockStoreMockRecorder) ListChallenges(ctx, pattern, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChallenges", reflect.TypeOf((*MockStore)(nil).ListChallenges), ctx, pattern, limit)
}

// Ping mocks base method.
func (m *MockStore) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStoreMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStore)(nil).Ping), ctx)
}

// SaveChallenge mocks base method.
func (m *MockStore) SaveChallenge(ctx context.Context, challenge *domain.Challenge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveChallenge", ctx, challenge)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveChallenge indicates an expected call of SaveChallenge.
func (mr *MockStoreMockRecorder) SaveChallenge(ctx, challenge interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveChallenge", reflect.TypeOf((*MockStore)(nil).SaveChallenge), ctx, challenge)
}

// SetTTL mocks base method.
func (m *MockStore) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTTL", ctx, key, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTTL indicates an expected call of SetTTL.
func (mr *MockStoreMockRecorder) SetTTL(ctx, key, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTTL", reflect.TypeOf((*MockStore)(nil).SetTTL), ctx, key, ttl)
}
