package hd2

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	// BaseURL is the root of the API.
	// e.g. https://1.2.3.4/api/
	BaseURL url.URL
	Client  *http.Client
}

// NewInsecureClient creates a Client that connects over HTTPS but does not verify the server certificate.
func NewInsecureClient(host string) *Client {
	return &Client{
		BaseURL: url.URL{
			Scheme: "https",
			Host:   host,
			Path:   "/api",
		},
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (c *Client) newRequest(method string, endpoint string) *http.Request {
	req := &http.Request{
		Method: method,
		URL:    c.BaseURL.JoinPath(endpoint),
		Header: make(http.Header),
	}
	return req
}

func handleResponse(res *http.Response, destPtr any) error {
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != 200 {
		return readError(res.Body)
	}
	rawJSON, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println("Got response " + string(rawJSON))
	return json.Unmarshal(rawJSON, destPtr)
}

type DeviceData struct {
	Entries []DeviceDataEntry
}

type DeviceDataEntry struct {
	Name string `json:"name"`
}

type AccessResponse struct {
	Rights string `json:"rights"`
}

func doGetRequest(conn *Client, target any, endpoint string) error {
	req := conn.newRequest("GET", endpoint)

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":Steinel123")))

	res, err := conn.Client.Do(req)
	if err != nil {
		return err
	}
	err = handleResponse(res, &target)
	return err
}
