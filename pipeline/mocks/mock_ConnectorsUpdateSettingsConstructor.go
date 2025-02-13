// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"
)

// ConnectorsUpdateSettingsConstructor is an autogenerated mock type for the ConnectorsUpdateSettingsConstructor type
type ConnectorsUpdateSettingsConstructor struct {
	mock.Mock
}

type ConnectorsUpdateSettingsConstructor_Expecter struct {
	mock *mock.Mock
}

func (_m *ConnectorsUpdateSettingsConstructor) EXPECT() *ConnectorsUpdateSettingsConstructor_Expecter {
	return &ConnectorsUpdateSettingsConstructor_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0, connector, pipeline_settings
func (_m *ConnectorsUpdateSettingsConstructor) Execute(_a0 common.Pipeline, connector common.Connector, pipeline_settings metadata.ReplicationSettingsMap) (metadata.ReplicationSettingsMap, error) {
	ret := _m.Called(_a0, connector, pipeline_settings)

	var r0 metadata.ReplicationSettingsMap
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Pipeline, common.Connector, metadata.ReplicationSettingsMap) (metadata.ReplicationSettingsMap, error)); ok {
		return rf(_a0, connector, pipeline_settings)
	}
	if rf, ok := ret.Get(0).(func(common.Pipeline, common.Connector, metadata.ReplicationSettingsMap) metadata.ReplicationSettingsMap); ok {
		r0 = rf(_a0, connector, pipeline_settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metadata.ReplicationSettingsMap)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Pipeline, common.Connector, metadata.ReplicationSettingsMap) error); ok {
		r1 = rf(_a0, connector, pipeline_settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConnectorsUpdateSettingsConstructor_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type ConnectorsUpdateSettingsConstructor_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 common.Pipeline
//   - connector common.Connector
//   - pipeline_settings metadata.ReplicationSettingsMap
func (_e *ConnectorsUpdateSettingsConstructor_Expecter) Execute(_a0 interface{}, connector interface{}, pipeline_settings interface{}) *ConnectorsUpdateSettingsConstructor_Execute_Call {
	return &ConnectorsUpdateSettingsConstructor_Execute_Call{Call: _e.mock.On("Execute", _a0, connector, pipeline_settings)}
}

func (_c *ConnectorsUpdateSettingsConstructor_Execute_Call) Run(run func(_a0 common.Pipeline, connector common.Connector, pipeline_settings metadata.ReplicationSettingsMap)) *ConnectorsUpdateSettingsConstructor_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Pipeline), args[1].(common.Connector), args[2].(metadata.ReplicationSettingsMap))
	})
	return _c
}

func (_c *ConnectorsUpdateSettingsConstructor_Execute_Call) Return(_a0 metadata.ReplicationSettingsMap, _a1 error) *ConnectorsUpdateSettingsConstructor_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ConnectorsUpdateSettingsConstructor_Execute_Call) RunAndReturn(run func(common.Pipeline, common.Connector, metadata.ReplicationSettingsMap) (metadata.ReplicationSettingsMap, error)) *ConnectorsUpdateSettingsConstructor_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewConnectorsUpdateSettingsConstructor creates a new instance of ConnectorsUpdateSettingsConstructor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConnectorsUpdateSettingsConstructor(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConnectorsUpdateSettingsConstructor {
	mock := &ConnectorsUpdateSettingsConstructor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
