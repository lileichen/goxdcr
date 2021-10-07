// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"
)

// PartsSettingsConstructor is an autogenerated mock type for the PartsSettingsConstructor type
type PartsSettingsConstructor struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, part, pipeline_settings, targetClusterref, ssl_port_map
func (_m *PartsSettingsConstructor) Execute(_a0 common.Pipeline, part common.Part, pipeline_settings metadata.ReplicationSettingsMap, targetClusterref *metadata.RemoteClusterReference, ssl_port_map map[string]uint16) (metadata.ReplicationSettingsMap, error) {
	ret := _m.Called(_a0, part, pipeline_settings, targetClusterref, ssl_port_map)

	var r0 metadata.ReplicationSettingsMap
	if rf, ok := ret.Get(0).(func(common.Pipeline, common.Part, metadata.ReplicationSettingsMap, *metadata.RemoteClusterReference, map[string]uint16) metadata.ReplicationSettingsMap); ok {
		r0 = rf(_a0, part, pipeline_settings, targetClusterref, ssl_port_map)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metadata.ReplicationSettingsMap)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, common.Part, metadata.ReplicationSettingsMap, *metadata.RemoteClusterReference, map[string]uint16) error); ok {
		r1 = rf(_a0, part, pipeline_settings, targetClusterref, ssl_port_map)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
