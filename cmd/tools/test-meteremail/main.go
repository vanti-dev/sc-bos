// Command test-occupancyemail tests the [meteremail] package, sending to a real email address.
package main

import (
	"context"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/meteremail"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"
)

var sampleNow = time.Date(2024, 01, 9, 0, 0, 0, 0, time.Local)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	root := node.New("test")
	alltraits.AddSupport(root)

	now, _ := time.Parse(time.DateTime, "2023-11-15 11:36:00")
	now = now.Round(time.Second) // get rid of millis, etc

	now = sampleNow

	serv := auto.Services{
		Logger: logger,
		Node:   root,
		Now: func() time.Time {
			return now.Add(-2 * time.Second)
		},
	}
	lifecycle := meteremail.Factory.New(serv)
	defer lifecycle.Stop()
	cfg := `{
  "name": "emails", "type": "meteremail",
  "source": {
    "name": "test",
    "title": "One Centenary Way"
  },
  "destination": {
    "host": "smtp.gmail.com",
    "from": "OCW Paradise Build <deanredfern26@gmail.com>",
    "to": ["Dean Redfern <dean.redfern@vanti.co.uk>"],
	"passwordFile" : ".localpassword",	
    "sendTime": "0 0 * * MON-FRI"
  },
	"serverAddr" : "172.16.100.10:23557",
	"meterNames" : ["uk-ocw/floors/01/devices/CE1-electric-meter/WestDBA/T1LightingTotalEnergy",
					"uk-ocw/floors/01/devices/CE1-electric-meter/WestDBA/T1TotalLoadTotalEnergy"
					]
}`
	// vantiocwdev@gmail.com
	// __Vanti0cwdeV!__
	_, err = lifecycle.Configure([]byte(cfg))
	if err != nil {
		panic(err)
	}
	_, err = lifecycle.Start()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
}
