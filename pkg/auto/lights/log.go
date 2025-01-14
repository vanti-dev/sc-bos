package lights

import (
	"strings"
	"time"

	"go.uber.org/zap"
)

// logProcessStart logs the start of a processState; the reason it is running.
func logProcessStart(logger *zap.Logger, oldState, newState *ReadState, reasons ...string) {
	merge := func(items ...string) []string {
		if len(reasons) == 0 {
			return items
		}
		if len(items) == 0 {
			return reasons
		}
		return append(reasons, items...)
	}
	if oldState == nil {
		logger.Debug("processState starting", zap.Strings("reasons", merge("initial state")))
		return
	}
	changes := newState.ChangesSince(oldState)
	if len(changes) == 0 {
		logger.Debug("processState starting", zap.Strings("reasons", merge("no changes")))
		return
	}
	logger.Debug("processState starting", zap.Strings("reasons", merge(changes...)))
}

// logProcessComplete logs the completion of a processState; the changes made and why those changes.
func logProcessComplete(logger *zap.Logger, state *ReadState, writeState *WriteState, counts *actionCounts, duration, ttl time.Duration, err error) {
	switch {
	case err != nil:
	// the caller will already have logged the error
	case counts.TotalWrites > 0:
		logger.Debug("processState completed",
			zap.Strings("changes", counts.Changes()),
			zap.Strings("reasons", writeState.Reasons),
			zap.Stringer("duration", formattedDuration(duration)),
			zap.Stringer("ttl", formattedDuration(ttl)),
		)
	case ttl > 0 && state.Config.LogTTLDelays:
		logger.Debug("processState completed",
			zap.Strings("reasons", writeState.Reasons),
			zap.Stringer("duration", formattedDuration(duration)),
			zap.Stringer("ttl", formattedDuration(ttl)),
		)
	case state.Config.LogEmptyChanges:
		logger.Debug("processState completed",
			zap.Strings("reasons", writeState.Reasons),
			zap.Stringer("duration", formattedDuration(duration)),
			zap.Stringer("ttl", formattedDuration(ttl)),
		)
	}
}

// formatDuration returns a more human and log compatible duration where exact values aren't required.
func formatDuration(d time.Duration) string {
	d2 := d
	switch {
	case d > time.Hour:
		d2 = d.Round(time.Hour)
	case d > time.Minute:
		d2 = d.Round(time.Minute)
	case d > time.Second:
		d2 = d.Round(time.Second)
	case d > time.Millisecond:
		d2 = d.Round(time.Millisecond)
	case d > time.Microsecond:
		d2 = d.Round(time.Microsecond)
	}

	s := d2.String()
	if d2 != d {
		s = "~" + s
	}
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

type formattedDuration time.Duration

func (f formattedDuration) String() string {
	return formatDuration(time.Duration(f))
}
