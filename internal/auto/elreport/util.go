package elreport

import (
	"sort"

	"github.com/vanti-dev/bsp-ew/pkg/gen"
)

// sorts (in numerical order) and deduplicates a list of faults.
// performed in place, returns the modified slice
func sortDeduplicateFaults(faults []gen.EmergencyLightFault) []gen.EmergencyLightFault {
	if len(faults) == 0 {
		return faults
	}

	sort.Slice(faults, func(i, j int) bool {
		return faults[i] < faults[j]
	})

	// the first element can't be a duplicate, it can stay where it is
	lastWritten := faults[0]
	writePos := 1

	// loop through the array, writing behind to remove duplicates
	// writePos can never get ahead of readPos, preventing data loss
	for readPos := 1; readPos < len(faults); readPos++ {
		if faults[readPos] == lastWritten {
			// skip duplicates
			continue
		}

		faults[writePos] = faults[readPos]
		lastWritten = faults[readPos]
		writePos++
	}

	// trim off the excess elements
	return faults[:writePos]
}

// performs a set comparison (order independent, count independent) between two slices of faults
func faultsEquivalent(a, b []gen.EmergencyLightFault) (equivalent bool) {
	aMap := make(map[gen.EmergencyLightFault]struct{}, len(a))
	bMap := make(map[gen.EmergencyLightFault]struct{}, len(b))
	for _, entry := range a {
		aMap[entry] = struct{}{}
	}
	for _, entry := range b {
		bMap[entry] = struct{}{}
	}

	for entry := range aMap {
		if _, present := bMap[entry]; !present {
			return false
		}
	}
	for entry := range bMap {
		if _, present := aMap[entry]; !present {
			return false
		}
	}
	return true
}
