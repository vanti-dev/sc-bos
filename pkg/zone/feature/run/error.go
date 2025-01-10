package run

import (
	"fmt"
)

// TagError returns a func that prefixes errors returned by fn with tag.
func TagError[T any](tag string, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		val, err := fn()
		if err != nil {
			err = fmt.Errorf("%s: %w", tag, err)
		}
		return val, err
	}
}
