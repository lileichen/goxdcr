package mocks

import mock "github.com/stretchr/testify/mock"

// RouterIface is an autogenerated mock type for the RouterIface type
type RouterIface struct {
	mock.Mock
}

// Start provides a mock function with given fields:
func (_m *RouterIface) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *RouterIface) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
