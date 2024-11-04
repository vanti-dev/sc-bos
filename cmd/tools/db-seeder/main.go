package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/app/appconf"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/driver/alldrivers"
	mockcfg "github.com/vanti-dev/sc-bos/pkg/driver/mock/config"
	"github.com/vanti-dev/sc-bos/pkg/history/pgxstore"
	airqualitycfg "github.com/vanti-dev/sc-bos/pkg/zone/feature/airquality/config"
	occupancycfg "github.com/vanti-dev/sc-bos/pkg/zone/feature/occupancy/config"

	"github.com/vanti-dev/sc-bos/pkg/zone/allzones"
	meterscfg "github.com/vanti-dev/sc-bos/pkg/zone/feature/meter/config"
)

var (
	lookBack time.Duration
	dbUrl    string
	app      string
)

func init() {
	flag.DurationVar(&lookBack, "look-back", time.Hour*24*30*2, "amount of time to populate database history for starting from now, going backwards")
	flag.StringVar(&dbUrl, "db-url", "postgres://postgres:postgres@localhost:5432/smart_core", "database url")
	flag.StringVar(&app, "appconf", "app.conf.json", "app configuration file")
}

func main() {
	flag.Parse()

	appConf, err := appconf.LoadLocalConfig(path.Dir(app), path.Base(app))
	if err != nil {
		panic(err)
	}

	aqs, occs, meters, err := parseZoneConfig(appConf)
	if err != nil {
		panic(err)
	}

	devices, err := parseDeviceConfig(appConf)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	conf, err := pgxpool.ParseConfig(dbUrl)

	if err != nil {
		panic(err)
	}

	conf.MaxConns = 4
	conf.MinConns = 1

	db, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// before running the seeding tasks in parallel, make sure the DB exists to avoid error on first run:
	// CREATE TABLE IF NOT EXISTS duplicate key value violates unique constraint
	_, err = pgxstore.SetupStoreFromPool(ctx, "dummy", db)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}

	wg.Add(4)

	go func() {
		defer wg.Done()
		for _, aq := range aqs {
			err = SeedAirQuality(ctx, db, aq.Name, lookBack)
			if err != nil {
				panic(err)
			}
			fmt.Printf("seeded air quality zone %s\n", aq.Name)
		}
	}()

	go func() {
		defer wg.Done()
		for _, occ := range occs {
			err = SeedOccupancy(ctx, db, occ.Name, lookBack)
			if err != nil {
				panic(err)
			}
			fmt.Printf("seeded occupancycfg zone %s\n", occ.Name)
		}
	}()

	go func() {
		defer wg.Done()
		for _, mtr := range meters {
			err = SeedMeter(ctx, db, mtr.Name, lookBack)
			if err != nil {
				panic(err)
			}
			fmt.Printf("seeded meter zone %s\n", mtr.Name)
		}
	}()

	go func() {
		defer wg.Done()
		for _, d := range devices {
			err = SeedAirQuality(ctx, db, d, lookBack)
			if err != nil {
				panic(err)
			}
			fmt.Printf("seeded device %s\n", d)
		}
	}()

	wg.Wait()
}

func loadSystemConfig() (sysconf.Config, error) {
	systemConfig := sysconf.Default()

	systemConfig.ZoneFactories = allzones.Factories()
	systemConfig.DriverFactories = alldrivers.Factories()

	err := sysconf.Load(&systemConfig)
	return systemConfig, err
}

func parseZoneConfig(appConf *appconf.Config) ([]*airqualitycfg.Root, []*occupancycfg.Root, []*meterscfg.Root, error) {
	var aqs []*airqualitycfg.Root
	var occs []*occupancycfg.Root
	var meters []*meterscfg.Root

	for _, conf := range appConf.Zones {
		if conf.Type != "area" {
			continue
		}
		aq := airqualitycfg.Root{}
		occ := occupancycfg.Root{}
		mtr := meterscfg.Root{}

		buf, err := conf.MarshalJSON()
		if err != nil {
			return nil, nil, nil, err
		}

		err = json.Unmarshal(buf, &aq)
		if err != nil {
			return nil, nil, nil, err
		}

		aqs = append(aqs, &aq)

		err = json.Unmarshal(buf, &occ)
		if err != nil {
			return nil, nil, nil, err
		}

		occs = append(occs, &occ)

		err = json.Unmarshal(buf, &mtr)
		if err != nil {
			return nil, nil, nil, err
		}

		meters = append(meters, &mtr)

	}

	return aqs, occs, meters, nil
}

func parseDeviceConfig(appConf *appconf.Config) ([]string, error) {
	var airqualityDevices []string

	for _, dr := range appConf.Drivers {
		buf, err := dr.MarshalJSON()
		if err != nil {
			return nil, err
		}

		var devices mockcfg.Root
		err = json.Unmarshal(buf, &devices)
		if err != nil {
			return nil, err
		}

		for _, device := range devices.Devices {
			for _, trt := range device.Traits {
				if trt.Name == trait.AirQualitySensor.String() {
					airqualityDevices = append(airqualityDevices, device.Name)
				}
			}
		}
	}

	return airqualityDevices, nil
}
