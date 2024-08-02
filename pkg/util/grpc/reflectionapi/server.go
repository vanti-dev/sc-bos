package reflectionapi

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	reflectionv1alphapb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// NewServer creates a *Server configured to use s and the global protoregistry types for reflection.
func NewServer(s *grpc.Server) *Server {
	return &Server{
		infoProvider: newServiceInfoProviderSet(s),
		descResolver: newDescResolverSet(protoregistry.GlobalFiles),
	}
}

// Server is an alternative to reflection.Register that supports adjusting the available services and descriptors at runtime.
// Call Register instead of reflection.Register to register the correct APIs with the grpc.Server.
// Add and Remove are safe to call from multiple goroutines, even after the grpc.Server has started.
type Server struct {
	infoProvider *serviceInfoProviderSet
	descResolver *descResolverSet

	mu      sync.Mutex
	records []record
}

// record keeps track of additional resources associated with a client connection.
type record struct {
	c       *grpc.ClientConn
	sip     reflection.ServiceInfoProvider // this also exists in Server.infoProvider
	dr      protodesc.Resolver             // this also exists in Server.descResolver
	drClose func() error
}

// Register registers the v1 and v1alpha reflection services on srv.
func (s *Server) Register(srv *grpc.Server) {
	gwrOps := reflection.ServerOptions{
		Services:           s.infoProvider,
		DescriptorResolver: s.descResolver,
	}
	reflectionv1alphapb.RegisterServerReflectionServer(srv, reflection.NewServer(gwrOps))
	reflectionpb.RegisterServerReflectionServer(srv, reflection.NewServerV1(gwrOps))
}

// Add adds services and descriptors associated with c to the reflection server.
func (s *Server) Add(c *grpc.ClientConn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := record{
		c:   c,
		sip: &serviceInfoProvider{client: reflectionpb.NewServerReflectionClient(c)},
	}
	r.dr, r.drClose = newDescResolver(reflectionpb.NewServerReflectionClient(c))
	s.infoProvider.Add(r.sip)
	s.descResolver.Add(r.dr)
	s.records = append(s.records, r)
}

// Remove removes services and descriptors associated with c from the reflection server, returning true if it was removed.
func (s *Server) Remove(c *grpc.ClientConn) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.records {
		if r.c == c {
			s.infoProvider.Remove(r.sip)
			s.descResolver.Remove(r.dr)
			_ = r.drClose()
			s.records = append(s.records[:i], s.records[i+1:]...)
			return true
		}
	}
	return false
}

// newServiceInfoProviderSet returns a new *serviceInfoProviderSet that aggregates multiple providers.
func newServiceInfoProviderSet(providers ...reflection.ServiceInfoProvider) *serviceInfoProviderSet {
	return &serviceInfoProviderSet{providers: providers}
}

// serviceInfoProviderSet is a reflection.ServiceInfoProvider that aggregates multiple providers.
// It is safe for concurrent use.
type serviceInfoProviderSet struct {
	mu        sync.RWMutex // guards providers
	providers []reflection.ServiceInfoProvider
}

// make sure we implement the required interface
var _ reflection.ServiceInfoProvider = (*serviceInfoProviderSet)(nil)

// GetServiceInfo implements reflection.serviceInfoProviderSet by collecting service information from all providers.
func (s *serviceInfoProviderSet) GetServiceInfo() map[string]grpc.ServiceInfo {
	s.mu.RLock()
	ps := s.providers[:]
	s.mu.RUnlock()

	m := make(map[string]grpc.ServiceInfo)
	for _, p := range ps {
		for k, v := range p.GetServiceInfo() {
			if _, ok := m[k]; ok {
				continue // use the first info we found
			}
			m[k] = v
		}
	}
	return m
}

// Add adds the provider p to the list of providers.
func (s *serviceInfoProviderSet) Add(p reflection.ServiceInfoProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providers = append(s.providers, p)
}

// Remove removes the provider `==` to p from the list of providers.
func (s *serviceInfoProviderSet) Remove(p reflection.ServiceInfoProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, v := range s.providers {
		if v == p {
			s.providers = append(s.providers[:i], s.providers[i+1:]...)
			return
		}
	}
}

// serviceInfoProvider is a reflection.ServiceInfoProvider that uses the reflection service to get service information from a server.
type serviceInfoProvider struct {
	client reflectionpb.ServerReflectionClient
}

func (r serviceInfoProvider) GetServiceInfo() map[string]grpc.ServiceInfo {
	m := make(map[string]grpc.ServiceInfo)
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	s, err := r.client.ServerReflectionInfo(ctx)
	if err != nil {
		return nil
	}

	services, err := ListServices(s)
	if err != nil {
		return nil
	}

	// collect all the file descriptors into a set we can query for type information
	seen := make(map[string]struct{})
	fdSet := &descriptorpb.FileDescriptorSet{}
	for _, service := range services {
		res, err := FileContainingSymbol(s, service.GetName())
		if err != nil {
			continue
		}
		for _, fdBytes := range res {
			fileDesc := &descriptorpb.FileDescriptorProto{}
			err := proto.Unmarshal(fdBytes, fileDesc)
			if err != nil {
				continue
			}
			if _, ok := seen[fileDesc.GetName()]; ok {
				continue // skip duplicates
			}
			fdSet.File = append(fdSet.File, fileDesc)
			seen[fileDesc.GetName()] = struct{}{}
		}
	}
	files, err := protodesc.NewFiles(fdSet)
	if err != nil {
		return nil
	}

	for _, service := range services {
		desc, err := files.FindDescriptorByName(protoreflect.FullName(service.GetName()))
		if err != nil {
			continue
		}
		// this should always be true as we're asking for descriptors for the service name
		serviceDesc, ok := desc.(protoreflect.ServiceDescriptor)
		if !ok {
			continue
		}
		si := grpc.ServiceInfo{}
		for i := 0; i < serviceDesc.Methods().Len(); i++ {
			m := serviceDesc.Methods().Get(i)
			si.Methods = append(si.Methods, grpc.MethodInfo{
				Name:           string(m.Name()),
				IsClientStream: m.IsStreamingClient(),
				IsServerStream: m.IsStreamingServer(),
			})
		}
		m[service.GetName()] = si
	}
	return m
}

// newDescResolverSet returns a new descResolverSet that aggregates multiple resolvers.
func newDescResolverSet(resolvers ...protodesc.Resolver) *descResolverSet {
	return &descResolverSet{resolvers: resolvers}
}

// descResolverSet is a protodesc.Resolver that aggregates multiple resolvers.
// It is safe for concurrent use.
type descResolverSet struct {
	mu        sync.RWMutex // guards resolvers
	resolvers []protodesc.Resolver
}

// make sure we implement the required interface
var _ protodesc.Resolver = (*descResolverSet)(nil)

func (d *descResolverSet) Add(r protodesc.Resolver) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resolvers = append(d.resolvers, r)
}

func (d *descResolverSet) Remove(r protodesc.Resolver) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, v := range d.resolvers {
		if v == r {
			d.resolvers = append(d.resolvers[:i], d.resolvers[i+1:]...)
			return
		}
	}
}

func (d *descResolverSet) FindFileByPath(s string) (protoreflect.FileDescriptor, error) {
	d.mu.RLock()
	rs := d.resolvers[:]
	d.mu.RUnlock()

	if len(rs) == 0 {
		return nil, protoregistry.NotFound
	}

	var errs []error
	for _, resolver := range rs {
		fd, err := resolver.FindFileByPath(s)
		if err == nil {
			return fd, nil
		}
		errs = append(errs, err)
	}
	return nil, errors.Join(errs...)
}

func (d *descResolverSet) FindDescriptorByName(fullName protoreflect.FullName) (protoreflect.Descriptor, error) {
	d.mu.RLock()
	rs := d.resolvers[:]
	d.mu.RUnlock()

	if len(rs) == 0 {
		return nil, protoregistry.NotFound
	}

	var errs []error
	for _, resolver := range rs {
		desc, err := resolver.FindDescriptorByName(fullName)
		if err == nil {
			return desc, nil
		}
		errs = append(errs, err)
	}
	return nil, errors.Join(errs...)
}

// newDescResolver returns a new protodesc.Resolver that uses the reflection (v1) service to resolve descriptors.
func newDescResolver(client reflectionpb.ServerReflectionClient) (r protodesc.Resolver, stop func() error) {
	res := &descResolver{client: client}
	return res, res.closeStream
}

type descResolver struct {
	client reflectionpb.ServerReflectionClient

	// Our files cache is only consistent while the stream remains valid due to potential load balancing.
	// We try to keep the stream open as long as we can to avoid re-fetching the same files.
	streamMu   sync.Mutex // guards the below
	stream     reflectionpb.ServerReflection_ServerReflectionInfoClient
	stopStream func()
	files      *files
}

func (r *descResolver) FindFileByPath(s string) (protoreflect.FileDescriptor, error) {
	session, err := r.openSession()
	if err != nil {
		return nil, err
	}

	desc, err := session.FindFileByPath(s)
	if err == nil {
		return desc, nil
	}

	fileDescs, err := FileByFilename(session.stream, s)
	if status.Code(err) == codes.NotFound {
		// the stream is still valid, the server just didn't know about the path
		return nil, err
	}
	if err != nil {
		_ = session.stop() // we think the stream is now invalid, we could retry but so could the caller
		return nil, err
	}
	if err := r.registerFiles(session.files, fileDescs); err != nil {
		return nil, err
	}

	return session.FindFileByPath(s)
}

func (r *descResolver) FindDescriptorByName(fullName protoreflect.FullName) (protoreflect.Descriptor, error) {
	session, err := r.openSession()
	if err != nil {
		return nil, err
	}

	desc, err := session.FindDescriptorByName(fullName)
	if err == nil {
		return desc, nil
	}

	// on err, try to find from the stream
	fileDescs, err := FileContainingSymbol(session.stream, string(fullName))
	if status.Code(err) == codes.NotFound {
		return nil, fmt.Errorf("remote: %w", err)
	}
	if err != nil {
		_ = r.closeStream() // it's invalid now, we think
		return nil, err
	}
	if err := r.registerFiles(session.files, fileDescs); err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return session.FindDescriptorByName(fullName)
}

func (r *descResolver) registerFiles(files *files, descs [][]byte) error {
	// note: the reflection API does not specify an order to the files it returns.
	// The following logic is inspired by protodesc.NewFiles and registers files with dependencies first.

	fdp := make(map[string]*descriptorpb.FileDescriptorProto)
	for _, fileDesc := range descs {
		descMsg := &descriptorpb.FileDescriptorProto{}
		err := proto.Unmarshal(fileDesc, descMsg)
		if err != nil {
			return err
		}
		fdp[descMsg.GetName()] = descMsg
	}

	r.streamMu.Lock()
	defer r.streamMu.Unlock()
	for _, descMsg := range fdp {
		if err := registerFileDeps(files, descMsg, fdp); err != nil {
			return err
		}
	}
	return nil
}

func registerFileDeps(r *files, fd *descriptorpb.FileDescriptorProto, files map[string]*descriptorpb.FileDescriptorProto) error {
	files[fd.GetName()] = nil // to detect cycles
	for _, dep := range fd.Dependency {
		depfd, ok := files[dep]
		if depfd == nil {
			if ok {
				return fmt.Errorf("cycle detected: %q", dep)
			}
			continue
		}
		if err := registerFileDeps(r, depfd, files); err != nil {
			return err
		}
	}

	delete(files, fd.GetName())
	f, err := protodesc.NewFile(fd, newDescResolverSet(protoregistry.GlobalFiles, r))
	if err != nil {
		return err
	}
	return r.RegisterFile(f)
}

func (r *descResolver) openSession() (session, error) {
	r.streamMu.Lock()
	defer r.streamMu.Unlock()
	if r.stream != nil {
		return r.newSessionLocked(), nil
	}
	ctx, stop := context.WithCancel(context.Background())
	stream, err := r.client.ServerReflectionInfo(ctx)
	if err != nil {
		stop()
		return session{}, err
	}
	r.stream = stream
	r.stopStream = stop
	r.files = &files{files: &protoregistry.Files{}}
	return r.newSessionLocked(), nil
}

func (r *descResolver) closeStream() error {
	r.streamMu.Lock()
	if r.stream == nil {
		r.streamMu.Unlock()
		return nil
	}
	s := r.newSessionLocked()
	r.streamMu.Unlock()
	return s.stop()
}

func (r *descResolver) newSessionLocked() session {
	stream := r.stream
	files := r.files
	stopStream := r.stopStream
	return session{
		stream: stream,
		files:  files,
		stop: func() error {
			r.streamMu.Lock()
			defer r.streamMu.Unlock()
			if r.stream != stream {
				return nil
			}
			r.stream = nil
			r.files = nil
			stopStream()
			return stream.CloseSend()
		},
	}
}

// files is like protoregistry.Files but safe for concurrent use.
type files struct {
	mu    sync.RWMutex
	files *protoregistry.Files
}

func (f *files) FindFileByPath(s string) (protoreflect.FileDescriptor, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.files.FindFileByPath(s)
}

func (f *files) FindDescriptorByName(fullName protoreflect.FullName) (protoreflect.Descriptor, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.files.FindDescriptorByName(fullName)
}

func (f *files) RegisterFile(fd protoreflect.FileDescriptor) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.files.RegisterFile(fd)
}

type session struct {
	stream reflectionpb.ServerReflection_ServerReflectionInfoClient
	stop   func() error
	*files
}
