package rpcutil

import (
	"context"
	"crypto/x509"
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

// VerifiedCertFromServerContext finds the certificate of the connection peer that was verified when the
// connection was established.
// Returns nil if no certificate was verified.
func VerifiedCertFromServerContext(ctx context.Context) *x509.Certificate {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil
	}

	tlsInfo, ok := peerInfo.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil
	}

	verifiedChains := tlsInfo.State.VerifiedChains
	if len(verifiedChains) == 0 {
		return nil
	}
	return verifiedChains[0][0]
}
