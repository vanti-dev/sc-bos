// Command sample-tenant provides an example tenant application, including auth.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/smart-core-os/sc-api/go/traits"
)

var (
	flagInsecure     bool
	flagCACert       string
	flagTokenURL     string
	flagGRPCAddr     string
	flagClientId     string
	flagClientSecret string
)

func init() {
	flag.BoolVar(&flagInsecure, "insecure", false, "don't verify TLS certificates")
	flag.StringVar(&flagCACert, "ca", "", "path to Root CA certificate(s), PEM format X.509")
	flag.StringVar(&flagTokenURL, "token-url", "https://localhost:8443/oauth2/token", "URL of OAuth2 token endpoint")
	flag.StringVar(&flagGRPCAddr, "addr", "localhost:23557", "host:port of gRPC server to call")
	flag.StringVar(&flagClientId, "client-id", "", "OAuth2 client ID")
	flag.StringVar(&flagClientSecret, "client-secret", "", "OAuth2 client secret")
}

func main() {
	flag.Parse()
	tlsConfig := &tls.Config{}
	if flagInsecure {
		tlsConfig.InsecureSkipVerify = true
	} else if flagCACert != "" {
		pem, err := os.ReadFile(flagCACert)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: read %q: %s", flagCACert, err.Error())
			os.Exit(1)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		tlsConfig.RootCAs = pool
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	resp, err := httpClient.PostForm(flagTokenURL, url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {flagClientId},
		"client_secret": {flagClientSecret},
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: get token: %s\n", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: token issue failed")
		_, _ = io.Copy(os.Stderr, resp.Body)
		_, _ = fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: read response body: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("response body: %s\n", string(respBody))
	var parsed struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(respBody, &parsed)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: invalid JSON: %s\n", err.Error())
		os.Exit(1)
	}

	name := "light-1"
	conn, err := grpc.Dial(flagGRPCAddr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: can't connect: %s\n", err.Error())
		os.Exit(1)
	}
	client := traits.NewLightApiClient(conn)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+parsed.AccessToken)

	res, err := client.GetBrightness(ctx, &traits.GetBrightnessRequest{Name: name})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: call GetBrightness error: %s\n", err.Error())
	} else {
		fmt.Printf("Get response: %v\n", res)
	}

	res, err = client.UpdateBrightness(ctx, &traits.UpdateBrightnessRequest{Name: name, Brightness: &traits.Brightness{LevelPercent: 75}})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: call UpdateBrightness error: %s\n", err.Error())
	} else {
		fmt.Printf("Update response: %v\n", res)
	}
}
