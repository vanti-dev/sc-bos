package http

import "go.uber.org/zap"

// logWrapper for go-retryablehttp Printf logging intercept
type logWrapper struct {
	*zap.Logger
}

func (lw *logWrapper) Printf(msg string, args ...any) {
	lw.Sugar().Debugf(msg, args...)
}
