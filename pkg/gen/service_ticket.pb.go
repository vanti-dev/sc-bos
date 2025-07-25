// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: service_ticket.proto

package gen

import (
	types "github.com/smart-core-os/sc-api/go/types"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

// Ticket represents a service ticket in a third party system.
type Ticket struct {
	state protoimpl.MessageState `protogen:"open.v1"`
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
	ExternalUrl   string `protobuf:"bytes,7,opt,name=external_url,json=externalUrl,proto3" json:"external_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Ticket) Reset() {
	*x = Ticket{}
	mi := &file_service_ticket_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ticket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ticket) ProtoMessage() {}

func (x *Ticket) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[0]
	if x != nil {
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Ticket        *Ticket                `protobuf:"bytes,2,opt,name=ticket,proto3" json:"ticket,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateTicketRequest) Reset() {
	*x = CreateTicketRequest{}
	mi := &file_service_ticket_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTicketRequest) ProtoMessage() {}

func (x *CreateTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[1]
	if x != nil {
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Ticket        *Ticket                `protobuf:"bytes,2,opt,name=ticket,proto3" json:"ticket,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateTicketRequest) Reset() {
	*x = UpdateTicketRequest{}
	mi := &file_service_ticket_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTicketRequest) ProtoMessage() {}

func (x *UpdateTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[2]
	if x != nil {
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DescribeTicketRequest) Reset() {
	*x = DescribeTicketRequest{}
	mi := &file_service_ticket_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DescribeTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeTicketRequest) ProtoMessage() {}

func (x *DescribeTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[3]
	if x != nil {
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
	state protoimpl.MessageState `protogen:"open.v1"`
	// How a named device supports read/write/pull apis
	ResourceSupport *types.ResourceSupport `protobuf:"bytes,1,opt,name=resource_support,json=resourceSupport,proto3" json:"resource_support,omitempty"`
	// The classifications supported by the implementing system.
	Classifications []*Ticket_Classification `protobuf:"bytes,2,rep,name=classifications,proto3" json:"classifications,omitempty"`
	// The severities supported by the implementing system.
	Severities    []*Ticket_Severity `protobuf:"bytes,3,rep,name=severities,proto3" json:"severities,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TicketSupport) Reset() {
	*x = TicketSupport{}
	mi := &file_service_ticket_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TicketSupport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TicketSupport) ProtoMessage() {}

func (x *TicketSupport) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[4]
	if x != nil {
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
	state protoimpl.MessageState `protogen:"open.v1"`
	// The title of the classification.
	// This is unique within the context of the implementing system.
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Optional. A more detailed description can be displayed to a user to help them decide the correct classification.
	Description   string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Ticket_Classification) Reset() {
	*x = Ticket_Classification{}
	mi := &file_service_ticket_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ticket_Classification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ticket_Classification) ProtoMessage() {}

func (x *Ticket_Classification) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[5]
	if x != nil {
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
	state protoimpl.MessageState `protogen:"open.v1"`
	// The title of the severity.
	// This is unique within the context of the implementing system.
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Optional. A more detailed description can be displayed to a user to help them decide the correct severity.
	Description   string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Ticket_Severity) Reset() {
	*x = Ticket_Severity{}
	mi := &file_service_ticket_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ticket_Severity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ticket_Severity) ProtoMessage() {}

func (x *Ticket_Severity) ProtoReflect() protoreflect.Message {
	mi := &file_service_ticket_proto_msgTypes[6]
	if x != nil {
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

const file_service_ticket_proto_rawDesc = "" +
	"\n" +
	"\x14service_ticket.proto\x12\rsmartcore.bos\x1a\x10types/info.proto\"\xb4\x03\n" +
	"\x06Ticket\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x18\n" +
	"\asummary\x18\x02 \x01(\tR\asummary\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12#\n" +
	"\rreporter_name\x18\x04 \x01(\tR\freporterName\x12L\n" +
	"\x0eclassification\x18\x05 \x01(\v2$.smartcore.bos.Ticket.ClassificationR\x0eclassification\x12:\n" +
	"\bseverity\x18\x06 \x01(\v2\x1e.smartcore.bos.Ticket.SeverityR\bseverity\x12!\n" +
	"\fexternal_url\x18\a \x01(\tR\vexternalUrl\x1aH\n" +
	"\x0eClassification\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x1aB\n" +
	"\bSeverity\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\"X\n" +
	"\x13CreateTicketRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12-\n" +
	"\x06ticket\x18\x02 \x01(\v2\x15.smartcore.bos.TicketR\x06ticket\"X\n" +
	"\x13UpdateTicketRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12-\n" +
	"\x06ticket\x18\x02 \x01(\v2\x15.smartcore.bos.TicketR\x06ticket\"+\n" +
	"\x15DescribeTicketRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"\xec\x01\n" +
	"\rTicketSupport\x12K\n" +
	"\x10resource_support\x18\x01 \x01(\v2 .smartcore.types.ResourceSupportR\x0fresourceSupport\x12N\n" +
	"\x0fclassifications\x18\x02 \x03(\v2$.smartcore.bos.Ticket.ClassificationR\x0fclassifications\x12>\n" +
	"\n" +
	"severities\x18\x03 \x03(\v2\x1e.smartcore.bos.Ticket.SeverityR\n" +
	"severities2\xac\x01\n" +
	"\x10ServiceTicketApi\x12K\n" +
	"\fCreateTicket\x12\".smartcore.bos.CreateTicketRequest\x1a\x15.smartcore.bos.Ticket\"\x00\x12K\n" +
	"\fUpdateTicket\x12\".smartcore.bos.UpdateTicketRequest\x1a\x15.smartcore.bos.Ticket\"\x002k\n" +
	"\x11ServiceTicketInfo\x12V\n" +
	"\x0eDescribeTicket\x12$.smartcore.bos.DescribeTicketRequest\x1a\x1c.smartcore.bos.TicketSupport\"\x00B%Z#github.com/vanti-dev/sc-bos/pkg/genb\x06proto3"

var (
	file_service_ticket_proto_rawDescOnce sync.Once
	file_service_ticket_proto_rawDescData []byte
)

func file_service_ticket_proto_rawDescGZIP() []byte {
	file_service_ticket_proto_rawDescOnce.Do(func() {
		file_service_ticket_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_service_ticket_proto_rawDesc), len(file_service_ticket_proto_rawDesc)))
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_service_ticket_proto_rawDesc), len(file_service_ticket_proto_rawDesc)),
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
	file_service_ticket_proto_goTypes = nil
	file_service_ticket_proto_depIdxs = nil
}
