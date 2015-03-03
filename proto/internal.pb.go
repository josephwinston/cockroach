// Code generated by protoc-gen-gogo.
// source: internal.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/gogo/protobuf/proto"
import math "math"

// discarding unused import gogoproto "github.com/gogo/protobuf/gogoproto/gogo.pb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

// ReplicaChangeType is a parameter of InternalChangeReplicasRequest.
type ReplicaChangeType int32

const (
	ADD_REPLICA    ReplicaChangeType = 0
	REMOVE_REPLICA ReplicaChangeType = 1
)

var ReplicaChangeType_name = map[int32]string{
	0: "ADD_REPLICA",
	1: "REMOVE_REPLICA",
}
var ReplicaChangeType_value = map[string]int32{
	"ADD_REPLICA":    0,
	"REMOVE_REPLICA": 1,
}

func (x ReplicaChangeType) Enum() *ReplicaChangeType {
	p := new(ReplicaChangeType)
	*p = x
	return p
}
func (x ReplicaChangeType) String() string {
	return proto1.EnumName(ReplicaChangeType_name, int32(x))
}
func (x *ReplicaChangeType) UnmarshalJSON(data []byte) error {
	value, err := proto1.UnmarshalJSONEnum(ReplicaChangeType_value, data, "ReplicaChangeType")
	if err != nil {
		return err
	}
	*x = ReplicaChangeType(value)
	return nil
}

// InternalValueType defines a set of string constants placed in the "tag" field
// of Value messages which are created internally. These are defined as a
// protocol buffer enumeration so that they can be used portably between our Go
// and C code.
type InternalValueType int32

const (
	// _CR_TS is applied to values which contain InternalTimeSeriesData. This
	// tag is used by the RocksDB Merge Operator to perform a specialized merge
	// for this data.
	_CR_TS InternalValueType = 1
)

var InternalValueType_name = map[int32]string{
	1: "_CR_TS",
}
var InternalValueType_value = map[string]int32{
	"_CR_TS": 1,
}

func (x InternalValueType) Enum() *InternalValueType {
	p := new(InternalValueType)
	*p = x
	return p
}
func (x InternalValueType) String() string {
	return proto1.EnumName(InternalValueType_name, int32(x))
}
func (x *InternalValueType) UnmarshalJSON(data []byte) error {
	value, err := proto1.UnmarshalJSONEnum(InternalValueType_value, data, "InternalValueType")
	if err != nil {
		return err
	}
	*x = InternalValueType(value)
	return nil
}

// An InternalRangeLookupRequest is arguments to the
// InternalRangeLookup() method. It specifies the key for which the
// containing range is being requested, and the maximum number of
// total range descriptors that should be returned, if there are
// additional consecutive addressable ranges. Specify max_ranges > 1
// to pre-fill the range descriptor cache.
type InternalRangeLookupRequest struct {
	RequestHeader    `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	MaxRanges        int32  `protobuf:"varint,2,opt,name=max_ranges" json:"max_ranges"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalRangeLookupRequest) Reset()         { *m = InternalRangeLookupRequest{} }
func (m *InternalRangeLookupRequest) String() string { return proto1.CompactTextString(m) }
func (*InternalRangeLookupRequest) ProtoMessage()    {}

func (m *InternalRangeLookupRequest) GetMaxRanges() int32 {
	if m != nil {
		return m.MaxRanges
	}
	return 0
}

// An InternalRangeLookupResponse is the return value from the
// InternalRangeLookup() method. It returns metadata for the range
// containing the requested key, optionally returning the metadata for
// additional consecutive ranges beyond the requested range to pre-fill
// the range descriptor cache.
type InternalRangeLookupResponse struct {
	ResponseHeader   `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	Ranges           []RangeDescriptor `protobuf:"bytes,2,rep,name=ranges" json:"ranges"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *InternalRangeLookupResponse) Reset()         { *m = InternalRangeLookupResponse{} }
func (m *InternalRangeLookupResponse) String() string { return proto1.CompactTextString(m) }
func (*InternalRangeLookupResponse) ProtoMessage()    {}

func (m *InternalRangeLookupResponse) GetRanges() []RangeDescriptor {
	if m != nil {
		return m.Ranges
	}
	return nil
}

// An InternalHeartbeatTxnRequest is arguments to the
// InternalHeartbeatTxn() method. It's sent by transaction
// coordinators to let the system know that the transaction is still
// ongoing. Note that this heartbeat message is different from the
// heartbeat message in the gossip protocol.
type InternalHeartbeatTxnRequest struct {
	RequestHeader    `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalHeartbeatTxnRequest) Reset()         { *m = InternalHeartbeatTxnRequest{} }
func (m *InternalHeartbeatTxnRequest) String() string { return proto1.CompactTextString(m) }
func (*InternalHeartbeatTxnRequest) ProtoMessage()    {}

// An InternalHeartbeatTxnResponse is the return value from the
// InternalHeartbeatTxn() method. It returns the transaction info in
// the response header. The returned transaction lets the coordinator
// know the disposition of the transaction (i.e. aborted, committed or
// pending).
type InternalHeartbeatTxnResponse struct {
	ResponseHeader   `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalHeartbeatTxnResponse) Reset()         { *m = InternalHeartbeatTxnResponse{} }
func (m *InternalHeartbeatTxnResponse) String() string { return proto1.CompactTextString(m) }
func (*InternalHeartbeatTxnResponse) ProtoMessage()    {}

// An InternalGCRequest is arguments to the InternalGC() method. It's
// sent by range leaders after scanning range data to find expired
// MVCC values.
type InternalGCRequest struct {
	RequestHeader    `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	ScanMeta         ScanMetadata              `protobuf:"bytes,2,opt,name=scan_meta" json:"scan_meta"`
	Keys             []InternalGCRequest_GCKey `protobuf:"bytes,3,rep,name=keys" json:"keys"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *InternalGCRequest) Reset()         { *m = InternalGCRequest{} }
func (m *InternalGCRequest) String() string { return proto1.CompactTextString(m) }
func (*InternalGCRequest) ProtoMessage()    {}

func (m *InternalGCRequest) GetScanMeta() ScanMetadata {
	if m != nil {
		return m.ScanMeta
	}
	return ScanMetadata{}
}

func (m *InternalGCRequest) GetKeys() []InternalGCRequest_GCKey {
	if m != nil {
		return m.Keys
	}
	return nil
}

type InternalGCRequest_GCKey struct {
	Key              Key       `protobuf:"bytes,1,opt,name=key,customtype=Key" json:"key"`
	Timestamp        Timestamp `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *InternalGCRequest_GCKey) Reset()         { *m = InternalGCRequest_GCKey{} }
func (m *InternalGCRequest_GCKey) String() string { return proto1.CompactTextString(m) }
func (*InternalGCRequest_GCKey) ProtoMessage()    {}

func (m *InternalGCRequest_GCKey) GetTimestamp() Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return Timestamp{}
}

// An InternalGCResponse is the return value from the InternalGC()
// method.
type InternalGCResponse struct {
	ResponseHeader   `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalGCResponse) Reset()         { *m = InternalGCResponse{} }
func (m *InternalGCResponse) String() string { return proto1.CompactTextString(m) }
func (*InternalGCResponse) ProtoMessage()    {}

// An InternalPushTxnRequest is arguments to the InternalPushTxn()
// method. It's sent by readers or writers which have encountered an
// "intent" laid down by another transaction. The goal is to resolve
// the conflict. Note that args.Key should be set to the txn ID of
// args.PusheeTxn, not args.Txn, as is usual. This RPC is addressed
// to the range which owns the pushee's txn record.
//
// Resolution is trivial if the txn which owns the intent has either
// been committed or aborted already. Otherwise, the existing txn can
// either be aborted (for write/write conflicts), or its commit
// timestamp can be moved forward (for read/write conflicts). The
// course of action is determined by the owning txn's status and also
// by comparing priorities.
type InternalPushTxnRequest struct {
	RequestHeader `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	PusheeTxn     Transaction `protobuf:"bytes,2,opt,name=pushee_txn" json:"pushee_txn"`
	// Set to true to request that the PushTxn be aborted if possible.
	// This is done in the event of a writer conflicting with PusheeTxn.
	// Readers set this to false and instead attempt to move PusheeTxn's
	// commit timestamp forward.
	Abort            bool   `protobuf:"varint,3,opt" json:"Abort"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalPushTxnRequest) Reset()         { *m = InternalPushTxnRequest{} }
func (m *InternalPushTxnRequest) String() string { return proto1.CompactTextString(m) }
func (*InternalPushTxnRequest) ProtoMessage()    {}

func (m *InternalPushTxnRequest) GetPusheeTxn() Transaction {
	if m != nil {
		return m.PusheeTxn
	}
	return Transaction{}
}

func (m *InternalPushTxnRequest) GetAbort() bool {
	if m != nil {
		return m.Abort
	}
	return false
}

// An InternalPushTxnResponse is the return value from the
// InternalPushTxn() method. It returns success and the resulting
// state of PusheeTxn if the conflict was resolved in favor of the
// caller; the caller should subsequently invoke
// InternalResolveIntent() on the conflicted key. It returns an error
// otherwise.
type InternalPushTxnResponse struct {
	ResponseHeader `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	// Txn is non-nil if the transaction could be heartbeat and contains
	// the current value of the transaction.
	PusheeTxn        *Transaction `protobuf:"bytes,2,opt,name=pushee_txn" json:"pushee_txn,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *InternalPushTxnResponse) Reset()         { *m = InternalPushTxnResponse{} }
func (m *InternalPushTxnResponse) String() string { return proto1.CompactTextString(m) }
func (*InternalPushTxnResponse) ProtoMessage()    {}

func (m *InternalPushTxnResponse) GetPusheeTxn() *Transaction {
	if m != nil {
		return m.PusheeTxn
	}
	return nil
}

// An InternalResolveIntentRequest is arguments to the
// InternalResolveIntent() method. It is sent by transaction
// coordinators and after success calling InternalPushTxn to clean up
// write intents: either to remove them or commit them.
type InternalResolveIntentRequest struct {
	RequestHeader    `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalResolveIntentRequest) Reset()         { *m = InternalResolveIntentRequest{} }
func (m *InternalResolveIntentRequest) String() string { return proto1.CompactTextString(m) }
func (*InternalResolveIntentRequest) ProtoMessage()    {}

// An InternalResolveIntentResponse is the return value from the
// InternalResolveIntent() method.
type InternalResolveIntentResponse struct {
	ResponseHeader   `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalResolveIntentResponse) Reset()         { *m = InternalResolveIntentResponse{} }
func (m *InternalResolveIntentResponse) String() string { return proto1.CompactTextString(m) }
func (*InternalResolveIntentResponse) ProtoMessage()    {}

// An InternalMergeRequest contains arguments to the InternalMerge() method. It
// specifies a key and a value which should be merged into the existing value at
// that key.
type InternalMergeRequest struct {
	RequestHeader    `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	Value            Value  `protobuf:"bytes,2,opt,name=value" json:"value"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalMergeRequest) Reset()         { *m = InternalMergeRequest{} }
func (m *InternalMergeRequest) String() string { return proto1.CompactTextString(m) }
func (*InternalMergeRequest) ProtoMessage()    {}

func (m *InternalMergeRequest) GetValue() Value {
	if m != nil {
		return m.Value
	}
	return Value{}
}

// InternalMergeResponse is the response to an InternalMerge() operation.
type InternalMergeResponse struct {
	ResponseHeader   `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalMergeResponse) Reset()         { *m = InternalMergeResponse{} }
func (m *InternalMergeResponse) String() string { return proto1.CompactTextString(m) }
func (*InternalMergeResponse) ProtoMessage()    {}

// InternalTruncateLogRequest is used to remove a prefix of the raft log. While there
// is no requirement for correctness that the raft log truncation be synchronized across
// replicas, it is nice to preserve the property that all replicas of a range are as close
// to identical as possible. The raft leader can also inform decisions about the cutoff point
// with its knowledge of the replicas' acknowledgement status.
type InternalTruncateLogRequest struct {
	RequestHeader `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	// Log entries < this index are to be discarded.
	Index            uint64 `protobuf:"varint,2,opt,name=index" json:"index"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalTruncateLogRequest) Reset()         { *m = InternalTruncateLogRequest{} }
func (m *InternalTruncateLogRequest) String() string { return proto1.CompactTextString(m) }
func (*InternalTruncateLogRequest) ProtoMessage()    {}

func (m *InternalTruncateLogRequest) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

// InternalTruncateLogResponse is the response to an InternalTruncateLog() operation.
type InternalTruncateLogResponse struct {
	ResponseHeader   `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalTruncateLogResponse) Reset()         { *m = InternalTruncateLogResponse{} }
func (m *InternalTruncateLogResponse) String() string { return proto1.CompactTextString(m) }
func (*InternalTruncateLogResponse) ProtoMessage()    {}

// InternalChangeReplicasRequest is used to add or remove a replica from a raft group.
// Only one ChangeReplicas operation can be in progress at a time; a proposed
// change will fail if the previous change has not yet completed.
type InternalChangeReplicasRequest struct {
	RequestHeader `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	NodeID        NodeID            `protobuf:"varint,2,opt,name=node_id,customtype=NodeID" json:"node_id"`
	StoreID       StoreID           `protobuf:"varint,3,opt,name=store_id,customtype=StoreID" json:"store_id"`
	ChangeType    ReplicaChangeType `protobuf:"varint,4,opt,name=change_type,enum=proto.ReplicaChangeType" json:"change_type"`
	// This field gets filled in as the request passes through raft.
	// It contains the current committed members of the group (equivalent to ConfState.Nodes)
	// to guard against any drift from incremental processing of changes.
	Nodes            []uint64 `protobuf:"varint,5,rep,name=nodes" json:"nodes,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *InternalChangeReplicasRequest) Reset()         { *m = InternalChangeReplicasRequest{} }
func (m *InternalChangeReplicasRequest) String() string { return proto1.CompactTextString(m) }
func (*InternalChangeReplicasRequest) ProtoMessage()    {}

func (m *InternalChangeReplicasRequest) GetChangeType() ReplicaChangeType {
	if m != nil {
		return m.ChangeType
	}
	return ADD_REPLICA
}

func (m *InternalChangeReplicasRequest) GetNodes() []uint64 {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type InternalChangeReplicasResponse struct {
	ResponseHeader   `protobuf:"bytes,1,opt,name=header,embedded=header" json:"header"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InternalChangeReplicasResponse) Reset()         { *m = InternalChangeReplicasResponse{} }
func (m *InternalChangeReplicasResponse) String() string { return proto1.CompactTextString(m) }
func (*InternalChangeReplicasResponse) ProtoMessage()    {}

// A ReadWriteCmdResponse is a union type containing instances of all
// mutating commands. Note that any entry added here must be handled
// in storage/engine/db.cc in GetResponseHeader().
type ReadWriteCmdResponse struct {
	Put                   *PutResponse                   `protobuf:"bytes,1,opt,name=put" json:"put,omitempty"`
	ConditionalPut        *ConditionalPutResponse        `protobuf:"bytes,2,opt,name=conditional_put" json:"conditional_put,omitempty"`
	Increment             *IncrementResponse             `protobuf:"bytes,3,opt,name=increment" json:"increment,omitempty"`
	Delete                *DeleteResponse                `protobuf:"bytes,4,opt,name=delete" json:"delete,omitempty"`
	DeleteRange           *DeleteRangeResponse           `protobuf:"bytes,5,opt,name=delete_range" json:"delete_range,omitempty"`
	EndTransaction        *EndTransactionResponse        `protobuf:"bytes,6,opt,name=end_transaction" json:"end_transaction,omitempty"`
	ReapQueue             *ReapQueueResponse             `protobuf:"bytes,7,opt,name=reap_queue" json:"reap_queue,omitempty"`
	EnqueueUpdate         *EnqueueUpdateResponse         `protobuf:"bytes,8,opt,name=enqueue_update" json:"enqueue_update,omitempty"`
	EnqueueMessage        *EnqueueMessageResponse        `protobuf:"bytes,9,opt,name=enqueue_message" json:"enqueue_message,omitempty"`
	InternalHeartbeatTxn  *InternalHeartbeatTxnResponse  `protobuf:"bytes,10,opt,name=internal_heartbeat_txn" json:"internal_heartbeat_txn,omitempty"`
	InternalPushTxn       *InternalPushTxnResponse       `protobuf:"bytes,11,opt,name=internal_push_txn" json:"internal_push_txn,omitempty"`
	InternalResolveIntent *InternalResolveIntentResponse `protobuf:"bytes,12,opt,name=internal_resolve_intent" json:"internal_resolve_intent,omitempty"`
	InternalMerge         *InternalMergeResponse         `protobuf:"bytes,13,opt,name=internal_merge" json:"internal_merge,omitempty"`
	InternalTruncateLog   *InternalTruncateLogResponse   `protobuf:"bytes,14,opt,name=internal_truncate_log" json:"internal_truncate_log,omitempty"`
	InternalGc            *InternalGCResponse            `protobuf:"bytes,15,opt,name=internal_gc" json:"internal_gc,omitempty"`
	XXX_unrecognized      []byte                         `json:"-"`
}

func (m *ReadWriteCmdResponse) Reset()         { *m = ReadWriteCmdResponse{} }
func (m *ReadWriteCmdResponse) String() string { return proto1.CompactTextString(m) }
func (*ReadWriteCmdResponse) ProtoMessage()    {}

func (m *ReadWriteCmdResponse) GetPut() *PutResponse {
	if m != nil {
		return m.Put
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetConditionalPut() *ConditionalPutResponse {
	if m != nil {
		return m.ConditionalPut
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetIncrement() *IncrementResponse {
	if m != nil {
		return m.Increment
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetDelete() *DeleteResponse {
	if m != nil {
		return m.Delete
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetDeleteRange() *DeleteRangeResponse {
	if m != nil {
		return m.DeleteRange
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetEndTransaction() *EndTransactionResponse {
	if m != nil {
		return m.EndTransaction
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetReapQueue() *ReapQueueResponse {
	if m != nil {
		return m.ReapQueue
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetEnqueueUpdate() *EnqueueUpdateResponse {
	if m != nil {
		return m.EnqueueUpdate
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetEnqueueMessage() *EnqueueMessageResponse {
	if m != nil {
		return m.EnqueueMessage
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetInternalHeartbeatTxn() *InternalHeartbeatTxnResponse {
	if m != nil {
		return m.InternalHeartbeatTxn
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetInternalPushTxn() *InternalPushTxnResponse {
	if m != nil {
		return m.InternalPushTxn
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetInternalResolveIntent() *InternalResolveIntentResponse {
	if m != nil {
		return m.InternalResolveIntent
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetInternalMerge() *InternalMergeResponse {
	if m != nil {
		return m.InternalMerge
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetInternalTruncateLog() *InternalTruncateLogResponse {
	if m != nil {
		return m.InternalTruncateLog
	}
	return nil
}

func (m *ReadWriteCmdResponse) GetInternalGc() *InternalGCResponse {
	if m != nil {
		return m.InternalGc
	}
	return nil
}

// An InternalRaftCommandUnion is the union of all commands which can be
// sent via raft.
type InternalRaftCommandUnion struct {
	// Non-batched external requests. This section is the same as RequestUnion.
	Contains       *ContainsRequest       `protobuf:"bytes,1,opt,name=contains" json:"contains,omitempty"`
	Get            *GetRequest            `protobuf:"bytes,2,opt,name=get" json:"get,omitempty"`
	Put            *PutRequest            `protobuf:"bytes,3,opt,name=put" json:"put,omitempty"`
	ConditionalPut *ConditionalPutRequest `protobuf:"bytes,4,opt,name=conditional_put" json:"conditional_put,omitempty"`
	Increment      *IncrementRequest      `protobuf:"bytes,5,opt,name=increment" json:"increment,omitempty"`
	Delete         *DeleteRequest         `protobuf:"bytes,6,opt,name=delete" json:"delete,omitempty"`
	DeleteRange    *DeleteRangeRequest    `protobuf:"bytes,7,opt,name=delete_range" json:"delete_range,omitempty"`
	Scan           *ScanRequest           `protobuf:"bytes,8,opt,name=scan" json:"scan,omitempty"`
	EndTransaction *EndTransactionRequest `protobuf:"bytes,9,opt,name=end_transaction" json:"end_transaction,omitempty"`
	ReapQueue      *ReapQueueRequest      `protobuf:"bytes,10,opt,name=reap_queue" json:"reap_queue,omitempty"`
	EnqueueUpdate  *EnqueueUpdateRequest  `protobuf:"bytes,11,opt,name=enqueue_update" json:"enqueue_update,omitempty"`
	EnqueueMessage *EnqueueMessageRequest `protobuf:"bytes,12,opt,name=enqueue_message" json:"enqueue_message,omitempty"`
	// Other requests. Allow a gap in tag numbers so the previous list can
	// be copy/pasted from RequestUnion.
	Batch                  *BatchRequest                  `protobuf:"bytes,30,opt,name=batch" json:"batch,omitempty"`
	InternalRangeLookup    *InternalRangeLookupRequest    `protobuf:"bytes,31,opt,name=internal_range_lookup" json:"internal_range_lookup,omitempty"`
	InternalHeartbeatTxn   *InternalHeartbeatTxnRequest   `protobuf:"bytes,32,opt,name=internal_heartbeat_txn" json:"internal_heartbeat_txn,omitempty"`
	InternalPushTxn        *InternalPushTxnRequest        `protobuf:"bytes,33,opt,name=internal_push_txn" json:"internal_push_txn,omitempty"`
	InternalResolveIntent  *InternalResolveIntentRequest  `protobuf:"bytes,34,opt,name=internal_resolve_intent" json:"internal_resolve_intent,omitempty"`
	InternalMergeResponse  *InternalMergeRequest          `protobuf:"bytes,35,opt,name=internal_merge_response" json:"internal_merge_response,omitempty"`
	InternalTruncateLog    *InternalTruncateLogRequest    `protobuf:"bytes,36,opt,name=internal_truncate_log" json:"internal_truncate_log,omitempty"`
	InternalGc             *InternalGCRequest             `protobuf:"bytes,37,opt,name=internal_gc" json:"internal_gc,omitempty"`
	InternalChangeReplicas *InternalChangeReplicasRequest `protobuf:"bytes,38,opt,name=internal_change_replicas" json:"internal_change_replicas,omitempty"`
	XXX_unrecognized       []byte                         `json:"-"`
}

func (m *InternalRaftCommandUnion) Reset()         { *m = InternalRaftCommandUnion{} }
func (m *InternalRaftCommandUnion) String() string { return proto1.CompactTextString(m) }
func (*InternalRaftCommandUnion) ProtoMessage()    {}

func (m *InternalRaftCommandUnion) GetContains() *ContainsRequest {
	if m != nil {
		return m.Contains
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetGet() *GetRequest {
	if m != nil {
		return m.Get
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetPut() *PutRequest {
	if m != nil {
		return m.Put
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetConditionalPut() *ConditionalPutRequest {
	if m != nil {
		return m.ConditionalPut
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetIncrement() *IncrementRequest {
	if m != nil {
		return m.Increment
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetDelete() *DeleteRequest {
	if m != nil {
		return m.Delete
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetDeleteRange() *DeleteRangeRequest {
	if m != nil {
		return m.DeleteRange
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetScan() *ScanRequest {
	if m != nil {
		return m.Scan
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetEndTransaction() *EndTransactionRequest {
	if m != nil {
		return m.EndTransaction
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetReapQueue() *ReapQueueRequest {
	if m != nil {
		return m.ReapQueue
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetEnqueueUpdate() *EnqueueUpdateRequest {
	if m != nil {
		return m.EnqueueUpdate
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetEnqueueMessage() *EnqueueMessageRequest {
	if m != nil {
		return m.EnqueueMessage
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetBatch() *BatchRequest {
	if m != nil {
		return m.Batch
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetInternalRangeLookup() *InternalRangeLookupRequest {
	if m != nil {
		return m.InternalRangeLookup
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetInternalHeartbeatTxn() *InternalHeartbeatTxnRequest {
	if m != nil {
		return m.InternalHeartbeatTxn
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetInternalPushTxn() *InternalPushTxnRequest {
	if m != nil {
		return m.InternalPushTxn
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetInternalResolveIntent() *InternalResolveIntentRequest {
	if m != nil {
		return m.InternalResolveIntent
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetInternalMergeResponse() *InternalMergeRequest {
	if m != nil {
		return m.InternalMergeResponse
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetInternalTruncateLog() *InternalTruncateLogRequest {
	if m != nil {
		return m.InternalTruncateLog
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetInternalGc() *InternalGCRequest {
	if m != nil {
		return m.InternalGc
	}
	return nil
}

func (m *InternalRaftCommandUnion) GetInternalChangeReplicas() *InternalChangeReplicasRequest {
	if m != nil {
		return m.InternalChangeReplicas
	}
	return nil
}

// An InternalRaftCommand is a command which can be serialized and
// sent via raft.
type InternalRaftCommand struct {
	RaftID           int64                    `protobuf:"varint,2,opt,name=raft_id" json:"raft_id"`
	Cmd              InternalRaftCommandUnion `protobuf:"bytes,3,opt,name=cmd" json:"cmd"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *InternalRaftCommand) Reset()         { *m = InternalRaftCommand{} }
func (m *InternalRaftCommand) String() string { return proto1.CompactTextString(m) }
func (*InternalRaftCommand) ProtoMessage()    {}

func (m *InternalRaftCommand) GetRaftID() int64 {
	if m != nil {
		return m.RaftID
	}
	return 0
}

func (m *InternalRaftCommand) GetCmd() InternalRaftCommandUnion {
	if m != nil {
		return m.Cmd
	}
	return InternalRaftCommandUnion{}
}

// InternalTimeSeriesData is a collection of data samples for some measurable
// value, where each sample is taken over a uniform time interval.
//
// The collection itself contains a start timestamp (in seconds since the unix
// epoch) and a sample duration (in milliseconds). Each sample in the collection
// will contain a positive integer offset that indicates the length of time
// between the start_timestamp of the collection and the time when the sample
// began, expressed as an whole number of sample intervals. For example, if the
// sample duration is 60000 (indicating 1 minute), then a contained sample with
// an offset value of 5 begins (5*60000ms = 300000ms = 5 minutes) after the
// start timestamp of this data.
//
// This is meant to be an efficient internal representation of time series data,
// ensuring that very little redundant data is stored on disk. With this goal in
// mind, this message does not identify the variable which is actually being
// measured; that information is expected be encoded in the key where this
// message is stored.
type InternalTimeSeriesData struct {
	// Holds a wall time, expressed as a unix epoch time in nanoseconds. This
	// represents the earliest possible timestamp for a sample within the
	// collection.
	StartTimestampNanos int64 `protobuf:"varint,1,opt,name=start_timestamp_nanos" json:"start_timestamp_nanos"`
	// The duration of each sample interval, expressed in nanoseconds.
	SampleDurationNanos int64 `protobuf:"varint,2,opt,name=sample_duration_nanos" json:"sample_duration_nanos"`
	// The actual data samples for this metric.
	Samples          []*InternalTimeSeriesSample `protobuf:"bytes,3,rep,name=samples" json:"samples,omitempty"`
	XXX_unrecognized []byte                      `json:"-"`
}

func (m *InternalTimeSeriesData) Reset()         { *m = InternalTimeSeriesData{} }
func (m *InternalTimeSeriesData) String() string { return proto1.CompactTextString(m) }
func (*InternalTimeSeriesData) ProtoMessage()    {}

func (m *InternalTimeSeriesData) GetStartTimestampNanos() int64 {
	if m != nil {
		return m.StartTimestampNanos
	}
	return 0
}

func (m *InternalTimeSeriesData) GetSampleDurationNanos() int64 {
	if m != nil {
		return m.SampleDurationNanos
	}
	return 0
}

func (m *InternalTimeSeriesData) GetSamples() []*InternalTimeSeriesSample {
	if m != nil {
		return m.Samples
	}
	return nil
}

// A InternalTimeSeriesSample represents data gathered from multiple
// measurements of a variable value over a given period of time. The length of
// that period of time is stored in an InternalTimeSeriesData message; a sample
// cannot be interpreted correctly without a start timestamp and sample
// duration.
//
// Each sample may contain data gathered from multiple measurements of the same
// variable, as long as all of those measurements occured within the sample
// period. The sample stores several aggregated values from these measurements:
// - The sum of all measured values
// - A count of all measurements taken
// - The maximum individual measurement seen
// - The minimum individual measurement seen
//
// If zero measurements are present in a sample, then it should be omitted
// entirely from any collection it would be a part of.
//
// If the count of measurements is 1, then max and min fields may be omitted
// and assumed equal to the sum field.
//
// The variable being measured may be either an integer or a floating point;
// therefore, there are two fields each for "sum", "max" and "min" to hold
// either an integer or floating point number. In practice, only one set of
// these fields should be present for any individual sample; however, int and
// float values are recorded in parallel, allowing clients to write both floats
// and integers to the same value. These are recorded separately to retain
// precision, but are easily combined by higher-level logic at query time.
type InternalTimeSeriesSample struct {
	// Temporal offset from the "start_timestamp" of the InternalTimeSeriesData
	// collection this data point is part in. The units of this value are
	// determined by the value of the "sample_duration_milliseconds" field of
	// the TimeSeriesData collection.
	Offset int32 `protobuf:"varint,1,opt,name=offset" json:"offset"`
	// Count of integer measurements taken within this sample.
	IntCount uint32 `protobuf:"varint,2,opt,name=int_count" json:"int_count"`
	// Sum of all integer measurements.
	IntSum *int64 `protobuf:"varint,3,opt,name=int_sum" json:"int_sum,omitempty"`
	// Maximum encountered integer measurement in this sample.
	IntMax *int64 `protobuf:"varint,4,opt,name=int_max" json:"int_max,omitempty"`
	// Minimum encountered integer measurement in this sample.
	IntMin *int64 `protobuf:"varint,5,opt,name=int_min" json:"int_min,omitempty"`
	// Count of floating point measurements taken within this sample.
	FloatCount uint32 `protobuf:"varint,6,opt,name=float_count" json:"float_count"`
	// Sum of all floating point measurements.
	FloatSum *float32 `protobuf:"fixed32,7,opt,name=float_sum" json:"float_sum,omitempty"`
	// Maximum encountered floating point measurement in this sample.
	FloatMax *float32 `protobuf:"fixed32,8,opt,name=float_max" json:"float_max,omitempty"`
	// Minimum encountered floating point measurement in this sample.
	FloatMin         *float32 `protobuf:"fixed32,9,opt,name=float_min" json:"float_min,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *InternalTimeSeriesSample) Reset()         { *m = InternalTimeSeriesSample{} }
func (m *InternalTimeSeriesSample) String() string { return proto1.CompactTextString(m) }
func (*InternalTimeSeriesSample) ProtoMessage()    {}

func (m *InternalTimeSeriesSample) GetOffset() int32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *InternalTimeSeriesSample) GetIntCount() uint32 {
	if m != nil {
		return m.IntCount
	}
	return 0
}

func (m *InternalTimeSeriesSample) GetIntSum() int64 {
	if m != nil && m.IntSum != nil {
		return *m.IntSum
	}
	return 0
}

func (m *InternalTimeSeriesSample) GetIntMax() int64 {
	if m != nil && m.IntMax != nil {
		return *m.IntMax
	}
	return 0
}

func (m *InternalTimeSeriesSample) GetIntMin() int64 {
	if m != nil && m.IntMin != nil {
		return *m.IntMin
	}
	return 0
}

func (m *InternalTimeSeriesSample) GetFloatCount() uint32 {
	if m != nil {
		return m.FloatCount
	}
	return 0
}

func (m *InternalTimeSeriesSample) GetFloatSum() float32 {
	if m != nil && m.FloatSum != nil {
		return *m.FloatSum
	}
	return 0
}

func (m *InternalTimeSeriesSample) GetFloatMax() float32 {
	if m != nil && m.FloatMax != nil {
		return *m.FloatMax
	}
	return 0
}

func (m *InternalTimeSeriesSample) GetFloatMin() float32 {
	if m != nil && m.FloatMin != nil {
		return *m.FloatMin
	}
	return 0
}

func init() {
	proto1.RegisterEnum("proto.ReplicaChangeType", ReplicaChangeType_name, ReplicaChangeType_value)
	proto1.RegisterEnum("proto.InternalValueType", InternalValueType_name, InternalValueType_value)
}
func (this *ReadWriteCmdResponse) GetValue() interface{} {
	if this.Put != nil {
		return this.Put
	}
	if this.ConditionalPut != nil {
		return this.ConditionalPut
	}
	if this.Increment != nil {
		return this.Increment
	}
	if this.Delete != nil {
		return this.Delete
	}
	if this.DeleteRange != nil {
		return this.DeleteRange
	}
	if this.EndTransaction != nil {
		return this.EndTransaction
	}
	if this.ReapQueue != nil {
		return this.ReapQueue
	}
	if this.EnqueueUpdate != nil {
		return this.EnqueueUpdate
	}
	if this.EnqueueMessage != nil {
		return this.EnqueueMessage
	}
	if this.InternalHeartbeatTxn != nil {
		return this.InternalHeartbeatTxn
	}
	if this.InternalPushTxn != nil {
		return this.InternalPushTxn
	}
	if this.InternalResolveIntent != nil {
		return this.InternalResolveIntent
	}
	if this.InternalMerge != nil {
		return this.InternalMerge
	}
	if this.InternalTruncateLog != nil {
		return this.InternalTruncateLog
	}
	if this.InternalGc != nil {
		return this.InternalGc
	}
	return nil
}

func (this *ReadWriteCmdResponse) SetValue(value interface{}) bool {
	switch vt := value.(type) {
	case *PutResponse:
		this.Put = vt
	case *ConditionalPutResponse:
		this.ConditionalPut = vt
	case *IncrementResponse:
		this.Increment = vt
	case *DeleteResponse:
		this.Delete = vt
	case *DeleteRangeResponse:
		this.DeleteRange = vt
	case *EndTransactionResponse:
		this.EndTransaction = vt
	case *ReapQueueResponse:
		this.ReapQueue = vt
	case *EnqueueUpdateResponse:
		this.EnqueueUpdate = vt
	case *EnqueueMessageResponse:
		this.EnqueueMessage = vt
	case *InternalHeartbeatTxnResponse:
		this.InternalHeartbeatTxn = vt
	case *InternalPushTxnResponse:
		this.InternalPushTxn = vt
	case *InternalResolveIntentResponse:
		this.InternalResolveIntent = vt
	case *InternalMergeResponse:
		this.InternalMerge = vt
	case *InternalTruncateLogResponse:
		this.InternalTruncateLog = vt
	case *InternalGCResponse:
		this.InternalGc = vt
	default:
		return false
	}
	return true
}
func (this *InternalRaftCommandUnion) GetValue() interface{} {
	if this.Contains != nil {
		return this.Contains
	}
	if this.Get != nil {
		return this.Get
	}
	if this.Put != nil {
		return this.Put
	}
	if this.ConditionalPut != nil {
		return this.ConditionalPut
	}
	if this.Increment != nil {
		return this.Increment
	}
	if this.Delete != nil {
		return this.Delete
	}
	if this.DeleteRange != nil {
		return this.DeleteRange
	}
	if this.Scan != nil {
		return this.Scan
	}
	if this.EndTransaction != nil {
		return this.EndTransaction
	}
	if this.ReapQueue != nil {
		return this.ReapQueue
	}
	if this.EnqueueUpdate != nil {
		return this.EnqueueUpdate
	}
	if this.EnqueueMessage != nil {
		return this.EnqueueMessage
	}
	if this.Batch != nil {
		return this.Batch
	}
	if this.InternalRangeLookup != nil {
		return this.InternalRangeLookup
	}
	if this.InternalHeartbeatTxn != nil {
		return this.InternalHeartbeatTxn
	}
	if this.InternalPushTxn != nil {
		return this.InternalPushTxn
	}
	if this.InternalResolveIntent != nil {
		return this.InternalResolveIntent
	}
	if this.InternalMergeResponse != nil {
		return this.InternalMergeResponse
	}
	if this.InternalTruncateLog != nil {
		return this.InternalTruncateLog
	}
	if this.InternalGc != nil {
		return this.InternalGc
	}
	if this.InternalChangeReplicas != nil {
		return this.InternalChangeReplicas
	}
	return nil
}

func (this *InternalRaftCommandUnion) SetValue(value interface{}) bool {
	switch vt := value.(type) {
	case *ContainsRequest:
		this.Contains = vt
	case *GetRequest:
		this.Get = vt
	case *PutRequest:
		this.Put = vt
	case *ConditionalPutRequest:
		this.ConditionalPut = vt
	case *IncrementRequest:
		this.Increment = vt
	case *DeleteRequest:
		this.Delete = vt
	case *DeleteRangeRequest:
		this.DeleteRange = vt
	case *ScanRequest:
		this.Scan = vt
	case *EndTransactionRequest:
		this.EndTransaction = vt
	case *ReapQueueRequest:
		this.ReapQueue = vt
	case *EnqueueUpdateRequest:
		this.EnqueueUpdate = vt
	case *EnqueueMessageRequest:
		this.EnqueueMessage = vt
	case *BatchRequest:
		this.Batch = vt
	case *InternalRangeLookupRequest:
		this.InternalRangeLookup = vt
	case *InternalHeartbeatTxnRequest:
		this.InternalHeartbeatTxn = vt
	case *InternalPushTxnRequest:
		this.InternalPushTxn = vt
	case *InternalResolveIntentRequest:
		this.InternalResolveIntent = vt
	case *InternalMergeRequest:
		this.InternalMergeResponse = vt
	case *InternalTruncateLogRequest:
		this.InternalTruncateLog = vt
	case *InternalGCRequest:
		this.InternalGc = vt
	case *InternalChangeReplicasRequest:
		this.InternalChangeReplicas = vt
	default:
		return false
	}
	return true
}
