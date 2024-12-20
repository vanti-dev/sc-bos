package devices

import (
	"time"
)

type Option func(*Server)

func WithNow(now func() time.Time) Option {
	return func(s *Server) {
		s.now = now
	}
}
