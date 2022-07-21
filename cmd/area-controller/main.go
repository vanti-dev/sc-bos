package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/pkg/pubcache"
	"go.etcd.io/bbolt"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cachePath        string
	serverAddress    string
	managementDevice string
	name             string
)

func init() {
	flag.StringVar(&name, "name", "test/area-controller-1", "Smart Core name of this area controller")
	flag.StringVar(&serverAddress, "server", "localhost:23557", "address (host:port) of the Smart Core server")
	flag.StringVar(&managementDevice, "management-device", "", "name of smart core device that manages this area controller")
	flag.StringVar(&cachePath, "cache", ".cache/area-controller.bolt", "path to cache database file")
}

func run(ctx context.Context) (errs error) {
	flag.Parse()

	// create DB dir if it doesn't exist
	dbDir := filepath.Dir(cachePath)
	err := os.MkdirAll(dbDir, 0750)
	if err != nil {
		errs = multierr.Append(errs, err)
	}

	// open database
	db, err := bbolt.Open(cachePath, 0640, nil)
	if err != nil {
		errs = multierr.Append(errs, err)
		return
	}
	pubStorage := pubcache.NewBoltStorage(db, []byte("publications"))

	// open connection to management server
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errs = multierr.Append(errs, err)
		return
	}

	defer func() {
		errs = multierr.Append(errs, conn.Close())
	}()

	pubClient := traits.NewPublicationApiClient(conn)
	pubID := fmt.Sprintf("%s:config", name)
	cache := pubcache.New(ctx, pubClient, managementDevice, pubID, pubStorage)

	for pub := range cache.Pull(ctx) {
		fmt.Printf("[%s] Publication %q update:\n", time.Now().String(), pubID)
		logPublication(pub)
	}

	return
}

func logPublication(pub *traits.Publication) {
	fmt.Printf("\tAudience: %q\n", pub.GetAudience())
	fmt.Printf("\tMedia Type: %q\n", pub.GetMediaType())
	fmt.Printf("\tVersion: %q\n", pub.GetVersion())
	body := pub.GetBody()
	fmt.Printf("\tBody (%d bytes):\n", len(body))

	bodyRunes := []rune(strings.ToValidUTF8(string(body), "."))
	for len(bodyRunes) > 0 {
		var lineRunes []rune
		if len(bodyRunes) >= 64 {
			lineRunes = bodyRunes[:64]
			bodyRunes = bodyRunes[64:]
		} else {
			lineRunes = bodyRunes
			bodyRunes = nil
		}

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
