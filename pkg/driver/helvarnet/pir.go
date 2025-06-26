package helvarnet

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/driver/helvarnet/config"
)

// Pir represents a single PIR sensor within the HelvarNet system.
type Pir struct {
	traits.UnimplementedOccupancySensorApiServer

	client    *tcpClient
	conf      *config.Device
	logger    *zap.Logger
	occupancy *resource.Value // *traits.Occupancy
}

func newPir(client *tcpClient, l *zap.Logger, conf *config.Device) *Pir {
	return &Pir{
		client:    client,
		conf:      conf,
		logger:    l,
		occupancy: resource.NewValue(resource.WithInitialValue(&traits.Occupancy{}), resource.WithNoDuplicates()),
	}
}

// refreshOccupancyStatus refreshes the occupancy by querying the input state of the sensor
func (p *Pir) refreshOccupancyStatus() error {
	command := queryInputState(p.conf.Address)
	want := "?" + command[1:len(command)-1]

	r, err := p.client.sendAndReceive(command, want)
	if err != nil {
		return err
	}

	split := strings.Split(r, "=")
	if len(split) < 2 {
		return fmt.Errorf("invalid response in getLastScene: %s", r)
	}

	s := strings.TrimSuffix(split[1], "#")
	statusInt, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	occupancy := traits.Occupancy_UNOCCUPIED
	if statusInt == 1 {
		occupancy = traits.Occupancy_OCCUPIED
	}

	// Update the occupancy status
	_, _ = p.occupancy.Set(&traits.Occupancy{
		State: occupancy,
	})

	return nil
}

func (p *Pir) GetOccupancy(context.Context, *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	value := p.occupancy.Get().(*traits.Occupancy)
	return value, nil
}

func (p *Pir) PullOccupancy(_ *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer) error {
	for value := range p.occupancy.Pull(server.Context()) {
		occupancy := value.Value.(*traits.Occupancy)
		err := server.Send(&traits.PullOccupancyResponse{Changes: []*traits.PullOccupancyResponse_Change{
			{
				Name:       p.conf.Name,
				ChangeTime: timestamppb.New(value.ChangeTime),
				Occupancy:  occupancy,
			},
		}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pir) runUpdateState(ctx context.Context, t time.Duration) error {

	err := p.refreshOccupancyStatus()
	if err != nil {
		p.logger.Error("failed to refresh occupancy status", zap.Error(err))
	}
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := p.refreshOccupancyStatus()
			if err != nil {
				p.logger.Error("failed to refresh occupancy status", zap.Error(err))
			}
		}
	}
}
