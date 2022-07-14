package policy

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/open-policy-agent/opa/rego"
	"github.com/vanti-dev/bsp-ew/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

	var authorization *auth.Authorization

	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		log.Printf("no request bearer token: %s", err.Error())
	}

	if token != "" {
		authorization, err = verifier.VerifyAccessToken(ctx, token)
		if err != nil {
			log.Printf("token failed verification: %s", err.Error())
		}
	}

	input := struct {
		Authorization *auth.Authorization
		Request       any
		Method        string
		Service       string
	}{
		Authorization: authorization,
		Method:        method,
		Service:       service,
		Request:       request,
	}

	query := fmt.Sprintf("data.%s.allow", service)
	r := rego.New(
		rego.Compiler(RegoCompiler),
		rego.Input(input),
		rego.Query(query),
	)
	result, err := r.Eval(ctx)
	if err != nil {
		return err
	}
	if !result.Allowed() {
		return status.Error(codes.PermissionDenied, "you are not authorized to perform this operation")
	}

	return nil
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
