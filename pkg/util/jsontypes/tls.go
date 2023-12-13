package jsontypes

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
)

// TLSConfig models a tls.Config as json.
// Call [Read] to convert it to a *tls.Config.
type TLSConfig struct {
	// IgnoreHub controls whether hub configured TLS settings are ignored.
	// If IgnoreHub is true then even if the client is enrolled with a hub, the hub provided TLS settings will be ignored.
	IgnoreHub bool `json:"ignoreHub,omitempty"`

	// These settings match their equivalents in [tls.Config].
	InsecureSkipVerify bool             `json:"insecureSkipVerify,omitempty"`
	Certificates       []TLSCertificate `json:"certificates,omitempty"`
	RootCAs            PEM              `json:"rootCAs,omitempty"`
}

// Read converts c to a tls.Config.
// Files read are resolved relative to base.
// If hubCfg is not nil - and IgnoreHub is false - then hubCfg is used as a base for the returned tls.Config.
func (c *TLSConfig) Read(base string, hubCfg *tls.Config) (*tls.Config, error) {
	if c == nil {
		return hubCfg, nil
	}
	var cfg *tls.Config
	if c.IgnoreHub || hubCfg == nil {
		cfg = &tls.Config{}
	} else {
		cfg = hubCfg.Clone()
	}
	if c.InsecureSkipVerify {
		cfg.InsecureSkipVerify = true
		cfg.VerifyConnection = nil
		cfg.GetConfigForClient = nil
	}
	if len(c.Certificates) > 0 {
		var certs []tls.Certificate
		for _, cert := range c.Certificates {
			c, err := cert.Read(base)
			if err != nil {
				return nil, err
			}
			// set c.Leaf to optimise handshakes
			if c.Leaf == nil && len(c.Certificate) > 0 {
				// guaranteed to succeed, any errors will be reported by cert.Read already
				c.Leaf, _ = x509.ParseCertificate(c.Certificate[0])
			}
			certs = append(certs, c)
		}
		cfg.Certificates = certs
		cfg.GetCertificate = nil
		cfg.GetClientCertificate = nil
		cfg.GetConfigForClient = nil
	}
	if len(c.RootCAs) > 0 {
		roots, err := c.RootCAs.Read(base)
		if err != nil {
			return nil, err
		}
		cfg.RootCAs = x509.NewCertPool()
		if !cfg.RootCAs.AppendCertsFromPEM(roots) {
			return nil, errors.New("failed to parse root CAs")
		}
		cfg.GetConfigForClient = nil
	}
	return cfg, nil
}

// TLSCertificate models a tls.Certificate as json.
type TLSCertificate struct {
	Certificate PEM `json:"certificate,omitempty"`
	PrivateKey  PEM `json:"privateKey,omitempty"`
}

// Read returns a tls.Certificate from the config in c.
// The certificate and private key will be validated as a matching pair, like tls.X509KeyPair.
func (c *TLSCertificate) Read(base string) (tls.Certificate, error) {
	cert, err := c.Certificate.Read(base)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("certificate %w", err)
	}
	key, err := c.PrivateKey.Read(base)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("privateKey %w", err)
	}
	return tls.X509KeyPair(cert, key)
}

// PEM represents one or more PEM encoded blocks.
// When encoded in json it is either a string of PEM encoded blocks, or a path to file containing PEM encoded blocks.
// See String.IsPath for how we tell the difference.
type PEM string

// Read reads the contents of p, the file or string.
// If p is a path, it is resolved relative to base.
func (p PEM) Read(base string) ([]byte, error) {
	if len(p) == 0 {
		return nil, errors.New("empty")
	}
	r, err := String(p).OpenBase(base)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
