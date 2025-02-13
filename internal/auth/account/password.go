package account

import (
	"github.com/vanti-dev/sc-bos/internal/util/pass"
)

const (
	minPasswordLength = 10
	maxPasswordLength = 100
)

func permitPassword(password string) bool {
	return len(password) >= minPasswordLength && len(password) <= maxPasswordLength
}

type PasswordHash []byte

func HashPassword(password string) (PasswordHash, error) {
	return pass.Hash([]byte(password))
}

func (h PasswordHash) Verify(password string) error {
	return pass.Compare(h, []byte(password))
}
