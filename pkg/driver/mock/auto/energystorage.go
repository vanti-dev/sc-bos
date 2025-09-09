package auto

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystoragepb"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

func EnergyStorage(model *energystoragepb.Model) service.Lifecycle {
	type state int
	const (
		stateIdle state = iota
		stateCharging
		stateDischarging
	)
	states := []state{stateIdle, stateCharging, stateDischarging}

	s := service.New(service.MonoApply(func(ctx context.Context, _ string) error {
		go func() {
			timer := time.NewTimer(durationBetween(30*time.Second, 2*time.Minute))

			// Initialize with a random starting percentage
			currentPercentage := float32Between(20, 80)
			var currentState state = stateIdle
			var stateStartTime *timestamppb.Timestamp

			for {
				// Randomly change state occasionally
				if randomBool(0.3) { // 30% chance to change state
					newState := oneOf(states...)
					if newState != currentState {
						currentState = newState
						stateStartTime = timestamppb.Now()
					}
				}

				// Calculate percentage change rates
				var chargeRate, dischargeRate float32

				// Update percentage based on current state
				switch currentState {
				case stateIdle:
					// Slight random drift when idle
					currentPercentage += float32Between(-0.5, 0.5)
				case stateCharging:
					// Increase by 1.0-3.0% per update when charging (faster)
					chargeRate = float32Between(1.0, 3.0)
					currentPercentage += chargeRate
				case stateDischarging:
					// Decrease by 0.3-1.5% per update when discharging (slower)
					dischargeRate = float32Between(0.3, 1.5)
					currentPercentage -= dischargeRate
				}

				// Clamp percentage to realistic bounds
				if currentPercentage < 5 {
					currentPercentage = 5
					currentState = stateCharging // Auto-charge when low
					stateStartTime = timestamppb.Now()
				} else if currentPercentage > 95 {
					currentPercentage = 95
					currentState = stateIdle // Stop charging when full
				}

				// Build the energy level state
				energyLevel := &traits.EnergyLevel{
					Quantity: &traits.EnergyLevel_Quantity{
						Percentage:  currentPercentage,
						EnergyKwh:   currentPercentage * 0.75, // Assume 75kWh max capacity
						Descriptive: getDescriptiveThreshold(currentPercentage),
						DistanceKm:  currentPercentage * 4.5,                          // Assume ~450km max range
						Voltage:     ptr(getVoltageFromPercentage(currentPercentage)), // Simulate realistic voltage
					},
					PluggedIn: currentState == stateCharging || randomBool(0.7), // Usually plugged in
				}

				// Set flow state
				switch currentState {
				case stateIdle:
					energyLevel.Flow = &traits.EnergyLevel_Idle{
						Idle: &traits.EnergyLevel_Steady{
							StartTime: stateStartTime,
						},
					}
				case stateCharging:
					energyLevel.Flow = &traits.EnergyLevel_Charge{
						Charge: &traits.EnergyLevel_Transfer{
							StartTime: stateStartTime,
							Speed:     getSpeedFromRate(chargeRate),
							Target: &traits.EnergyLevel_Quantity{
								Percentage: float32Between(85, 100),
							},
						},
					}
				case stateDischarging:
					energyLevel.Flow = &traits.EnergyLevel_Discharge{
						Discharge: &traits.EnergyLevel_Transfer{
							StartTime:  stateStartTime,
							DistanceKm: float32Between(50, 200), // Trip distance
							Speed:      getSpeedFromRate(dischargeRate),
							Target: &traits.EnergyLevel_Quantity{
								Percentage: float32Between(5, 25), // Target low percentage when discharging
							},
						},
					}
				}

				_, _ = model.UpdateEnergyLevel(energyLevel)

				select {
				case <-ctx.Done():
					return
				case <-timer.C:
					timer = time.NewTimer(durationBetween(30*time.Second, 2*time.Minute))
				}
			}
		}()
		return nil
	}))
	_, _ = s.Configure([]byte(`""`))
	return s
}

func getDescriptiveThreshold(percentage float32) traits.EnergyLevel_Quantity_Threshold {
	switch {
	case percentage < 10:
		return traits.EnergyLevel_Quantity_CRITICALLY_LOW
	case percentage < 20:
		return traits.EnergyLevel_Quantity_LOW
	case percentage < 40:
		return traits.EnergyLevel_Quantity_MEDIUM
	case percentage < 80:
		return traits.EnergyLevel_Quantity_HIGH
	case percentage >= 95:
		return traits.EnergyLevel_Quantity_FULL
	default:
		return traits.EnergyLevel_Quantity_HIGH
	}
}

func getSpeedFromRate(rate float32) traits.EnergyLevel_Transfer_Speed {
	switch {
	case rate < 0.5:
		return traits.EnergyLevel_Transfer_EXTRA_SLOW
	case rate < 1.0:
		return traits.EnergyLevel_Transfer_SLOW
	case rate < 2.0:
		return traits.EnergyLevel_Transfer_NORMAL
	case rate < 2.5:
		return traits.EnergyLevel_Transfer_FAST
	default:
		return traits.EnergyLevel_Transfer_EXTRA_FAST
	}
}

func getVoltageFromPercentage(percentage float32) float32 {
	// Simulate realistic EV battery voltage curve (typical range 350-420V for high voltage batteries)
	// Voltage generally increases with charge level but not linearly
	minVoltage := float32(350.0)
	maxVoltage := float32(420.0)

	// Add some non-linear curve and randomness to make it realistic
	normalized := percentage / 100.0
	voltageRange := maxVoltage - minVoltage
	baseVoltage := minVoltage + (voltageRange * normalized)

	// Add some realistic variation (+/- 5V)
	variation := float32Between(-5.0, 5.0)
	return baseVoltage + variation
}
