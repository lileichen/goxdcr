// Copyright (c) 2013 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package pipeline_utils

import (
	"errors"
	"github.com/couchbase/goxdcr/common"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/parts"
	"github.com/couchbase/goxdcr/service_def"
	"strings"
)

var ErrorNoSourceKV = errors.New("Invalid configuration. No source kv node is found.")

func GetSourceVBListPerPipeline(pipeline common.Pipeline) []uint16 {
	ret := []uint16{}
	sourceNozzles := pipeline.Sources()
	for _, sourceNozzle := range sourceNozzles {
		ret = append(ret, sourceNozzle.(*parts.DcpNozzle).GetVBList()...)
	}
	return ret
}

func GetSourceVBMapForReplication(cluster_info_svc service_def.ClusterInfoSvc, xdcr_topology_svc service_def.XDCRCompTopologySvc,
	spec *metadata.ReplicationSpecification, logger *log.CommonLogger) (map[string][]uint16, error) {
	kv_vb_map := make(map[string][]uint16)
	server_vbmap, err := cluster_info_svc.GetServerVBucketsMap(xdcr_topology_svc, spec.SourceBucketName)
	if err != nil {
		return nil, err
	}

	logger.Debugf("server_vbmap=%v\n", server_vbmap)
	nodes, err := xdcr_topology_svc.MyKVNodes()
	if err != nil {
		logger.Errorf("Failed to get my KV nodes, err=%v\n", err)
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, ErrorNoSourceKV
	}

	for _, node := range nodes {
		var kvaddr string = ""
		var vbnos []uint16
		// iterate through serverVBMap and look for server addr that starts with "kvHost:"
		for kvaddr_iter, vbnos_iter := range server_vbmap {
			if kvaddr_iter == node {
				kvaddr = kvaddr_iter
				vbnos = vbnos_iter
				break
			}
		}
		if kvaddr != "" {
			kv_vb_map[kvaddr] = vbnos
		}
	}
	return kv_vb_map, nil
}

// checks if target cluster supports ssl over memcached
func HasSSLOverMemSupport(cluster_info_svc service_def.ClusterInfoSvc, targetClusterRef *metadata.RemoteClusterReference) (bool, error) {
	return cluster_info_svc.IsClusterCompatible(targetClusterRef, []int{3, 0})
}

// the element here can be ComponentEventListener, ComponentEventHandler, or other elements in pipeiline
func GetElementIdFromName(pipeline common.Pipeline, name string) string {
	return pipeline.Topic() + "_" + name
}

func GetElementNameFromId(id string) string {
	parts := strings.Split(id, "_")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func GetAsyncComponentEventListener(listenerMap map[string]common.AsyncComponentEventListener, name string) common.AsyncComponentEventListener {
	for id, listener := range listenerMap {
		if name == GetElementNameFromId(id) {
			return listener
		}
	}
	return nil
}
