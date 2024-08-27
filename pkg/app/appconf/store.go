package appconf

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/vanti-dev/sc-bos/pkg/block"
)

type Store struct {
	dir string
}

// SwapLocalConfig loads the contents of the local config cache, returning it, and then replaces it with the provided config.
// If the local config cache is empty, returns a nil old config and no error.
// If err is not nil, the state of the local config cache is undefined.
func (s *Store) SwapLocalConfig(new *Config) (old *Config, err error) {
	localPath := filepath.Join(s.dir, "local.json")
	// if the file simply doesn't exist, that's fine, we will continue with old == nil
	old, err = configFromFile(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	err = writeJSONAtomic(s.dir, "local.json", new)
	return old, err
}

func (s *Store) GetActiveConfig() (*Config, error) {
	return configFromFile(filepath.Join(s.dir, "active.json"))
}

func (s *Store) SetActiveConfig(c *Config) error {
	return writeJSONAtomic(s.dir, "active.json", c)
}

func BootConfig(localConfDir, localConfFile string, store *Store, schema []block.Block) (*Config, error) {
	newLocal, err := LoadLocalConfig(localConfDir, localConfFile)
	if err != nil {
		return nil, err
	}

	oldLocal, err := store.SwapLocalConfig(newLocal)
	if err != nil {
		return nil, err
	}

	patches, err := block.Diff(oldLocal, newLocal, schema)
	if err != nil {
		return nil, err
	}

	oldActive, err := store.GetActiveConfig()
	if err != nil {
		return nil, err
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

func writeJSONAtomic(dir, filename string, data any) error {
	tmpFile, err := os.CreateTemp(dir, filename)
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	err = json.NewEncoder(tmpFile).Encode(data)
	if err != nil {
		return err
	}
	err = tmpFile.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmpFile.Name(), filepath.Join(dir, filename))
}
