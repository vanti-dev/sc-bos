package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"sc-bos-db-seeder/pkg/seeder/air_quality"
	"sc-bos-db-seeder/pkg/seeder/meter"
	"sc-bos-db-seeder/pkg/seeder/occupancy"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/app/appconf"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/driver/alldrivers"
	mockConfig "github.com/vanti-dev/sc-bos/pkg/driver/mock/config"
	airqualityConfig "github.com/vanti-dev/sc-bos/pkg/zone/feature/airquality/config"
	occupancyConfig "github.com/vanti-dev/sc-bos/pkg/zone/feature/occupancy/config"

	"github.com/vanti-dev/sc-bos/pkg/zone/allzones"
	metersConfig "github.com/vanti-dev/sc-bos/pkg/zone/feature/meter/config"
)

var (
	lookBack time.Duration
	dbUrl    string
)

func init() {
	flag.DurationVar(&lookBack, "look-back", time.Hour*24*30*2, "amount of time to populate database history for starting from now, going backwards")
	flag.StringVar(&dbUrl, "db-url", "postgres://postgres:postgres@localhost:5432/smart_core", "database url")
}

func main() {
	sysConf, err := loadSystemConfig()

	if err != nil {
		panic(err)
	}

	appConf, err := appconf.LoadLocalConfig(path.Dir(sysConf.AppConfig[0]), path.Base(sysConf.AppConfig[0]))
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

	db, err := pgxpool.NewWithConfig(ctx, conf)
	defer db.Close()

	wg := &sync.WaitGroup{}

	wg.Add(4)

	go func() {
		defer wg.Done()
		for _, aq := range aqs {
			err = air_quality.Seed(ctx, db, aq.Name, lookBack)
			if err != nil {
				panic(err)
			}
			fmt.Printf("seeded air quality zone %s\n", aq.Name)
		}
	}()

	go func() {
		defer wg.Done()
		for _, occ := range occs {
			err = occupancy.Seed(ctx, db, occ.Name, lookBack)
			if err != nil {
				panic(err)
			}
			fmt.Printf("seeded occupancy zone %s\n", occ.Name)
		}
	}()

	go func() {
		defer wg.Done()
		for _, mtr := range meters {
			err = meter.Seed(ctx, db, mtr.Name, lookBack)
			if err != nil {
				panic(err)
			}
			fmt.Printf("seeded meter zone %s\n", mtr.Name)
		}
	}()

	go func() {
		defer wg.Done()
		for _, d := range devices {
			err = air_quality.Seed(ctx, db, d, lookBack)
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

func parseZoneConfig(appConf *appconf.Config) ([]*airqualityConfig.Root, []*occupancyConfig.Root, []*metersConfig.Root, error) {
	var aqs []*airqualityConfig.Root
	var occs []*occupancyConfig.Root
	var meters []*metersConfig.Root

	for _, conf := range appConf.Zones {
		if conf.Type != "area" {
			continue
		}
		aq := airqualityConfig.Root{}
		occ := occupancyConfig.Root{}
		mtr := metersConfig.Root{}

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

		var devices mockConfig.Root
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
