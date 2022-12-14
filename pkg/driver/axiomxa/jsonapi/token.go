package jsonapi

import (
	"context"
	"errors"
	"time"
)

var ErrUnauthorizedResponse = errors.New("401 Unauthorized")

const fetchTokenTimeout = 10 * time.Second

// withToken runs the given func passing the current auth token.
// If the func returns ErrUnauthorizedResponse then the token will be refreshed and the func will be run again.
// The passed context is used only when waiting for the token to become available.
// The func will not be called if getting the token results in an error, that error will be returned.
func (c *Client) withToken(ctx context.Context, do func(token string) error) error {
	token, err := c.token(ctx)
	if err != nil {
		return err
	}

	err = do(token)
	if errors.Is(err, ErrUnauthorizedResponse) {
		c.forgetToken()
		token, err = c.token(ctx)
		if err != nil {
			return err
		}

		err = do(token)
	}
	return err
}

// token attempts to login using Username and Password waiting until ctx is done for the login to complete.
// token may be called from multiple go routines at the same time, only one login request will be performed
func (c *Client) token(ctx context.Context) (string, error) {
	c.fetchTokenIfNeeded()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-c.tokenReady:
		return c.loginResponse.Key, c.loginErr
	}
}

func (c *Client) fetchTokenIfNeeded() {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	if c.tokenReady == nil {
		tokenReady := make(chan struct{})
		c.tokenReady = tokenReady
		// It shouldn't take more than this time to login
		ctx, stop := context.WithTimeout(context.Background(), fetchTokenTimeout)
		c.loginStop = stop
		go func() {
			defer stop()
			defer func() {
				c.tokenMu.Lock()
				defer c.tokenMu.Unlock()
				close(tokenReady)
			}()

			// we don't need to lock here because loginResponse and loginErr are only accessed when tokenReady is closed
			// which implies this go routine has returned
			c.loginResponse, c.loginErr = c.Login(ctx)
		}()
	}
}

func (c *Client) forgetToken() {
	c.tokenMu.Lock()
	c.tokenReady = nil
	if c.loginStop != nil {
		c.loginStop()
	}
	c.tokenMu.Unlock()

	// There is a non critical race here.
	// When we write to these fields because of a Login request we don't lock.
	// This is fine when accessing the token via c.token and related methods because they rely
	// on c.tokenReady instead which is guarded.
	// The worst case is we end up keeping the response or error in memory longer than needed
	c.loginResponse = LoginResponse{}
	c.loginErr = nil
}
