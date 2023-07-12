package statusalerts

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/statusalerts/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func analyseStatusLogs(ctx context.Context, source config.Source, c <-chan *gen.StatusLog, name string, client gen.AlertAdminApiClient, logger *zap.Logger) error {
	var failedLog *gen.StatusLog
	var failedCount int

	// a stopped timer
	retryTimer := newStoppedTimer()
	nextAttemptDelay := 200 * time.Millisecond
	var firstAttemptTime time.Time
	const nextAttemptScale = 1.2
	const maxAttemptDelay = 10 * time.Second

	recordResult := func(msg *gen.StatusLog, err error) {
		switch {
		case err == nil && failedLog == nil: // last attempt worked, this attempt worked too
		case err == nil && failedLog != nil:
			if failedCount > 5 {
				logger.Debug("alert saved successfully after previous attempt", zap.Int("attempts", failedCount))
			}

			failedLog = nil
			failedCount = 0
			if !retryTimer.Stop() {
				<-retryTimer.C
			}
			nextAttemptDelay = 200 * time.Millisecond
		case err != nil:
			if failedLog == nil {
				firstAttemptTime = time.Now()
			}
			failedLog = msg
			failedCount++
			// setup the next attempt to send the msg
			retryTimer.Reset(nextAttemptDelay)
			nextAttemptDelay = time.Duration(float64(nextAttemptDelay) * nextAttemptScale)
			if nextAttemptDelay > maxAttemptDelay {
				nextAttemptDelay = maxAttemptDelay
			}

			switch {
			case failedCount == 5:
				logger.Warn("failed to save alert, will retry", zap.Int("attempts", failedCount), zap.Error(err))
			case failedCount == 20:
				logger.Warn("failed to save alert, reducing logging", zap.Int("attempts", failedCount), zap.Error(err))
			case failedCount%100 == 0:
				logger.Debug("failed to save alert, will retry", zap.Int("attempts", failedCount), zap.Time("since", firstAttemptTime))
			}
		}
	}

	for {
		var msg *gen.StatusLog
		select {
		case <-retryTimer.C:
			msg = failedLog
		case m, ok := <-c:
			if !ok {
				return ctx.Err()
			}
			msg = m
		}

		switch {
		case msg.Level == gen.StatusLog_NOMINAL:
			_, err := client.ResolveAlert(ctx, &gen.ResolveAlertRequest{
				Name:         name,
				Alert:        &gen.Alert{Source: source.Name},
				AllowMissing: true,
			})
			recordResult(msg, err)
		default:
			_, err := client.CreateAlert(ctx, &gen.CreateAlertRequest{
				Name: name,
				Alert: &gen.Alert{
					Description: logToDescription(msg),
					Severity:    levelToSeverity(msg.Level),
					Floor:       source.Floor,
					Zone:        source.Zone,
					Source:      source.Name,
				},
				MergeSource: true,
			})
			recordResult(msg, err)
		}
	}
}

func newStoppedTimer() *time.Timer {
	t := time.NewTimer(0)
	if !t.Stop() {
		<-t.C
	}
	return t
}

func levelToSeverity(level gen.StatusLog_Level) gen.Alert_Severity {
	switch level {
	case gen.StatusLog_NOMINAL:
		return gen.Alert_SEVERITY_UNSPECIFIED
	case gen.StatusLog_NOTICE:
		return gen.Alert_INFO
	case gen.StatusLog_REDUCED_FUNCTION:
		return gen.Alert_WARNING
	case gen.StatusLog_NON_FUNCTIONAL, gen.StatusLog_OFFLINE:
		return gen.Alert_SEVERE
	default:
		return gen.Alert_WARNING
	}
}

func logToDescription(log *gen.StatusLog) string {
	return log.Description
}
