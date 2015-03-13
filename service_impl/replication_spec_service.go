// Copyright (c) 2013 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

// metadata service implementation leveraging gometa
package service_impl

import (
	"encoding/json"
	"fmt"
	"github.com/couchbase/cbauth/metakv"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/service_def"
	"sync"
	"errors"
	"strings"
)

const (
	// the key to the metadata that stores the keys of all Replication Specs
	ReplicationSpecsCatalogKey = metadata.ReplicationSpecKeyPrefix
)

var ReplicationSpecAlreadyExistErrorMessage = "Replication to the same remote cluster and bucket already exists"
var ReplicationSpecNotFoundErrorMessage = "Requested resource not found"

type ReplicationSpecService struct {
	metadata_svc       service_def.MetadataSvc
	uilog_svc          service_def.UILogSvc
	remote_cluster_svc service_def.RemoteClusterSvc
	call_back          service_def.SpecChangedCallback
	failure_call_back  service_def.SpecChangeListenerFailureCallBack
	logger             *log.CommonLogger
}

func NewReplicationSpecService(uilog_svc service_def.UILogSvc, remote_cluster_svc service_def.RemoteClusterSvc, metadata_svc service_def.MetadataSvc, logger_ctx *log.LoggerContext) *ReplicationSpecService {
	return &ReplicationSpecService{
		metadata_svc:       metadata_svc,
		uilog_svc:          uilog_svc,
		remote_cluster_svc: remote_cluster_svc,
		logger:             log.NewLogger("ReplicationSpecService", logger_ctx),
	}
}

func (service *ReplicationSpecService) StartSpecChangedCallBack(call_back service_def.SpecChangedCallback, failure_call_back service_def.SpecChangeListenerFailureCallBack,
	cancel <-chan struct{}, waitGrp *sync.WaitGroup) error {
	// start listening to changed to specs
	service.call_back = call_back
	service.failure_call_back = failure_call_back

	replSpecCatalogPath := GetCatalogPathFromCatalogKey(ReplicationSpecsCatalogKey)

	waitGrp.Add(1)
	go service.observeChildren(replSpecCatalogPath, cancel, waitGrp)

	return nil
}

func (service *ReplicationSpecService) observeChildren(dirpath string, cancel <-chan struct{}, waitGrp *sync.WaitGroup) {
	defer waitGrp.Done()
	err := metakv.RunObserveChildren(dirpath, service.metakvCallback, cancel)
	// call failure call back only when there are real errors
	// err may be nil when observeChildren is canceled, in which case there is no need to call failure call back
	service.failure_call_back(err)
}

func (service *ReplicationSpecService) ReplicationSpec(replicationId string) (*metadata.ReplicationSpecification, error) {
	result, rev, err := service.metadata_svc.Get(replicationId)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		service.logger.Errorf("Failed to get metadata %v, err=%v\n", replicationId, err)
		return nil, errors.New(ReplicationSpecNotFoundErrorMessage)
	}
	return constructReplicationSpec(result, rev)
}

func (service *ReplicationSpecService) AddReplicationSpec(spec *metadata.ReplicationSpecification) error {
	service.logger.Infof("Start AddReplicationSpec, spec=%v\n", spec)
	
	key := spec.Id
	_, err := service.ReplicationSpec(key)
	if err == nil {
		return errors.New(ReplicationSpecAlreadyExistErrorMessage)
	}
	
	value, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	service.logger.Info("Adding it to metadata store...")
	err = service.metadata_svc.AddWithCatalog(ReplicationSpecsCatalogKey, key, value)
	if err != nil {
		return err
	}
	service.logger.Info("log it with ale logger...")
	service.writeUiLog(spec, "created")
	service.logger.Info("Done with logging...")
	return nil
}

func (service *ReplicationSpecService) SetReplicationSpec(spec *metadata.ReplicationSpecification) error {
	key := spec.Id
	value, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	return service.metadata_svc.Set(key, value, spec.Revision)
}

func (service *ReplicationSpecService) DelReplicationSpec(replicationId string) (*metadata.ReplicationSpecification, error) {
	spec, err := service.ReplicationSpec(replicationId)
	if err != nil {
		return nil, errors.New(ReplicationSpecNotFoundErrorMessage)
	}

	err = service.metadata_svc.DelWithCatalog(ReplicationSpecsCatalogKey, replicationId, spec.Revision)
	if err != nil {
		return nil, err
	}

	service.writeUiLog(spec, "removed")
	return spec, nil
}

func (service *ReplicationSpecService) AllReplicationSpecs() (map[string]*metadata.ReplicationSpecification, error) {
	specs := make(map[string]*metadata.ReplicationSpecification, 0)

	entries, err := service.metadata_svc.GetAllMetadataFromCatalog(ReplicationSpecsCatalogKey)
	if err != nil {
		service.logger.Errorf("Failed to get all entries, err=%v\n", err)
		return nil, err
	}

	for _, entry := range entries {
		spec, err := constructReplicationSpec(entry.Value, entry.Rev)
		if err != nil {
			return nil, err
		}
		specs[entry.Key] = spec
	}

	return specs, nil
}

func (service *ReplicationSpecService) AllReplicationSpecIds() ([]string, error) {
	repIds, err := service.metadata_svc.GetAllKeysFromCatalog(ReplicationSpecsCatalogKey)
	if err != nil {
		return nil, err
	}
	return repIds, nil
}

func (service *ReplicationSpecService) AllReplicationSpecIdsForBucket(bucket string) ([]string, error) {
	var repIds []string
	allRepIds, err := service.AllReplicationSpecIds()
	if err != nil {
		return nil, err
	}

	if allRepIds != nil {
		for _, repId := range allRepIds {
			if metadata.IsReplicationIdForSourceBucket(repId, bucket) {
				repIds = append(repIds, repId)
			}
		}
	}
	return repIds, nil
}

func constructReplicationSpec(value []byte, rev interface{}) (*metadata.ReplicationSpecification, error) {
	if value == nil {
		return nil, nil
	}

	spec := &metadata.ReplicationSpecification{}
	err := json.Unmarshal(value, spec)
	if err != nil {
		return nil, err
	}
	spec.Revision = rev
	return spec, nil
}

// Implement callback function for metakv
func (service *ReplicationSpecService) metakvCallback(path string, value []byte, rev interface{}) error {
	service.logger.Infof("metakvCallback called on path = %v\n", path)

	if service.call_back != nil {
		spec, err := constructReplicationSpec(value, rev)
		if err != nil {
			// should never get here.
			service.logger.Errorf("Error marshaling replication spec. value=%v, err=%v\n", string(value), err)
			return err
		}

		err = service.call_back(GetKeyFromPath(path), spec)
		if err != nil {
			service.logger.Errorf("Replication spec change call back returned err=%v\n", err)
		}
		// do not return err since we do not want RunObserveChildren to abort
		return nil
	}

	return nil

}

func (service *ReplicationSpecService) writeUiLog(spec *metadata.ReplicationSpecification, action string) {
	if service.uilog_svc != nil {
		remoteClusterName := service.remote_cluster_svc.GetRemoteClusterNameFromClusterUuid(spec.TargetClusterUUID)

		uiLogMsg := fmt.Sprintf("Replication from bucket \"%s\" to bucket \"%s\" on cluster \"%s\" %s.", spec.SourceBucketName, spec.TargetBucketName, remoteClusterName, action)
		service.uilog_svc.Write(uiLogMsg)
	}
}

func (service *ReplicationSpecService) IsReplicationValidationError(err error) bool {
	if err != nil {
		return strings.HasPrefix(err.Error(), ReplicationSpecAlreadyExistErrorMessage) || strings.HasPrefix(err.Error(), ReplicationSpecNotFoundErrorMessage)
	} else {
		return false
	}
}