package jsonapi

import (
	"net/http"
	"sync"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client

	Username, Password string

	tokenMu       sync.Mutex    // guards tokenReady assignment and close interactions
	tokenReady    chan struct{} // nil when not getting, blocks when getting, closed when token or error has been got
	loginStop     func()        // call this func to stop any running login requests
	loginResponse LoginResponse
	loginErr      error
}

func NewClient(baseURL, username, password string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{},
	}
}
