// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: button.proto

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

// Instantaneous button state.
type ButtonState_Press int32

const (
	ButtonState_PRESS_UNSPECIFIED ButtonState_Press = 0
	// Button is in its neutral position.
	ButtonState_UNPRESSED ButtonState_Press = 1
	// Button is being pushed in.
	ButtonState_PRESSED ButtonState_Press = 2
)

// Enum value maps for ButtonState_Press.
var (
	ButtonState_Press_name = map[int32]string{
		0: "PRESS_UNSPECIFIED",
		1: "UNPRESSED",
		2: "PRESSED",
	}
	ButtonState_Press_value = map[string]int32{
		"PRESS_UNSPECIFIED": 0,
		"UNPRESSED":         1,
		"PRESSED":           2,
	}
)

func (x ButtonState_Press) Enum() *ButtonState_Press {
	p := new(ButtonState_Press)
	*p = x
	return p
}

func (x ButtonState_Press) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ButtonState_Press) Descriptor() protoreflect.EnumDescriptor {
	return file_button_proto_enumTypes[0].Descriptor()
}

func (ButtonState_Press) Type() protoreflect.EnumType {
	return &file_button_proto_enumTypes[0]
}

func (x ButtonState_Press) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ButtonState_Press.Descriptor instead.
func (ButtonState_Press) EnumDescriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{0, 0}
}

type ButtonState_Gesture_Kind int32

const (
	ButtonState_Gesture_KIND_UNSPECIFIED ButtonState_Gesture_Kind = 0
	// One or more short press-and-release actions.
	// Clicks in short succession may be fused into double-clicks, triple-clicks etc. - in this case, the number
	// of fused clicks is stored in the count field.
	// When clicks are fused in this way, the gesture will not appear at all until the final click has finished -
	// it's not possible for a single gesture to be first reported as a single click, and then modified to a double click.
	ButtonState_Gesture_CLICK ButtonState_Gesture_Kind = 1
	// Button is kept in the pressed state for an extended period.
	// Buttons may support repeat events, in which case the count will increment for each repeat event, keeping id
	// the same because it's part of the same gesture.
	// For HOLD gestures, the end_time is not set until the button has been released, allowing the client to determine
	// when the gesture has ended.
	ButtonState_Gesture_HOLD ButtonState_Gesture_Kind = 2
)

// Enum value maps for ButtonState_Gesture_Kind.
var (
	ButtonState_Gesture_Kind_name = map[int32]string{
		0: "KIND_UNSPECIFIED",
		1: "CLICK",
		2: "HOLD",
	}
	ButtonState_Gesture_Kind_value = map[string]int32{
		"KIND_UNSPECIFIED": 0,
		"CLICK":            1,
		"HOLD":             2,
	}
)

func (x ButtonState_Gesture_Kind) Enum() *ButtonState_Gesture_Kind {
	p := new(ButtonState_Gesture_Kind)
	*p = x
	return p
}

func (x ButtonState_Gesture_Kind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ButtonState_Gesture_Kind) Descriptor() protoreflect.EnumDescriptor {
	return file_button_proto_enumTypes[1].Descriptor()
}

func (ButtonState_Gesture_Kind) Type() protoreflect.EnumType {
	return &file_button_proto_enumTypes[1]
}

func (x ButtonState_Gesture_Kind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ButtonState_Gesture_Kind.Descriptor instead.
func (ButtonState_Gesture_Kind) EnumDescriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{0, 0, 0}
}

type ButtonState struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	State ButtonState_Press      `protobuf:"varint,1,opt,name=state,proto3,enum=smartcore.bos.ButtonState_Press" json:"state,omitempty"`
	// The time that state changed to its present value.
	StateChangeTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=state_change_time,json=stateChangeTime,proto3" json:"state_change_time,omitempty"`
	// The gesture that is currently in progress, or finished most recently.
	// May be absent, if there is no gesture recorded for this button.
	MostRecentGesture *ButtonState_Gesture `protobuf:"bytes,3,opt,name=most_recent_gesture,json=mostRecentGesture,proto3" json:"most_recent_gesture,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *ButtonState) Reset() {
	*x = ButtonState{}
	mi := &file_button_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ButtonState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ButtonState) ProtoMessage() {}

func (x *ButtonState) ProtoReflect() protoreflect.Message {
	mi := &file_button_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ButtonState.ProtoReflect.Descriptor instead.
func (*ButtonState) Descriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{0}
}

func (x *ButtonState) GetState() ButtonState_Press {
	if x != nil {
		return x.State
	}
	return ButtonState_PRESS_UNSPECIFIED
}

func (x *ButtonState) GetStateChangeTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StateChangeTime
	}
	return nil
}

func (x *ButtonState) GetMostRecentGesture() *ButtonState_Gesture {
	if x != nil {
		return x.MostRecentGesture
	}
	return nil
}

type GetButtonStateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask      *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetButtonStateRequest) Reset() {
	*x = GetButtonStateRequest{}
	mi := &file_button_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetButtonStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetButtonStateRequest) ProtoMessage() {}

func (x *GetButtonStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_button_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetButtonStateRequest.ProtoReflect.Descriptor instead.
func (*GetButtonStateRequest) Descriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{1}
}

func (x *GetButtonStateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetButtonStateRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

type PullButtonStateRequest struct {
	state    protoimpl.MessageState `protogen:"open.v1"`
	Name     string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ReadMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=read_mask,json=readMask,proto3" json:"read_mask,omitempty"`
	// By default, PullButtonState sends the initial ButtonState when the stream opens, followed by changes.
	// Setting updates_only true will disable this behaviour, sending only when the ButtonState changes.
	UpdatesOnly   bool `protobuf:"varint,3,opt,name=updates_only,json=updatesOnly,proto3" json:"updates_only,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullButtonStateRequest) Reset() {
	*x = PullButtonStateRequest{}
	mi := &file_button_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullButtonStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullButtonStateRequest) ProtoMessage() {}

func (x *PullButtonStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_button_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullButtonStateRequest.ProtoReflect.Descriptor instead.
func (*PullButtonStateRequest) Descriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{2}
}

func (x *PullButtonStateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullButtonStateRequest) GetReadMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.ReadMask
	}
	return nil
}

func (x *PullButtonStateRequest) GetUpdatesOnly() bool {
	if x != nil {
		return x.UpdatesOnly
	}
	return false
}

type PullButtonStateResponse struct {
	state         protoimpl.MessageState            `protogen:"open.v1"`
	Changes       []*PullButtonStateResponse_Change `protobuf:"bytes,1,rep,name=changes,proto3" json:"changes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullButtonStateResponse) Reset() {
	*x = PullButtonStateResponse{}
	mi := &file_button_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullButtonStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullButtonStateResponse) ProtoMessage() {}

func (x *PullButtonStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_button_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullButtonStateResponse.ProtoReflect.Descriptor instead.
func (*PullButtonStateResponse) Descriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{3}
}

func (x *PullButtonStateResponse) GetChanges() []*PullButtonStateResponse_Change {
	if x != nil {
		return x.Changes
	}
	return nil
}

type UpdateButtonStateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	UpdateMask    *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	ButtonState   *ButtonState           `protobuf:"bytes,3,opt,name=button_state,json=buttonState,proto3" json:"button_state,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateButtonStateRequest) Reset() {
	*x = UpdateButtonStateRequest{}
	mi := &file_button_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateButtonStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateButtonStateRequest) ProtoMessage() {}

func (x *UpdateButtonStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_button_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateButtonStateRequest.ProtoReflect.Descriptor instead.
func (*UpdateButtonStateRequest) Descriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateButtonStateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateButtonStateRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

func (x *UpdateButtonStateRequest) GetButtonState() *ButtonState {
	if x != nil {
		return x.ButtonState
	}
	return nil
}

// A representation of user intent, deduced from a pattern of button presses.
// The way that the device converts button presses into gestures is implementation-defined.
// There may be a delay between the button presses and the registration of a gesture.
type ButtonState_Gesture struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Opaque identifier changes each time a new gesture begins.
	// The gesture will remain in the ButtonState even when the client has already seen it; the client can use the id
	// to detect when a new gesture has begun.
	Id   string                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Kind ButtonState_Gesture_Kind `protobuf:"varint,2,opt,name=kind,proto3,enum=smartcore.bos.ButtonState_Gesture_Kind" json:"kind,omitempty"`
	// A counter for sub-events that occur within a single gesture. See the Kind for details of meaning.
	Count int32 `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
	// The time when the gesture was first recognised.
	StartTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// The time when the gesture was recognised as completed. For HOLD gestures, this remains unset until the button
	// is released.
	EndTime       *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ButtonState_Gesture) Reset() {
	*x = ButtonState_Gesture{}
	mi := &file_button_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ButtonState_Gesture) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ButtonState_Gesture) ProtoMessage() {}

func (x *ButtonState_Gesture) ProtoReflect() protoreflect.Message {
	mi := &file_button_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ButtonState_Gesture.ProtoReflect.Descriptor instead.
func (*ButtonState_Gesture) Descriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ButtonState_Gesture) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ButtonState_Gesture) GetKind() ButtonState_Gesture_Kind {
	if x != nil {
		return x.Kind
	}
	return ButtonState_Gesture_KIND_UNSPECIFIED
}

func (x *ButtonState_Gesture) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *ButtonState_Gesture) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *ButtonState_Gesture) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

type PullButtonStateResponse_Change struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ChangeTime    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=change_time,json=changeTime,proto3" json:"change_time,omitempty"`
	ButtonState   *ButtonState           `protobuf:"bytes,3,opt,name=button_state,json=buttonState,proto3" json:"button_state,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PullButtonStateResponse_Change) Reset() {
	*x = PullButtonStateResponse_Change{}
	mi := &file_button_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PullButtonStateResponse_Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullButtonStateResponse_Change) ProtoMessage() {}

func (x *PullButtonStateResponse_Change) ProtoReflect() protoreflect.Message {
	mi := &file_button_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullButtonStateResponse_Change.ProtoReflect.Descriptor instead.
func (*PullButtonStateResponse_Change) Descriptor() ([]byte, []int) {
	return file_button_proto_rawDescGZIP(), []int{3, 0}
}

func (x *PullButtonStateResponse_Change) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PullButtonStateResponse_Change) GetChangeTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ChangeTime
	}
	return nil
}

func (x *PullButtonStateResponse_Change) GetButtonState() *ButtonState {
	if x != nil {
		return x.ButtonState
	}
	return nil
}

var File_button_proto protoreflect.FileDescriptor

const file_button_proto_rawDesc = "" +
	"\n" +
	"\fbutton.proto\x12\rsmartcore.bos\x1a\x1fgoogle/protobuf/timestamp.proto\x1a google/protobuf/field_mask.proto\"\xb1\x04\n" +
	"\vButtonState\x126\n" +
	"\x05state\x18\x01 \x01(\x0e2 .smartcore.bos.ButtonState.PressR\x05state\x12F\n" +
	"\x11state_change_time\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x0fstateChangeTime\x12R\n" +
	"\x13most_recent_gesture\x18\x03 \x01(\v2\".smartcore.bos.ButtonState.GestureR\x11mostRecentGesture\x1a\x91\x02\n" +
	"\aGesture\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12;\n" +
	"\x04kind\x18\x02 \x01(\x0e2'.smartcore.bos.ButtonState.Gesture.KindR\x04kind\x12\x14\n" +
	"\x05count\x18\x03 \x01(\x05R\x05count\x129\n" +
	"\n" +
	"start_time\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\tstartTime\x125\n" +
	"\bend_time\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\aendTime\"1\n" +
	"\x04Kind\x12\x14\n" +
	"\x10KIND_UNSPECIFIED\x10\x00\x12\t\n" +
	"\x05CLICK\x10\x01\x12\b\n" +
	"\x04HOLD\x10\x02\":\n" +
	"\x05Press\x12\x15\n" +
	"\x11PRESS_UNSPECIFIED\x10\x00\x12\r\n" +
	"\tUNPRESSED\x10\x01\x12\v\n" +
	"\aPRESSED\x10\x02\"d\n" +
	"\x15GetButtonStateRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x127\n" +
	"\tread_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\breadMask\"\x88\x01\n" +
	"\x16PullButtonStateRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x127\n" +
	"\tread_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\breadMask\x12!\n" +
	"\fupdates_only\x18\x03 \x01(\bR\vupdatesOnly\"\xfd\x01\n" +
	"\x17PullButtonStateResponse\x12G\n" +
	"\achanges\x18\x01 \x03(\v2-.smartcore.bos.PullButtonStateResponse.ChangeR\achanges\x1a\x98\x01\n" +
	"\x06Change\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12;\n" +
	"\vchange_time\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"changeTime\x12=\n" +
	"\fbutton_state\x18\x03 \x01(\v2\x1a.smartcore.bos.ButtonStateR\vbuttonState\"\xaa\x01\n" +
	"\x18UpdateButtonStateRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12;\n" +
	"\vupdate_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\n" +
	"updateMask\x12=\n" +
	"\fbutton_state\x18\x03 \x01(\v2\x1a.smartcore.bos.ButtonStateR\vbuttonState2\x9d\x02\n" +
	"\tButtonApi\x12R\n" +
	"\x0eGetButtonState\x12$.smartcore.bos.GetButtonStateRequest\x1a\x1a.smartcore.bos.ButtonState\x12b\n" +
	"\x0fPullButtonState\x12%.smartcore.bos.PullButtonStateRequest\x1a&.smartcore.bos.PullButtonStateResponse0\x01\x12X\n" +
	"\x11UpdateButtonState\x12'.smartcore.bos.UpdateButtonStateRequest\x1a\x1a.smartcore.bos.ButtonStateB%Z#github.com/vanti-dev/sc-bos/pkg/genb\x06proto3"

var (
	file_button_proto_rawDescOnce sync.Once
	file_button_proto_rawDescData []byte
)

func file_button_proto_rawDescGZIP() []byte {
	file_button_proto_rawDescOnce.Do(func() {
		file_button_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_button_proto_rawDesc), len(file_button_proto_rawDesc)))
	})
	return file_button_proto_rawDescData
}

var file_button_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_button_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_button_proto_goTypes = []any{
	(ButtonState_Press)(0),                 // 0: smartcore.bos.ButtonState.Press
	(ButtonState_Gesture_Kind)(0),          // 1: smartcore.bos.ButtonState.Gesture.Kind
	(*ButtonState)(nil),                    // 2: smartcore.bos.ButtonState
	(*GetButtonStateRequest)(nil),          // 3: smartcore.bos.GetButtonStateRequest
	(*PullButtonStateRequest)(nil),         // 4: smartcore.bos.PullButtonStateRequest
	(*PullButtonStateResponse)(nil),        // 5: smartcore.bos.PullButtonStateResponse
	(*UpdateButtonStateRequest)(nil),       // 6: smartcore.bos.UpdateButtonStateRequest
	(*ButtonState_Gesture)(nil),            // 7: smartcore.bos.ButtonState.Gesture
	(*PullButtonStateResponse_Change)(nil), // 8: smartcore.bos.PullButtonStateResponse.Change
	(*timestamppb.Timestamp)(nil),          // 9: google.protobuf.Timestamp
	(*fieldmaskpb.FieldMask)(nil),          // 10: google.protobuf.FieldMask
}
var file_button_proto_depIdxs = []int32{
	0,  // 0: smartcore.bos.ButtonState.state:type_name -> smartcore.bos.ButtonState.Press
	9,  // 1: smartcore.bos.ButtonState.state_change_time:type_name -> google.protobuf.Timestamp
	7,  // 2: smartcore.bos.ButtonState.most_recent_gesture:type_name -> smartcore.bos.ButtonState.Gesture
	10, // 3: smartcore.bos.GetButtonStateRequest.read_mask:type_name -> google.protobuf.FieldMask
	10, // 4: smartcore.bos.PullButtonStateRequest.read_mask:type_name -> google.protobuf.FieldMask
	8,  // 5: smartcore.bos.PullButtonStateResponse.changes:type_name -> smartcore.bos.PullButtonStateResponse.Change
	10, // 6: smartcore.bos.UpdateButtonStateRequest.update_mask:type_name -> google.protobuf.FieldMask
	2,  // 7: smartcore.bos.UpdateButtonStateRequest.button_state:type_name -> smartcore.bos.ButtonState
	1,  // 8: smartcore.bos.ButtonState.Gesture.kind:type_name -> smartcore.bos.ButtonState.Gesture.Kind
	9,  // 9: smartcore.bos.ButtonState.Gesture.start_time:type_name -> google.protobuf.Timestamp
	9,  // 10: smartcore.bos.ButtonState.Gesture.end_time:type_name -> google.protobuf.Timestamp
	9,  // 11: smartcore.bos.PullButtonStateResponse.Change.change_time:type_name -> google.protobuf.Timestamp
	2,  // 12: smartcore.bos.PullButtonStateResponse.Change.button_state:type_name -> smartcore.bos.ButtonState
	3,  // 13: smartcore.bos.ButtonApi.GetButtonState:input_type -> smartcore.bos.GetButtonStateRequest
	4,  // 14: smartcore.bos.ButtonApi.PullButtonState:input_type -> smartcore.bos.PullButtonStateRequest
	6,  // 15: smartcore.bos.ButtonApi.UpdateButtonState:input_type -> smartcore.bos.UpdateButtonStateRequest
	2,  // 16: smartcore.bos.ButtonApi.GetButtonState:output_type -> smartcore.bos.ButtonState
	5,  // 17: smartcore.bos.ButtonApi.PullButtonState:output_type -> smartcore.bos.PullButtonStateResponse
	2,  // 18: smartcore.bos.ButtonApi.UpdateButtonState:output_type -> smartcore.bos.ButtonState
	16, // [16:19] is the sub-list for method output_type
	13, // [13:16] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_button_proto_init() }
func file_button_proto_init() {
	if File_button_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_button_proto_rawDesc), len(file_button_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_button_proto_goTypes,
		DependencyIndexes: file_button_proto_depIdxs,
		EnumInfos:         file_button_proto_enumTypes,
		MessageInfos:      file_button_proto_msgTypes,
	}.Build()
	File_button_proto = out.File
	file_button_proto_goTypes = nil
	file_button_proto_depIdxs = nil
}
