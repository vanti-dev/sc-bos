package pgxhub

import (
	"go.uber.org/zap"
)

type Option func(server *Server)

func WithLogger(logger *zap.Logger) Option {
	return func(server *Server) {
		server.logger = logger
	}
}
