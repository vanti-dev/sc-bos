package pki

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"log"
	"os"
	"sync"

	"go.uber.org/multierr"
)

// Source defines a source of certificate information.
type Source interface {
	// Certs returns the currently active certificate and pool, loading it as needed.
	Certs() (cert *tls.Certificate, roots []*x509.Certificate, err error)

	// Note: We use []*x509.Certificate instead of x509.CertPool because CertPool does not allow you to get at the
	// certs once it's created, specifically for concatenation and encoding to PEM. Both of these things are useful to
	// us.
	// This does mean we are slightly less memory efficient as CertPool can represent unparsed certs, it's worth it though.
	// The alternative would be to use [][]byte, but we think that's a step too far.
}

// ErrNoCertOrErr is returned when a Source returns neither a tls.Certificate nor an error when invoked.
var ErrNoCertOrErr = errors.New("pki: no cert from source")

// Expiry is called to know if certificates need to be reloaded.
// Either cert is nil or err is nil but never both.
// roots may be empty.
//
// Expiry will not be called from concurrent go routines.
type Expiry func(cert *tls.Certificate, roots []*x509.Certificate, err error) bool

// CacheSource wraps source such that it is only ever called when expiry returns true.
// See the package pki/expiry for a collection of bundled expiration functions.
func CacheSource(source Source, expiry Expiry, opts ...CacheSourceOpt) Source {
	cs := &cachedSource{source: source, expiry: expiry}
	for _, opt := range opts {
		opt(cs)
	}
	return cs
}

type CacheSourceOpt func(cs *cachedSource)

// WithFSCache instructs the source to cache certs and roots to disk when source is invoked.
// When CacheSource Certs is called for the first time, an attempt will be made to load these files before invoking the
// underlying source.
func WithFSCache(certPath, rootsPath string, key PrivateKey) CacheSourceOpt {
	return func(cs *cachedSource) {
		cs.store = fsCacheStore{certPath: certPath, rootsPath: rootsPath, key: key}
	}
}

type cachedSource struct {
	source Source
	expiry Expiry

	mu    sync.Mutex // guards the following
	ran   bool
	cert  *tls.Certificate
	roots []*x509.Certificate
	err   error
	store cacheStore
}

func (c *cachedSource) Certs() (*tls.Certificate, []*x509.Certificate, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.ran && c.store != nil {
		c.ran = true
		c.cert, c.roots, c.err = c.store.Certs()
		if errors.Is(c.err, os.ErrNotExist) {
			// reset the error, not found if a fine starting place
			c.ran = false
			c.err = nil
		}
	}
	if c.expired() {
		cert, roots, err := c.source.Certs()
		if cert == nil && err == nil {
			roots = nil
			err = ErrNoCertOrErr
		}
		c.cert, c.roots, c.err = cert, roots, err
		// save the cert for future runs
		if err == nil && c.store != nil {
			err := c.store.Write(cert, roots)
			if err != nil {
				log.Printf("Unable to write certs to file %v", err)
			}
		}
	}
	return c.cert, c.roots, c.err
}

// expired must be called while holding mu.
func (c *cachedSource) expired() bool {
	if c.ran {
		return c.expiry(c.cert, c.roots, c.err)
	}
	c.ran = true
	return true
}

type cacheStore interface {
	Source
	Write(cert *tls.Certificate, roots []*x509.Certificate) error
}

type fsCacheStore struct {
	certPath, rootsPath string
	key                 PrivateKey
}

func (f fsCacheStore) Certs() (cert *tls.Certificate, roots []*x509.Certificate, err error) {
	return LoadCertAndRootsWithKey(f.certPath, f.rootsPath, f.key)
}

func (f fsCacheStore) Write(cert *tls.Certificate, roots []*x509.Certificate) error {
	if cert.PrivateKey != f.key {
		return errors.New("cert private key doesn't match")
	}
	var err error
	if f.certPath != "" {
		certPEM := EncodePEMSequence(cert.Certificate, "CERTIFICATE")
		err = multierr.Append(err, os.WriteFile(f.certPath, certPEM, 0664))
	}
	if f.rootsPath != "" {
		rootsPEM := EncodeCertificates(roots)
		err = multierr.Append(err, os.WriteFile(f.rootsPath, rootsPEM, 0664))
	}
	return err
}

// FuncSource adapts f to implement Source.
func FuncSource(f func() (cert *tls.Certificate, roots []*x509.Certificate, err error)) Source {
	return funcSource(f)
}

type funcSource func() (cert *tls.Certificate, roots []*x509.Certificate, err error)

func (f funcSource) Certs() (cert *tls.Certificate, roots []*x509.Certificate, err error) {
	return f()
}

// LazySource returns certs from f().Certs(), only invoking f the first time the returns Source.Certs is called.
func LazySource(f func() (Source, error)) Source {
	var once sync.Once
	var err error
	var source Source
	return FuncSource(func() (*tls.Certificate, []*x509.Certificate, error) {
		once.Do(func() {
			source, err = f()
		})
		if err != nil {
			return nil, nil, err
		}
		return source.Certs()
	})
}

// FSSource returns a Source that reads the cert+private keypair and roots from files in PEM format on the filesystem.
// Each call to Certs will read the files.
// rootsFile can be empty in which case no roots will be read or returned from Source.Certs.
// certFile should contain the leaf as the first entry followed by any intermediate certificates linking leaf with roots.
func FSSource(certFile, keyFile, rootsFile string) Source {
	return fsSource{certFile, keyFile, rootsFile}
}

type fsSource [3]string // [certFile, keyFile, rootsFile]

func (fsp fsSource) Certs() (*tls.Certificate, []*x509.Certificate, error) {
	certFile, keyFile, rootsFile := fsp[0], fsp[1], fsp[2]

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, err
	}

	if len(rootsFile) == 0 {
		return &cert, nil, nil
	}

	rootsPEM, err := os.ReadFile(rootsFile)
	if err != nil {
		return nil, nil, err
	}
	roots, err := ParseCertificatesPEM(rootsPEM)

	return &cert, roots, err
}

// FSKeySource is like FSSource but the key is already loaded.
// See LoadCertAndRootsWithKey.
func FSKeySource(certFile string, key PrivateKey, rootsFile string) Source {
	return fsKeySource{certFile, key, rootsFile}
}

type fsKeySource struct {
	certFile  string
	key       PrivateKey
	rootsFile string
}

func (fsp fsKeySource) Certs() (*tls.Certificate, []*x509.Certificate, error) {
	return LoadCertAndRootsWithKey(fsp.certFile, fsp.rootsFile, fsp.key)
}

// SelfSignedSource returns a Source backed by CreateSelfSignedCertificate and a basic certificate template.
func SelfSignedSource(key PrivateKey, opts ...CSROption) Source {
	return &ssSource{
		template: func() *x509.Certificate {
			return &x509.Certificate{
				Subject:               pkix.Name{CommonName: "localhost"},
				KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
				BasicConstraintsValid: true,
			}
		},
		key:  key,
		opts: opts,
	}
}

// SelfSignedSourceT is like SelfSignedSource but using the given template.
func SelfSignedSourceT(key PrivateKey, template *x509.Certificate, opts ...CSROption) Source {
	return &ssSource{
		template: func() *x509.Certificate {
			return template
		},
		key:  key,
		opts: opts,
	}
}

type ssSource struct {
	template func() *x509.Certificate
	key      PrivateKey
	opts     []CSROption
}

func (s *ssSource) Certs() (*tls.Certificate, []*x509.Certificate, error) {
	certDER, err := CreateSelfSignedCertificate(s.template(), s.key, s.opts...)
	if err != nil {
		return nil, nil, err
	}
	leaf, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, err
	}
	cert := &tls.Certificate{
		Certificate: [][]byte{certDER},
		Leaf:        leaf,
		PrivateKey:  s.key,
	}

	var roots []*x509.Certificate
	if leaf.IsCA {
		// assume that if we're making a self signed ca, then we want that ca to exist in the roots.
		roots = append(roots, leaf)
	}
	return cert, roots, nil
}

// SourceSet is a Source that will return certs from the first of sources to return a non-nil cert and nil err.
type SourceSet []Source

func (ss *SourceSet) Certs() (cert *tls.Certificate, roots []*x509.Certificate, err error) {
	if ss == nil || len(*ss) == 0 {
		return nil, nil, ErrNoCertOrErr
	}

	var rootErr error
	for _, source := range *ss {
		certs, roots, err := source.Certs()
		if certs == nil && err == nil {
			err = ErrNoCertOrErr
		}
		if err != nil {
			rootErr = multierr.Append(rootErr, err)
			continue
		}
		return certs, roots, nil
	}
	return nil, nil, rootErr
}

// Append adds source to c.
func (ss *SourceSet) Append(source Source) {
	*ss = append(*ss, source)
}

// Delete deletes source from c.
func (ss *SourceSet) Delete(source Source) {
	for i, item := range *ss {
		if source == item {
			*ss = append((*ss)[:i], (*ss)[i+1:]...)
		}
	}
}

// DirectSource adapts a tls.Certificate and some roots into a Source that always returns these values.
func DirectSource(cert *tls.Certificate, roots []*x509.Certificate) Source {
	return &directSource{cert, roots}
}

type directSource struct {
	cert  *tls.Certificate
	roots []*x509.Certificate
}

func (c *directSource) Certs() (cert *tls.Certificate, roots []*x509.Certificate, err error) {
	return c.cert, c.roots, nil
}

// AuthoritySource is like AuthoritySourceFn using a fixed id and key.
func AuthoritySource(authority Source, id *x509.Certificate, key PrivateKey, csrOpts ...CSROption) Source {
	return AuthoritySourceFn(authority, func() (*x509.Certificate, PrivateKey, error) {
		return id, key, nil
	}, csrOpts...)
}

// AuthoritySourceFn returns a Source that mints a new tls.Certificate based on the given signing authority using id and
// key funcs.
func AuthoritySourceFn(authority Source, keyPair func() (*x509.Certificate, PrivateKey, error), csrOpts ...CSROption) Source {
	return &authoritySource{
		authority: authority,
		csr:       newCSR(csrOpts...),
		keyPair:   keyPair,
	}
}

type authoritySource struct {
	authority Source
	csr       csr

	keyPair func() (*x509.Certificate, PrivateKey, error)
}

func (a *authoritySource) Certs() (cert *tls.Certificate, roots []*x509.Certificate, err error) {
	authority, roots, err := a.authority.Certs()
	if err != nil {
		return nil, nil, err
	}
	id, key, err := a.keyPair()
	if err != nil {
		return nil, nil, err
	}

	leafDer, err := createCertificate(a.csr, authority, id, key.Public())
	if err != nil {
		return nil, nil, err
	}

	// leaf should parse because we just created it
	leaf, _ := x509.ParseCertificate(leafDer)

	cert = &tls.Certificate{
		Certificate: make([][]byte, 0, len(authority.Certificate)+1),
		PrivateKey:  key,
		Leaf:        leaf,
	}

	// populate intermediates
	cert.Certificate = append(cert.Certificate, leafDer)
	cert.Certificate = append(cert.Certificate, authority.Certificate...)

	return cert, roots, nil
}
