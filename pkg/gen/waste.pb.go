// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: waste.proto

package gen

import (
	types "github.com/smart-core-os/sc-api/go/types"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Describes the disposal method of the waste. Recycled, general
type WasteRecord_DisposalMethod int32

const (
	WasteRecord_DISPOSAL_METHOD_UNSPECIFIED WasteRecord_DisposalMethod = 0
	WasteRecord_GENERAL_WASTE               WasteRecord_DisposalMethod = 1
	WasteRecord_MIXED_RECYCLING             WasteRecord_DisposalMethod = 2
)

// Enum value maps for WasteRecord_DisposalMethod.
var (
	WasteRecord_DisposalMethod_name = map[int32]string{
		0: "DISPOSAL_METHOD_UNSPECIFIED",
		1: "GENERAL_WASTE",
		2: "MIXED_RECYCLING",
	}
	WasteRecord_DisposalMethod_value = map[string]int32{
		"DISPOSAL_METHOD_UNSPECIFIED": 0,
		"GENERAL_WASTE":               1,
		"MIXED_RECYCLING":             2,
	}
)

func (x WasteRecord_DisposalMethod) Enum() *WasteRecord_DisposalMethod {
	p := new(WasteRecord_DisposalMethod)
	*p = x
	return p
}

func (x WasteRecord_DisposalMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WasteRecord_DisposalMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_waste_proto_enumTypes[0].Descriptor()
}

func (WasteRecord_DisposalMethod) Type() protoreflect.EnumType {
	return &file_waste_proto_enumTypes[0]
}

func (x WasteRecord_DisposalMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WasteRecord_DisposalMethod.Descriptor instead.
func (WasteRecord_DisposalMethod) EnumDescriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{0, 0}
}

// WasteRecord is a record of a unit of waste produced by the building
type WasteRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The id of the waste record assigned by the external waste management system
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The time the record was created, this is not the time the waste was created
	RecordCreateTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=record_create_time,json=recordCreateTime,proto3" json:"record_create_time,omitempty"`
	// The weight of the waste. Use the InfoApi to find the actual unit the device uses.
	Weight float32 `protobuf:"fixed32,3,opt,name=weight,proto3" json:"weight,omitempty"`
	// The system the waste was collected from. Used for presentation not to identify the waste
	System         string                     `protobuf:"bytes,4,opt,name=system,proto3" json:"system,omitempty"`
	DisposalMethod WasteRecord_DisposalMethod `protobuf:"varint,5,opt,name=disposal_method,json=disposalMethod,proto3,enum=smartcore.bos.WasteRecord_DisposalMethod" json:"disposal_method,omitempty"`
	// Area the waste was collected from, e.g. tenant x, common area. Used for presentation not to identify the waste
	Area string `protobuf:"bytes,6,opt,name=area,proto3" json:"area,omitempty"`
	// The date/time the waste was created
	WasteCreateTime *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=waste_create_time,json=wasteCreateTime,proto3" json:"waste_create_time,omitempty"`
	// The stream the waste was collected from. For example, general waste, mixed recycling etc.
	// Used for presentation not to identify the waste
	Stream string `protobuf:"bytes,8,opt,name=stream,proto3" json:"stream,omitempty"`
}

func (x *WasteRecord) Reset() {
	*x = WasteRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waste_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WasteRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WasteRecord) ProtoMessage() {}

func (x *WasteRecord) ProtoReflect() protoreflect.Message {
	mi := &file_waste_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WasteRecord.ProtoReflect.Descriptor instead.
func (*WasteRecord) Descriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{0}
}

func (x *WasteRecord) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *WasteRecord) GetRecordCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.RecordCreateTime
	}
	return nil
}

func (x *WasteRecord) GetWeight() float32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *WasteRecord) GetSystem() string {
	if x != nil {
		return x.System
	}
	return ""
}

func (x *WasteRecord) GetDisposalMethod() WasteRecord_DisposalMethod {
	if x != nil {
		return x.DisposalMethod
	}
	return WasteRecord_DISPOSAL_METHOD_UNSPECIFIED
}

func (x *WasteRecord) GetArea() string {
	if x != nil {
		return x.Area
	}
	return ""
}

func (x *WasteRecord) GetWasteCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.WasteCreateTime
	}
	return nil
}

func (x *WasteRecord) GetStream() string {
	if x != nil {
		return x.Stream
	}
	return ""
}

type ListWasteRecordsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WasteRecords []*WasteRecord `protobuf:"bytes,1,rep,name=wasteRecords,proto3" json:"wasteRecords,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page.
	// If this field is omitted, there are no subsequent pages.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	// If non-zero this is the total number of alerts after filtering is applied.
	// This may be an estimate.
	TotalSize int32 `protobuf:"varint,3,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`
}

func (x *ListWasteRecordsResponse) Reset() {
	*x = ListWasteRecordsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waste_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWasteRecordsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWasteRecordsResponse) ProtoMessage() {}

func (x *ListWasteRecordsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_waste_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWasteRecordsResponse.ProtoReflect.Descriptor instead.
func (*ListWasteRecordsResponse) Descriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{1}
}

func (x *ListWasteRecordsResponse) GetWasteRecords() []*WasteRecord {
	if x != nil {
		return x.WasteRecords
	}
	return nil
}

func (x *ListWasteRecordsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

func (x *ListWasteRecordsResponse) GetTotalSize() int32 {
	if x != nil {
		return x.TotalSize
	}
	return 0
}

type ListWasteRecordsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	// The maximum number of WasteRecords to return.
	// The service may return fewer than this value.
	// If unspecified, at most 50 items will be returned.
	// The maximum value is 1000; values above 1000 will be coerced to 1000.
	PageSize int32 `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// A page token, received from a previous `ListWasteRecordsResponse` call.
	// Provide this to retrieve the subsequent page.
	PageToken string `protobuf:"bytes,4,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *ListWasteRecordsRequest) Reset() {
	*x = ListWasteRecordsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waste_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWasteRecordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWasteRecordsRequest) ProtoMessage() {}

func (x *ListWasteRecordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_waste_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWasteRecordsRequest.ProtoReflect.Descriptor instead.
func (*ListWasteRecordsRequest) Descriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{2}
}

func (x *ListWasteRecordsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListWasteRecordsRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

func (x *ListWasteRecordsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListWasteRecordsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type PullWasteRecordsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask    *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	UpdatesOnly bool                   `protobuf:"varint,3,opt,name=updates_only,json=updatesOnly,proto3" json:"updates_only,omitempty"`
}

func (x *PullWasteRecordsRequest) Reset() {
	*x = PullWasteRecordsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waste_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullWasteRecordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullWasteRecordsRequest) ProtoMessage() {}

func (x *PullWasteRecordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_waste_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullWasteRecordsRequest.ProtoReflect.Descriptor instead.
func (*PullWasteRecordsRequest) Descriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{3}
}

func (x *PullWasteRecordsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullWasteRecordsRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

func (x *PullWasteRecordsRequest) GetUpdatesOnly() bool {
	if x != nil {
		return x.UpdatesOnly
	}
	return false
}

type PullWasteRecordsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Changes []*PullWasteRecordsResponse_Change `protobuf:"bytes,1,rep,name=changes,proto3" json:"changes,omitempty"`
}

func (x *PullWasteRecordsResponse) Reset() {
	*x = PullWasteRecordsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waste_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullWasteRecordsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullWasteRecordsResponse) ProtoMessage() {}

func (x *PullWasteRecordsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_waste_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullWasteRecordsResponse.ProtoReflect.Descriptor instead.
func (*PullWasteRecordsResponse) Descriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{4}
}

func (x *PullWasteRecordsResponse) GetChanges() []*PullWasteRecordsResponse_Change {
	if x != nil {
		return x.Changes
	}
	return nil
}

type DescribeWasteRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DescribeWasteRecordRequest) Reset() {
	*x = DescribeWasteRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waste_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeWasteRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeWasteRecordRequest) ProtoMessage() {}

func (x *DescribeWasteRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_waste_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeWasteRecordRequest.ProtoReflect.Descriptor instead.
func (*DescribeWasteRecordRequest) Descriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{5}
}

func (x *DescribeWasteRecordRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// WasteRecordSupport describes the capabilities of devices implementing this trait
type WasteRecordSupport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// How a named device supports read/write/pull apis
	ResourceSupport *types.ResourceSupport `protobuf:"bytes,1,opt,name=resource_support,json=resourceSupport,proto3" json:"resource_support,omitempty"`
	// The unit associated with the weight value
	Unit string `protobuf:"bytes,2,opt,name=unit,proto3" json:"unit,omitempty"`
}

func (x *WasteRecordSupport) Reset() {
	*x = WasteRecordSupport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waste_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WasteRecordSupport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WasteRecordSupport) ProtoMessage() {}

func (x *WasteRecordSupport) ProtoReflect() protoreflect.Message {
	mi := &file_waste_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WasteRecordSupport.ProtoReflect.Descriptor instead.
func (*WasteRecordSupport) Descriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{6}
}

func (x *WasteRecordSupport) GetResourceSupport() *types.ResourceSupport {
	if x != nil {
		return x.ResourceSupport
	}
	return nil
}

func (x *WasteRecordSupport) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

type PullWasteRecordsResponse_Change struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ChangeTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=change_time,json=changeTime,proto3" json:"change_time,omitempty"`
	NewValue   *WasteRecord           `protobuf:"bytes,3,opt,name=new_value,json=newValue,proto3" json:"new_value,omitempty"`
	OldValue   *WasteRecord           `protobuf:"bytes,4,opt,name=old_value,json=oldValue,proto3" json:"old_value,omitempty"`
	Type       types.ChangeType       `protobuf:"varint,5,opt,name=type,proto3,enum=smartcore.types.ChangeType" json:"type,omitempty"`
}

func (x *PullWasteRecordsResponse_Change) Reset() {
	*x = PullWasteRecordsResponse_Change{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waste_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullWasteRecordsResponse_Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullWasteRecordsResponse_Change) ProtoMessage() {}

func (x *PullWasteRecordsResponse_Change) ProtoReflect() protoreflect.Message {
	mi := &file_waste_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullWasteRecordsResponse_Change.ProtoReflect.Descriptor instead.
func (*PullWasteRecordsResponse_Change) Descriptor() ([]byte, []int) {
	return file_waste_proto_rawDescGZIP(), []int{4, 0}
}

func (x *PullWasteRecordsResponse_Change) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullWasteRecordsResponse_Change) GetChangeTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ChangeTime
	}
	return nil
}

func (x *PullWasteRecordsResponse_Change) GetNewValue() *WasteRecord {
	if x != nil {
		return x.NewValue
	}
	return nil
}

func (x *PullWasteRecordsResponse_Change) GetOldValue() *WasteRecord {
	if x != nil {
		return x.OldValue
	}
	return nil
}

func (x *PullWasteRecordsResponse_Change) GetType() types.ChangeType {
	if x != nil {
		return x.Type
	}
	return types.ChangeType(0)
}

var File_waste_proto protoreflect.FileDescriptor

var file_waste_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x77, 0x61, 0x73, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73,
	0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x1a, 0x20, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x10, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x12, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xba, 0x03, 0x0a, 0x0b, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x48, 0x0a, 0x12, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x10, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12,
	0x52, 0x0a, 0x0f, 0x64, 0x69, 0x73, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x5f, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x52, 0x0e, 0x64, 0x69, 0x73, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x65, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x72, 0x65, 0x61, 0x12, 0x46, 0x0a, 0x11, 0x77, 0x61, 0x73, 0x74, 0x65,
	0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f,
	0x77, 0x61, 0x73, 0x74, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x59, 0x0a, 0x0e, 0x44, 0x69, 0x73, 0x70, 0x6f,
	0x73, 0x61, 0x6c, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1f, 0x0a, 0x1b, 0x44, 0x49, 0x53,
	0x50, 0x4f, 0x53, 0x41, 0x4c, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x47, 0x45,
	0x4e, 0x45, 0x52, 0x41, 0x4c, 0x5f, 0x57, 0x41, 0x53, 0x54, 0x45, 0x10, 0x01, 0x12, 0x13, 0x0a,
	0x0f, 0x4d, 0x49, 0x58, 0x45, 0x44, 0x5f, 0x52, 0x45, 0x43, 0x59, 0x43, 0x4c, 0x49, 0x4e, 0x47,
	0x10, 0x02, 0x22, 0xa1, 0x01, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x61, 0x73, 0x74, 0x65,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3e, 0x0a, 0x0c, 0x77, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x52, 0x0c, 0x77, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12,
	0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61,
	0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x22, 0xa2, 0x01, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x57,
	0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6d,
	0x61, 0x73, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x12,
	0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x89, 0x01, 0x0a, 0x17,
	0x50, 0x75, 0x6c, 0x6c, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x72,
	0x65, 0x61, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64,
	0x4d, 0x61, 0x73, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x5f,
	0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x73, 0x4f, 0x6e, 0x6c, 0x79, 0x22, 0xe3, 0x02, 0x0a, 0x18, 0x50, 0x75, 0x6c, 0x6c,
	0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x1a, 0xfc,
	0x01, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3b, 0x0a,
	0x0b, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x6e, 0x65,
	0x77, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x57, 0x61,
	0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x08, 0x6e, 0x65, 0x77, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x6f, 0x6c, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x52, 0x08, 0x6f, 0x6c, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2f, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x73, 0x6d, 0x61,
	0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x30, 0x0a,
	0x1a, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x75, 0x0a, 0x12, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x4b, 0x0a, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x75, 0x70, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x32, 0xd6, 0x01, 0x0a, 0x08, 0x57, 0x61, 0x73, 0x74, 0x65,
	0x41, 0x70, 0x69, 0x12, 0x63, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x61, 0x73, 0x74, 0x65,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x26, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x61, 0x73, 0x74,
	0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x27, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x10, 0x50, 0x75, 0x6c, 0x6c,
	0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x26, 0x2e, 0x73,
	0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c,
	0x6c, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x32,
	0x70, 0x0a, 0x09, 0x57, 0x61, 0x73, 0x74, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x63, 0x0a, 0x13,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x57, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x29, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x62, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x57, 0x61, 0x73, 0x74,
	0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21,
	0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x57,
	0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x76, 0x61, 0x6e, 0x74, 0x69, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x73, 0x63, 0x2d, 0x62, 0x6f, 0x73,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_waste_proto_rawDescOnce sync.Once
	file_waste_proto_rawDescData = file_waste_proto_rawDesc
)

func file_waste_proto_rawDescGZIP() []byte {
	file_waste_proto_rawDescOnce.Do(func() {
		file_waste_proto_rawDescData = protoimpl.X.CompressGZIP(file_waste_proto_rawDescData)
	})
	return file_waste_proto_rawDescData
}

var file_waste_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_waste_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_waste_proto_goTypes = []any{
	(WasteRecord_DisposalMethod)(0),         // 0: smartcore.bos.WasteRecord.DisposalMethod
	(*WasteRecord)(nil),                     // 1: smartcore.bos.WasteRecord
	(*ListWasteRecordsResponse)(nil),        // 2: smartcore.bos.ListWasteRecordsResponse
	(*ListWasteRecordsRequest)(nil),         // 3: smartcore.bos.ListWasteRecordsRequest
	(*PullWasteRecordsRequest)(nil),         // 4: smartcore.bos.PullWasteRecordsRequest
	(*PullWasteRecordsResponse)(nil),        // 5: smartcore.bos.PullWasteRecordsResponse
	(*DescribeWasteRecordRequest)(nil),      // 6: smartcore.bos.DescribeWasteRecordRequest
	(*WasteRecordSupport)(nil),              // 7: smartcore.bos.WasteRecordSupport
	(*PullWasteRecordsResponse_Change)(nil), // 8: smartcore.bos.PullWasteRecordsResponse.Change
	(*timestamppb.Timestamp)(nil),           // 9: google.protobuf.Timestamp
	(*fieldmaskpb.FieldMask)(nil),           // 10: google.protobuf.FieldMask
	(*types.ResourceSupport)(nil),           // 11: smartcore.types.ResourceSupport
	(types.ChangeType)(0),                   // 12: smartcore.types.ChangeType
}
var file_waste_proto_depIdxs = []int32{
	9,  // 0: smartcore.bos.WasteRecord.record_create_time:type_name -> google.protobuf.Timestamp
	0,  // 1: smartcore.bos.WasteRecord.disposal_method:type_name -> smartcore.bos.WasteRecord.DisposalMethod
	9,  // 2: smartcore.bos.WasteRecord.waste_create_time:type_name -> google.protobuf.Timestamp
	1,  // 3: smartcore.bos.ListWasteRecordsResponse.wasteRecords:type_name -> smartcore.bos.WasteRecord
	10, // 4: smartcore.bos.ListWasteRecordsRequest.read_mask:type_name -> google.protobuf.FieldMask
	10, // 5: smartcore.bos.PullWasteRecordsRequest.read_mask:type_name -> google.protobuf.FieldMask
	8,  // 6: smartcore.bos.PullWasteRecordsResponse.changes:type_name -> smartcore.bos.PullWasteRecordsResponse.Change
	11, // 7: smartcore.bos.WasteRecordSupport.resource_support:type_name -> smartcore.types.ResourceSupport
	9,  // 8: smartcore.bos.PullWasteRecordsResponse.Change.change_time:type_name -> google.protobuf.Timestamp
	1,  // 9: smartcore.bos.PullWasteRecordsResponse.Change.new_value:type_name -> smartcore.bos.WasteRecord
	1,  // 10: smartcore.bos.PullWasteRecordsResponse.Change.old_value:type_name -> smartcore.bos.WasteRecord
	12, // 11: smartcore.bos.PullWasteRecordsResponse.Change.type:type_name -> smartcore.types.ChangeType
	3,  // 12: smartcore.bos.WasteApi.ListWasteRecords:input_type -> smartcore.bos.ListWasteRecordsRequest
	4,  // 13: smartcore.bos.WasteApi.PullWasteRecords:input_type -> smartcore.bos.PullWasteRecordsRequest
	6,  // 14: smartcore.bos.WasteInfo.DescribeWasteRecord:input_type -> smartcore.bos.DescribeWasteRecordRequest
	2,  // 15: smartcore.bos.WasteApi.ListWasteRecords:output_type -> smartcore.bos.ListWasteRecordsResponse
	5,  // 16: smartcore.bos.WasteApi.PullWasteRecords:output_type -> smartcore.bos.PullWasteRecordsResponse
	7,  // 17: smartcore.bos.WasteInfo.DescribeWasteRecord:output_type -> smartcore.bos.WasteRecordSupport
	15, // [15:18] is the sub-list for method output_type
	12, // [12:15] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_waste_proto_init() }
func file_waste_proto_init() {
	if File_waste_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_waste_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*WasteRecord); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_waste_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ListWasteRecordsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_waste_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ListWasteRecordsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_waste_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*PullWasteRecordsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_waste_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*PullWasteRecordsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_waste_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DescribeWasteRecordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_waste_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*WasteRecordSupport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_waste_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*PullWasteRecordsResponse_Change); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_waste_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_waste_proto_goTypes,
		DependencyIndexes: file_waste_proto_depIdxs,
		EnumInfos:         file_waste_proto_enumTypes,
		MessageInfos:      file_waste_proto_msgTypes,
	}.Build()
	File_waste_proto = out.File
	file_waste_proto_rawDesc = nil
	file_waste_proto_goTypes = nil
	file_waste_proto_depIdxs = nil
}
