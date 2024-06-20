package appconf

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
)

// unmarshalExt reads files and includes from external config files.
func unmarshalExt(dst *Config, dir string, paths ...string) ([]string, error) {
	return loadIncludes(dir, dst, paths, nil)
}

// readSplits reads split data from a file.
func readSplits(file string) ([]split, error) {
	var splits []split
	data, err := readFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &splits)
	if err != nil {
		return nil, err
	}
	return splits, nil
}

// writeSplits writes split data to a file.
func writeSplits(file string, splits []split) error {
	data, err := json.Marshal(splits)
	if err != nil {
		return err
	}
	return writeFile(file, data, 0664)
}

func splitsEqual(a, b []split) bool {
	return true
}

func paginate(c *Config, splits []split) []page {
	return nil
}

// writePages writes pages to files based in the directory file is in.
// Unchanged pages when compared with the filesystem are skipped.
// All written pages are returned.
func writePages(file string, pages []page) ([]page, error) {
	// todo: add a .sum file for pages
	// The above will improve performance by reducing reads but also
	// allow us to remove files that are no longer part of the config.

	dir := filepath.Dir(file)
	if err := mkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	var (
		changed []page
		errs    error
	)
	for _, p := range pages {
		fullPath := filepath.Join(dir, p.Path)
		var oldHash, newHash []byte

		// hash the old page data
		f, err := openFile(fullPath)
		switch {
		case errors.Is(err, os.ErrNotExist):
			changed = append(changed, p)
		case err != nil:
			// don't report this error on the assumption that write will also fail,
			// and the write error is more informative
		default: // err == nil
			h := sha256.New()
			_, err = copyFile(h, f)
			f.Close()

			if err == nil {
				oldHash = h.Sum(nil)
			}
		}

		// hash the new page data
		h := sha256.New()
		_, err = h.Write(p.JSON)
		if err != nil {
			// unexpected error
			errs = multierr.Append(errs, fmt.Errorf("hashing page %q %w", p.Path, err))
			continue
		}
		newHash = h.Sum(nil)

		if bytes.Equal(oldHash, newHash) {
			continue // nothing has changed, skip writing
		}

		// write our changes
		if err := mkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		err = writeFile(filepath.Join(dir, p.Path), p.JSON, 0664)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		changed = append(changed, p)
	}
	return changed, nil
}

func unmarshalPages(dst *Config, file string) error {
	// todo: I can't work out if it's safe to blindly follow $refs in these documents or not.
	//  An alternative to this would be to accept a []split and only read $refs that come from there.
	//  This might be required if the config schema also uses $ref for something else.
	// todo: it's probably not a good idea to load all the json into memory at once,
	//  but we don't have streaming json until json/v2

	jsonMap := make(map[string]any)
	data, err := readFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return err
	}
	return nil
}

func readJSONMapRefs(file string) (any, error) {
	// read the file into a map
	var jsonData any
	data, err := readFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	err = replaceRefs(filepath.Dir(file), jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func replaceRefs(dir string, jsonData any) error {
	switch v := jsonData.(type) {
	case map[string]any:
		if ref, ok := v["$ref"]; ok {
			if len(v) != 1 {
				return nil // our logic only creates single key ref objects, this isn't our $ref
			}
			refStr, ok := ref.(string)
			if !ok {
				return nil // not a string ref, not our $ref
			}
			path := filepath.Join(dir, refStr)
			newDir := filepath.Dir(path)

		}
	case []any:
		for i, item := range v {
			err := replaceRefs(dir, item)
			if err != nil {
				return err
			}
			v[i] = item
		}
	default:
		// nothing to change, other types of value remain the same
		return nil
	}
}

type split struct {
	Path        string             `json:"path,omitempty"`
	Key         string             `json:"key,omitempty"`
	SplitKey    string             `json:"splitKey,omitempty"`
	SplitsByKey map[string][]split `json:"splitsByKey,omitempty"`
}

type page struct {
	Path string
	JSON json.RawMessage
}

type bootConfig struct {
	extDir           string
	extCacheRootFile string
	splitCacheFile   string
	dbRootFile       string
	liveSplits       func() ([]split, error)
}

func (b bootConfig) unmarshalBootConfig(dst *Config, paths ...string) error {
	splits, err := readSplits(b.splitCacheFile)
	if errors.Is(err, os.ErrNotExist) {
		// if the splits don't exist then ext cache and db shouldn't either, aka first boot
		err := b.unmarshalFirstBootConfig(dst, paths...)
		if err != nil {
			return fmt.Errorf("first boot %w", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("read splits: %w", err)
	}

	extCfg := &Config{}
	_, err = unmarshalExt(extCfg, b.extDir, paths...)
	if err != nil {
		return fmt.Errorf("ext unmarshal: %w", err)
	}

	extPages := paginate(extCfg, splits)
	extChanges, err := writePages(b.extCacheRootFile, extPages)
	if err != nil {
		return fmt.Errorf("write ext cache: %w", err)
	}
	if len(extChanges) > 0 {
		_, err = writePages(b.dbRootFile, extChanges)
		if err != nil {
			return fmt.Errorf("write ext changes to db: %w", err)
		}
	}

	err = unmarshalPages(dst, b.dbRootFile)
	if err != nil {
		return fmt.Errorf("unmarshal db: %w", err)
	}

	// If the cached splits have changed when compared to live splits,
	// then we need to update the files we have on disk to match the new format.
	// Splits can change if drivers are updated or added (or removed) from the system.
	liveSplits, err := b.liveSplits()
	if err != nil {
		return fmt.Errorf("get live splits: %w", err)
	}
	if !splitsEqual(splits, liveSplits) {
		extPages = paginate(extCfg, liveSplits)
		_, err = writePages(b.extCacheRootFile, extPages)
		if err != nil {
			return fmt.Errorf("update ext cache: %w", err)
		}

		dbPages := paginate(dst, liveSplits)
		_, err = writePages(b.dbRootFile, dbPages)
		if err != nil {
			return fmt.Errorf("update db: %w", err)
		}

		err = writeSplits(b.splitCacheFile, liveSplits)
		if err != nil {
			return fmt.Errorf("write live splits: %w", err)
		}
	}

	return nil
}

func (b bootConfig) unmarshalFirstBootConfig(dst *Config, paths ...string) error {
	_, err := unmarshalExt(dst, b.extDir, paths...)
	if err != nil {
		return fmt.Errorf("ext unmarshal: %w", err)
	}

	splits, err := b.liveSplits()
	if err != nil {
		return fmt.Errorf("get live splits: %w", err)
	}

	pages := paginate(dst, splits)

	_, err = writePages(b.dbRootFile, pages)
	if err != nil {
		return fmt.Errorf("write db: %w", err)
	}
	_, err = writePages(b.extCacheRootFile, pages)
	if err != nil {
		return fmt.Errorf("write ext cache: %w", err)
	}
	err = writeSplits(b.splitCacheFile, splits)
	if err != nil {
		return fmt.Errorf("write splits: %w", err)
	}
	return nil
}
