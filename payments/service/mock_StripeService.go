// Code generated by mockery v2.12.1. DO NOT EDIT.

package service

import (
	mock "github.com/stretchr/testify/mock"
	stripe "github.com/stripe/stripe-go/v72"

	testing "testing"
)

// MockStripeService is an autogenerated mock type for the StripeService type
type MockStripeService struct {
	mock.Mock
}

// HandleCharge provides a mock function with given fields: amount, token
func (_m *MockStripeService) HandleCharge(amount uint, token string) (*stripe.Charge, error) {
	ret := _m.Called(amount, token)

	var r0 *stripe.Charge
	if rf, ok := ret.Get(0).(func(uint, string) *stripe.Charge); ok {
		r0 = rf(amount, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stripe.Charge)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, string) error); ok {
		r1 = rf(amount, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockStripeService creates a new instance of MockStripeService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockStripeService(t testing.TB) *MockStripeService {
	mock := &MockStripeService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
