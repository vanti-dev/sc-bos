// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: meter.proto

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

type MeterReading struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Usage records the current value of the meter.
	// The unit is unspecified, use device documentation or MeterInfo to discover it.
	// This value is a total recorded between the start and end times.
	Usage float32 `protobuf:"fixed32,1,opt,name=usage,proto3" json:"usage,omitempty"`
	// The start period usage is recorded relative to. Typically the installation date, but not required to be.
	// The start time can be reset and updated by the device if it is serviced or updated.
	StartTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// The end of the period usage is recorded relative to, i.e. the time the reading was taken.
	// This time might not be now if the device has a low resolution for taking readings.
	EndTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (x *MeterReading) Reset() {
	*x = MeterReading{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MeterReading) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MeterReading) ProtoMessage() {}

func (x *MeterReading) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MeterReading.ProtoReflect.Descriptor instead.
func (*MeterReading) Descriptor() ([]byte, []int) {
	return file_meter_proto_rawDescGZIP(), []int{0}
}

func (x *MeterReading) GetUsage() float32 {
	if x != nil {
		return x.Usage
	}
	return 0
}

func (x *MeterReading) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *MeterReading) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

// MeterReadingSupport describes the capabilities of devices implementing this trait
type MeterReadingSupport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// How a named device supports read/write/pull apis
	ResourceSupport *types.ResourceSupport `protobuf:"bytes,1,opt,name=resource_support,json=resourceSupport,proto3" json:"resource_support,omitempty"`
	// The unit associated with the usage value
	Unit string `protobuf:"bytes,2,opt,name=unit,proto3" json:"unit,omitempty"`
}

func (x *MeterReadingSupport) Reset() {
	*x = MeterReadingSupport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MeterReadingSupport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MeterReadingSupport) ProtoMessage() {}

func (x *MeterReadingSupport) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MeterReadingSupport.ProtoReflect.Descriptor instead.
func (*MeterReadingSupport) Descriptor() ([]byte, []int) {
	return file_meter_proto_rawDescGZIP(), []int{1}
}

func (x *MeterReadingSupport) GetResourceSupport() *types.ResourceSupport {
	if x != nil {
		return x.ResourceSupport
	}
	return nil
}

func (x *MeterReadingSupport) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

type GetMeterReadingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
}

func (x *GetMeterReadingRequest) Reset() {
	*x = GetMeterReadingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMeterReadingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeterReadingRequest) ProtoMessage() {}

func (x *GetMeterReadingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMeterReadingRequest.ProtoReflect.Descriptor instead.
func (*GetMeterReadingRequest) Descriptor() ([]byte, []int) {
	return file_meter_proto_rawDescGZIP(), []int{2}
}

func (x *GetMeterReadingRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetMeterReadingRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

type PullMeterReadingsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask    *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	UpdatesOnly bool                   `protobuf:"varint,3,opt,name=updates_only,json=updatesOnly,proto3" json:"updates_only,omitempty"`
}

func (x *PullMeterReadingsRequest) Reset() {
	*x = PullMeterReadingsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullMeterReadingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullMeterReadingsRequest) ProtoMessage() {}

func (x *PullMeterReadingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullMeterReadingsRequest.ProtoReflect.Descriptor instead.
func (*PullMeterReadingsRequest) Descriptor() ([]byte, []int) {
	return file_meter_proto_rawDescGZIP(), []int{3}
}

func (x *PullMeterReadingsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullMeterReadingsRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

func (x *PullMeterReadingsRequest) GetUpdatesOnly() bool {
	if x != nil {
		return x.UpdatesOnly
	}
	return false
}

type PullMeterReadingsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Changes []*PullMeterReadingsResponse_Change `protobuf:"bytes,1,rep,name=changes,proto3" json:"changes,omitempty"`
}

func (x *PullMeterReadingsResponse) Reset() {
	*x = PullMeterReadingsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullMeterReadingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullMeterReadingsResponse) ProtoMessage() {}

func (x *PullMeterReadingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullMeterReadingsResponse.ProtoReflect.Descriptor instead.
func (*PullMeterReadingsResponse) Descriptor() ([]byte, []int) {
	return file_meter_proto_rawDescGZIP(), []int{4}
}

func (x *PullMeterReadingsResponse) GetChanges() []*PullMeterReadingsResponse_Change {
	if x != nil {
		return x.Changes
	}
	return nil
}

type DescribeMeterReadingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DescribeMeterReadingRequest) Reset() {
	*x = DescribeMeterReadingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeMeterReadingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeMeterReadingRequest) ProtoMessage() {}

func (x *DescribeMeterReadingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeMeterReadingRequest.ProtoReflect.Descriptor instead.
func (*DescribeMeterReadingRequest) Descriptor() ([]byte, []int) {
	return file_meter_proto_rawDescGZIP(), []int{5}
}

func (x *DescribeMeterReadingRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type PullMeterReadingsResponse_Change struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ChangeTime   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=change_time,json=changeTime,proto3" json:"change_time,omitempty"`
	MeterReading *MeterReading          `protobuf:"bytes,3,opt,name=meter_reading,json=meterReading,proto3" json:"meter_reading,omitempty"`
}

func (x *PullMeterReadingsResponse_Change) Reset() {
	*x = PullMeterReadingsResponse_Change{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meter_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullMeterReadingsResponse_Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullMeterReadingsResponse_Change) ProtoMessage() {}

func (x *PullMeterReadingsResponse_Change) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullMeterReadingsResponse_Change.ProtoReflect.Descriptor instead.
func (*PullMeterReadingsResponse_Change) Descriptor() ([]byte, []int) {
	return file_meter_proto_rawDescGZIP(), []int{4, 0}
}

func (x *PullMeterReadingsResponse_Change) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullMeterReadingsResponse_Change) GetChangeTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ChangeTime
	}
	return nil
}

func (x *PullMeterReadingsResponse_Change) GetMeterReading() *MeterReading {
	if x != nil {
		return x.MeterReading
	}
	return nil
}

var File_meter_proto protoreflect.FileDescriptor

var file_meter_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73,
	0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x1a, 0x20, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x10, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x96, 0x01, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69,
	0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x76, 0x0a, 0x13, 0x4d, 0x65,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x4b, 0x0a, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x73, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x73, 0x6d,
	0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x0f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x6e,
	0x69, 0x74, 0x22, 0x65, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x61, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x37, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52,
	0x08, 0x72, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x22, 0x8a, 0x01, 0x0a, 0x18, 0x50, 0x75,
	0x6c, 0x6c, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x72, 0x65,
	0x61, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x4d,
	0x61, 0x73, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x5f, 0x6f,
	0x6e, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x73, 0x4f, 0x6e, 0x6c, 0x79, 0x22, 0x84, 0x02, 0x0a, 0x19, 0x50, 0x75, 0x6c, 0x6c, 0x4d,
	0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x1a,
	0x9b, 0x01, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3b,
	0x0a, 0x0b, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0a, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x40, 0x0a, 0x0d, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62,
	0x6f, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x52,
	0x0c, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x22, 0x31, 0x0a,
	0x1b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x61, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x32, 0xcb, 0x01, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x41, 0x70, 0x69, 0x12, 0x55, 0x0a,
	0x0f, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67,
	0x12, 0x25, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73,
	0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x12, 0x68, 0x0a, 0x11, 0x50, 0x75, 0x6c, 0x6c, 0x4d, 0x65, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x27, 0x2e, 0x73, 0x6d, 0x61, 0x72,
	0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x4d, 0x65,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x28, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62,
	0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x32, 0x73,
	0x0a, 0x09, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x66, 0x0a, 0x14, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x12, 0x2a, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x62, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4d, 0x65, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x22, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e,
	0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x75, 0x70, 0x70,
	0x6f, 0x72, 0x74, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x76, 0x61, 0x6e, 0x74, 0x69, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x73, 0x63, 0x2d, 0x62,
	0x6f, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_meter_proto_rawDescOnce sync.Once
	file_meter_proto_rawDescData = file_meter_proto_rawDesc
)

func file_meter_proto_rawDescGZIP() []byte {
	file_meter_proto_rawDescOnce.Do(func() {
		file_meter_proto_rawDescData = protoimpl.X.CompressGZIP(file_meter_proto_rawDescData)
	})
	return file_meter_proto_rawDescData
}

var file_meter_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_meter_proto_goTypes = []interface{}{
	(*MeterReading)(nil),                     // 0: smartcore.bos.MeterReading
	(*MeterReadingSupport)(nil),              // 1: smartcore.bos.MeterReadingSupport
	(*GetMeterReadingRequest)(nil),           // 2: smartcore.bos.GetMeterReadingRequest
	(*PullMeterReadingsRequest)(nil),         // 3: smartcore.bos.PullMeterReadingsRequest
	(*PullMeterReadingsResponse)(nil),        // 4: smartcore.bos.PullMeterReadingsResponse
	(*DescribeMeterReadingRequest)(nil),      // 5: smartcore.bos.DescribeMeterReadingRequest
	(*PullMeterReadingsResponse_Change)(nil), // 6: smartcore.bos.PullMeterReadingsResponse.Change
	(*timestamppb.Timestamp)(nil),            // 7: google.protobuf.Timestamp
	(*types.ResourceSupport)(nil),            // 8: smartcore.types.ResourceSupport
	(*fieldmaskpb.FieldMask)(nil),            // 9: google.protobuf.FieldMask
}
var file_meter_proto_depIdxs = []int32{
	7,  // 0: smartcore.bos.MeterReading.start_time:type_name -> google.protobuf.Timestamp
	7,  // 1: smartcore.bos.MeterReading.end_time:type_name -> google.protobuf.Timestamp
	8,  // 2: smartcore.bos.MeterReadingSupport.resource_support:type_name -> smartcore.types.ResourceSupport
	9,  // 3: smartcore.bos.GetMeterReadingRequest.read_mask:type_name -> google.protobuf.FieldMask
	9,  // 4: smartcore.bos.PullMeterReadingsRequest.read_mask:type_name -> google.protobuf.FieldMask
	6,  // 5: smartcore.bos.PullMeterReadingsResponse.changes:type_name -> smartcore.bos.PullMeterReadingsResponse.Change
	7,  // 6: smartcore.bos.PullMeterReadingsResponse.Change.change_time:type_name -> google.protobuf.Timestamp
	0,  // 7: smartcore.bos.PullMeterReadingsResponse.Change.meter_reading:type_name -> smartcore.bos.MeterReading
	2,  // 8: smartcore.bos.MeterApi.GetMeterReading:input_type -> smartcore.bos.GetMeterReadingRequest
	3,  // 9: smartcore.bos.MeterApi.PullMeterReadings:input_type -> smartcore.bos.PullMeterReadingsRequest
	5,  // 10: smartcore.bos.MeterInfo.DescribeMeterReading:input_type -> smartcore.bos.DescribeMeterReadingRequest
	0,  // 11: smartcore.bos.MeterApi.GetMeterReading:output_type -> smartcore.bos.MeterReading
	4,  // 12: smartcore.bos.MeterApi.PullMeterReadings:output_type -> smartcore.bos.PullMeterReadingsResponse
	1,  // 13: smartcore.bos.MeterInfo.DescribeMeterReading:output_type -> smartcore.bos.MeterReadingSupport
	11, // [11:14] is the sub-list for method output_type
	8,  // [8:11] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_meter_proto_init() }
func file_meter_proto_init() {
	if File_meter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_meter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MeterReading); i {
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
		file_meter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MeterReadingSupport); i {
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
		file_meter_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMeterReadingRequest); i {
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
		file_meter_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullMeterReadingsRequest); i {
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
		file_meter_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullMeterReadingsResponse); i {
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
		file_meter_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeMeterReadingRequest); i {
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
		file_meter_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullMeterReadingsResponse_Change); i {
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
			RawDescriptor: file_meter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_meter_proto_goTypes,
		DependencyIndexes: file_meter_proto_depIdxs,
		MessageInfos:      file_meter_proto_msgTypes,
	}.Build()
	File_meter_proto = out.File
	file_meter_proto_rawDesc = nil
	file_meter_proto_goTypes = nil
	file_meter_proto_depIdxs = nil
}
