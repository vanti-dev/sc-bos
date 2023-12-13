package hpd3

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func testCtx(t *testing.T) (context.Context, context.CancelFunc) {
	deadline, ok := t.Deadline()
	if ok {
		return context.WithDeadline(context.Background(), deadline)
	} else {
		return context.WithCancel(context.Background())
	}
}

func testLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}
