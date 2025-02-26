package auto

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/transport"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

func getFloorName(f int) string {
	if f == 0 {
		return "GF"
	} else {
		return fmt.Sprintf("%2d", f)
	}
}

func getNextDestination(maxFloor int) string {
	f := 1 + rand.Int()%maxFloor
	return getFloorName(f)
}

func initTransport(load *float32) *gen.Transport {
	return &gen.Transport{
		ActualPosition: &gen.Transport_Location{
			Floor: "GF",
			Id:    "GF",
			Title: "Ground Floor",
		},
		NextDestinations: []*gen.Transport_Location{
			{
				Floor: "GF",
				Id:    "GF",
				Title: "Ground Floor",
			},
		},
		MovingDirection: gen.Transport_NO_DIRECTION,
		Load:            load,
		Doors: []*gen.Transport_Door{
			{
				Deck:   0,
				Title:  "Main Door",
				Status: gen.Transport_Door_CLOSED,
			},
		},
	}
}

func openCloseDoor(t *gen.Transport, nextDest string, load float32) {
	t.MovingDirection = gen.Transport_NO_DIRECTION
	if t.Doors[0].Status == gen.Transport_Door_CLOSED {
		t.Doors[0].Status = gen.Transport_Door_OPENING
	} else if t.Doors[0].Status == gen.Transport_Door_OPENING {
		t.Doors[0].Status = gen.Transport_Door_OPEN
	} else if t.Doors[0].Status == gen.Transport_Door_OPEN {
		t.Doors[0].Status = gen.Transport_Door_CLOSING
		t.Load = &load
	} else if t.Doors[0].Status == gen.Transport_Door_CLOSING {
		t.Doors[0].Status = gen.Transport_Door_CLOSED
		t.NextDestinations[0].Floor = nextDest
	}
}

func TransportAuto(model *transport.Model, maxFloor int) *service.Service[string] {
	slc := service.New(service.MonoApply(func(ctx context.Context, _ string) error {
		ticker := time.NewTicker(2 * time.Second)
		go func() {
			gfMaxWait := 3
			gfCurrentWait := 0
			currentFloor := 0
			load := float32(0)
			t := initTransport(&load)
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					// a basic simulation of a lift moving up and down.
					// if we are on the ground floor, choose a random floor between 1 and maxFloor and
					// and move to that floor, travelling at 1 floor per interval.
					if t.ActualPosition.Floor == "GF" && t.NextDestinations[0].Floor == "GF" {
						t.MovingDirection = gen.Transport_NO_DIRECTION
						if gfCurrentWait >= gfMaxWait {
							// we have waited long enough, move to a new floor
							load = float32(rand.Int() % 1000)
							openCloseDoor(t, getNextDestination(maxFloor), load)
						} else {
							gfCurrentWait++
						}
					} else if t.ActualPosition.Floor == t.NextDestinations[0].Floor { // we have arrived
						t.MovingDirection = gen.Transport_NO_DIRECTION
						gfCurrentWait = 0
						load = 0
						openCloseDoor(t, "GF", load)
					} else { // we are on the move,
						if t.NextDestinations[0].Floor == "GF" {
							t.MovingDirection = gen.Transport_DOWN
							currentFloor--
							t.ActualPosition.Floor = getFloorName(currentFloor)
						} else {
							t.MovingDirection = gen.Transport_UP
							currentFloor++
							t.ActualPosition.Floor = getFloorName(currentFloor)
						}
					}
					_, _ = model.UpdateTransport(t)
				}
			}
		}()
		return nil
	}), service.WithParser(func(data []byte) (string, error) {
		return string(data), nil
	}))
	_, _ = slc.Configure([]byte{})
	return slc
}
