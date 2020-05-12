package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"
)

// PipelineService is an autogenerated mock type for the PipelineService type
type PipelineService struct {
	mock.Mock
}

// Attach provides a mock function with given fields: pipeline
func (_m *PipelineService) Attach(pipeline common.Pipeline) error {
	ret := _m.Called(pipeline)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(pipeline)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Detach provides a mock function with given fields: pipeline
func (_m *PipelineService) Detach(pipeline common.Pipeline) error {
	ret := _m.Called(pipeline)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(pipeline)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsSharable provides a mock function with given fields:
func (_m *PipelineService) IsSharable() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *PipelineService) Start(_a0 metadata.ReplicationSettingsMap) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *PipelineService) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSettings provides a mock function with given fields: settings
func (_m *PipelineService) UpdateSettings(settings metadata.ReplicationSettingsMap) error {
	ret := _m.Called(settings)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) error); ok {
		r0 = rf(settings)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
