package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/vanti-dev/bsp-ew/pkg/testgen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	clientId  = "d0d14b7a-4bab-4cb6-a9a2-800af405cfae"
	authority = "https://login.microsoftonline.com/0551decc-7aac-47f9-9f75-06d0fc432ce1/"
	apiServer = "localhost:9000"
	message   = "Client was here!"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := run(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	secret, ok := os.LookupEnv("CLIENT_SECRET")
	if !ok {
		return errors.New("no CLIENT_SECRET environment variable provided")
	}

	cred, err := confidential.NewCredFromSecret(secret)
	if err != nil {
		return err
	}

	msalClient, err := confidential.New(clientId, cred, confidential.WithAuthority(authority))
	if err != nil {
		return err
	}

	scopes := []string{
		"api://10c0b693-5131-4d1f-840f-cbaafc4d72fb/.default",
	}

	result, err := msalClient.AcquireTokenByCredential(ctx, scopes)
	if err != nil {
		return err
	}

	if len(result.DeclinedScopes) != 0 {
		return fmt.Errorf("authentication server declined scopes %s", strings.Join(result.DeclinedScopes, ", "))
	}

	fmt.Printf("got a token %s\n", result.AccessToken)

	conn, err := grpc.DialContext(ctx, apiServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	testClient := testgen.NewTestApiClient(conn)
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+result.AccessToken)

	before, err := testClient.GetTest(ctx, &testgen.GetTestRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("current test data is %q\n", before.GetData())

	fmt.Printf("updating test data to %q\n", message)
	updated, err := testClient.UpdateTest(ctx, &testgen.UpdateTestRequest{
		Test: &testgen.Test{Data: message},
	})
	if err != nil {
		return err
	}
	fmt.Printf("updated test data is %q\n", updated.GetData())

	return nil
}
