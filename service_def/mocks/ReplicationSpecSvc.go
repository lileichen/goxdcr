package mocks

import base "github.com/couchbase/goxdcr/base"
import metadata "github.com/couchbase/goxdcr/metadata"
import mock "github.com/stretchr/testify/mock"

// ReplicationSpecSvc is an autogenerated mock type for the ReplicationSpecSvc type
type ReplicationSpecSvc struct {
	mock.Mock
}

// AddReplicationSpec provides a mock function with given fields: spec
func (_m *ReplicationSpecSvc) AddReplicationSpec(spec *metadata.ReplicationSpecification) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AllReplicationSpecIds provides a mock function with given fields:
func (_m *ReplicationSpecSvc) AllReplicationSpecIds() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
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

// AllReplicationSpecIdsForBucket provides a mock function with given fields: bucket
func (_m *ReplicationSpecSvc) AllReplicationSpecIdsForBucket(bucket string) ([]string, error) {
	ret := _m.Called(bucket)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(bucket)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllReplicationSpecs provides a mock function with given fields:
func (_m *ReplicationSpecSvc) AllReplicationSpecs() (map[string]*metadata.ReplicationSpecification, error) {
	ret := _m.Called()

	var r0 map[string]*metadata.ReplicationSpecification
	if rf, ok := ret.Get(0).(func() map[string]*metadata.ReplicationSpecification); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*metadata.ReplicationSpecification)
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

// ConstructNewReplicationSpec provides a mock function with given fields: sourceBucketName, targetClusterUUID, targetBucketName
func (_m *ReplicationSpecSvc) ConstructNewReplicationSpec(sourceBucketName string, targetClusterUUID string, targetBucketName string) (*metadata.ReplicationSpecification, error) {
	ret := _m.Called(sourceBucketName, targetClusterUUID, targetBucketName)

	var r0 *metadata.ReplicationSpecification
	if rf, ok := ret.Get(0).(func(string, string, string) *metadata.ReplicationSpecification); ok {
		r0 = rf(sourceBucketName, targetClusterUUID, targetBucketName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.ReplicationSpecification)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(sourceBucketName, targetClusterUUID, targetBucketName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DelReplicationSpec provides a mock function with given fields: replicationId
func (_m *ReplicationSpecSvc) DelReplicationSpec(replicationId string) (*metadata.ReplicationSpecification, error) {
	ret := _m.Called(replicationId)

	var r0 *metadata.ReplicationSpecification
	if rf, ok := ret.Get(0).(func(string) *metadata.ReplicationSpecification); ok {
		r0 = rf(replicationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.ReplicationSpecification)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(replicationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DelReplicationSpecWithReason provides a mock function with given fields: replicationId, reason
func (_m *ReplicationSpecSvc) DelReplicationSpecWithReason(replicationId string, reason string) (*metadata.ReplicationSpecification, error) {
	ret := _m.Called(replicationId, reason)

	var r0 *metadata.ReplicationSpecification
	if rf, ok := ret.Get(0).(func(string, string) *metadata.ReplicationSpecification); ok {
		r0 = rf(replicationId, reason)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.ReplicationSpecification)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(replicationId, reason)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDerivedObj provides a mock function with given fields: specId
func (_m *ReplicationSpecSvc) GetDerivedObj(specId string) (interface{}, error) {
	ret := _m.Called(specId)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(specId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(specId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsReplicationValidationError provides a mock function with given fields: err
func (_m *ReplicationSpecSvc) IsReplicationValidationError(err error) bool {
	ret := _m.Called(err)

	var r0 bool
	if rf, ok := ret.Get(0).(func(error) bool); ok {
		r0 = rf(err)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ReplicationSpec provides a mock function with given fields: replicationId
func (_m *ReplicationSpecSvc) ReplicationSpec(replicationId string) (*metadata.ReplicationSpecification, error) {
	ret := _m.Called(replicationId)

	var r0 *metadata.ReplicationSpecification
	if rf, ok := ret.Get(0).(func(string) *metadata.ReplicationSpecification); ok {
		r0 = rf(replicationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.ReplicationSpecification)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(replicationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReplicationSpecServiceCallback provides a mock function with given fields: path, value, rev
func (_m *ReplicationSpecSvc) ReplicationSpecServiceCallback(path string, value []byte, rev interface{}) error {
	ret := _m.Called(path, value, rev)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, interface{}) error); ok {
		r0 = rf(path, value, rev)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetDerivedObj provides a mock function with given fields: specId, derivedObj
func (_m *ReplicationSpecSvc) SetDerivedObj(specId string, derivedObj interface{}) error {
	ret := _m.Called(specId, derivedObj)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(specId, derivedObj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetMetadataChangeHandlerCallback provides a mock function with given fields: callBack
func (_m *ReplicationSpecSvc) SetMetadataChangeHandlerCallback(callBack base.MetadataChangeHandlerCallback) {
	_m.Called(callBack)
}

// SetReplicationSpec provides a mock function with given fields: spec
func (_m *ReplicationSpecSvc) SetReplicationSpec(spec *metadata.ReplicationSpecification) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateAndGC provides a mock function with given fields: spec
func (_m *ReplicationSpecSvc) ValidateAndGC(spec *metadata.ReplicationSpecification) {
	_m.Called(spec)
}

// ValidateNewReplicationSpec provides a mock function with given fields: sourceBucket, targetCluster, targetBucket, settings
func (_m *ReplicationSpecSvc) ValidateNewReplicationSpec(sourceBucket string, targetCluster string, targetBucket string, settings map[string]interface{}) (string, string, *metadata.RemoteClusterReference, map[string]error) {
	ret := _m.Called(sourceBucket, targetCluster, targetBucket, settings)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, map[string]interface{}) string); ok {
		r0 = rf(sourceBucket, targetCluster, targetBucket, settings)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string, string, string, map[string]interface{}) string); ok {
		r1 = rf(sourceBucket, targetCluster, targetBucket, settings)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 *metadata.RemoteClusterReference
	if rf, ok := ret.Get(2).(func(string, string, string, map[string]interface{}) *metadata.RemoteClusterReference); ok {
		r2 = rf(sourceBucket, targetCluster, targetBucket, settings)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*metadata.RemoteClusterReference)
		}
	}

	var r3 map[string]error
	if rf, ok := ret.Get(3).(func(string, string, string, map[string]interface{}) map[string]error); ok {
		r3 = rf(sourceBucket, targetCluster, targetBucket, settings)
	} else {
		if ret.Get(3) != nil {
			r3 = ret.Get(3).(map[string]error)
		}
	}

	return r0, r1, r2, r3
}