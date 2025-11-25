package adapt

import (
	"context"
	"errors"
	"fmt"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/statuspb"
)

var (
	ErrNoDefault    = errors.New("no default adaptation for object")
	ErrNoAdaptation = errors.New("no adaptation from object to trait")
)

func updateRequestErrorStatus(statuses *statuspb.Map, name, request string, err error) {
	problemName := fmt.Sprintf("%s.%s", name, "request")
	switch {
	case err == nil:
		statuses.UpdateProblem(name, &gen.StatusLog_Problem{
			Name:        problemName,
			Level:       gen.StatusLog_NOMINAL,
			Description: fmt.Sprintf("request success %s", request),
		})
	case errors.Is(err, context.DeadlineExceeded):
		statuses.UpdateProblem(name, &gen.StatusLog_Problem{
			Name:        problemName,
			Level:       gen.StatusLog_OFFLINE,
			Description: fmt.Sprintf("timeout during %s", request),
		})
	}
}
