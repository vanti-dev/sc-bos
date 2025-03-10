package modbusip

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
		handler.BaudRate = baudRate
	}
}

func WithRTUDataBits(dataBits int) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		handler.DataBits = dataBits
	}
}

func WithRTUStopBits(stopBits int) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
		handler.StopBits = stopBits
	}
}

func WithRTUParity(parity string) RTUOption {
	return func(handler *modbus.RTUClientHandler) {
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
