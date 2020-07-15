package mocks

import (
	net "net"

	mock "github.com/stretchr/testify/mock"
)

// DnsSrvHelperIface is an autogenerated mock type for the DnsSrvHelperIface type
type DnsSrvHelperIface struct {
	mock.Mock
}

// DnsSrvLookup provides a mock function with given fields: hostname
func (_m *DnsSrvHelperIface) DnsSrvLookup(hostname string) ([]*net.SRV, error) {
	ret := _m.Called(hostname)

	var r0 []*net.SRV
	if rf, ok := ret.Get(0).(func(string) []*net.SRV); ok {
		r0 = rf(hostname)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*net.SRV)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(hostname)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}