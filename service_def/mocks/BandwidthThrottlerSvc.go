package mocks

import mock "github.com/stretchr/testify/mock"

// BandwidthThrottlerSvc is an autogenerated mock type for the BandwidthThrottlerSvc type
type BandwidthThrottlerSvc struct {
	mock.Mock
}

// Throttle provides a mock function with given fields: numberOfBytes, minNumberOfBytes, numberOfBytesOfFirstItem
func (_m *BandwidthThrottlerSvc) Throttle(numberOfBytes int64, minNumberOfBytes int64, numberOfBytesOfFirstItem int64) (int64, int64) {
	ret := _m.Called(numberOfBytes, minNumberOfBytes, numberOfBytesOfFirstItem)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, int64, int64) int64); ok {
		r0 = rf(numberOfBytes, minNumberOfBytes, numberOfBytesOfFirstItem)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(int64, int64, int64) int64); ok {
		r1 = rf(numberOfBytes, minNumberOfBytes, numberOfBytesOfFirstItem)
	} else {
		r1 = ret.Get(1).(int64)
	}

	return r0, r1
}

// Wait provides a mock function with given fields:
func (_m *BandwidthThrottlerSvc) Wait() {
	_m.Called()
}
