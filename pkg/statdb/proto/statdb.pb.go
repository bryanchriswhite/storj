// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: statdb.proto

package statdb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// CreateRequest is a request message for the Create rpc call
type CreateRequest struct {
	Node                 *Node      `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"`
	Stats                *NodeStats `protobuf:"bytes,2,opt,name=stats" json:"stats,omitempty"`
	APIKey               []byte     `protobuf:"bytes,3,opt,name=APIKey,proto3" json:"APIKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{0}
}
func (m *CreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRequest.Unmarshal(m, b)
}
func (m *CreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRequest.Marshal(b, m, deterministic)
}
func (dst *CreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRequest.Merge(dst, src)
}
func (m *CreateRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRequest.Size(m)
}
func (m *CreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRequest proto.InternalMessageInfo

func (m *CreateRequest) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *CreateRequest) GetStats() *NodeStats {
	if m != nil {
		return m.Stats
	}
	return nil
}

func (m *CreateRequest) GetAPIKey() []byte {
	if m != nil {
		return m.APIKey
	}
	return nil
}

// CreateResponse is a response message for the Create rpc call
type CreateResponse struct {
	Stats                *NodeStats `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *CreateResponse) Reset()         { *m = CreateResponse{} }
func (m *CreateResponse) String() string { return proto.CompactTextString(m) }
func (*CreateResponse) ProtoMessage()    {}
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{1}
}
func (m *CreateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateResponse.Unmarshal(m, b)
}
func (m *CreateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateResponse.Marshal(b, m, deterministic)
}
func (dst *CreateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateResponse.Merge(dst, src)
}
func (m *CreateResponse) XXX_Size() int {
	return xxx_messageInfo_CreateResponse.Size(m)
}
func (m *CreateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateResponse proto.InternalMessageInfo

func (m *CreateResponse) GetStats() *NodeStats {
	if m != nil {
		return m.Stats
	}
	return nil
}

// GetRequest is a request message for the Get rpc call
type GetRequest struct {
	NodeId               NodeID   `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3,customtype=NodeID" json:"node_id"`
	APIKey               []byte   `protobuf:"bytes,2,opt,name=APIKey,proto3" json:"APIKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{2}
}
func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (dst *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(dst, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetAPIKey() []byte {
	if m != nil {
		return m.APIKey
	}
	return nil
}

// GetResponse is a response message for the Get rpc call
type GetResponse struct {
	Stats                *NodeStats `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{3}
}
func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (dst *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(dst, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

func (m *GetResponse) GetStats() *NodeStats {
	if m != nil {
		return m.Stats
	}
	return nil
}

// FindValidNodesRequest is a request message for the FindValidNodes rpc call
type FindValidNodesRequest struct {
	NodeIds              NodeIDList `protobuf:"bytes,1,opt,name=node_ids,json=nodeIds,proto3,casttype=NodeIDList" json:"node_ids,omitempty"`
	MinStats             *NodeStats `protobuf:"bytes,2,opt,name=min_stats,json=minStats" json:"min_stats,omitempty"`
	APIKey               []byte     `protobuf:"bytes,3,opt,name=APIKey,proto3" json:"APIKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *FindValidNodesRequest) Reset()         { *m = FindValidNodesRequest{} }
func (m *FindValidNodesRequest) String() string { return proto.CompactTextString(m) }
func (*FindValidNodesRequest) ProtoMessage()    {}
func (*FindValidNodesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{4}
}
func (m *FindValidNodesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindValidNodesRequest.Unmarshal(m, b)
}
func (m *FindValidNodesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindValidNodesRequest.Marshal(b, m, deterministic)
}
func (dst *FindValidNodesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindValidNodesRequest.Merge(dst, src)
}
func (m *FindValidNodesRequest) XXX_Size() int {
	return xxx_messageInfo_FindValidNodesRequest.Size(m)
}
func (m *FindValidNodesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FindValidNodesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FindValidNodesRequest proto.InternalMessageInfo

func (m *FindValidNodesRequest) GetNodeIds() NodeIDList {
	if m != nil {
		return m.NodeIds
	}
	return nil
}

func (m *FindValidNodesRequest) GetMinStats() *NodeStats {
	if m != nil {
		return m.MinStats
	}
	return nil
}

func (m *FindValidNodesRequest) GetAPIKey() []byte {
	if m != nil {
		return m.APIKey
	}
	return nil
}

// FindValidNodesResponse is a response message for the FindValidNodes rpc call
type FindValidNodesResponse struct {
	PassedIds            NodeIDList `protobuf:"bytes,1,opt,name=passed_ids,json=passedIds,proto3,casttype=NodeIDList" json:"passed_ids,omitempty"`
	FailedIds            NodeIDList `protobuf:"bytes,2,opt,name=failed_ids,json=failedIds,proto3,casttype=NodeIDList" json:"failed_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *FindValidNodesResponse) Reset()         { *m = FindValidNodesResponse{} }
func (m *FindValidNodesResponse) String() string { return proto.CompactTextString(m) }
func (*FindValidNodesResponse) ProtoMessage()    {}
func (*FindValidNodesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{5}
}
func (m *FindValidNodesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindValidNodesResponse.Unmarshal(m, b)
}
func (m *FindValidNodesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindValidNodesResponse.Marshal(b, m, deterministic)
}
func (dst *FindValidNodesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindValidNodesResponse.Merge(dst, src)
}
func (m *FindValidNodesResponse) XXX_Size() int {
	return xxx_messageInfo_FindValidNodesResponse.Size(m)
}
func (m *FindValidNodesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FindValidNodesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FindValidNodesResponse proto.InternalMessageInfo

func (m *FindValidNodesResponse) GetPassedIds() NodeIDList {
	if m != nil {
		return m.PassedIds
	}
	return nil
}

func (m *FindValidNodesResponse) GetFailedIds() NodeIDList {
	if m != nil {
		return m.FailedIds
	}
	return nil
}

// UpdateRequest is a request message for the Update rpc call
type UpdateRequest struct {
	Node                 *Node    `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"`
	APIKey               []byte   `protobuf:"bytes,2,opt,name=APIKey,proto3" json:"APIKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateRequest) Reset()         { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()    {}
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{6}
}
func (m *UpdateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateRequest.Unmarshal(m, b)
}
func (m *UpdateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateRequest.Marshal(b, m, deterministic)
}
func (dst *UpdateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateRequest.Merge(dst, src)
}
func (m *UpdateRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateRequest.Size(m)
}
func (m *UpdateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateRequest proto.InternalMessageInfo

func (m *UpdateRequest) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *UpdateRequest) GetAPIKey() []byte {
	if m != nil {
		return m.APIKey
	}
	return nil
}

// UpdateRequest is a response message for the Update rpc call
type UpdateResponse struct {
	Stats                *NodeStats `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *UpdateResponse) Reset()         { *m = UpdateResponse{} }
func (m *UpdateResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateResponse) ProtoMessage()    {}
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{7}
}
func (m *UpdateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateResponse.Unmarshal(m, b)
}
func (m *UpdateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateResponse.Marshal(b, m, deterministic)
}
func (dst *UpdateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateResponse.Merge(dst, src)
}
func (m *UpdateResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateResponse.Size(m)
}
func (m *UpdateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateResponse proto.InternalMessageInfo

func (m *UpdateResponse) GetStats() *NodeStats {
	if m != nil {
		return m.Stats
	}
	return nil
}

// UpdateBatchRequest is a request message for the UpdateBatch rpc call
type UpdateBatchRequest struct {
	NodeList             []*Node  `protobuf:"bytes,1,rep,name=node_list,json=nodeList" json:"node_list,omitempty"`
	APIKey               []byte   `protobuf:"bytes,2,opt,name=APIKey,proto3" json:"APIKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateBatchRequest) Reset()         { *m = UpdateBatchRequest{} }
func (m *UpdateBatchRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateBatchRequest) ProtoMessage()    {}
func (*UpdateBatchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{8}
}
func (m *UpdateBatchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateBatchRequest.Unmarshal(m, b)
}
func (m *UpdateBatchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateBatchRequest.Marshal(b, m, deterministic)
}
func (dst *UpdateBatchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateBatchRequest.Merge(dst, src)
}
func (m *UpdateBatchRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateBatchRequest.Size(m)
}
func (m *UpdateBatchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateBatchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateBatchRequest proto.InternalMessageInfo

func (m *UpdateBatchRequest) GetNodeList() []*Node {
	if m != nil {
		return m.NodeList
	}
	return nil
}

func (m *UpdateBatchRequest) GetAPIKey() []byte {
	if m != nil {
		return m.APIKey
	}
	return nil
}

// UpdateBatchResponse is a response message for the UpdateBatch rpc call
type UpdateBatchResponse struct {
	StatsList            []*NodeStats `protobuf:"bytes,1,rep,name=stats_list,json=statsList" json:"stats_list,omitempty"`
	FailedNodes          []*Node      `protobuf:"bytes,2,rep,name=failed_nodes,json=failedNodes" json:"failed_nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *UpdateBatchResponse) Reset()         { *m = UpdateBatchResponse{} }
func (m *UpdateBatchResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateBatchResponse) ProtoMessage()    {}
func (*UpdateBatchResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{9}
}
func (m *UpdateBatchResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateBatchResponse.Unmarshal(m, b)
}
func (m *UpdateBatchResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateBatchResponse.Marshal(b, m, deterministic)
}
func (dst *UpdateBatchResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateBatchResponse.Merge(dst, src)
}
func (m *UpdateBatchResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateBatchResponse.Size(m)
}
func (m *UpdateBatchResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateBatchResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateBatchResponse proto.InternalMessageInfo

func (m *UpdateBatchResponse) GetStatsList() []*NodeStats {
	if m != nil {
		return m.StatsList
	}
	return nil
}

func (m *UpdateBatchResponse) GetFailedNodes() []*Node {
	if m != nil {
		return m.FailedNodes
	}
	return nil
}

// CreateEntryIfNotExistsRequest is a request message for the CreateEntryIfNotExists rpc call
type CreateEntryIfNotExistsRequest struct {
	Node                 *Node    `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"`
	APIKey               []byte   `protobuf:"bytes,2,opt,name=APIKey,proto3" json:"APIKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateEntryIfNotExistsRequest) Reset()         { *m = CreateEntryIfNotExistsRequest{} }
func (m *CreateEntryIfNotExistsRequest) String() string { return proto.CompactTextString(m) }
func (*CreateEntryIfNotExistsRequest) ProtoMessage()    {}
func (*CreateEntryIfNotExistsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{10}
}
func (m *CreateEntryIfNotExistsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateEntryIfNotExistsRequest.Unmarshal(m, b)
}
func (m *CreateEntryIfNotExistsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateEntryIfNotExistsRequest.Marshal(b, m, deterministic)
}
func (dst *CreateEntryIfNotExistsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateEntryIfNotExistsRequest.Merge(dst, src)
}
func (m *CreateEntryIfNotExistsRequest) XXX_Size() int {
	return xxx_messageInfo_CreateEntryIfNotExistsRequest.Size(m)
}
func (m *CreateEntryIfNotExistsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateEntryIfNotExistsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateEntryIfNotExistsRequest proto.InternalMessageInfo

func (m *CreateEntryIfNotExistsRequest) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *CreateEntryIfNotExistsRequest) GetAPIKey() []byte {
	if m != nil {
		return m.APIKey
	}
	return nil
}

// CreateEntryIfNotExistsResponse is a response message for the CreateEntryIfNotExists rpc call
type CreateEntryIfNotExistsResponse struct {
	Stats                *NodeStats `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *CreateEntryIfNotExistsResponse) Reset()         { *m = CreateEntryIfNotExistsResponse{} }
func (m *CreateEntryIfNotExistsResponse) String() string { return proto.CompactTextString(m) }
func (*CreateEntryIfNotExistsResponse) ProtoMessage()    {}
func (*CreateEntryIfNotExistsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_statdb_5594d60637806120, []int{11}
}
func (m *CreateEntryIfNotExistsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateEntryIfNotExistsResponse.Unmarshal(m, b)
}
func (m *CreateEntryIfNotExistsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateEntryIfNotExistsResponse.Marshal(b, m, deterministic)
}
func (dst *CreateEntryIfNotExistsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateEntryIfNotExistsResponse.Merge(dst, src)
}
func (m *CreateEntryIfNotExistsResponse) XXX_Size() int {
	return xxx_messageInfo_CreateEntryIfNotExistsResponse.Size(m)
}
func (m *CreateEntryIfNotExistsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateEntryIfNotExistsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateEntryIfNotExistsResponse proto.InternalMessageInfo

func (m *CreateEntryIfNotExistsResponse) GetStats() *NodeStats {
	if m != nil {
		return m.Stats
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateRequest)(nil), "statdb.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "statdb.CreateResponse")
	proto.RegisterType((*GetRequest)(nil), "statdb.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "statdb.GetResponse")
	proto.RegisterType((*FindValidNodesRequest)(nil), "statdb.FindValidNodesRequest")
	proto.RegisterType((*FindValidNodesResponse)(nil), "statdb.FindValidNodesResponse")
	proto.RegisterType((*UpdateRequest)(nil), "statdb.UpdateRequest")
	proto.RegisterType((*UpdateResponse)(nil), "statdb.UpdateResponse")
	proto.RegisterType((*UpdateBatchRequest)(nil), "statdb.UpdateBatchRequest")
	proto.RegisterType((*UpdateBatchResponse)(nil), "statdb.UpdateBatchResponse")
	proto.RegisterType((*CreateEntryIfNotExistsRequest)(nil), "statdb.CreateEntryIfNotExistsRequest")
	proto.RegisterType((*CreateEntryIfNotExistsResponse)(nil), "statdb.CreateEntryIfNotExistsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for StatDB service

type StatDBClient interface {
	// Create a db entry for the provided storagenode ID
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	// Get uses a storagenode ID to get that storagenode's stats
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	// FindValidNodes gets a subset of storagenodes that fit minimum reputation args
	FindValidNodes(ctx context.Context, in *FindValidNodesRequest, opts ...grpc.CallOption) (*FindValidNodesResponse, error)
	// Update updates storagenode stats for a single storagenode
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	// UpdateBatch updates storagenode stats for multiple farmers at a time
	UpdateBatch(ctx context.Context, in *UpdateBatchRequest, opts ...grpc.CallOption) (*UpdateBatchResponse, error)
	// CreateEntryIfNotExists creates a db entry if it didn't exist
	CreateEntryIfNotExists(ctx context.Context, in *CreateEntryIfNotExistsRequest, opts ...grpc.CallOption) (*CreateEntryIfNotExistsResponse, error)
}

type statDBClient struct {
	cc *grpc.ClientConn
}

func NewStatDBClient(cc *grpc.ClientConn) StatDBClient {
	return &statDBClient{cc}
}

func (c *statDBClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/statdb.StatDB/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statDBClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/statdb.StatDB/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statDBClient) FindValidNodes(ctx context.Context, in *FindValidNodesRequest, opts ...grpc.CallOption) (*FindValidNodesResponse, error) {
	out := new(FindValidNodesResponse)
	err := c.cc.Invoke(ctx, "/statdb.StatDB/FindValidNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statDBClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/statdb.StatDB/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statDBClient) UpdateBatch(ctx context.Context, in *UpdateBatchRequest, opts ...grpc.CallOption) (*UpdateBatchResponse, error) {
	out := new(UpdateBatchResponse)
	err := c.cc.Invoke(ctx, "/statdb.StatDB/UpdateBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statDBClient) CreateEntryIfNotExists(ctx context.Context, in *CreateEntryIfNotExistsRequest, opts ...grpc.CallOption) (*CreateEntryIfNotExistsResponse, error) {
	out := new(CreateEntryIfNotExistsResponse)
	err := c.cc.Invoke(ctx, "/statdb.StatDB/CreateEntryIfNotExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for StatDB service

type StatDBServer interface {
	// Create a db entry for the provided storagenode ID
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// Get uses a storagenode ID to get that storagenode's stats
	Get(context.Context, *GetRequest) (*GetResponse, error)
	// FindValidNodes gets a subset of storagenodes that fit minimum reputation args
	FindValidNodes(context.Context, *FindValidNodesRequest) (*FindValidNodesResponse, error)
	// Update updates storagenode stats for a single storagenode
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	// UpdateBatch updates storagenode stats for multiple farmers at a time
	UpdateBatch(context.Context, *UpdateBatchRequest) (*UpdateBatchResponse, error)
	// CreateEntryIfNotExists creates a db entry if it didn't exist
	CreateEntryIfNotExists(context.Context, *CreateEntryIfNotExistsRequest) (*CreateEntryIfNotExistsResponse, error)
}

func RegisterStatDBServer(s *grpc.Server, srv StatDBServer) {
	s.RegisterService(&_StatDB_serviceDesc, srv)
}

func _StatDB_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatDBServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statdb.StatDB/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatDBServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatDB_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatDBServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statdb.StatDB/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatDBServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatDB_FindValidNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindValidNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatDBServer).FindValidNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statdb.StatDB/FindValidNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatDBServer).FindValidNodes(ctx, req.(*FindValidNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatDB_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatDBServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statdb.StatDB/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatDBServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatDB_UpdateBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatDBServer).UpdateBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statdb.StatDB/UpdateBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatDBServer).UpdateBatch(ctx, req.(*UpdateBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatDB_CreateEntryIfNotExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEntryIfNotExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatDBServer).CreateEntryIfNotExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statdb.StatDB/CreateEntryIfNotExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatDBServer).CreateEntryIfNotExists(ctx, req.(*CreateEntryIfNotExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StatDB_serviceDesc = grpc.ServiceDesc{
	ServiceName: "statdb.StatDB",
	HandlerType: (*StatDBServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _StatDB_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _StatDB_Get_Handler,
		},
		{
			MethodName: "FindValidNodes",
			Handler:    _StatDB_FindValidNodes_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _StatDB_Update_Handler,
		},
		{
			MethodName: "UpdateBatch",
			Handler:    _StatDB_UpdateBatch_Handler,
		},
		{
			MethodName: "CreateEntryIfNotExists",
			Handler:    _StatDB_CreateEntryIfNotExists_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statdb.proto",
}

func init() { proto.RegisterFile("statdb.proto", fileDescriptor_statdb_5594d60637806120) }

var fileDescriptor_statdb_5594d60637806120 = []byte{
	// 531 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xdf, 0x6e, 0x12, 0x41,
	0x14, 0xc6, 0x5d, 0xd0, 0x6d, 0xf9, 0xa0, 0x98, 0x4c, 0x53, 0x42, 0xd6, 0x14, 0x08, 0x49, 0xb5,
	0x26, 0xc2, 0x45, 0x35, 0xe1, 0x5a, 0x6c, 0x4b, 0x36, 0x6a, 0x35, 0xdb, 0x54, 0x2f, 0x9b, 0xad,
	0x33, 0xc5, 0x49, 0xe8, 0x2e, 0x32, 0xa3, 0xb1, 0x6f, 0xe0, 0xc3, 0xf8, 0x20, 0x3e, 0x83, 0x17,
	0x3c, 0x88, 0x57, 0x66, 0xfe, 0xac, 0xfb, 0x47, 0x16, 0x43, 0xbc, 0xdb, 0x3d, 0xe7, 0x9b, 0x6f,
	0x7e, 0x73, 0xce, 0x99, 0x41, 0x43, 0xc8, 0x50, 0xd2, 0xab, 0xe1, 0x7c, 0x11, 0xcb, 0x98, 0xb8,
	0xe6, 0xcf, 0xc3, 0x34, 0x9e, 0xc6, 0x26, 0xe6, 0x21, 0x8a, 0x29, 0x33, 0xdf, 0xfd, 0x08, 0x3b,
	0x2f, 0x16, 0x2c, 0x94, 0x2c, 0x60, 0x9f, 0x3e, 0x33, 0x21, 0x49, 0x07, 0x77, 0x55, 0xba, 0xed,
	0xf4, 0x9c, 0xc3, 0xfa, 0x11, 0x86, 0x5a, 0x7b, 0x16, 0x53, 0x16, 0xe8, 0x38, 0x39, 0xc0, 0x3d,
	0x65, 0x29, 0xda, 0x15, 0x2d, 0xb8, 0x9f, 0x0a, 0xce, 0x55, 0x38, 0x30, 0x59, 0xd2, 0x82, 0xfb,
	0xfc, 0xad, 0xff, 0x92, 0xdd, 0xb6, 0xab, 0x3d, 0xe7, 0xb0, 0x11, 0xd8, 0xbf, 0xfe, 0x08, 0xcd,
	0x64, 0x3f, 0x31, 0x8f, 0x23, 0x91, 0x31, 0x74, 0xd6, 0x19, 0xf6, 0x5f, 0x03, 0x13, 0x26, 0x13,
	0xca, 0x47, 0xd8, 0x52, 0xb2, 0x4b, 0x4e, 0xf5, 0xb2, 0xc6, 0xb8, 0xf9, 0x63, 0xd9, 0xbd, 0xf3,
	0x73, 0xd9, 0x75, 0xd5, 0x42, 0xff, 0x38, 0x70, 0x55, 0xda, 0xa7, 0x19, 0x8e, 0x4a, 0x8e, 0xe3,
	0x19, 0xea, 0xda, 0x6e, 0x33, 0x88, 0x6f, 0x0e, 0xf6, 0x4e, 0x79, 0x44, 0xdf, 0x85, 0x33, 0x4e,
	0x55, 0x56, 0x24, 0x40, 0x8f, 0xb1, 0x6d, 0x81, 0x44, 0x42, 0xf4, 0x6b, 0xd9, 0x85, 0xa1, 0x79,
	0xc5, 0x85, 0x0c, 0xb6, 0x0c, 0x91, 0x20, 0x4f, 0x50, 0xbb, 0xe1, 0xd1, 0xe5, 0xda, 0x2a, 0x6e,
	0xdf, 0xf0, 0xe8, 0x7c, 0x6d, 0x21, 0xbf, 0xa0, 0x55, 0x24, 0xb1, 0x67, 0x19, 0x00, 0xf3, 0x50,
	0x08, 0x46, 0xd7, 0xc0, 0xd4, 0x8c, 0x42, 0xe1, 0x0c, 0x80, 0xeb, 0x90, 0xcf, 0xac, 0xbc, 0xb2,
	0x5a, 0x6e, 0x14, 0x3e, 0x15, 0xfd, 0x09, 0x76, 0x2e, 0xe6, 0x74, 0x83, 0x81, 0x29, 0xeb, 0xc0,
	0x08, 0xcd, 0xc4, 0x68, 0xb3, 0x26, 0x5c, 0x80, 0x98, 0x85, 0xe3, 0x50, 0x7e, 0xf8, 0x98, 0x4e,
	0x44, 0x4d, 0x37, 0x60, 0xc6, 0x85, 0x6c, 0x3b, 0xbd, 0x6a, 0x81, 0x45, 0x77, 0x47, 0x9d, 0xa5,
	0x94, 0x47, 0x62, 0x37, 0x67, 0x6b, 0xa1, 0x86, 0x80, 0xde, 0x36, 0x6b, 0xfc, 0x17, 0x59, 0x4d,
	0x4b, 0xb4, 0xfd, 0x00, 0x0d, 0x5b, 0x4e, 0xa5, 0x51, 0x05, 0x2d, 0xa2, 0xd4, 0x4d, 0x5e, 0x37,
	0xad, 0xff, 0x1e, 0xfb, 0xe6, 0x3e, 0x9c, 0x44, 0x72, 0x71, 0xeb, 0x5f, 0x9f, 0xc5, 0xf2, 0xe4,
	0x2b, 0x17, 0x52, 0xfc, 0x6f, 0x79, 0x27, 0xe8, 0x94, 0x19, 0x6f, 0x54, 0xee, 0xa3, 0xef, 0x55,
	0xb8, 0x2a, 0x70, 0x3c, 0x26, 0x23, 0xb8, 0xc6, 0x93, 0xec, 0x0d, 0xed, 0x2b, 0x93, 0x7b, 0x3c,
	0xbc, 0x56, 0x31, 0xfc, 0xa7, 0x88, 0xd5, 0x09, 0x93, 0x84, 0x24, 0xe9, 0xf4, 0x26, 0x7b, 0xbb,
	0xb9, 0x98, 0xd5, 0xbf, 0x41, 0x33, 0x3f, 0xdc, 0x64, 0x3f, 0x91, 0xad, 0xbc, 0x7e, 0x5e, 0xa7,
	0x2c, 0x6d, 0x0d, 0x47, 0x70, 0x4d, 0x73, 0x53, 0xf2, 0xdc, 0x14, 0xa7, 0xe4, 0x85, 0x99, 0x3c,
	0x45, 0x3d, 0x33, 0x15, 0xc4, 0xcb, 0xcb, 0xb2, 0x13, 0xe8, 0x3d, 0x58, 0x99, 0xb3, 0x3e, 0x53,
	0xb4, 0x56, 0xb7, 0x83, 0x1c, 0xe4, 0x6b, 0x56, 0x32, 0x07, 0xde, 0xc3, 0x7f, 0xc9, 0xcc, 0x46,
	0x57, 0xae, 0x7e, 0xd7, 0x9f, 0xfe, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xc5, 0x29, 0xe3, 0x04, 0x07,
	0x06, 0x00, 0x00,
}
