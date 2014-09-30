package main

import (
	"flag"
	"fmt"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr/pipeline_manager"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/metadata"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/parts"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/replication_manager"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/utils"
	//	c "github.com/couchbase/indexing/secondary/common"
	"github.com/couchbaselabs/go-couchbase"
	"log"
	"net/http"
	"os"
	"time"
)

import _ "net/http/pprof"

var options struct {
	source_bucket           string // source bucket
	target_bucket           string //target bucket
	source_cluster_addr     string //source connect string
	target_cluster_addr     string //target connect string
	source_cluster_username string //source cluster username
	source_cluster_password string //source cluster password
	target_cluster_username string //target cluster username
	target_cluster_password string //target cluster password
	target_bucket_password  string //target bucket password
	nozzles_per_node_source int    // number of nozzles per source node
	nozzles_per_node_target int    // number of nozzles per target node
	//	username        string //username
	//	password        string //password
	//	maxVbno         int    // maximum number of vbuckets
}

const (
	NUM_SOURCE_CONN = 2
	NUM_TARGET_CONN = 3
)

func argParse() {

	flag.StringVar(&options.source_cluster_addr, "source_cluster_addr", "127.0.0.1:9000",
		"source cluster address")
	flag.StringVar(&options.source_bucket, "source_bucket", "default",
		"bucket to replicate from")
	flag.StringVar(&options.target_cluster_addr, "target_cluster_addr", "127.0.0.1:9000",
		"target cluster address")
	flag.StringVar(&options.target_bucket, "target_bucket", "target",
		"bucket to replicate to")
	flag.StringVar(&options.source_cluster_username, "source_cluster_username", "Administrator",
		"user name to use for logging into source cluster")
	flag.StringVar(&options.source_cluster_password, "source_cluster_password", "welcome",
		"password to use for logging into source cluster")
	flag.StringVar(&options.target_cluster_username, "target_cluster_username", "Administrator",
		"user name to use for logging into target cluster")
	flag.StringVar(&options.target_cluster_password, "target_cluster_password", "welcome",
		"password to use for logging into target cluster")
	flag.StringVar(&options.target_bucket_password, "target_bucket_password", "welcome",
		"password to use for accessing target bucket")
	flag.IntVar(&options.nozzles_per_node_source, "nozzles_per_node_source", NUM_SOURCE_CONN,
		"number of nozzles per source node")
	flag.IntVar(&options.nozzles_per_node_target, "nozzles_per_node_target", NUM_TARGET_CONN,
		"number of nozzles per target node")

	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage : %s [OPTIONS]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	go func() {
		log.Println("Try to start pprof...")
		err := http.ListenAndServe("localhost:7000", nil)
		if err != nil {
			panic(err)
		} else {
			log.Println("Http server for pprof is started")
		}
	}()

	//	c.SetLogLevel(c.LogLevelTrace)
	fmt.Println("Start Testing ...")
	argParse()
	setup()
	test()
	verify()
}

func setup() {
	flushTargetBkt()
	fmt.Println("Finish setup")
	replication_manager.Initialize(NewMockMetadataSvc(), &mockClusterInfoSvc{}, &mockXDCRTopologySvc{}, nil)
	return
}

func test() {
	topic, err := replication_manager.CreateReplication(options.source_cluster_addr, options.source_bucket, options.target_cluster_addr, options.target_bucket, options.target_bucket_password, nil)
	if err != nil {
		fail(fmt.Sprintf("%v", err))
	}
	time.Sleep(1 * time.Second)
	replication_manager.PauseReplication(topic)

	err = replication_manager.SetPipelineLogLevel(topic, "Error")
	if err != nil {
		fail(fmt.Sprintf("%v", err))
	}
	fmt.Printf("Replication %s is paused\n", topic)
	time.Sleep(100 * time.Millisecond)
	err = replication_manager.ResumeReplication(topic)
	if err != nil {
		fail(fmt.Sprintf("%v", err))
	}
	fmt.Printf("Replication %s is resumed\n", topic)
	time.Sleep(2 * time.Second)
	summary(topic)

	err = replication_manager.DeleteReplication(topic)
	if err != nil {
		fail(fmt.Sprintf("%v", err))
	}
	fmt.Printf("Replication %s is deleted\n", topic)
	time.Sleep(4 * time.Second)

}

func summary(topic string) {
	pipeline := pipeline_manager.Pipeline(topic)
	targetNozzles := pipeline.Targets()
	for _, targetNozzle := range targetNozzles {
		fmt.Println(targetNozzle.(*parts.XmemNozzle).StatusSummary())
	}
}
func fail(msg string) {
	panic(fmt.Sprintf("TEST FAILED - %s", msg))
}

func verify() {
	sourceDocCount := getDocCounts(options.source_cluster_addr, options.source_bucket, "")
	targetDocCount := getDocCounts(options.target_cluster_addr, options.target_bucket, options.target_bucket_password)

	if sourceDocCount == targetDocCount {
		fmt.Println("TEST SUCCESS")
	} else {
		fmt.Printf("TEST FAILED\n")
		fmt.Printf("Source doc count=%d; target doc count=%d\n", sourceDocCount, targetDocCount)
	}
}

func getDocCounts(clusterAddress string, bucketName string, password string) int {
	output := &utils.CouchBucket{}
	baseURL, err := couchbase.ParseURL("http://" + clusterAddress)

	if err == nil {
		err = utils.QueryRestAPI(baseURL,
			"/pools/default/buckets/"+bucketName,
			bucketName,
			password,
			"GET",
			output, nil)
	}
	if err != nil {
		panic(err)
	}
	log.Printf("name=%s itemCount=%d\n, out=%v\n", output.Name, output.Stat.ItemCount, output)
	return output.Stat.ItemCount

}

func flushTargetBkt() {
	//flush the target bucket
	baseURL, err := couchbase.ParseURL("http://" + options.target_bucket + ":" + options.target_bucket_password + "@" + options.target_cluster_addr)

	if err == nil {
		err = utils.QueryRestAPI(baseURL,
			"/pools/default/buckets/"+options.target_bucket+"/controller/doFlush",
			options.target_cluster_username,
			options.target_cluster_password,
			"POST",
			nil, nil)
	}

	if err != nil {
		log.Printf("Setup error=%v\n", err)
	} else {
		log.Println("Setup is done")
	}

}

type mockMetadataSvc struct {
	specs map[string]metadata.ReplicationSpecification
}

func NewMockMetadataSvc() *mockMetadataSvc {
	return &mockMetadataSvc{specs: make(map[string]metadata.ReplicationSpecification)}
}
func (mock_meta_svc *mockMetadataSvc) ReplicationSpec(replicationId string) (*metadata.ReplicationSpecification, error) {
	spec, ok := mock_meta_svc.specs[replicationId]
	if !ok {
		spec_ptr := metadata.NewReplicationSpecification(options.source_cluster_addr, options.source_bucket, options.target_cluster_addr, options.target_bucket, "")
		settings := spec_ptr.Settings()
		settings.SetTargetNozzlesPerNode(options.nozzles_per_node_target)
		settings.SetSourceNozzlesPerNode(options.nozzles_per_node_source)
		mock_meta_svc.specs[replicationId] = *spec_ptr
		return spec_ptr, nil
	}else {
		return &spec, nil
	}
}

func (mock_meta_svc *mockMetadataSvc) AddReplicationSpec(spec metadata.ReplicationSpecification) error {
	mock_meta_svc.specs[spec.Id()] = spec
	return nil
}

func (mock_meta_svc *mockMetadataSvc) SetReplicationSpec(spec metadata.ReplicationSpecification) error {
	mock_meta_svc.specs[spec.Id()] = spec
	return nil
}

func (mock_meta_svc *mockMetadataSvc) DelReplicationSpec(replicationId string) error {
	delete(mock_meta_svc.specs, replicationId)
	return nil
}

type mockClusterInfoSvc struct {
}

func (mock_ci_svc *mockClusterInfoSvc) GetClusterConnectionStr(ClusterUUID string) (string, error) {
	return options.source_cluster_addr, nil
}

func (mock_ci_svc *mockClusterInfoSvc) GetMyActiveVBuckets(ClusterUUID string, bucketName string, NodeId string) ([]uint16, error) {
	sourceCluster, err := mock_ci_svc.GetClusterConnectionStr(ClusterUUID)
	if err != nil {
		return nil, err
	}
	b, err := utils.Bucket(sourceCluster, bucketName, options.source_cluster_username, options.source_cluster_password)
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

func (mock_ci_svc *mockClusterInfoSvc) GetServerList(ClusterUUID string, bucketName string) ([]string, error) {
	cluster, err := mock_ci_svc.GetClusterConnectionStr(ClusterUUID)
	if err != nil {
		return nil, err
	}
	bucket, err := utils.Bucket(cluster, bucketName, options.source_cluster_username, options.source_cluster_password)
	if err != nil {
		return nil, err
	}

	// in test env, there should be only one kv in bucket server list
	serverlist := bucket.VBServerMap().ServerList

	return serverlist, nil
}

func (mock_ci_svc *mockClusterInfoSvc) GetServerVBucketsMap(ClusterUUID string, bucketName string) (map[string][]uint16, error) {
	cluster, err := mock_ci_svc.GetClusterConnectionStr(ClusterUUID)
	fmt.Printf("cluster=%s\n", cluster)
	if err != nil {
		return nil, err
	}
	bucket, err := utils.Bucket(cluster, bucketName, options.source_cluster_username, options.source_cluster_password)
	if err != nil {
		return nil, err
	}
	fmt.Printf("ServerList=%v\n", bucket.VBServerMap().ServerList)
	serverVBMap, err := bucket.GetVBmap(bucket.VBServerMap().ServerList)
	fmt.Printf("ServerVBMap=%v\n", serverVBMap)
	return serverVBMap, err
}

func (mock_ci_svc *mockClusterInfoSvc) IsNodeCompatible(node string, version string) (bool, error) {
	return true, nil
}

func (mock_ci_svc *mockClusterInfoSvc) GetBucket(clusterUUID, bucketName string) (*couchbase.Bucket, error) {
	clusterConnStr, err := mock_ci_svc.GetClusterConnectionStr(clusterUUID)
	if err != nil {
		return nil, err
	}
	return utils.Bucket(clusterConnStr, bucketName, options.source_cluster_username, options.source_cluster_password)
}

type mockXDCRTopologySvc struct {
}

func (mock_top_svc *mockXDCRTopologySvc) MyHost() (string, error) {
	return options.source_cluster_addr, nil
}

func (mock_top_svc *mockXDCRTopologySvc) MyAdminPort() (uint16, error) {
	return 0, nil
}

func (mock_top_svc *mockXDCRTopologySvc) MyKVNodes() ([]string, error) {
	mock_ci_svc := &mockClusterInfoSvc{}
	nodes, err := mock_ci_svc.GetServerList(options.source_cluster_addr, "default")
	return nodes, err
}

func (mock_top_svc *mockXDCRTopologySvc) XDCRTopology() (map[string]uint16, error) {
	retmap := make(map[string]uint16)
	return retmap, nil
}

func (mock_top_svc *mockXDCRTopologySvc) XDCRCompToKVNodeMap() (map[string][]string, error) {
	retmap := make(map[string][]string)
	return retmap, nil
}

func (mock_top_svc *mockXDCRTopologySvc) MyCluster() (string, error) {
	return options.source_cluster_addr, nil
}