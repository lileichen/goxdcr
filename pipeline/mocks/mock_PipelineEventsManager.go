// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"

	pipeline "github.com/couchbase/goxdcr/pipeline"
)

// PipelineEventsManager is an autogenerated mock type for the PipelineEventsManager type
type PipelineEventsManager struct {
	mock.Mock
}

type PipelineEventsManager_Expecter struct {
	mock *mock.Mock
}

func (_m *PipelineEventsManager) EXPECT() *PipelineEventsManager_Expecter {
	return &PipelineEventsManager_Expecter{mock: &_m.Mock}
}

// AddEvent provides a mock function with given fields: eventType, eventDesc, eventExtras, hint
func (_m *PipelineEventsManager) AddEvent(eventType base.EventInfoType, eventDesc string, eventExtras base.EventsMap, hint interface{}) int64 {
	ret := _m.Called(eventType, eventDesc, eventExtras, hint)

	var r0 int64
	if rf, ok := ret.Get(0).(func(base.EventInfoType, string, base.EventsMap, interface{}) int64); ok {
		r0 = rf(eventType, eventDesc, eventExtras, hint)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// PipelineEventsManager_AddEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddEvent'
type PipelineEventsManager_AddEvent_Call struct {
	*mock.Call
}

// AddEvent is a helper method to define mock.On call
//   - eventType base.EventInfoType
//   - eventDesc string
//   - eventExtras base.EventsMap
//   - hint interface{}
func (_e *PipelineEventsManager_Expecter) AddEvent(eventType interface{}, eventDesc interface{}, eventExtras interface{}, hint interface{}) *PipelineEventsManager_AddEvent_Call {
	return &PipelineEventsManager_AddEvent_Call{Call: _e.mock.On("AddEvent", eventType, eventDesc, eventExtras, hint)}
}

func (_c *PipelineEventsManager_AddEvent_Call) Run(run func(eventType base.EventInfoType, eventDesc string, eventExtras base.EventsMap, hint interface{})) *PipelineEventsManager_AddEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(base.EventInfoType), args[1].(string), args[2].(base.EventsMap), args[3].(interface{}))
	})
	return _c
}

func (_c *PipelineEventsManager_AddEvent_Call) Return(eventId int64) *PipelineEventsManager_AddEvent_Call {
	_c.Call.Return(eventId)
	return _c
}

func (_c *PipelineEventsManager_AddEvent_Call) RunAndReturn(run func(base.EventInfoType, string, base.EventsMap, interface{}) int64) *PipelineEventsManager_AddEvent_Call {
	_c.Call.Return(run)
	return _c
}

// BackfillUpdateCb provides a mock function with given fields: diffPair, srcManifestsDelta
func (_m *PipelineEventsManager) BackfillUpdateCb(diffPair *metadata.CollectionNamespaceMappingsDiffPair, srcManifestsDelta []*metadata.CollectionsManifest) error {
	ret := _m.Called(diffPair, srcManifestsDelta)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.CollectionNamespaceMappingsDiffPair, []*metadata.CollectionsManifest) error); ok {
		r0 = rf(diffPair, srcManifestsDelta)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineEventsManager_BackfillUpdateCb_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BackfillUpdateCb'
type PipelineEventsManager_BackfillUpdateCb_Call struct {
	*mock.Call
}

// BackfillUpdateCb is a helper method to define mock.On call
//   - diffPair *metadata.CollectionNamespaceMappingsDiffPair
//   - srcManifestsDelta []*metadata.CollectionsManifest
func (_e *PipelineEventsManager_Expecter) BackfillUpdateCb(diffPair interface{}, srcManifestsDelta interface{}) *PipelineEventsManager_BackfillUpdateCb_Call {
	return &PipelineEventsManager_BackfillUpdateCb_Call{Call: _e.mock.On("BackfillUpdateCb", diffPair, srcManifestsDelta)}
}

func (_c *PipelineEventsManager_BackfillUpdateCb_Call) Run(run func(diffPair *metadata.CollectionNamespaceMappingsDiffPair, srcManifestsDelta []*metadata.CollectionsManifest)) *PipelineEventsManager_BackfillUpdateCb_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*metadata.CollectionNamespaceMappingsDiffPair), args[1].([]*metadata.CollectionsManifest))
	})
	return _c
}

func (_c *PipelineEventsManager_BackfillUpdateCb_Call) Return(_a0 error) *PipelineEventsManager_BackfillUpdateCb_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineEventsManager_BackfillUpdateCb_Call) RunAndReturn(run func(*metadata.CollectionNamespaceMappingsDiffPair, []*metadata.CollectionsManifest) error) *PipelineEventsManager_BackfillUpdateCb_Call {
	_c.Call.Return(run)
	return _c
}

// ClearNonBrokenMapEvents provides a mock function with given fields:
func (_m *PipelineEventsManager) ClearNonBrokenMapEvents() {
	_m.Called()
}

// PipelineEventsManager_ClearNonBrokenMapEvents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClearNonBrokenMapEvents'
type PipelineEventsManager_ClearNonBrokenMapEvents_Call struct {
	*mock.Call
}

// ClearNonBrokenMapEvents is a helper method to define mock.On call
func (_e *PipelineEventsManager_Expecter) ClearNonBrokenMapEvents() *PipelineEventsManager_ClearNonBrokenMapEvents_Call {
	return &PipelineEventsManager_ClearNonBrokenMapEvents_Call{Call: _e.mock.On("ClearNonBrokenMapEvents")}
}

func (_c *PipelineEventsManager_ClearNonBrokenMapEvents_Call) Run(run func()) *PipelineEventsManager_ClearNonBrokenMapEvents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PipelineEventsManager_ClearNonBrokenMapEvents_Call) Return() *PipelineEventsManager_ClearNonBrokenMapEvents_Call {
	_c.Call.Return()
	return _c
}

func (_c *PipelineEventsManager_ClearNonBrokenMapEvents_Call) RunAndReturn(run func()) *PipelineEventsManager_ClearNonBrokenMapEvents_Call {
	_c.Call.Return(run)
	return _c
}

// ClearNonBrokenMapEventsWithString provides a mock function with given fields: substr
func (_m *PipelineEventsManager) ClearNonBrokenMapEventsWithString(substr string) {
	_m.Called(substr)
}

// PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClearNonBrokenMapEventsWithString'
type PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call struct {
	*mock.Call
}

// ClearNonBrokenMapEventsWithString is a helper method to define mock.On call
//   - substr string
func (_e *PipelineEventsManager_Expecter) ClearNonBrokenMapEventsWithString(substr interface{}) *PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call {
	return &PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call{Call: _e.mock.On("ClearNonBrokenMapEventsWithString", substr)}
}

func (_c *PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call) Run(run func(substr string)) *PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call) Return() *PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call {
	_c.Call.Return()
	return _c
}

func (_c *PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call) RunAndReturn(run func(string)) *PipelineEventsManager_ClearNonBrokenMapEventsWithString_Call {
	_c.Call.Return(run)
	return _c
}

// ContainsEvent provides a mock function with given fields: eventId
func (_m *PipelineEventsManager) ContainsEvent(eventId int) bool {
	ret := _m.Called(eventId)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int) bool); ok {
		r0 = rf(eventId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PipelineEventsManager_ContainsEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ContainsEvent'
type PipelineEventsManager_ContainsEvent_Call struct {
	*mock.Call
}

// ContainsEvent is a helper method to define mock.On call
//   - eventId int
func (_e *PipelineEventsManager_Expecter) ContainsEvent(eventId interface{}) *PipelineEventsManager_ContainsEvent_Call {
	return &PipelineEventsManager_ContainsEvent_Call{Call: _e.mock.On("ContainsEvent", eventId)}
}

func (_c *PipelineEventsManager_ContainsEvent_Call) Run(run func(eventId int)) *PipelineEventsManager_ContainsEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *PipelineEventsManager_ContainsEvent_Call) Return(_a0 bool) *PipelineEventsManager_ContainsEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineEventsManager_ContainsEvent_Call) RunAndReturn(run func(int) bool) *PipelineEventsManager_ContainsEvent_Call {
	_c.Call.Return(run)
	return _c
}

// DismissEvent provides a mock function with given fields: eventId
func (_m *PipelineEventsManager) DismissEvent(eventId int) error {
	ret := _m.Called(eventId)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(eventId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineEventsManager_DismissEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DismissEvent'
type PipelineEventsManager_DismissEvent_Call struct {
	*mock.Call
}

// DismissEvent is a helper method to define mock.On call
//   - eventId int
func (_e *PipelineEventsManager_Expecter) DismissEvent(eventId interface{}) *PipelineEventsManager_DismissEvent_Call {
	return &PipelineEventsManager_DismissEvent_Call{Call: _e.mock.On("DismissEvent", eventId)}
}

func (_c *PipelineEventsManager_DismissEvent_Call) Run(run func(eventId int)) *PipelineEventsManager_DismissEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *PipelineEventsManager_DismissEvent_Call) Return(_a0 error) *PipelineEventsManager_DismissEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineEventsManager_DismissEvent_Call) RunAndReturn(run func(int) error) *PipelineEventsManager_DismissEvent_Call {
	_c.Call.Return(run)
	return _c
}

// GetCurrentEvents provides a mock function with given fields:
func (_m *PipelineEventsManager) GetCurrentEvents() *pipeline.PipelineEventList {
	ret := _m.Called()

	var r0 *pipeline.PipelineEventList
	if rf, ok := ret.Get(0).(func() *pipeline.PipelineEventList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pipeline.PipelineEventList)
		}
	}

	return r0
}

// PipelineEventsManager_GetCurrentEvents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentEvents'
type PipelineEventsManager_GetCurrentEvents_Call struct {
	*mock.Call
}

// GetCurrentEvents is a helper method to define mock.On call
func (_e *PipelineEventsManager_Expecter) GetCurrentEvents() *PipelineEventsManager_GetCurrentEvents_Call {
	return &PipelineEventsManager_GetCurrentEvents_Call{Call: _e.mock.On("GetCurrentEvents")}
}

func (_c *PipelineEventsManager_GetCurrentEvents_Call) Run(run func()) *PipelineEventsManager_GetCurrentEvents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PipelineEventsManager_GetCurrentEvents_Call) Return(_a0 *pipeline.PipelineEventList) *PipelineEventsManager_GetCurrentEvents_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineEventsManager_GetCurrentEvents_Call) RunAndReturn(run func() *pipeline.PipelineEventList) *PipelineEventsManager_GetCurrentEvents_Call {
	_c.Call.Return(run)
	return _c
}

// LoadLatestBrokenMap provides a mock function with given fields: mapping
func (_m *PipelineEventsManager) LoadLatestBrokenMap(mapping metadata.CollectionNamespaceMapping) {
	_m.Called(mapping)
}

// PipelineEventsManager_LoadLatestBrokenMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadLatestBrokenMap'
type PipelineEventsManager_LoadLatestBrokenMap_Call struct {
	*mock.Call
}

// LoadLatestBrokenMap is a helper method to define mock.On call
//   - mapping metadata.CollectionNamespaceMapping
func (_e *PipelineEventsManager_Expecter) LoadLatestBrokenMap(mapping interface{}) *PipelineEventsManager_LoadLatestBrokenMap_Call {
	return &PipelineEventsManager_LoadLatestBrokenMap_Call{Call: _e.mock.On("LoadLatestBrokenMap", mapping)}
}

func (_c *PipelineEventsManager_LoadLatestBrokenMap_Call) Run(run func(mapping metadata.CollectionNamespaceMapping)) *PipelineEventsManager_LoadLatestBrokenMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.CollectionNamespaceMapping))
	})
	return _c
}

func (_c *PipelineEventsManager_LoadLatestBrokenMap_Call) Return() *PipelineEventsManager_LoadLatestBrokenMap_Call {
	_c.Call.Return()
	return _c
}

func (_c *PipelineEventsManager_LoadLatestBrokenMap_Call) RunAndReturn(run func(metadata.CollectionNamespaceMapping)) *PipelineEventsManager_LoadLatestBrokenMap_Call {
	_c.Call.Return(run)
	return _c
}

// ResetDismissedHistory provides a mock function with given fields:
func (_m *PipelineEventsManager) ResetDismissedHistory() {
	_m.Called()
}

// PipelineEventsManager_ResetDismissedHistory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetDismissedHistory'
type PipelineEventsManager_ResetDismissedHistory_Call struct {
	*mock.Call
}

// ResetDismissedHistory is a helper method to define mock.On call
func (_e *PipelineEventsManager_Expecter) ResetDismissedHistory() *PipelineEventsManager_ResetDismissedHistory_Call {
	return &PipelineEventsManager_ResetDismissedHistory_Call{Call: _e.mock.On("ResetDismissedHistory")}
}

func (_c *PipelineEventsManager_ResetDismissedHistory_Call) Run(run func()) *PipelineEventsManager_ResetDismissedHistory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PipelineEventsManager_ResetDismissedHistory_Call) Return() *PipelineEventsManager_ResetDismissedHistory_Call {
	_c.Call.Return()
	return _c
}

func (_c *PipelineEventsManager_ResetDismissedHistory_Call) RunAndReturn(run func()) *PipelineEventsManager_ResetDismissedHistory_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateEvent provides a mock function with given fields: eventId, eventDesc, eventExtras
func (_m *PipelineEventsManager) UpdateEvent(eventId int64, eventDesc string, eventExtras *base.EventsMap) error {
	ret := _m.Called(eventId, eventDesc, eventExtras)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, string, *base.EventsMap) error); ok {
		r0 = rf(eventId, eventDesc, eventExtras)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineEventsManager_UpdateEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateEvent'
type PipelineEventsManager_UpdateEvent_Call struct {
	*mock.Call
}

// UpdateEvent is a helper method to define mock.On call
//   - eventId int64
//   - eventDesc string
//   - eventExtras *base.EventsMap
func (_e *PipelineEventsManager_Expecter) UpdateEvent(eventId interface{}, eventDesc interface{}, eventExtras interface{}) *PipelineEventsManager_UpdateEvent_Call {
	return &PipelineEventsManager_UpdateEvent_Call{Call: _e.mock.On("UpdateEvent", eventId, eventDesc, eventExtras)}
}

func (_c *PipelineEventsManager_UpdateEvent_Call) Run(run func(eventId int64, eventDesc string, eventExtras *base.EventsMap)) *PipelineEventsManager_UpdateEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(string), args[2].(*base.EventsMap))
	})
	return _c
}

func (_c *PipelineEventsManager_UpdateEvent_Call) Return(_a0 error) *PipelineEventsManager_UpdateEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineEventsManager_UpdateEvent_Call) RunAndReturn(run func(int64, string, *base.EventsMap) error) *PipelineEventsManager_UpdateEvent_Call {
	_c.Call.Return(run)
	return _c
}

// NewPipelineEventsManager creates a new instance of PipelineEventsManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPipelineEventsManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *PipelineEventsManager {
	mock := &PipelineEventsManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
