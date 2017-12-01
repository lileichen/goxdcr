package mocks

import metadata "github.com/couchbase/goxdcr/metadata"
import mock "github.com/stretchr/testify/mock"

// ReplicationSettingsSvc is an autogenerated mock type for the ReplicationSettingsSvc type
type ReplicationSettingsSvc struct {
	mock.Mock
}

// GetDefaultReplicationSettings provides a mock function with given fields:
func (_m *ReplicationSettingsSvc) GetDefaultReplicationSettings() (*metadata.ReplicationSettings, error) {
	ret := _m.Called()

	var r0 *metadata.ReplicationSettings
	if rf, ok := ret.Get(0).(func() *metadata.ReplicationSettings); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.ReplicationSettings)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetDefaultReplicationSettings provides a mock function with given fields: _a0
func (_m *ReplicationSettingsSvc) SetDefaultReplicationSettings(_a0 *metadata.ReplicationSettings) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSettings) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
