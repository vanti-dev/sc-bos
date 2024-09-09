package appconf

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/renameio/v2/maybe"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/vanti-dev/sc-bos/internal/confmerge"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/block"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Store struct {
	logger  *zap.Logger
	backing confmerge.Store

	m      sync.Mutex
	active Config
}

func LoadStore(external Config, schema Schema, storeDir string, logger *zap.Logger) (*Store, error) {
	store := confmerge.NewDirStore(storeDir)
	active, patches, err := confmerge.Merge(external, store, schema.Blocks())
	err = multierr.Append(err, saveConfigPatches(patches, filepath.Join(storeDir, patchDirName), logger))
	if err != nil {
		return nil, err
	}

	return &Store{
		logger:  logger,
		backing: store,
		active:  active,
	}, nil
}

func (s *Store) Active() Config {
	s.m.Lock()
	defer s.m.Unlock()
	return s.active.clone()
}

func (s *Store) Drivers() *DriverStore {
	return &DriverStore{store: s}
}

func (s *Store) Automations() *AutomationStore {
	return &AutomationStore{store: s}
}

func (s *Store) Zones() *ZoneStore {
	return &ZoneStore{store: s}
}

func (s *Store) save(updated Config) error {
	encoded, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		return err
	}
	err = s.backing.SetActiveConfig(encoded)
	if err != nil {
		return err
	}
	s.active = updated
	return nil
}

type DriverStore struct {
	store *Store
}

func (ds *DriverStore) SaveConfig(_ context.Context, name string, typ string, data []byte) error {
	ds.store.m.Lock()
	defer ds.store.m.Unlock()

	if name == "" {
		return errors.New("driver name is required")
	}

	// insert the updated driver into a copy of the config
	updated := ds.store.active.clone()
	idx := slices.IndexFunc(updated.Drivers, func(c driver.RawConfig) bool {
		return c.Name == name
	})
	var driverCfg driver.RawConfig
	if idx < 0 {
		// add new driver
		if typ == "" {
			return errors.New("driver type is required to add a new driver")
		}
		driverCfg = driver.RawConfig{
			BaseConfig: driver.BaseConfig{
				Name:     name,
				Type:     typ,
				Disabled: false,
			},
			Raw: data,
		}
		updated.Drivers = append(updated.Drivers, driverCfg)
	} else {
		existing := updated.Drivers[idx]
		if typ != "" && typ != existing.Type {
			return errors.New("driver type cannot be changed")
		}
		driverCfg = existing
		driverCfg.Raw = data
		updated.Drivers[idx] = driverCfg
	}

	// test that the new/updated driver config marshalls successfully - we don't want to store a bad config
	_, err := json.Marshal(driverCfg)
	if err != nil {
		return err
	}

	return ds.store.save(updated)
}

type AutomationStore struct {
	store *Store
}

func (as *AutomationStore) SaveConfig(_ context.Context, name string, typ string, data []byte) error {
	as.store.m.Lock()
	defer as.store.m.Unlock()

	if name == "" {
		return errors.New("automation name is required")
	}

	// insert the updated driver into a copy of the config
	updated := as.store.active.clone()
	idx := slices.IndexFunc(updated.Automation, func(c auto.RawConfig) bool {
		return c.Name == name
	})
	var autoCfg auto.RawConfig
	if idx < 0 {
		// add new driver
		if typ == "" {
			return errors.New("automation type is required to add a new automation")
		}
		autoCfg = auto.RawConfig{
			Config: auto.Config{
				Name:     name,
				Type:     typ,
				Disabled: false,
			},
			Raw: data,
		}
		updated.Automation = append(updated.Automation, autoCfg)
	} else {
		existing := updated.Automation[idx]
		if typ != "" && typ != existing.Type {
			return errors.New("automation type cannot be changed")
		}
		autoCfg = existing
		autoCfg.Raw = data
		updated.Automation[idx] = autoCfg
	}

	// test that the updated config marshalls successfully - we don't want to store a bad config
	_, err := json.Marshal(updated)
	if err != nil {
		return err
	}

	return as.store.save(updated)
}

type ZoneStore struct {
	store *Store
}

func (zs *ZoneStore) SaveConfig(_ context.Context, name string, typ string, data []byte) error {
	zs.store.m.Lock()
	defer zs.store.m.Unlock()

	if name == "" {
		return errors.New("zone name is required")
	}

	// insert the updated driver into a copy of the config
	updated := zs.store.active.clone()
	idx := slices.IndexFunc(updated.Zones, func(c zone.RawConfig) bool {
		return c.Name == name
	})
	var zoneCfg zone.RawConfig
	if idx < 0 {
		// add new driver
		if typ == "" {
			return errors.New("zone type is required to add a new zone")
		}
		zoneCfg = zone.RawConfig{
			Config: zone.Config{
				Name:     name,
				Type:     typ,
				Disabled: false,
			},
			Raw: data,
		}
		updated.Zones = append(updated.Zones, zoneCfg)
	} else {
		existing := updated.Zones[idx]
		if typ != "" && typ != existing.Type {
			return errors.New("zone type cannot be changed")
		}
		zoneCfg = existing
		zoneCfg.Raw = data
		updated.Zones[idx] = zoneCfg
	}

	// test that the updated config marshalls successfully - we don't want to store a bad config
	_, err := json.Marshal(updated)
	if err != nil {
		return err
	}

	return zs.store.save(updated)
}

const patchDirName = "patches"

func saveConfigPatches(patches []block.Patch, dir string, logger *zap.Logger) error {
	if len(patches) == 0 {
		return nil
	}

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	filename := filepath.Join(dir, "patch-"+time.Now().UTC().Format("20060102-150405")+".json")
	raw, err := json.MarshalIndent(patches, "", "  ")
	if err != nil {
		return err
	}
	err = maybe.WriteFile(filename, raw, 0644)
	if err != nil {
		return err
	}

	logger.Info("applied patches based on external config",
		zap.Int("count", len(patches)),
		zap.String("patchLogFile", filename),
	)
	return nil
}

type Schema struct {
	Drivers     map[string][]block.Block
	Automations map[string][]block.Block
	Zones       map[string][]block.Block
}

func (s Schema) Blocks() []block.Block {
	defaultBlocks := []block.Block{
		{
			Path: []string{"disabled"},
		},
	}

	return []block.Block{
		{
			Path:         []string{"drivers"},
			Key:          "name",
			TypeKey:      "type",
			BlocksByType: s.Drivers,
			Blocks:       defaultBlocks,
		},
		{
			Path:         []string{"automation"},
			Key:          "name",
			TypeKey:      "type",
			BlocksByType: s.Automations,
			Blocks:       defaultBlocks,
		},
		{
			Path:         []string{"zones"},
			Key:          "name",
			TypeKey:      "type",
			BlocksByType: s.Zones,
			Blocks:       defaultBlocks,
		},
		{
			Path: []string{"metadata"},
		},
	}
}
