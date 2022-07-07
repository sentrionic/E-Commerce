// Code generated by mockery v2.12.1. DO NOT EDIT.

package listeners

import (
	context "context"
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// MockProductListener is an autogenerated mock type for the ProductListener type
type MockProductListener struct {
	mock.Mock
}

// ProductCreatedListener provides a mock function with given fields: ctx
func (_m *MockProductListener) ProductCreatedListener(ctx context.Context) {
	_m.Called(ctx)
}

// NewMockProductListener creates a new instance of MockProductListener. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockProductListener(t testing.TB) *MockProductListener {
	mock := &MockProductListener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
