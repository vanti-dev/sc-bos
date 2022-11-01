package bridge

import (
	"fmt"

	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
)

type notification struct {
	Valid          bool         `tc3ads:"valid"`
	Sequence       uint64       `tc3ads:"sequence"`
	NumInputEvents uint8        `tc3ads:"nInputEvents"`
	InputEvents    []inputEvent `tc3ads:"inputEvents"`
}

func (n *notification) Decode() ([]dali.InputEvent, error) {
	if !n.Valid {
		return nil, ErrInvalid
	}
	if int(n.NumInputEvents) > len(n.InputEvents) {
		return nil, ErrMalformed
	}
	rawEvents := n.InputEvents[:n.NumInputEvents]
	var events []dali.InputEvent
	for _, rawEvent := range rawEvents {
		events = append(events, rawEvent.Decode())
	}
	return events, nil
}

type inputEvent struct {
	Parameters dali.InputEventParameters `tc3ads:"parameters"`
	Error      bool                      `tc3ads:"error"`
	Status     uint32                    `tc3ads:"status"`
	Message    string                    `tc3ads:"message"`
	Data       uint16                    `tc3ads:"data"`
}

func (e *inputEvent) Decode() dali.InputEvent {
	var err error
	if e.Error {
		err = fmt.Errorf("%d: %s", e.Status, e.Message)
	}
	return dali.InputEvent{
		InputEventParameters: e.Parameters,
		Err:                  err,
		Data:                 e.Data,
	}
}
