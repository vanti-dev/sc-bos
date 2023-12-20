package se_wiser_knx

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	// BaseURL is the root of the API.
	// e.g. https://1.2.3.4/scada-remote
	BaseURL url.URL
	Client  *http.Client
}

// NewInsecureClient creates a Client that connects over HTTPS but does not verify the server certificate.
func NewInsecureClient(host string, username string, password string) *Client {
	return &Client{
		BaseURL: url.URL{
			Scheme: "https",
			User:   url.UserPassword(username, password),
			Host:   host,
			Path:   "/scada-remote",
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

func (c *Client) newRequest(method string, query url.Values) *http.Request {
	_url := c.BaseURL
	query.Add("m", "json")
	_url.RawQuery = query.Encode()
	req := &http.Request{
		Method: method,
		URL:    &_url,
		Header: make(http.Header),
	}
	return req
}

func handleResponse(res *http.Response, destPtr any) error {
	defer func() {
		_ = res.Body.Close()
	}()
	switch res.StatusCode {
	case 200: // continue
	case 401:
		return status.Error(codes.FailedPrecondition, "credentials are invalid")
	default:
		return readError(res.Body)
	}
	rawJSON, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawJSON) == 0 {
		// an empty response is not a valid json payload,
		// so we ignore it to avoid incorrect errors being reported
		return nil
	}
	return json.Unmarshal(rawJSON, destPtr)
}

func doPost(client *Client, query url.Values, target any) error {
	req := client.newRequest(http.MethodPost, query)
	log.Printf("POST %s", req.URL.String())
	res, err := client.Client.Do(req)
	if err != nil {
		return err
	}
	return handleResponse(res, target)
}

type Object struct {
	Id       int32       `json:"id"`
	Address  string      `json:"address"`
	Name     string      `json:"name"`
	Data     interface{} `json:"data"`
	DataType string      `json:"dataType"`
	Time     int32       `json:"time"`
	Date     string      `json:"date"`
	Comment  string      `json:"comment"`
}

func QueryObjects(client *Client) ([]Object, error) {
	var objects []Object
	query := make(url.Values)
	query.Add("r", "objects")
	err := doPost(client, query, &objects)
	return objects, err
}

func GetValue(client *Client, address string) (string, error) {
	query := url.Values{}
	query.Add("r", "grp")
	query.Add("fn", "read")
	query.Add("alias", address)
	var value string
	err := doPost(client, query, &value)

	if err != nil {
		return "", err
	}
	return value, nil
}

func SetValue(client *Client, address string, value any) error {
	query := url.Values{}
	query.Add("r", "grp")
	query.Add("fn", "write")
	query.Add("alias", address)
	query.Add("value", fmt.Sprintf("%s", value))
	var t bool
	err := doPost(client, query, &t)
	if err != nil {
		return err
	} else if !t {
		return status.Error(codes.Internal, "write failed")
	}
	return nil
}
