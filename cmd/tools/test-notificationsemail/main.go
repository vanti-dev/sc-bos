package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/notificationsemail"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/alert"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

var sampleNow = time.Date(2024, 01, 19, 0, 0, 0, 0, time.Local)

func addDummyAlerts(m *alert.Model) {

	sevs := []gen.Alert_Severity{
		gen.Alert_SEVERITY_UNSPECIFIED,
		gen.Alert_INFO,
		gen.Alert_WARNING,
		gen.Alert_SEVERE,
		gen.Alert_LIFE_SAFETY,
	}

	zones := []string{"East", "West"}

	for i := 0; i < 100; i++ {
		m.AddAlert(&gen.Alert{
			Id:          fmt.Sprintf("alert Id : %d", i),
			Severity:    sevs[i%5],
			Description: fmt.Sprintf("test: %ds", i),
			Zone:        zones[i%2],
			Floor:       fmt.Sprintf("%d", i/4),
			Source:      "manual",
			CreateTime:  timestamppb.New(time.Now().Add(-24 * 30 * time.Hour)),
			ResolveTime: nil,
		})
	}

	m.AddAlert(&gen.Alert{
		Id:          fmt.Sprintf("alert Id : %d", 555),
		Severity:    sevs[0],
		Description: fmt.Sprintf("test: %ds", 555),
		Zone:        zones[0],
		Floor:       fmt.Sprintf("%d", 1),
		Source:      "manual",
		CreateTime:  timestamppb.New(time.Now()),
		ResolveTime: timestamppb.Now(),
	})
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	root := node.New("testdevice01")

	m := alert.NewModel()
	addDummyAlerts(m)
	client := gen.WrapAlertApi(alert.NewModelServer(m))
	root.Announce(root.Name(), node.HasClient(client))
	alertApiRouter := gen.NewAlertApiRouter()
	alertApiRouter.AddAlertApiClient("testdevice01", client)

	root.Support(
		node.Routing(alertApiRouter), node.Clients(client),
	)

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

	lifecycle := notificationsemail.Factory.New(serv)
	defer lifecycle.Stop()
	cfg := `{
	"name": "emails", "type": "meteremail",
	"destination": {
	"host": "smtp.gmail.com",
	"from": "OCW Paradise Build <vantiocwdev@gmail.com>",
	"to": ["Dean Redfern <dean.redfern@vanti.co.uk>"],
	"passwordFile" : ".localpassword",
	"sendTime": "* * * * MON-FRI"
	},
	"subject" : "test alerts",
	"source" : "testdevice01"
}`

	_, err = lifecycle.Configure([]byte(cfg))
	if err != nil {
		panic(err)
	}
	_, err = lifecycle.Start()
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
}
