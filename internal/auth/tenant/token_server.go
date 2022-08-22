package tenant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func OAuth2TokenHandler(secrets SecretSource, source *TokenSource) http.Handler {
	if secrets == nil || source == nil {
		panic("parameters must be non-nil")
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// time out operations after one minute
		ctx, cancel := context.WithTimeout(request.Context(), time.Minute)
		defer cancel()

		writer.Header().Set("Cache-Control", "no-store")

		// check and extract request parameters
		if request.Method != http.MethodPost {
			writeTokenError(writer, errInvalidRequest)
			return
		}
		if err := request.ParseForm(); err != nil {
			writeTokenError(writer, errInvalidRequest)
			log.Printf("form parse error: %s", err.Error())
			return
		}
		if request.PostForm.Get("grant_type") != "client_credentials" {
			writeTokenError(writer, errUnsupportedGrantType)
			return
		}
		if !request.PostForm.Has("client_id") || !request.PostForm.Has("client_secret") {
			writeTokenError(writer, errInvalidRequest)
		}
		clientId := request.PostForm.Get("client_id")
		clientSecret := request.PostForm.Get("client_secret")

		// lookup secret, and ensure it's for the matching client
		secretData, err := secrets.Verify(ctx, clientSecret)
		if err != nil || secretData.ClientID != clientId {
			writeTokenError(writer, errInvalidClient)
			return
		}

		// generate an access token for the client
		token, err := source.GenerateAccessToken(secretData.ClientID, "sc-api")
		if err != nil {
			writeTokenError(writer, errors.New("failed to generate token"))
			log.Printf("failed to generate token for %q: %s", secretData.ClientID, err.Error())
			return
		}

		// send response to the client
		response := tokenSuccessResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   int(source.Validity.Seconds()),
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			writeTokenError(writer, errors.New("failed to marshal response"))
			log.Printf("failed to marshal token response for %q: %s", secretData.ClientID, err.Error())
			return
		}

		_, err = writer.Write(responseBytes)
		if err != nil {
			log.Printf("failed to write response body: %s", err.Error())
		}
	})
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

func writeTokenError(writer http.ResponseWriter, err error) {
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
		log.Printf("failed to marshal error response: %s", marshalErr.Error())
		body = []byte(`{"error": "internal"}`)
	}

	writer.WriteHeader(tokErr.Code)
	_, writeErr := writer.Write(body)
	if writeErr != nil {
		log.Printf("failed to write error body: %s", writeErr.Error())
	}
}
