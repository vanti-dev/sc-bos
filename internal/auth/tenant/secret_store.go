package tenant

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
)

type SecretSource interface {
	Verify(ctx context.Context, secret string) (data SecretData, err error)
}

type SecretStore interface {
	SecretSource
	Enroll(ctx context.Context, data SecretData) (secret string, err error)
	Invalidate(ctx context.Context, secret string) (present bool, err error)
	InvalidateClient(ctx context.Context, clientID string) error
}

// MemorySecretStore implements a primitive, in memory store for client secrets.
// A zero MemorySecretStore is ready to use as an empty store. Don't copy once accessed.
// In production, you'd want to store the secrets hashed in a database, so don't use this!
type MemorySecretStore struct {
	m     sync.RWMutex
	store map[string]SecretData // map from secret to associated data
}

func NewMemorySecretStore(m map[string]SecretData) *MemorySecretStore {
	copied := make(map[string]SecretData, len(m))
	for k, v := range m {
		copied[k] = v
	}

	return &MemorySecretStore{
		store: copied,
	}
}

func (s *MemorySecretStore) Enroll(_ context.Context, data SecretData) (secret string, err error) {
	s.m.Lock()
	defer s.m.Unlock()

	s.ensureInitialised()

	// the case where a duplicate secret is generated is vanishingly unlikely, but check for it anyway
	for {
		secret = genSecret()
		_, present := s.store[secret]
		if present {
			continue
		}

		s.store[secret] = data
		return
	}
}

func (s *MemorySecretStore) Verify(_ context.Context, secret string) (data SecretData, err error) {
	s.m.RLock()
	defer s.m.RUnlock()

	data, ok := s.store[secret]
	if !ok {
		err = errors.New("secret not found")
	}
	return
}

func (s *MemorySecretStore) Invalidate(_ context.Context, secret string) (present bool, err error) {
	s.m.Lock()
	defer s.m.Unlock()

	s.ensureInitialised()

	_, present = s.store[secret]
	delete(s.store, secret)
	return
}

func (s *MemorySecretStore) InvalidateClient(_ context.Context, clientID string) error {
	s.m.Lock()
	defer s.m.Unlock()

	for secret, data := range s.store {
		if data.ClientID == clientID {
			delete(s.store, secret)
		}
	}
	return nil
}

func (s *MemorySecretStore) ensureInitialised() {
	if s.store == nil {
		s.store = make(map[string]SecretData)
	}
}

type SecretData struct {
	ClientID string
}

func genSecret() (secret string) {
	// generate 256 bits of cryptographic-strength random data to use as the secret
	raw := make([]byte, 256/8)
	_, err := rand.Read(raw)
	if err != nil {
		// on a sane platform, random numbers should always be available
		panic(err)
	}

	// encode the secret with URL-safe base64 encoding
	var buffer bytes.Buffer
	encoder := base64.NewEncoder(base64.URLEncoding, &buffer)
	_, err = encoder.Write(raw)
	if err != nil {
		// the base64 encoder should work with any data
		panic(err)
	}
	err = encoder.Close()
	if err != nil {
		panic(err)
	}

	return buffer.String()
}
