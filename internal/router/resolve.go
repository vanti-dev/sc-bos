package router

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Method describes an RPC we are dynamically referencing.
type Method struct {
	Desc     protoreflect.MethodDescriptor
	Resolver ConnResolver
}

// ConnResolver resolves a grpc.ClientConnInterface from an incoming message.
type ConnResolver interface {
	// ResolveConn resolves a grpc.ClientConnInterface to proxy a request to.
	// The MsgRecver mr may be used at most once, to retrieve the request message, or the first message in the
	// client-to-server stream.
	ResolveConn(mr MsgRecver) (grpc.ClientConnInterface, error)
}

// ConnResolverFunc is a ConnResolver implemented as a function.
type ConnResolverFunc func(mr MsgRecver) (grpc.ClientConnInterface, error)

func (rf ConnResolverFunc) ResolveConn(mr MsgRecver) (grpc.ClientConnInterface, error) {
	return rf(mr)
}

type KeyFunc func(mr MsgRecver) (string, error)

// FieldKey returns a KeyFunc that extracts a named property value from messages that are described by msgDesc.
// The key field must have the following properties:
//
//   - The property is of type string.
//   - The property is not repeated.
//
// Safe for concurrent use.
func FieldKey(msgDesc protoreflect.MessageDescriptor, field string) (KeyFunc, error) {
	nameFieldDesc := msgDesc.Fields().ByName(protoreflect.Name(field))
	if nameFieldDesc == nil {
		return nil, fmt.Errorf("no %s field found in %q", field, msgDesc.FullName())
	}
	if nameFieldDesc.Kind() != protoreflect.StringKind {
		return nil, fmt.Errorf("%s field in %q is not a string", field, msgDesc.FullName())
	}
	if nameFieldDesc.Cardinality() == protoreflect.Repeated {
		return nil, fmt.Errorf("%s name field in %q is repeated", field, msgDesc.FullName())
	}
	return func(mr MsgRecver) (string, error) {
		m := dynamicpb.NewMessage(msgDesc)
		if err := mr.RecvMsg(m); err != nil {
			return "", err
		}
		key := m.Get(nameFieldDesc).String()
		return key, nil
	}, nil
}

// captureMsgRecver is a MsgRecver that captures the message received.
type captureMsgRecver struct {
	MsgRecver
	msg any
}

func (cmr *captureMsgRecver) RecvMsg(m any) error {
	err := cmr.MsgRecver.RecvMsg(m)
	cmr.msg = m
	return err
}
