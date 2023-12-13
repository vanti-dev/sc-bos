package jsontypes

import (
	"fmt"
	"os"
	"strings"
)

// Password allows specifying a password either directly or via a password file.
type Password struct {
	Password     string `json:"password,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}

// Read returns the password, either from Password or PasswordFile.
func (c Password) Read() (string, error) {
	if c.Password != "" {
		return c.Password, nil
	}
	bs, err := os.ReadFile(c.PasswordFile)
	if err != nil {
		return "", fmt.Errorf("%w: %q", err, c.PasswordFile)
	}
	return strings.TrimSpace(string(bs)), nil
}
