package hpd

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var wg sync.WaitGroup

type MockExportServer struct {
	ContextField context.Context
	Response     gen.PullExportMessagesResponse
}

type TestValues struct {
	tempCelcius float64
	humidity    float32
	co2         float32
	voc         float32
	people      int32
	state       traits.Occupancy_State
}

func (mes *MockExportServer) SetHeader(md metadata.MD) error {
	//TODO implement mes
	panic("implement mes")
}

func (mes *MockExportServer) SendHeader(md metadata.MD) error {
	//TODO implement mes
	panic("implement mes")
}

func (mes *MockExportServer) SetTrailer(md metadata.MD) {
	//TODO implement mes
	panic("implement mes")
}

func (mes *MockExportServer) Context() context.Context {
	return mes.ContextField
}

func (mes *MockExportServer) SendMsg(m any) error {
	//TODO implement mes
	panic("implement mes")
}

func (mes *MockExportServer) RecvMsg(m any) error {
	//TODO implement mes
	panic("implement mes")
}

func (mes *MockExportServer) Send(res *gen.PullExportMessagesResponse) error {
	mes.Response = *res
	return nil
}

// update the given resources with new values
func setNewAqValues(cancel context.CancelFunc, aq *resource.Value, newValues *TestValues) {
	defer wg.Done()
	// sleep a little before setting the value
	time.Sleep(1000 * time.Millisecond)

	aq.Set(&traits.AirQuality{
		CarbonDioxideLevel:       &newValues.co2,
		VolatileOrganicCompounds: &newValues.voc,
	})

	// sleep again so that we dont cancel the context before PullExportMessages has a chance to capture the changes
	time.Sleep(1000 * time.Millisecond)
	cancel()
}

func setNewThValues(cancel context.CancelFunc, t *resource.Value, newValues *TestValues) {
	defer wg.Done()
	// sleep a little before setting the value
	time.Sleep(1000 * time.Millisecond)
	t.Set(&traits.AirTemperature{
		Mode:               0,
		TemperatureGoal:    nil,
		AmbientTemperature: &types.Temperature{ValueCelsius: newValues.tempCelcius},
		AmbientHumidity:    &newValues.humidity,
		DewPoint:           nil,
	})

	// sleep again so that we dont cancel the context before PullExportMessages has a chance to capture the changes
	time.Sleep(1000 * time.Millisecond)
	cancel()
}

func setNewOccValues(cancel context.CancelFunc, o *resource.Value, newValues *TestValues) {
	defer wg.Done()
	// sleep a little before setting the value
	time.Sleep(1000 * time.Millisecond)

	o.Set(&traits.Occupancy{
		PeopleCount: newValues.people,
		State:       newValues.state,
	})

	// sleep again so that we dont cancel the context before PullExportMessages has a chance to capture the changes
	time.Sleep(1000 * time.Millisecond)
	cancel()
}

func Test_PullExportMessages_AirQuality(t *testing.T) {

	co2 := float32(0)
	voc := float32(0)
	humidity := float32(0)
	aq := resource.NewValue(resource.WithInitialValue(&traits.AirQuality{CarbonDioxideLevel: &co2, VolatileOrganicCompounds: &voc}), resource.WithNoDuplicates())
	o := resource.NewValue(resource.WithInitialValue(&traits.Occupancy{PeopleCount: 0, State: traits.Occupancy_OCCUPIED}), resource.WithNoDuplicates())
	temp := resource.NewValue(resource.WithInitialValue(&traits.AirTemperature{AmbientTemperature: &types.Temperature{ValueCelsius: 0}, AmbientHumidity: &humidity}), resource.WithNoDuplicates())

	server := NewUdmiServiceServer(nil, aq, o, temp, "prefix")

	req := &gen.PullExportMessagesRequest{
		Name: "test",
	}

	// create some new values that we check against at the end
	newValues := TestValues{
		co2:         123.4,
		voc:         345.6,
		humidity:    98.7,
		tempCelcius: 765.4,
		state:       traits.Occupancy_OCCUPIED,
		people:      459,
	}

	// before we call function under test, start a go routine which will change the values after a short delay
	myContext, cancel := context.WithCancel(context.Background())
	pullServer := MockExportServer{ContextField: myContext}
	wg.Add(1)
	go setNewAqValues(cancel, aq, &newValues)
	err := server.PullExportMessages(req, &pullServer)
	wg.Wait()

	if err != nil {
		// We are cancelling the context so this is expected
		if err.Error() != "context canceled" {
			t.Errorf("server.PullExportMessages(req, &pullServer) returned error")
		}
	}

	if pullServer.Response.Message == nil {
		t.Errorf("server response is nil")
	}

	// take the response payload which should be a valid PointsetEventMessage
	var pointSetMessage PointsetEventMessage
	err = json.Unmarshal([]byte(pullServer.Response.Message.Payload), &pointSetMessage)

	if err != nil {
		t.Errorf("json.Unmarshal failed")
	}

	if (pointSetMessage.Points.Co2Level.PresentValue != newValues.co2) ||
		(pointSetMessage.Points.VocLevel.PresentValue != newValues.voc) {
		t.Errorf("AirQuality trait does not match")
	}
}

func Test_PullExportMessages_TempHumidity(t *testing.T) {

	co2 := float32(0)
	voc := float32(0)
	humidity := float32(0)
	aq := resource.NewValue(resource.WithInitialValue(&traits.AirQuality{CarbonDioxideLevel: &co2, VolatileOrganicCompounds: &voc}), resource.WithNoDuplicates())
	o := resource.NewValue(resource.WithInitialValue(&traits.Occupancy{PeopleCount: 0, State: traits.Occupancy_OCCUPIED}), resource.WithNoDuplicates())
	temp := resource.NewValue(resource.WithInitialValue(&traits.AirTemperature{AmbientTemperature: &types.Temperature{ValueCelsius: 0}, AmbientHumidity: &humidity}), resource.WithNoDuplicates())

	server := NewUdmiServiceServer(nil, aq, o, temp, "prefix")

	req := &gen.PullExportMessagesRequest{
		Name: "test",
	}

	// create some new values that we check against at the end
	newValues := TestValues{
		co2:         123.4,
		voc:         345.6,
		humidity:    98.7,
		tempCelcius: 765.4,
		state:       traits.Occupancy_OCCUPIED,
		people:      459,
	}

	// before we call function under test, start a go routine which will change the values after a short delay
	myContext, cancel := context.WithCancel(context.Background())
	pullServer := MockExportServer{ContextField: myContext}
	wg.Add(1)
	go setNewThValues(cancel, temp, &newValues)
	err := server.PullExportMessages(req, &pullServer)
	wg.Wait()

	if err != nil {
		// We are cancelling the context so this is expected
		if err.Error() != "context canceled" {
			t.Errorf("server.PullExportMessages(req, &pullServer) returned error")
		}
	}

	if pullServer.Response.Message == nil {
		t.Errorf("server response is nil")
	}

	// take the response payload which should be a valid PointsetEventMessage
	var pointSetMessage PointsetEventMessage
	err = json.Unmarshal([]byte(pullServer.Response.Message.Payload), &pointSetMessage)

	if err != nil {
		t.Errorf("json.Unmarshal failed")
	}

	if (pointSetMessage.Points.Humidity.PresentValue != newValues.humidity) ||
		(pointSetMessage.Points.Temperature.PresentValue != newValues.tempCelcius) {
		t.Errorf("Temperature trait does not match")
	}
}

func Test_PullExportMessages_Occupancy(t *testing.T) {

	co2 := float32(0)
	voc := float32(0)
	humidity := float32(0)
	aq := resource.NewValue(resource.WithInitialValue(&traits.AirQuality{CarbonDioxideLevel: &co2, VolatileOrganicCompounds: &voc}), resource.WithNoDuplicates())
	o := resource.NewValue(resource.WithInitialValue(&traits.Occupancy{PeopleCount: 0, State: traits.Occupancy_OCCUPIED}), resource.WithNoDuplicates())
	temp := resource.NewValue(resource.WithInitialValue(&traits.AirTemperature{AmbientTemperature: &types.Temperature{ValueCelsius: 0}, AmbientHumidity: &humidity}), resource.WithNoDuplicates())

	server := NewUdmiServiceServer(nil, aq, o, temp, "prefix")

	req := &gen.PullExportMessagesRequest{
		Name: "test",
	}

	// create some new values that we check against at the end
	newValues := TestValues{
		co2:         123.4,
		voc:         345.6,
		humidity:    98.7,
		tempCelcius: 765.4,
		state:       traits.Occupancy_OCCUPIED,
		people:      459,
	}

	// before we call function under test, start a go routine which will change the values after a short delay
	myContext, cancel := context.WithCancel(context.Background())
	pullServer := MockExportServer{ContextField: myContext}
	wg.Add(1)
	go setNewOccValues(cancel, o, &newValues)
	err := server.PullExportMessages(req, &pullServer)
	wg.Wait()

	if err != nil {
		// We are cancelling the context so this is expected
		if err.Error() != "context canceled" {
			t.Errorf("server.PullExportMessages(req, &pullServer) returned error")
		}
	}

	if pullServer.Response.Message == nil {
		t.Errorf("server response is nil")
	}

	// take the response payload which should be a valid PointsetEventMessage
	var pointSetMessage PointsetEventMessage
	err = json.Unmarshal([]byte(pullServer.Response.Message.Payload), &pointSetMessage)

	if err != nil {
		t.Errorf("json.Unmarshal failed")
	}

	if (pointSetMessage.Points.PeopleCount.PresentValue != newValues.people) ||
		(pointSetMessage.Points.OccupancyState.PresentValue != newValues.state.String()) {
		t.Errorf("Occupancy trait does not match")
	}
}
