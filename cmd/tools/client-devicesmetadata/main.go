// Command client-devicesmetadata provides a CLI tool for interacting with the [gen.DevicesApiClient].
package main

import (
	"context"
	"flag"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/util/client"
)

var clientConfig client.Config

func init() {
	flag.StringVar(&clientConfig.Endpoint, "endpoint", "localhost:23557", "smart core endpoint")
	flag.BoolVar(&clientConfig.Get, "get", true, "perform a get request")
	flag.BoolVar(&clientConfig.Pull, "pull", false, "pull changes")
	flag.BoolVar(&clientConfig.TLS.InsecureNoClientCert, "insecure-no-client-cert", false, "")
	flag.BoolVar(&clientConfig.TLS.InsecureSkipVerify, "insecure-skip-verify", false, "")
}

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

func run() error {
	log.Printf("dialling: %s", clientConfig.Endpoint)
	conn, err := client.NewConnection(clientConfig)
	if err != nil {
		return err
	}
	log.Printf("dialled")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := gen.NewDevicesApiClient(conn)

	includes := &gen.DevicesMetadata_Include{Fields: []string{"metadata.membership.subsystem"}}
	get := func(c context.Context) error {
		log.Printf("GetDevicesMetadata")
		res, err := client.GetDevicesMetadata(ctx, &gen.GetDevicesMetadataRequest{Includes: includes})
		if err != nil {
			return err
		}
		log.Printf("got state %s", res)
		return nil
	}

	pull := func(c context.Context) error {
		log.Printf("PullDevicesMetadata")
		stream, err := client.PullDevicesMetadata(ctx, &gen.PullDevicesMetadataRequest{Includes: includes})
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		for {
			res, err := stream.Recv()
			if err != nil {
				return err
			}
			for _, change := range res.Changes {
				log.Printf("got change %s", change.DevicesMetadata)
			}
		}
	}

	grp, ctx := errgroup.WithContext(ctx)
	if clientConfig.Get {
		grp.Go(func() error {
			return get(ctx)
		})
	}
	if clientConfig.Pull {
		grp.Go(func() error {
			return pull(ctx)
		})
	}

	return grp.Wait()
}
