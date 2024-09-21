// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/cache/user.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/cache/user.go -package=cachemocks -destination=./internal/repository/cache/mocks/user.mock.go
//

// Package cachemocks is a generated GoMock package.
package cachemocks

import (
	context "context"
	domain "learn_go/webook/internal/domain"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockUserCache is a mock of UserCache interface.
type MockUserCache struct {
	ctrl     *gomock.Controller
	recorder *MockUserCacheMockRecorder
}

// MockUserCacheMockRecorder is the mock recorder for MockUserCache.
type MockUserCacheMockRecorder struct {
	mock *MockUserCache
}

// NewMockUserCache creates a new mock instance.
func NewMockUserCache(ctrl *gomock.Controller) *MockUserCache {
	mock := &MockUserCache{ctrl: ctrl}
	mock.recorder = &MockUserCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserCache) EXPECT() *MockUserCacheMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserCacheMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserCache)(nil).Get), ctx, id)
}

// Set mocks base method.
func (m *MockUserCache) Set(ctx context.Context, user domain.User, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, user, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockUserCacheMockRecorder) Set(ctx, user, expiration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockUserCache)(nil).Set), ctx, user, expiration)
}
