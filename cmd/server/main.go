package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/ew-config-poc/pkg/db"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	flagListenGRPC       string
	flagPostgresAddress  string
	flagPostgresUsername string
	flagPostgresPassword string
	flagPostgresDatabase string
	flagPopulateDatabase bool
)

func init() {
	flag.StringVar(&flagListenGRPC, "listen-grpc", "localhost:23557", "port to host gRPC server on")
	flag.StringVar(&flagPostgresAddress, "postgres-addr", "localhost:5432", "host:port to connect to postgres on")
	flag.StringVar(&flagPostgresUsername, "postgres-user", "postgres", "username for authenticating with postgres")
	flag.StringVar(&flagPostgresPassword, "postgres-password", "postgres", "password for authenticating with postgres")
	flag.StringVar(&flagPostgresDatabase, "postgres-db", "config_poc", "database name for connecting to postgres")
	flag.BoolVar(&flagPopulateDatabase, "populate-db", false, "inserts some test data into the database and exits")
}

func run(ctx context.Context) error {
	flag.Parse()

	group, ctx := errgroup.WithContext(ctx)

	dbConn, err := connectDB(ctx)
	if err != nil {
		return err
	}
	if flagPopulateDatabase {
		return populateDB(ctx, dbConn)
	}

	pubServer := &PublicationServer{conn: dbConn}

	group.Go(func() error {
		return serveGRPC(ctx, pubServer)
	})

	return group.Wait()
}

func serveGRPC(ctx context.Context, pubServer *PublicationServer) error {
	group, ctx := errgroup.WithContext(ctx)

	listener, err := net.Listen("tcp", flagListenGRPC)
	if err != nil {
		return fmt.Errorf("can't listen on %q: %w", flagListenGRPC, err)
	}

	server := grpc.NewServer()
	traits.RegisterPublicationApiServer(server, pubServer)
	reflection.Register(server)

	// serve gRPC on the listener
	group.Go(func() error {
		return server.Serve(listener)
	})

	// immediately attempt graceful shutdown when context cancelled
	stopped := make(chan struct{})
	group.Go(func() error {
		<-ctx.Done()
		server.GracefulStop()
		close(stopped)
		return nil
	})
	// force shutdown 5s after context cancelled
	group.Go(func() error {
		<-ctx.Done()
		log.Println("gRPC server will be force-closed in 5 seconds")
		select {
		case <-time.After(5 * time.Second):
			server.Stop()
		case <-stopped:
		}
		log.Println("gRPC server shut down")
		return nil
	})

	log.Printf("insecure gRPC server listening on %s", listener.Addr().String())

	return group.Wait()
}

func connectDB(ctx context.Context) (*pgx.Conn, error) {
	connectURL := url.URL{
		Scheme: "postgres",
		Host:   flagPostgresAddress,
		User:   url.UserPassword(flagPostgresUsername, flagPostgresPassword),
		Path:   "/" + flagPostgresDatabase,
	}
	conn, err := pgx.Connect(ctx, connectURL.String())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func populateDB(ctx context.Context, conn *pgx.Conn) error {
	deviceNames := []string{
		"vanti/floors/1/area-controller",
		"vanti/floors/2/area-controller",
		"vanti/floors/3/area-controller",
	}

	baseTime := time.Date(2022, 7, 6, 11, 18, 0, 0, time.UTC)

	err := conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		var errs error
		for _, name := range deviceNames {
			// register a publication
			id := name + ":config"
			errs = multierr.Append(errs, db.RegisterPublication(ctx, tx, id, name))

			// add some versions to it
			for i := 1; i <= 3; i++ {
				payload := struct {
					Device      string `json:"device"`
					Publication string `json:"publication"`
					Sequence    int    `json:"sequence"`
				}{
					Device:      name,
					Publication: id,
					Sequence:    i,
				}

				encoded, err := json.Marshal(payload)
				if err != nil {
					errs = multierr.Append(errs, err)
					continue
				}

				_, err = db.AddPublicationVersion(ctx, tx, db.PublicationVersion{
					PublicationID: id,
					PublishTime:   baseTime.Add(time.Duration(i) * time.Hour),
					Body:          encoded,
					MediaType:     "application/json",
					Changelog:     fmt.Sprintf("auto-populated revision %d", i),
				})
				errs = multierr.Append(errs, err)
			}
		}

		return errs
	})

	if err != nil {
		log.Printf("failed to populate database: %s", err.Error())
	} else {
		log.Println("database populated")
	}
	return err
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
