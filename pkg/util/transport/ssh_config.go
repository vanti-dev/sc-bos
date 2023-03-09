package transport

// SshConfig is used to configure an Ssh transport
type SshConfig struct {
	ConnectionConfig
	// Ip is the IP address or hostname of the endpoint
	Ip string
	// Port is the port number to connect to, defaults to 22
	Port int

	Username string
	Password string

	IgnoreHostKey bool
	// todo: support host keys (and remove default IgnoreHostKey value below)
}

func (c *SshConfig) defaults() {
	c.IgnoreHostKey = true
	if c.Port == 0 {
		c.Port = 22
	}
	c.ConnectionConfig.defaults()
}
