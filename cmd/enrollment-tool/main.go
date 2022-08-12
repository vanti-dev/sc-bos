package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/vanti-dev/bsp-ew/internal/enrollment"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	flagAddr        string
	flagCert        string
	flagKey         string
	flagCA          string
	flagName        string
	flagManagerName string
	flagManagerAddr string
)

func init() {
	flag.StringVar(&flagAddr, "addr", "localhost:23557", "address (host:port) of Area Controller gRPC server")
	flag.StringVar(&flagCert, "cert", "", "path to App Server intermediate CA certificate")
	flag.StringVar(&flagKey, "key", "", "path to App Server intermediate CA private key")
	flag.StringVar(&flagCA, "ca", "", "path to root CA")
	flag.StringVar(&flagName, "name", "", "name to assign to the enrolled node")
	flag.StringVar(&flagManagerName, "manager-name", "", "Smart Core name of the node's new manager")
	flag.StringVar(&flagManagerAddr, "manager-addr", "", "Address (host:port) of the node's new manager")
}

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	if flagCert == "" || flagKey == "" {
		return errors.New("cert and key options must be specified")
	}
	pair, err := tls.LoadX509KeyPair(flagCert, flagKey)
	if err != nil {
		return fmt.Errorf("load cert and key: %w", err)
	}
	pair.Leaf, err = x509.ParseCertificate(pair.Certificate[0])
	if err != nil {
		return err
	}

	rootsPEM, err := os.ReadFile(flagCA)
	if err != nil {
		return err
	}

	ca := &enrollment.CA{
		Certificate:   pair.Leaf,
		PrivateKey:    pair.PrivateKey,
		Intermediates: pair.Certificate[1:],
		Now:           time.Now,
		Validity:      30 * 24 * time.Hour,
	}
	en := &gen.Enrollment{
		TargetName:     flagName,
		TargetAddress:  flagAddr,
		ManagerName:    flagManagerName,
		ManagerAddress: flagManagerAddr,
		RootCas:        rootsPEM,
	}
	en, err = enrollment.EnrollAreaController(context.Background(), en, ca)
	fmt.Println(protojson.Format(en))
	return err
}
