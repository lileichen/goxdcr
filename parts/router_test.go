// +build !pcre

package parts

import (
	"encoding/binary"
	"fmt"
	"github.com/couchbase/gomemcached"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/common"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/service_def"
	service_def_mocks "github.com/couchbase/goxdcr/service_def/mocks"
	utilities "github.com/couchbase/goxdcr/utils"
	UtilitiesMock "github.com/couchbase/goxdcr/utils/mocks"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"io/ioutil"
	"testing"
)

var dummyDownStream string = "dummy"

func setupBoilerPlateRouter() (routerId string, downStreamParts map[string]common.Part,
	routingMap map[uint16]string, crMode base.ConflictResolutionMode, loggerCtx *log.LoggerContext,
	utilsMock utilities.UtilsIface, throughputThrottlerSvc service_def.ThroughputThrottlerSvc,
	needToThrottle bool, expDelMode base.FilterExpDelType, collectionsManifestSvc *service_def_mocks.CollectionsManifestSvc,
	spec *metadata.ReplicationSpecification) {
	routerId = "routerUnitTest"

	downStreamParts = make(map[string]common.Part)
	downStreamParts[dummyDownStream] = nil
	routingMap = make(map[uint16]string)
	crMode = base.CRMode_RevId
	loggerCtx = log.DefaultLoggerContext
	utilsMock = &UtilitiesMock.UtilsIface{}
	throughputThrottlerSvc = &service_def_mocks.ThroughputThrottlerSvc{}
	needToThrottle = false
	expDelMode = base.FilterExpDelNone
	collectionsManifestSvc = &service_def_mocks.CollectionsManifestSvc{}
	spec, _ = metadata.NewReplicationSpecification("srcBucket", "srcBucketUUID", "targetClusterUUID", "tgtBucket", "tgtBucketUUID")

	return
}

func setupCollectionManifestsSvcRouter(collectionsManifestSvc *service_def_mocks.CollectionsManifestSvc) (pair metadata.CollectionsManifestPair) {
	manifestFileDir := "../metadata/testdata"
	manifestFileName := "provisionedManifest.json"

	data, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", manifestFileDir, manifestFileName))
	if err != nil {
		panic(err.Error())
	}
	manifest, err := metadata.NewCollectionsManifestFromBytes(data)
	if err != nil {
		panic(err.Error())
	}
	collectionsManifestSvc.On("GetLatestManifests", mock.Anything).Return(&manifest, &manifest, nil)
	pair.Source = &manifest
	pair.Target = &manifest
	return
}

func setupCollectionManifestsSvcRouterWithDefaultTarget(collectionsManifestSvc *service_def_mocks.CollectionsManifestSvc) {
	defaultManifest := metadata.NewDefaultCollectionsManifest()
	manifestFileDir := "../metadata/testdata"
	manifestFileName := "provisionedManifest.json"

	data, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", manifestFileDir, manifestFileName))
	if err != nil {
		panic(err.Error())
	}
	manifest, err := metadata.NewCollectionsManifestFromBytes(data)
	if err != nil {
		panic(err.Error())
	}
	collectionsManifestSvc.On("GetLatestManifests", mock.Anything).Return(&manifest, &defaultManifest, nil)
}

func TestRouterRouteFunc(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterRouteFunc =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	uprEvent, err := RetrieveUprFile("./testdata/uprEventDeletion.json")
	assert.Nil(err)
	assert.NotNil(uprEvent)
	wrappedEvent := &base.WrappedUprEvent{UprEvent: uprEvent}

	// Deletion does not contain any flags
	wrappedMCRequest, err := router.ComposeMCRequest(wrappedEvent)
	assert.Nil(err)
	assert.NotNil(wrappedMCRequest)
	checkUint := binary.BigEndian.Uint32(wrappedMCRequest.Req.Extras[0:4])
	assert.False(checkUint&base.IS_EXPIRATION > 0)

	// Expiration has to set a flag of 0x10
	uprEvent, err = RetrieveUprFile("./testdata/uprEventExpiration.json")
	assert.Nil(err)
	assert.NotNil(uprEvent)
	wrappedEvent = &base.WrappedUprEvent{UprEvent: uprEvent}

	wrappedMCRequest, err = router.ComposeMCRequest(wrappedEvent)
	assert.Nil(err)
	assert.NotNil(wrappedMCRequest)
	checkUint = binary.BigEndian.Uint32(wrappedMCRequest.Req.Extras[24:28])
	assert.True(checkUint&base.IS_EXPIRATION > 0)

	fmt.Println("============== Test case end: TestRouterRouteFunc =================")
}

func TestRouterInitialNone(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterInitialNone =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	assert.Equal(base.FilterExpDelNone, router.expDelMode.Get())

	fmt.Println("============== Test case end: TestRouterInitialNone =================")
}

func TestRouterSkipDeletion(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterSkipDeletion =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelSkipDeletes

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	assert.NotEqual(base.FilterExpDelNone, router.expDelMode.Get())

	delEvent, err := RetrieveUprFile("./testdata/uprEventDeletion.json")
	assert.Nil(err)
	assert.NotNil(delEvent)

	expEvent, err := RetrieveUprFile("./testdata/uprEventExpiration.json")
	assert.Nil(err)
	assert.NotNil(expEvent)

	shouldContinue := router.ProcessExpDelTTL(delEvent)
	assert.False(shouldContinue)

	shouldContinue = router.ProcessExpDelTTL(expEvent)
	assert.True(shouldContinue)

	fmt.Println("============== Test case end: TestRouterSkipDeletion =================")
}

func TestRouterSkipExpiration(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterSkipExpiration =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelSkipExpiration

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	assert.NotEqual(base.FilterExpDelNone, router.expDelMode.Get())

	delEvent, err := RetrieveUprFile("./testdata/uprEventDeletion.json")
	assert.Nil(err)
	assert.NotNil(delEvent)

	expEvent, err := RetrieveUprFile("./testdata/uprEventExpiration.json")
	assert.Nil(err)
	assert.NotNil(expEvent)

	shouldContinue := router.ProcessExpDelTTL(delEvent)
	assert.True(shouldContinue)

	shouldContinue = router.ProcessExpDelTTL(expEvent)
	assert.False(shouldContinue)

	fmt.Println("============== Test case end: TestRouterSkipExpiration =================")
}

func TestRouterSkipDeletesStripTTL(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterSkipExpiryStripTTL =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelSkipExpiration | base.FilterExpDelStripExpiration

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	assert.NotEqual(base.FilterExpDelNone, router.expDelMode.Get())

	// delEvent contains expiry
	delEvent, err := RetrieveUprFile("./testdata/uprEventDeletion.json")
	assert.Nil(err)
	assert.NotNil(delEvent)

	expEvent, err := RetrieveUprFile("./testdata/uprEventExpiration.json")
	assert.Nil(err)
	assert.NotNil(expEvent)

	assert.NotEqual(0, int(delEvent.Expiry))
	shouldContinue := router.ProcessExpDelTTL(delEvent)
	assert.True(shouldContinue)
	assert.Equal(0, int(delEvent.Expiry))

	shouldContinue = router.ProcessExpDelTTL(expEvent)
	assert.False(shouldContinue)

	fmt.Println("============== Test case end: TestRouterSkipExpiryStripTTL =================")
}

func TestRouterExpDelAllMode(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterExpDelAllMode =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelAll

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	assert.NotEqual(base.FilterExpDelNone, router.expDelMode.Get())

	delEvent, err := RetrieveUprFile("./testdata/uprEventDeletion.json")
	assert.Nil(err)
	assert.NotNil(delEvent)

	mutEvent, err := RetrieveUprFile("./testdata/perfDataExpiry.json")
	assert.Nil(err)
	assert.NotNil(mutEvent)

	expEvent, err := RetrieveUprFile("./testdata/uprEventExpiration.json")
	assert.Nil(err)
	assert.NotNil(expEvent)

	assert.NotEqual(0, int(mutEvent.Expiry))
	shouldContinue := router.ProcessExpDelTTL(mutEvent)
	assert.True(shouldContinue)
	assert.Equal(0, int(mutEvent.Expiry))

	shouldContinue = router.ProcessExpDelTTL(expEvent)
	assert.False(shouldContinue)

	shouldContinue = router.ProcessExpDelTTL(delEvent)
	assert.False(shouldContinue)
	fmt.Println("============== Test case end: TestRouterExpDelAllMode =================")
}

var testDir string = "../metadata/testData/"

var targetv8 string = testDir + "diffTargetv8.json"
var targetv9 string = testDir + "diffTargetv9.json"

func TestRouterManifestChange(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterManifestChange =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	setupCollectionManifestsSvcRouter(collectionsManifestSvc)

	expDelMode = base.FilterExpDelAll

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	assert.Nil(router.Start())
	collectionsRouter := router.collectionsRouting[dummyDownStream]
	assert.NotNil(collectionsRouter)

	// routing updater receiver
	newRoutingUpdater := func(info CollectionsRoutingInfo) error {
		// This test will show no broken map but one fixed map
		assert.Equal(0, len(info.BrokenMap))
		assert.Equal(1, len(info.BackfillMap))
		return nil
	}
	collectionsRouter.routingUpdater = newRoutingUpdater

	data, _ := ioutil.ReadFile(targetv8)
	targetv8Manifest, _ := metadata.NewCollectionsManifestFromBytes(data)

	data, _ = ioutil.ReadFile(targetv9)
	targetv9Manifest, _ := metadata.NewCollectionsManifestFromBytes(data)

	assert.Nil(collectionsRouter.handleNewManifestChanges(&targetv8Manifest))
	// Force a manual brokenmap. V9 will have the following fixed
	implicitNamespace := &base.CollectionNamespace{"S2", "col3"}
	collectionsRouter.brokenDenyMtx.Lock()
	collectionsRouter.brokenMapping.AddSingleMapping(implicitNamespace, implicitNamespace)
	collectionsRouter.brokenDenyMtx.Unlock()

	assert.Nil(collectionsRouter.handleNewManifestChanges(&targetv9Manifest))

	collectionsRouter.brokenDenyMtx.Lock()
	assert.Equal(0, len(collectionsRouter.brokenMapping))
	collectionsRouter.brokenDenyMtx.Unlock()

	fmt.Println("============== Test case end: TestRouterManifestChange =================")
}

func mockTargetCollectionDNE(collectionsManifestSvc *service_def_mocks.CollectionsManifestSvc) {
	// Target DNE simply means an empty manifest
	defaultManifest := metadata.NewDefaultCollectionsManifest()
	setupCollectionManifestsSvcRouter(collectionsManifestSvc)
	collectionsManifestSvc.On("GetSpecificTargetManifest", mock.Anything, mock.Anything).Return(&defaultManifest, nil)
}

func TestRouterTargetCollectionDNE(t *testing.T) {
	fmt.Println("============== Test case start: TargetCollectionDNE =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelAll

	mockTargetCollectionDNE(collectionsManifestSvc)

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	assert.Nil(router.Start())
	collectionsRouter := router.collectionsRouting[dummyDownStream]
	assert.NotNil(collectionsRouter)

	// routing updater receiver
	newRoutingUpdater := func(info CollectionsRoutingInfo) error {
		// This test will show one broken map and no fixed map
		assert.Equal(1, len(info.BrokenMap))
		assert.Equal(0, len(info.BackfillMap))
		// For now say persist is fine
		return nil
	}

	var ignoreCnt int
	ignoreFunc := func(*base.WrappedMCRequest) {
		ignoreCnt++
	}

	collectionsRouter.routingUpdater = newRoutingUpdater
	collectionsRouter.ignoreDataFunc = ignoreFunc

	implicitNamespace := &base.CollectionNamespace{"S2", "col3"}
	dummyData := &base.WrappedMCRequest{ColNamespace: implicitNamespace,
		ColInfo: &base.TargetCollectionInfo{},
	}

	assert.Equal(base.ErrorIgnoreRequest, router.RouteCollection(dummyData, dummyDownStream))
	assert.Equal(1, ignoreCnt)
	fmt.Println("============== Test case end: TargetCollectionDNE =================")
}

func TestRouterTargetCollectionDNEPersistErr(t *testing.T) {
	fmt.Println("============== Test case start: TargetCollectionDNEPersistErr =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelAll

	mockTargetCollectionDNE(collectionsManifestSvc)

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	assert.NotNil(router)

	assert.Nil(router.Start())
	collectionsRouter := router.collectionsRouting[dummyDownStream]
	assert.NotNil(collectionsRouter)

	// routing updater receiver
	newRoutingUpdater := func(info CollectionsRoutingInfo) error {
		// This test will show one broken map and no fixed map
		assert.Equal(1, len(info.BrokenMap))
		assert.Equal(0, len(info.BackfillMap))
		// For now say persist is not fine
		return fmt.Errorf("Dummy persist err")
	}

	var ignoreCnt int
	ignoreFunc := func(*base.WrappedMCRequest) {
		ignoreCnt++
	}

	collectionsRouter.routingUpdater = newRoutingUpdater
	collectionsRouter.ignoreDataFunc = ignoreFunc

	implicitNamespace := &base.CollectionNamespace{"S2", "col3"}
	dummyData := &base.WrappedMCRequest{ColNamespace: implicitNamespace,
		ColInfo: &base.TargetCollectionInfo{},
	}

	// Even if persist has problem, routeCollection should return a non-nil error to prevent forwarding to xmem
	assert.Equal(base.ErrorIgnoreRequest, router.RouteCollection(dummyData, dummyDownStream))
	// The ignore count should be 0 to indicate that throughSeqno will not move foward
	assert.Equal(0, ignoreCnt)
	fmt.Println("============== Test case end: TargetCollectionDNEPersistErr =================")
}

func TestRouterExplicitMode(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterExplicitMode =================")
	defer fmt.Println("============== Test case end: TestRouterExplicitMode =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	mappingMode := spec.Settings.GetCollectionModes()
	mappingMode.SetExplicitMapping(true)
	rules := spec.Settings.GetCollectionsRoutingRules()
	rules["S1"] = "S2"
	updatedMap := make(map[string]interface{})
	updatedMap[metadata.CollectionsMgtMultiKey] = mappingMode
	updatedMap[metadata.CollectionsMappingRulesKey] = rules
	_, errMap := spec.Settings.UpdateSettingsFromMap(updatedMap)
	assert.Equal(0, len(errMap))

	setupCollectionManifestsSvcRouterWithDefaultTarget(collectionsManifestSvc)

	router, err := NewRouter(routerId, spec, downStreamParts, routingMap, crMode, loggerCtx, utilsMock, throughputThrottlerSvc, needToThrottle, expDelMode, collectionsManifestSvc, nil, nil)

	assert.Nil(err)
	modes := router.collectionModes.Get()
	assert.True(modes.IsExplicitMapping())

	collectionsRouter := router.collectionsRouting[dummyDownStream]
	assert.NotNil(collectionsRouter)
	assert.Nil(collectionsRouter.Start())

	// routing updater receiver
	newRoutingUpdater := func(info CollectionsRoutingInfo) error {
		return nil
	}

	var ignoreCnt int
	ignoreFunc := func(*base.WrappedMCRequest) {
		ignoreCnt++
	}

	collectionsRouter.routingUpdater = newRoutingUpdater
	collectionsRouter.ignoreDataFunc = ignoreFunc

	sourceNs := &base.CollectionNamespace{"S1", "col1"}
	mcReq := &gomemcached.MCRequest{
		Key:    []byte("testKey"),
		Keylen: len("testKey"),
	}
	dummyData := &base.WrappedMCRequest{
		Req:          mcReq,
		ColNamespace: sourceNs,
		ColInfo:      &base.TargetCollectionInfo{},
	}
	err = router.RouteCollection(dummyData, dummyDownStream)
	// Right now target cluster gave back "default manifest"
	assert.Equal(base.ErrorIgnoreRequest, err)
	assert.Equal(1, ignoreCnt)
	assert.Equal(1, len(collectionsRouter.brokenMapping))
	// Second one should be already ignored
	err = router.RouteCollection(dummyData, dummyDownStream)
	assert.Equal(base.ErrorRequestAlreadyIgnored, err)
	assert.Equal(2, ignoreCnt)
	assert.Equal(1, len(collectionsRouter.brokenMapping))

	collectionsManifestSvc = &service_def_mocks.CollectionsManifestSvc{}
	pair := setupCollectionManifestsSvcRouter(collectionsManifestSvc)
	collectionsRouter.collectionsManifestSvc = collectionsManifestSvc

	// TODO - change this to use actual update call
	collectionsRouter.explicitMappings, err = metadata.NewCollectionNamespaceMappingFromRules(pair, mappingMode, rules)
	assert.Nil(err)

	err = router.RouteCollection(dummyData, dummyDownStream)
	assert.Nil(err)
	//targetNs := &base.CollectionNamespace{"S2", "col1"}

	//// Even if persist has problem, routeCollection should return a non-nil error to prevent forwarding to xmem
	//assert.Equal(base.ErrorIgnoreRequest, router.RouteCollection(dummyData, dummyDownStream))
	//// The ignore count should be 0 to indicate that throughSeqno will not move foward
}
