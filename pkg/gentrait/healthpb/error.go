package healthpb

import (
	"fmt"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// ErrorCheck updates a health check based on a general error value.
type ErrorCheck struct {
	*checkBase
}

// newErrorCheck creates a new ErrorCheck for the given health check.
func newErrorCheck(c *gen.HealthCheck) (*ErrorCheck, error) {
	if err := validateErrorCheck(c); err != nil {
		return nil, err
	}
	return &ErrorCheck{checkBase: &checkBase{check: c}}, nil
}

func validateErrorCheck(c *gen.HealthCheck) error {
	// error checks shouldn't have bounds or a current value
	if c.GetCheck().GetBounds() != nil {
		return fmt.Errorf("bounds should be absent")
	}
	if c.GetCheck().GetCurrentValue() != nil {
		return fmt.Errorf("current_value should be absent")
	}
	return nil
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
