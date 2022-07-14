package auth

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/vanti-dev/ew-auth-poc/pkg/fetch"
	"go.uber.org/multierr"
)

var (
	ErrTokenNotSigned          = errors.New("token is not signed")
	ErrTokenMultipleSignatures = errors.New("token has multiple signatures")
	ErrKeyNotFound             = errors.New("signing key not found in key set")
	ErrUpdateTooSoon           = errors.New("trying to update too soon since last update")
)

const minimumUpdateInterval = time.Minute

type KeySet interface {
	VerifySignature(ctx context.Context, jws string) (payload []byte, err error)
}

type LocalKeySet struct {
	keys jose.JSONWebKeySet
}

func NewLocalKeySet(keys jose.JSONWebKeySet) *LocalKeySet {
	return &LocalKeySet{keys: keys}
}

func (ks *LocalKeySet) VerifySignature(_ context.Context, jws string) (payload []byte, err error) {
	sig, err := jose.ParseSigned(jws)
	if err != nil {
		return nil, err
	}
	return verifyJWSWithKeySet(sig, ks.keys)
}

func NewRemoteKeySet(background context.Context, url string) *RemoteKeySet {
	return &RemoteKeySet{
		url:        url,
		background: background,
	}
}

// RemoteKeySet handles verification of JSON Web Signatures based on public keys from a remote JWKS URL.
// If a verification of a signature fails because the signing key is unknown, the RemoteKeySet will automatically
// query the remote JWKS url for new keys.
type RemoteKeySet struct {
	url        string
	background context.Context

	m               sync.RWMutex
	cache           jose.JSONWebKeySet
	pending         *keySetFetchJob
	noUpdatesBefore time.Time
}

// VerifySignature will check that the provided JWS has a valid signature from a key included in this
// RemoteKeySet. Returns nil if the signature is valid, or a non-nil error otherwise.
// This function may make a network request to refresh the local cache of the remote key set, if the local cache
// cannot verify the token.
// It verifies only the signature - it does not verify any claims in the payload, or inspect the payload in any way!
func (ks *RemoteKeySet) VerifySignature(ctx context.Context, jws string) (payload []byte, err error) {
	sig, err := jose.ParseSigned(jws)
	if err != nil {
		return nil, err
	}

	ks.m.RLock()
	keySet := ks.cache
	ks.m.RUnlock()

	payload, err = verifyJWSWithKeySet(sig, keySet)
	if err == nil {
		return payload, nil
	} else if !errors.Is(err, ErrKeyNotFound) {
		// the JWS failed to verify, and it wasn't because of a missing key
		return nil, err
	}

	// The JWS failed to verify because the key we need is not cached; try getting it
	keySet, err = ks.updateKeys(ctx)
	if errors.Is(err, ErrUpdateTooSoon) {
		// We can't update at the moment, so just report that the key wasn't found.
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}

	// try verification again
	return verifyJWSWithKeySet(sig, keySet)
}

// updates the key cache, or waits for the update to complete if one is already in progress.
func (ks *RemoteKeySet) updateKeys(ctx context.Context) (jose.JSONWebKeySet, error) {
	// get the current fetch job, starting one if necessary
	ks.m.Lock()
	if time.Now().Before(ks.noUpdatesBefore) {
		ks.m.Unlock()
		return jose.JSONWebKeySet{}, ErrUpdateTooSoon
	}
	if ks.pending == nil {
		ks.pending = ks.startKeyFetchJob()
	}
	pending := ks.pending
	ks.m.Unlock()

	// wait for the job to complete
	return pending.wait(ctx)
}

func (ks *RemoteKeySet) startKeyFetchJob() *keySetFetchJob {
	done := make(chan struct{})
	job := &keySetFetchJob{done: done}

	go func() {
		var keySet jose.JSONWebKeySet
		err := fetch.JSON(ks.background, ks.url, &keySet)
		job.complete(keySet, err)

		ks.m.Lock()
		defer ks.m.Unlock()
		if err == nil {
			// save results in the cache
			ks.cache = keySet
		}
		// job is no longer pending
		ks.pending = nil
		ks.noUpdatesBefore = time.Now().Add(minimumUpdateInterval)
	}()

	return job
}

type keySetFetchJob struct {
	done   chan struct{}
	result jose.JSONWebKeySet
	err    error
}

func (job *keySetFetchJob) wait(ctx context.Context) (jose.JSONWebKeySet, error) {
	select {
	case <-ctx.Done():
		return jose.JSONWebKeySet{}, ctx.Err()
	case <-job.done:
		return job.result, job.err
	}
}

// must only be called once!
func (job *keySetFetchJob) complete(result jose.JSONWebKeySet, err error) {
	job.result = result
	job.err = err
	close(job.done)
}

// verifies that a JWS is signed by a key from a key set.
// The key ID recorded in the JWS must match one of the keys in the key set.
//
func verifyJWSWithKeySet(jws *jose.JSONWebSignature, keySet jose.JSONWebKeySet) (payload []byte, err error) {
	// check there is exactly one signature
	switch len(jws.Signatures) {
	case 0:
		return nil, ErrTokenNotSigned
	case 1: // this is fine
	default:
		return nil, ErrTokenMultipleSignatures
	}

	var errs error

	// search through all keys matching the key ID
	keyId := jws.Signatures[0].Header.KeyID
	for _, key := range keySet.Key(keyId) {
		payload, err = jws.Verify(key)
		if err == nil {
			return payload, nil
		}
		// record the reason this key failed to verify
		errs = multierr.Append(errs, err)
	}

	if errs != nil {
		// some key ID(s) matched, but failed to verify
		return nil, errs
	} else {
		// no keys matched
		return nil, ErrKeyNotFound
	}
}
