package router

import (
	"context"
	"errors"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MsgRecver interface {
	RecvMsg(into any) error
}

type Router struct {
	m             sync.RWMutex
	services      map[string]*Service
	routes        map[routeID]grpc.ClientConnInterface
	defaultRoutes map[string]grpc.ClientConnInterface // map of service name to default conn
}

func New() *Router {
	return &Router{
		services: make(map[string]*Service),
		routes:   make(map[routeID]grpc.ClientConnInterface),
	}
}

// SupportService registers the service with the router.
// This allows the router to handle requests for this service.
// If a service with the same fully qualified name already exists, it is not added again, and true is returned.
//
// SupportService does not add any routes for the service, so by default all methods are unimplemented.
func (r *Router) SupportService(s *Service) (exists bool) {
	r.m.Lock()
	defer r.m.Unlock()

	name := s.Name()
	if _, exists := r.services[name]; exists {
		return true
	}
	r.services[name] = s
	return false
}

func (r *Router) SupportsService(name string) bool {
	r.m.RLock()
	defer r.m.RUnlock()
	_, exists := r.services[name]
	return exists
}

func (r *Router) GetServiceInfo() map[string]grpc.ServiceInfo {
	r.m.RLock()
	defer r.m.RUnlock()

	services := make(map[string]grpc.ServiceInfo)
	for name, s := range r.services {
		desc := s.Descriptor()
		var methodInfos []grpc.MethodInfo
		for i := 0; i < desc.Methods().Len(); i++ {
			method := desc.Methods().Get(i)
			methodInfos = append(methodInfos, grpc.MethodInfo{
				Name:           string(method.Name()),
				IsClientStream: method.IsStreamingClient(),
				IsServerStream: method.IsStreamingServer(),
			})
		}
		services[name] = grpc.ServiceInfo{
			Methods: methodInfos,
		}
	}
	return services
}

func (r *Router) Invoke(ctx context.Context, method string, args, reply any, opts ...grpc.CallOption) error {
	panic("not implemented")
}

func (r *Router) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	panic("not implemented")
}

// AddRoute registers a target connection to be used for a specific combination of service and key.
// This route can only be matched if the corresponding service supplied to SupportService supports routing by key.
func (r *Router) AddRoute(service, key string, target grpc.ClientConnInterface) error {
	r.m.Lock()
	defer r.m.Unlock()

	// don't allow routes for services that we don't know about
	if _, exists := r.services[service]; !exists {
		return ErrUnknownService
	}

	id := routeID{Service: service, Key: key}
	if _, exists := r.routes[id]; exists {
		return ErrRouteExists
	}
	r.routes[id] = target
	return nil
}

func (r *Router) DeleteRoute(service, key string) (exists bool) {
	r.m.Lock()
	defer r.m.Unlock()

	id := routeID{Service: service, Key: key}
	_, exists = r.routes[id]
	if exists {
		delete(r.routes, id)
	}
	return exists
}

func (r *Router) AddDefaultRoute(service string, target grpc.ClientConnInterface) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, exists := r.services[service]; !exists {
		return ErrUnknownService
	}
	_, exists := r.defaultRoutes[service]
	if exists {
		return ErrRouteExists
	}
	r.defaultRoutes[service] = target
	return nil
}

func (r *Router) DeleteDefaultRoute(service string) (exists bool) {
	r.m.Lock()
	defer r.m.Unlock()

	_, exists = r.defaultRoutes[service]
	if exists {
		delete(r.defaultRoutes, service)
	}
	return exists
}

func (r *Router) ResolveMethod(fullName string) (Method, bool) {
	r.m.RLock()
	defer r.m.RUnlock()

	serviceName, methodName, ok := parseMethod(fullName)
	if !ok {
		return Method{}, false
	}
	// if the service is not registered then we certainly can't resolve the method
	service, exists := r.services[serviceName]
	if !exists {
		return Method{}, false
	}
	methodDesc := service.descriptor.Methods().ByName(protoreflect.Name(methodName))
	if methodDesc == nil {
		return Method{}, false
	}

	connResolver := ResolverFunc(func(mr MsgRecver) (grpc.ClientConnInterface, error) {
		// first try to route based on key
		keyFunc, exists := service.keys[methodName]
		if exists {
			key, err := keyFunc(mr)
			if err != nil {
				return nil, err
			}
			id := routeID{Service: serviceName, Key: key}
			if conn, exists := r.routes[id]; exists {
				return conn, nil
			}
		}

		// if we cannot route by key, then try routing to the default conn for the service
		conn, exists := r.defaultRoutes[serviceName]
		if exists {
			return conn, nil
		}

		return nil, status.Error(codes.NotFound, "no route found")
	})

	return Method{
		StreamDesc: descriptorToStreamDesc(methodDesc),
		Resolver:   connResolver,
	}, true
}

type Service struct {
	// immutable after creation
	descriptor protoreflect.ServiceDescriptor
	keys       map[string]KeyFunc // map from unqualified method name to key func
}

// NewUnroutedService creates a new service that does not support routing by key.
//
// All requests to this service will be directed to the default route.
func NewUnroutedService(desc protoreflect.ServiceDescriptor) *Service {
	return &Service{descriptor: desc}
}

func NewRoutedService(desc protoreflect.ServiceDescriptor, keyField string) *Service {
	return &Service{
		descriptor: desc,
		keys:       make(map[string]KeyFunc),
	}
}

func (s *Service) Name() string {
	return string(s.descriptor.FullName())
}

func (s *Service) Descriptor() protoreflect.ServiceDescriptor {
	return s.descriptor
}

var (
	ErrRouteExists    = errors.New("route already exists")
	ErrUnknownService = errors.New("unknown service")
)

type routeID struct {
	Key     string
	Service string
}

func descriptorToStreamDesc(desc protoreflect.MethodDescriptor) grpc.StreamDesc {
	return grpc.StreamDesc{
		StreamName:    string(desc.Name()),
		ClientStreams: desc.IsStreamingClient(),
		ServerStreams: desc.IsStreamingServer(),
	}
}

func parseMethod(fullMethod string) (service, method string, ok bool) {
	// strip leading /
	if !strings.HasPrefix(fullMethod, "/") {
		return "", "", false
	}
	fullMethod = fullMethod[1:]

	idx := strings.LastIndex(fullMethod, "/")
	if idx < 0 {
		return "", "", false
	}

	service = fullMethod[:idx]
	method = fullMethod[idx+1:]
	if service == "" || method == "" {
		return "", "", false
	}
	return service, method, true
}
