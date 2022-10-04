package tenant

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"go.uber.org/multierr"
	"sync"
)

type SecretSource interface {
	VerifySecret(ctx context.Context, secret string) (data SecretData, err error)
}

type SecretData struct {
	TenantID string
	Zones    []string
}

// SecretSourceFunc adapts an ordinary func to implement SecretSource.
type SecretSourceFunc func(ctx context.Context, secret string) (data SecretData, err error)

func (s SecretSourceFunc) VerifySecret(ctx context.Context, secret string) (data SecretData, err error) {
	return s(ctx, secret)
}

// FirstSuccessfulSecret implements SecretSource returning the first successful response from member SecretSource.
// Each SecretSource will be invoked in separate go routines in parallel.
// The first SecretSource to return a non-error will attempt to cancel the remaining SecretSource.VerifySecret invocations.
// If all members return errors then these will be combined and returned from this call.
type FirstSuccessfulSecret []SecretSource

func (s FirstSuccessfulSecret) VerifySecret(ctx context.Context, secret string) (data SecretData, err error) {
	if len(s) == 1 {
		return s[0].VerifySecret(ctx, secret)
	}

	success := make(chan SecretData, 1)
	errs := make([]error, len(s))

	ctx, cancel := context.WithCancel(ctx)
	// Technically we only need to cancel on success as returning due to error means all members have completed and the context is no longer in use.
	// But we do it on any return anyway as it doesn't really hurt anything and removes lint warnings :/
	defer cancel()

	allDone := make(chan struct{})
	var outstandingTasks sync.WaitGroup
	outstandingTasks.Add(len(s))

	for i, source := range s {
		i, source := i, source
		go func() {
			defer outstandingTasks.Done()
			data, err := source.VerifySecret(ctx, secret)
			if err != nil {
				errs[i] = err
				return
			}

			// We do this in a select with default in case more than one source succeeds.
			// The receiver on success only accepts 1 response and we don't want the go routine to block.
			select {
			case success <- data:
			default:
			}
		}()
	}

	go func() {
		outstandingTasks.Wait()
		close(allDone)
	}()

	select {
	case data := <-success:
		return data, nil
	case <-allDone:
		// We select on success again here in case there was a race between all tasks completing and one of them succeeding.
		// This can happen depending on the order go routines are woken up.
		// If all tasks complete, at least one of which succeeded, before this go routine is woken
		// then there's a race between this go routine and the outstandingTasks go routine above. If the later is
		// woken first then both the success chan and allDone chan are 'active' and there's no way to know
		// which case will be triggered in that situation. So we check them both.
		select {
		case data := <-success:
			return data, nil
		default:
		}

		return SecretData{}, multierr.Combine(errs...)
	}
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

func (s *MemorySecretStore) VerifySecret(_ context.Context, secret string) (data SecretData, err error) {
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
		if data.TenantID == clientID {
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
