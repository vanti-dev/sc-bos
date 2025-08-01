// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
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
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MeterReading struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Usage records the current energy consumed value of the meter.
	// The unit is unspecified, use device documentation or MeterInfo to discover it.
	// This value is a total recorded between the start and end times.
	Usage float32 `protobuf:"fixed32,1,opt,name=usage,proto3" json:"usage,omitempty"`
	// The start period usage is recorded relative to. Typically the installation date, but not required to be.
	// The start time can be reset and updated by the device if it is serviced or updated.
	StartTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// The end of the period usage is recorded relative to, i.e. the time the reading was taken.
	// This time might not be now if the device has a low resolution for taking readings.
	EndTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	// Produced records the current energy produced value of the meter.
	// The unit is unspecified, use device documentation or MeterInfo to discover it.
	// This value is a total recorded between the start and end times.
	Produced      float32 `protobuf:"fixed32,4,opt,name=produced,proto3" json:"produced,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MeterReading) Reset() {
	*x = MeterReading{}
	mi := &file_meter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MeterReading) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MeterReading) ProtoMessage() {}

func (x *MeterReading) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[0]
	if x != nil {
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

func (x *MeterReading) GetProduced() float32 {
	if x != nil {
		return x.Produced
	}
	return 0
}

// MeterReadingSupport describes the capabilities of devices implementing this trait
type MeterReadingSupport struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// How a named device supports read/write/pull apis
	ResourceSupport *types.ResourceSupport `protobuf:"bytes,1,opt,name=resource_support,json=resourceSupport,proto3" json:"resource_support,omitempty"`
	// The unit associated with the usage value
	UsageUnit string `protobuf:"bytes,2,opt,name=usage_unit,json=usageUnit,proto3" json:"usage_unit,omitempty"`
	// The unit associated with the produced value
	ProducedUnit  string `protobuf:"bytes,3,opt,name=produced_unit,json=producedUnit,proto3" json:"produced_unit,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MeterReadingSupport) Reset() {
	*x = MeterReadingSupport{}
	mi := &file_meter_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MeterReadingSupport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MeterReadingSupport) ProtoMessage() {}

func (x *MeterReadingSupport) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[1]
	if x != nil {
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

func (x *MeterReadingSupport) GetUsageUnit() string {
	if x != nil {
		return x.UsageUnit
	}
	return ""
}

func (x *MeterReadingSupport) GetProducedUnit() string {
	if x != nil {
		return x.ProducedUnit
	}
	return ""
}

type GetMeterReadingRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask      *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMeterReadingRequest) Reset() {
	*x = GetMeterReadingRequest{}
	mi := &file_meter_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMeterReadingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeterReadingRequest) ProtoMessage() {}

func (x *GetMeterReadingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[2]
	if x != nil {
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask      *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	UpdatesOnly   bool                   `protobuf:"varint,3,opt,name=updates_only,json=updatesOnly,proto3" json:"updates_only,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullMeterReadingsRequest) Reset() {
	*x = PullMeterReadingsRequest{}
	mi := &file_meter_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullMeterReadingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullMeterReadingsRequest) ProtoMessage() {}

func (x *PullMeterReadingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[3]
	if x != nil {
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
	state         protoimpl.MessageState              `protogen:"open.v1"`
	Changes       []*PullMeterReadingsResponse_Change `protobuf:"bytes,1,rep,name=changes,proto3" json:"changes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullMeterReadingsResponse) Reset() {
	*x = PullMeterReadingsResponse{}
	mi := &file_meter_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullMeterReadingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullMeterReadingsResponse) ProtoMessage() {}

func (x *PullMeterReadingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[4]
	if x != nil {
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DescribeMeterReadingRequest) Reset() {
	*x = DescribeMeterReadingRequest{}
	mi := &file_meter_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DescribeMeterReadingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeMeterReadingRequest) ProtoMessage() {}

func (x *DescribeMeterReadingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[5]
	if x != nil {
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ChangeTime    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=change_time,json=changeTime,proto3" json:"change_time,omitempty"`
	MeterReading  *MeterReading          `protobuf:"bytes,3,opt,name=meter_reading,json=meterReading,proto3" json:"meter_reading,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullMeterReadingsResponse_Change) Reset() {
	*x = PullMeterReadingsResponse_Change{}
	mi := &file_meter_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullMeterReadingsResponse_Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullMeterReadingsResponse_Change) ProtoMessage() {}

func (x *PullMeterReadingsResponse_Change) ProtoReflect() protoreflect.Message {
	mi := &file_meter_proto_msgTypes[6]
	if x != nil {
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

const file_meter_proto_rawDesc = "" +
	"\n" +
	"\vmeter.proto\x12\rsmartcore.bos\x1a google/protobuf/field_mask.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x10types/info.proto\"\xb2\x01\n" +
	"\fMeterReading\x12\x14\n" +
	"\x05usage\x18\x01 \x01(\x02R\x05usage\x129\n" +
	"\n" +
	"start_time\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\tstartTime\x125\n" +
	"\bend_time\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\aendTime\x12\x1a\n" +
	"\bproduced\x18\x04 \x01(\x02R\bproduced\"\xa6\x01\n" +
	"\x13MeterReadingSupport\x12K\n" +
	"\x10resource_support\x18\x01 \x01(\v2 .smartcore.types.ResourceSupportR\x0fresourceSupport\x12\x1d\n" +
	"\n" +
	"usage_unit\x18\x02 \x01(\tR\tusageUnit\x12#\n" +
	"\rproduced_unit\x18\x03 \x01(\tR\fproducedUnit\"e\n" +
	"\x16GetMeterReadingRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x127\n" +
	"\tread_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\breadMask\"\x8a\x01\n" +
	"\x18PullMeterReadingsRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x127\n" +
	"\tread_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\breadMask\x12!\n" +
	"\fupdates_only\x18\x03 \x01(\bR\vupdatesOnly\"\x84\x02\n" +
	"\x19PullMeterReadingsResponse\x12I\n" +
	"\achanges\x18\x01 \x03(\v2/.smartcore.bos.PullMeterReadingsResponse.ChangeR\achanges\x1a\x9b\x01\n" +
	"\x06Change\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12;\n" +
	"\vchange_time\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"changeTime\x12@\n" +
	"\rmeter_reading\x18\x03 \x01(\v2\x1b.smartcore.bos.MeterReadingR\fmeterReading\"1\n" +
	"\x1bDescribeMeterReadingRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name2\xcb\x01\n" +
	"\bMeterApi\x12U\n" +
	"\x0fGetMeterReading\x12%.smartcore.bos.GetMeterReadingRequest\x1a\x1b.smartcore.bos.MeterReading\x12h\n" +
	"\x11PullMeterReadings\x12'.smartcore.bos.PullMeterReadingsRequest\x1a(.smartcore.bos.PullMeterReadingsResponse0\x012s\n" +
	"\tMeterInfo\x12f\n" +
	"\x14DescribeMeterReading\x12*.smartcore.bos.DescribeMeterReadingRequest\x1a\".smartcore.bos.MeterReadingSupportB%Z#github.com/vanti-dev/sc-bos/pkg/genb\x06proto3"

var (
	file_meter_proto_rawDescOnce sync.Once
	file_meter_proto_rawDescData []byte
)

func file_meter_proto_rawDescGZIP() []byte {
	file_meter_proto_rawDescOnce.Do(func() {
		file_meter_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_meter_proto_rawDesc), len(file_meter_proto_rawDesc)))
	})
	return file_meter_proto_rawDescData
}

var file_meter_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_meter_proto_goTypes = []any{
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_meter_proto_rawDesc), len(file_meter_proto_rawDesc)),
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
	file_meter_proto_goTypes = nil
	file_meter_proto_depIdxs = nil
}
