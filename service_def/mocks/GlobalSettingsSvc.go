package mocks

import base "github.com/couchbase/goxdcr/base"
import metadata "github.com/couchbase/goxdcr/metadata"
import mock "github.com/stretchr/testify/mock"

// GlobalSettingsSvc is an autogenerated mock type for the GlobalSettingsSvc type
type GlobalSettingsSvc struct {
	mock.Mock
}

// GetDefaultGlobalSettings provides a mock function with given fields:
func (_m *GlobalSettingsSvc) GetDefaultGlobalSettings() (*metadata.GlobalSettings, error) {
	ret := _m.Called()

	var r0 *metadata.GlobalSettings
	if rf, ok := ret.Get(0).(func() *metadata.GlobalSettings); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.GlobalSettings)
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

// GlobalSettingsServiceCallback provides a mock function with given fields: path, value, rev
func (_m *GlobalSettingsSvc) GlobalSettingsServiceCallback(path string, value []byte, rev interface{}) error {
	ret := _m.Called(path, value, rev)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, interface{}) error); ok {
		r0 = rf(path, value, rev)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetDefaultGlobalSettings provides a mock function with given fields: _a0
func (_m *GlobalSettingsSvc) SetDefaultGlobalSettings(_a0 *metadata.GlobalSettings) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.GlobalSettings) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetMetadataChangeHandlerCallback provides a mock function with given fields: callBack
func (_m *GlobalSettingsSvc) SetMetadataChangeHandlerCallback(callBack base.MetadataChangeHandlerCallback) {
	_m.Called(callBack)
}
