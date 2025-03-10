package pass

import "golang.org/x/crypto/bcrypt"

var ErrMismatchedHashAndPassword = bcrypt.ErrMismatchedHashAndPassword

// Compare compares a pass hash with a pass returning an error if they do not match.
func Compare(hash, secret []byte) error {
	return bcrypt.CompareHashAndPassword(hash, secret)
}
