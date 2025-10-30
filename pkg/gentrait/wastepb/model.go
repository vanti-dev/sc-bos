package wastepb

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var wasteTypes = [9]gen.WasteRecord_Type{
	gen.WasteRecord_TYPE_UNSPECIFIED,
	gen.WasteRecord_MIXED_RECYCLING,
	gen.WasteRecord_GENERAL_WASTE,
	gen.WasteRecord_ELECTRONICS,
	gen.WasteRecord_CHEMICAL,
	gen.WasteRecord_FOOD,
	gen.WasteRecord_PAPER,
	gen.WasteRecord_GLASS,
	gen.WasteRecord_PLASTIC,
}

var areas = [3]string{"Area 1", "Area 2", "Area 3"}
var systems = [3]string{"System 1", "System 2", "System 3"}
var streams = [3]string{"Stream 1", "Stream 2", "Stream 3"}

type Model struct {
	mu              sync.Mutex // guards allWasteRecords and genId
	allWasteRecords []*gen.WasteRecord
	genId           int

	lastWasteRecord *resource.Value // of *gen.WasteRecord
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.WasteRecord{})}
	opts = append(defaultOpts, opts...)

	m := &Model{
		lastWasteRecord: resource.NewValue(opts...),
	}

	// let's add some records to start with so we can test the list method without waiting
	startTime := time.Now().Add(-100 * time.Minute)
	for i := 0; i < 100; i++ {
		_, _ = m.GenerateWasteRecord(timestamppb.New(startTime))
		startTime = startTime.Add(time.Minute)
	}
	return m
}

// AddWasteRecord manually adds a waste record to the model
func (m *Model) AddWasteRecord(wr *gen.WasteRecord, opts ...resource.WriteOption) (*gen.WasteRecord, error) {
	v, err := m.lastWasteRecord.Set(wr, opts...)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.allWasteRecords = append(m.allWasteRecords, wr)
	m.genId++
	return v.(*gen.WasteRecord), nil
}

// GenerateWasteRecord generates a new waste record with the given timestamp and adds it to the model
func (m *Model) GenerateWasteRecord(ts *timestamppb.Timestamp) (*gen.WasteRecord, error) {

	wr := &gen.WasteRecord{
		WasteCreateTime:  ts,
		RecordCreateTime: ts,
		Id:               strconv.Itoa(m.genId),
		Weight:           rand.Float32() * 1000,
		Type:             wasteTypes[m.genId%len(wasteTypes)],
		Area:             areas[m.genId%len(areas)],
		System:           systems[m.genId%len(systems)],
		Stream:           streams[m.genId%len(streams)],
	}

	if m.genId%5 != 0 {
		// don't always add this info, it is not always available
		co2 := rand.Float32() * 100
		land := rand.Float32() * 100
		trees := rand.Float32() * 5

		wr.Co2Saved = &co2
		wr.LandSaved = &land
		wr.TreesSaved = &trees
	}

	return m.AddWasteRecord(wr)
}

func (m *Model) GetWasteRecordCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.allWasteRecords)
}

func (m *Model) ListWasteRecords(start, count int) []*gen.WasteRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	var wasteRecords []*gen.WasteRecord
	// reverse to retrieve the latest wasteRecords first
	for i := start - 1; i >= 0; i-- {
		wasteRecords = append(wasteRecords, m.allWasteRecords[i])
		if len(wasteRecords) >= count {
			break
		}
	}
	return wasteRecords
}

func (m *Model) pullWasteRecordsWrapper(request *gen.PullWasteRecordsRequest, server gen.WasteApi_PullWasteRecordsServer) error {
	if !request.UpdatesOnly {
		m.mu.Lock()
		i := len(m.allWasteRecords) - 50
		if i < 0 {
			i = 0
		}
		for ; i < len(m.allWasteRecords)-1; i++ {
			change := &gen.PullWasteRecordsResponse_Change{
				Name:       request.Name,
				NewValue:   m.allWasteRecords[i],
				ChangeTime: m.allWasteRecords[i].WasteCreateTime,
				Type:       types.ChangeType_ADD,
			}
			if err := server.Send(&gen.PullWasteRecordsResponse{Changes: []*gen.PullWasteRecordsResponse_Change{change}}); err != nil {
				m.mu.Unlock()
				return err
			}
		}
		m.mu.Unlock()
	}
	for change := range m.PullWasteRecords(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullWasteRecordsResponse{}
		msg.Changes = append(msg.Changes, change)
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) PullWasteRecords(ctx context.Context, opts ...resource.ReadOption) <-chan *gen.PullWasteRecordsResponse_Change {
	send := make(chan *gen.PullWasteRecordsResponse_Change)
	recv := m.lastWasteRecord.Pull(ctx, opts...)
	go func() {
		defer close(send)
		for change := range recv {
			wr := change.Value.(*gen.WasteRecord)
			change := &gen.PullWasteRecordsResponse_Change{
				Name:       "Waste Record",
				NewValue:   wr, // the mock driver only generates new waste records and does not delete them
				ChangeTime: wr.WasteCreateTime,
				Type:       types.ChangeType_ADD,
			}
			send <- change
		}
	}()

	return send
}
