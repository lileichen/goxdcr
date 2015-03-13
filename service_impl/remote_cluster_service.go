// Copyright (c) 2013 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the refific language governing permissions
// and limitations under the License.

// metadata service implementation leveraging gometa
package service_impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/service_def"
	"github.com/couchbase/goxdcr/utils"
	"reflect"
	"strings"
	"sync"
	"net/http"
)

const (
	// the key to the metadata that stores the keys of all remote clusters
	RemoteClustersCatalogKey = metadata.RemoteClusterKeyPrefix
)

var InvalidRemoteClusterErrorMessage = "Invalid remote cluster. "
var UnknownRemoteClusterErrorMessage = "unknown remote cluster"

type remoteClusterCache struct {
	key                 string
	nodes_connectionstr []string
	ref                 *metadata.RemoteClusterReference
}

type RemoteClusterService struct {
	metadata_svc service_def.MetadataSvc
	uilog_svc    service_def.UILogSvc
	xdcr_topology_svc service_def.XDCRCompTopologySvc
	logger       *log.CommonLogger
	cache_lock   *sync.RWMutex
	cache_map    map[string]*remoteClusterCache
}

func NewRemoteClusterService(uilog_svc service_def.UILogSvc, metadata_svc service_def.MetadataSvc, xdcr_topology_svc service_def.XDCRCompTopologySvc, logger_ctx *log.LoggerContext) *RemoteClusterService {
	return &RemoteClusterService{
		metadata_svc: metadata_svc,
		cache_lock:   &sync.RWMutex{},
		cache_map:    make(map[string]*remoteClusterCache),
		uilog_svc:    uilog_svc,
		xdcr_topology_svc: xdcr_topology_svc,
		logger:       log.NewLogger("RemoteClusterService", logger_ctx),
	}
}

func (service *RemoteClusterService) RemoteClusterByRefId(refId string, refresh bool) (*metadata.RemoteClusterReference, error) {
	result, rev, err := service.metadata_svc.Get(refId)
	if err != nil {
		return nil, err
	}
	ref, err := service.constructRemoteClusterReference(result, rev)
	if err != nil {
		service.logger.Errorf("Failed to get remote cluster refId=%v, err=%v\n", refId, err)
		return nil, err
	}

	if refresh {
		ref, err = service.refresh(ref)
	}

	return ref, err
}

func (service *RemoteClusterService) RemoteClusterByUuid(uuid string, refresh bool) (*metadata.RemoteClusterReference, error) {
	return service.RemoteClusterByRefId(metadata.RemoteClusterRefId(uuid), refresh)
}

func (service *RemoteClusterService) RemoteClusterByRefName(refName string, refresh bool) (*metadata.RemoteClusterReference, error) {
	var ref *metadata.RemoteClusterReference
	results, err := service.RemoteClusters(false)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if result.Name == refName {
			ref = result
			break
		}
	}

	if ref == nil {
		return nil, errors.New(UnknownRemoteClusterErrorMessage)
	} else {
		if refresh {
			ref, err = service.refresh(ref)
		}
		return ref, err
	}
}

func (service *RemoteClusterService) AddRemoteCluster(ref *metadata.RemoteClusterReference) error {
	service.logger.Infof("Adding remote cluster with referenceId %v\n", ref.Id)

	err := service.ValidateRemoteCluster(ref)
	if err != nil {
		return err
	}

	err = service.addRemoteCluster(ref)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	if service.uilog_svc != nil {
		uiLogMsg := fmt.Sprintf("Created remote cluster reference \"%s\" via %s.", ref.Name, ref.HostName)
		service.uilog_svc.Write(uiLogMsg)
	}
	return nil
}

func (service *RemoteClusterService) updateRemoteCluster(ref *metadata.RemoteClusterReference, revision interface{}) error {
	key := ref.Id
	value, err := json.Marshal(ref)
	if err != nil {
		return err
	}
	service.logger.Debugf("Remote cluster being changed: key=%v, value=%v\n", key, string(value))
	err = service.metadata_svc.Set(key, value, revision)
	if err != nil {
		return err
	}

	err = service.validateCache(ref)

	return err
}

func (service *RemoteClusterService) SetRemoteCluster(refName string, ref *metadata.RemoteClusterReference) error {
	service.logger.Infof("Setting remote cluster with refName %v\n", refName)

	err := service.ValidateRemoteCluster(ref)
	if err != nil {
		return err
	}

	oldRef, err := service.RemoteClusterByRefName(refName, false)
	if err != nil {
		return err
	}

	if ref.Id == oldRef.Id {
		// if the id of the remote cluster reference has not been changed, simply update the existing reference
		err = service.updateRemoteCluster(ref, oldRef.Revision)
		if err != nil {
			return err
		}
		return nil

	} else {
		// if id of the remote cluster reference has been changed, delete the existing reference and create a new one
		err = service.metadata_svc.DelWithCatalog(RemoteClustersCatalogKey, oldRef.Id, oldRef.Revision)
		if err != nil {
			return err
		}
		err = service.addRemoteCluster(ref)
		if err != nil {
			return err
		}
	}

	if service.uilog_svc != nil {
		nameChangeMsg := ""
		hostnameChangeMsg := ""
		if oldRef.Name != ref.Name {
			nameChangeMsg = fmt.Sprintf(" New name is \"%s\".", ref.Name)
		}
		if oldRef.HostName != ref.HostName {
			hostnameChangeMsg = fmt.Sprintf(" New contact point is %s.", ref.HostName)
		}
		uiLogMsg := fmt.Sprintf("Remote cluster reference \"%s\" updated.%s%s", oldRef.Name, nameChangeMsg, hostnameChangeMsg)
		service.uilog_svc.Write(uiLogMsg)
	}
	return nil
}

func (service *RemoteClusterService) DelRemoteCluster(refName string) (*metadata.RemoteClusterReference, error) {
	service.logger.Infof("Deleting remote cluster with reference name=%v\n", refName)
	ref, err := service.RemoteClusterByRefName(refName, false)
	if err != nil {
		return nil, err
	}
	key := ref.Id

	err = service.metadata_svc.DelWithCatalog(RemoteClustersCatalogKey, key, ref.Revision)
	if err != nil {
		return nil, err
	}

	service.invalidateCache(key)

	if service.uilog_svc != nil {
		uiLogMsg := fmt.Sprintf("Remote cluster reference \"%s\" known via %s removed.", ref.Name, ref.HostName)
		service.uilog_svc.Write(uiLogMsg)
	}
	return ref, nil
}

func (service *RemoteClusterService) RemoteClusters(refresh bool) (map[string]*metadata.RemoteClusterReference, error) {
	service.logger.Infof("Getting remote clusters")
	refs := make(map[string]*metadata.RemoteClusterReference, 0)

	entries, err := service.metadata_svc.GetAllMetadataFromCatalog(RemoteClustersCatalogKey)
	service.logger.Debugf("entries for remote clusters %v\n", entries)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		ref, err := service.constructRemoteClusterReference(entry.Value, entry.Rev)
		if err != nil {
			return nil, err
		}

		if refresh {
			ref, err = service.refresh(ref)
			if err != nil {
				return nil, err
			}
		}
		refs[entry.Key] = ref
	}

	return refs, nil
}

// validate remote cluster info and retrieve actual uuid
func (service *RemoteClusterService) ValidateRemoteCluster(ref *metadata.RemoteClusterReference) error {
	isEnterprise, err := service.xdcr_topology_svc.IsMyClusterEnterprise()
	if err != nil {
		return err
	}
	if ref.DemandEncryption && !isEnterprise {
		return wrapAsInvalidRemoteClusterError(errors.New("encryption can only be used in enterprise edition when the entire cluster is running at least 2.5 version of Couchbase Server"))
	}

	hostName := utils.GetHostName(ref.HostName)
	port, err := utils.GetPortNumber(ref.HostName)
	if err != nil {
		return wrapAsInvalidRemoteClusterError(errors.New(fmt.Sprintf("Failed to resolve address for \"%v\". The hostname may be incorrect or not resolvable.", ref.HostName)))
	}
		
	var hostAddr string
	var isInternalError bool
	if ref.DemandEncryption {
		hostAddr, err, isInternalError = service.httpsHostAddress(ref.HostName, ref.UserName, ref.Password)
		if err != nil {
			if isInternalError {
				return err
			} else {
				return wrapAsInvalidRemoteClusterError(errors.New(fmt.Sprintf("Could not connect to \"%v\" on port %v. This could be due to an incorrect host/port combination or a firewall in place between the servers.", hostName, port)))
			}
		}
	} else {
		hostAddr = ref.HostName
	}
	var poolsInfo map[string]interface{}

	if ref.DemandEncryption {
		hostAddr = utils.EnforcePrefix("https://", hostAddr)
	} else {
		hostAddr = utils.EnforcePrefix("http://", hostAddr)
	}

	err, statusCode := utils.QueryRestApiWithAuth(hostAddr, base.PoolsPath, false, ref.UserName, ref.Password, ref.Certificate, base.MethodGet, "", nil, 0, &poolsInfo, service.logger)
	service.logger.Infof("Result from validate remote cluster call: err=%v, statusCode=%v\n", err, statusCode)
	if err != nil || statusCode != http.StatusOK {
		if statusCode == http.StatusUnauthorized {
			return wrapAsInvalidRemoteClusterError(errors.New(fmt.Sprintf("Authentication failed. Verify username and password. Got HTTP status %v from REST call get to %v%v. Body was: []", statusCode, hostAddr, base.PoolsPath)))
		} else if !ref.DemandEncryption {
			// if encryption is not on, most likely the error is caused by incorrect hostname or firewall. 
			return wrapAsInvalidRemoteClusterError(errors.New(fmt.Sprintf("Could not connect to \"%v\" on port %v. This could be due to an incorrect host/port combination or a firewall in place between the servers.", hostName, port)))
		} else {
			// if encryption is on, several different errors could be returned here, e.g., invalid hostname, invalid certificate, certificate by unknown authority, etc.
			// just return the err
			return wrapAsInvalidRemoteClusterError(err)	
		}
	}

	// get remote cluster uuid from the map
	actualUuid, ok := poolsInfo[base.RemoteClusterUuid]
	if !ok {
		// should never get here
		return errors.New("Could not get uuid of remote cluster.")
	}
	actualUuidStr, ok := actualUuid.(string)
	if !ok {
		// should never get here
		return errors.New(fmt.Sprintf("uuid of remote cluster is of wrong type. Expected type: string; Actual type: %s", reflect.TypeOf(actualUuid)))
	}

	// update uuid in ref to real value
	ref.Uuid = actualUuidStr
	ref.Id = metadata.RemoteClusterRefId(ref.Uuid)

	return nil
}

func (service *RemoteClusterService) httpsHostAddress(hostName, userName, password string) (string, error, bool) {
	sslPort, err, isInternalError := utils.GetXDCRSSLPort(hostName, userName, password, service.logger)
	if err != nil {
		return "", err, isInternalError
	}

	hostNode := strings.Split(hostName, base.UrlPortNumberDelimiter)[0]
	newHostName := utils.GetHostAddr(hostNode, sslPort)
	return newHostName, nil, false
}

// this internal api differs from AddRemoteCluster in that it does not perform validation
func (service *RemoteClusterService) addRemoteCluster(ref *metadata.RemoteClusterReference) error {

	key := ref.Id
	value, err := json.Marshal(ref)
	if err != nil {
		return err
	}
	service.logger.Debugf("Remote cluster being added: key=%v, value=%v\n", key, string(value))
	err = service.metadata_svc.AddWithCatalog(RemoteClustersCatalogKey, key, value)
	if err != nil {
		return err
	}
	err = service.validateCache(ref)
	return err
}

func (service *RemoteClusterService) constructRemoteClusterReference(value []byte, rev interface{}) (*metadata.RemoteClusterReference, error) {
	if len(value) == 0 {
		return nil, errors.New("Failed to construct remote cluster value returned by metakv is empty")
	}

	ref := &metadata.RemoteClusterReference{}
	err := json.Unmarshal(value, ref)
	if err != nil {
		return nil, err
	}
	ref.Revision = rev

	return ref, err
}

func (service *RemoteClusterService) invalidateCache(key string) {
	service.cache_lock.Lock()
	defer service.cache_lock.Unlock()

	_, ok := service.cache_map[key]
	if ok {
		delete(service.cache_map, key)
		service.logger.Infof("Remote cluster reference %v is removed from the cache", key)

	}
}

func (service *RemoteClusterService) validateCache(ref *metadata.RemoteClusterReference) error {
	service.cache_lock.Lock()
	defer service.cache_lock.Unlock()

	service.logger.Infof("Remote Cluster reference %v need to be updated", ref.HostName)
	err := service.cacheRef(ref)
	if err != nil {
		service.logger.Infof("Didn't update the cache for remote cluster %v, err=%v\n", ref.Id, err)
	}
	return err
}

func (service *RemoteClusterService) cacheRef(ref *metadata.RemoteClusterReference) error {
	username, password, err := ref.MyCredentials()
	if err != nil {
		return err
	}
	connStr, err := ref.MyConnectionStr()
	if err != nil {
		return err
	}

	pool, err := utils.RemotePool(connStr, username, password)

	if err == nil {
		service.logger.Infof("Pool.Nodes=%v\n", pool.Nodes)

		nodes_connStrs := []string{}
		for _, node := range pool.Nodes {
			if node.Hostname != connStr {
				nodes_connStrs = append(nodes_connStrs, node.Hostname)
			}
		}

		ref_cache := &remoteClusterCache{key: ref.Id,
			nodes_connectionstr: nodes_connStrs,
			ref:                 ref}
		service.cache_map[ref.Id] = ref_cache

		service.logger.Infof("Remote cluster reference %v is cached, cache =%v", ref.Id, ref_cache)
	} else {
		service.logger.Infof("Remote cluster reference %v has a bad connectivity, didn't add to cache. err=%v", ref.Id, err)
	}
	return err
}

func (service *RemoteClusterService) getCache(key string) (*remoteClusterCache, bool) {
	service.cache_lock.RLock()
	defer service.cache_lock.RUnlock()
	cache, ok := service.cache_map[key]
	return cache, ok
}

func (service *RemoteClusterService) refresh(ref *metadata.RemoteClusterReference) (*metadata.RemoteClusterReference, error) {
	service.logger.Infof("Refresh remote cluster reference %v\n", ref.Id)

	err := service.validateCache(ref)

	if err == nil {
		return ref, nil
	}

	connStr, err := ref.MyConnectionStr()
	if err != nil {
		return nil, err
	}
	username, password, err := ref.MyCredentials()
	if err != nil {
		return nil, err
	}

	service.logger.Infof("Connstr %v in remote cluster reference failed to connect. Try to use alternative connStr", connStr)

	//the hostname in the remote cluster reference doesn't work
	//try on other connection strings in the cache
	ref_cache, ok := service.getCache(ref.Id)
	if !ok {
		service.logger.Errorf("Failed to connect to cluster reference %v - %v doesn't work, no alternative connStr\n", ref.Id, connStr)
		return nil, errors.New(fmt.Sprintf("Failed to connect to cluster reference %v\n", ref.Id))
	}

	service.logger.Infof("ref_cache=%v\n", ref_cache)


	var working_conn_str string = ""
	for _, alt_conn_str := range ref_cache.nodes_connectionstr {
		_, err = utils.RemotePool(alt_conn_str, username, password)
		if err == nil {
			working_conn_str = alt_conn_str
			break
		}
	}

	if working_conn_str != "" {
		service.logger.Infof("Found a working alternative connStr %v", working_conn_str)
		//update the ref
		ref.HostName = working_conn_str

		//persisted
		err = service.updateRemoteCluster(ref, ref.Revision)
		if err != nil {
			return ref, err
		}

		if err != nil {
			return ref, err
		}
		return ref, nil
	}

	return nil, errors.New(fmt.Sprintf("Failed to connect to cluster reference %v\n", ref.Id))

}

//get remote cluster name from remote cluster uuid. Return unknown if remote cluster cannot be found
func (service *RemoteClusterService) GetRemoteClusterNameFromClusterUuid(uuid string) string {
	remoteClusterRef, err := service.RemoteClusterByUuid(uuid, false)
	if err != nil || remoteClusterRef == nil {
		errMsg := fmt.Sprintf("Error getting the name of the remote cluster with uuid=%v.", uuid)
		if err != nil {
			errMsg += fmt.Sprintf(" err=%v", err)
		} else {
			errMsg += " The remote cluster may have been deleted."
		}
		service.logger.Error(errMsg)
		return base.UnknownRemoteClusterName
	}
	return remoteClusterRef.Name
}

// wrap/mark an error as invalid remote cluster error - by adding "invalid remote cluster" message to the front
func wrapAsInvalidRemoteClusterError(err error) error {
	return errors.New(InvalidRemoteClusterErrorMessage + err.Error())
}

func (service *RemoteClusterService) CheckAndUnwrapRemoteClusterError(err error) (bool, error) {
	if err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(errMsg, InvalidRemoteClusterErrorMessage) {
			return true, errors.New(errMsg[len(InvalidRemoteClusterErrorMessage):])
		} else if strings.HasPrefix(err.Error(), UnknownRemoteClusterErrorMessage) {
			return true, err
		} else {
			return false, err
		}		
	} else {
		return false, nil
	}
}