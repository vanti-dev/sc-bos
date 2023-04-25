package merge

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
)

var (
	ErrTraitNotSupported = errors.New("trait not supported")
)

type ErrReadProperty struct {
	Prop  string
	Cause error
}

func (e ErrReadProperty) PropName() string {
	return e.Prop
}

func (e ErrReadProperty) Error() string {
	return fmt.Sprintf("read property %q: %v", e.Prop, e.Cause)
}

func (e ErrReadProperty) Unwrap() error {
	return e.Cause
}

func LogPollError(logger *zap.Logger, msg string, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, context.Canceled) {
		return
	}

	logger.Warn(msg, zap.Stringer("err", stringerFunc(func() string {
		errs := multierr.Errors(err)
		var notFoundProps []string
		var unnamedNotFoundProps int
		var otherIssues []string
		for _, err := range errs {
			switch {
			case errors.Is(err, known.ErrNotFound):
				if prop, ok := err.(interface{ PropName() string }); ok {
					notFoundProps = append(notFoundProps, prop.PropName())
				} else {
					unnamedNotFoundProps++
				}
			default:
				otherIssues = append(otherIssues, err.Error())
			}
		}
		var out strings.Builder
		if totalNotFound := len(notFoundProps) + unnamedNotFoundProps; totalNotFound > 0 {
			fmt.Fprintf(&out, "%d not available", totalNotFound)
			if lenProps := len(notFoundProps); lenProps > 0 {
				sort.Strings(notFoundProps)
				fmt.Fprint(&out, " ")
				if lenProps < totalNotFound {
					fmt.Fprint(&out, "including ")
				}
				if lenProps > 6 {
					notFoundProps = append(notFoundProps[:6], "...")
				}
				for i, prop := range notFoundProps {
					if i > 0 {
						fmt.Fprint(&out, ",")
					}
					fmt.Fprintf(&out, "%s", prop)
				}
			}
		}
		if len(otherIssues) > 0 {
			if out.Len() > 0 {
				fmt.Fprint(&out, "; ")
			}
			fmt.Fprint(&out, strings.Join(otherIssues, "; "))
		}
		return out.String()
	})))
}

type stringerFunc func() string

func (s stringerFunc) String() string {
	return s()
}

func updatePollErrorStatus(statuses *statuspb.Map, name string, requests int, errs ...error) {
	problemName := fmt.Sprintf("%s.%s", name, "poll")

	allFailed := len(errs) == requests
	someOffline, allOffline := isOfflineError(errs...)

	if !someOffline {
		statuses.UpdateProblem(name, &gen.StatusLog_Problem{
			Name:        problemName,
			Level:       gen.StatusLog_NOMINAL,
			Description: fmt.Sprintf("poll success"),
		})
		return
	}

	level := gen.StatusLog_REDUCED_FUNCTION
	if allOffline && allFailed {
		level = gen.StatusLog_OFFLINE
	}
	statuses.UpdateProblem(name, &gen.StatusLog_Problem{
		Name:        problemName,
		Level:       level,
		Description: fmt.Sprintf("poll timeout"),
	})
}

func isOfflineError(errs ...error) (some, all bool) {
	offlineCount := 0
	for _, err := range errs {
		if errors.Is(err, context.DeadlineExceeded) {
			offlineCount++
		}
	}
	return offlineCount > 0, offlineCount == len(errs)
}
