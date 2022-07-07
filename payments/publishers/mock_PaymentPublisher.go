// Code generated by mockery v2.12.1. DO NOT EDIT.

package publishers

import (
	ent "github.com/sentrionic/ecommerce/payments/ent"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// MockPaymentPublisher is an autogenerated mock type for the PaymentPublisher type
type MockPaymentPublisher struct {
	mock.Mock
}

// PublishPaymentCreated provides a mock function with given fields: payment
func (_m *MockPaymentPublisher) PublishPaymentCreated(payment *ent.Payment) {
	_m.Called(payment)
}

// NewMockPaymentPublisher creates a new instance of MockPaymentPublisher. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPaymentPublisher(t testing.TB) *MockPaymentPublisher {
	mock := &MockPaymentPublisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
