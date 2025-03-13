// Command client-meter provides a CLI tool for interacting with the [gen.MeterApiClient].
package main

import (
	"context"
	"flag"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/util/client"
)

var clientConfig client.Config

func init() {
	flag.StringVar(&clientConfig.Endpoint, "endpoint", "localhost:23557", "smart core endpoint")
	flag.BoolVar(&clientConfig.Get, "get", true, "perform a get request")
	flag.BoolVar(&clientConfig.Pull, "pull", false, "pull changes")
	flag.StringVar(&clientConfig.Name, "name", "", "smart core name for requests")
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
	em := gen.NewMeterApiClient(conn)

	get := func(c context.Context, name string) error {
		log.Printf("GetMeterReading %s", name)
		res, err := em.GetMeterReading(ctx, &gen.GetMeterReadingRequest{Name: name})
		if err != nil {
			return err
		}
		log.Printf("got usage %v for %s", res.Usage, name)
		log.Printf("got production %v for %s", res.Produced, name)
		return nil
	}

	pull := func(c context.Context, name string) error {
		log.Printf("PullMeterReadings %s", name)
		stream, err := em.PullMeterReadings(ctx, &gen.PullMeterReadingsRequest{Name: name})
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
				log.Printf("got change for %s: %v", name, change.MeterReading.Usage)
			}
		}
	}

	grp, ctx := errgroup.WithContext(ctx)
	if clientConfig.Get {
		grp.Go(func() error {
			return get(ctx, clientConfig.Name)
		})
	}
	if clientConfig.Pull {
		grp.Go(func() error {
			return pull(ctx, clientConfig.Name)
		})
	}

	return grp.Wait()
}
