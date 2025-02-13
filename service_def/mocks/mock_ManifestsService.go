// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	metadata "github.com/couchbase/goxdcr/metadata"
	mock "github.com/stretchr/testify/mock"
)

// ManifestsService is an autogenerated mock type for the ManifestsService type
type ManifestsService struct {
	mock.Mock
}

type ManifestsService_Expecter struct {
	mock *mock.Mock
}

func (_m *ManifestsService) EXPECT() *ManifestsService_Expecter {
	return &ManifestsService_Expecter{mock: &_m.Mock}
}

// DelManifests provides a mock function with given fields: replSpec
func (_m *ManifestsService) DelManifests(replSpec *metadata.ReplicationSpecification) error {
	ret := _m.Called(replSpec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) error); ok {
		r0 = rf(replSpec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManifestsService_DelManifests_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DelManifests'
type ManifestsService_DelManifests_Call struct {
	*mock.Call
}

// DelManifests is a helper method to define mock.On call
//   - replSpec *metadata.ReplicationSpecification
func (_e *ManifestsService_Expecter) DelManifests(replSpec interface{}) *ManifestsService_DelManifests_Call {
	return &ManifestsService_DelManifests_Call{Call: _e.mock.On("DelManifests", replSpec)}
}

func (_c *ManifestsService_DelManifests_Call) Run(run func(replSpec *metadata.ReplicationSpecification)) *ManifestsService_DelManifests_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*metadata.ReplicationSpecification))
	})
	return _c
}

func (_c *ManifestsService_DelManifests_Call) Return(_a0 error) *ManifestsService_DelManifests_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManifestsService_DelManifests_Call) RunAndReturn(run func(*metadata.ReplicationSpecification) error) *ManifestsService_DelManifests_Call {
	_c.Call.Return(run)
	return _c
}

// GetSourceManifests provides a mock function with given fields: replSpec
func (_m *ManifestsService) GetSourceManifests(replSpec *metadata.ReplicationSpecification) (*metadata.ManifestsList, error) {
	ret := _m.Called(replSpec)

	var r0 *metadata.ManifestsList
	var r1 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) (*metadata.ManifestsList, error)); ok {
		return rf(replSpec)
	}
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) *metadata.ManifestsList); ok {
		r0 = rf(replSpec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.ManifestsList)
		}
	}

	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification) error); ok {
		r1 = rf(replSpec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ManifestsService_GetSourceManifests_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSourceManifests'
type ManifestsService_GetSourceManifests_Call struct {
	*mock.Call
}

// GetSourceManifests is a helper method to define mock.On call
//   - replSpec *metadata.ReplicationSpecification
func (_e *ManifestsService_Expecter) GetSourceManifests(replSpec interface{}) *ManifestsService_GetSourceManifests_Call {
	return &ManifestsService_GetSourceManifests_Call{Call: _e.mock.On("GetSourceManifests", replSpec)}
}

func (_c *ManifestsService_GetSourceManifests_Call) Run(run func(replSpec *metadata.ReplicationSpecification)) *ManifestsService_GetSourceManifests_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*metadata.ReplicationSpecification))
	})
	return _c
}

func (_c *ManifestsService_GetSourceManifests_Call) Return(_a0 *metadata.ManifestsList, _a1 error) *ManifestsService_GetSourceManifests_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ManifestsService_GetSourceManifests_Call) RunAndReturn(run func(*metadata.ReplicationSpecification) (*metadata.ManifestsList, error)) *ManifestsService_GetSourceManifests_Call {
	_c.Call.Return(run)
	return _c
}

// GetTargetManifests provides a mock function with given fields: replSpec
func (_m *ManifestsService) GetTargetManifests(replSpec *metadata.ReplicationSpecification) (*metadata.ManifestsList, error) {
	ret := _m.Called(replSpec)

	var r0 *metadata.ManifestsList
	var r1 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) (*metadata.ManifestsList, error)); ok {
		return rf(replSpec)
	}
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) *metadata.ManifestsList); ok {
		r0 = rf(replSpec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.ManifestsList)
		}
	}

	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification) error); ok {
		r1 = rf(replSpec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ManifestsService_GetTargetManifests_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTargetManifests'
type ManifestsService_GetTargetManifests_Call struct {
	*mock.Call
}

// GetTargetManifests is a helper method to define mock.On call
//   - replSpec *metadata.ReplicationSpecification
func (_e *ManifestsService_Expecter) GetTargetManifests(replSpec interface{}) *ManifestsService_GetTargetManifests_Call {
	return &ManifestsService_GetTargetManifests_Call{Call: _e.mock.On("GetTargetManifests", replSpec)}
}

func (_c *ManifestsService_GetTargetManifests_Call) Run(run func(replSpec *metadata.ReplicationSpecification)) *ManifestsService_GetTargetManifests_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*metadata.ReplicationSpecification))
	})
	return _c
}

func (_c *ManifestsService_GetTargetManifests_Call) Return(_a0 *metadata.ManifestsList, _a1 error) *ManifestsService_GetTargetManifests_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ManifestsService_GetTargetManifests_Call) RunAndReturn(run func(*metadata.ReplicationSpecification) (*metadata.ManifestsList, error)) *ManifestsService_GetTargetManifests_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertSourceManifests provides a mock function with given fields: replSpec, src
func (_m *ManifestsService) UpsertSourceManifests(replSpec *metadata.ReplicationSpecification, src *metadata.ManifestsList) error {
	ret := _m.Called(replSpec, src)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification, *metadata.ManifestsList) error); ok {
		r0 = rf(replSpec, src)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManifestsService_UpsertSourceManifests_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertSourceManifests'
type ManifestsService_UpsertSourceManifests_Call struct {
	*mock.Call
}

// UpsertSourceManifests is a helper method to define mock.On call
//   - replSpec *metadata.ReplicationSpecification
//   - src *metadata.ManifestsList
func (_e *ManifestsService_Expecter) UpsertSourceManifests(replSpec interface{}, src interface{}) *ManifestsService_UpsertSourceManifests_Call {
	return &ManifestsService_UpsertSourceManifests_Call{Call: _e.mock.On("UpsertSourceManifests", replSpec, src)}
}

func (_c *ManifestsService_UpsertSourceManifests_Call) Run(run func(replSpec *metadata.ReplicationSpecification, src *metadata.ManifestsList)) *ManifestsService_UpsertSourceManifests_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*metadata.ReplicationSpecification), args[1].(*metadata.ManifestsList))
	})
	return _c
}

func (_c *ManifestsService_UpsertSourceManifests_Call) Return(_a0 error) *ManifestsService_UpsertSourceManifests_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManifestsService_UpsertSourceManifests_Call) RunAndReturn(run func(*metadata.ReplicationSpecification, *metadata.ManifestsList) error) *ManifestsService_UpsertSourceManifests_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertTargetManifests provides a mock function with given fields: replSpec, tgt
func (_m *ManifestsService) UpsertTargetManifests(replSpec *metadata.ReplicationSpecification, tgt *metadata.ManifestsList) error {
	ret := _m.Called(replSpec, tgt)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification, *metadata.ManifestsList) error); ok {
		r0 = rf(replSpec, tgt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManifestsService_UpsertTargetManifests_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertTargetManifests'
type ManifestsService_UpsertTargetManifests_Call struct {
	*mock.Call
}

// UpsertTargetManifests is a helper method to define mock.On call
//   - replSpec *metadata.ReplicationSpecification
//   - tgt *metadata.ManifestsList
func (_e *ManifestsService_Expecter) UpsertTargetManifests(replSpec interface{}, tgt interface{}) *ManifestsService_UpsertTargetManifests_Call {
	return &ManifestsService_UpsertTargetManifests_Call{Call: _e.mock.On("UpsertTargetManifests", replSpec, tgt)}
}

func (_c *ManifestsService_UpsertTargetManifests_Call) Run(run func(replSpec *metadata.ReplicationSpecification, tgt *metadata.ManifestsList)) *ManifestsService_UpsertTargetManifests_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*metadata.ReplicationSpecification), args[1].(*metadata.ManifestsList))
	})
	return _c
}

func (_c *ManifestsService_UpsertTargetManifests_Call) Return(_a0 error) *ManifestsService_UpsertTargetManifests_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManifestsService_UpsertTargetManifests_Call) RunAndReturn(run func(*metadata.ReplicationSpecification, *metadata.ManifestsList) error) *ManifestsService_UpsertTargetManifests_Call {
	_c.Call.Return(run)
	return _c
}

// NewManifestsService creates a new instance of ManifestsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewManifestsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ManifestsService {
	mock := &ManifestsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
