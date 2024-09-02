package confmerge

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"

	"github.com/google/renameio/v2/maybe"

	"github.com/vanti-dev/sc-bos/pkg/block"
)

// Store is the set of operations required for a config store, as used by Merge.
//
// DirStore is the canonical implementation of this interface.
type Store interface {
	// GetLocalConfig returns the local config from the store, or nil if there is no local config.
	GetLocalConfig() ([]byte, error)
	// SetLocalConfig replaces the local config in the store with the provided config.
	// c must not be nil.
	SetLocalConfig(c []byte) error
	// GetActiveConfig returns the active config from the store, or nil if there is no active config.
	GetActiveConfig() ([]byte, error)
	// SetActiveConfig replaces the active config in the store with the provided config.
	// c must not be nil.
	SetActiveConfig(c []byte) error
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

func (s *DirStore) GetLocalConfig() ([]byte, error) {
	return s.read(localConfigFilename)
}

func (s *DirStore) SetLocalConfig(c []byte) error {
	return s.write(localConfigFilename, c)
}

func (s *DirStore) GetActiveConfig() ([]byte, error) {
	return s.read(activeConfigFilename)
}

func (s *DirStore) SetActiveConfig(c []byte) error {
	return s.write(activeConfigFilename, c)
}

func (s *DirStore) SavePatches(patches []block.Patch) (ref string, err error) {
	err = s.ensureDirExists()
	if err != nil {
		return "", err
	}
	name := "patch-" + time.Now().UTC().Format("20060102-150405") + ".json"
	raw, err := json.Marshal(patches)
	if err != nil {
		return "", err
	}
	err = s.write(name, raw)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *DirStore) ensureDirExists() error {
	return os.MkdirAll(s.dir, 0755)
}

func (s *DirStore) read(name string) ([]byte, error) {
	raw, err := os.ReadFile(filepath.Join(s.dir, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return raw, nil
}

func (s *DirStore) write(name string, data []byte) error {
	if err := s.ensureDirExists(); err != nil {
		return err
	}
	return maybe.WriteFile(filepath.Join(s.dir, name), data, 0644)
}

// Merge detects changes to local configuration, and applies those changes to the active configuration.
//
// Returns the new active configuration.
// newLocal should be the fully loaded Config from outside the system.
// The provided newLocal config and resolved active config are stored in the provided Store.
// The store may be empty:
//   - If the store has no saved local config, then newLocal is compared against the zero value of T.
//   - If the store has no saved active config, the newLocal config is used verbatim as the active config.
//
// schema describes the structure of T object, used to produce patches. This controls the granularity of patches
// that will be applied to the active config. For full detail of this, see the block package.
//
// Configs are stored as JSON, so T must marshal and unmarshal to/from JSON correctly.
func Merge[T any](newLocal *T, store Store, schema []block.Block, logger *zap.Logger) (*T, error) {
	oldLocal, err := getLocalJSON[T](store)
	if err != nil {
		return nil, err
	}
	if oldLocal == nil {
		logger.Debug("no local config cache found, treating as empty")
		var empty T
		oldLocal = &empty
	}

	oldActive, err := getActiveJSON[T](store)
	if err != nil {
		return nil, err
	}
	var newActive *T
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

	err = setActiveJSON(store, newActive)
	if err != nil {
		return nil, err
	}
	err = setLocalJSON(store, newLocal)
	if err != nil {
		return nil, err
	}
	return newActive, nil
}

func getLocalJSON[T any](store Store) (*T, error) {
	raw, err := store.GetLocalConfig()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	var c T
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func getActiveJSON[T any](store Store) (*T, error) {
	raw, err := store.GetActiveConfig()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	var c T
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func setLocalJSON[T any](store Store, c *T) error {
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return store.SetLocalConfig(raw)
}

func setActiveJSON[T any](store Store, c *T) error {
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return store.SetActiveConfig(raw)
}
