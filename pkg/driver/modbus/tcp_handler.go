package modbus

import (
	"time"

	"github.com/goburrow/modbus"
	"go.uber.org/zap"
)

type TcpHandler struct {
	Handle
}

func NewTCPClientHandler(address string, opts ...TCPOption) *TcpHandler {
	handler := modbus.NewTCPClientHandler(address)

	for _, opt := range opts {
		opt(handler)
	}

	return &TcpHandler{handler}
}

// Connect the handler to the address supplied
func (t *TcpHandler) Connect() error {
	return t.Connect()
}

// Close  when done using the handler
func (t *TcpHandler) Close() error {
	return t.Close()
}

type TCPOption func(*modbus.TCPClientHandler)

func WithTCPTimeout(timeout time.Duration) TCPOption {
	return func(handler *modbus.TCPClientHandler) {
		handler.Timeout = timeout
	}
}

func WithTCPSlaveId(slaveId byte) TCPOption {
	return func(handler *modbus.TCPClientHandler) {
		handler.SlaveId = slaveId
	}
}

func WithTCPLogger(l *zap.Logger) TCPOption {
	return func(handler *modbus.TCPClientHandler) {
		var err error
		handler.Logger, err = zap.NewStdLogAt(l, zap.DebugLevel)

		if err != nil {
			l.Error("modbus tcp logger init failed", zap.Error(err))
		}
	}
}
