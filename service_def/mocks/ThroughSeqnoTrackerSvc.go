// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	common "github.com/couchbase/goxdcr/common"

	mock "github.com/stretchr/testify/mock"
)

// ThroughSeqnoTrackerSvc is an autogenerated mock type for the ThroughSeqnoTrackerSvc type
type ThroughSeqnoTrackerSvc struct {
	mock.Mock
}

// Attach provides a mock function with given fields: pipeline
func (_m *ThroughSeqnoTrackerSvc) Attach(pipeline common.Pipeline) error {
	ret := _m.Called(pipeline)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(pipeline)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetThroughSeqno provides a mock function with given fields: vbno
func (_m *ThroughSeqnoTrackerSvc) GetThroughSeqno(vbno uint16) uint64 {
	ret := _m.Called(vbno)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(uint16) uint64); ok {
		r0 = rf(vbno)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// GetThroughSeqnos provides a mock function with given fields:
func (_m *ThroughSeqnoTrackerSvc) GetThroughSeqnos() map[uint16]uint64 {
	ret := _m.Called()

	var r0 map[uint16]uint64
	if rf, ok := ret.Get(0).(func() map[uint16]uint64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uint16]uint64)
		}
	}

	return r0
}

// GetThroughSeqnosAndManifestIds provides a mock function with given fields:
func (_m *ThroughSeqnoTrackerSvc) GetThroughSeqnosAndManifestIds() (map[uint16]uint64, map[uint16]uint64, map[uint16]uint64) {
	ret := _m.Called()

	var r0 map[uint16]uint64
	if rf, ok := ret.Get(0).(func() map[uint16]uint64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uint16]uint64)
		}
	}

	var r1 map[uint16]uint64
	if rf, ok := ret.Get(1).(func() map[uint16]uint64); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[uint16]uint64)
		}
	}

	var r2 map[uint16]uint64
	if rf, ok := ret.Get(2).(func() map[uint16]uint64); ok {
		r2 = rf()
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(map[uint16]uint64)
		}
	}

	return r0, r1, r2
}

// PrintStatusSummary provides a mock function with given fields:
func (_m *ThroughSeqnoTrackerSvc) PrintStatusSummary() {
	_m.Called()
}

// SetStartSeqno provides a mock function with given fields: vbno, seqno, manifestIds
func (_m *ThroughSeqnoTrackerSvc) SetStartSeqno(vbno uint16, seqno uint64, manifestIds base.CollectionsManifestIdPair) {
	_m.Called(vbno, seqno, manifestIds)
}
