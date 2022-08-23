package tenant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vanti-dev/bsp-ew/internal/auth"
	"github.com/vanti-dev/bsp-ew/internal/util/rpcutil"
	"go.uber.org/zap"
)

type TokenServer struct {
	logger  *zap.Logger
	secrets SecretSource
	tokens  *TokenSource
}

func NewTokenSever(secrets SecretSource, name string, validity time.Duration, logger *zap.Logger) (*TokenServer, error) {
	key, err := generateKey()
	if err != nil {
		return nil, err
	}
	tokens := &TokenSource{
		Key:      key,
		Issuer:   name,
		Validity: validity,
		Now:      time.Now,
	}

	return &TokenServer{
		logger:  logger,
		secrets: secrets,
		tokens:  tokens,
	}, nil
}

func (s *TokenServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// time out operations after one minute
	ctx, cancel := context.WithTimeout(request.Context(), time.Minute)
	defer cancel()
	logger := rpcutil.HTTPLogger(request, s.logger)

	writer.Header().Set("Cache-Control", "no-store")

	// check and extract request parameters
	if request.Method != http.MethodPost {
		writeTokenError(writer, errInvalidRequest, logger)
		return
	}
	if err := request.ParseForm(); err != nil {
		writeTokenError(writer, errInvalidRequest, logger)
		logger.Error("form parse error", zap.Error(err))
		return
	}
	if request.PostForm.Get("grant_type") != "client_credentials" {
		writeTokenError(writer, errUnsupportedGrantType, logger)
		return
	}
	if !request.PostForm.Has("client_id") || !request.PostForm.Has("client_secret") {
		writeTokenError(writer, errInvalidRequest, logger)
	}
	clientId := request.PostForm.Get("client_id")
	clientSecret := request.PostForm.Get("client_secret")

	// lookup secret, and ensure it's for the matching client
	secretData, err := s.secrets.VerifySecret(ctx, clientSecret)
	if err != nil || secretData.TenantID != clientId {
		writeTokenError(writer, errInvalidClient, logger)
		return
	}
	logger = logger.With(zap.String("tenant", secretData.TenantID))

	// generate an access token for the client
	token, err := s.tokens.GenerateAccessToken(secretData)
	if err != nil {
		writeTokenError(writer, errors.New("failed to generate token"), logger)
		logger.Error("token generation failed", zap.Error(err))
		return
	}

	// send response to the client
	response := tokenSuccessResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.tokens.Validity.Seconds()),
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		writeTokenError(writer, errors.New("failed to marshal response"), logger)
		logger.Error("failed to marshal token response", zap.Error(err))
		return
	}

	_, err = writer.Write(responseBytes)
	if err != nil {
		logger.Error("failed to write response body", zap.Error(err))
	}
}

func (s *TokenServer) TokenValidator() auth.TokenValidator {
	return s.tokens
}

func writeTokenError(writer http.ResponseWriter, err error, logger *zap.Logger) {
	tokErr, ok := err.(tokenError)
	if !ok {
		tokErr = tokenError{
			Code:             500,
			ErrorName:        "internal",
			ErrorDescription: err.Error(),
		}
	}

	body, marshalErr := json.Marshal(tokErr)
	if marshalErr != nil {
		logger.Error("failed to marshal error response", zap.Error(err))
		body = []byte(`{"error": "internal"}`)
	}

	writer.WriteHeader(tokErr.Code)
	_, writeErr := writer.Write(body)
	if writeErr != nil {
		logger.Error("failed to write error body", zap.Error(err))
	}
}

type tokenSuccessResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // in seconds
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type tokenError struct {
	Code             int    `json:"-"` // http status code to use
	ErrorName        string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
}

func (te tokenError) Error() string {
	if te.ErrorDescription != "" {
		return fmt.Sprintf("%s: %s", te.ErrorName, te.ErrorDescription)
	} else {
		return te.ErrorName
	}
}

var (
	errInvalidRequest       = tokenError{Code: 400, ErrorName: "invalid_request"}
	errInvalidClient        = tokenError{Code: 401, ErrorName: "invalid_client"}
	errUnsupportedGrantType = tokenError{Code: 400, ErrorName: "unsupported_grant_type"}
)
