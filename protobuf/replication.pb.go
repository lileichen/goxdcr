// Code generated by protoc-gen-go.
// source: replication.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	replication.proto

It has these top-level messages:
	CreateReplicationRequest
	CreateReplicationResponse
	DeleteReplicationRequest
	ViewSettingsRequest
	Settings
	ChangeGlobalSettingsRequest
	ChangeReplicationSettingsRequest
	ChangeInternalSettingsRequest
	GetStatisticsRequest
*/
package protobuf

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type CreateReplicationRequest_Mode int32

const (
	CreateReplicationRequest_capi CreateReplicationRequest_Mode = 0
	CreateReplicationRequest_xmem CreateReplicationRequest_Mode = 1
)

var CreateReplicationRequest_Mode_name = map[int32]string{
	0: "capi",
	1: "xmem",
}
var CreateReplicationRequest_Mode_value = map[string]int32{
	"capi": 0,
	"xmem": 1,
}

func (x CreateReplicationRequest_Mode) Enum() *CreateReplicationRequest_Mode {
	p := new(CreateReplicationRequest_Mode)
	*p = x
	return p
}
func (x CreateReplicationRequest_Mode) String() string {
	return proto.EnumName(CreateReplicationRequest_Mode_name, int32(x))
}
func (x *CreateReplicationRequest_Mode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(CreateReplicationRequest_Mode_value, data, "CreateReplicationRequest_Mode")
	if err != nil {
		return err
	}
	*x = CreateReplicationRequest_Mode(value)
	return nil
}

type GetStatisticsRequest_Stats int32

const (
	GetStatisticsRequest_docs_written           GetStatisticsRequest_Stats = 0
	GetStatisticsRequest_data_replicated        GetStatisticsRequest_Stats = 1
	GetStatisticsRequest_changes_left           GetStatisticsRequest_Stats = 2
	GetStatisticsRequest_docs_checked           GetStatisticsRequest_Stats = 3
	GetStatisticsRequest_num_checkpoints        GetStatisticsRequest_Stats = 4
	GetStatisticsRequest_num_failedckpts        GetStatisticsRequest_Stats = 5
	GetStatisticsRequest_size_rep_queue         GetStatisticsRequest_Stats = 6
	GetStatisticsRequest_time_committing        GetStatisticsRequest_Stats = 7
	GetStatisticsRequest_bandwidth_usage        GetStatisticsRequest_Stats = 8
	GetStatisticsRequest_docs_lanecy_aggr       GetStatisticsRequest_Stats = 9
	GetStatisticsRequest_docs_latency_wt        GetStatisticsRequest_Stats = 10
	GetStatisticsRequest_docs_req_queue         GetStatisticsRequest_Stats = 11
	GetStatisticsRequest_meta_latency_aggr      GetStatisticsRequest_Stats = 12
	GetStatisticsRequest_meta_latency_wt        GetStatisticsRequest_Stats = 13
	GetStatisticsRequest_rate_replication       GetStatisticsRequest_Stats = 14
	GetStatisticsRequest_docs_opt_repd          GetStatisticsRequest_Stats = 15
	GetStatisticsRequest_active_vbreps          GetStatisticsRequest_Stats = 16
	GetStatisticsRequest_waiting_vbreps         GetStatisticsRequest_Stats = 17
	GetStatisticsRequest_time_working           GetStatisticsRequest_Stats = 18
	GetStatisticsRequest_timeout_percentage_map GetStatisticsRequest_Stats = 19
)

var GetStatisticsRequest_Stats_name = map[int32]string{
	0:  "docs_written",
	1:  "data_replicated",
	2:  "changes_left",
	3:  "docs_checked",
	4:  "num_checkpoints",
	5:  "num_failedckpts",
	6:  "size_rep_queue",
	7:  "time_committing",
	8:  "bandwidth_usage",
	9:  "docs_lanecy_aggr",
	10: "docs_latency_wt",
	11: "docs_req_queue",
	12: "meta_latency_aggr",
	13: "meta_latency_wt",
	14: "rate_replication",
	15: "docs_opt_repd",
	16: "active_vbreps",
	17: "waiting_vbreps",
	18: "time_working",
	19: "timeout_percentage_map",
}
var GetStatisticsRequest_Stats_value = map[string]int32{
	"docs_written":           0,
	"data_replicated":        1,
	"changes_left":           2,
	"docs_checked":           3,
	"num_checkpoints":        4,
	"num_failedckpts":        5,
	"size_rep_queue":         6,
	"time_committing":        7,
	"bandwidth_usage":        8,
	"docs_lanecy_aggr":       9,
	"docs_latency_wt":        10,
	"docs_req_queue":         11,
	"meta_latency_aggr":      12,
	"meta_latency_wt":        13,
	"rate_replication":       14,
	"docs_opt_repd":          15,
	"active_vbreps":          16,
	"waiting_vbreps":         17,
	"time_working":           18,
	"timeout_percentage_map": 19,
}

func (x GetStatisticsRequest_Stats) Enum() *GetStatisticsRequest_Stats {
	p := new(GetStatisticsRequest_Stats)
	*p = x
	return p
}
func (x GetStatisticsRequest_Stats) String() string {
	return proto.EnumName(GetStatisticsRequest_Stats_name, int32(x))
}
func (x *GetStatisticsRequest_Stats) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(GetStatisticsRequest_Stats_value, data, "GetStatisticsRequest_Stats")
	if err != nil {
		return err
	}
	*x = GetStatisticsRequest_Stats(value)
	return nil
}

type CreateReplicationRequest struct {
	FromBucket       *string                        `protobuf:"bytes,1,req,name=fromBucket" json:"fromBucket,omitempty"`
	ToCluster        *string                        `protobuf:"bytes,2,req,name=toCluster" json:"toCluster,omitempty"`
	ToBucket         *string                        `protobuf:"bytes,3,req,name=toBucket" json:"toBucket,omitempty"`
	FilterName       *string                        `protobuf:"bytes,4,req,name=filterName" json:"filterName,omitempty"`
	Mode             *CreateReplicationRequest_Mode `protobuf:"varint,5,req,name=mode,enum=protobuf.CreateReplicationRequest_Mode,def=0" json:"mode,omitempty"`
	Settings         *Settings                      `protobuf:"bytes,6,opt,name=settings" json:"settings,omitempty"`
	Forward          *bool                          `protobuf:"varint,7,opt,name=forward,def=1" json:"forward,omitempty"`
	XXX_unrecognized []byte                         `json:"-"`
}

func (m *CreateReplicationRequest) Reset()         { *m = CreateReplicationRequest{} }
func (m *CreateReplicationRequest) String() string { return proto.CompactTextString(m) }
func (*CreateReplicationRequest) ProtoMessage()    {}

const Default_CreateReplicationRequest_Mode CreateReplicationRequest_Mode = CreateReplicationRequest_capi
const Default_CreateReplicationRequest_Forward bool = true

func (m *CreateReplicationRequest) GetFromBucket() string {
	if m != nil && m.FromBucket != nil {
		return *m.FromBucket
	}
	return ""
}

func (m *CreateReplicationRequest) GetToCluster() string {
	if m != nil && m.ToCluster != nil {
		return *m.ToCluster
	}
	return ""
}

func (m *CreateReplicationRequest) GetToBucket() string {
	if m != nil && m.ToBucket != nil {
		return *m.ToBucket
	}
	return ""
}

func (m *CreateReplicationRequest) GetFilterName() string {
	if m != nil && m.FilterName != nil {
		return *m.FilterName
	}
	return ""
}

func (m *CreateReplicationRequest) GetMode() CreateReplicationRequest_Mode {
	if m != nil && m.Mode != nil {
		return *m.Mode
	}
	return Default_CreateReplicationRequest_Mode
}

func (m *CreateReplicationRequest) GetSettings() *Settings {
	if m != nil {
		return m.Settings
	}
	return nil
}

func (m *CreateReplicationRequest) GetForward() bool {
	if m != nil && m.Forward != nil {
		return *m.Forward
	}
	return Default_CreateReplicationRequest_Forward
}

type CreateReplicationResponse struct {
	Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
	Database         *string `protobuf:"bytes,2,req,name=database" json:"database,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CreateReplicationResponse) Reset()         { *m = CreateReplicationResponse{} }
func (m *CreateReplicationResponse) String() string { return proto.CompactTextString(m) }
func (*CreateReplicationResponse) ProtoMessage()    {}

func (m *CreateReplicationResponse) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *CreateReplicationResponse) GetDatabase() string {
	if m != nil && m.Database != nil {
		return *m.Database
	}
	return ""
}

type DeleteReplicationRequest struct {
	Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
	Forward          *bool   `protobuf:"varint,2,opt,name=forward,def=1" json:"forward,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DeleteReplicationRequest) Reset()         { *m = DeleteReplicationRequest{} }
func (m *DeleteReplicationRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteReplicationRequest) ProtoMessage()    {}

const Default_DeleteReplicationRequest_Forward bool = true

func (m *DeleteReplicationRequest) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *DeleteReplicationRequest) GetForward() bool {
	if m != nil && m.Forward != nil {
		return *m.Forward
	}
	return Default_DeleteReplicationRequest_Forward
}

type ViewSettingsRequest struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *ViewSettingsRequest) Reset()         { *m = ViewSettingsRequest{} }
func (m *ViewSettingsRequest) String() string { return proto.CompactTextString(m) }
func (*ViewSettingsRequest) ProtoMessage()    {}

type Settings struct {
	Protocol                       *string `protobuf:"bytes,1,opt,name=protocol" json:"protocol,omitempty"`
	FilterExpression               *string `protobuf:"bytes,2,opt,name=filterExpression" json:"filterExpression,omitempty"`
	CheckpointInterval             *uint32 `protobuf:"varint,3,opt,name=checkpointInterval" json:"checkpointInterval,omitempty"`
	WorkerBatchSize                *uint32 `protobuf:"varint,4,opt,name=workerBatchSize" json:"workerBatchSize,omitempty"`
	DocBatchSizeKb                 *uint32 `protobuf:"varint,5,opt,name=docBatchSizeKb" json:"docBatchSizeKb,omitempty"`
	FailureRestartInterval         *uint32 `protobuf:"varint,6,opt,name=failureRestartInterval" json:"failureRestartInterval,omitempty"`
	OptimisticReplicationThreshold *uint32 `protobuf:"varint,7,opt,name=optimisticReplicationThreshold" json:"optimisticReplicationThreshold,omitempty"`
	HttpConnections                *uint32 `protobuf:"varint,8,opt,name=httpConnections" json:"httpConnections,omitempty"`
	SourceNozzlePerNode            *uint32 `protobuf:"varint,9,opt,name=sourceNozzlePerNode" json:"sourceNozzlePerNode,omitempty"`
	TargetNozzlePerNode            *uint32 `protobuf:"varint,10,opt,name=targetNozzlePerNode" json:"targetNozzlePerNode,omitempty"`
	MaxExpectedReplicationLag      *uint32 `protobuf:"varint,11,opt,name=maxExpectedReplicationLag" json:"maxExpectedReplicationLag,omitempty"`
	TimeoutPercentageCap           *uint32 `protobuf:"varint,12,opt,name=timeoutPercentageCap" json:"timeoutPercentageCap,omitempty"`
	XXX_unrecognized               []byte  `json:"-"`
}

func (m *Settings) Reset()         { *m = Settings{} }
func (m *Settings) String() string { return proto.CompactTextString(m) }
func (*Settings) ProtoMessage()    {}

func (m *Settings) GetProtocol() string {
	if m != nil && m.Protocol != nil {
		return *m.Protocol
	}
	return ""
}

func (m *Settings) GetFilterExpression() string {
	if m != nil && m.FilterExpression != nil {
		return *m.FilterExpression
	}
	return ""
}

func (m *Settings) GetCheckpointInterval() uint32 {
	if m != nil && m.CheckpointInterval != nil {
		return *m.CheckpointInterval
	}
	return 0
}

func (m *Settings) GetWorkerBatchSize() uint32 {
	if m != nil && m.WorkerBatchSize != nil {
		return *m.WorkerBatchSize
	}
	return 0
}

func (m *Settings) GetDocBatchSizeKb() uint32 {
	if m != nil && m.DocBatchSizeKb != nil {
		return *m.DocBatchSizeKb
	}
	return 0
}

func (m *Settings) GetFailureRestartInterval() uint32 {
	if m != nil && m.FailureRestartInterval != nil {
		return *m.FailureRestartInterval
	}
	return 0
}

func (m *Settings) GetOptimisticReplicationThreshold() uint32 {
	if m != nil && m.OptimisticReplicationThreshold != nil {
		return *m.OptimisticReplicationThreshold
	}
	return 0
}

func (m *Settings) GetHttpConnections() uint32 {
	if m != nil && m.HttpConnections != nil {
		return *m.HttpConnections
	}
	return 0
}

func (m *Settings) GetSourceNozzlePerNode() uint32 {
	if m != nil && m.SourceNozzlePerNode != nil {
		return *m.SourceNozzlePerNode
	}
	return 0
}

func (m *Settings) GetTargetNozzlePerNode() uint32 {
	if m != nil && m.TargetNozzlePerNode != nil {
		return *m.TargetNozzlePerNode
	}
	return 0
}

func (m *Settings) GetMaxExpectedReplicationLag() uint32 {
	if m != nil && m.MaxExpectedReplicationLag != nil {
		return *m.MaxExpectedReplicationLag
	}
	return 0
}

func (m *Settings) GetTimeoutPercentageCap() uint32 {
	if m != nil && m.TimeoutPercentageCap != nil {
		return *m.TimeoutPercentageCap
	}
	return 0
}

// request to change global settings of all replications in the cluster
type ChangeGlobalSettingsRequest struct {
	Settings         *Settings `protobuf:"bytes,1,req,name=settings" json:"settings,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *ChangeGlobalSettingsRequest) Reset()         { *m = ChangeGlobalSettingsRequest{} }
func (m *ChangeGlobalSettingsRequest) String() string { return proto.CompactTextString(m) }
func (*ChangeGlobalSettingsRequest) ProtoMessage()    {}

func (m *ChangeGlobalSettingsRequest) GetSettings() *Settings {
	if m != nil {
		return m.Settings
	}
	return nil
}

// request to change settings of an individual replication
type ChangeReplicationSettingsRequest struct {
	Id               *string   `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
	Settings         *Settings `protobuf:"bytes,2,req,name=settings" json:"settings,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *ChangeReplicationSettingsRequest) Reset()         { *m = ChangeReplicationSettingsRequest{} }
func (m *ChangeReplicationSettingsRequest) String() string { return proto.CompactTextString(m) }
func (*ChangeReplicationSettingsRequest) ProtoMessage()    {}

func (m *ChangeReplicationSettingsRequest) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *ChangeReplicationSettingsRequest) GetSettings() *Settings {
	if m != nil {
		return m.Settings
	}
	return nil
}

// request to change internal settings of all replications in the cluster.
// effectively the same as ChangeGlobalSettingsRequest but with a different url
type ChangeInternalSettingsRequest struct {
	Settings         *Settings `protobuf:"bytes,1,req,name=settings" json:"settings,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *ChangeInternalSettingsRequest) Reset()         { *m = ChangeInternalSettingsRequest{} }
func (m *ChangeInternalSettingsRequest) String() string { return proto.CompactTextString(m) }
func (*ChangeInternalSettingsRequest) ProtoMessage()    {}

func (m *ChangeInternalSettingsRequest) GetSettings() *Settings {
	if m != nil {
		return m.Settings
	}
	return nil
}

type GetStatisticsRequest struct {
	Uuid             *string                     `protobuf:"bytes,1,req,name=uuid" json:"uuid,omitempty"`
	FromBucket       *string                     `protobuf:"bytes,2,req,name=fromBucket" json:"fromBucket,omitempty"`
	ToBucket         *string                     `protobuf:"bytes,3,req,name=toBucket" json:"toBucket,omitempty"`
	Stats            *GetStatisticsRequest_Stats `protobuf:"varint,4,req,name=stats,enum=protobuf.GetStatisticsRequest_Stats" json:"stats,omitempty"`
	XXX_unrecognized []byte                      `json:"-"`
}

func (m *GetStatisticsRequest) Reset()         { *m = GetStatisticsRequest{} }
func (m *GetStatisticsRequest) String() string { return proto.CompactTextString(m) }
func (*GetStatisticsRequest) ProtoMessage()    {}

func (m *GetStatisticsRequest) GetUuid() string {
	if m != nil && m.Uuid != nil {
		return *m.Uuid
	}
	return ""
}

func (m *GetStatisticsRequest) GetFromBucket() string {
	if m != nil && m.FromBucket != nil {
		return *m.FromBucket
	}
	return ""
}

func (m *GetStatisticsRequest) GetToBucket() string {
	if m != nil && m.ToBucket != nil {
		return *m.ToBucket
	}
	return ""
}

func (m *GetStatisticsRequest) GetStats() GetStatisticsRequest_Stats {
	if m != nil && m.Stats != nil {
		return *m.Stats
	}
	return GetStatisticsRequest_docs_written
}

func init() {
	proto.RegisterEnum("protobuf.CreateReplicationRequest_Mode", CreateReplicationRequest_Mode_name, CreateReplicationRequest_Mode_value)
	proto.RegisterEnum("protobuf.GetStatisticsRequest_Stats", GetStatisticsRequest_Stats_name, GetStatisticsRequest_Stats_value)
}