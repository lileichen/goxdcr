// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	mock "github.com/stretchr/testify/mock"
)

// BackfillMgrComponentListenerGetter is an autogenerated mock type for the BackfillMgrComponentListenerGetter type
type BackfillMgrComponentListenerGetter struct {
	mock.Mock
}

// GetComponentEventListener provides a mock function with given fields: pipeline
func (_m *BackfillMgrComponentListenerGetter) GetComponentEventListener(pipeline common.Pipeline) (common.ComponentEventListener, error) {
	ret := _m.Called(pipeline)

	var r0 common.ComponentEventListener
	if rf, ok := ret.Get(0).(func(common.Pipeline) common.ComponentEventListener); ok {
		r0 = rf(pipeline)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.ComponentEventListener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline) error); ok {
		r1 = rf(pipeline)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
