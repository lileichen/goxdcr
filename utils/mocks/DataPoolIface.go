// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"

// DataPoolIface is an autogenerated mock type for the DataPoolIface type
type DataPoolIface struct {
	mock.Mock
}

// GetByteSlice provides a mock function with given fields: sizeRequested
func (_m *DataPoolIface) GetByteSlice(sizeRequested uint64) ([]byte, error) {
	ret := _m.Called(sizeRequested)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(uint64) []byte); ok {
		r0 = rf(sizeRequested)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(sizeRequested)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutByteSlice provides a mock function with given fields: doneSlice
func (_m *DataPoolIface) PutByteSlice(doneSlice []byte) {
	_m.Called(doneSlice)
}
