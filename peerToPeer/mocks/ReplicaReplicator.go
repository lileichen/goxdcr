// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	metadata "github.com/couchbase/goxdcr/metadata"
	mock "github.com/stretchr/testify/mock"
)

// ReplicaReplicator is an autogenerated mock type for the ReplicaReplicator type
type ReplicaReplicator struct {
	mock.Mock
}

// HandleSpecChange provides a mock function with given fields: oldSpec, newSpec
func (_m *ReplicaReplicator) HandleSpecChange(oldSpec *metadata.ReplicationSpecification, newSpec *metadata.ReplicationSpecification) {
	_m.Called(oldSpec, newSpec)
}

// HandleSpecCreation provides a mock function with given fields: spec
func (_m *ReplicaReplicator) HandleSpecCreation(spec *metadata.ReplicationSpecification) {
	_m.Called(spec)
}

// HandleSpecDeletion provides a mock function with given fields: oldSpec
func (_m *ReplicaReplicator) HandleSpecDeletion(oldSpec *metadata.ReplicationSpecification) {
	_m.Called(oldSpec)
}

// Start provides a mock function with given fields:
func (_m *ReplicaReplicator) Start() {
	_m.Called()
}

// Stop provides a mock function with given fields:
func (_m *ReplicaReplicator) Stop() {
	_m.Called()
}