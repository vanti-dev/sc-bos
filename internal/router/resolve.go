package router

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Method describes an RPC we are dynamically referencing.
type Method struct {
	StreamDesc grpc.StreamDesc
	Resolver   Resolver
}

// Resolver resolves a *grpc.ClientConn from an incoming message.
// Resolvers Resolve method must be safe to be called from multiple goroutines.
type Resolver interface {
	Resolve(mr MsgRecver) (grpc.ClientConnInterface, error)
}

// ResolverFunc is a Resolver that resolves a *grpc.ClientConn using a function.
type ResolverFunc func(mr MsgRecver) (grpc.ClientConnInterface, error)

func (rf ResolverFunc) Resolve(mr MsgRecver) (grpc.ClientConnInterface, error) {
	return rf(mr)
}

type KeyFunc func(mr MsgRecver) (string, error)

// NameKey returns a KeyFunc that extracts "name" keys from messages that are described by msgDesc.
// "name" keys have the following properties:
//
//   - The property is named "name".
//   - The property is of type string.
//   - The property is not repeated.
//
// NameKey is safe for concurrent use.
func NameKey(msgDesc protoreflect.MessageDescriptor) (KeyFunc, error) {
	nameFieldDesc := msgDesc.Fields().ByName("name")
	if nameFieldDesc == nil {
		return nil, fmt.Errorf("no name field found in %q", msgDesc.FullName())
	}
	if nameFieldDesc.Kind() != protoreflect.StringKind {
		return nil, fmt.Errorf("name field in %q is not a string", msgDesc.FullName())
	}
	if nameFieldDesc.Cardinality() == protoreflect.Repeated {
		return nil, fmt.Errorf("name field in %q is repeated", msgDesc.FullName())
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
