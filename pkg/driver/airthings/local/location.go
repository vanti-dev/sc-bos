// Package local implements a local model of the AirThings api.
// The types in this package decouple the needs of Smart Core from any limitations in the AirThings api.
package local

import (
	"fmt"
	"sync"

	"github.com/olebedev/emitter"

	"github.com/smart-core-os/sc-bos/pkg/driver/airthings/api"
)

// Location holds information about devices in a particular location.
type Location struct {
	mu                    sync.RWMutex // guards the below
	latestSamplesByDevice map[string]api.DeviceSampleResponseEnriched
	bus                   *emitter.Emitter
}

func NewLocation() *Location {
	return &Location{
		latestSamplesByDevice: make(map[string]api.DeviceSampleResponseEnriched),
		bus:                   emitter.New(0),
	}
}

// GetLatestSample returns the latest sample for the given device, or false if there is none.
func (m *Location) GetLatestSample(deviceID string) (api.DeviceSampleResponseEnriched, bool) {
	m.mu.RLock()
	sample, ok := m.latestSamplesByDevice[deviceID]
	m.mu.RUnlock()
	return sample, ok
}

// PullLatestSamples subscribes to changes to the latest sample for the given device.
// Changes will be published to the returned chan, the current value will be returned immediately.
// Call the returned func to unsubscribe, the chan will be closed.
// It is safe to unsubscribe on a different goroutine.
func (m *Location) PullLatestSamples(deviceID string) (api.DeviceSampleResponseEnriched, <-chan api.DeviceSampleResponseEnriched, func()) {
	topic := fmt.Sprintf("sample/%s/change", deviceID)
	m.mu.RLock()
	latestSample, _ := m.latestSamplesByDevice[deviceID]
	stream := m.bus.On(topic)
	m.mu.RUnlock()

	ch := make(chan api.DeviceSampleResponseEnriched)
	off := func() {
		m.bus.Off(topic, stream)
		// don't close ch here, wait for the translation goroutine to close it,
		// this avoids races between send and close

		// There's a possible leak in the below go routine:
		// if the caller calls off() while there's still a sample to send on ch,
		// then the ch <- sample will block indefinitely.
		// We <-ch here to unblock that.
		//  - If ch got closed (because the goroutine returned) then all is good
		//  - If ch has a sample waiting, then we'll get it rather than the caller,
		//    which is also fine because they called off anyway.
		select {
		case <-ch:
		default:
		}
	}

	go func() {
		defer close(ch)
		for event := range stream {
			sample, ok := event.Args[0].(api.DeviceSampleResponseEnriched)
			if !ok {
				continue
			}
			ch <- sample
		}
	}()

	return latestSample, ch, off
}

// UpdateLatestSamples writes updates and notifies subscribers.
// The returned chan will be closed when all subscribers have been notified.
func (m *Location) UpdateLatestSamples(samples api.GetLocationSamplesResponseEnriched) <-chan struct{} {
	var emitChans []<-chan struct{}

	m.mu.Lock()
	for _, sample := range samples.GetDevices() {
		id := m.sampleDeviceID(sample)
		m.latestSamplesByDevice[id] = sample
		ch := m.emitSampleChange(sample)
		emitChans = append(emitChans, ch)
	}
	m.mu.Unlock()

	ch := make(chan struct{})
	go func() {
		for _, emitChan := range emitChans {
			<-emitChan
		}
		close(ch)
	}()
	return ch
}

func (m *Location) emitSampleChange(sample api.DeviceSampleResponseEnriched) <-chan struct{} {
	id := m.sampleDeviceID(sample)
	return m.bus.Emit(fmt.Sprintf("sample/%s/change", id), sample)
}

func (*Location) sampleDeviceID(sample api.DeviceSampleResponseEnriched) string {
	segment, ok := sample.GetSegmentOk()
	if !ok {
		return ""
	}
	return segment.GetId()
}
