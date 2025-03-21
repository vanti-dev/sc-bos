// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.1
// source: service_ticket.proto

package gen

import (
	types "github.com/smart-core-os/sc-api/go/types"
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

// Ticket represents a service ticket in a third party system.
type Ticket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Id is blank when creating a ticket, and is filled in by the external system. The ID is then used to update the ticket.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Short summary of the issue being reported.
	Summary string `protobuf:"bytes,2,opt,name=summary,proto3" json:"summary,omitempty"`
	// Full description on the issue being reported. This should include all the available information to help resolve the issue.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// Who reported the issue, this could be a person or an automation. e.g. "Bob" or "Cleaning Assistant".
	ReporterName string `protobuf:"bytes,4,opt,name=reporter_name,json=reporterName,proto3" json:"reporter_name,omitempty"`
	// What type of issue is this. e.g. "Fault", "Cleaning", "Maintenance".
	// Values supported by the implementing system are discovered via the ServiceTicketInfo service.
	Classification *Ticket_Classification `protobuf:"bytes,5,opt,name=classification,proto3" json:"classification,omitempty"`
	// How severe is the issue. e.g. "Critical", "High", "Medium", "Low".
	// Values supported by the implementing system are discovered via the ServiceTicketInfo service.
	Severity *Ticket_Severity `protobuf:"bytes,6,opt,name=severity,proto3" json:"severity,omitempty"`
	// Optional. A url that points to more information on this ticket
	ExternalUrl string `protobuf:"bytes,7,opt,name=external_url,json=externalUrl,proto3" json:"external_url,omitempty"`
}

func (x *Ticket) Reset() {
	*x = Ticket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_ticket_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ticket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ticket) ProtoMessage() {}

func (x *Ticket) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ticket.ProtoReflect.Descriptor instead.
func (*Ticket) Descriptor() ([]byte, []int) {
	return file_service_ticket_proto_rawDescGZIP(), []int{0}
}

func (x *Ticket) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Ticket) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *Ticket) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Ticket) GetReporterName() string {
	if x != nil {
		return x.ReporterName
	}
	return ""
}

func (x *Ticket) GetClassification() *Ticket_Classification {
	if x != nil {
		return x.Classification
	}
	return nil
}

func (x *Ticket) GetSeverity() *Ticket_Severity {
	if x != nil {
		return x.Severity
	}
	return nil
}

func (x *Ticket) GetExternalUrl() string {
	if x != nil {
		return x.ExternalUrl
	}
	return ""
}

// CreateTicketRequest is the request to create a ticket in the external system.
type CreateTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Ticket *Ticket `protobuf:"bytes,2,opt,name=ticket,proto3" json:"ticket,omitempty"`
}

func (x *CreateTicketRequest) Reset() {
	*x = CreateTicketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_ticket_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTicketRequest) ProtoMessage() {}

func (x *CreateTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTicketRequest.ProtoReflect.Descriptor instead.
func (*CreateTicketRequest) Descriptor() ([]byte, []int) {
	return file_service_ticket_proto_rawDescGZIP(), []int{1}
}

func (x *CreateTicketRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateTicketRequest) GetTicket() *Ticket {
	if x != nil {
		return x.Ticket
	}
	return nil
}

// UpdateTicketRequest is the request to update a ticket in the external system. The ticket ID must be set.
type UpdateTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Ticket *Ticket `protobuf:"bytes,2,opt,name=ticket,proto3" json:"ticket,omitempty"`
}

func (x *UpdateTicketRequest) Reset() {
	*x = UpdateTicketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_ticket_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTicketRequest) ProtoMessage() {}

func (x *UpdateTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTicketRequest.ProtoReflect.Descriptor instead.
func (*UpdateTicketRequest) Descriptor() ([]byte, []int) {
	return file_service_ticket_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateTicketRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateTicketRequest) GetTicket() *Ticket {
	if x != nil {
		return x.Ticket
	}
	return nil
}

type DescribeTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DescribeTicketRequest) Reset() {
	*x = DescribeTicketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_ticket_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeTicketRequest) ProtoMessage() {}

func (x *DescribeTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeTicketRequest.ProtoReflect.Descriptor instead.
func (*DescribeTicketRequest) Descriptor() ([]byte, []int) {
	return file_service_ticket_proto_rawDescGZIP(), []int{3}
}

func (x *DescribeTicketRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type TicketSupport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// How a named device supports read/write/pull apis
	ResourceSupport *types.ResourceSupport `protobuf:"bytes,1,opt,name=resource_support,json=resourceSupport,proto3" json:"resource_support,omitempty"`
	// The classifications supported by the implementing system.
	Classifications []*Ticket_Classification `protobuf:"bytes,2,rep,name=classifications,proto3" json:"classifications,omitempty"`
	// The severities supported by the implementing system.
	Severities []*Ticket_Severity `protobuf:"bytes,3,rep,name=severities,proto3" json:"severities,omitempty"`
}

func (x *TicketSupport) Reset() {
	*x = TicketSupport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_ticket_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TicketSupport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TicketSupport) ProtoMessage() {}

func (x *TicketSupport) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TicketSupport.ProtoReflect.Descriptor instead.
func (*TicketSupport) Descriptor() ([]byte, []int) {
	return file_service_ticket_proto_rawDescGZIP(), []int{4}
}

func (x *TicketSupport) GetResourceSupport() *types.ResourceSupport {
	if x != nil {
		return x.ResourceSupport
	}
	return nil
}

func (x *TicketSupport) GetClassifications() []*Ticket_Classification {
	if x != nil {
		return x.Classifications
	}
	return nil
}

func (x *TicketSupport) GetSeverities() []*Ticket_Severity {
	if x != nil {
		return x.Severities
	}
	return nil
}

type Ticket_Classification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The title of the classification.
	// This is unique within the context of the implementing system.
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Optional. A more detailed description can be displayed to a user to help them decide the correct classification.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Ticket_Classification) Reset() {
	*x = Ticket_Classification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_ticket_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ticket_Classification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ticket_Classification) ProtoMessage() {}

func (x *Ticket_Classification) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ticket_Classification.ProtoReflect.Descriptor instead.
func (*Ticket_Classification) Descriptor() ([]byte, []int) {
	return file_service_ticket_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Ticket_Classification) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Ticket_Classification) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type Ticket_Severity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The title of the severity.
	// This is unique within the context of the implementing system.
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Optional. A more detailed description can be displayed to a user to help them decide the correct severity.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Ticket_Severity) Reset() {
	*x = Ticket_Severity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_ticket_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ticket_Severity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ticket_Severity) ProtoMessage() {}

func (x *Ticket_Severity) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ticket_Severity.ProtoReflect.Descriptor instead.
func (*Ticket_Severity) Descriptor() ([]byte, []int) {
	return file_service_ticket_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Ticket_Severity) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Ticket_Severity) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

var File_service_ticket_proto protoreflect.FileDescriptor

var file_service_ticket_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x62, 0x6f, 0x73, 0x1a, 0x10, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x66,
	0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb4, 0x03, 0x0a, 0x06, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23,
	0x0a, 0x0d, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x4c, 0x0a, 0x0e, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x73, 0x6d,
	0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x74, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0e, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x3a, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x62, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x53, 0x65, 0x76, 0x65, 0x72,
	0x69, 0x74, 0x79, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x21, 0x0a,
	0x0c, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x55, 0x72, 0x6c,
	0x1a, 0x48, 0x0a, 0x0e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x42, 0x0a, 0x08, 0x53, 0x65,
	0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x58,
	0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x74, 0x69, 0x63,
	0x6b, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x6d, 0x61, 0x72,
	0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x52, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x58, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x62, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x74, 0x69, 0x63, 0x6b,
	0x65, 0x74, 0x22, 0x2b, 0x0a, 0x15, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x69,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0xec, 0x01, 0x0a, 0x0d, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x4b, 0x0a, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x73, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x73, 0x6d,
	0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x0f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x4e,
	0x0a, 0x0f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x43,
	0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x63,
	0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3e,
	0x0a, 0x0a, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62,
	0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69,
	0x74, 0x79, 0x52, 0x0a, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x32, 0xac,
	0x01, 0x0a, 0x10, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x41, 0x70, 0x69, 0x12, 0x4b, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x63,
	0x6b, 0x65, 0x74, 0x12, 0x22, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x62, 0x6f, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x00,
	0x12, 0x4b, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x12, 0x22, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x00, 0x32, 0x6b, 0x0a,
	0x11, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x56, 0x0a, 0x0e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x69,
	0x63, 0x6b, 0x65, 0x74, 0x12, 0x24, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x69, 0x63,
	0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x73, 0x6d, 0x61,
	0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x62, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65,
	0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x00, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x61, 0x6e, 0x74, 0x69, 0x2d, 0x64,
	0x65, 0x76, 0x2f, 0x73, 0x63, 0x2d, 0x62, 0x6f, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65,
	0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_ticket_proto_rawDescOnce sync.Once
	file_service_ticket_proto_rawDescData = file_service_ticket_proto_rawDesc
)

func file_service_ticket_proto_rawDescGZIP() []byte {
	file_service_ticket_proto_rawDescOnce.Do(func() {
		file_service_ticket_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_ticket_proto_rawDescData)
	})
	return file_service_ticket_proto_rawDescData
}

var file_service_ticket_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_service_ticket_proto_goTypes = []any{
	(*Ticket)(nil),                // 0: smartcore.bos.Ticket
	(*CreateTicketRequest)(nil),   // 1: smartcore.bos.CreateTicketRequest
	(*UpdateTicketRequest)(nil),   // 2: smartcore.bos.UpdateTicketRequest
	(*DescribeTicketRequest)(nil), // 3: smartcore.bos.DescribeTicketRequest
	(*TicketSupport)(nil),         // 4: smartcore.bos.TicketSupport
	(*Ticket_Classification)(nil), // 5: smartcore.bos.Ticket.Classification
	(*Ticket_Severity)(nil),       // 6: smartcore.bos.Ticket.Severity
	(*types.ResourceSupport)(nil), // 7: smartcore.types.ResourceSupport
}
var file_service_ticket_proto_depIdxs = []int32{
	5,  // 0: smartcore.bos.Ticket.classification:type_name -> smartcore.bos.Ticket.Classification
	6,  // 1: smartcore.bos.Ticket.severity:type_name -> smartcore.bos.Ticket.Severity
	0,  // 2: smartcore.bos.CreateTicketRequest.ticket:type_name -> smartcore.bos.Ticket
	0,  // 3: smartcore.bos.UpdateTicketRequest.ticket:type_name -> smartcore.bos.Ticket
	7,  // 4: smartcore.bos.TicketSupport.resource_support:type_name -> smartcore.types.ResourceSupport
	5,  // 5: smartcore.bos.TicketSupport.classifications:type_name -> smartcore.bos.Ticket.Classification
	6,  // 6: smartcore.bos.TicketSupport.severities:type_name -> smartcore.bos.Ticket.Severity
	1,  // 7: smartcore.bos.ServiceTicketApi.CreateTicket:input_type -> smartcore.bos.CreateTicketRequest
	2,  // 8: smartcore.bos.ServiceTicketApi.UpdateTicket:input_type -> smartcore.bos.UpdateTicketRequest
	3,  // 9: smartcore.bos.ServiceTicketInfo.DescribeTicket:input_type -> smartcore.bos.DescribeTicketRequest
	0,  // 10: smartcore.bos.ServiceTicketApi.CreateTicket:output_type -> smartcore.bos.Ticket
	0,  // 11: smartcore.bos.ServiceTicketApi.UpdateTicket:output_type -> smartcore.bos.Ticket
	4,  // 12: smartcore.bos.ServiceTicketInfo.DescribeTicket:output_type -> smartcore.bos.TicketSupport
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_service_ticket_proto_init() }
func file_service_ticket_proto_init() {
	if File_service_ticket_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_ticket_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Ticket); i {
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
		file_service_ticket_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateTicketRequest); i {
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
		file_service_ticket_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateTicketRequest); i {
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
		file_service_ticket_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*DescribeTicketRequest); i {
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
		file_service_ticket_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*TicketSupport); i {
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
		file_service_ticket_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Ticket_Classification); i {
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
		file_service_ticket_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*Ticket_Severity); i {
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
			RawDescriptor: file_service_ticket_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_service_ticket_proto_goTypes,
		DependencyIndexes: file_service_ticket_proto_depIdxs,
		MessageInfos:      file_service_ticket_proto_msgTypes,
	}.Build()
	File_service_ticket_proto = out.File
	file_service_ticket_proto_rawDesc = nil
	file_service_ticket_proto_goTypes = nil
	file_service_ticket_proto_depIdxs = nil
}
