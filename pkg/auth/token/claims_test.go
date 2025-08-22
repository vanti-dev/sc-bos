package token

import (
	"encoding/json"
	"testing"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func TestResourceType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rt      ResourceType
		want    string
		wantErr bool
	}{
		{"zero", ResourceType(gen.RoleAssignment_RESOURCE_TYPE_UNSPECIFIED), `"RESOURCE_TYPE_UNSPECIFIED"`, false},
		{"value", ResourceType(gen.RoleAssignment_NAMED_RESOURCE), `"NAMED_RESOURCE"`, false},
		{"unknown", ResourceType(999), `999`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.rt)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ResourceType
		wantErr bool
	}{
		{"empty", "", ResourceType(0), true},
		{"string known", `"NAMED_RESOURCE"`, ResourceType(gen.RoleAssignment_NAMED_RESOURCE), false},
		{"string invalid", `"INVALID"`, ResourceType(0), true},
		{"int known", `1`, ResourceType(gen.RoleAssignment_NAMED_RESOURCE), false},
		{"int unknown", `999`, ResourceType(999), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rt ResourceType
			if err := json.Unmarshal([]byte(tt.data), &rt); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rt != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", rt, tt.want)
			}
		})
	}
}
