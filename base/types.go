// Copyright 2013-Present Couchbase, Inc.
//
// Use of this software is governed by the Business Source License included in
// the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
// file, in accordance with the Business Source License, use of this software
// will be governed by the Apache License, Version 2.0, included in the file
// licenses/APL2.txt.

package base

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	mrand "math/rand"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/snappy"

	"github.com/couchbase/gomemcached"
	mcc "github.com/couchbase/gomemcached/client"
	"github.com/google/uuid"
)

type SettingDef struct {
	Data_type reflect.Type
	Required  bool
}

func NewSettingDef(data_type reflect.Type, bReq bool) *SettingDef {
	return &SettingDef{Data_type: data_type, Required: bReq}
}

type SettingDefinitions map[string]*SettingDef

type ErrorMap map[string]error

// Not thread safe - callers need to sync
func (errMap *ErrorMap) AddErrors(otherMap ErrorMap) {
	MergeErrorMaps(*errMap, otherMap, true /*overwrite*/)
}

func (errMap *ErrorMap) HasError(err error) bool {
	for _, v := range *errMap {
		if v == err {
			return true
		}
	}
	return false
}

type SSLPortMap map[string]uint16

type SettingsError struct {
	err_map ErrorMap
}

func (se SettingsError) Error() string {
	var buffer bytes.Buffer
	for key, err := range se.err_map {
		errStr := fmt.Sprintf("setting=%s; err=%s\n", key, err.Error())
		buffer.WriteString(errStr)
	}
	return buffer.String()
}

func NewSettingsError() *SettingsError {
	return &SettingsError{make(ErrorMap)}
}

func (se SettingsError) Add(key string, err error) {
	se.err_map[key] = err
}

type CollectionsManifestIdPair struct {
	SourceManifestId uint64
	TargetManifestId uint64
}

// Ensure that this current c can accomodate the other range
func (c *CollectionsManifestIdPair) Accomodate(other CollectionsManifestIdPair) {
	if c.SourceManifestId > other.SourceManifestId {
		c.SourceManifestId = other.SourceManifestId
	}
	if c.TargetManifestId < other.TargetManifestId {
		c.TargetManifestId = other.TargetManifestId
	}
}

// timestamp for a specific vb
type VBTimestamp struct {
	Vbno          uint16
	Vbuuid        uint64
	Seqno         uint64
	SnapshotStart uint64
	SnapshotEnd   uint64
	ManifestIDs   CollectionsManifestIdPair
}

var emptyVBts VBTimestamp

func (vbts *VBTimestamp) String() string {
	return fmt.Sprintf("[vbno=%v, uuid=%v, seqno=%v, sn_start=%v, sn_end=%v, srcManId=%v, tgtManId=%v]",
		vbts.Vbno, vbts.Vbuuid, vbts.Seqno, vbts.SnapshotStart, vbts.SnapshotEnd, vbts.ManifestIDs.SourceManifestId, vbts.ManifestIDs.TargetManifestId)
}

func (vbts VBTimestamp) IsEmpty() bool {
	return vbts == emptyVBts
}

func (vbts *VBTimestamp) SameAs(other *VBTimestamp) bool {
	if vbts == nil && other == nil {
		return true
	} else if vbts == nil && other != nil {
		return false
	} else if vbts != nil && other == nil {
		return false
	}

	return vbts.Vbno == other.Vbno && vbts.Vbuuid == other.Vbuuid && vbts.Seqno == other.Seqno &&
		vbts.SnapshotStart == other.SnapshotStart && vbts.SnapshotEnd == other.SnapshotEnd &&
		vbts.ManifestIDs.SourceManifestId == other.ManifestIDs.SourceManifestId &&
		vbts.ManifestIDs.TargetManifestId == other.ManifestIDs.TargetManifestId
}

func (vbts VBTimestamp) Clone() VBTimestamp {
	var clonedTs VBTimestamp
	clonedTs.Vbno = vbts.Vbno
	clonedTs.Vbuuid = vbts.Vbuuid
	clonedTs.Seqno = vbts.Seqno
	clonedTs.SnapshotStart = vbts.SnapshotStart
	clonedTs.SnapshotEnd = vbts.SnapshotEnd
	clonedTs.ManifestIDs.SourceManifestId = vbts.ManifestIDs.SourceManifestId
	clonedTs.ManifestIDs.TargetManifestId = vbts.ManifestIDs.TargetManifestId
	return clonedTs
}

func (vbts *VBTimestamp) Sanitize() {
	if vbts.Seqno == 0 {
		vbts.SnapshotStart = 0
		vbts.SnapshotEnd = 0
	} else {
		if vbts.SnapshotStart > vbts.Seqno || vbts.SnapshotStart == 0 {
			vbts.SnapshotStart = vbts.Seqno
		}
		// Sanitize should only be used on a non-checkpoint VBTimestamp
		// Since the timestamp is not used as a checkpoint to resume, just set the endTs to be the same
		// so that DCP resume can work and let DCP roll us back if necessary
		// Let's not pretend that there is a legit snapshot end that we know of
		if vbts.SnapshotEnd != vbts.Seqno {
			vbts.SnapshotEnd = vbts.Seqno
		}
	}
}

type ClusterConnectionInfoProvider interface {
	MyConnectionStr() (string, error)
	// returns username, password, http auth mechanism, certificate, whether certificate contains SAN, client certificate, client key
	MyCredentials() (string, string, HttpAuthMech, []byte, bool, []byte, []byte, error)
}

type ReplicationInfo struct {
	Id        string
	StatsMap  map[string]interface{}
	ErrorList []ErrorInfo
}

type EventInfoType int

const (
	HighPriorityMsg       EventInfoType = iota
	LowPriorityMsg        EventInfoType = iota
	PersistentMsg         EventInfoType = iota
	BrokenMappingInfoType EventInfoType = iota
)

func (e EventInfoType) CanDismiss() bool {
	switch e {
	case HighPriorityMsg:
		return false
	case LowPriorityMsg:
		return true
	case BrokenMappingInfoType:
		return true
	case PersistentMsg:
		return false
	default:
		panic("Implement me")
	}
}

func (e EventInfoType) String() string {
	switch e {
	case HighPriorityMsg:
		return "HighPriorityMsg"
	case LowPriorityMsg:
		return "LowPriorityMsg"
	case BrokenMappingInfoType:
		return "BrokenMappingInfo"
	case PersistentMsg:
		return "PersistentMsg"
	default:
		return "?? (EventInfoType)"
	}
}

type EventInfoComparator interface {
	SameAs(other interface{}) bool
}

type EventsMap struct {
	EventsMap   map[int64]interface{}
	eventMapMtx *sync.RWMutex // Not exporting because of JSON marshaller
}

func NewEventsMap() EventsMap {
	return EventsMap{
		EventsMap:   make(map[int64]interface{}),
		eventMapMtx: &sync.RWMutex{},
	}
}

func (c *EventsMap) IsNil() bool {
	return c.eventMapMtx == nil && c.EventsMap == nil
}

func (c *EventsMap) IsEmpty() bool {
	if !c.IsNil() {
		c.eventMapMtx.RLock()
		defer c.eventMapMtx.RUnlock()
		return len(c.EventsMap) == 0
	}
	return true
}

// Assumes RLock his held - only safe for 1-writer multi-readers
func (c *EventsMap) UpgradeLock() func() {
	c.eventMapMtx.RUnlock()
	c.eventMapMtx.Lock()
	return func() {
		c.eventMapMtx.Unlock()
		c.eventMapMtx.RLock()
	}
}

// Not concurrency safe
func (c *EventsMap) Init() {
	c.eventMapMtx = &sync.RWMutex{}
	c.EventsMap = make(map[int64]interface{})
}

func (c *EventsMap) GetRWLock() *sync.RWMutex {
	return c.eventMapMtx
}

func (c *EventsMap) MarshalJSON() ([]byte, error) {
	if c == nil {
		return nil, ErrorNilPtr
	}

	c.eventMapMtx.RLock()
	defer c.eventMapMtx.RUnlock()

	return json.Marshal(c.EventsMap)
}

func (c *EventsMap) UnmarshalJSON(b []byte) error {
	if c == nil {
		return ErrorNilPtr
	}

	c.eventMapMtx.Lock()
	defer c.eventMapMtx.Unlock()

	return json.Unmarshal(b, &c.EventsMap)
}

func (c *EventsMap) String() string {
	if c == nil {
		return ""
	}

	c.eventMapMtx.RLock()
	defer c.eventMapMtx.RUnlock()

	return fmt.Sprintf("%v", c.EventsMap)
}

func (c *EventsMap) ContainsEvent(eventId int) bool {
	if c == nil {
		return false
	}

	c.eventMapMtx.RLock()
	defer c.eventMapMtx.RUnlock()

	for k, v := range c.EventsMap {
		if int(k) == eventId {
			return true
		}
		if event, ok := v.(*EventInfo); ok {
			if event.ContainsEvent(eventId) {
				return true
			}
		}
	}
	return false
}

func (c *EventsMap) SameAs(other EventsMap) bool {
	if c == nil {
		return false
	}

	c.eventMapMtx.RLock()
	defer c.eventMapMtx.RUnlock()

	other.eventMapMtx.RLock()
	defer other.eventMapMtx.RUnlock()

	if len(c.EventsMap) != len(other.EventsMap) {
		return false
	}

	for k, v := range c.EventsMap {
		compareV, ok := other.EventsMap[k]
		if !ok {
			return false
		}
		vComparator := v.(EventInfoComparator)
		if !vComparator.SameAs(compareV) {
			return false
		}
	}

	return true
}

func (c *EventsMap) Len() int {
	if c == nil {
		return 0
	}

	c.eventMapMtx.RLock()
	defer c.eventMapMtx.RUnlock()
	return len(c.EventsMap)
}

/*
 * Usage for EventInfo data structure:
 *
 * For high and low priority messages:
 *      - ErrDesc == message to be communicated
 *		- EventExtras not used
 *
 * For broken map events:
 * Level 0: Complete brokenmap event
 *      - ErrDesc == ""
 *		- Use EventExtras with KV pairs. Each KV pair contains a Level 1 event
 * Level 1: Source Scope level event
 *      - ErrDesc == a single source scope name
 *		- Use EventExtras with KV pairs. Each KV pair contains a Level 2 event
 * Level 2: Source namespace event
 *      - ErrDesc == "scopeName.collectionName"
 *		- Use EventExtras with KV pairs. Each KV pair contains a Level 3 event
 * Level 3: Target namespace event
 *      - ErrDesc == "targetScopeName.targetCollectionName"
 *		- EventExtras not used
 *
 */

type EventInfo struct {
	EventId   int64
	EventType EventInfoType
	EventDesc string

	// Keyed by subEventId
	// Value depends on EventType
	EventExtras EventsMap

	hintMtx *sync.RWMutex
	hint    interface{}
}

func NewEventInfo() *EventInfo {
	return &EventInfo{
		EventExtras: NewEventsMap(),
		hintMtx:     &sync.RWMutex{},
	}
}

func (e *EventInfo) SetHint(ns interface{}) {
	e.hintMtx.Lock()
	defer e.hintMtx.Unlock()
	e.hint = ns
}

func (e *EventInfo) GetHint() interface{} {
	e.hintMtx.RLock()
	defer e.hintMtx.RUnlock()
	return e.hint
}

func (e *EventInfo) SameAs(otherRaw interface{}) bool {
	otherE, ok := otherRaw.(EventInfo)
	if ok {
		return e.EventId == otherE.EventId &&
			e.EventType == otherE.EventType &&
			e.EventDesc == otherE.EventDesc &&
			e.EventExtras.SameAs(otherE.EventExtras)
	}

	otherEPtr, ok := otherRaw.(*EventInfo)
	if ok {
		return e.EventId == otherEPtr.EventId &&
			e.EventType == otherEPtr.EventType &&
			e.EventDesc == otherEPtr.EventDesc &&
			e.EventExtras.SameAs(otherEPtr.EventExtras)
	}

	return false
}

func (e *EventInfo) String() string {
	return fmt.Sprintf("EventId: %v type: %v desc: %v extras: %v hint: %v", e.EventId, e.EventType.String(), e.EventDesc, e.EventExtras.String(), e.hint)
}

func (e *EventInfo) ContainsEvent(eventId int) bool {
	if int(e.EventId) == eventId {
		return true
	}

	return e.EventExtras.ContainsEvent(eventId)
}

func (e *EventInfo) GetSubEvent(eventId int) (*EventInfo, error) {
	e.EventExtras.eventMapMtx.RLock()
	defer e.EventExtras.eventMapMtx.RUnlock()

	for eId, eventRaw := range e.EventExtras.EventsMap {
		eventInfo, ok := eventRaw.(*EventInfo)
		if int(eId) == eventId {
			if ok {
				return eventInfo, nil
			} else {
				return nil, fmt.Errorf("Wrong type: %v", reflect.TypeOf(eventRaw))
			}
		}
		if ok {
			recurseSearch, reErr := eventInfo.GetSubEvent(eventId)
			if reErr == nil && recurseSearch != nil {
				return recurseSearch, reErr
			}
		}
	}
	return nil, ErrorNotFound
}

func (e *EventInfo) GetLegacyErrMsg() string {
	if e == nil {
		return ""
	}
	switch e.EventType {
	case HighPriorityMsg:
		return e.EventDesc
	case LowPriorityMsg:
		return e.EventDesc
	case PersistentMsg:
		return e.EventDesc
	case BrokenMappingInfoType:
		hint, isString := e.GetHint().(string)
		if isString {
			// Only the top level brokenMapping event should have a string hint
			return hint
		} else {
			return ""
		}
	default:
		return "?? (EventInfo)"
	}
}

type ErrorInfo struct {
	*EventInfo
	// Time is the number of nano seconds elapsed since 1/1/1970 UTC
	Time     int64
	ErrorMsg string
}

// Internally, XDCR eventIDs are kept as int64
// Externally, since JSON doesn't play nice with int64, we will use reg int32
// The Getter function here should ensure that if it exceeds MaxInt32, it starts over from
// the beginning
func GetEventIdFromWell(eventIdWell *int64) int64 {
	for {
		retVal := atomic.AddInt64(eventIdWell, 1)
		if retVal > math.MaxInt32 {
			swapped := atomic.CompareAndSwapInt64(eventIdWell, atomic.LoadInt64(eventIdWell), 0)
			if swapped {
				return atomic.AddInt64(eventIdWell, 1)
			}
		} else {
			return retVal
		}
	}
}

func NewErrorInfo(time int64, errorMsg string, eventIdWell *int64) ErrorInfo {
	eventInfo := NewEventInfo()
	eventInfo.EventType = HighPriorityMsg
	eventInfo.EventId = GetEventIdFromWell(eventIdWell)

	errInfo := ErrorInfo{
		EventInfo: eventInfo,
		Time:      time,
		ErrorMsg:  errorMsg,
	}
	return errInfo
}

func NewErrorInfoFromEventInfo(event *EventInfo, t int64) ErrorInfo {
	return ErrorInfo{
		EventInfo: event,
		ErrorMsg:  event.GetLegacyErrMsg(),
		Time:      t,
	}
}

// These should be RO once they are created
// because they are potentially used as constant keys for CollectionNamespaceMapping
type CollectionNamespace struct {
	ScopeName      string
	CollectionName string
}

var DefaultCollectionNamespace = CollectionNamespace{
	ScopeName:      DefaultScopeCollectionName,
	CollectionName: DefaultScopeCollectionName,
}

// The specificStr should follow the format of: "<scope>.<collection>"
func NewCollectionNamespaceFromString(specificStr string) (CollectionNamespace, error) {
	if !CollectionNamespaceRegex.MatchString(specificStr) {
		if len(specificStr) > MaxCollectionNameBytes {
			return CollectionNamespace{}, ErrorLengthExceeded
		}
		return CollectionNamespace{}, ErrorInvalidColNamespaceFormat
	} else {
		names := CollectionNamespaceRegex.FindStringSubmatch(specificStr)
		if len(names) != 3 {
			return CollectionNamespace{}, fmt.Errorf("Invalid capture for CollectionNameSpaces")
		}
		for i := 1; i <= 2; i++ {
			valid := CollectionNameValidationRegex.MatchString(names[i])
			if !valid {
				return CollectionNamespace{}, ErrorInvalidColNamespaceFormat
			}
			if len(names[i]) > MaxCollectionNameBytes {
				return CollectionNamespace{}, ErrorLengthExceeded
			}
		}

		return CollectionNamespace{
			ScopeName:      names[1],
			CollectionName: names[2],
		}, nil
	}
}

func NewDefaultCollectionNamespace() *CollectionNamespace {
	ns, _ := NewCollectionNamespaceFromString(fmt.Sprintf("%v%v%v", DefaultScopeCollectionName, ScopeCollectionDelimiter, DefaultScopeCollectionName))
	return &ns
}

// To index string showcases a simple to understand of "scope.collection"
func (c CollectionNamespace) ToIndexString() string {
	return fmt.Sprintf("%v%v%v", c.ScopeName, ScopeCollectionDelimiter, c.CollectionName)
}

func (c *CollectionNamespace) IsDefault() bool {
	return c.ScopeName == DefaultScopeCollectionName && c.CollectionName == DefaultScopeCollectionName
}

func (c *CollectionNamespace) IsSystemScope() bool {
	return c.ScopeName == SystemScopeName
}

func (c *CollectionNamespace) IsEmpty() bool {
	return c.ScopeName == "" && c.CollectionName == ""
}

func (c CollectionNamespace) LessThan(other CollectionNamespace) bool {
	if c.ScopeName == other.ScopeName {
		return c.CollectionName < other.CollectionName
	} else {
		return c.ScopeName < other.ScopeName
	}
}

func (c CollectionNamespace) IsSameAs(other CollectionNamespace) bool {
	return c.ScopeName == other.ScopeName && c.CollectionName == other.CollectionName
}

func (c CollectionNamespace) Clone() CollectionNamespace {
	return CollectionNamespace{
		ScopeName:      c.ScopeName,
		CollectionName: c.CollectionName,
	}
}

type CollectionNamespacePtrList []*CollectionNamespace

func (c CollectionNamespacePtrList) Len() int           { return len(c) }
func (c CollectionNamespacePtrList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CollectionNamespacePtrList) Less(i, j int) bool { return (*(c[i])).LessThan(*(c[j])) }

func SortCollectionNamespacePtrList(list CollectionNamespacePtrList) CollectionNamespacePtrList {
	sort.Sort(list)
	return list
}

type WrappedUprEvent struct {
	UprEvent          *mcc.UprEvent
	ColNamespace      *CollectionNamespace
	Flags             WrappedUprEventFlag
	ByteSliceGetter   func(uint64) ([]byte, error)
	DecompressedValue []byte
}

type WrappedUprEventFlag uint64

const (
	WrappedUprCollectionDNE        = 0x1
	WrappedUprValueUseDecompressed = 0x2
)

func (w WrappedUprEventFlag) IsSet() bool {
	return uint64(w) > uint64(0)
}

func (w WrappedUprEventFlag) CollectionDNE() bool {
	return w&WrappedUprCollectionDNE > 0
}

func (w *WrappedUprEventFlag) SetCollectionDNE() {
	*w |= WrappedUprCollectionDNE
}

func (w WrappedUprEventFlag) ShouldUseDecompressedValue() bool {
	return w&WrappedUprValueUseDecompressed > 0
}

func (w *WrappedUprEventFlag) SetShouldUseDecompressedValue() {
	*w |= WrappedUprValueUseDecompressed
}

type TargetCollectionInfo struct {
	// The manifestID used when retrying
	ManifestId uint64
	ColId      uint32
	// Temporary storage space to avoid memmove
	ColIDPrefixedKey []byte
	// The len is needed because slice returned from datapool can have garbage
	ColIDPrefixedKeyLen int

	// For migration, since given WrappedMCRequest's SrcColNamespace, there is no way to look up the actual
	// target namespace given migration rules, this is needed for error handling
	// For non-migration use cases, this is not to be used
	TargetNamespace *CollectionNamespace
}

type WrappedMCRequest struct {
	Seqno                      uint64
	Req                        *gomemcached.MCRequest
	Start_time                 time.Time
	UniqueKey                  string
	SrcColNamespace            *CollectionNamespace
	SrcColNamespaceMtx         sync.RWMutex
	ColInfo                    *TargetCollectionInfo
	ColInfoMtx                 sync.RWMutex
	SlicesToBeReleasedByXmem   [][]byte
	SlicesToBeReleasedByRouter [][]byte
	SlicesToBeReleasedMtx      sync.Mutex
	ReqBytesCachedMtx          sync.RWMutex
	ReqBytesCached             []byte

	// If a single source mutation is translated to multiple target requests, the additional ones are listed here
	SiblingReqs    []*WrappedMCRequest
	SiblingReqsMtx sync.RWMutex
	RetryCRCount   int
	Cloned         bool
	ClonedSyncCh   chan bool
}

func (req *WrappedMCRequest) GetReqBytes() []byte {
	req.ReqBytesCachedMtx.RLock()
	retBytes := req.ReqBytesCached
	if retBytes != nil {
		req.ReqBytesCachedMtx.RUnlock()
		return retBytes
	}
	req.ReqBytesCachedMtx.RUnlock()
	return req.populateReqCache()
}

func (req *WrappedMCRequest) populateReqCache() []byte {
	req.ReqBytesCachedMtx.Lock()
	defer req.ReqBytesCachedMtx.Unlock()
	req.ReqBytesCached = req.Req.Bytes()
	return req.ReqBytesCached
}

func (req *WrappedMCRequest) UpdateReqBytes() {
	req.ReqBytesCachedMtx.RLock()
	curReqBytesLen := len(req.ReqBytesCached)
	req.ReqBytesCachedMtx.RUnlock()

	if curReqBytesLen == 0 {
		// Never established
		req.populateReqCache()
	} else {
		req.ReqBytesCachedMtx.Lock()
		defer req.ReqBytesCachedMtx.Unlock()
		if cap(req.ReqBytesCached) >= req.Req.Size() {
			// Repopulate cache
			req.Req.BytesPreallocated(req.ReqBytesCached)
		} else {
			// Too small, reallocate
			req.ReqBytesCached = req.Req.Bytes()
		}
		// Trim just in case
		req.ReqBytesCached = req.ReqBytesCached[0:req.Req.Size()]
	}
}

func (req *WrappedMCRequest) ConstructUniqueKey() {
	var buffer bytes.Buffer
	buffer.Write(req.Req.Key)
	buffer.Write(req.Req.Extras[8:16])
	req.UniqueKey = buffer.String()

	req.SiblingReqsMtx.RLock()
	if len(req.SiblingReqs) > 0 {
		for _, sibling := range req.SiblingReqs {
			if sibling != nil {
				sibling.ConstructUniqueKey()
			}
		}
	}
	req.SiblingReqsMtx.RUnlock()
}

// Returns 0 if no collection is used
func (req *WrappedMCRequest) GetManifestId() uint64 {
	req.ColInfoMtx.RLock()
	defer req.ColInfoMtx.RUnlock()
	if req != nil && req.ColInfo != nil {
		return req.ColInfo.ManifestId
	} else {
		return 0
	}
}

func (req *WrappedMCRequest) GetSourceCollectionNamespace() *CollectionNamespace {
	req.SrcColNamespaceMtx.RLock()
	defer req.SrcColNamespaceMtx.RUnlock()
	return req.SrcColNamespace
}

func (req *WrappedMCRequest) GetBodySize() int {
	if req == nil || req.Req == nil {
		return 0
	}

	return len(req.Req.Body)
}

func (req *WrappedMCRequest) GetUncompressedBodySize() int {
	if req == nil || req.Req == nil {
		return 0
	}

	if req.Req.DataType&mcc.SnappyDataType > 0 {
		decodedLen, err := snappy.DecodedLen(req.Req.Body)
		if err != nil {
			return req.GetBodySize()
		}
		return decodedLen
	} else {
		return req.GetBodySize()
	}
}

type McRequestMap map[string]*WrappedMCRequest

type MCResponseMap map[string]*gomemcached.MCResponse

type MetadataChangeListener interface {
	Id() string
	Start() error
}

type ObjectWithLock struct {
	Object interface{}
	Lock   *sync.RWMutex
}

type SeqnoWithLock struct {
	seqno uint64
	lock  *sync.RWMutex
}

func NewSeqnoWithLock() *SeqnoWithLock {
	return &SeqnoWithLock{0, &sync.RWMutex{}}
}

// the following methods should be sufficient for typical get/set operations
func (seqno_obj *SeqnoWithLock) GetSeqno() uint64 {
	return seqno_obj.getSeqno(true)
}

func (seqno_obj *SeqnoWithLock) SetSeqno(seqno uint64) {
	seqno_obj.setSeqno(seqno, true)
}

// the following methods allow more than one get/set operations to be done in one transation
func (seqno_obj *SeqnoWithLock) Lock() {
	seqno_obj.lock.Lock()
}

func (seqno_obj *SeqnoWithLock) Unlock() {
	seqno_obj.lock.Unlock()
}

func (seqno_obj *SeqnoWithLock) GetSeqnoWithoutLock() uint64 {
	return seqno_obj.getSeqno(false)
}

func (seqno_obj *SeqnoWithLock) SetSeqnoWithoutLock(seqno uint64) {
	seqno_obj.setSeqno(seqno, false)
}

func (seqno_obj *SeqnoWithLock) getSeqno(lock bool) uint64 {
	if lock {
		seqno_obj.lock.RLock()
		defer seqno_obj.lock.RUnlock()
	}
	return seqno_obj.seqno
}

func (seqno_obj *SeqnoWithLock) setSeqno(seqno uint64, lock bool) {
	if lock {
		seqno_obj.lock.Lock()
		defer seqno_obj.lock.Unlock()
	}
	seqno_obj.seqno = seqno
}

// Callback function, which typically is a method in metadata service, which translates metakv call back parameters into metadata objects
// The function may do something addtional that is specific to the metadata service, e.g., caching the new metadata
type MetadataServiceCallback func(path string, value []byte, rev interface{}) error

// Callback function for the handling of metadata changed event
type MetadataChangeHandlerCallback func(metadataId string, oldMetadata interface{}, newMetadata interface{}) error

// Certain callbacks are done in parallel and so may need a waitGroup to ensure synchronization is done
type MetadataChangeHandlerCallbackWithWg func(metadataId string, oldMetadata interface{}, newMetadata interface{}, wg *sync.WaitGroup) error

type MetadataChangeHandlerPriority int

// When callbacks are to be called concurrently, certain callbacks may have higher priority to be called than another
// when certain metadata is being created or deleted
const (
	MetadataChangeHighPrioriy MetadataChangeHandlerPriority = iota
	MetadataChangeMedPrioriy  MetadataChangeHandlerPriority = iota
	MetadataChangeLowPrioriy  MetadataChangeHandlerPriority = iota
)

type DataObjRecycler func(topic string, dataObj *WrappedMCRequest)

type VBErrorType int

const (
	// vb error caused by source topology change
	VBErrorType_Source VBErrorType = iota
	// vb error caused by target topology change
	VBErrorType_Target VBErrorType = iota
)

type VBErrorEventAdditional struct {
	Vbno      uint16
	Error     error
	ErrorType VBErrorType
}

type ConflictResolutionMode int

const (
	CRMode_RevId  ConflictResolutionMode = iota
	CRMode_LWW    ConflictResolutionMode = iota
	CRMode_Custom ConflictResolutionMode = iota
)

// stack implementation
type Stack []interface{}

func (s *Stack) Empty() bool        { return len(*s) == 0 }
func (s *Stack) Peek() interface{}  { return (*s)[len(*s)-1] }
func (s *Stack) Push(i interface{}) { (*s) = append((*s), i) }
func (s *Stack) Pop() interface{} {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
}

type InterfaceMap map[string]interface{}

// Shallow clone
func (in InterfaceMap) Clone() InterfaceMap {
	clonedMap := make(InterfaceMap)
	for k, v := range in {
		clonedMap[k] = v
	}
	return clonedMap
}

// Redact only works on keys
func (in InterfaceMap) Redact() InterfaceMap {
	for k, v := range in {
		if !IsStringRedacted(k) {
			in[TagUD(k)] = v
			delete(in, k)
		}
	}
	return in
}

// bucket info map
type BucketKVVbMap map[string][]uint16

// Usually, server vbucket maps are consisted of host:directPort
// To support External hosts, we need to manually look up external hosts if they exist, as well as
// the external direct port if they exist.
func (kvVbMap BucketKVVbMap) ReplaceInternalWithExternalHosts(translationMap map[string]string) {
	var keysToDelete []string
	// bucketKVVBMap's key is usually [internalHostname:internalPort]
	// The translated map will look like [internalHostname:internalPort] -> [externalHostName:externalPort/internalPort]
	for internalHostAndPort, vbSlice := range kvVbMap {
		var externalHostAndPort string
		externalHostAndPort, internalExists := translationMap[internalHostAndPort]
		if internalExists && internalHostAndPort != externalHostAndPort {
			// replace the internal key with external key and the same value
			kvVbMap[externalHostAndPort] = vbSlice
			keysToDelete = append(keysToDelete, internalHostAndPort)
		}
	}
	for _, key := range keysToDelete {
		delete(kvVbMap, key)
	}
}

// Compression Section
type CompressionType int

type UserAuthMode int

const (
	// no user auth
	UserAuthModeNone UserAuthMode = iota
	// use implicit local user auth
	UserAuthModeLocal UserAuthMode = iota
	// basic user auth using passed in username and password - so far this applies to user in remote clusters only
	UserAuthModeBasic UserAuthMode = iota
)

// http authentication mode
type HttpAuthMech int

const (
	// plain http authentication
	HttpAuthMechPlain HttpAuthMech = iota
	// scram sha authentication
	HttpAuthMechScramSha HttpAuthMech = iota
	// https authentication
	HttpAuthMechHttps HttpAuthMech = iota
)

// for logging
func (httpAuthMech HttpAuthMech) String() string {
	switch httpAuthMech {
	case HttpAuthMechPlain:
		return "Plain"
	case HttpAuthMechScramSha:
		return "ScramSha"
	case HttpAuthMechHttps:
		return "Https"
	default:
		return "Unknown"
	}
}

// StringPair holds a pair of strings
type StringPair [2]string

func (sp StringPair) GetFirstString() string {
	return sp[0]
}

func (sp StringPair) GetSecondString() string {
	return sp[1]
}

// StringPairList is used mainly to facilitate sorting
// StringPairList can be sorted by first string through sort.Sort(list)
type StringPairList []StringPair

func (spl StringPairList) Len() int           { return len(spl) }
func (spl StringPairList) Swap(i, j int)      { spl[i], spl[j] = spl[j], spl[i] }
func (spl StringPairList) Less(i, j int) bool { return spl[i][0] < spl[j][0] }

// get a list of string1
func (spl StringPairList) GetListOfFirstString() []string {
	result := make([]string, len(spl))

	for i := 0; i < len(spl); i++ {
		result[i] = spl[i][0]
	}

	return result
}

func ShuffleStringPairList(list StringPairList) {
	r := mrand.New(mrand.NewSource(time.Now().Unix()))
	// Start at the end of the slice, go backwards and scramble
	for i := len(list); i > 1; i-- {
		randIndex := r.Intn(i)
		// Swap values and continue until we're done
		if (i - 1) != randIndex {
			list[i-1], list[randIndex] = list[randIndex], list[i-1]
		}
	}
}

func DeepCopyStringPairList(list StringPairList) StringPairList {
	if list == nil {
		return nil
	}

	out := make([]StringPair, len(list))
	copy(out, list)
	return out
}

type ComponentError struct {
	ComponentId string
	Err         error
}

func FormatErrMsg(errCh chan ComponentError) map[string]error {
	// use -1 to indicate no upper limit on number of errors to return
	return FormatErrMsgWithUpperLimit(errCh, -1)
}

func FormatErrMsgWithUpperLimit(errCh chan ComponentError, maxNumberOfErrorsToTrack int) map[string]error {
	errMap := make(map[string]error)
	for {
		select {
		case componentErr := <-errCh:
			errMap[componentErr.ComponentId] = componentErr.Err
			if maxNumberOfErrorsToTrack > 0 && len(errMap) >= maxNumberOfErrorsToTrack {
				return errMap
			}
		default:
			return errMap
		}
	}
	return errMap
}

type PriorityType int

const (
	PriorityTypeHigh   PriorityType = iota
	PriorityTypeMedium PriorityType = iota
	PriorityTypeLow    PriorityType = iota
)

// priority shown on UI and rest api
const (
	PriorityHighString   string = "High"
	PriorityMediumString string = "Medium"
	PriorityLowString    string = "Low"
)

func PriorityTypeFromStr(priorityStr string) (PriorityType, error) {
	var priority PriorityType
	switch priorityStr {
	case PriorityHighString:
		priority = PriorityTypeHigh
	case PriorityMediumString:
		priority = PriorityTypeMedium
	case PriorityLowString:
		priority = PriorityTypeLow
	default:
		return -1, errors.New(fmt.Sprintf("%v is not a valid priority", priorityStr))
	}
	return priority, nil
}

func (priority PriorityType) String() string {
	switch priority {
	case PriorityTypeHigh:
		return PriorityHighString
	case PriorityTypeMedium:
		return PriorityMediumString
	case PriorityTypeLow:
		return PriorityLowString
	}
	return "Unknown"
}

type DpGetterFunc func(uint64) ([]byte, error)

// atomic boolean type which uses a uint32 integer to store boolean value
type AtomicBooleanType struct {
	val uint32
}

const (
	AtomicBooleanTrue  uint32 = 1
	AtomicBooleanFalse uint32 = 0
)

func NewAtomicBooleanType(value bool) *AtomicBooleanType {
	atomicBoolean := &AtomicBooleanType{}
	atomicBoolean.Set(value)
	return atomicBoolean
}

func (a *AtomicBooleanType) Set(value bool) {
	if value {
		a.SetTrue()
	} else {
		a.SetFalse()
	}
}

func (a *AtomicBooleanType) SetTrue() {
	atomic.StoreUint32(&a.val, AtomicBooleanTrue)
}

func (a *AtomicBooleanType) SetFalse() {
	atomic.StoreUint32(&a.val, AtomicBooleanFalse)
}

func (a *AtomicBooleanType) Get() bool {
	return atomic.LoadUint32(&a.val) == AtomicBooleanTrue
}

type FilterExpDelType int

const (
	filterExpDelStripN              = 0
	filterExpDelSkipDelN            = 1
	filterExpDelSkipExpN            = 2
	filterExpDelSkipUncommittedTxnN = 3
	filterExpDelSkipBinaryN         = 4
)

var FilterExpDelNone FilterExpDelType = 0x0
var FilterExpDelStripExpiration FilterExpDelType = 1 << filterExpDelStripN               // 0x1
var FilterExpDelSkipDeletes FilterExpDelType = 1 << filterExpDelSkipDelN                 // 0x2
var FilterExpDelSkipExpiration FilterExpDelType = 1 << filterExpDelSkipExpN              // 0x4
var FilterSkipReplUncommittedTxn FilterExpDelType = 1 << filterExpDelSkipUncommittedTxnN // 0x8
var FilterSkipBinary FilterExpDelType = 1 << filterExpDelSkipBinaryN                     // 0x16
var FilterExpDelAllFiltered = FilterExpDelStripExpiration | FilterExpDelSkipDeletes | FilterExpDelSkipExpiration
var FilterExpDelMax = FilterExpDelStripExpiration | FilterExpDelSkipDeletes | FilterExpDelSkipExpiration | FilterSkipReplUncommittedTxn

func (a *FilterExpDelType) IsStripExpirationSet() bool {
	return *a&FilterExpDelStripExpiration > 0
}

func (a *FilterExpDelType) IsSkipDeletesSet() bool {
	return *a&FilterExpDelSkipDeletes > 0
}

func (a *FilterExpDelType) IsSkipExpirationSet() bool {
	return *a&FilterExpDelSkipExpiration > 0
}

func (a *FilterExpDelType) IsSkipReplicateUncommittedTxnSet() bool {
	return *a&FilterSkipReplUncommittedTxn > 0
}

func (a *FilterExpDelType) IsSkipBinarySet() bool {
	return *a&FilterSkipBinary > 0
}

func (a *FilterExpDelType) SetStripExpiration(setVal bool) {
	curValue := *a&FilterExpDelStripExpiration > 0
	if curValue != setVal {
		*a ^= 1 << filterExpDelStripN
	}
}

func (a *FilterExpDelType) SetSkipDeletes(setVal bool) {
	curValue := *a&FilterExpDelSkipDeletes > 0
	if curValue != setVal {
		*a ^= 1 << filterExpDelSkipDelN
	}
}

func (a *FilterExpDelType) SetSkipExpiration(setVal bool) {
	curValue := *a&FilterExpDelSkipExpiration > 0
	if curValue != setVal {
		*a ^= 1 << filterExpDelSkipExpN
	}
}

func (a *FilterExpDelType) SetSkipReplicateUncommittedTxn(setVal bool) {
	curValue := *a&FilterSkipReplUncommittedTxn > 0
	if curValue != setVal {
		*a ^= 1 << filterExpDelSkipUncommittedTxnN
	}
}

func (a *FilterExpDelType) SetSkipBinary(setVal bool) {
	curValue := *a&FilterSkipBinary > 0
	if curValue != setVal {
		*a ^= 1 << filterExpDelSkipBinaryN
	}
}

func (a FilterExpDelType) String() string {
	return fmt.Sprintf("%d", a)
}

func (a FilterExpDelType) LogString() string {
	return fmt.Sprintf("StripTTL(%v), SkipDeletes(%v), SkipExpiration(%v), SkipUncommittedTxn(%v), SkipBinary(%v)",
		a&FilterExpDelStripExpiration > 0, a&FilterExpDelSkipDeletes > 0, a&FilterExpDelSkipExpiration > 0,
		a&FilterSkipReplUncommittedTxn > 0, a&FilterSkipBinary > 0)
}

type CollectionsMgtType int

const (
	CollectionsMappingKey         = "collectionsExplicitMapping"
	CollectionsMirrorKey          = "collectionsMirroringMode"
	CollectionsMigrateKey         = "collectionsMigrationMode"
	CollectionsOsoKey             = "collectionsOSOMode"
	CollectionsMappingRulesKey    = "colMappingRules"
	CollectionsSkipSourceCheckKey = "collectionsSkipSrcValidation"
)

const (
	colMgtMappingN   = 0 // Non-Set bit if implicit, Set bit if explicit
	colMgtMirroringN = 1 // Non-Set bit if mirroring off, set bit if mirroring on
	colMgtMigrateN   = 2 // Non-set if traditional, set bit if migration mode
	colMgtOsoN       = 3 // Non-set if traditional, set bit to opt into OSO mode
)

var CollectionsExplicitBit CollectionsMgtType = 1 << colMgtMappingN
var CollectionsMirroringBit CollectionsMgtType = 1 << colMgtMirroringN
var CollectionsMigrationBit CollectionsMgtType = 1 << colMgtMigrateN
var CollectionsOSOBit CollectionsMgtType = 1 << colMgtOsoN

var CollectionsMgtDefault = CollectionsOSOBit // Implicit, no mirror, no migration, OSO

var CollectionsMgtMax = CollectionsExplicitBit | CollectionsMirroringBit | CollectionsMigrationBit | CollectionsOSOBit

func (c CollectionsMgtType) String() string {
	var output []string
	output = append(output, "ExplicitMapping:")
	output = append(output, fmt.Sprintf("%v", c.IsExplicitMapping()))
	output = append(output, "Mirroring:")
	output = append(output, fmt.Sprintf("%v", c.IsMirroringOn()))
	output = append(output, "Migration:")
	output = append(output, fmt.Sprintf("%v", c.IsMigrationOn()))
	output = append(output, "OSO:")
	output = append(output, fmt.Sprintf("%v", c.IsOsoOn()))
	return strings.Join(output, " ")
}

func (c *CollectionsMgtType) IsImplicitMapping() bool {
	return !c.IsExplicitMapping() && !c.IsMigrationOn()
}

func (c *CollectionsMgtType) IsExplicitMapping() bool {
	return *c&CollectionsExplicitBit > 0
}

func (c *CollectionsMgtType) SetExplicitMapping(val bool) {
	if c.IsExplicitMapping() != val {
		*c ^= 1 << colMgtMappingN
	}
}

func (c *CollectionsMgtType) IsMirroringOn() bool {
	return *c&CollectionsMirroringBit > 0
}

func (c *CollectionsMgtType) SetMirroring(val bool) {
	if c.IsMirroringOn() != val {
		*c ^= 1 << colMgtMirroringN
	}
}

func (c *CollectionsMgtType) IsMigrationOn() bool {
	return *c&CollectionsMigrationBit > 0
}

func (c *CollectionsMgtType) SetMigration(val bool) {
	if c.IsMigrationOn() != val {
		*c ^= 1 << colMgtMigrateN
	}
}

func (c *CollectionsMgtType) IsOsoOn() bool {
	return *c&CollectionsOSOBit > 0
}

func (c *CollectionsMgtType) SetOSO(val bool) {
	if c.IsOsoOn() != val {
		*c ^= 1 << colMgtOsoN
	}
}

type MergeFunctionMappingType map[string]string

func (mf MergeFunctionMappingType) SameAs(otherRaw interface{}) bool {
	other, ok := otherRaw.(MergeFunctionMappingType)
	if !ok {
		return false
	}
	if len(mf) != len(other) {
		return false
	}
	for k, v := range mf {
		v2, exists := other[k]
		if !exists {
			return false
		}
		if v != v2 {
			return false
		}
	}
	return true
}

func ValidateAndConvertStringToMergeFunctionMappingType(value string) (MergeFunctionMappingType, error) {
	jsonMap := make(map[string]string)
	err := json.Unmarshal([]byte(value), &jsonMap)
	if err != nil {
		return MergeFunctionMappingType{}, err
	}
	// Check for duplicated keys
	res, err := JsonStringReEncodeTest(value)
	if err != nil {
		return MergeFunctionMappingType{}, err
	}
	if res == false {
		return MergeFunctionMappingType{}, fmt.Errorf("JSON string passed in for merge function did not pass re-encode test. Are there potentially duplicated keys?")
	}
	mergeFunc := make(MergeFunctionMappingType, len(jsonMap))
	for k, v := range jsonMap {
		mergeFunc[k] = v
	}
	return mergeFunc, nil
}

// Return true if re-encode has the same length as passed in value with space trimmed
func JsonStringReEncodeTest(value string) (bool, error) {
	trimmedValue := []byte(strings.Replace(value, " ", "", -1))
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(trimmedValue), &jsonMap)
	if err != nil {
		return false, err
	}
	testMarshal, err := json.Marshal(jsonMap)
	if err != nil {
		return false, err
	}
	if len(testMarshal) != len([]byte(trimmedValue)) {
		return false, nil
	} else {
		return true, nil
	}
}

func GetXattrSize(body []byte) (uint32, error) {
	if len(body) < 4 {
		return 0, fmt.Errorf("body is too short to contain valid xattrs")
	}
	xattrSize := binary.BigEndian.Uint32(body[0:4])
	// Couchbase doc size is max of 20MB. Xattribute count against this limit.
	// So if total xattr size is greater than this limit, then something is wrong
	if xattrSize > MaxDocSizeByte {
		return 0, fmt.Errorf("xattrs size %v exceeds max doc size", xattrSize)
	}
	return xattrSize, nil
}

// NOTE: Assumes that Xattribute is present
func StripXattrAndGetBody(bodyWithXattr []byte) ([]byte, error) {
	xattrSize, err := GetXattrSize(bodyWithXattr)
	if err != nil {
		return nil, err
	}
	return bodyWithXattr[xattrSize+4:], nil
}

// An xattribute composer will compose given xattr key value pairs, and then format it properly
// prior to attaching the rest of the document body/value
type XattrComposer struct {
	// Precompiled length slice of data to include
	// - Xattr section
	// - Document Value after the Xattr section
	body []byte

	// cursor
	pos int

	rawMode         composerMode
	rawModeStartPos int

	atLeastOneXattr bool
}

type composerMode int

var errComposerIncorrectMode = fmt.Errorf("composer is in an incorrect mode for this operation")

const (
	compNonRawMode    composerMode = iota
	compKeyMode       composerMode = iota
	compValueMode     composerMode = iota
	compValueDoneMode composerMode = iota
)

func NewXattrComposer(dst []byte) *XattrComposer {
	return &XattrComposer{
		body: dst,
		pos:  4, /* Start at 4 since first uint32 indicates size for entire xattr section */
	}
}

// Raw composer is used in the middle of an xattr such as for testing
func NewXattrRawComposer(dst []byte) *XattrComposer {
	return &XattrComposer{
		body: dst,
		pos:  0,
	}
}

func (x *XattrComposer) WriteKV(key, value []byte) error {
	if x.rawMode != compNonRawMode {
		return errComposerIncorrectMode
	}
	x.atLeastOneXattr = true
	WriteXattrKVPair(x.body, key, value, &x.pos)
	return nil
}

// Note - The destination must be big enough
func WriteXattrKVPair(dst, key, value []byte, pos *int) {
	// Two bytes are for the '\x00' below
	binary.BigEndian.PutUint32(dst[*pos:*pos+4], uint32(len(key)+len(value)+2))
	*pos = *pos + 4
	copy(dst[*pos:*pos+len(key)], key)
	*pos = *pos + len(key)
	dst[*pos] = '\x00'
	*pos++
	copy(dst[*pos:*pos+len(value)], value)
	*pos = *pos + len(value)
	dst[*pos] = '\x00'
	*pos++
}

// Raw mode allows Key and Values to be written in sequence and then calculate the size of the KV pair
// at the end once both K and V have been written
func (x *XattrComposer) StartRawMode() error {
	if x.rawMode != compNonRawMode {
		return errComposerIncorrectMode
	}
	x.rawMode = compKeyMode

	// Reserve the first 4 bytes for total length of the xattribute once value is committed
	x.rawModeStartPos = x.pos
	x.pos += 4
	return nil
}

func (x *XattrComposer) RawWriteKey(key []byte) error {
	if x.rawMode != compKeyMode {
		return errComposerIncorrectMode
	}

	x.atLeastOneXattr = true
	copy(x.body[x.pos:x.pos+len(key)], key)
	x.pos += len(key)
	x.body[x.pos] = '\x00'
	x.pos++
	x.rawMode = compValueMode
	return nil
}

// Allow the caller to gain access to destination slice and position pointer
// Caller must increment the position pointer appropriately
func (x *XattrComposer) RawHijackValue() (body []byte, pos *int, err error) {
	if x.rawMode != compValueMode {
		return nil, nil, errComposerIncorrectMode
	}
	x.rawMode = compValueDoneMode

	return x.body, &x.pos, nil
}

func (x *XattrComposer) CommitRawKVPair() (int, error) {
	if x.rawMode != compValueDoneMode {
		return -1, errComposerIncorrectMode
	}

	x.body[x.pos] = '\x00'
	x.pos++

	// Now, knowing key and value size, mark the beginning
	binary.BigEndian.PutUint32(x.body[x.rawModeStartPos:x.rawModeStartPos+4], uint32(x.pos-(x.rawModeStartPos+4)))
	x.rawMode = compNonRawMode

	return x.pos, nil
}

// Once all the Xattributes are finished, calculate the whole xattr section and append doc value
// Remember to set Xattr flag
func (x *XattrComposer) FinishAndAppendDocValue(val []byte) ([]byte, bool) {
	if !x.atLeastOneXattr {
		// No xattr written - do not do anything
		copy(x.body[0:len(val)], val)
		return x.body[0:len(val)], x.atLeastOneXattr
	}

	// Write the entire Xattr size at the beginning of body
	binary.BigEndian.PutUint32(x.body[0:4], uint32(x.pos-4))
	copy(x.body[x.pos:], val)
	x.pos = x.pos + len(val)
	return x.body[0:x.pos], x.atLeastOneXattr
}

type XattrIterator struct {
	body []byte
	// end position of xattrs
	endPos uint32
	// current cursor position
	pos uint32
}

func NewXattrIterator(body []byte) (*XattrIterator, error) {
	xattrSize, err := GetXattrSize(body)
	if err != nil {
		return nil, err
	}
	return &XattrIterator{
		body:   body,
		endPos: xattrSize + 4,
		// jump to the beginning of the first xattr
		pos: 4,
	}, nil
}

func (xi *XattrIterator) HasNext() bool {
	return xi.pos < xi.endPos
}

func (xi *XattrIterator) Next() ([]byte, []byte, error) {
	// xattrs pattern: uint32 -> key -> NUL -> value -> NUL (repeat)

	xi.pos += 4
	var separator uint32

	// Search for end of key
	for separator = xi.pos; xi.body[separator] != '\x00'; separator++ {
		if separator >= xi.endPos {
			return nil, nil, fmt.Errorf("Error parsing xattr key")
		}
	}

	key := xi.body[xi.pos:separator]

	xi.pos = separator + 1

	// Search for end of value
	for separator = xi.pos; xi.body[separator] != '\x00'; separator++ {
		if separator >= xi.endPos {
			return nil, nil, fmt.Errorf("Error parsing xattr value")
		}
	}

	value := xi.body[xi.pos:separator]

	xi.pos = separator + 1

	return key, value, nil
}

type CCRXattrFieldIterator struct {
	xattr []byte
	pos   int
}

func NewCCRXattrFieldIterator(xattr []byte) (*CCRXattrFieldIterator, error) {
	length := len(xattr)
	if xattr[0] != '{' || xattr[length-1] != '}' {
		return nil, fmt.Errorf("Invalid format for XATTR: %v", xattr)
	}
	return &CCRXattrFieldIterator{
		xattr: xattr,
		pos:   1,
	}, nil
}

func (xfi *CCRXattrFieldIterator) HasNext() bool {
	return xfi.pos < len(xfi.xattr)
}

// Expected format is: {"key":"value","key":"value","key":{...},key:{...}}
func (xfi *CCRXattrFieldIterator) Next() (key, value []byte, err error) {
	beginQuote := xfi.pos
	colon := bytes.Index(xfi.xattr[beginQuote:], []byte{':'})
	endQuote := beginQuote + colon - 1
	if colon == -1 || xfi.xattr[beginQuote] != '"' || xfi.xattr[endQuote] != '"' {
		err = fmt.Errorf("XATTR %s invalid format at pos=%v, char '%c,, endQuote =%v", xfi.xattr, xfi.pos, xfi.xattr[xfi.pos], endQuote)
		return
	}
	key = xfi.xattr[beginQuote+1 : endQuote]

	beginValue := endQuote + 2
	var endValue int
	if xfi.xattr[beginValue] == '"' {
		endValue = beginValue + 1 + bytes.Index(xfi.xattr[beginValue+1:], []byte{'"'})
		if endValue < beginValue {
			err = fmt.Errorf("XATTR %s invalid format searching for key '%s' at pos=%v, beginValue pos %v, char '%c', endValue pos %v, char '%c'", xfi.xattr, key, xfi.pos, beginValue, xfi.xattr[beginValue], endValue, xfi.xattr[endValue])
		}
		// The value is a string
		value = xfi.xattr[beginValue+1 : endValue]
	} else if xfi.xattr[beginValue] == '{' {
		endValue = beginValue + 1 + bytes.Index(xfi.xattr[beginValue+1:], []byte{'}'})
		if endValue < beginValue {
			err = fmt.Errorf("XATTR %s invalid format searching for key '%s' at pos=%v, beginValue pos %v, char '%c', endValue pos %v, char '%c'", xfi.xattr, key, xfi.pos, beginValue, xfi.xattr[beginValue], endValue, xfi.xattr[endValue])
		}
		value = xfi.xattr[beginValue : endValue+1]
	} else {
		err = fmt.Errorf("XATTR %s invalid format searching for key '%s' at pos=%v, beginValue pos %v, char '%c', endValue pos %v, char '%c'", xfi.xattr, key, xfi.pos, beginValue, xfi.xattr[beginValue], endValue, xfi.xattr[endValue])
	}
	xfi.pos = endValue + 2
	return
}

type SortedSeqnoListWithLock struct {
	seqno_list []uint64
	lock       *sync.RWMutex
}

func NewSortedSeqnoListWithLock() *SortedSeqnoListWithLock {
	return &SortedSeqnoListWithLock{make([]uint64, 0), &sync.RWMutex{}}
}

// when needToSort is true, sort the internal seqno_list before returning it
// sorting is needed only when seqno_list is not already sorted, which is the case only for sent_seqno_list
// in other words, needToSort should be set to true only when operating on sent_seqno_list
func (list_obj *SortedSeqnoListWithLock) GetSortedSeqnoList(needToSort bool) []uint64 {
	if needToSort {
		list_obj.lock.Lock()
		defer list_obj.lock.Unlock()

		SortUint64List(list_obj.seqno_list)
		return DeepCopyUint64Array(list_obj.seqno_list)
	} else {
		list_obj.lock.RLock()
		defer list_obj.lock.RUnlock()

		return DeepCopyUint64Array(list_obj.seqno_list)
	}
}

func (list_obj *SortedSeqnoListWithLock) GetLengthOfSeqnoList() int {
	list_obj.lock.RLock()
	defer list_obj.lock.RUnlock()

	return len(list_obj.seqno_list)
}

// append seqno to the end of seqno_list
func (list_obj *SortedSeqnoListWithLock) AppendSeqno(seqno uint64) {
	list_obj.lock.Lock()
	defer list_obj.lock.Unlock()
	list_obj.seqno_list = append(list_obj.seqno_list, seqno)
}

// truncate all seqnos that are no larger than passed in through_seqno
func (list_obj *SortedSeqnoListWithLock) TruncateSeqnos(through_seqno uint64) {
	list_obj.lock.Lock()
	defer list_obj.lock.Unlock()
	seqno_list := list_obj.seqno_list
	index, found := SearchUint64List(seqno_list, through_seqno)
	if found {
		list_obj.seqno_list = seqno_list[index+1:]
	} else if index > 0 {
		list_obj.seqno_list = seqno_list[index:]
	}
}

type Uleb128 []byte

type ExplicitMappingValidator struct {
	// S.C -> S.C || S.C -> nil
	oneToOneRules map[string]interface{}
	// S -> S || S -> nil
	scopeToScopeRules map[string]interface{}
	// Ensure tgt namespaces are not being reused
	valueDuplicatorCheckMap map[string]interface{}
}

type explicitValidatingType int

const (
	explicitRuleInvalid          explicitValidatingType = iota
	explicitRuleInvalidScopeName explicitValidatingType = iota
	explicitRuleOneToOne         explicitValidatingType = iota
	explicitRuleScopeToScope     explicitValidatingType = iota
	explicitRuleInvalidType      explicitValidatingType = iota
	explicitRuleEmptyString      explicitValidatingType = iota
	explicitRuleStringTooLong    explicitValidatingType = iota
	explicitRuleTargetDuplicate  explicitValidatingType = iota
	explicitRuleSystemScope      explicitValidatingType = iota
)

func NewExplicitMappingValidator() *ExplicitMappingValidator {
	return &ExplicitMappingValidator{
		oneToOneRules:           make(map[string]interface{}),
		scopeToScopeRules:       make(map[string]interface{}),
		valueDuplicatorCheckMap: make(map[string]interface{}),
	}
}

func (e *ExplicitMappingValidator) parseRule(k string, v interface{}) explicitValidatingType {
	if k == "" {
		return explicitRuleInvalid
	}
	vStr, vIsString := v.(string)

	kNameSpace, err := NewCollectionNamespaceFromString(k)
	if err == nil {
		// value must be the same or nil
		if v == nil {
			return explicitRuleOneToOne
		}
		if kNameSpace.IsSystemScope() {
			return explicitRuleSystemScope
		}
		if !vIsString {
			return explicitRuleInvalidType
		}
		vNameSpace, err := NewCollectionNamespaceFromString(vStr)
		if err != nil {
			if err == ErrorLengthExceeded {
				return explicitRuleStringTooLong
			}
			return explicitRuleInvalid
		}
		if vNameSpace.IsSystemScope() {
			return explicitRuleSystemScope
		}
		if vIsString {
			_, exists := e.valueDuplicatorCheckMap[vStr]
			if exists {
				return explicitRuleTargetDuplicate
			}
			e.valueDuplicatorCheckMap[vStr] = true
		}

		return explicitRuleOneToOne
	}
	if err != nil && err == ErrorLengthExceeded {
		return explicitRuleStringTooLong
	}

	// At this point, name has to be a scope name
	if len(k) > MaxCollectionNameBytes {
		return explicitRuleStringTooLong
	}

	if k == SystemScopeName || vIsString && vStr == SystemScopeName {
		return explicitRuleSystemScope
	}
	matched := k == DefaultScopeCollectionName
	if !matched {
		matched = CollectionNameValidationRegex.MatchString(k)
	}
	if !matched {
		if strings.Contains(k, ScopeCollectionDelimiter) {
			return explicitRuleInvalid
		} else if k == "" {
			return explicitRuleEmptyString
		} else {
			return explicitRuleInvalidScopeName
		}
	}

	if vIsString {
		if len(vStr) > MaxCollectionNameBytes {
			return explicitRuleStringTooLong
		}
		matched = vStr == DefaultScopeCollectionName
		if !matched {
			matched = CollectionNameValidationRegex.MatchString(vStr)
		}
		if !matched {
			if strings.Contains(k, ScopeCollectionDelimiter) {
				return explicitRuleInvalid
			} else if vStr == "" {
				return explicitRuleEmptyString
			} else {
				return explicitRuleInvalidScopeName
			}
		}
		_, exists := e.valueDuplicatorCheckMap[vStr]
		if exists {
			return explicitRuleTargetDuplicate
		}
		e.valueDuplicatorCheckMap[vStr] = true
	} else if v != nil {
		// Can only be string type or nil type
		return explicitRuleInvalidType
	}
	return explicitRuleScopeToScope
}

// Rules in descending priority:
// 1. S.C -> S'.C'
// 2. S.C -> null
// 3. S -> S'
// 4. S -> null
//
// Given the same scope:
// 1 can coexist with 3 if C' != C
// 1 can coexist with 4
// 2 can coexist with 3
func (e *ExplicitMappingValidator) ValidateKV(k string, v interface{}) error {
	ruleType := e.parseRule(k, v)
	switch ruleType {
	case explicitRuleInvalid:
		return fmt.Errorf("invalid rule: %v:%v", k, v)
	case explicitRuleInvalidScopeName:
		return fmt.Errorf("Invalid scope or collection name for one or both of %v:%v", k, v)
	case explicitRuleInvalidType:
		return fmt.Errorf("Rule must be either string or null types. Received: %v", reflect.TypeOf(v))
	case explicitRuleEmptyString:
		errStr := "Rule must not contain empty string(s):"
		if k == "" {
			errStr += "key is empty "
		}
		if vStr, ok := v.(string); ok && vStr == "" {
			errStr += " value is empty "
		}
		return fmt.Errorf(errStr)
	case explicitRuleOneToOne:
		// Shouldn't have duplicated keys, but check anyway
		_, exists := e.oneToOneRules[k]
		if exists {
			return fmt.Errorf("Collection To Collection rule with source namespace %v already exists", k)
		}
		e.oneToOneRules[k] = v

		// Check for duplicity
		submatches := CollectionNamespaceRegex.FindStringSubmatch(k)
		oneToOneRuleSourceScope := submatches[1]
		oneToOneRuleSourceCollection := submatches[2]
		scopeRuleTargetScopeName, scopeRuleExists := e.scopeToScopeRules[oneToOneRuleSourceScope]
		if v == nil {
			if scopeRuleExists && scopeRuleTargetScopeName == nil {
				// S -> null already exists. S.C -> null is redundant
				return fmt.Errorf("The rule %v%v%v is redundant", k, JsonDelimiter, v)
			}
		} else if scopeRuleExists {
			// A non-deny list rule, i.e. S.C -> ST.CT
			submatches2 := CollectionNamespaceRegex.FindStringSubmatch(v.(string))
			oneToOneRuleTargetScope := submatches2[1]
			oneToOneRuleTargetCollection := submatches2[2]

			if scopeRuleTargetScopeName == oneToOneRuleTargetScope &&
				oneToOneRuleSourceCollection == oneToOneRuleTargetCollection {
				// This means S -> ST exists (since scopeRuleExists) and collection name mapping
				// is following the rule S.C -> ST.C
				// S -> S2 already exists, S.C -> S2.C is redundant
				return fmt.Errorf("The rule %v%v%v is redundant", k, JsonDelimiter, v)
			}
		}
	case explicitRuleScopeToScope:
		// Shouldn't have duplicated keys, but check anyway
		_, exists := e.scopeToScopeRules[k]
		if exists {
			return fmt.Errorf("Scope To Scope rule with source scope %v already exists", k)
		}
		e.scopeToScopeRules[k] = v

		// Check for duplicity
		for checkK, checkV := range e.oneToOneRules {
			submatches := CollectionNamespaceRegex.FindStringSubmatch(checkK)
			sourceScopeName := submatches[1]
			sourceCollectionName := submatches[2]
			if sourceScopeName == k {
				if v == nil && checkV == nil {
					// S1.C1 -> nil already exists
					return fmt.Errorf("The rule %v%v%v is redundant", checkK, JsonDelimiter, checkV)
				} else if v != nil && checkV != nil {
					submatches2 := CollectionNamespaceRegex.FindStringSubmatch(checkV.(string))
					targetColName := submatches2[2]
					if sourceCollectionName == targetColName {
						// S1.C1 -> S2.C1 exists, and trying to enter rule S1 -> S2
						return fmt.Errorf("The rule %v%v%v is redundant", checkK, JsonDelimiter, checkV)
					} else {
						// S1.C1 -> S2.C3 exists, and can coexist with S1 -> S2
					}
				}
			}
		}
	case explicitRuleStringTooLong:
		return fmt.Errorf("%v - maximum length for scope or collection names is %v", ErrorLengthExceeded, MaxCollectionNameBytes)
	case explicitRuleTargetDuplicate:
		return fmt.Errorf("%v is reused as a target namespace", v.(string))
	case explicitRuleSystemScope:
		return ErrorSystemScopeMapped
	default:
		panic(fmt.Sprintf("Unhandled rule type: %v", ruleType))
	}
	return nil
}

type StoppedPipelineCallback func() error
type StoppedPipelineErrCallback func(err error, cbCalled bool)

type PipelineMgrStopCbType func(string, StoppedPipelineCallback, StoppedPipelineErrCallback) error

type CustomCRMeta struct {
	senderId []byte // The source cluster that sent the document to XDCR
	cas      uint64 // The Cas in the document metadata
	cvId     []byte // _xdcr.id in XATTR
	cv       uint64 // _xdcr.cv in XATTR
	pcas     []byte // _xdcr.pc in XATTR
	mv       []byte // _xdcr.mv in XATTR
}

func NewCustomCRMeta(sourceId []byte, cas uint64, id, cvHex, pcas, mv []byte) (*CustomCRMeta, error) {
	if len(cvHex) == 0 {
		return &CustomCRMeta{
			senderId: sourceId,
			cas:      cas,
			cvId:     id,
			cv:       0,
			pcas:     pcas,
			mv:       mv,
		}, nil
	}
	cv, err := HexLittleEndianToUint64(cvHex)
	if err != nil {
		return nil, err
	}
	if cv > cas {
		// cv is initially the same as cas. Cas may increase later. So cv should never be greather than cas
		return nil, fmt.Errorf("cv>cas, cv=%v,cas=%v,cvHex=%s,pcas=%s,mv=%s", cv, cas, cvHex, pcas, mv)
	}
	return &CustomCRMeta{
		senderId: sourceId,
		cas:      cas,
		cvId:     id,
		cv:       cv,
		pcas:     pcas,
		mv:       mv,
	}, nil
}

func (meta *CustomCRMeta) GetCas() uint64 {
	return meta.cas
}

func (meta *CustomCRMeta) GetCv() uint64 {
	return meta.cv
}
func (meta *CustomCRMeta) GetMv() []byte {
	return meta.mv
}

func (meta *CustomCRMeta) GetPcas() []byte {
	return meta.pcas
}

func (meta *CustomCRMeta) GetCvId() []byte {
	return meta.cvId
}

// Returns the id of the cluster that last updated this document
func (meta *CustomCRMeta) DocumentSourceId() []byte {
	if meta.cas > meta.cv {
		// New change in sending cluster so that's the id
		return meta.senderId
	} else {
		return meta.cvId
	}
}

func (meta *CustomCRMeta) IsMergedDoc() bool {
	if meta.mv != nil && meta.cas == meta.cv {
		return true
	} else {
		return false
	}
}

func (meta CustomCRMeta) ContainsVersion(targetId []byte, targetCas uint64) (bool, error) {
	if meta.cas > meta.cv && bytes.Equal(meta.senderId, targetId) {
		if meta.cas >= targetCas {
			return true, nil
		} else {
			return false, nil
		}
	}
	if bytes.Equal(meta.cvId, targetId) {
		if meta.cv >= targetCas {
			return true, nil
		} else {
			return false, nil
		}
	}
	if meta.pcas != nil {
		casInVV, err := findItemInVV(meta.pcas, targetId)
		if err != nil {
			return false, err
		}
		if casInVV >= targetCas {
			return true, nil
		}
	}
	if meta.mv != nil {
		casInVV, err := findItemInVV(meta.mv, targetId)
		if err != nil {
			return false, err
		}
		if casInVV >= targetCas {
			return true, nil
		}
	}
	return false, nil
}

// This routine checks if meta contains the version vector passed in.
func (meta *CustomCRMeta) ContainsVV(vv []byte) (bool, error) {
	it, err := NewCCRXattrFieldIterator(vv)
	if err != nil {
		return false, err
	}
	for it.HasNext() {
		var key, value []byte
		key, value, err = it.Next()
		if err != nil {
			return false, err
		}
		cas, err := Base64ToUint64(value)
		if err != nil {
			return false, err
		}
		if res, err := meta.ContainsVersion(key, cas); err != nil || res == false {
			return false, err
		}
	}
	// Looped through all items in vv and meta contains all of them
	return true, nil
}

// This routine returns a version vector but in the form map[clusterId]CAS.
// If the document is a merge, the map will have a length > 1 where all source of the merge will be included. This is the MV
// If the document is not a merge, the map will have a length of 1 with the cluster ID/CAS of the last update.
func (meta *CustomCRMeta) currentVersions() (map[string]uint64, error) {
	currentVersions := make(map[string]uint64)
	if meta.cas > meta.cv {
		currentVersions[string(meta.senderId)] = meta.cas
	} else if meta.mv == nil {
		currentVersions[string(meta.cvId)] = meta.cv
	} else {
		it, err := NewCCRXattrFieldIterator(meta.mv)
		if err != nil {
			return nil, err
		}
		for it.HasNext() {
			k, v, err := it.Next()
			if err != nil {
				return nil, err
			}
			cas, err := Base64ToUint64(v)
			if err != nil {
				return nil, err
			}
			currentVersions[string(k)] = cas
		}
	}
	return currentVersions, nil
}

// This routine returns the updated PCAS in the form map[clusterId]CAS
func (meta *CustomCRMeta) previousVersions() (map[string]uint64, error) {
	if meta.cv == 0 {
		return nil, nil
	}
	previousVersions := make(map[string]uint64)
	if meta.pcas != nil {
		it, err := NewCCRXattrFieldIterator(meta.pcas)
		if err != nil {
			return nil, err
		}
		for it.HasNext() {
			k, v, err := it.Next()
			if err != nil {
				return nil, err
			}
			cas, err := Base64ToUint64(v)
			if err != nil {
				return nil, err
			}
			previousVersions[string(k)] = cas
		}
	}
	if meta.cas > meta.cv {
		if meta.mv != nil {
			// Update since last merge. Need to add mv to previous versions
			it, err := NewCCRXattrFieldIterator(meta.mv)
			if err != nil {
				return nil, err
			}
			for it.HasNext() {
				k, v, err := it.Next()
				if err != nil {
					return nil, err
				}
				cas, err := Base64ToUint64(v)
				if err != nil {
					return nil, err
				}
				i, ok := previousVersions[string(k)]
				if !ok || cas > i {
					previousVersions[string(k)] = cas
				}
			}
		}
		// Update since cv was set. Need to add cv/cvid to previousVersions.
		// Missing items from previousVersions will cause false conflict and possibly ping-pong
		v, ok := previousVersions[string(meta.cvId)]
		if !ok || meta.cv > v {
			previousVersions[string(meta.cvId)] = meta.cv
		}
	}
	return previousVersions, nil
}

// {"id":...;"cv":...
func (meta *CustomCRMeta) formatIdCv(body []byte, pos int) int {
	body, pos = WriteJsonRawMsg(body, []byte(XATTR_ID), pos, WriteJsonKey, len(XATTR_ID), true /*firstKey*/)
	if meta.cas <= meta.cv {
		// Since there is no new update, id/cv stays the same
		body, pos = WriteJsonRawMsg(body, meta.cvId, pos, WriteJsonValue, len(meta.cvId), false /*firstKey*/)
		cvHex := Uint64ToHexLittleEndian(meta.cv)
		body, pos = WriteJsonRawMsg(body, []byte(XATTR_CV), pos, WriteJsonKey, len(XATTR_CV), false /*firstKey*/)
		body, pos = WriteJsonRawMsg(body, cvHex, pos, WriteJsonValue, len(cvHex), false /*firstKey*/)
	} else {
		// Set id to sender
		body, pos = WriteJsonRawMsg(body, meta.senderId, pos, WriteJsonValue, len(meta.senderId), false /*firstKey*/)
		// Set cv to the new CAS
		cvHex := Uint64ToHexLittleEndian(meta.cas)
		body, pos = WriteJsonRawMsg(body, []byte(XATTR_CV), pos, WriteJsonKey, len(XATTR_CV), false /*firstKey*/)
		body, pos = WriteJsonRawMsg(body, cvHex, pos, WriteJsonValue, len(cvHex), false /*firstKey*/)
	}
	return pos
}

// If key:value exceeds pruning window, it will not be added
// No pruning if pruningWindow == 0
// If an entry is pruned, its CAS will be returned.
func (meta *CustomCRMeta) addEntryToBody(key, value, body []byte, pos int, empty bool, pruningWindow time.Duration) (uint64 /*cas*/, int /*pos*/, bool /*empty*/, error) {
	if pruningWindow > 0 {
		// Convert value to CAS to check for pruning
		cas, err := Base64ToUint64(value)
		if err != nil {
			return 0, pos, empty, err
		}
		if CasDuration(cas, meta.cas) >= pruningWindow {
			return cas, pos, empty, nil
		}
	}
	body, pos = WriteJsonRawMsg(body, key, pos, WriteJsonKey, len(key), empty /*firstKey*/)
	body, pos = WriteJsonRawMsg(body, value, pos, WriteJsonValue, len(value), false /*firstKey*/)
	return 0, pos, false, nil
}

func (meta *CustomCRMeta) formatPv(body []byte, pos int, pruningWindow time.Duration) (int, error) {
	empty := true
	startPos := pos
	body, pos = WriteJsonRawMsg(body, []byte(XATTR_PCAS), pos, WriteJsonKey, len(XATTR_PCAS), false /*firstKey*/)
	if meta.cas > meta.cv {
		// New change, move mv into pv
		if meta.mv != nil {
			it, err := NewCCRXattrFieldIterator(meta.mv)
			if err != nil {
				return startPos, err
			}
			for it.HasNext() {
				key, value, err := it.Next()
				if err != nil {
					return startPos, err
				}
				// Use 0 for pruning window since items from MV will not be pruned
				_, pos, empty, err = meta.addEntryToBody(key, value, body, pos, empty, 0)
			}
		} else if len(meta.cvId) > 0 {
			// New change but no MV, put id/cv into PV
			var err error
			cvBase64 := Uint64ToBase64(meta.cv)
			_, pos, empty, err = meta.addEntryToBody(meta.cvId, cvBase64, body, pos, empty, 0)
			if err != nil {
				return startPos, err
			}
		}
	}
	if meta.pcas != nil {
		var savedKey, savedVal []byte
		var maxCas uint64 = 0
		it, err := NewCCRXattrFieldIterator(meta.pcas)
		if err != nil {
			return startPos, err
		}
		for it.HasNext() {
			key, value, err := it.Next()
			if err != nil {
				return startPos, err
			}
			// The document source will be set as id/cv, so it doesn't need to be in PV
			if bytes.Equal(meta.DocumentSourceId(), key) {
				continue
			}
			var cas uint64
			cas, pos, empty, err = meta.addEntryToBody(key, value, body, pos, empty, pruningWindow)
			if err != nil {
				return startPos, err
			}
			// No value has been written to body and pruning is currently taking place
			if empty && cas > maxCas {
				maxCas = cas
				savedKey = key
				savedVal = value
			}
		}
		if empty && !meta.IsMergedDoc() && meta.pcas != nil && savedKey != nil && savedVal != nil {
			// We have to have one PV entry for future CR
			_, pos, empty, err = meta.addEntryToBody(savedKey, savedVal, body, pos, empty, 0)
		}
	}
	if empty {
		return startPos, nil
	}
	body[pos] = '}'
	pos++
	return pos, nil
}

// This routine construct XATTR _xdcr{...} based on meta. The constructed XATTRs includes updates for
// new change (meta.cas > meta.cv) and pruning in PV
func (meta *CustomCRMeta) ConstructCustomCRXattrForSetMeta(pruningWindow time.Duration, xattrComposer *XattrComposer) (int, error) {
	err := xattrComposer.StartRawMode()
	if err != nil {
		return -1, err
	}

	err = xattrComposer.RawWriteKey([]byte(XATTR_XDCR))
	if err != nil {
		return -1, err
	}

	body, pos, err := xattrComposer.RawHijackValue()
	if err != nil {
		return -1, err
	}

	*pos = meta.formatIdCv(body, *pos)
	if meta.cas <= meta.cv && meta.mv != nil {
		// Since there is no new change, mv stays the same
		body, *pos = WriteJsonRawMsg(body, []byte(XATTR_MV), *pos, WriteJsonKey, len(XATTR_MV), false /*firstKey*/)
		// copy the mv since it is already in json format
		copy(body[*pos:], meta.mv)
		*pos = *pos + len(meta.mv)
	}
	*pos, err = meta.formatPv(body, *pos, pruningWindow)
	if err != nil {
		return -1, err
	}
	body[*pos] = '}'
	*pos++

	return xattrComposer.CommitRawKVPair()
}

// Called when target has a smaller CAS but its xattrs dominates
// This routine receives the target CustomCRMeta and returns an updated pcas and mv to be used to send subdoc_multi_mutation to source.
// Target with smaller CAS can only dominate source if source document is a merged document.
func (meta *CustomCRMeta) UpdateMetaForSetBack() (pcas, mv []byte, err error) {
	if meta.IsMergedDoc() {
		// This is a merged doc. Setback will just take the same pcas and mv. The calling routine will generate new CV/CvID/CAS at source
		// as if source is doing the same merge
		return meta.GetPcas(), meta.GetMv(), nil
	} else {
		// There has been changes since previous merges. The document history needs to be all in PCAS.
		// This history includes everything in PCAS, MV, cv/cvid, and cas/senderId
		previousVersions, err := meta.previousVersions()
		if err != nil {
			return nil, nil, err
		}
		// Need to add the current sender/CAS and cv/cvId to previous version since setback will generate new CAS/cv
		v, ok := previousVersions[string(meta.cvId)]
		if !ok || v < meta.cv {
			previousVersions[string(meta.cvId)] = meta.cv
		}
		v, ok = previousVersions[string(meta.senderId)]
		if ok || v < meta.cas {
			previousVersions[string(meta.senderId)] = meta.cas
		}
		pcaslen := 0
		for key, _ := range previousVersions {
			pcaslen = pcaslen + len(key) + MaxBase64CASLength + 6 // quotes and sepeartors
		}
		pcaslen = pcaslen + 2 // { and }
		// TODO(MB-41808): data pool
		pcas = make([]byte, pcaslen)
		firstKey := true
		pos := 0
		for key, cas := range previousVersions {
			value := Uint64ToBase64(cas)
			pcas, pos = WriteJsonRawMsg(pcas, []byte(key), pos, WriteJsonKey, len(key), firstKey /*firstKey*/)
			pcas, pos = WriteJsonRawMsg(pcas, value, pos, WriteJsonValue, len(value), false /*firstKey*/)
			firstKey = false
		}
		pcas[pos] = '}'
		pos++
		return pcas[:pos], nil, nil
	}
}

// This routine combines source and target meta and puts the new MV/PCAS in the mergedMvSlice/mergedPcasSlice
//
//	It finds the current versions (MV) of the two documents, combine them into new MV.
//	It finds the previous versions (PCAS) of the two documents, combine them into new PCAS
//	It then remove any duplicates and format them in the provided slices
func (meta *CustomCRMeta) MergeMeta(targetMeta *CustomCRMeta, mergedMvSlice, mergedPcasSlice []byte, pruningWindow time.Duration) (mvlen int, pcaslen int, err error) {
	currentVersions, err := meta.currentVersions()
	if err != nil {
		return 0, 0, err
	}
	targetCurrentVersions, err := targetMeta.currentVersions()
	if err != nil {
		return 0, 0, err
	}
	// merge current versions
	for key, value := range targetCurrentVersions {
		v, ok := currentVersions[key]
		if !ok || value > v {
			currentVersions[key] = value
		}
	}
	previousVersions, err := meta.previousVersions()
	if err != nil {
		return 0, 0, err
	}
	targetPreviousVersions, err := targetMeta.previousVersions()
	if err != nil {
		return 0, 0, err
	}
	if previousVersions == nil {
		previousVersions = targetPreviousVersions
	} else {
		// merge previous versions
		for key, value := range targetPreviousVersions {
			v, ok := previousVersions[key]
			if !ok || value > v {
				previousVersions[key] = value
			}
		}
	}
	// Remove any entries in previousVersions that are already in currentVersions or exceed pruningWindow
	for key, value := range previousVersions {
		if _, ok := currentVersions[key]; ok {
			delete(previousVersions, key)
		} else {
			if pruningWindow > 0 && CasDuration(value, meta.cas) >= pruningWindow {
				delete(previousVersions, key)
			}
		}
	}
	// Construct MV
	firstKey := true
	for key, cas := range currentVersions {
		value := Uint64ToBase64(cas)
		mergedMvSlice, mvlen = WriteJsonRawMsg(mergedMvSlice, []byte(key), mvlen, WriteJsonKey, len(key), firstKey /*firstKey*/)
		mergedMvSlice, mvlen = WriteJsonRawMsg(mergedMvSlice, value, mvlen, WriteJsonValue, len(value), false /*firstKey*/)
		firstKey = false
	}
	if mvlen > 0 {
		// We have MV. Finish it with '}'
		mergedMvSlice[mvlen] = '}'
		mvlen++
	}
	// Construct PCAS
	firstKey = true
	for key, cas := range previousVersions {
		value := Uint64ToBase64(cas)
		mergedPcasSlice, pcaslen = WriteJsonRawMsg(mergedPcasSlice, []byte(key), pcaslen, WriteJsonKey, len(key), firstKey /*firstKey*/)
		mergedPcasSlice, pcaslen = WriteJsonRawMsg(mergedPcasSlice, value, pcaslen, WriteJsonValue, len(value), false /*firstKey*/)
		firstKey = false
	}
	if pcaslen > 0 {
		// We have PCAS. Finish it with '}
		mergedPcasSlice[pcaslen] = '}'
		pcaslen++
	}
	return mvlen, pcaslen, nil
}

func (meta *CustomCRMeta) NeedUpdate(pruningWindow time.Duration) (bool, error) {
	// We need to update if CAS has changed from cv.CAS
	if meta.GetCas() > meta.GetCv() {
		return true, nil
	}

	// The CAS has not changed. But do we need to prune PV?

	// No pruning if its setting is 0 or there is no pv to prune
	if pruningWindow == 0 || meta.pcas == nil {
		return false, nil
	}
	it, err := NewCCRXattrFieldIterator(meta.pcas)
	if err != nil {
		return false, err
	}
	for it.HasNext() {
		_, v, err := it.Next()
		if err != nil {
			return false, err
		}
		cas, err := Base64ToUint64(v)
		if err != nil {
			return false, err
		}
		// Check for pruning window
		if CasDuration(cas, meta.cas) >= pruningWindow {
			return true, nil
		}
	}
	return false, nil
}

func (meta *CustomCRMeta) String() string {
	return fmt.Sprintf("senderId: %s, cas: %v, cvId: %s, cv: %v, pcas: %s, mv: %s",
		meta.senderId, meta.cas, meta.cvId, meta.cv, meta.pcas, meta.mv)
}

func MergedPcasLength(sourceMeta, targetMeta *CustomCRMeta) int {
	return len(sourceMeta.GetPcas()) + len(targetMeta.GetPcas()) + 2 + // maximum combined PCAS + 2 seperators ','
		len(sourceMeta.GetMv()) + len(targetMeta.GetMv()) + 2 + // maximum combined MV rolled into to PCAS + 2 seperators ','
		len(sourceMeta.GetCvId()) + 3 + // quotes and separators "...":
		len(targetMeta.GetCvId()) + 3 + // quotes and separators "...":
		2*(MaxBase64CASLength+3) + 2 // quotes and separators "...", plus '{' and '}'
}

var MaxBase64CASLength = len(Uint64ToBase64(math.MaxUint64))

func MergedMvLength(sourceMeta, targetMeta *CustomCRMeta) int {
	return len(sourceMeta.GetMv()) + len(targetMeta.GetMv()) + 2 + // maximumn combined MV + 2 seperators ','
		len(sourceMeta.DocumentSourceId()) + 3 + // quotes and separators "...":
		len(targetMeta.DocumentSourceId()) + 3 + // quotes and separators "...":
		2*(MaxBase64CASLength+3) + 2 // quotes and separators "...", plus '{' and '}'
}

type ConflictParams struct {
	Source         *WrappedMCRequest
	Target         *SubdocLookupResponse
	SourceId       []byte
	TargetId       []byte
	MergeFunction  string
	ResultNotifier MergeResultNotifier
	TimeoutMs      uint32
	ObjectRecycler func(request *WrappedMCRequest) // for source WrappedMCRequest cleanup
}

type ConflictManagerAction int

const (
	SetMergeToSource  ConflictManagerAction = iota
	SetTargetToSource ConflictManagerAction = iota
	SetToConflictFeed ConflictManagerAction = iota
)

type MergeInputAndResult struct {
	Action ConflictManagerAction
	Input  *ConflictParams
	Result interface{}
	Err    error
}

type MergeResultNotifier interface {
	NotifyMergeResult(input *ConflictParams, mergeResult interface{}, mergeError error)
}

// So on pool details and bucket details "status":
// "healthy" currently means node recently sent heartbeat. It indicates that distributed erlang facility works between current
//
//	node and node which status is shown. And that erlang functions reasonably well there.
//	If there are buckets defined it also implies that all buckets are ready state. I.e. available for data ops.
//
// "status": "unhealthy" means that we haven't seen heartbeats from this node. Either node is down or it's sufficiently
//
//	unhealthy to be unable to send heartbeats. Or network is down.
//
// "status": "warmup" means that there are bucket(s) on this node that are either being warmed up or otherwise unavailable.
//
//	Internally we're doing gen_server call with timeout of 5 seconds to find out if per-bucket memcached can respond and if it thinks that bucket is warmed up. If call either times out or ns_memcached is dead (not created yet or recently crashed) or if ns_memcached is fine but bucket is warming up, then we'll mark entire node as "warmup".
type HeartbeatStatus int

const (
	HeartbeatHealthy   HeartbeatStatus = iota
	HeartbeatUnhealthy HeartbeatStatus = iota
	HeartbeatWarmup    HeartbeatStatus = iota
	HeartbeatInvalid   HeartbeatStatus = iota
)

func NewHeartbeatStatusFromString(statusString string) (HeartbeatStatus, error) {
	switch statusString {
	case "healthy":
		return HeartbeatHealthy, nil
	case "warmup":
		return HeartbeatWarmup, nil
	case "unhealthy":
		return HeartbeatUnhealthy, nil
	default:
		return HeartbeatInvalid, fmt.Errorf("invalid heartbeat string: %v", statusString)
	}
}

func (h HeartbeatStatus) String() string {
	switch h {
	case HeartbeatHealthy:
		return "healthy"
	case HeartbeatWarmup:
		return "warmup"
	case HeartbeatUnhealthy:
		return "unhealthy"
	default:
		return "?? (HeartbeatStatus)"
	}
}

type KvVBMapType map[string][]uint16

func (m *KvVBMapType) GetKeyList() []string {
	if m == nil {
		return nil
	}
	var retList []string
	for k, _ := range *m {
		retList = append(retList, k)
	}
	return retList
}

func (k *KvVBMapType) CompileLookupIndex() map[uint16]string {
	retMap := make(map[uint16]string)

	for kv, vbs := range *k {
		for _, vb := range vbs {
			retMap[vb] = kv
		}
	}

	return retMap
}

func (k *KvVBMapType) GetSortedVBList() []uint16 {
	if k == nil {
		return nil
	}

	indexMap := k.CompileLookupIndex()

	var vbList []uint16
	for vb, _ := range indexMap {
		vbList = append(vbList, vb)
	}
	return SortUint16List(vbList)
}

func (k *KvVBMapType) Clone() KvVBMapType {
	clonedMap := make(KvVBMapType)
	if k == nil {
		return clonedMap
	}

	for key, val := range *k {
		clonedList := make([]uint16, len(val))
		for i, vbno := range val {
			clonedList[i] = vbno
		}
		clonedMap[key] = clonedList
	}
	return clonedMap
}

func (k *KvVBMapType) GreenClone(poolGet func(keys []string) *KvVBMapType) *KvVBMapType {
	if k == nil || *k == nil {
		return nil
	}

	recycledMap := poolGet(k.GetKeyList())
	for key, val := range *k {
		(*recycledMap)[key] = CloneUint16List(val)
	}

	return recycledMap
}

func GetKeysListFromStrStrMap(a map[string]string) []string {
	var retList []string
	for k, _ := range a {
		retList = append(retList, k)
	}
	return retList
}

func GetKeysListFromStrBoolMap(a map[string]bool) []string {
	var retList []string
	for k, _ := range a {
		retList = append(retList, k)
	}
	return retList
}

type DcpStatsMapType map[string]*StringStringMap

func (t *DcpStatsMapType) GetKeyList() []string {
	if t == nil || *t == nil {
		return nil
	}
	var retList []string
	for k, _ := range *t {
		retList = append(retList, k)
	}
	return retList
}

func (t *DcpStatsMapType) Clone() DcpStatsMapType {
	if t == nil || *t == nil {
		return nil
	}

	clonedMap := make(DcpStatsMapType)

	for k, vMap := range *t {
		if vMap == nil {
			clonedMap[k] = nil
		} else {
			vMapClone := make(StringStringMap)
			for k2, v2 := range *vMap {
				vMapClone[k2] = v2
			}
			clonedMap[k] = &vMapClone
		}
	}
	return clonedMap
}

func (t *DcpStatsMapType) GreenClone(dcpStatsPool func(keys []string) *DcpStatsMapType, strStrPool func(keys []string) *StringStringMap) *DcpStatsMapType {
	if t == nil || *t == nil {
		return nil
	}

	recycledMap := dcpStatsPool(t.GetKeyList())
	for key, vMap := range *t {
		if vMap == nil {
			(*recycledMap)[key] = nil
		} else {
			recycledStrStr := strStrPool(GetKeysListFromStrStrMap(*vMap))
			for key2, val2 := range *vMap {
				(*recycledStrStr)[key2] = val2
			}
			(*recycledMap)[key] = recycledStrStr
		}
	}
	return recycledMap
}

type BucketInfoMapType map[string]interface{}

// Shallow copy clone of the values
func (t *BucketInfoMapType) Clone() BucketInfoMapType {
	if t == nil || *t == nil {
		return nil
	}

	clonedMap := make(BucketInfoMapType)

	for k, v := range *t {
		clonedMap[k] = v
	}

	return clonedMap
}

func (t *BucketInfoMapType) GetKeyList() []string {
	if t == nil || *t == nil {
		return nil
	}

	var retList []string
	for key, _ := range *t {
		retList = append(retList, key)
	}
	return retList
}

type HighSeqnosMapType map[string]*map[uint16]uint64

func GetVBListFromSeqnosMap(a map[uint16]uint64) []uint16 {
	var retList []uint16
	for vbno, _ := range a {
		retList = append(retList, vbno)
	}
	return retList
}

func (h HighSeqnosMapType) Clone() HighSeqnosMapType {
	clonedMap := make(HighSeqnosMapType)

	for k, vMap := range h {
		if vMap == nil {
			clonedMap[k] = nil
		} else {
			clonedVMap := make(map[uint16]uint64)
			for k2, v2 := range *vMap {
				clonedVMap[k2] = v2
			}
			clonedMap[k] = &clonedVMap
		}
	}
	return clonedMap
}

func (h *HighSeqnosMapType) GetKeyList() []string {
	if h == nil || *h == nil {
		return nil
	}
	var retList []string
	for k, _ := range *h {
		retList = append(retList, k)
	}
	return retList
}

func (h *HighSeqnosMapType) GreenClone(highSeqnoPool func(keys []string) *HighSeqnosMapType, vbSeqnoPool func(vbnos []uint16) *map[uint16]uint64) *HighSeqnosMapType {
	if h == nil || *h == nil {
		return nil
	}

	recycledVbSeqnoMap := highSeqnoPool(h.GetKeyList())
	for key, vbSeqnoMap := range *h {
		if vbSeqnoMap == nil {
			(*recycledVbSeqnoMap)[key] = nil
		} else {
			recycledSeqnoMap := vbSeqnoPool(GetVBListFromSeqnosMap(*vbSeqnoMap))
			for vbno, seqno := range *vbSeqnoMap {
				(*recycledSeqnoMap)[vbno] = seqno
			}
			(*recycledVbSeqnoMap)[key] = recycledSeqnoMap
		}
	}
	return recycledVbSeqnoMap
}

type VbHostsMapType map[uint16]*[]string

func (v *VbHostsMapType) Clone() VbHostsMapType {
	if v == nil {
		return nil
	}

	cloneMap := make(VbHostsMapType)

	for key, val := range *v {
		cloneList := make([]string, len(*val))
		for i, hostName := range *val {
			cloneList[i] = hostName
		}
		cloneMap[key] = &cloneList
	}
	return cloneMap
}

func (v *VbHostsMapType) GreenClone(pool func(vbnos []uint16) *VbHostsMapType, stringSlicePool func() *[]string) *VbHostsMapType {
	if v == nil || *v == nil {
		return nil
	}

	recycledMap := pool(v.GetKeyList())
	for key, val := range *v {
		for _, strVal := range *val {
			if (*recycledMap)[key] == nil {
				(*recycledMap)[key] = stringSlicePool()
			}
			*(*recycledMap)[key] = append(*(*recycledMap)[key], strVal)
		}
	}

	return recycledMap
}

func (v *VbHostsMapType) GetKeyList() []uint16 {
	if v == nil {
		return nil
	}
	var retList []uint16
	for vbno, _ := range *v {
		retList = append(retList, vbno)
	}
	return retList
}

func (v *VbHostsMapType) Clear() *VbHostsMapType {
	if v == nil {
		return nil
	}
	for k, _ := range *v {
		delete(*v, k)
	}
	return v
}

type StringStringMap map[string]string

func (s StringStringMap) Clone() StringStringMap {
	cloneMap := make(StringStringMap)
	for k, v := range s {
		cloneMap[k] = v
	}
	return cloneMap
}

type UIWarningImpl struct {
	GenericWarnings []string
	FieldWarnings   map[string]interface{}
}

func NewUIWarning() *UIWarningImpl {
	return &UIWarningImpl{
		FieldWarnings: make(map[string]interface{}),
	}
}

func (u *UIWarningImpl) isEmpty() bool {
	return len(u.GenericWarnings) == 0 && len(u.FieldWarnings) == 0
}

func (u *UIWarningImpl) AppendGeneric(warning string) {
	u.GenericWarnings = append(u.GenericWarnings, warning)
}

func (u *UIWarningImpl) AddWarning(key string, val string) {
	u.FieldWarnings[key] = val
}

func (u *UIWarningImpl) String() string {
	if u == nil || u.isEmpty() {
		return ""
	}
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("Internal error: %v", err)
	}
	return string(out)
}

// Successful Creation warnings from UI are essentially pop up boxes with strings in each box
func (u *UIWarningImpl) GetSuccessfulWarningStrings() []string {
	if u == nil {
		return nil
	}

	var retVal []string
	for _, warning := range u.GenericWarnings {
		retVal = append(retVal, warning)
	}

	for fieldKey, warning := range u.FieldWarnings {
		retVal = append(retVal, fmt.Sprintf("%v:%v", fieldKey, warning))
	}

	return retVal
}

func (u *UIWarningImpl) Len() int {
	if u == nil {
		return 0
	}

	var count int
	count += len(u.FieldWarnings)
	count += len(u.GenericWarnings)
	return count
}

func (u *UIWarningImpl) GetFieldWarningsOnly() map[string]interface{} {
	if u == nil {
		return nil
	}

	return u.FieldWarnings
}

func (s *StringStringMap) GetKeyList() []string {
	if s == nil || *s == nil {
		return nil
	}
	var retList []string
	for k, _ := range *s {
		retList = append(retList, k)
	}
	return retList
}

func (s *StringStringMap) GreenClone(strStrPool func(keys []string) *StringStringMap) *StringStringMap {
	if s == nil || *s == nil {
		return nil
	}

	recycledMap := strStrPool(s.GetKeyList())
	for key, val := range *s {
		(*recycledMap)[key] = val
	}
	return recycledMap
}

func (s *StringStringMap) Clear() *StringStringMap {
	if s == nil {
		return nil
	}

	for k, _ := range *s {
		delete(*s, k)
	}
	return s
}

func NewUuid() string {
	return uuid.New().String()
}

type PipelineStatusGauge int

const (
	PipelineStatusPaused  PipelineStatusGauge = iota
	PipelineStatusRunning PipelineStatusGauge = iota
	PipelineStatusError   PipelineStatusGauge = iota
	PipelineStatusMax     PipelineStatusGauge = iota // Boundary
)

func (p PipelineStatusGauge) String() string {
	switch p {
	case PipelineStatusPaused:
		return "paused"
	case PipelineStatusRunning:
		return "running"
	case PipelineStatusError:
		return "error"
	default:
		return ""
	}
}

// hostname -> list of connection errors mapping
type HostToErrorsMapType map[string][]string

// hostAddr -> <portKey -> port> mapping
type HostPortMapType map[string]map[string]uint16

// sourceNode -> <targetNode -> list of connection errors> mapping
type ConnectionErrMapType map[string]map[string][]string

type PortType int

type FilteringStatusType int

// Stats per vbucket
type VBCountMetricMap map[string]int64
