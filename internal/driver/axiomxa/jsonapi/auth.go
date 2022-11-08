package jsonapi

import (
	"context"
)

type LoginRequest struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type LoginResponse struct {
	Key string `json:"Key,omitempty"`
}

func (c *Client) Login(ctx context.Context) (LoginResponse, error) {
	var resBody LoginResponse
	err := c.postNoAuth(ctx, "/login", LoginRequest{Username: c.Username, Password: c.Password}, &resBody)
	return resBody, err
}
