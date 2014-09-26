// mock services for unit tests
package common

import (
	"fmt"
	//	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/base"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/metadata"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/utils"
	//	"github.com/couchbaselabs/go-couchbase"
	"github.com/couchbaselabs/go-couchbase"
)

var options struct {
	sourceBucket    string // source bucket
	targetBucket    string //target bucket
	connectStr      string //connect string
	numConnPerKV    int    // number of connections per source KV node
	numOutgoingConn int    // number of connections to target cluster
	username        string //username
	password        string //password
	maxVbno         int    // maximum number of vbuckets
}

func SetTestOptions(sourceBucket, targetBucket, connectStr, username, password string, numConnPerKV, numOutgoingConn int) {
	options.connectStr = connectStr
	options.sourceBucket = sourceBucket
	options.targetBucket = targetBucket
	options.username = username
	options.password = password
	options.numConnPerKV = numConnPerKV
	options.numOutgoingConn = numOutgoingConn 
}

type MockMetadataSvc struct {
}

func (mock_meta_svc *MockMetadataSvc) ReplicationSpec(replicationId string) (*metadata.ReplicationSpecification, error) {
	spec := metadata.NewReplicationSpecification(options.connectStr, options.sourceBucket, options.connectStr, options.targetBucket, "")
	settings := spec.Settings()
	settings.SetTargetNozzlesPerNode(options.numOutgoingConn)
	settings.SetSourceNozzlesPerNode(options.numConnPerKV)
	return spec, nil
}

func (mock_meta_svc *MockMetadataSvc) AddReplicationSpec(spec metadata.ReplicationSpecification) error {
	return nil
}

func (mock_meta_svc *MockMetadataSvc) SetReplicationSpec(spec metadata.ReplicationSpecification) error {
	return nil
}

func (mock_meta_svc *MockMetadataSvc) DelReplicationSpec(replicationId string) error {
	return nil
}

type MockClusterInfoSvc struct {
}

func (mock_ci_svc *MockClusterInfoSvc) GetClusterConnectionStr(ClusterUUID string) (string, error) {
	return options.connectStr, nil
}

func (mock_ci_svc *MockClusterInfoSvc) GetMyActiveVBuckets(ClusterUUID string, bucketName string, NodeId string) ([]uint16, error) {
	sourceCluster, err := mock_ci_svc.GetClusterConnectionStr(ClusterUUID)
	if err != nil {
		return nil, err
	}
	b, err := utils.Bucket(sourceCluster, bucketName, options.username, options.password)
	if err != nil {
		return nil, err
	}

	// in test env, there should be only one kv in bucket server list
	kvaddr := b.VBServerMap().ServerList[0]

	m, err := b.GetVBmap([]string{kvaddr})
	if err != nil {
		return nil, err
	}

	vbList := m[kvaddr]

	return vbList, nil
}

func (mock_ci_svc *MockClusterInfoSvc) GetServerList(ClusterUUID string, bucketName string) ([]string, error) {
	cluster, err := mock_ci_svc.GetClusterConnectionStr(ClusterUUID)
	if err != nil {
		return nil, err
	}
	bucket, err := utils.Bucket(cluster, bucketName, options.username, options.password)
	if err != nil {
		return nil, err
	}

	// in test env, there should be only one kv in bucket server list
	serverlist := bucket.VBServerMap().ServerList

	return serverlist, nil
}

func (mock_ci_svc *MockClusterInfoSvc) GetServerVBucketsMap(ClusterUUID string, bucketName string) (map[string][]uint16, error) {
	cluster, err := mock_ci_svc.GetClusterConnectionStr(ClusterUUID)
	fmt.Printf("cluster=%s\n", cluster)
	if err != nil {
		return nil, err
	}
	bucket, err := utils.Bucket(cluster, bucketName, options.username, options.password)
	if err != nil {
		return nil, err
	}
	fmt.Printf("ServerList=%v\n", bucket.VBServerMap().ServerList)
	serverVBMap, err := bucket.GetVBmap(bucket.VBServerMap().ServerList)
	fmt.Printf("ServerVBMap=%v\n", serverVBMap)
	return serverVBMap, err
}

func (mock_ci_svc *MockClusterInfoSvc) IsNodeCompatible(node string, version string) (bool, error) {
	return true, nil
}

func (mock_ci_svc *MockClusterInfoSvc) GetBucket(clusterUUID, bucketName string) (*couchbase.Bucket, error) {
	clusterConnStr, err := mock_ci_svc.GetClusterConnectionStr(clusterUUID)
	if err != nil {
		return nil, err
	}
	return utils.Bucket(clusterConnStr, bucketName, options.username, options.password)
}


type MockXDCRTopologySvc struct {
}

func (mock_top_svc *MockXDCRTopologySvc) MyHost() (string, error) {
	return options.connectStr, nil
}

func (mock_top_svc *MockXDCRTopologySvc) MyAdminPort() (uint16, error) {
	return 0, nil
}

func (mock_top_svc *MockXDCRTopologySvc) MyKVNodes() ([]string, error) {
	mock_ci_svc := &MockClusterInfoSvc{}
	nodes, err := mock_ci_svc.GetServerList(options.connectStr, "default")
	return nodes, err
}

func (mock_top_svc *MockXDCRTopologySvc) XDCRTopology() (map[string]uint16, error) {
	retmap := make(map[string]uint16)
	return retmap, nil
}

func (mock_top_svc *MockXDCRTopologySvc) XDCRCompToKVNodeMap() (map[string][]string, error) {
	retmap := make(map[string][]string)
	return retmap, nil
}

func (mock_top_svc *MockXDCRTopologySvc) MyCluster() (string, error) {
	return options.connectStr, nil
}

type MockReplicationSettingsSvc struct {
}

func (mock_repl_settings_svc *MockReplicationSettingsSvc) GetReplicationSettings() (*metadata.ReplicationSettings, error) {
	return metadata.DefaultSettings(), nil
}
	
func (mock_repl_settings_svc *MockReplicationSettingsSvc) SetReplicationSettings(*metadata.ReplicationSettings) error {
	return nil
}