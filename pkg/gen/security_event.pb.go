// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: security_event.proto

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

type SecurityEvent_State int32

const (
	// The state of the security event is unknown.
	SecurityEvent_STATE_UNKNOWN SecurityEvent_State = 0
	// The security event has not been acknowledged.
	SecurityEvent_UNACKNOWLEDGED SecurityEvent_State = 1
	// The security event has been acknowledged.
	SecurityEvent_ACKNOWLEDGED SecurityEvent_State = 2
	// The security event has been resolved.
	SecurityEvent_RESOLVED SecurityEvent_State = 3
)

// Enum value maps for SecurityEvent_State.
var (
	SecurityEvent_State_name = map[int32]string{
		0: "STATE_UNKNOWN",
		1: "UNACKNOWLEDGED",
		2: "ACKNOWLEDGED",
		3: "RESOLVED",
	}
	SecurityEvent_State_value = map[string]int32{
		"STATE_UNKNOWN":  0,
		"UNACKNOWLEDGED": 1,
		"ACKNOWLEDGED":   2,
		"RESOLVED":       3,
	}
)

func (x SecurityEvent_State) Enum() *SecurityEvent_State {
	p := new(SecurityEvent_State)
	*p = x
	return p
}

func (x SecurityEvent_State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SecurityEvent_State) Descriptor() protoreflect.EnumDescriptor {
	return file_security_event_proto_enumTypes[0].Descriptor()
}

func (SecurityEvent_State) Type() protoreflect.EnumType {
	return &file_security_event_proto_enumTypes[0]
}

func (x SecurityEvent_State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SecurityEvent_State.Descriptor instead.
func (SecurityEvent_State) EnumDescriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{0, 0}
}

type SecurityEvent_EventType int32

const (
	// The event type is unknown
	SecurityEvent_EVENT_TYPE_UNKNOWN SecurityEvent_EventType = 0
	// A security device has been tampered with
	SecurityEvent_TAMPER SecurityEvent_EventType = 1
	// A tamper detection has returned to normal
	SecurityEvent_TAMPER_CLEAR SecurityEvent_EventType = 2
	// A device has stopped responding
	SecurityEvent_DEVICE_OFFLINE SecurityEvent_EventType = 3
	// A card has reported an error
	SecurityEvent_CARD_ERROR SecurityEvent_EventType = 4
	// A maintenance warning has occurred
	SecurityEvent_MAINTENANCE_WARNING SecurityEvent_EventType = 5
	// A maintenance error has occurred
	SecurityEvent_MAINTENANCE_ERROR SecurityEvent_EventType = 6
	// The state of an alarm zone has changed
	SecurityEvent_ALARM_ZONE_STATE_CHANGE SecurityEvent_EventType = 7
	// An incorrect pin has been entered
	SecurityEvent_INCORRECT_PIN SecurityEvent_EventType = 8
	// Access has been denied
	SecurityEvent_ACCESS_DENIED SecurityEvent_EventType = 9
	// Access has been granted
	SecurityEvent_ACCESS_GRANTED SecurityEvent_EventType = 10
	// The source experienced duress
	SecurityEvent_DURESS SecurityEvent_EventType = 11
	// A card event has occurred
	SecurityEvent_CARD_EVENT SecurityEvent_EventType = 12
	// A door status has been reported
	SecurityEvent_DOOR_STATUS SecurityEvent_EventType = 13
	// A door has been open too long
	SecurityEvent_DOOR_OPEN_TOO_LONG SecurityEvent_EventType = 14
	// A door has been forced open
	SecurityEvent_DOOR_FORCED_OPEN SecurityEvent_EventType = 15
	// A door has not been locked
	SecurityEvent_DOOR_NOT_LOCKED SecurityEvent_EventType = 16
	// A power failure has occurred
	SecurityEvent_POWER_FAILURE SecurityEvent_EventType = 17
	// An invalid logon attempt has occurred
	SecurityEvent_INVALID_LOGON_ATTEMPT SecurityEvent_EventType = 18
	// A network attack has been detected
	SecurityEvent_NETWORK_ATTACK SecurityEvent_EventType = 19
	// Locker status
	SecurityEvent_LOCKER_STATUS SecurityEvent_EventType = 20
	// A locker has been open too long
	SecurityEvent_LOCKER_OPEN_TOO_LONG SecurityEvent_EventType = 21
	// A locker has been forced open
	SecurityEvent_LOCKER_FORCED_OPEN SecurityEvent_EventType = 22
	// A locker has not been locked
	SecurityEvent_LOCKER_NOT_LOCKED SecurityEvent_EventType = 23
)

// Enum value maps for SecurityEvent_EventType.
var (
	SecurityEvent_EventType_name = map[int32]string{
		0:  "EVENT_TYPE_UNKNOWN",
		1:  "TAMPER",
		2:  "TAMPER_CLEAR",
		3:  "DEVICE_OFFLINE",
		4:  "CARD_ERROR",
		5:  "MAINTENANCE_WARNING",
		6:  "MAINTENANCE_ERROR",
		7:  "ALARM_ZONE_STATE_CHANGE",
		8:  "INCORRECT_PIN",
		9:  "ACCESS_DENIED",
		10: "ACCESS_GRANTED",
		11: "DURESS",
		12: "CARD_EVENT",
		13: "DOOR_STATUS",
		14: "DOOR_OPEN_TOO_LONG",
		15: "DOOR_FORCED_OPEN",
		16: "DOOR_NOT_LOCKED",
		17: "POWER_FAILURE",
		18: "INVALID_LOGON_ATTEMPT",
		19: "NETWORK_ATTACK",
		20: "LOCKER_STATUS",
		21: "LOCKER_OPEN_TOO_LONG",
		22: "LOCKER_FORCED_OPEN",
		23: "LOCKER_NOT_LOCKED",
	}
	SecurityEvent_EventType_value = map[string]int32{
		"EVENT_TYPE_UNKNOWN":      0,
		"TAMPER":                  1,
		"TAMPER_CLEAR":            2,
		"DEVICE_OFFLINE":          3,
		"CARD_ERROR":              4,
		"MAINTENANCE_WARNING":     5,
		"MAINTENANCE_ERROR":       6,
		"ALARM_ZONE_STATE_CHANGE": 7,
		"INCORRECT_PIN":           8,
		"ACCESS_DENIED":           9,
		"ACCESS_GRANTED":          10,
		"DURESS":                  11,
		"CARD_EVENT":              12,
		"DOOR_STATUS":             13,
		"DOOR_OPEN_TOO_LONG":      14,
		"DOOR_FORCED_OPEN":        15,
		"DOOR_NOT_LOCKED":         16,
		"POWER_FAILURE":           17,
		"INVALID_LOGON_ATTEMPT":   18,
		"NETWORK_ATTACK":          19,
		"LOCKER_STATUS":           20,
		"LOCKER_OPEN_TOO_LONG":    21,
		"LOCKER_FORCED_OPEN":      22,
		"LOCKER_NOT_LOCKED":       23,
	}
)

func (x SecurityEvent_EventType) Enum() *SecurityEvent_EventType {
	p := new(SecurityEvent_EventType)
	*p = x
	return p
}

func (x SecurityEvent_EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SecurityEvent_EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_security_event_proto_enumTypes[1].Descriptor()
}

func (SecurityEvent_EventType) Type() protoreflect.EnumType {
	return &file_security_event_proto_enumTypes[1]
}

func (x SecurityEvent_EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SecurityEvent_EventType.Descriptor instead.
func (SecurityEvent_EventType) EnumDescriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{0, 1}
}

// SecurityEvent describes a security event that has occurred.
// At a minimum this should define the time the event occurred, a description of the event
// and a unique ID for the event, typically derived from the originating system.
// Ideally, this will contain all the relevant information we know about the event.
type SecurityEvent struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The time the security event occurred.
	SecurityEventTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=security_event_time,json=securityEventTime,proto3" json:"security_event_time,omitempty"`
	// A description of the security event.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// The ID of the event in the source system.
	Id string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	// The source the security event originated from, typically a door or other access point.
	Source *SecurityEvent_Source `protobuf:"bytes,4,opt,name=source,proto3" json:"source,omitempty"`
	// Optional. What actor caused this event to be triggered, an object, like a car, or a person or animal.
	// Sometimes, the event will not require an actor e.g. a door has been open too long.
	// Other times, an actor will cause a security event to occur at a source, e.g. a person has entered a door.
	Actor *Actor `protobuf:"bytes,5,opt,name=actor,proto3" json:"actor,omitempty"`
	// Optional. The state of the security event, unacknowledged, acknowledged etc.
	State SecurityEvent_State `protobuf:"varint,6,opt,name=state,proto3,enum=smartcore.bos.SecurityEvent_State" json:"state,omitempty"`
	// Optional. The priority of the security event
	Priority int32 `protobuf:"varint,7,opt,name=priority,proto3" json:"priority,omitempty"`
	// Optional. The type of the security event
	EventType     SecurityEvent_EventType `protobuf:"varint,8,opt,name=event_type,json=eventType,proto3,enum=smartcore.bos.SecurityEvent_EventType" json:"event_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SecurityEvent) Reset() {
	*x = SecurityEvent{}
	mi := &file_security_event_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SecurityEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecurityEvent) ProtoMessage() {}

func (x *SecurityEvent) ProtoReflect() protoreflect.Message {
	mi := &file_security_event_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecurityEvent.ProtoReflect.Descriptor instead.
func (*SecurityEvent) Descriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{0}
}

func (x *SecurityEvent) GetSecurityEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.SecurityEventTime
	}
	return nil
}

func (x *SecurityEvent) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SecurityEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SecurityEvent) GetSource() *SecurityEvent_Source {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *SecurityEvent) GetActor() *Actor {
	if x != nil {
		return x.Actor
	}
	return nil
}

func (x *SecurityEvent) GetState() SecurityEvent_State {
	if x != nil {
		return x.State
	}
	return SecurityEvent_STATE_UNKNOWN
}

func (x *SecurityEvent) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *SecurityEvent) GetEventType() SecurityEvent_EventType {
	if x != nil {
		return x.EventType
	}
	return SecurityEvent_EVENT_TYPE_UNKNOWN
}

type ListSecurityEventsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The name of the device to get security events for.
	Name     string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	// The maximum number of SecurityEvents to return.
	// The service may return fewer than this value.
	// If unspecified, at most 50 items will be returned.
	// The maximum value is 1000; values above 1000 will be coerced to 1000.
	PageSize int32 `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// A page token, received from a previous `ListSecurityEventsResponse` call.
	// Provide this to retrieve the subsequent page.
	PageToken     string `protobuf:"bytes,4,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSecurityEventsRequest) Reset() {
	*x = ListSecurityEventsRequest{}
	mi := &file_security_event_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSecurityEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSecurityEventsRequest) ProtoMessage() {}

func (x *ListSecurityEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_security_event_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSecurityEventsRequest.ProtoReflect.Descriptor instead.
func (*ListSecurityEventsRequest) Descriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{1}
}

func (x *ListSecurityEventsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListSecurityEventsRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

func (x *ListSecurityEventsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListSecurityEventsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type ListSecurityEventsResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	SecurityEvents []*SecurityEvent       `protobuf:"bytes,1,rep,name=security_events,json=securityEvents,proto3" json:"security_events,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page.
	// If this field is omitted, there are no subsequent pages.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	// If non-zero this is the total number of alerts after filtering is applied.
	// This may be an estimate.
	TotalSize     int32 `protobuf:"varint,3,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSecurityEventsResponse) Reset() {
	*x = ListSecurityEventsResponse{}
	mi := &file_security_event_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSecurityEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSecurityEventsResponse) ProtoMessage() {}

func (x *ListSecurityEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_security_event_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSecurityEventsResponse.ProtoReflect.Descriptor instead.
func (*ListSecurityEventsResponse) Descriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{2}
}

func (x *ListSecurityEventsResponse) GetSecurityEvents() []*SecurityEvent {
	if x != nil {
		return x.SecurityEvents
	}
	return nil
}

func (x *ListSecurityEventsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

func (x *ListSecurityEventsResponse) GetTotalSize() int32 {
	if x != nil {
		return x.TotalSize
	}
	return 0
}

type PullSecurityEventsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask      *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	UpdatesOnly   bool                   `protobuf:"varint,3,opt,name=updates_only,json=updatesOnly,proto3" json:"updates_only,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullSecurityEventsRequest) Reset() {
	*x = PullSecurityEventsRequest{}
	mi := &file_security_event_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullSecurityEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullSecurityEventsRequest) ProtoMessage() {}

func (x *PullSecurityEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_security_event_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullSecurityEventsRequest.ProtoReflect.Descriptor instead.
func (*PullSecurityEventsRequest) Descriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{3}
}

func (x *PullSecurityEventsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullSecurityEventsRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

func (x *PullSecurityEventsRequest) GetUpdatesOnly() bool {
	if x != nil {
		return x.UpdatesOnly
	}
	return false
}

type PullSecurityEventsResponse struct {
	state         protoimpl.MessageState               `protogen:"open.v1"`
	Changes       []*PullSecurityEventsResponse_Change `protobuf:"bytes,1,rep,name=changes,proto3" json:"changes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullSecurityEventsResponse) Reset() {
	*x = PullSecurityEventsResponse{}
	mi := &file_security_event_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullSecurityEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullSecurityEventsResponse) ProtoMessage() {}

func (x *PullSecurityEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_security_event_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullSecurityEventsResponse.ProtoReflect.Descriptor instead.
func (*PullSecurityEventsResponse) Descriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{4}
}

func (x *PullSecurityEventsResponse) GetChanges() []*PullSecurityEventsResponse_Change {
	if x != nil {
		return x.Changes
	}
	return nil
}

// A source is the device, cardholder or any other object where the security event originated from.
type SecurityEvent_Source struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ID. The unique identifier of the source.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name. The human readable name of the source.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Subsystem. The subsystem this event originated from, access control, cctv etc.
	Subsystem string `protobuf:"bytes,3,opt,name=subsystem,proto3" json:"subsystem,omitempty"`
	// Floor. The floor the event originated from.
	Floor string `protobuf:"bytes,4,opt,name=floor,proto3" json:"floor,omitempty"`
	// Zone. The zone the event originated from.
	Zone          string `protobuf:"bytes,5,opt,name=zone,proto3" json:"zone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SecurityEvent_Source) Reset() {
	*x = SecurityEvent_Source{}
	mi := &file_security_event_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SecurityEvent_Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecurityEvent_Source) ProtoMessage() {}

func (x *SecurityEvent_Source) ProtoReflect() protoreflect.Message {
	mi := &file_security_event_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecurityEvent_Source.ProtoReflect.Descriptor instead.
func (*SecurityEvent_Source) Descriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{0, 0}
}

func (x *SecurityEvent_Source) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SecurityEvent_Source) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SecurityEvent_Source) GetSubsystem() string {
	if x != nil {
		return x.Subsystem
	}
	return ""
}

func (x *SecurityEvent_Source) GetFloor() string {
	if x != nil {
		return x.Floor
	}
	return ""
}

func (x *SecurityEvent_Source) GetZone() string {
	if x != nil {
		return x.Zone
	}
	return ""
}

type PullSecurityEventsResponse_Change struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ChangeTime    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=change_time,json=changeTime,proto3" json:"change_time,omitempty"`
	NewValue      *SecurityEvent         `protobuf:"bytes,3,opt,name=new_value,json=newValue,proto3" json:"new_value,omitempty"`
	OldValue      *SecurityEvent         `protobuf:"bytes,4,opt,name=old_value,json=oldValue,proto3" json:"old_value,omitempty"`
	Type          types.ChangeType       `protobuf:"varint,5,opt,name=type,proto3,enum=smartcore.types.ChangeType" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullSecurityEventsResponse_Change) Reset() {
	*x = PullSecurityEventsResponse_Change{}
	mi := &file_security_event_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullSecurityEventsResponse_Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullSecurityEventsResponse_Change) ProtoMessage() {}

func (x *PullSecurityEventsResponse_Change) ProtoReflect() protoreflect.Message {
	mi := &file_security_event_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullSecurityEventsResponse_Change.ProtoReflect.Descriptor instead.
func (*PullSecurityEventsResponse_Change) Descriptor() ([]byte, []int) {
	return file_security_event_proto_rawDescGZIP(), []int{4, 0}
}

func (x *PullSecurityEventsResponse_Change) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullSecurityEventsResponse_Change) GetChangeTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ChangeTime
	}
	return nil
}

func (x *PullSecurityEventsResponse_Change) GetNewValue() *SecurityEvent {
	if x != nil {
		return x.NewValue
	}
	return nil
}

func (x *PullSecurityEventsResponse_Change) GetOldValue() *SecurityEvent {
	if x != nil {
		return x.OldValue
	}
	return nil
}

func (x *PullSecurityEventsResponse_Change) GetType() types.ChangeType {
	if x != nil {
		return x.Type
	}
	return types.ChangeType(0)
}

var File_security_event_proto protoreflect.FileDescriptor

const file_security_event_proto_rawDesc = "" +
	"\n" +
	"\x14security_event.proto\x12\rsmartcore.bos\x1a google/protobuf/field_mask.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x12types/change.proto\x1a\vactor.proto\"\xd6\b\n" +
	"\rSecurityEvent\x12J\n" +
	"\x13security_event_time\x18\x01 \x01(\v2\x1a.google.protobuf.TimestampR\x11securityEventTime\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12\x0e\n" +
	"\x02id\x18\x03 \x01(\tR\x02id\x12;\n" +
	"\x06source\x18\x04 \x01(\v2#.smartcore.bos.SecurityEvent.SourceR\x06source\x12*\n" +
	"\x05actor\x18\x05 \x01(\v2\x14.smartcore.bos.ActorR\x05actor\x128\n" +
	"\x05state\x18\x06 \x01(\x0e2\".smartcore.bos.SecurityEvent.StateR\x05state\x12\x1a\n" +
	"\bpriority\x18\a \x01(\x05R\bpriority\x12E\n" +
	"\n" +
	"event_type\x18\b \x01(\x0e2&.smartcore.bos.SecurityEvent.EventTypeR\teventType\x1at\n" +
	"\x06Source\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x1c\n" +
	"\tsubsystem\x18\x03 \x01(\tR\tsubsystem\x12\x14\n" +
	"\x05floor\x18\x04 \x01(\tR\x05floor\x12\x12\n" +
	"\x04zone\x18\x05 \x01(\tR\x04zone\"N\n" +
	"\x05State\x12\x11\n" +
	"\rSTATE_UNKNOWN\x10\x00\x12\x12\n" +
	"\x0eUNACKNOWLEDGED\x10\x01\x12\x10\n" +
	"\fACKNOWLEDGED\x10\x02\x12\f\n" +
	"\bRESOLVED\x10\x03\"\xfa\x03\n" +
	"\tEventType\x12\x16\n" +
	"\x12EVENT_TYPE_UNKNOWN\x10\x00\x12\n" +
	"\n" +
	"\x06TAMPER\x10\x01\x12\x10\n" +
	"\fTAMPER_CLEAR\x10\x02\x12\x12\n" +
	"\x0eDEVICE_OFFLINE\x10\x03\x12\x0e\n" +
	"\n" +
	"CARD_ERROR\x10\x04\x12\x17\n" +
	"\x13MAINTENANCE_WARNING\x10\x05\x12\x15\n" +
	"\x11MAINTENANCE_ERROR\x10\x06\x12\x1b\n" +
	"\x17ALARM_ZONE_STATE_CHANGE\x10\a\x12\x11\n" +
	"\rINCORRECT_PIN\x10\b\x12\x11\n" +
	"\rACCESS_DENIED\x10\t\x12\x12\n" +
	"\x0eACCESS_GRANTED\x10\n" +
	"\x12\n" +
	"\n" +
	"\x06DURESS\x10\v\x12\x0e\n" +
	"\n" +
	"CARD_EVENT\x10\f\x12\x0f\n" +
	"\vDOOR_STATUS\x10\r\x12\x16\n" +
	"\x12DOOR_OPEN_TOO_LONG\x10\x0e\x12\x14\n" +
	"\x10DOOR_FORCED_OPEN\x10\x0f\x12\x13\n" +
	"\x0fDOOR_NOT_LOCKED\x10\x10\x12\x11\n" +
	"\rPOWER_FAILURE\x10\x11\x12\x19\n" +
	"\x15INVALID_LOGON_ATTEMPT\x10\x12\x12\x12\n" +
	"\x0eNETWORK_ATTACK\x10\x13\x12\x11\n" +
	"\rLOCKER_STATUS\x10\x14\x12\x18\n" +
	"\x14LOCKER_OPEN_TOO_LONG\x10\x15\x12\x16\n" +
	"\x12LOCKER_FORCED_OPEN\x10\x16\x12\x15\n" +
	"\x11LOCKER_NOT_LOCKED\x10\x17\"\xa4\x01\n" +
	"\x19ListSecurityEventsRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x127\n" +
	"\tread_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\breadMask\x12\x1b\n" +
	"\tpage_size\x18\x03 \x01(\x05R\bpageSize\x12\x1d\n" +
	"\n" +
	"page_token\x18\x04 \x01(\tR\tpageToken\"\xaa\x01\n" +
	"\x1aListSecurityEventsResponse\x12E\n" +
	"\x0fsecurity_events\x18\x01 \x03(\v2\x1c.smartcore.bos.SecurityEventR\x0esecurityEvents\x12&\n" +
	"\x0fnext_page_token\x18\x02 \x01(\tR\rnextPageToken\x12\x1d\n" +
	"\n" +
	"total_size\x18\x03 \x01(\x05R\ttotalSize\"\x8b\x01\n" +
	"\x19PullSecurityEventsRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x127\n" +
	"\tread_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\breadMask\x12!\n" +
	"\fupdates_only\x18\x03 \x01(\bR\vupdatesOnly\"\xeb\x02\n" +
	"\x1aPullSecurityEventsResponse\x12J\n" +
	"\achanges\x18\x01 \x03(\v20.smartcore.bos.PullSecurityEventsResponse.ChangeR\achanges\x1a\x80\x02\n" +
	"\x06Change\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12;\n" +
	"\vchange_time\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"changeTime\x129\n" +
	"\tnew_value\x18\x03 \x01(\v2\x1c.smartcore.bos.SecurityEventR\bnewValue\x129\n" +
	"\told_value\x18\x04 \x01(\v2\x1c.smartcore.bos.SecurityEventR\boldValue\x12/\n" +
	"\x04type\x18\x05 \x01(\x0e2\x1b.smartcore.types.ChangeTypeR\x04type2\xee\x01\n" +
	"\x10SecurityEventApi\x12k\n" +
	"\x12ListSecurityEvents\x12(.smartcore.bos.ListSecurityEventsRequest\x1a).smartcore.bos.ListSecurityEventsResponse\"\x00\x12m\n" +
	"\x12PullSecurityEvents\x12(.smartcore.bos.PullSecurityEventsRequest\x1a).smartcore.bos.PullSecurityEventsResponse\"\x000\x01B%Z#github.com/vanti-dev/sc-bos/pkg/genb\x06proto3"

var (
	file_security_event_proto_rawDescOnce sync.Once
	file_security_event_proto_rawDescData []byte
)

func file_security_event_proto_rawDescGZIP() []byte {
	file_security_event_proto_rawDescOnce.Do(func() {
		file_security_event_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_security_event_proto_rawDesc), len(file_security_event_proto_rawDesc)))
	})
	return file_security_event_proto_rawDescData
}

var file_security_event_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_security_event_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_security_event_proto_goTypes = []any{
	(SecurityEvent_State)(0),                  // 0: smartcore.bos.SecurityEvent.State
	(SecurityEvent_EventType)(0),              // 1: smartcore.bos.SecurityEvent.EventType
	(*SecurityEvent)(nil),                     // 2: smartcore.bos.SecurityEvent
	(*ListSecurityEventsRequest)(nil),         // 3: smartcore.bos.ListSecurityEventsRequest
	(*ListSecurityEventsResponse)(nil),        // 4: smartcore.bos.ListSecurityEventsResponse
	(*PullSecurityEventsRequest)(nil),         // 5: smartcore.bos.PullSecurityEventsRequest
	(*PullSecurityEventsResponse)(nil),        // 6: smartcore.bos.PullSecurityEventsResponse
	(*SecurityEvent_Source)(nil),              // 7: smartcore.bos.SecurityEvent.Source
	(*PullSecurityEventsResponse_Change)(nil), // 8: smartcore.bos.PullSecurityEventsResponse.Change
	(*timestamppb.Timestamp)(nil),             // 9: google.protobuf.Timestamp
	(*Actor)(nil),                             // 10: smartcore.bos.Actor
	(*fieldmaskpb.FieldMask)(nil),             // 11: google.protobuf.FieldMask
	(types.ChangeType)(0),                     // 12: smartcore.types.ChangeType
}
var file_security_event_proto_depIdxs = []int32{
	9,  // 0: smartcore.bos.SecurityEvent.security_event_time:type_name -> google.protobuf.Timestamp
	7,  // 1: smartcore.bos.SecurityEvent.source:type_name -> smartcore.bos.SecurityEvent.Source
	10, // 2: smartcore.bos.SecurityEvent.actor:type_name -> smartcore.bos.Actor
	0,  // 3: smartcore.bos.SecurityEvent.state:type_name -> smartcore.bos.SecurityEvent.State
	1,  // 4: smartcore.bos.SecurityEvent.event_type:type_name -> smartcore.bos.SecurityEvent.EventType
	11, // 5: smartcore.bos.ListSecurityEventsRequest.read_mask:type_name -> google.protobuf.FieldMask
	2,  // 6: smartcore.bos.ListSecurityEventsResponse.security_events:type_name -> smartcore.bos.SecurityEvent
	11, // 7: smartcore.bos.PullSecurityEventsRequest.read_mask:type_name -> google.protobuf.FieldMask
	8,  // 8: smartcore.bos.PullSecurityEventsResponse.changes:type_name -> smartcore.bos.PullSecurityEventsResponse.Change
	9,  // 9: smartcore.bos.PullSecurityEventsResponse.Change.change_time:type_name -> google.protobuf.Timestamp
	2,  // 10: smartcore.bos.PullSecurityEventsResponse.Change.new_value:type_name -> smartcore.bos.SecurityEvent
	2,  // 11: smartcore.bos.PullSecurityEventsResponse.Change.old_value:type_name -> smartcore.bos.SecurityEvent
	12, // 12: smartcore.bos.PullSecurityEventsResponse.Change.type:type_name -> smartcore.types.ChangeType
	3,  // 13: smartcore.bos.SecurityEventApi.ListSecurityEvents:input_type -> smartcore.bos.ListSecurityEventsRequest
	5,  // 14: smartcore.bos.SecurityEventApi.PullSecurityEvents:input_type -> smartcore.bos.PullSecurityEventsRequest
	4,  // 15: smartcore.bos.SecurityEventApi.ListSecurityEvents:output_type -> smartcore.bos.ListSecurityEventsResponse
	6,  // 16: smartcore.bos.SecurityEventApi.PullSecurityEvents:output_type -> smartcore.bos.PullSecurityEventsResponse
	15, // [15:17] is the sub-list for method output_type
	13, // [13:15] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_security_event_proto_init() }
func file_security_event_proto_init() {
	if File_security_event_proto != nil {
		return
	}
	file_actor_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_security_event_proto_rawDesc), len(file_security_event_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_security_event_proto_goTypes,
		DependencyIndexes: file_security_event_proto_depIdxs,
		EnumInfos:         file_security_event_proto_enumTypes,
		MessageInfos:      file_security_event_proto_msgTypes,
	}.Build()
	File_security_event_proto = out.File
	file_security_event_proto_goTypes = nil
	file_security_event_proto_depIdxs = nil
}
