package dali

import (
	"errors"
	"fmt"
)

type Error struct {
	Status  uint32
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d) %s", e.Status, e.Message)
}

var (
	ErrCommandUnimplemented = errors.New("command unimplemented")
)
