// Command test-statusemail tests the [statusemail] package, sending to a real email address.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/statusemail"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	root := node.New("test")

	var models []*statuspb.Model
	deviceCount := 100
	for i := 0; i < deviceCount; i++ {
		m := statuspb.NewModel()
		m.UpdateProblem(&gen.StatusLog_Problem{Name: "test", Level: gen.StatusLog_OFFLINE})
		models = append(models, m)
		client := node.WithClients(gen.WrapStatusApi(statuspb.NewModelServer(m)))
		root.Announce(fmt.Sprintf("device-%d", i), node.HasTrait(statuspb.TraitName, client),
			node.HasMetadata(&traits.Metadata{
				Appearance: &traits.Metadata_Appearance{Title: fmt.Sprintf("Device %d", i)},
				Location:   &traits.Metadata_Location{Floor: fmt.Sprintf("Floor %d", i%10), Zone: fmt.Sprintf("Room %d", i%5)},
				Membership: &traits.Metadata_Membership{Subsystem: "bms"},
			}))
	}

	serv := auto.Services{
		Logger: logger,
		Node:   root,
	}
	lifecycle := statusemail.Factory.New(serv)
	defer lifecycle.Stop()
	cfg := `{
  "name": "emails", "type": "statusemail",
  "discoverSources": true,
  "destination": {
    "host": "smtp.gmail.com",
    "from": "Enterprise Wharf <no-reply@enterprisewharf.co.uk>",
    "to": ["Matt Nathan <matt.nathan@vanti.co.uk>"],
    "passwordFile": ".secrets/ew-email-pass",
    "minInterval": "30s"
  }
}`
	_, err = lifecycle.Configure([]byte(cfg))
	if err != nil {
		panic(err)
	}
	_, err = lifecycle.Start()
	if err != nil {
		panic(err)
	}

	// this mimics how drivers work
	time.Sleep(2 * time.Second)
	for _, model := range models {
		model.UpdateProblem(&gen.StatusLog_Problem{Name: "test", Level: gen.StatusLog_NOMINAL})
	}

	levels := []gen.StatusLog_Level{
		gen.StatusLog_NOMINAL,
		gen.StatusLog_NOTICE,
		gen.StatusLog_REDUCED_FUNCTION,
		gen.StatusLog_NON_FUNCTIONAL,
		gen.StatusLog_OFFLINE,
	}
	for i := 0; i < 100; i++ {
		mi := rand.Int31n(int32(len(models)))
		m := models[mi]
		l := levels[rand.Int31n(int32(len(levels)))]
		n := fmt.Sprintf("device-%d", mi)
		log.Println("updating level for", n, "to", l)
		_, err := m.UpdateProblem(&gen.StatusLog_Problem{Name: "test", Level: l, Description: "test message"})
		if err != nil {
			panic(err)
		}
		d := time.Duration(rand.Int31n(10)) * time.Second
		time.Sleep(d)
	}
}
