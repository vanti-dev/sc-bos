package healthpb

import (
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// ErrorCheck updates a health check based on a general error value.
type ErrorCheck struct {
	checkBase
}

// NewErrorCheck creates a new ErrorCheck for the given health check.
func NewErrorCheck(c *gen.HealthCheck) *ErrorCheck {
	return &ErrorCheck{checkBase: checkBase{check: c}}
}

// UpdateError updates the health check state based on the given error.
// See [ErrorToProto] for the mapping from error to health check state.
func (c *ErrorCheck) UpdateError(err error) {
	c.UpdateErrorPb(ErrorToProto(err))
}

// UpdateErrorPb updates the health check state based on the given HealthCheck_Error.
func (c *ErrorCheck) UpdateErrorPb(err *gen.HealthCheck_Error) {
	c.write(func(dst *gen.HealthCheck) {
		check := dst.GetCheck()
		oldState := check.GetState()
		newState := gen.HealthCheck_Check_NORMAL
		if err != nil {
			newState = gen.HealthCheck_Check_ABNORMAL
			check.LastError = err
		}
		check.State = newState
		updateStateTimes(check, oldState, newState)
	})
}
