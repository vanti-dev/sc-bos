package main

import (
	"context"
	"flag"
	"os"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/building"
)

var (
	flagConfigDir  string
	flagPopulateDB bool
)

func init() {
	flag.StringVar(&flagConfigDir, "config-dir", ".data/building-controller", "path to the configuration directory")
	flag.BoolVar(&flagPopulateDB, "populate-db", false, "inserts some test data into the database and exits")
}

func run(ctx context.Context) error {
	flag.Parse()
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	return building.RunController(ctx, logger, flagConfigDir, flagPopulateDB)
}

func main() {
	os.Exit(app.RunUntilInterrupt(run))
}
