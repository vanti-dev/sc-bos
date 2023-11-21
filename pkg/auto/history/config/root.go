package config

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Root struct {
	auto.Config
	Source  *Source  `json:"source,omitempty"`
	Storage *Storage `json:"storage,omitempty"`
}

type Source struct {
	Name  string     `json:"name,omitempty"`
	Trait trait.Name `json:"trait,omitempty"`
	// ReadMask instructs the history service to only read the specified fields
	ReadMask *FieldMask `json:"readMask,omitempty"`
}

func (s Source) SourceName() string {
	return fmt.Sprintf("%s[%s]", s.Name, s.Trait)
}

type FieldMask fieldmaskpb.FieldMask

func (f *FieldMask) PB() *fieldmaskpb.FieldMask {
	if f == nil {
		return nil
	}
	return (*fieldmaskpb.FieldMask)(f)
}

func (f *FieldMask) UnmarshalJSON(bytes []byte) error {
	return protojson.Unmarshal(bytes, (*fieldmaskpb.FieldMask)(f))
}

func (f *FieldMask) MarshalJSON() ([]byte, error) {
	return protojson.Marshal((*fieldmaskpb.FieldMask)(f))
}

type Storage struct {
	Type string `json:"type,omitempty"`
	pgxutil.ConnectConfig
	Name string `json:"name,omitempty"`
	TTL  *TTL   `json:"ttl,omitempty"`
	// Retention is the minimum time records should be stored for. Zero-value (not-specified) means "forever".
	// Records can be deleted after this period, but may be kept longer depending on the cleanup cycle (e.g. if records
	// are only pruned once a day, a record could be kept for retention + 1day). Not all storage types might support this.
	Retention jsontypes.Duration `json:"retention,omitempty"`
}

type TTL struct {
	MaxAge   jsontypes.Duration `json:"maxAge,omitempty"`
	MaxCount int64              `json:"maxCount,omitempty"`
}
