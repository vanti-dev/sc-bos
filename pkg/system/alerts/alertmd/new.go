package alertmd

import (
	"github.com/smart-core-os/sc-bos/pkg/gen"
)

func New() *gen.AlertMetadata {
	return &gen.AlertMetadata{
		AcknowledgedCounts:   make(map[bool]uint32),
		FloorCounts:          make(map[string]uint32),
		ZoneCounts:           make(map[string]uint32),
		SeverityCounts:       make(map[int32]uint32),
		ResolvedCounts:       make(map[bool]uint32),
		NeedsAttentionCounts: make(map[string]uint32),
		SubsystemCounts:      make(map[string]uint32),
	}
}
