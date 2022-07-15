package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddress string
	deviceName    string
)

func init() {
	flag.StringVar(&serverAddress, "server", "localhost:23557", "address (host:port) of the Smart Core server")
	flag.StringVar(&deviceName, "device", "", "name of smart core device to pull publications from")
}

func run(ctx context.Context) (errs error) {
	flag.Parse()
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errs = multierr.Append(errs, err)
		return
	}

	defer func() {
		errs = multierr.Append(errs, conn.Close())
	}()

	pubClient := traits.NewPublicationApiClient(conn)
	stream, err := pubClient.PullPublications(ctx, &traits.PullPublicationsRequest{Name: deviceName})
	if err != nil {
		errs = multierr.Append(errs, err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			errs = multierr.Append(errs, err)
			return
		}

		for _, change := range res.Changes {
			logChange(change)
		}
	}
}

func logChange(change *traits.PullPublicationsResponse_Change) {
	t := change.ChangeTime.AsTime()

	switch change.Type {
	case types.ChangeType_REMOVE:
		fmt.Printf("[%s] Publication %q removed.\n", t.String(), change.OldValue.GetId())
	case types.ChangeType_ADD:
		fmt.Printf("[%s] Publication %q added:\n", t.String(), change.NewValue.GetId())
		logPublication(change.NewValue)
	case types.ChangeType_UPDATE, types.ChangeType_REPLACE:
		fmt.Printf("[%s] Publication %q updated:\n", t.String(), change.NewValue.GetId())
		logPublication(change.NewValue)
	}
}

func logPublication(pub *traits.Publication) {
	fmt.Printf("\tAudience: %q\n", pub.GetAudience())
	fmt.Printf("\tMedia Type: %q\n", pub.GetMediaType())
	fmt.Printf("\tVersion: %q\n", pub.GetVersion())
	body := pub.GetBody()
	fmt.Printf("\tBody (%d bytes):\n", len(body))

	bodyRunes := []rune(strings.ToValidUTF8(string(body), "."))
	for len(bodyRunes) > 0 {
		lineRunes := bodyRunes[:64]
		bodyRunes = bodyRunes[64:]

		fmt.Printf("\t\t%s\n", string(lineRunes))
	}
	fmt.Println()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errs := multierr.Errors(run(ctx))

	var code int
	switch len(errs) {
	case 0:
	case 1:
		_, _ = fmt.Fprintf(os.Stderr, "fatal error: %s\n", errs[0].Error())
		code = 1
	default:
		_, _ = fmt.Fprintln(os.Stderr, "fatal errors:")
		for _, err := range errs {
			_, _ = fmt.Fprintf(os.Stderr, "\t%s\n", err.Error())
		}
		code = 1
	}

	os.Exit(code)
}
