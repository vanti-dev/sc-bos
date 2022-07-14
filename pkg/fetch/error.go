package fetch

import (
	"fmt"
)

type HTTPError struct {
	code        int
	description string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP Code %d: %s", e.code, e.description)
}
