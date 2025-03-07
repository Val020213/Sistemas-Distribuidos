// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v4.25.3
// source: chord.proto

package chordpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_chord_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{0}
}

type Node struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Address       string                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Node) Reset() {
	*x = Node{}
	mi := &file_chord_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{1}
}

func (x *Node) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Node) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type FindSuccessorRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           uint64                 `protobuf:"varint,1,opt,name=key,proto3" json:"key,omitempty"`
	Hops          int32                  `protobuf:"varint,2,opt,name=hops,proto3" json:"hops,omitempty"`
	Visited       map[uint64]bool        `protobuf:"bytes,3,rep,name=visited,proto3" json:"visited,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FindSuccessorRequest) Reset() {
	*x = FindSuccessorRequest{}
	mi := &file_chord_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FindSuccessorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSuccessorRequest) ProtoMessage() {}

func (x *FindSuccessorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSuccessorRequest.ProtoReflect.Descriptor instead.
func (*FindSuccessorRequest) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{2}
}

func (x *FindSuccessorRequest) GetKey() uint64 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *FindSuccessorRequest) GetHops() int32 {
	if x != nil {
		return x.Hops
	}
	return 0
}

func (x *FindSuccessorRequest) GetVisited() map[uint64]bool {
	if x != nil {
		return x.Visited
	}
	return nil
}

type GetSuccessorsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Successors    []*Node                `protobuf:"bytes,1,rep,name=successors,proto3" json:"successors,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSuccessorsResponse) Reset() {
	*x = GetSuccessorsResponse{}
	mi := &file_chord_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSuccessorsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSuccessorsResponse) ProtoMessage() {}

func (x *GetSuccessorsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSuccessorsResponse.ProtoReflect.Descriptor instead.
func (*GetSuccessorsResponse) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{3}
}

func (x *GetSuccessorsResponse) GetSuccessors() []*Node {
	if x != nil {
		return x.Successors
	}
	return nil
}

type Successful struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Successful    bool                   `protobuf:"varint,1,opt,name=successful,proto3" json:"successful,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Successful) Reset() {
	*x = Successful{}
	mi := &file_chord_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Successful) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Successful) ProtoMessage() {}

func (x *Successful) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Successful.ProtoReflect.Descriptor instead.
func (*Successful) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{4}
}

func (x *Successful) GetSuccessful() bool {
	if x != nil {
		return x.Successful
	}
	return false
}

type HealthResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Address       string                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	mi := &file_chord_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthResponse.ProtoReflect.Descriptor instead.
func (*HealthResponse) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{5}
}

func (x *HealthResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *HealthResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type StoreDataRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []*Data                `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StoreDataRequest) Reset() {
	*x = StoreDataRequest{}
	mi := &file_chord_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StoreDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreDataRequest) ProtoMessage() {}

func (x *StoreDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreDataRequest.ProtoReflect.Descriptor instead.
func (*StoreDataRequest) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{6}
}

func (x *StoreDataRequest) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type Data struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           uint64                 `protobuf:"varint,1,opt,name=key,proto3" json:"key,omitempty"`
	Url           string                 `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Status        string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Content       []byte                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Data) Reset() {
	*x = Data{}
	mi := &file_chord_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{7}
}

func (x *Data) GetKey() uint64 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *Data) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Data) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Data) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *Data) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Data) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type Id struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Id) Reset() {
	*x = Id{}
	mi := &file_chord_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Id) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Id) ProtoMessage() {}

func (x *Id) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Id.ProtoReflect.Descriptor instead.
func (*Id) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{8}
}

func (x *Id) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type State struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr          string                 `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Data          []*Data                `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	Finger        []*Node                `protobuf:"bytes,4,rep,name=finger,proto3" json:"finger,omitempty"`
	Successors    []*Node                `protobuf:"bytes,5,rep,name=successors,proto3" json:"successors,omitempty"`
	Predecessor   *Node                  `protobuf:"bytes,6,opt,name=predecessor,proto3" json:"predecessor,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *State) Reset() {
	*x = State{}
	mi := &file_chord_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{9}
}

func (x *State) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *State) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *State) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *State) GetFinger() []*Node {
	if x != nil {
		return x.Finger
	}
	return nil
}

func (x *State) GetSuccessors() []*Node {
	if x != nil {
		return x.Successors
	}
	return nil
}

func (x *State) GetPredecessor() *Node {
	if x != nil {
		return x.Predecessor
	}
	return nil
}

type CreateDataRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *Data                  `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateDataRequest) Reset() {
	*x = CreateDataRequest{}
	mi := &file_chord_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDataRequest) ProtoMessage() {}

func (x *CreateDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDataRequest.ProtoReflect.Descriptor instead.
func (*CreateDataRequest) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{10}
}

func (x *CreateDataRequest) GetData() *Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type GetNodeDataRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PredecesorId  uint64                 `protobuf:"varint,1,opt,name=predecesorId,proto3" json:"predecesorId,omitempty"`
	Id            uint64                 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetNodeDataRequest) Reset() {
	*x = GetNodeDataRequest{}
	mi := &file_chord_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetNodeDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNodeDataRequest) ProtoMessage() {}

func (x *GetNodeDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNodeDataRequest.ProtoReflect.Descriptor instead.
func (*GetNodeDataRequest) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{11}
}

func (x *GetNodeDataRequest) GetPredecesorId() uint64 {
	if x != nil {
		return x.PredecesorId
	}
	return 0
}

func (x *GetNodeDataRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Successors    []*Node                `protobuf:"bytes,1,rep,name=successors,proto3" json:"successors,omitempty"`
	Data          []*Data                `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	mi := &file_chord_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chord_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file_chord_proto_rawDescGZIP(), []int{12}
}

func (x *ListResponse) GetSuccessors() []*Node {
	if x != nil {
		return x.Successors
	}
	return nil
}

func (x *ListResponse) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_chord_proto protoreflect.FileDescriptor

var file_chord_proto_rawDesc = string([]byte{
	0x0a, 0x0b, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x63,
	0x68, 0x6f, 0x72, 0x64, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x30,
	0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x22, 0xbc, 0x01, 0x0a, 0x14, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x6f, 0x70, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x68, 0x6f, 0x70, 0x73, 0x12,
	0x42, 0x0a, 0x07, 0x76, 0x69, 0x73, 0x69, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x28, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x56, 0x69,
	0x73, 0x69, 0x74, 0x65, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x76, 0x69, 0x73, 0x69,
	0x74, 0x65, 0x64, 0x1a, 0x3a, 0x0a, 0x0c, 0x56, 0x69, 0x73, 0x69, 0x74, 0x65, 0x64, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x44, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x0a, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63,
	0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x6f, 0x72, 0x73, 0x22, 0x2c, 0x0a, 0x0a, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x66, 0x75, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x66, 0x75, 0x6c, 0x22, 0x3a, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22,
	0x33, 0x0a, 0x10, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0xd2, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39,
	0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x14, 0x0a, 0x02, 0x49, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22,
	0xcd, 0x01, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x1f, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x68,
	0x6f, 0x72, 0x64, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x23,
	0x0a, 0x06, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x06, 0x66, 0x69, 0x6e,
	0x67, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x52, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73,
	0x12, 0x2d, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x22,
	0x34, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x48, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x64, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x70,
	0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x6f, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x6f, 0x72, 0x49, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x5c, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2b, 0x0a, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x52, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x12, 0x1f, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x68, 0x6f,
	0x72, 0x64, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xec, 0x04,
	0x0a, 0x0c, 0x43, 0x68, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x28,
	0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64,
	0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x1a, 0x11, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x12, 0x2d, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x12, 0x0c, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x15, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0d, 0x46, 0x69, 0x6e, 0x64, 0x53,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x1b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64,
	0x2e, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x12, 0x2b, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65,
	0x73, 0x73, 0x6f, 0x72, 0x12, 0x0c, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x12,
	0x3b, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73,
	0x12, 0x0c, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1c,
	0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x09,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x17, 0x2e, 0x63, 0x68, 0x6f, 0x72,
	0x64, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x11, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x66, 0x75, 0x6c, 0x12, 0x2a, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x09, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x49, 0x64, 0x1a, 0x11,
	0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75,
	0x6c, 0x12, 0x28, 0x0a, 0x0a, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x0c, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e,
	0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x26, 0x0a, 0x0c, 0x52,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x09, 0x2e, 0x63, 0x68,
	0x6f, 0x72, 0x64, 0x2e, 0x49, 0x64, 0x1a, 0x0b, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x39, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x18, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x63, 0x68,
	0x6f, 0x72, 0x64, 0x2e, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x12, 0x41,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x19, 0x2e,
	0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64,
	0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x29, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x0c, 0x2e, 0x63, 0x68, 0x6f, 0x72,
	0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13, 0x2e, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0b, 0x5a, 0x09,
	0x2e, 0x2f, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_chord_proto_rawDescOnce sync.Once
	file_chord_proto_rawDescData []byte
)

func file_chord_proto_rawDescGZIP() []byte {
	file_chord_proto_rawDescOnce.Do(func() {
		file_chord_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_chord_proto_rawDesc), len(file_chord_proto_rawDesc)))
	})
	return file_chord_proto_rawDescData
}

var file_chord_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_chord_proto_goTypes = []any{
	(*Empty)(nil),                 // 0: chord.Empty
	(*Node)(nil),                  // 1: chord.Node
	(*FindSuccessorRequest)(nil),  // 2: chord.FindSuccessorRequest
	(*GetSuccessorsResponse)(nil), // 3: chord.GetSuccessorsResponse
	(*Successful)(nil),            // 4: chord.Successful
	(*HealthResponse)(nil),        // 5: chord.HealthResponse
	(*StoreDataRequest)(nil),      // 6: chord.StoreDataRequest
	(*Data)(nil),                  // 7: chord.Data
	(*Id)(nil),                    // 8: chord.Id
	(*State)(nil),                 // 9: chord.State
	(*CreateDataRequest)(nil),     // 10: chord.CreateDataRequest
	(*GetNodeDataRequest)(nil),    // 11: chord.GetNodeDataRequest
	(*ListResponse)(nil),          // 12: chord.ListResponse
	nil,                           // 13: chord.FindSuccessorRequest.VisitedEntry
	(*timestamppb.Timestamp)(nil), // 14: google.protobuf.Timestamp
}
var file_chord_proto_depIdxs = []int32{
	13, // 0: chord.FindSuccessorRequest.visited:type_name -> chord.FindSuccessorRequest.VisitedEntry
	1,  // 1: chord.GetSuccessorsResponse.successors:type_name -> chord.Node
	7,  // 2: chord.StoreDataRequest.data:type_name -> chord.Data
	14, // 3: chord.Data.created_at:type_name -> google.protobuf.Timestamp
	14, // 4: chord.Data.updated_at:type_name -> google.protobuf.Timestamp
	7,  // 5: chord.State.data:type_name -> chord.Data
	1,  // 6: chord.State.finger:type_name -> chord.Node
	1,  // 7: chord.State.successors:type_name -> chord.Node
	1,  // 8: chord.State.predecessor:type_name -> chord.Node
	7,  // 9: chord.CreateDataRequest.data:type_name -> chord.Data
	1,  // 10: chord.ListResponse.successors:type_name -> chord.Node
	7,  // 11: chord.ListResponse.data:type_name -> chord.Data
	1,  // 12: chord.ChordService.Notify:input_type -> chord.Node
	0,  // 13: chord.ChordService.Health:input_type -> chord.Empty
	2,  // 14: chord.ChordService.FindSuccessor:input_type -> chord.FindSuccessorRequest
	0,  // 15: chord.ChordService.GetPredecessor:input_type -> chord.Empty
	0,  // 16: chord.ChordService.GetSuccessors:input_type -> chord.Empty
	6,  // 17: chord.ChordService.StoreData:input_type -> chord.StoreDataRequest
	8,  // 18: chord.ChordService.DeleteData:input_type -> chord.Id
	0,  // 19: chord.ChordService.PrintState:input_type -> chord.Empty
	8,  // 20: chord.ChordService.RetrieveData:input_type -> chord.Id
	10, // 21: chord.ChordService.CreateData:input_type -> chord.CreateDataRequest
	11, // 22: chord.ChordService.GetNodeData:input_type -> chord.GetNodeDataRequest
	0,  // 23: chord.ChordService.List:input_type -> chord.Empty
	4,  // 24: chord.ChordService.Notify:output_type -> chord.Successful
	5,  // 25: chord.ChordService.Health:output_type -> chord.HealthResponse
	1,  // 26: chord.ChordService.FindSuccessor:output_type -> chord.Node
	1,  // 27: chord.ChordService.GetPredecessor:output_type -> chord.Node
	3,  // 28: chord.ChordService.GetSuccessors:output_type -> chord.GetSuccessorsResponse
	4,  // 29: chord.ChordService.StoreData:output_type -> chord.Successful
	4,  // 30: chord.ChordService.DeleteData:output_type -> chord.Successful
	9,  // 31: chord.ChordService.PrintState:output_type -> chord.State
	7,  // 32: chord.ChordService.RetrieveData:output_type -> chord.Data
	4,  // 33: chord.ChordService.CreateData:output_type -> chord.Successful
	6,  // 34: chord.ChordService.GetNodeData:output_type -> chord.StoreDataRequest
	12, // 35: chord.ChordService.List:output_type -> chord.ListResponse
	24, // [24:36] is the sub-list for method output_type
	12, // [12:24] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_chord_proto_init() }
func file_chord_proto_init() {
	if File_chord_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_chord_proto_rawDesc), len(file_chord_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chord_proto_goTypes,
		DependencyIndexes: file_chord_proto_depIdxs,
		MessageInfos:      file_chord_proto_msgTypes,
	}.Build()
	File_chord_proto = out.File
	file_chord_proto_goTypes = nil
	file_chord_proto_depIdxs = nil
}
