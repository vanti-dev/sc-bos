package router

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

type MsgRecver interface {
	RecvMsg(into any) error
}

// Router implements MethodResolver and can route requests to a grpc.ClientConnInterface based on a combination
// of service name and a key which can be extracted from the request.
//
// Router only supports services which are registered using AddService. Key-based routes need the registered
// service to support routing by key (such as services constructed with NewRoutedService). If a key-based route is
// registered for a service that does not support routing by key, it will never match.
//
// Router supports four kinds of routes:
//  1. Service-and-key routes: Matches when both the service and key extracted from the request match.
//  2. Key-only route: Matches any key-routable service when the key extracted from the request matches.
//  3. Service-only route: Matches any method for a specific service. The service does not need to be key-routable.
//  4. Default route: Matches when no other route matches.
//
// The list above is in order of precedence.
type Router struct {
	keyInterceptor KeyInterceptor

	m        sync.RWMutex
	services map[string]*Service
	routes   map[routeID]grpc.ClientConnInterface
}

func New(opts ...Option) *Router {
	r := &Router{
		services: make(map[string]*Service),
		routes:   make(map[routeID]grpc.ClientConnInterface),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// AddService adds support for the service with the router.
// This allows the router to handle requests for this service.
// If a service with the same fully qualified name already exists, it is not added again, and ErrServiceExists is returned.
//
// AddService does not add any routes for the service, so by default all methods are unimplemented.
func (r *Router) AddService(s *Service) error {
	r.m.Lock()
	defer r.m.Unlock()

	name := s.Name()
	if _, exists := r.services[name]; exists {
		return ErrServiceExists
	}
	r.services[name] = s
	return nil
}

// DeleteService removes support for the service with the given fully qualified name from the router.
//
// It does not delete routes that name this service, but those routes will no longer match.
func (r *Router) DeleteService(name string) (exists bool) {
	r.m.Lock()
	defer r.m.Unlock()

	_, exists = r.services[name]
	if exists {
		delete(r.services, name)
	}
	return exists
}

// GetService returns the service registered with the given name.
// Returns nil if no service is registered with the given name.
func (r *Router) GetService(name string) *Service {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.services[name]
}

// RegisterService is a convenience implementation of grpc.ServiceRegistrar.
//
// It adds the service to the router (if not already present) as an unrouted service, and adds a service-only route
// backed by the given implementation.
// The service must be present in the global protobuf registry, or RegisterService will panic.
// It will also panic if a service-only route already exists for the service.
func (r *Router) RegisterService(desc *grpc.ServiceDesc, impl any) {
	reflectDesc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(desc.ServiceName))
	if err != nil {
		panic(err)
	}
	reflectSrvDesc := reflectDesc.(protoreflect.ServiceDescriptor)

	err = r.AddService(NewUnroutedService(reflectSrvDesc))
	if err != nil && !errors.Is(err, ErrServiceExists) {
		panic(err)
	}

	conn := wrap.ServerToClient(*desc, impl)

	err = r.AddRoute(desc.ServiceName, "", conn)
	if err != nil {
		panic(err)
	}
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

// AddRoute registers a target connection to be used when matching a specified route.
//
// The four kinds of routes are supported:
// 1. Service-and-key routes: both service and key are non-empty; the registered service must be key-routable.
// 2. Key-only route: service is empty, key is non-empty; the registered service must be key-routable.
// 3. Service-only route: service is non-empty, key is empty; the registered service does not need to be key-routable.
// 4. Default route: both service and key are empty.
//
// Returns ErrUnknownService if a service is specified but not registered.
// Returns ErrRouteExists if the same route is already registered.
func (r *Router) AddRoute(service, key string, target grpc.ClientConnInterface) error {
	r.m.Lock()
	defer r.m.Unlock()

	// don't allow routes for services that we don't know about
	if service != "" {
		if _, exists := r.services[service]; !exists {
			return ErrUnknownService
		}
	}

	id := routeID{Service: service, Key: key}
	if _, exists := r.routes[id]; exists {
		return ErrRouteExists
	}
	r.routes[id] = target
	return nil
}

// DeleteRoute removes a route.
// The service and key paremeters are intepreted the same way as in AddRoute.
//
// Returns true if the route existed and was removed.
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
		var candidates []routeID // what routes should we try to match?

		if keyFunc, exists := service.keys[methodName]; exists {
			// we can route by key
			key, err := keyFunc(mr)
			if err != nil {
				return nil, err
			}
			if r.keyInterceptor != nil {
				key, err = r.keyInterceptor(key)
				if err != nil {
					return nil, err
				}
			}
			candidates = []routeID{
				{Service: serviceName, Key: key},
				{Service: "", Key: key},
				{Service: serviceName, Key: ""},
				{Service: "", Key: ""},
			}
		} else {
			// can't route by key, we can only try routes that don't involve a key
			candidates = []routeID{
				{Service: serviceName, Key: ""},
				{Service: "", Key: ""},
			}
		}

		r.m.RLock()
		defer r.m.RUnlock()
		for _, candidate := range candidates {
			if conn, exists := r.routes[candidate]; exists {
				return conn, nil
			}
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

func NewRoutedService(desc protoreflect.ServiceDescriptor, keyField string) (*Service, error) {
	keys := make(map[string]KeyFunc)
	for i := 0; i < desc.Methods().Len(); i++ {
		method := desc.Methods().Get(i)
		keyFunc, err := FieldKey(method.Input(), keyField)
		if err != nil {
			return nil, fmt.Errorf("method %s is not routable: %w", method.Name(), err)
		}
		keys[string(method.Name())] = keyFunc
	}
	return &Service{
		descriptor: desc,
		keys:       keys,
	}, nil
}

func (s *Service) Name() string {
	return string(s.descriptor.FullName())
}

func (s *Service) Descriptor() protoreflect.ServiceDescriptor {
	return s.descriptor
}

// KeyRoutable returns true if the service supports routing by key.
func (s *Service) KeyRoutable() bool {
	return len(s.keys) > 0
}

var (
	ErrRouteExists    = errors.New("route already exists in router")
	ErrServiceExists  = errors.New("service already exists in router")
	ErrUnknownService = errors.New("unknown service")
)

// used as a map key for storing and looking up routes
// an empty field is interpreted as a wildcard
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

type Option func(router *Router)

func WithKeyInterceptor(interceptor KeyInterceptor) Option {
	return func(router *Router) {
		router.keyInterceptor = interceptor
	}
}

type KeyInterceptor func(key string) (mappedKey string, err error)
