package transport

// TcpConfig is used to configure a Tcp transport
type TcpConfig struct {
	ConnectionConfig
	// Ip is the IP address or hostname of the endpoint
	Ip string
	// Port is the port number to connect to
	Port int
}

func (c *TcpConfig) defaults() {
	c.ConnectionConfig.defaults()
}
