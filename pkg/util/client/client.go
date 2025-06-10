package client

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewConnection(conf Config) (*grpc.ClientConn, error) {
	tlsConfig := &tls.Config{}
	if conf.TLS.InsecureNoClientCert || conf.TLS.InsecureSkipVerify {
		if conf.TLS.InsecureSkipVerify {
			tlsConfig.InsecureSkipVerify = true
			tlsConfig.VerifyConnection = nil
		}
		if conf.TLS.InsecureNoClientCert {
			tlsConfig.Certificates = nil
			tlsConfig.GetClientCertificate = nil
		}
	}

	return grpc.NewClient(conf.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}
