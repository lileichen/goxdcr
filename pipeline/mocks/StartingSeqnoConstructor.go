// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	mock "github.com/stretchr/testify/mock"
)

// StartingSeqnoConstructor is an autogenerated mock type for the StartingSeqnoConstructor type
type StartingSeqnoConstructor struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *StartingSeqnoConstructor) Execute(_a0 common.Pipeline) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
