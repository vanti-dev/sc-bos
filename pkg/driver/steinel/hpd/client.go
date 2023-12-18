package hpd

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	// BaseURL is the root of the API.
	// e.g. https://1.2.3.4/api/
	BaseURL  url.URL
	Client   *http.Client
	Password string `default:""`
}

// NewInsecureClient creates a Client that connects over HTTPS but does not verify the server certificate.
func NewInsecureClient(host string, password string) *Client {
	client := &Client{
		BaseURL: url.URL{
			Scheme: "https",
			Host:   host,
			Path:   "/rest",
		},
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
	if len(password) > 0 {
		client.Password = password
	} else {
		client.Password = "Steinel123"
	}
	return client
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
	// fmt.Println("Got response " + string(rawJSON))
	return json.Unmarshal(rawJSON, destPtr)
}

type DeviceData struct {
	Entries []DeviceDataEntry
}

type DeviceDataEntry struct {
	Name string `json:"name"`
}

type SensorResponse struct {
	SensorName                   string  `json:"SensorName"`
	SensorSWVersion              string  `json:"SensorSWVersion"`
	DetectorFWVersion            string  `json:"DetectorFWVersion"`
	Motion1                      bool    `json:"Motion1"`
	Presence1                    bool    `json:"Presence1"`
	TruePresence1                bool    `json:"TruePresence1"`
	Brightness1                  int     `json:"Brightness1"`
	Temperature                  float64 `json:"Temperature"`
	Humidity                     float64 `json:"Humidity"`
	VOC                          int     `json:"VOC"`
	CO2                          int     `json:"CO2"`
	AirPressure                  float64 `json:"AirPressure,omitempty"`
	Noise                        int     `json:"Noise,omitempty"`
	AerosolStaleAirStatus        int     `json:"AerosolStaleAirStatus,omitempty"`
	AerosolRiskOfInfectionStatus int     `json:"AerosolRiskOfInfectionStatus,omitempty"`
	ComfortZone                  bool    `json:"ComfortZone,omitempty"`
	DewPoint                     float64 `json:"DewPoint,omitempty"`
	AerosolStaleAir              int     `json:"AerosolStaleAir,omitempty"`
	AerosolRiskOfInfection       int     `json:"AerosolRiskOfInfection,omitempty"`
	ZonePeople0                  int     `json:"ZonePeople0,omitempty"`
	IAQ                          int     `json:"IAQ,omitempty"`
}

func doGetRequest(conn *Client, target any, endpoint string) error {
	req := conn.newRequest("GET", endpoint)

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+conn.Password)))

	res, err := conn.Client.Do(req)
	if err != nil {
		return err
	}
	err = handleResponse(res, &target)
	return err
}
