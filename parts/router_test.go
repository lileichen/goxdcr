// +build !pcre

package parts

import (
	"encoding/binary"
	"fmt"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/common"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/service_def"
	service_def_mocks "github.com/couchbase/goxdcr/service_def/mocks"
	utilities "github.com/couchbase/goxdcr/utils"
	UtilitiesMock "github.com/couchbase/goxdcr/utils/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dummyDownStream string = "dummy"

func setupBoilerPlateRouter() (routerId string, downStreamParts map[string]common.Part,
	routingMap map[uint16]string, crMode base.ConflictResolutionMode, loggerCtx *log.LoggerContext,
	req_creater ReqCreator, utilsMock utilities.UtilsIface, throughputThrottlerSvc service_def.ThroughputThrottlerSvc,
	needToThrottle bool, expDelMode base.FilterExpDelType, collectionsManifestSvc service_def.CollectionsManifestSvc,
	spec *metadata.ReplicationSpecification) {
	routerId = "routerUnitTest"

	downStreamParts = make(map[string]common.Part)
	downStreamParts[dummyDownStream] = nil
	routingMap = make(map[uint16]string)
	crMode = base.CRMode_RevId
	loggerCtx = log.DefaultLoggerContext
	utilsMock = &UtilitiesMock.UtilsIface{}
	throughputThrottlerSvc = &service_def_mocks.ThroughputThrottlerSvc{}
	req_creater = nil
	needToThrottle = false
	expDelMode = base.FilterExpDelNone
	collectionsManifestSvc = &service_def_mocks.CollectionsManifestSvc{}
	spec, _ = metadata.NewReplicationSpecification("srcBucket", "srcBucketUUID", "targetClusterUUID", "tgtBucket", "tgtBucketUUID")

	return
}

func TestRouterRouteFunc(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterRouteFunc =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		req_creater, utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	router, err := NewRouter(routerId, spec, downStreamParts,
		routingMap, crMode, loggerCtx, req_creater, utilsMock, throughputThrottlerSvc, needToThrottle,
		expDelMode, collectionsManifestSvc, nil /*objRecycler*/)

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
		req_creater, utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	router, err := NewRouter(routerId, spec, downStreamParts,
		routingMap, crMode, loggerCtx, req_creater, utilsMock, throughputThrottlerSvc, needToThrottle,
		expDelMode, collectionsManifestSvc, nil /*objRecycler*/)

	assert.Nil(err)
	assert.NotNil(router)

	assert.Equal(base.FilterExpDelNone, router.expDelMode.Get())

	fmt.Println("============== Test case end: TestRouterInitialNone =================")
}

func TestRouterSkipDeletion(t *testing.T) {
	fmt.Println("============== Test case start: TestRouterSkipDeletion =================")
	assert := assert.New(t)

	routerId, downStreamParts, routingMap, crMode, loggerCtx,
		req_creater, utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelSkipDeletes

	router, err := NewRouter(routerId, spec, downStreamParts,
		routingMap, crMode, loggerCtx, req_creater, utilsMock, throughputThrottlerSvc, needToThrottle,
		expDelMode, collectionsManifestSvc, nil /*objRecycler*/)

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
		req_creater, utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelSkipExpiration

	router, err := NewRouter(routerId, spec, downStreamParts,
		routingMap, crMode, loggerCtx, req_creater, utilsMock, throughputThrottlerSvc, needToThrottle,
		expDelMode, collectionsManifestSvc, nil /*objRecycler*/)

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
		req_creater, utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelSkipExpiration | base.FilterExpDelStripExpiration

	router, err := NewRouter(routerId, spec, downStreamParts,
		routingMap, crMode, loggerCtx, req_creater, utilsMock, throughputThrottlerSvc, needToThrottle,
		expDelMode, collectionsManifestSvc, nil /*objRecycler*/)

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
		req_creater, utilsMock, throughputThrottlerSvc,
		needToThrottle, expDelMode, collectionsManifestSvc, spec := setupBoilerPlateRouter()

	expDelMode = base.FilterExpDelAll

	router, err := NewRouter(routerId, spec, downStreamParts,
		routingMap, crMode, loggerCtx, req_creater, utilsMock, throughputThrottlerSvc, needToThrottle,
		expDelMode, collectionsManifestSvc, nil /*objRecycler*/)

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
