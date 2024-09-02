package appconf

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/block"
)

// Store is the set of operations required for a config store, as used by BootConfig.
//
// DirStore is the canonical implementation of this interface.
type Store interface {
	// GetLocalConfig returns the local config from the store, or nil if there is no local config.
	GetLocalConfig() (*Config, error)
	// SetLocalConfig replaces the local config in the store with the provided config.
	// c must not be nil.
	SetLocalConfig(c *Config) error
	// GetActiveConfig returns the active config from the store, or nil if there is no active config.
	GetActiveConfig() (*Config, error)
	// SetActiveConfig replaces the active config in the store with the provided config.
	// c must not be nil.
	SetActiveConfig(c *Config) error
	// SavePatches saves a set of patches to the patch log, returning a unique name for the log entry.
	SavePatches(patches []block.Patch) (ref string, err error)
}

const (
	localConfigFilename  = "local.json"
	activeConfigFilename = "active.json"
)

// DirStore is a Store that stores config files in a directory on disk.
type DirStore struct {
	dir string
}

func NewDirStore(dir string) *DirStore {
	return &DirStore{dir: dir}
}

func (s *DirStore) GetLocalConfig() (*Config, error) {
	local, err := configFromFile(filepath.Join(s.dir, localConfigFilename))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return local, err
}

func (s *DirStore) SetLocalConfig(c *Config) error {
	if err := s.ensureDirExists(); err != nil {
		return err
	}
	return writeJSONAtomic(s.dir, localConfigFilename, c)
}

func (s *DirStore) GetActiveConfig() (*Config, error) {
	conf, err := configFromFile(filepath.Join(s.dir, activeConfigFilename))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return conf, err
}

func (s *DirStore) SetActiveConfig(c *Config) error {
	if err := s.ensureDirExists(); err != nil {
		return err
	}
	return writeJSONAtomic(s.dir, activeConfigFilename, c)
}

func (s *DirStore) SavePatches(patches []block.Patch) (ref string, err error) {
	err = s.ensureDirExists()
	if err != nil {
		return "", err
	}
	name := "patch-" + time.Now().UTC().Format("20060102-150405") + ".json"
	err = writeJSONAtomic(s.dir, name, patches)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *DirStore) ensureDirExists() error {
	return os.MkdirAll(s.dir, 0755)
}

// BootConfig resolves the active Config for a node.
// newLocal should be the fully loaded Config from outside the system.
// The provided newLocal config and resolved active config are stored in the provided Store.
// The store may be empty:
//   - If the store has no saved local config, this is treated the same as an empty saved local config.
//   - If the store has no saved active config, the newLocal config is used verbatim as the active config.
//
// schema describes the structure of the Config object, used to produce patches. Blocks can be used to generate an appropriate schema.
func BootConfig(newLocal *Config, store Store, schema []block.Block, logger *zap.Logger) (*Config, error) {
	oldLocal, err := store.GetLocalConfig()
	if err != nil {
		return nil, err
	}
	if oldLocal == nil {
		logger.Debug("no local config cache found, treating as empty")
		oldLocal = &Config{}
	}

	oldActive, err := store.GetActiveConfig()
	if err != nil {
		return nil, err
	}
	var newActive *Config
	if oldActive == nil {
		// no active config, just use the provided local config
		newActive = newLocal
	} else {
		patches, err := block.Diff(oldLocal, newLocal, schema)
		if err != nil {
			return nil, err
		}
		if len(patches) > 0 {
			patchRef, err := store.SavePatches(patches)
			if err != nil {
				return nil, err
			}
			logger.Info("applied config patch", zap.String("ref", patchRef))
		}

		newActive, err = block.ApplyPatches(oldActive, patches)
		if err != nil {
			return nil, err
		}
	}

	err = store.SetActiveConfig(newActive)
	if err != nil {
		return nil, err
	}
	err = store.SetLocalConfig(newLocal)
	if err != nil {
		return nil, err
	}
	return newActive, nil
}

func writeFileAtomic(dir, filename string, data []byte) (err error) {
	tmpFile, err := os.CreateTemp(dir, filename)
	if err != nil {
		return err
	}
	_, err = tmpFile.Write(data)
	if err != nil {
		_ = tmpFile.Close()
		return err
	}
	err = tmpFile.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmpFile.Name(), filepath.Join(dir, filename))
}

func writeJSONAtomic(dir, filename string, data any) error {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(dir, filename, encoded)
}

// Blocks returns a set of block.Block that represent the structure of a Config object.
// The parameters describe blocks for driver, automations and zones, keyed by type.
func Blocks(driverBlocks, autoBlocks, zoneBlocks map[string][]block.Block) []block.Block {
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
			BlocksByType: driverBlocks,
			Blocks:       defaultBlocks,
		},
		{
			Path:         []string{"automation"},
			Key:          "name",
			TypeKey:      "type",
			BlocksByType: autoBlocks,
			Blocks:       defaultBlocks,
		},
		{
			Path:         []string{"zones"},
			Key:          "name",
			TypeKey:      "type",
			BlocksByType: zoneBlocks,
			Blocks:       defaultBlocks,
		},
		{
			Path: []string{"metadata"},
		},
	}
}
