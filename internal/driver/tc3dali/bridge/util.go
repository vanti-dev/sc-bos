package bridge

import (
	"context"
	"io"
)

func contextToCloser(cancel context.CancelFunc) io.Closer {
	return contextCloser(cancel)
}

type contextCloser context.CancelFunc

func (c contextCloser) Close() error {
	c()
	return nil
}
