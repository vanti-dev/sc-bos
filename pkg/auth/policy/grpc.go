package policy

import (
	"context"
	"crypto/x509"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/internal/util/rpcutil"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
)

type Interceptor struct {
	logger   *zap.Logger
	policy   Policy
	verifier token.Validator
}

func NewInterceptor(policy Policy, opts ...InterceptorOption) *Interceptor {
	interceptor := &Interceptor{
		logger: zap.NewNop(),
		policy: policy,
	}
	for _, o := range opts {
		o(interceptor)
	}
	return interceptor
}

func (i *Interceptor) GRPCUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp any, err error) {
		_, err = i.checkPolicyGrpc(ctx, nil, req, StreamAttributes{
			IsServerStream: false,
			IsClientStream: false,
			Open:           false,
		})
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (i *Interceptor) GRPCStreamingInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// for client / bidirectional streams we don't have a request, so we'll evaluate the policy without one
		// to check if it's OK to open the stream.
		// This isn't necessary for server streams; the client will immediately send the request message, and the
		// generated code will call RecvMsg to get this *before* control is transferred to the service implementation,
		// so we can check then using the serverStreamInterceptor
		var cachedCreds *verifiedCreds
		if info.IsClientStream {
			var err error
			cachedCreds, err = i.checkPolicyGrpc(ss.Context(), nil, nil, StreamAttributes{
				IsServerStream: info.IsServerStream,
				IsClientStream: info.IsClientStream,
				Open:           false,
			})
			if err != nil {
				return err
			}
		}

		cb := func(msg any) error {
			streamAttrs := StreamAttributes{
				IsServerStream: info.IsServerStream,
				IsClientStream: info.IsClientStream,
				// Only client/bidirectional streams cause policies to be evaluated once the stream is already
				// open.
				Open: info.IsClientStream,
			}

			_, err := i.checkPolicyGrpc(ss.Context(), cachedCreds, msg, streamAttrs)
			return err
		}
		wrapped := &serverStreamInterceptor{
			ServerStream: ss,
			cb:           cb,
		}
		return handler(srv, wrapped)
	}
}

// Returns a set of verified credentials that can be used to speed up future calls to checkPolicyGrpc for the same
// call (useful for streams). Pass nil creds the first time, then cache the creds.
func (i *Interceptor) checkPolicyGrpc(ctx context.Context, creds *verifiedCreds, req any, stream StreamAttributes) (*verifiedCreds, error) {
	service, method, ok := rpcutil.ServiceMethod(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to resolve method")
	}

	if creds == nil {
		tkn, err := grpc_auth.AuthFromMD(ctx, "Bearer")
		var tokenClaims *token.Claims
		if err == nil && tkn != "" && i.verifier != nil {
			tokenClaims, err = i.verifier.ValidateAccessToken(ctx, tkn)
			if err != nil {
				tokenClaims = nil
				log.Printf("token failed verification: %s", err.Error())
			}
		}

		cert, valid := rpcutil.CertFromServerContext(ctx)

		creds = &verifiedCreds{
			cert:        cert,
			certValid:   valid,
			token:       tkn,
			tokenClaims: tokenClaims,
		}
	}

	input := Attributes{
		Service:            service,
		Method:             method,
		Stream:             stream,
		Request:            req,
		CertificatePresent: creds.cert != nil,
		CertificateValid:   creds.certValid,
		Certificate:        creds.cert,
		TokenPresent:       creds.token != "",
		TokenValid:         creds.tokenClaims != nil,
		TokenClaims:        creds.tokenClaims,
	}

	queries, err := Validate(ctx, i.policy, input)
	addr := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr.String()
	}
	if err != nil {
		i.logger.Debug("request blocked by policy",
			zap.Any("attributes", input),
			zap.String("addr", addr),
			zap.Strings("queries", queries),
		)
	}
	return creds, err
}

type InterceptorOption func(interceptor *Interceptor)

func WithLogger(logger *zap.Logger) InterceptorOption {
	return func(interceptor *Interceptor) {
		interceptor.logger = logger
	}
}

func WithTokenVerifier(tv token.Validator) InterceptorOption {
	return func(interceptor *Interceptor) {
		interceptor.verifier = tv
	}
}

type verifiedCreds struct {
	cert        *x509.Certificate
	certValid   bool
	token       string
	tokenClaims *token.Claims
}

// if we want to get the request of a server-to-client streaming call from within an interceptor, we need a way to
// intercept the RecvMsg call.
// serverStreamInterceptor will run cb on all messages received through the stream. In a server-streaming RPC,
// the first message will be the request message. If cb returns a non-nil error, then that call to RecvMsg
// will return the error from cb.
type serverStreamInterceptor struct {
	grpc.ServerStream
	cb func(m any) error
}

func (ss *serverStreamInterceptor) RecvMsg(m any) error {
	err := ss.ServerStream.RecvMsg(m)
	if err != nil {
		return err
	}

	return ss.cb(m)
}
