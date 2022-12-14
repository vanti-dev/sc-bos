package jsonapi

import (
	"net/http"
	"sync"
)

type Client struct {
	baseURL    string
	httpClient *http.Client

	Username, Password string

	tokenMu       sync.Mutex    // guards tokenReady assignment and close interactions
	tokenReady    chan struct{} // nil when not getting, blocks when getting, closed when token or error has been got
	loginStop     func()        // call this func to stop any running login requests
	loginResponse LoginResponse
	loginErr      error
}
