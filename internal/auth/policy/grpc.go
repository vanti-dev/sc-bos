package policy

import (
	"context"
	"crypto/x509"
	"log"
	"regexp"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/vanti-dev/bsp-ew/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func GRPCUnaryInterceptor(verifier auth.TokenVerifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		err = checkPolicyGrpc(ctx, verifier, req)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func GRPCStreamingInterceptor(verifier auth.TokenVerifier) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := checkPolicyGrpc(ss.Context(), verifier, nil)
		if err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func checkPolicyGrpc(ctx context.Context, verifier auth.TokenVerifier, request any) error {
	service, method, ok := getGrpcServiceMethod(ctx)
	if !ok {
		return status.Error(codes.Internal, "failed to resolve method")
	}

	var tokenClaims *auth.TokenClaims

	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		log.Printf("no request bearer token: %s", err.Error())
	}

	if token != "" && verifier != nil {
		tokenClaims, err = verifier.VerifyAccessToken(ctx, token)
		if err != nil {
			tokenClaims = nil
			log.Printf("token failed verification: %s", err.Error())
		}
	}

	cert := getConnectionVerifiedCertificate(ctx)

	input := Attributes{
		Service:          service,
		Method:           method,
		Request:          request,
		CertificateValid: cert != nil,
		Certificate:      cert,
		TokenValid:       tokenClaims != nil,
		TokenClaims:      tokenClaims,
	}

	return CheckAttributes(ctx, input)
}

var grpcMethodRegexp = regexp.MustCompile("^/([^/]*)/([^/]*)$")

func getGrpcServiceMethod(ctx context.Context) (service, method string, ok bool) {
	full, ok := grpc.Method(ctx)
	if !ok {
		return
	}

	groups := grpcMethodRegexp.FindStringSubmatch(full)
	if len(groups) != 3 {
		ok = false
		return
	}

	service = groups[1]
	method = groups[2]
	ok = true
	return
}

// find the certificate of the connection peer that was verified when the connection was established.
// returns nil if no certificate was verified
func getConnectionVerifiedCertificate(ctx context.Context) *x509.Certificate {
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
