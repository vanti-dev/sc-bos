package client

// Config represents the options for use with a simple Smart Core client
// this is useful for building small (CLI) tools for interacting with Smart Core endpoints
type Config struct {
	Endpoint string
	TLS      TLS

	Get  bool
	Pull bool
	Name string
}

type TLS struct {
	InsecureNoClientCert bool `json:"insecureNoClientCert,omitempty"` // don't present a client certificate when connecting to proxy servers
	InsecureSkipVerify   bool `json:"insecureSkipVerify,omitempty"`   // don't verify proxy server certificates
}
