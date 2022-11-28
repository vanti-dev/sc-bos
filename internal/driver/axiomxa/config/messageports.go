package config

type MessagePorts struct {
	AccessEvents *MessagePortStream
	FaultEvents  *MessagePortStream
}

type MessagePortStream struct {
	LocalAddress string `json:"localAddress,omitempty"` // [<host>]:<port>, as accepted by tcp.Listen
	Format       string `json:"format,omitempty"`       // Matches the MessagePort format configured in AxiomXa
}
