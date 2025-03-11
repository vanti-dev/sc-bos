package modbus

import (
	"time"

	"github.com/goburrow/modbus"
	"go.uber.org/zap"
)

type RtuHandler struct {
	Handle
}

func NewRTUClientHandler(address string, opts ...RTUOption) *RtuHandler {
	handler := modbus.NewRTUClientHandler(address)

	for _, opt := range opts {
		opt(handler)
	}

	return &RtuHandler{handler}
}

type RTUOption func(*modbus.RTUClientHandler)

func (r *RtuHandler) Connect() error {
	return r.Connect()
}

func (r *RtuHandler) Close() error {
	return r.Close()
}

func WithRTUBaudRate(baudRate int) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		if baudRate == 0 {
			handler.BaudRate = 19_200
			return
		}
		handler.BaudRate = baudRate
	}
}

func WithRTUDataBits(dataBits int) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		if dataBits == 0 {
			handler.DataBits = 8
			return
		}
		handler.DataBits = dataBits
	}
}

func WithRTUStopBits(stopBits int) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		if stopBits == 0 {
			handler.StopBits = 1
			return
		}
		handler.StopBits = stopBits
	}
}

func WithRTUParity(parity string) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		if parity == "" {
			handler.Parity = "E"
			return
		}
		handler.Parity = parity
	}
}

func WithRTUTimeout(timeout time.Duration) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		handler.Timeout = timeout
	}
}
func WithRTUSlaveId(slaveId byte) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		handler.SlaveId = slaveId
	}
}

func WithRTULogger(l *zap.Logger) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		var err error
		handler.Logger, err = zap.NewStdLogAt(l, zap.DebugLevel)

		if err != nil {
			l.Error("modbus rtu logger init failed", zap.Error(err))
		}
	}
}
