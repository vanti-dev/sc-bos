package pass

import "golang.org/x/crypto/bcrypt"

// Hash hashes a pass ready for storage.
// Calling Hash multiple times on the same pass will return different hashes.
func Hash(secret []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
}
