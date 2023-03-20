package remote

import (
	"context"
	"crypto/tls"
	"encoding/pem"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// Inspect connects to a remote node returning its public certs and metadata.
func Inspect(ctx context.Context, address string) (*gen.HubNodeInspection, error) {
	tlsConfig := &tls.Config{
		// We're actively trying to connect to the remote and fetch their cert
		InsecureSkipVerify: true,
	}

	// capture the server cert which we'll eventually return to the caller
	creds := &certInterceptor{TransportCredentials: credentials.NewTLS(tlsConfig)}
	conn, err := grpc.DialContext(ctx, address,
		// capture the tls cert
		grpc.WithTransportCredentials(creds),
		// block so we know if we've got the cert or not
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
	)
	if err != nil {
		return nil, err
	}

	out := &gen.HubNodeInspection{}

	// if the remote has certs then encode them as PEM in the response
	if certs, ok := creds.PeerCertificates(); ok && len(certs) > 0 {
		for _, cert := range certs {
			pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
			out.PublicCerts = append(out.PublicCerts, string(pemBytes))
		}
	}

	client := traits.NewMetadataApiClient(conn)
	md, err := client.GetMetadata(ctx, &traits.GetMetadataRequest{})
	if err != nil {
		return nil, err
	}
	out.Metadata = md

	return out, nil
}
