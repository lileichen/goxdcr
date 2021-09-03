// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	mock "github.com/stretchr/testify/mock"
)

// CollectionNamespacePoolIface is an autogenerated mock type for the CollectionNamespacePoolIface type
type CollectionNamespacePoolIface struct {
	mock.Mock
}

// Get provides a mock function with given fields:
func (_m *CollectionNamespacePoolIface) Get() *base.CollectionNamespace {
	ret := _m.Called()

	var r0 *base.CollectionNamespace
	if rf, ok := ret.Get(0).(func() *base.CollectionNamespace); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*base.CollectionNamespace)
		}
	}

	return r0
}

// Put provides a mock function with given fields: ns
func (_m *CollectionNamespacePoolIface) Put(ns *base.CollectionNamespace) {
	_m.Called(ns)
}
