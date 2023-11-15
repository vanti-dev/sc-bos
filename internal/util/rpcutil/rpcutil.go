package rpcutil

import (
	"context"
	"crypto/x509"
	"net/http"
	"regexp"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// ServerLogger creates a sub-logger with fields derived from a gRPC server context.
func ServerLogger(ctx context.Context, base *zap.Logger) *zap.Logger {
	return base.With(ServerLoggerFields(ctx)...)
}

func ServerLoggerFields(ctx context.Context) []zap.Field {
	fields := []zap.Field{zap.Namespace("grpc")}
	if service, method, ok := ServiceMethod(ctx); ok {
		fields = append(fields,
			zap.String("service", service),
			zap.String("method", method),
		)
	}
	if p, ok := peer.FromContext(ctx); ok {
		auth := "none"
		if p.AuthInfo != nil {
			auth = p.AuthInfo.AuthType()
		}
		fields = append(fields,
			zap.String("addr", p.Addr.String()),
			zap.String("auth", auth),
		)
	}

	return fields
}

func HTTPLoggerFields(request *http.Request) []zap.Field {
	fields := []zap.Field{
		zap.Namespace("http"),
		zap.String("method", request.Method),
		zap.String("url", request.RequestURI),
		zap.String("peer", request.RemoteAddr),
	}
	return fields
}

func HTTPLogger(request *http.Request, base *zap.Logger) *zap.Logger {
	return base.With(HTTPLoggerFields(request)...)
}

func SplitMethodPath(methodPath string) (service, method string, ok bool) {
	groups := grpcMethodRegexp.FindStringSubmatch(methodPath)
	if len(groups) != 3 {
		ok = false
		return
	}

	service = groups[1]
	method = groups[2]
	ok = true
	return
}

var grpcMethodRegexp = regexp.MustCompile("^/([^/]*)/([^/]*)$")

// ServiceMethod retrieves the service name and method name from a gRPC server context.
// Works like grpc.Method, but splits the path into service and method parts.
func ServiceMethod(ctx context.Context) (service, method string, ok bool) {
	full, ok := grpc.Method(ctx)
	if !ok {
		return
	}

	return SplitMethodPath(full)
}

// CertFromServerContext returns the peer provided certificate any whether it id valid or not.
func CertFromServerContext(ctx context.Context) (cert *x509.Certificate, valid bool) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, false
	}

	tlsInfo, ok := peerInfo.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, false
	}

	peerCerts := tlsInfo.State.PeerCertificates
	if len(peerCerts) == 0 {
		return nil, false
	}
	return peerCerts[0], len(tlsInfo.State.VerifiedChains) > 0
}
