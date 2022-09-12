package bridge

import (
	"fmt"
)

type InstanceType byte

type notification struct {
	Valid          bool         `tc3ads:"valid"`
	Sequence       uint64       `tc3ads:"sequence"`
	NumInputEvents uint8        `tc3ads:"nInputEvents"`
	InputEvents    []inputEvent `tc3ads:"inputEvents"`
}

func (n *notification) Decode() ([]InputEvent, error) {
	if !n.Valid {
		return nil, ErrInvalid
	}
	if int(n.NumInputEvents) > len(n.InputEvents) {
		return nil, ErrMalformed
	}
	rawEvents := n.InputEvents[:n.NumInputEvents]
	var events []InputEvent
	for _, rawEvent := range rawEvents {
		events = append(events, rawEvent.Decode())
	}
	return events, nil
}

type inputEvent struct {
	Parameters InputEventParameters `tc3ads:"parameters"`
	Error      bool                 `tc3ads:"error"`
	Status     uint32               `tc3ads:"status"`
	Message    string               `tc3ads:"message"`
	Data       uint16               `tc3ads:"data"`
}

func (e *inputEvent) Decode() InputEvent {
	var err error
	if e.Error {
		err = fmt.Errorf("%d: %s", e.Status, e.Message)
	}
	return InputEvent{
		InputEventParameters: e.Parameters,
		Err:                  err,
		Data:                 e.Data,
	}
}
