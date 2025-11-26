package sccexporter

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-bos/pkg/auto/sccexporter/config"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

type Device struct {
	Name string `json:"name"`
	// Data is a map of device data, the key is the trait name and the value is JSON data. Could be data or metadata depending on the message type
	Data map[trait.Name]string `json:"data,omitempty"`
}

type message struct {
	Agent     string    `json:"agent"`
	Device    Device    `json:"device"`
	Timestamp time.Time `json:"timestamp"`
}

type sccConnector struct {
	logger     *zap.Logger
	messagesCh chan message
	mqttCfg    config.Mqtt
	mqttClient mqtt.Client
}

func newSccConnector(logger *zap.Logger, mqttCfg config.Mqtt, client mqtt.Client) *sccConnector {
	s := &sccConnector{
		logger:     logger,
		messagesCh: make(chan message, 100),
		mqttCfg:    mqttCfg,
		mqttClient: client,
	}
	return s
}

// publishToScc listens on the messages channel and publishes messages to the SCC MQTT broker
// do not call this more than once per sccConnector instance
func (s *sccConnector) publishToScc(ctx context.Context) error {

	token := s.mqttClient.Connect()
	if !token.WaitTimeout(s.mqttCfg.ConnectTimeout.Duration) {
		return fmt.Errorf("timeout connecting to mqtt broker")
	}
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, stop publishing
			return nil
		case m, ok := <-s.messagesCh:
			if !ok {
				// Channel closed, stop publishing
				return nil
			}
			bytes, err := json.Marshal(m)
			if err != nil {
				s.logger.Error("failed to marshal scc payload", zap.Error(err))
				continue
			}
			t := s.mqttClient.Publish(s.mqttCfg.Topic, byte(*s.mqttCfg.Qos), false, bytes)
			if !t.WaitTimeout(s.mqttCfg.PublishTimeout.Duration) {
				s.logger.Warn("timeout publishing message for device", zap.Duration("timeout", s.mqttCfg.PublishTimeout.Duration))
				s.mqttClient.Disconnect(500)
				if token := s.mqttClient.Connect(); !token.WaitTimeout(s.mqttCfg.ConnectTimeout.Duration) {
					// this might be transient, so just log and try again next time
					s.logger.Error("failed to connect to mqtt broker", zap.Error(token.Error()))
				}
			}
		}
	}
}

// newMqttClient creates a new MQTT client with TLS configuration
// assumes the config contains valid paths with valid certs.
func newMqttClient(cfg config.Mqtt) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Host)
	opts.SetClientID(cfg.ClientId)
	opts.SetOrderMatters(false)

	cert, err := tls.LoadX509KeyPair(cfg.ClientCertPath, cfg.ClientKeyPath)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(cfg.CaCertPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	opts.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	return mqtt.NewClient(opts), nil
}
