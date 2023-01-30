package config

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MQTTBroker struct {
	Host         string `json:"host,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}

func (b MQTTBroker) ClientOptions() (*mqtt.ClientOptions, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(b.Host)
	opts.SetUsername(b.Username)
	opts.SetOrderMatters(false)
	return opts, nil
}

type MqttServiceSource struct {
	Source

	// the names to use for rpc requests
	RpcNames []string `json:"rpcNames,omitempty"`
}
