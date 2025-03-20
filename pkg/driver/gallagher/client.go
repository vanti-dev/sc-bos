package gallagher

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	ApiKey     string
}

func newHttpClient(baseURL string, apiKey string, caPath string, certPath string, keyPath string) (*Client, error) {

	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, _ := tls.LoadX509KeyPair(certPath, keyPath)

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{clientCert},
				},
			},
		},
		ApiKey: apiKey,
	}, nil
}

func (c *Client) getUrl(p string) string {
	return fmt.Sprintf("%s/%s", c.BaseURL, p)
}

func (c *Client) doRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "GGL-API-KEY "+c.ApiKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: %d %s", resp.StatusCode, resp.Status)
	}
	return bytes, nil
}
