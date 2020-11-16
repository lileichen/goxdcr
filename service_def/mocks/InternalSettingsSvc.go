// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"
)

// InternalSettingsSvc is an autogenerated mock type for the InternalSettingsSvc type
type InternalSettingsSvc struct {
	mock.Mock
}

// GetInternalSettings provides a mock function with given fields:
func (_m *InternalSettingsSvc) GetInternalSettings() *metadata.InternalSettings {
	ret := _m.Called()

	var r0 *metadata.InternalSettings
	if rf, ok := ret.Get(0).(func() *metadata.InternalSettings); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.InternalSettings)
		}
	}

	return r0
}

// InternalSettingsServiceCallback provides a mock function with given fields: path, value, rev
func (_m *InternalSettingsSvc) InternalSettingsServiceCallback(path string, value []byte, rev interface{}) error {
	ret := _m.Called(path, value, rev)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, interface{}) error); ok {
		r0 = rf(path, value, rev)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetMetadataChangeHandlerCallback provides a mock function with given fields: callBack
func (_m *InternalSettingsSvc) SetMetadataChangeHandlerCallback(callBack base.MetadataChangeHandlerCallback) {
	_m.Called(callBack)
}

// UpdateInternalSettings provides a mock function with given fields: settingsMap
func (_m *InternalSettingsSvc) UpdateInternalSettings(settingsMap map[string]interface{}) (*metadata.InternalSettings, map[string]error, error) {
	ret := _m.Called(settingsMap)

	var r0 *metadata.InternalSettings
	if rf, ok := ret.Get(0).(func(map[string]interface{}) *metadata.InternalSettings); ok {
		r0 = rf(settingsMap)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.InternalSettings)
		}
	}

	var r1 map[string]error
	if rf, ok := ret.Get(1).(func(map[string]interface{}) map[string]error); ok {
		r1 = rf(settingsMap)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]error)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(map[string]interface{}) error); ok {
		r2 = rf(settingsMap)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
