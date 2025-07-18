// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: report.proto

package gen

import (
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

// ListReportsRequest contains a list of reports that are available for download.
type ListReportsRequest struct {
	state    protoimpl.MessageState `protogen:"open.v1"`
	Name     string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	// The maximum number of Reports to return.
	// The service may return fewer than this value.
	// If unspecified, at most 50 items will be returned.
	// The maximum value is 1000; values above 1000 will be coerced to 1000.
	PageSize int32 `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// A page token, received from a previous `ListReportsResponse` call.
	// Provide this to retrieve the subsequent page.
	PageToken     string `protobuf:"bytes,4,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListReportsRequest) Reset() {
	*x = ListReportsRequest{}
	mi := &file_report_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListReportsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListReportsRequest) ProtoMessage() {}

func (x *ListReportsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_report_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListReportsRequest.ProtoReflect.Descriptor instead.
func (*ListReportsRequest) Descriptor() ([]byte, []int) {
	return file_report_proto_rawDescGZIP(), []int{0}
}

func (x *ListReportsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListReportsRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

func (x *ListReportsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListReportsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

// ListReportsResponse contains a list of reports.
type ListReportsResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The list of reports.
	Reports []*Report `protobuf:"bytes,1,rep,name=reports,proto3" json:"reports,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page.
	// If this field is omitted, there are no subsequent pages.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	// If non-zero this is the total number of Reports after filtering is applied.
	// This may be an estimate.
	TotalSize     int32 `protobuf:"varint,3,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListReportsResponse) Reset() {
	*x = ListReportsResponse{}
	mi := &file_report_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListReportsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListReportsResponse) ProtoMessage() {}

func (x *ListReportsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_report_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListReportsResponse.ProtoReflect.Descriptor instead.
func (*ListReportsResponse) Descriptor() ([]byte, []int) {
	return file_report_proto_rawDescGZIP(), []int{1}
}

func (x *ListReportsResponse) GetReports() []*Report {
	if x != nil {
		return x.Reports
	}
	return nil
}

func (x *ListReportsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

func (x *ListReportsResponse) GetTotalSize() int32 {
	if x != nil {
		return x.TotalSize
	}
	return 0
}

// Report represents a report in the system.
type Report struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique identifier of the report. Can be the filename or a UUID.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The title of the report.
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	// The description of the report.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// The timestamp when the report was created.
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// The media type of the report, e.g., "application/pdf", "text/csv".
	MediaType     string `protobuf:"bytes,5,opt,name=media_type,json=mediaType,proto3" json:"media_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Report) Reset() {
	*x = Report{}
	mi := &file_report_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Report) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Report) ProtoMessage() {}

func (x *Report) ProtoReflect() protoreflect.Message {
	mi := &file_report_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Report.ProtoReflect.Descriptor instead.
func (*Report) Descriptor() ([]byte, []int) {
	return file_report_proto_rawDescGZIP(), []int{2}
}

func (x *Report) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Report) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Report) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Report) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Report) GetMediaType() string {
	if x != nil {
		return x.MediaType
	}
	return ""
}

// GetDownloadReportUrlRequest is used to request a specific report by its ID.
type GetDownloadReportUrlRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Name  string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The unique identifier of the report to retrieve.
	Id            string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDownloadReportUrlRequest) Reset() {
	*x = GetDownloadReportUrlRequest{}
	mi := &file_report_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDownloadReportUrlRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDownloadReportUrlRequest) ProtoMessage() {}

func (x *GetDownloadReportUrlRequest) ProtoReflect() protoreflect.Message {
	mi := &file_report_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDownloadReportUrlRequest.ProtoReflect.Descriptor instead.
func (*GetDownloadReportUrlRequest) Descriptor() ([]byte, []int) {
	return file_report_proto_rawDescGZIP(), []int{3}
}

func (x *GetDownloadReportUrlRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetDownloadReportUrlRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DownloadReportUrl struct {
	state    protoimpl.MessageState `protogen:"open.v1"`
	Url      string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Filename string                 `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
	// The media type of the report, e.g., "application/pdf", "text/csv".
	MediaType string `protobuf:"bytes,3,opt,name=media_type,json=mediaType,proto3" json:"media_type,omitempty"`
	// The latest time the url will be valid for, you will not be able to use the url after this time.
	ExpireAfterTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=expire_after_time,json=expireAfterTime,proto3" json:"expire_after_time,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *DownloadReportUrl) Reset() {
	*x = DownloadReportUrl{}
	mi := &file_report_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DownloadReportUrl) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadReportUrl) ProtoMessage() {}

func (x *DownloadReportUrl) ProtoReflect() protoreflect.Message {
	mi := &file_report_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadReportUrl.ProtoReflect.Descriptor instead.
func (*DownloadReportUrl) Descriptor() ([]byte, []int) {
	return file_report_proto_rawDescGZIP(), []int{4}
}

func (x *DownloadReportUrl) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *DownloadReportUrl) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *DownloadReportUrl) GetMediaType() string {
	if x != nil {
		return x.MediaType
	}
	return ""
}

func (x *DownloadReportUrl) GetExpireAfterTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpireAfterTime
	}
	return nil
}

var File_report_proto protoreflect.FileDescriptor

const file_report_proto_rawDesc = "" +
	"\n" +
	"\freport.proto\x12\rsmartcore.bos\x1a google/protobuf/field_mask.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\x9d\x01\n" +
	"\x12ListReportsRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x127\n" +
	"\tread_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\breadMask\x12\x1b\n" +
	"\tpage_size\x18\x03 \x01(\x05R\bpageSize\x12\x1d\n" +
	"\n" +
	"page_token\x18\x04 \x01(\tR\tpageToken\"\x8d\x01\n" +
	"\x13ListReportsResponse\x12/\n" +
	"\areports\x18\x01 \x03(\v2\x15.smartcore.bos.ReportR\areports\x12&\n" +
	"\x0fnext_page_token\x18\x02 \x01(\tR\rnextPageToken\x12\x1d\n" +
	"\n" +
	"total_size\x18\x03 \x01(\x05R\ttotalSize\"\xac\x01\n" +
	"\x06Report\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12;\n" +
	"\vcreate_time\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"createTime\x12\x1d\n" +
	"\n" +
	"media_type\x18\x05 \x01(\tR\tmediaType\"A\n" +
	"\x1bGetDownloadReportUrlRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\tR\x02id\"\xa8\x01\n" +
	"\x11DownloadReportUrl\x12\x10\n" +
	"\x03url\x18\x01 \x01(\tR\x03url\x12\x1a\n" +
	"\bfilename\x18\x02 \x01(\tR\bfilename\x12\x1d\n" +
	"\n" +
	"media_type\x18\x03 \x01(\tR\tmediaType\x12F\n" +
	"\x11expire_after_time\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\x0fexpireAfterTime2\xc7\x01\n" +
	"\tReportApi\x12T\n" +
	"\vListReports\x12!.smartcore.bos.ListReportsRequest\x1a\".smartcore.bos.ListReportsResponse\x12d\n" +
	"\x14GetDownloadReportUrl\x12*.smartcore.bos.GetDownloadReportUrlRequest\x1a .smartcore.bos.DownloadReportUrlB%Z#github.com/vanti-dev/sc-bos/pkg/genb\x06proto3"

var (
	file_report_proto_rawDescOnce sync.Once
	file_report_proto_rawDescData []byte
)

func file_report_proto_rawDescGZIP() []byte {
	file_report_proto_rawDescOnce.Do(func() {
		file_report_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_report_proto_rawDesc), len(file_report_proto_rawDesc)))
	})
	return file_report_proto_rawDescData
}

var file_report_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_report_proto_goTypes = []any{
	(*ListReportsRequest)(nil),          // 0: smartcore.bos.ListReportsRequest
	(*ListReportsResponse)(nil),         // 1: smartcore.bos.ListReportsResponse
	(*Report)(nil),                      // 2: smartcore.bos.Report
	(*GetDownloadReportUrlRequest)(nil), // 3: smartcore.bos.GetDownloadReportUrlRequest
	(*DownloadReportUrl)(nil),           // 4: smartcore.bos.DownloadReportUrl
	(*fieldmaskpb.FieldMask)(nil),       // 5: google.protobuf.FieldMask
	(*timestamppb.Timestamp)(nil),       // 6: google.protobuf.Timestamp
}
var file_report_proto_depIdxs = []int32{
	5, // 0: smartcore.bos.ListReportsRequest.read_mask:type_name -> google.protobuf.FieldMask
	2, // 1: smartcore.bos.ListReportsResponse.reports:type_name -> smartcore.bos.Report
	6, // 2: smartcore.bos.Report.create_time:type_name -> google.protobuf.Timestamp
	6, // 3: smartcore.bos.DownloadReportUrl.expire_after_time:type_name -> google.protobuf.Timestamp
	0, // 4: smartcore.bos.ReportApi.ListReports:input_type -> smartcore.bos.ListReportsRequest
	3, // 5: smartcore.bos.ReportApi.GetDownloadReportUrl:input_type -> smartcore.bos.GetDownloadReportUrlRequest
	1, // 6: smartcore.bos.ReportApi.ListReports:output_type -> smartcore.bos.ListReportsResponse
	4, // 7: smartcore.bos.ReportApi.GetDownloadReportUrl:output_type -> smartcore.bos.DownloadReportUrl
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_report_proto_init() }
func file_report_proto_init() {
	if File_report_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_report_proto_rawDesc), len(file_report_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_report_proto_goTypes,
		DependencyIndexes: file_report_proto_depIdxs,
		MessageInfos:      file_report_proto_msgTypes,
	}.Build()
	File_report_proto = out.File
	file_report_proto_goTypes = nil
	file_report_proto_depIdxs = nil
}
