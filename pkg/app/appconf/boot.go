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

type Store interface {
	SwapLocalConfig(new *Config) (old *Config, err error)
	GetActiveConfig() (*Config, error)
	SetActiveConfig(c *Config) error
	SavePatches(patches []block.Patch) (ref string, err error)
}

type DirStore struct {
	dir string
}

func NewDirStore(dir string) *DirStore {
	return &DirStore{dir: dir}
}

// SwapLocalConfig loads the contents of the local config cache, returning it, and then replaces it with the provided config.
// If the local config cache is empty, returns a nil old config and no error.
// If err is not nil, the state of the local config cache is undefined.
func (s *DirStore) SwapLocalConfig(new *Config) (old *Config, err error) {
	localPath := filepath.Join(s.dir, "local.json")
	// if the file simply doesn't exist, that's fine, we will continue with old == nil
	old, err = configFromFile(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	err = s.ensureDirExists()
	if err != nil {
		return nil, err
	}
	err = writeJSONAtomic(s.dir, "local.json", new)
	return old, err
}

func (s *DirStore) GetActiveConfig() (*Config, error) {
	return configFromFile(filepath.Join(s.dir, "active.json"))
}

func (s *DirStore) SetActiveConfig(c *Config) error {
	if err := s.ensureDirExists(); err != nil {
		return err
	}
	return writeJSONAtomic(s.dir, "active.json", c)
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

func BootConfig(newLocal *Config, store Store, schema []block.Block, logger *zap.Logger) (*Config, error) {
	oldLocal, err := store.SwapLocalConfig(newLocal)
	if err != nil {
		return nil, err
	}
	if oldLocal == nil {
		logger.Debug("no local config cache found, treating as empty")
		oldLocal = &Config{}
	}

	oldActive, err := store.GetActiveConfig()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if oldActive == nil {
		// no active config, just use the provided local config
		err = store.SetActiveConfig(newLocal)
		if err != nil {
			return nil, err
		}
		return newLocal, nil
	}

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

	newActive, err := block.ApplyPatches(oldActive, patches)
	if err != nil {
		return nil, err
	}
	err = store.SetActiveConfig(newActive)
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
