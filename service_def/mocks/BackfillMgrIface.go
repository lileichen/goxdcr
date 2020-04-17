package mocks

import mock "github.com/stretchr/testify/mock"

// BackfillMgrIface is an autogenerated mock type for the BackfillMgrIface type
type BackfillMgrIface struct {
	mock.Mock
}

// Start provides a mock function with given fields:
func (_m *BackfillMgrIface) Start() error {
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
func (_m *BackfillMgrIface) Stop() {
	_m.Called()
}
