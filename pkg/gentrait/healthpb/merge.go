package healthpb

import (
	"github.com/vanti-dev/sc-bos/internal/health/healthmerge"
)

var (
	MergeChecks = healthmerge.Checks
	MergeCheck  = healthmerge.Check
	RemoveCheck = healthmerge.Remove
)
