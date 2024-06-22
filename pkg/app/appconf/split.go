package appconf

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"go.uber.org/multierr"
)

// Split The top level struct of the *.split.json file which defines how to split the *.json file
// recursive
type Split struct {
	SplitKey string  `json:"splitKey"`
	Key      string  `json:"key"`
	Path     string  `json:"path"`
	Splits   []Split `json:"splits"`
}

// readSplits reads split data from a file.
func readSplits(file string) ([]Split, error) {
	f, err := readFile(file)

	if err != nil {
		return nil, err
	}

	var splitFile []Split
	err = json.Unmarshal(f, &splitFile)
	if err != nil {
		return nil, err
	}
	return splitFile, nil
}

// writeSplit writes the split structure recursively
func writeSplit(path string, split Split) error {
	if len(split.Splits) == 0 {
		// recursion over, write the child nodes & return
		path = filepath.Join(path, split.Path)
		// leave them empty on init, then if empty
		err := writeFile(path, []byte{}, 0664)
		if err != nil {
			return err
		}
	} else {
		path = filepath.Join(path, split.Path)
		if err := mkdirAll(path, 0755); err != nil {
			return err
		}
		for _, s := range split.Splits {
			err := writeSplit(path, s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// writeSplitStructure write the ext cache & config db files describing the config structure i.e.
// */metadata/location/floor
// */metadata/product/manufacturer
// */metadata/product/model
// this only needs to be called on first boot to set up the structure
func writeSplitStructure(root string, splits []Split) error {

	// first write the root directory
	if err := mkdirAll(root, 0755); err != nil {
		return err
	}

	// then let write the structure recursively
	for _, s := range splits {
		err := writeSplit(root, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeSplits writes split data to a file.
func writeSplits(file string, splits []split) error {
	data, err := json.Marshal(splits)
	if err != nil {
		return err
	}
	return writeFile(file, data, 0664)
}

func isPrimitiveType(t string) bool {
	return t == "int" || t == "string" || t == "float32" || t == "float64" || t == "bool"
}

func setValue(value *reflect.Value, toSet any) {
	typeOfT := value.Type()

	switch typeOfT.Name() {
	case "string":
		value.SetString(toSet.(string))
	case "int":
		value.SetInt(int64(toSet.(int)))
	case "float32":
		value.SetFloat(float64(toSet.(float32)))
	case "float64":
		value.SetFloat(toSet.(float64))
	case "bool":
		value.SetBool(toSet.(bool))
	}
}

// crap name recurses through the
func mergeField(
	value reflect.Value, path string) error {

	if value.Kind() != reflect.Ptr {
		return errors.New("Not a pointer value")
	}

	// we now have a struct, which should resemble the given directory
	// reflect on the struct, compare with the dir tree and when we reach primitive level
	// update the value of the struct with the value of the file

	value = reflect.Indirect(value)
	typeOfT := value.Type()

	kind := value.Kind()
	switch kind {
	case reflect.Int:
		file, _ := readFile(path)
		i, err := strconv.ParseInt(string(file), 10, 0)
		if err != nil {
			return err
		}
		value.SetInt(i)
	case reflect.String:
		file, _ := readFile(path)
		value.SetString(string(file))
	case reflect.Bool:
		file, _ := readFile(path)
		b, err := strconv.ParseBool(string(file))
		if err != nil {
			return err
		}
		value.SetBool(b)
	case reflect.Float32:
	case reflect.Float64:
		file, _ := readFile(path)
		f, err := strconv.ParseFloat(string(file), 64)
		if err != nil {
			return err
		}
		value.SetFloat(f)
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			directory, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for _, d := range directory {
				if strings.EqualFold(d.Name(), typeOfT.Field(i).Name) {
					path := filepath.Join(path, d.Name())

					if value.Field(i).Kind() == reflect.Ptr {
						err = mergeField(value.Field(i), path)
					} else {
						err = mergeField(value.Field(i).Addr(), path)
					}
					if err != nil {
						return err
					}
				} else {
					// do nothing, if the directory name doesn't match the struct field then move on to the next
				}
			}
		}
	default:
	}

	return nil
}

// mergeDbWithExtConfig reads the ext config and merges changes from the DB into it
func mergeDbWithExtConfig(appConfig *Config, dbRoot string) error {

	// so at this point we have the ext config loaded into appConfig & splits loaded into splits
	// now we need to look through the db and merge the changes into appConfig

	// open the root of the db, it should look like appConfig.Config top level, so:
	// /metadata /drivers /automations /zones
	directory, err := os.ReadDir(dbRoot)
	if err != nil {
		return err
	}

	for _, d := range directory {
		switch d.Name() {
		case "metadata":
			path := filepath.Join(dbRoot, d.Name())
			err := mergeField(reflect.ValueOf(appConfig.Metadata), path)
			if err != nil {
				return err
			}
		}
	}

	return nil
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
