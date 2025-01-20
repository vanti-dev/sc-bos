package merge

import (
	"encoding/json"
	"testing"

	"github.com/smart-core-os/sc-api/go/traits"
)

// test the emergencyImpl alarmConfig for int values
//   - ok value good
//   - ok value bad
//   - okAbove good
//   - okAbove bad
func Test_emergency_int(t *testing.T) {
	impl := &emergencyImpl{}
	okValue := int64(234)

	tests := []struct {
		name       string
		alarmConf  AlarmConfig
		checkFunc  func(any) (*traits.Emergency, error)
		pointValue any
		wantLevel  traits.Emergency_Level
		wantReason string
		wantErr    bool
	}{
		{
			name: "ok value good",
			alarmConf: AlarmConfig{
				OkValue: &okValue,
			},
			checkFunc:  impl.checkIntValueForEmergency,
			pointValue: int32(234),
			wantLevel:  traits.Emergency_OK,
			wantReason: "",
			wantErr:    false,
		},
		{
			name: "ok value bad",
			alarmConf: AlarmConfig{
				OkValue:     &okValue,
				AlarmReason: "Sensor is not reading the expected value",
			},
			checkFunc:  impl.checkIntValueForEmergency,
			pointValue: int32(777),
			wantLevel:  traits.Emergency_EMERGENCY,
			wantReason: "Sensor is not reading the expected value",
			wantErr:    false,
		},
		{
			name: "okAbove good",
			alarmConf: AlarmConfig{
				OkAbove: &okValue,
			},
			checkFunc:  impl.checkIntValueForEmergency,
			pointValue: int32(235),
			wantLevel:  traits.Emergency_OK,
			wantReason: "",
			wantErr:    false,
		},
		{
			name: "okAbove bad",
			alarmConf: AlarmConfig{
				OkAbove:     &okValue,
				AlarmReason: "Sensor is not reading the expected value",
			},
			checkFunc:  impl.checkIntValueForEmergency,
			pointValue: int32(233),
			wantLevel:  traits.Emergency_EMERGENCY,
			wantReason: "Sensor is not reading the expected value",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl.config.AlarmConfig = &tt.alarmConf
			emergency, err := tt.checkFunc(tt.pointValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if emergency.Level != tt.wantLevel {
				t.Errorf("%s: got = %v, want %v", tt.name, emergency.Level, tt.wantLevel)
			}
			if emergency.Reason != tt.wantReason {
				t.Errorf("%s: got = %v, want %v", tt.name, emergency.Reason, tt.wantReason)
			}
		})
	}
}

// test the emergencyImpl alarmConfig for float values
//   - okAbove good
//   - okAbove bad
func Test_emergency_float(t *testing.T) {
	impl := &emergencyImpl{}
	okValue := int64(234)

	tests := []struct {
		name       string
		alarmConf  AlarmConfig
		checkFunc  func(any) (*traits.Emergency, error)
		pointValue any
		wantLevel  traits.Emergency_Level
		wantReason string
		wantErr    bool
	}{
		{
			name: "okAbove good",
			alarmConf: AlarmConfig{
				OkAbove: &okValue,
			},
			checkFunc:  impl.checkFloatValueForEmergency,
			pointValue: float32(235),
			wantLevel:  traits.Emergency_OK,
			wantReason: "",
			wantErr:    false,
		},
		{
			name: "okAbove bad",
			alarmConf: AlarmConfig{
				OkAbove:     &okValue,
				AlarmReason: "Sensor is not reading the expected value",
			},
			checkFunc:  impl.checkFloatValueForEmergency,
			pointValue: float32(233),
			wantLevel:  traits.Emergency_EMERGENCY,
			wantReason: "Sensor is not reading the expected value",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl.config.AlarmConfig = &tt.alarmConf
			emergency, err := tt.checkFunc(tt.pointValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if emergency.Level != tt.wantLevel {
				t.Errorf("%s: got = %v, want %v", tt.name, emergency.Level, tt.wantLevel)
			}
			if emergency.Reason != tt.wantReason {
				t.Errorf("%s: got = %v, want %v", tt.name, emergency.Reason, tt.wantReason)
			}
		})
	}
}

// test create a new emergencyImpl with both okValue and okAbove set
func Test_emergency_both(t *testing.T) {
	okValue := int64(234)
	okAbove := int64(235)

	emergencyCfg := emergencyConfig{
		AlarmConfig: &AlarmConfig{
			OkValue: &okValue,
			OkAbove: &okAbove,
		},
	}

	bytes, err := json.Marshal(emergencyCfg)
	if err != nil {
		t.Errorf("error marshalling emergencyCfg: %v", err)
	}

	_, err = readEmergencyConfig(bytes)
	if err != nil {
		if err.Error() != "cannot set both okValue and okAbove" {
			t.Errorf("cannot set both okValue and okAbove: %v", err)
		}
	}
}
