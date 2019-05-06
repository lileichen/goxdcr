package mocks

import memcached "github.com/couchbase/gomemcached/client"
import mock "github.com/stretchr/testify/mock"

// FilterIface is an autogenerated mock type for the FilterIface type
type FilterIface struct {
	mock.Mock
}

// FilterUprEvent provides a mock function with given fields: uprEvent
func (_m *FilterIface) FilterUprEvent(uprEvent *memcached.UprEvent) bool {
	ret := _m.Called(uprEvent)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*memcached.UprEvent) bool); ok {
		r0 = rf(uprEvent)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
