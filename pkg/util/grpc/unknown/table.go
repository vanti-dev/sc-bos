package unknown

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// NewMethodTable returns a new *MethodTable.
func NewMethodTable() *MethodTable {
	return &MethodTable{byMethod: make(map[string]Method)}
}

// MethodTable holds the mapping from gRPC method to downstream target.
// The methods of MethodTable are safe for concurrent use.
type MethodTable struct {
	mu       sync.RWMutex // guards byMethod
	byMethod map[string]Method
}

// Get returns the downstream target for the given gRPC method, if any.
func (r *MethodTable) Get(method string) (Method, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.byMethod[method]
	return m, ok
}

// Add registers a downstream target for the given gRPC method, returning true if it was added.
func (r *MethodTable) Add(method string, target Method) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.byMethod[method]
	if !ok {
		r.byMethod[method] = target
	}
	return !ok
}

// Set registers a downstream target for the given gRPC method, replacing any existing value.
func (r *MethodTable) Set(method string, target Method) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byMethod[method] = target
}

// Delete deletes the downstream target for the given gRPC method, returning true if it was deleted.
func (r *MethodTable) Delete(method string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.byMethod[method]
	if ok {
		delete(r.byMethod, method)
	}
	return ok
}

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

// MsgRecver is the interface for [grpc.ServerStream.RecvMsg]
type MsgRecver interface {
	RecvMsg(any) error
}

// ResolverFunc is a Resolver that resolves a *grpc.ClientConn using a function.
type ResolverFunc func(mr MsgRecver) (grpc.ClientConnInterface, error)

func (rf ResolverFunc) Resolve(mr MsgRecver) (grpc.ClientConnInterface, error) {
	return rf(mr)
}

// NewFixedResolver returns a Resolver that always returns the same *grpc.ClientConn.
func NewFixedResolver(cc grpc.ClientConnInterface) Resolver {
	return ResolverFunc(func(_ MsgRecver) (grpc.ClientConnInterface, error) {
		return cc, nil
	})
}

// KeyFunc extracts a key from an incoming message.
type KeyFunc func(MsgRecver) (string, error)

// NewKeyResolver returns a *KeyResolver that resolves a *grpc.ClientConn based on a key extracted from the incoming message.
func NewKeyResolver(key KeyFunc) *KeyResolver {
	return &KeyResolver{key: key, byKey: make(map[string]grpc.ClientConnInterface)}
}

// KeyResolver resolves a *grpc.ClientConn based on a key extracted from the incoming message.
// KeyResolver is safe for concurrent use iff the KeyFunc is safe for concurrent use.
type KeyResolver struct {
	mu    sync.RWMutex // guards byKey
	byKey map[string]grpc.ClientConnInterface
	key   KeyFunc
}

// Set registers a *grpc.ClientConn with the given key, replacing any existing value.
func (k *KeyResolver) Set(key string, cc grpc.ClientConnInterface) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.byKey[key] = cc
}

// Delete removes the *grpc.ClientConn associated with the given key, returning true if it was deleted.
func (k *KeyResolver) Delete(key string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	_, ok := k.byKey[key]
	if ok {
		delete(k.byKey, key)
	}
	return ok
}

func (k *KeyResolver) Resolve(mr MsgRecver) (grpc.ClientConnInterface, error) {
	key, err := k.key(mr)
	if err != nil {
		return nil, err
	}
	k.mu.RLock()
	cc, ok := k.byKey[key]
	k.mu.RUnlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot resolve %q", key)
	}
	return cc, nil
}

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
