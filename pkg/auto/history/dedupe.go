package history

import (
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-golang/pkg/cmp"
)

type deduper[T proto.Message] struct {
	last  T
	equal cmp.Message
}

func newDeduper[T proto.Message](message cmp.Message) *deduper[T] {
	return &deduper[T]{equal: message}
}

// Changed checks if the message has changed compared to the last one.
// If the messages are equal, it returns false.
// If the messages are not equal, it returns true and updates the last message.
// It uses the provided equal comparator to determine equality, or a default one if not provided.
func (d *deduper[T]) Changed(m T) bool {
	// set a default equal comparator if not provided
	if d.equal == nil {
		d.equal = cmp.Equal()
	}

	if d.equal(d.last, m) {
		return false
	}

	d.last = m

	return true
}
