// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: anpr_camera.proto

package gen

import (
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

// AnprEvent describes a registration plate detection event from an ANPR camera.
// Includes registration plate information but also vehicle information if the camera supports it.
type AnprEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The time the detection events occurred.
	EventTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	// The registration plate detected.
	RegistrationPlate string `protobuf:"bytes,2,opt,name=registration_plate,json=registrationPlate,proto3" json:"registration_plate,omitempty"`
	// Country of the detected registration plate.
	Country string `protobuf:"bytes,3,opt,name=country,proto3" json:"country,omitempty"`
	// Optional. The area of the detected registration plate. i.e. for UAE this could be Abu Dhabi, Dubai, etc.
	Area *string `protobuf:"bytes,4,opt,name=area,proto3,oneof" json:"area,omitempty"`
	// Optional. The confidence level of the detection as a percentage. If omitted, means unknown confidence level.
	Confidence *int32 `protobuf:"varint,5,opt,name=confidence,proto3,oneof" json:"confidence,omitempty"`
	// Optional. The type of plate, e.g. standard, personalised, etc.
	PlateType *string `protobuf:"bytes,6,opt,name=plate_type,json=plateType,proto3,oneof" json:"plate_type,omitempty"`
	// Optional. The year of the vehicle.
	Year        *string                `protobuf:"bytes,7,opt,name=year,proto3,oneof" json:"year,omitempty"`
	VehicleInfo *AnprEvent_VehicleInfo `protobuf:"bytes,8,opt,name=vehicle_info,json=vehicleInfo,proto3,oneof" json:"vehicle_info,omitempty"`
}

func (x *AnprEvent) Reset() {
	*x = AnprEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anpr_camera_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnprEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnprEvent) ProtoMessage() {}

func (x *AnprEvent) ProtoReflect() protoreflect.Message {
	mi := &file_anpr_camera_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnprEvent.ProtoReflect.Descriptor instead.
func (*AnprEvent) Descriptor() ([]byte, []int) {
	return file_anpr_camera_proto_rawDescGZIP(), []int{0}
}

func (x *AnprEvent) GetEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *AnprEvent) GetRegistrationPlate() string {
	if x != nil {
		return x.RegistrationPlate
	}
	return ""
}

func (x *AnprEvent) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *AnprEvent) GetArea() string {
	if x != nil && x.Area != nil {
		return *x.Area
	}
	return ""
}

func (x *AnprEvent) GetConfidence() int32 {
	if x != nil && x.Confidence != nil {
		return *x.Confidence
	}
	return 0
}

func (x *AnprEvent) GetPlateType() string {
	if x != nil && x.PlateType != nil {
		return *x.PlateType
	}
	return ""
}

func (x *AnprEvent) GetYear() string {
	if x != nil && x.Year != nil {
		return *x.Year
	}
	return ""
}

func (x *AnprEvent) GetVehicleInfo() *AnprEvent_VehicleInfo {
	if x != nil {
		return x.VehicleInfo
	}
	return nil
}

type GetLastEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the device to get the last event for.
	Name     string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
}

func (x *GetLastEventRequest) Reset() {
	*x = GetLastEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anpr_camera_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLastEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLastEventRequest) ProtoMessage() {}

func (x *GetLastEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_anpr_camera_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLastEventRequest.ProtoReflect.Descriptor instead.
func (*GetLastEventRequest) Descriptor() ([]byte, []int) {
	return file_anpr_camera_proto_rawDescGZIP(), []int{1}
}

func (x *GetLastEventRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetLastEventRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

type PullEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the device to pull events for.
	Name        string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask    *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	UpdatesOnly bool                   `protobuf:"varint,3,opt,name=updates_only,json=updatesOnly,proto3" json:"updates_only,omitempty"`
}

func (x *PullEventsRequest) Reset() {
	*x = PullEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anpr_camera_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullEventsRequest) ProtoMessage() {}

func (x *PullEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_anpr_camera_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullEventsRequest.ProtoReflect.Descriptor instead.
func (*PullEventsRequest) Descriptor() ([]byte, []int) {
	return file_anpr_camera_proto_rawDescGZIP(), []int{2}
}

func (x *PullEventsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullEventsRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

func (x *PullEventsRequest) GetUpdatesOnly() bool {
	if x != nil {
		return x.UpdatesOnly
	}
	return false
}

type PullEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Changes []*PullEventsResponse_Change `protobuf:"bytes,1,rep,name=changes,proto3" json:"changes,omitempty"`
}

func (x *PullEventsResponse) Reset() {
	*x = PullEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anpr_camera_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullEventsResponse) ProtoMessage() {}

func (x *PullEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_anpr_camera_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullEventsResponse.ProtoReflect.Descriptor instead.
func (*PullEventsResponse) Descriptor() ([]byte, []int) {
	return file_anpr_camera_proto_rawDescGZIP(), []int{3}
}

func (x *PullEventsResponse) GetChanges() []*PullEventsResponse_Change {
	if x != nil {
		return x.Changes
	}
	return nil
}

// Optional. Information about the vehicle itself.
type AnprEvent_VehicleInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional. The type of vehicle, e.g. car, truck, etc.
	VehicleType string `protobuf:"bytes,1,opt,name=vehicle_type,json=vehicleType,proto3" json:"vehicle_type,omitempty"`
	// Optional. The colour of the vehicle.
	Colour string `protobuf:"bytes,2,opt,name=colour,proto3" json:"colour,omitempty"`
	// Optional. The make of the vehicle.
	Make string `protobuf:"bytes,3,opt,name=make,proto3" json:"make,omitempty"`
	// Optional. The model of the vehicle.
	Model string `protobuf:"bytes,4,opt,name=model,proto3" json:"model,omitempty"`
}

func (x *AnprEvent_VehicleInfo) Reset() {
	*x = AnprEvent_VehicleInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anpr_camera_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnprEvent_VehicleInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnprEvent_VehicleInfo) ProtoMessage() {}

func (x *AnprEvent_VehicleInfo) ProtoReflect() protoreflect.Message {
	mi := &file_anpr_camera_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnprEvent_VehicleInfo.ProtoReflect.Descriptor instead.
func (*AnprEvent_VehicleInfo) Descriptor() ([]byte, []int) {
	return file_anpr_camera_proto_rawDescGZIP(), []int{0, 0}
}

func (x *AnprEvent_VehicleInfo) GetVehicleType() string {
	if x != nil {
		return x.VehicleType
	}
	return ""
}

func (x *AnprEvent_VehicleInfo) GetColour() string {
	if x != nil {
		return x.Colour
	}
	return ""
}

func (x *AnprEvent_VehicleInfo) GetMake() string {
	if x != nil {
		return x.Make
	}
	return ""
}

func (x *AnprEvent_VehicleInfo) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

// The detection event.
type PullEventsResponse_Change struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ChangeTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=change_time,json=changeTime,proto3" json:"change_time,omitempty"`
	AnprEvent  *AnprEvent             `protobuf:"bytes,3,opt,name=anpr_event,json=anprEvent,proto3" json:"anpr_event,omitempty"`
}

func (x *PullEventsResponse_Change) Reset() {
	*x = PullEventsResponse_Change{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anpr_camera_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullEventsResponse_Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullEventsResponse_Change) ProtoMessage() {}

func (x *PullEventsResponse_Change) ProtoReflect() protoreflect.Message {
	mi := &file_anpr_camera_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullEventsResponse_Change.ProtoReflect.Descriptor instead.
func (*PullEventsResponse_Change) Descriptor() ([]byte, []int) {
	return file_anpr_camera_proto_rawDescGZIP(), []int{3, 0}
}

func (x *PullEventsResponse_Change) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullEventsResponse_Change) GetChangeTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ChangeTime
	}
	return nil
}

func (x *PullEventsResponse_Change) GetAnprEvent() *AnprEvent {
	if x != nil {
		return x.AnprEvent
	}
	return nil
}

var File_anpr_camera_proto protoreflect.FileDescriptor

var file_anpr_camera_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x6e, 0x70, 0x72, 0x5f, 0x63, 0x61, 0x6d, 0x65, 0x72, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62,
	0x6f, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8d, 0x04, 0x0a, 0x09, 0x41, 0x6e, 0x70, 0x72, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x2d,
	0x0a, 0x12, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x17, 0x0a, 0x04, 0x61, 0x72, 0x65, 0x61, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x61, 0x72, 0x65, 0x61, 0x88, 0x01, 0x01,
	0x12, 0x23, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x01, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x09, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x79, 0x65, 0x61,
	0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x04, 0x79, 0x65, 0x61, 0x72, 0x88,
	0x01, 0x01, 0x12, 0x4c, 0x0a, 0x0c, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x41, 0x6e, 0x70, 0x72, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x2e, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x04,
	0x52, 0x0b, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x88, 0x01, 0x01,
	0x1a, 0x72, 0x0a, 0x0b, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x21, 0x0a, 0x0c, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x6f, 0x75, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x6f, 0x75, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x61,
	0x6b, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x61, 0x6b, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x61, 0x72, 0x65, 0x61, 0x42, 0x0d, 0x0a,
	0x0b, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x42, 0x0d, 0x0a, 0x0b,
	0x5f, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x79, 0x65, 0x61, 0x72, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x62, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x73, 0x74,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x37, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52,
	0x08, 0x72, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x22, 0x83, 0x01, 0x0a, 0x11, 0x50, 0x75,
	0x6c, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61,
	0x73, 0x6b, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x12, 0x21, 0x0a, 0x0c,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x4f, 0x6e, 0x6c, 0x79, 0x22,
	0xed, 0x01, 0x0a, 0x12, 0x50, 0x75, 0x6c, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x1a, 0x92, 0x01, 0x0a, 0x06, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x61, 0x6e, 0x70, 0x72, 0x5f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73, 0x6d, 0x61,
	0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x41, 0x6e, 0x70, 0x72, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x61, 0x6e, 0x70, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x32,
	0xb2, 0x01, 0x0a, 0x0d, 0x41, 0x6e, 0x70, 0x72, 0x43, 0x61, 0x6d, 0x65, 0x72, 0x61, 0x41, 0x70,
	0x69, 0x12, 0x4a, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x22, 0x2e,
	0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x47, 0x65,
	0x74, 0x4c, 0x61, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x18, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f,
	0x73, 0x2e, 0x41, 0x6e, 0x70, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x55, 0x0a,
	0x0a, 0x50, 0x75, 0x6c, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x20, 0x2e, 0x73, 0x6d,
	0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c, 0x6c,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e,
	0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75,
	0x6c, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x30, 0x01, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x76, 0x61, 0x6e, 0x74, 0x69, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x73, 0x63, 0x2d,
	0x62, 0x6f, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_anpr_camera_proto_rawDescOnce sync.Once
	file_anpr_camera_proto_rawDescData = file_anpr_camera_proto_rawDesc
)

func file_anpr_camera_proto_rawDescGZIP() []byte {
	file_anpr_camera_proto_rawDescOnce.Do(func() {
		file_anpr_camera_proto_rawDescData = protoimpl.X.CompressGZIP(file_anpr_camera_proto_rawDescData)
	})
	return file_anpr_camera_proto_rawDescData
}

var file_anpr_camera_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_anpr_camera_proto_goTypes = []any{
	(*AnprEvent)(nil),                 // 0: smartcore.bos.AnprEvent
	(*GetLastEventRequest)(nil),       // 1: smartcore.bos.GetLastEventRequest
	(*PullEventsRequest)(nil),         // 2: smartcore.bos.PullEventsRequest
	(*PullEventsResponse)(nil),        // 3: smartcore.bos.PullEventsResponse
	(*AnprEvent_VehicleInfo)(nil),     // 4: smartcore.bos.AnprEvent.VehicleInfo
	(*PullEventsResponse_Change)(nil), // 5: smartcore.bos.PullEventsResponse.Change
	(*timestamppb.Timestamp)(nil),     // 6: google.protobuf.Timestamp
	(*fieldmaskpb.FieldMask)(nil),     // 7: google.protobuf.FieldMask
}
var file_anpr_camera_proto_depIdxs = []int32{
	6, // 0: smartcore.bos.AnprEvent.event_time:type_name -> google.protobuf.Timestamp
	4, // 1: smartcore.bos.AnprEvent.vehicle_info:type_name -> smartcore.bos.AnprEvent.VehicleInfo
	7, // 2: smartcore.bos.GetLastEventRequest.read_mask:type_name -> google.protobuf.FieldMask
	7, // 3: smartcore.bos.PullEventsRequest.read_mask:type_name -> google.protobuf.FieldMask
	5, // 4: smartcore.bos.PullEventsResponse.changes:type_name -> smartcore.bos.PullEventsResponse.Change
	6, // 5: smartcore.bos.PullEventsResponse.Change.change_time:type_name -> google.protobuf.Timestamp
	0, // 6: smartcore.bos.PullEventsResponse.Change.anpr_event:type_name -> smartcore.bos.AnprEvent
	1, // 7: smartcore.bos.AnprCameraApi.GetEvent:input_type -> smartcore.bos.GetLastEventRequest
	2, // 8: smartcore.bos.AnprCameraApi.PullEvents:input_type -> smartcore.bos.PullEventsRequest
	0, // 9: smartcore.bos.AnprCameraApi.GetEvent:output_type -> smartcore.bos.AnprEvent
	3, // 10: smartcore.bos.AnprCameraApi.PullEvents:output_type -> smartcore.bos.PullEventsResponse
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_anpr_camera_proto_init() }
func file_anpr_camera_proto_init() {
	if File_anpr_camera_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_anpr_camera_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*AnprEvent); i {
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
		file_anpr_camera_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetLastEventRequest); i {
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
		file_anpr_camera_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PullEventsRequest); i {
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
		file_anpr_camera_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*PullEventsResponse); i {
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
		file_anpr_camera_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*AnprEvent_VehicleInfo); i {
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
		file_anpr_camera_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*PullEventsResponse_Change); i {
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
	file_anpr_camera_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_anpr_camera_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_anpr_camera_proto_goTypes,
		DependencyIndexes: file_anpr_camera_proto_depIdxs,
		MessageInfos:      file_anpr_camera_proto_msgTypes,
	}.Build()
	File_anpr_camera_proto = out.File
	file_anpr_camera_proto_rawDesc = nil
	file_anpr_camera_proto_goTypes = nil
	file_anpr_camera_proto_depIdxs = nil
}
