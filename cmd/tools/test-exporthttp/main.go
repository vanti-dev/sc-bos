package main

import (
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/exporthttp"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func main() {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	announcer := node.New("exporthttp-test")

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

	lifecycle := exporthttp.Factory.New(srv)
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

	// wait for all automations in exporthttp to finish
	time.Sleep(3 * time.Minute)
}

const (
	cfg = `
{
  "name":     "exporthttp-test",
  "type":     "exporthttp",
  "disabled": false,
  "baseUrl": "",
  "site":     "abc-test1",
  "authentication": {
	"type":       "Bearer",
	"secretFile": ""
  },
  "sources": {
	"occupancy":    {
	  "path": "occupancy",
	  "schedule": "0/1 * * * *",
	  "sensors":  ["pir/01","pir/02","pir/03"]
	},
	"temperature": {
	  "path": "temperature",
	  "schedule": "0/2 * * * *",
	  "sensors": ["FCU/01","FCU/02"]
	},
	"energy":       {
	  "path": "energy",
	  "schedule": "0/1 * * * *",
	  "meters": ["smart-core/meters/01"]
	},
	"airQuality":  {
	  "path": "air_quality",
	  "schedule": "0/1 * * * *",
	  "sensors": ["smart-core/iaq/01"]
	},
	"water": {
	  "path": "water",
	  "schedule": "0/1 * * * *",
	  "meters" : ["smart-core/meters/03"]
	}
  }
}`
)
