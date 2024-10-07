package main

import (
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/wordpress"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func main() {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	announcer := node.New("wordpress-test")

	err = announceOccupancy(announcer, "pir/01", 8)
	if err != nil {
		panic(err)
	}
	err = announceOccupancy(announcer, "pir/02", 2)
	if err != nil {
		panic(err)
	}
	err = announceOccupancy(announcer, "pir/03", 4)
	if err != nil {
		panic(err)
	}
	err = announceAirQuality(announcer, "smart-core/iaq/01", 112.1)
	if err != nil {
		panic(err)
	}

	err = announceTemperature(announcer, "FCU/01", 19.1)
	if err != nil {
		panic(err)
	}
	err = announceTemperature(announcer, "FCU/02", 21.3)
	if err != nil {
		panic(err)
	}

	err = announceMeter(announcer, "smart-core/meters/01", "mWh", []float32{0, 1, 2, 12, 54, 100, 222, 654, 900, 1122, 1543})

	if err != nil {
		panic(err)
	}

	err = announceMeter(announcer, "smart-core/meters/03", "litres", []float32{0, 1, 11, 111, 222, 433, 566, 888, 1002, 1023, 2000})

	if err != nil {
		panic(err)
	}

	srv := auto.Services{
		Logger: logger,
		Node:   announcer,
		Now:    func() time.Time { return time.Now() },
	}

	lifecycle := wordpress.Factory.New(srv)
	_, err = lifecycle.Configure([]byte(cfg))

	if err != nil {
		panic(err)
	}

	_, err = lifecycle.Start()
	if err != nil {
		panic(err)
	}

	defer func() {
		_, err = lifecycle.Stop()
		if err != nil {
			panic(err)
		}
	}()

	// wait for all automations for wordpress to finish
	time.Sleep(15 * time.Second)
}

const (
	cfg = `
{
  "name":     "wordpress-test",
  "type":     "wordpress",
  "disabled": false,
  "baseUrl": "https://vanti-plugin-test-com.stackstaging.com/wp-json/recording/v1/create",
  "site":     "abc-test1",
  "authentication": {
	"type":       "Bearer",
	"secretFile": "./.data/secrets/wordpress"
  },
  "sources": {
	"occupancy":    {
	  "path": "occupancy",
	  "interval": "10s",
	  "sensors":  ["pir/01","pir/02","pir/03"]
	},
	"temperature": {
	  "path": "temperature",
	  "interval": "11s",
	  "sensors": ["FCU/01","FCU/02"]
	},
	"energy":       {
	  "path": "energy",
	  "interval": "10s",
	  "meters": ["smart-core/meters/01"]
	},
	"airQuality":  {
	  "path": "air_quality",
	  "interval": "10s",
	  "sensors": ["smart-core/iaq/01"]
	},
	"water": {
	  "path": "water",
	  "interval": "10s",
	  "meters" : ["smart-core/meters/03"]
	}
  }
}`
)
