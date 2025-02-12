// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	metadata "github.com/couchbase/goxdcr/metadata"
	mock "github.com/stretchr/testify/mock"
)

// PeerManifestsGetter is an autogenerated mock type for the PeerManifestsGetter type
type PeerManifestsGetter struct {
	mock.Mock
}

type PeerManifestsGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *PeerManifestsGetter) EXPECT() *PeerManifestsGetter_Expecter {
	return &PeerManifestsGetter_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: specId, specInternalId
func (_m *PeerManifestsGetter) Execute(specId string, specInternalId string) (*metadata.CollectionsManifestPair, error) {
	ret := _m.Called(specId, specInternalId)

	var r0 *metadata.CollectionsManifestPair
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*metadata.CollectionsManifestPair, error)); ok {
		return rf(specId, specInternalId)
	}
	if rf, ok := ret.Get(0).(func(string, string) *metadata.CollectionsManifestPair); ok {
		r0 = rf(specId, specInternalId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.CollectionsManifestPair)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(specId, specInternalId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PeerManifestsGetter_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type PeerManifestsGetter_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - specId string
//   - specInternalId string
func (_e *PeerManifestsGetter_Expecter) Execute(specId interface{}, specInternalId interface{}) *PeerManifestsGetter_Execute_Call {
	return &PeerManifestsGetter_Execute_Call{Call: _e.mock.On("Execute", specId, specInternalId)}
}

func (_c *PeerManifestsGetter_Execute_Call) Run(run func(specId string, specInternalId string)) *PeerManifestsGetter_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *PeerManifestsGetter_Execute_Call) Return(_a0 *metadata.CollectionsManifestPair, _a1 error) *PeerManifestsGetter_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PeerManifestsGetter_Execute_Call) RunAndReturn(run func(string, string) (*metadata.CollectionsManifestPair, error)) *PeerManifestsGetter_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewPeerManifestsGetter creates a new instance of PeerManifestsGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPeerManifestsGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *PeerManifestsGetter {
	mock := &PeerManifestsGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
