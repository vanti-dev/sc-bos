package router

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/smart-core-os/sc-api/go/traits"
)

func TestCopyRecver_RecvMsg(t *testing.T) {
	msg := func(s string) proto.Message {
		return &traits.Metadata{Name: s}
	}
	msg2 := func(s string) proto.Message {
		return &traits.TraitMetadata{Name: s}
	}
	// simulate descriptors from the reflection package
	mkDynamic := func(m proto.Message) proto.Message {
		// convert the descriptor to an equivalent one via the proto file representation
		descriptor := m.ProtoReflect().Descriptor()
		file, err := protodesc.NewFile(protodesc.ToFileDescriptorProto(descriptor.ParentFile()), protoregistry.GlobalFiles)
		if err != nil {
			t.Fatalf("NewFile() error = %v", err)
		}
		msgDesc := file.Messages().ByName(descriptor.Name())
		if msgDesc == nil {
			t.Fatalf("message not found in file: %v", descriptor.Name())
		}
		dm := dynamicpb.NewMessage(msgDesc)

		// copy over the data, don't use proto.Merge!
		data, err := proto.Marshal(m)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if err := proto.Unmarshal(data, dm); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		return dm
	}

	tests := []struct {
		name               string
		src, dst, want     proto.Message
		wantErr, wantPanic bool
	}{
		{"msg,msg", msg("a"), msg("b"), msg("a"), false, false},
		{"msg,dynamic", msg("a"), mkDynamic(msg("b")), msg("a"), false, false},
		{"dynamic,msg", mkDynamic(msg("a")), msg("b"), msg("a"), false, false},
		{"dynamic,dynamic", mkDynamic(msg("a")), mkDynamic(msg("b")), msg("a"), false, false},
		{"msg1,msg2", msg("1"), msg2("2"), msg2("2"), false, true},
		{"msg2,msg1", msg2("2"), msg("1"), msg("1"), false, true},
		{"msg1,dynamic2", msg("1"), mkDynamic(msg2("2")), msg2("2"), false, true},
		{"dynamic2,msg1", mkDynamic(msg2("2")), msg("1"), msg("1"), false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					if !tt.wantPanic {
						t.Errorf("unexpected panic: %v", r)
					}
				} else if tt.wantPanic {
					t.Error("expected panic")
				}
			}()

			cr := copyRecver{from: tt.src}
			if err := cr.RecvMsg(tt.dst); (err != nil) != tt.wantErr {
				t.Errorf("RecvMsg() error = %v, wantErr %v", err, tt.wantErr)
			}

			// proto.Equal doesn't work for messages with different descriptors
			if diff := cmp.Diff(tt.want, tt.dst, protocmp.Transform()); diff != "" {
				t.Errorf("RecvMsg() (-want,+got)\n%s", diff)
			}
		})
	}
}
