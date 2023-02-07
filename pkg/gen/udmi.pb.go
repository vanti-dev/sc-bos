// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.2
// source: udmi.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DescribeTopicsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DescribeTopicsRequest) Reset() {
	*x = DescribeTopicsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_udmi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeTopicsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeTopicsRequest) ProtoMessage() {}

func (x *DescribeTopicsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_udmi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeTopicsRequest.ProtoReflect.Descriptor instead.
func (*DescribeTopicsRequest) Descriptor() ([]byte, []int) {
	return file_udmi_proto_rawDescGZIP(), []int{0}
}

func (x *DescribeTopicsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DescribeTopicsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Topics []string `protobuf:"bytes,2,rep,name=topics,proto3" json:"topics,omitempty"`
}

func (x *DescribeTopicsResponse) Reset() {
	*x = DescribeTopicsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_udmi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeTopicsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeTopicsResponse) ProtoMessage() {}

func (x *DescribeTopicsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_udmi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeTopicsResponse.ProtoReflect.Descriptor instead.
func (*DescribeTopicsResponse) Descriptor() ([]byte, []int) {
	return file_udmi_proto_rawDescGZIP(), []int{1}
}

func (x *DescribeTopicsResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DescribeTopicsResponse) GetTopics() []string {
	if x != nil {
		return x.Topics
	}
	return nil
}

type OnMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Message *MqttMessage `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *OnMessageRequest) Reset() {
	*x = OnMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_udmi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnMessageRequest) ProtoMessage() {}

func (x *OnMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_udmi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnMessageRequest.ProtoReflect.Descriptor instead.
func (*OnMessageRequest) Descriptor() ([]byte, []int) {
	return file_udmi_proto_rawDescGZIP(), []int{2}
}

func (x *OnMessageRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *OnMessageRequest) GetMessage() *MqttMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

type OnMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *OnMessageResponse) Reset() {
	*x = OnMessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_udmi_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnMessageResponse) ProtoMessage() {}

func (x *OnMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_udmi_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnMessageResponse.ProtoReflect.Descriptor instead.
func (*OnMessageResponse) Descriptor() ([]byte, []int) {
	return file_udmi_proto_rawDescGZIP(), []int{3}
}

func (x *OnMessageResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type PullExportMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *PullExportMessagesRequest) Reset() {
	*x = PullExportMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_udmi_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullExportMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullExportMessagesRequest) ProtoMessage() {}

func (x *PullExportMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_udmi_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullExportMessagesRequest.ProtoReflect.Descriptor instead.
func (*PullExportMessagesRequest) Descriptor() ([]byte, []int) {
	return file_udmi_proto_rawDescGZIP(), []int{4}
}

func (x *PullExportMessagesRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type PullExportMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Message *MqttMessage `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PullExportMessagesResponse) Reset() {
	*x = PullExportMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_udmi_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullExportMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullExportMessagesResponse) ProtoMessage() {}

func (x *PullExportMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_udmi_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullExportMessagesResponse.ProtoReflect.Descriptor instead.
func (*PullExportMessagesResponse) Descriptor() ([]byte, []int) {
	return file_udmi_proto_rawDescGZIP(), []int{5}
}

func (x *PullExportMessagesResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullExportMessagesResponse) GetMessage() *MqttMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

type MqttMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic   string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Payload string `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"` // JSON payload
}

func (x *MqttMessage) Reset() {
	*x = MqttMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_udmi_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MqttMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MqttMessage) ProtoMessage() {}

func (x *MqttMessage) ProtoReflect() protoreflect.Message {
	mi := &file_udmi_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MqttMessage.ProtoReflect.Descriptor instead.
func (*MqttMessage) Descriptor() ([]byte, []int) {
	return file_udmi_proto_rawDescGZIP(), []int{6}
}

func (x *MqttMessage) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *MqttMessage) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

var File_udmi_proto protoreflect.FileDescriptor

var file_udmi_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x64, 0x6d, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73, 0x6d,
	0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x22, 0x2b, 0x0a, 0x15, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x44, 0x0a, 0x16, 0x44, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x22, 0x5c,
	0x0a, 0x10, 0x4f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x4d, 0x71, 0x74, 0x74, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x27, 0x0a, 0x11,
	0x4f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2f, 0x0a, 0x19, 0x50, 0x75, 0x6c, 0x6c, 0x45, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x66, 0x0a, 0x1a, 0x50, 0x75, 0x6c, 0x6c, 0x45, 0x78,
	0x70, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x6d, 0x61, 0x72,
	0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x4d, 0x71, 0x74, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3d,
	0x0a, 0x0b, 0x4d, 0x71, 0x74, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x32, 0xa9, 0x02,
	0x0a, 0x0b, 0x55, 0x64, 0x6d, 0x69, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5d, 0x0a,
	0x0e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x12,
	0x24, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f,
	0x70, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x09,
	0x4f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x2e, 0x73, 0x6d, 0x61, 0x72,
	0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x4f, 0x6e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x73, 0x6d, 0x61,
	0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x4f, 0x6e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6b, 0x0a, 0x12,
	0x50, 0x75, 0x6c, 0x6c, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x12, 0x28, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62,
	0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x73,
	0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x50, 0x75, 0x6c,
	0x6c, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x61, 0x6e, 0x74, 0x69, 0x2d, 0x64, 0x65,
	0x76, 0x2f, 0x73, 0x63, 0x2d, 0x62, 0x6f, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65, 0x6e,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_udmi_proto_rawDescOnce sync.Once
	file_udmi_proto_rawDescData = file_udmi_proto_rawDesc
)

func file_udmi_proto_rawDescGZIP() []byte {
	file_udmi_proto_rawDescOnce.Do(func() {
		file_udmi_proto_rawDescData = protoimpl.X.CompressGZIP(file_udmi_proto_rawDescData)
	})
	return file_udmi_proto_rawDescData
}

var file_udmi_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_udmi_proto_goTypes = []interface{}{
	(*DescribeTopicsRequest)(nil),      // 0: smartcore.bos.DescribeTopicsRequest
	(*DescribeTopicsResponse)(nil),     // 1: smartcore.bos.DescribeTopicsResponse
	(*OnMessageRequest)(nil),           // 2: smartcore.bos.OnMessageRequest
	(*OnMessageResponse)(nil),          // 3: smartcore.bos.OnMessageResponse
	(*PullExportMessagesRequest)(nil),  // 4: smartcore.bos.PullExportMessagesRequest
	(*PullExportMessagesResponse)(nil), // 5: smartcore.bos.PullExportMessagesResponse
	(*MqttMessage)(nil),                // 6: smartcore.bos.MqttMessage
}
var file_udmi_proto_depIdxs = []int32{
	6, // 0: smartcore.bos.OnMessageRequest.message:type_name -> smartcore.bos.MqttMessage
	6, // 1: smartcore.bos.PullExportMessagesResponse.message:type_name -> smartcore.bos.MqttMessage
	0, // 2: smartcore.bos.UdmiService.DescribeTopics:input_type -> smartcore.bos.DescribeTopicsRequest
	2, // 3: smartcore.bos.UdmiService.OnMessage:input_type -> smartcore.bos.OnMessageRequest
	4, // 4: smartcore.bos.UdmiService.PullExportMessages:input_type -> smartcore.bos.PullExportMessagesRequest
	1, // 5: smartcore.bos.UdmiService.DescribeTopics:output_type -> smartcore.bos.DescribeTopicsResponse
	3, // 6: smartcore.bos.UdmiService.OnMessage:output_type -> smartcore.bos.OnMessageResponse
	5, // 7: smartcore.bos.UdmiService.PullExportMessages:output_type -> smartcore.bos.PullExportMessagesResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_udmi_proto_init() }
func file_udmi_proto_init() {
	if File_udmi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_udmi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeTopicsRequest); i {
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
		file_udmi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeTopicsResponse); i {
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
		file_udmi_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnMessageRequest); i {
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
		file_udmi_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnMessageResponse); i {
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
		file_udmi_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullExportMessagesRequest); i {
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
		file_udmi_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullExportMessagesResponse); i {
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
		file_udmi_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MqttMessage); i {
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
			RawDescriptor: file_udmi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_udmi_proto_goTypes,
		DependencyIndexes: file_udmi_proto_depIdxs,
		MessageInfos:      file_udmi_proto_msgTypes,
	}.Build()
	File_udmi_proto = out.File
	file_udmi_proto_rawDesc = nil
	file_udmi_proto_goTypes = nil
	file_udmi_proto_depIdxs = nil
}
