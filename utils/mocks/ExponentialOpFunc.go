// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ExponentialOpFunc is an autogenerated mock type for the ExponentialOpFunc type
type ExponentialOpFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields:
func (_m *ExponentialOpFunc) Execute() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
