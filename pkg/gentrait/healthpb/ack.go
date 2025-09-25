package healthpb

import (
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// AckBitset returns the bitwise OR of the ack transition enum numbers.
// This can be used to produce values suitable for [gen.HealthCheck.AckRequired] and [gen.HealthCheck.AckExpected].
func AckBitset(acks ...gen.HealthCheck_HealthChange) int32 {
	var ack int32
	for _, a := range acks {
		ack |= int32(a.Number())
	}
	return ack
}
