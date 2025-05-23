package hikcentral

import (
	"testing"
	"time"
)

func Test_marshalUDMIPayload(t *testing.T) {
	msg := &CameraState{
		CamState:     true,
		CamFlt:       false,
		CamAim:       nil,
		CamOcc:       "",
		CamVideo:     "",
		CamStateTime: time.Date(2023, 6, 15, 17, 32, 0, 0, time.UTC),
		CamFltTime:   time.Time{},
	}
	want := `{` +
		`"CamFlt":{"present_value":false},` +
		`"CamFltTime":{"present_value":"0001-01-01T00:00:00Z"},` +
		`"CamState":{"present_value":true},` +
		`"CamStateTime":{"present_value":"2023-06-15T17:32:00Z"}` +
		`}`
	got, err := marshalUDMIPayload(msg)
	if err != nil {
		t.Fatalf("marshalUDMIPayload() error = %v", err)
	}
	if string(got) != want {
		t.Errorf("marshalUDMIPayload() got = %s, want %v", got, want)
	}
}
