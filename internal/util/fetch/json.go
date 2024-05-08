package fetch

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
)

func JSON(ctx context.Context, url string, into any, options ...Option) error {
	o := resolveOpts(options...)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := o.httpClient
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {

				roots := x509.NewCertPool()

				// this loop performs the normal verification of all certs in the chain starting with the root
				for i := len(rawCerts) - 1; i >= 0; i-- {
					rawCert := rawCerts[i]
					c, _ := x509.ParseCertificate(rawCert)
					certItem, _ := x509.ParseCertificate(rawCert)

					if i == len(rawCerts)-1 {
						// this is the root cert, verify this using the system defaults, no custom Roots
						opts := x509.VerifyOptions{}
						if _, err := certItem.Verify(opts); err != nil {
							return err
						}
					} else {
						opts := x509.VerifyOptions{
							Roots: roots,
						}
						if _, err := certItem.Verify(opts); err != nil {
							return err
						}
					}

					roots.AddCert(c)
				}

				// now we verify the host name on the leaf, allowing for a wildcard which spans n-level subdomains
				leafCert, _ := x509.ParseCertificate(rawCerts[0])
				if err := leafCert.VerifyHostname(url); err != nil {
					if strings.HasPrefix(leafCert.Subject.CommonName, "*") {
						domainSuffix := leafCert.Subject.CommonName[1:]
						trimmedProtocol := strings.TrimPrefix(url, "https://")
						host, _, _ := net.SplitHostPort(trimmedProtocol)

						if !strings.HasSuffix(host, domainSuffix) {
							return err
						}
					} else {
						return err
					}
				}

				return nil
			},
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return HTTPError{response.StatusCode, response.Status}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, into)
}
