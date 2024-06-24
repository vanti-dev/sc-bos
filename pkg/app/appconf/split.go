package appconf

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// Split The top level struct of the *.split.json file which defines how to split the *.json file
// recursive
type Split struct {
	SplitKey string  `json:"splitKey"`
	Key      string  `json:"key"`
	Path     string  `json:"path"`
	Splits   []Split `json:"splits"`
}

func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
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

// crap name recurses through the
func mergeField(
	v reflect.Value, path string) error {

	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer value")
	}
	value := reflect.Indirect(v)

	// we now have a struct and a path to an element of that struct
	// if the path points to a file, we are ready to update the struct
	// else if the path is a dir, we keep recursing until we find a file
	isDrcty, err := isDir(path)
	if err != nil {
		return err
	}
	if isDrcty {
		directory, err := readDir(path)
		if err != nil {
			return err
		}
		typeOfT := value.Type()
		for _, d := range directory {
			for i := 0; i < value.NumField(); i++ {
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
	} else {
		kind := value.Kind()
		file, _ := readFile(path)
		switch kind {
		case reflect.Int:
			i, err := strconv.ParseInt(string(file), 10, 0)
			if err != nil {
				return err
			}
			value.SetInt(i)
		case reflect.String:
			value.SetString(string(file))
		case reflect.Bool:
			b, err := strconv.ParseBool(string(file))
			if err != nil {
				return err
			}
			value.SetBool(b)
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			f, err := strconv.ParseFloat(string(file), 64)
			if err != nil {
				return err
			}
			value.SetFloat(f)
		case reflect.Struct:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Map:
			err = json.Unmarshal(file, v.Interface())
			if err != nil {
				return err
			}
		default:
		}
	}

	return nil
}

// mergeDbWithExtConfig reads the ext config and merges changes from the DB into it
func mergeDbWithExtConfig(appConfig *Config, dbRoot string) error {

	// so at this point we have the ext config loaded into appConfig & splits loaded into splits
	// now we need to look through the db and merge the changes into appConfig

	// open the root of the db, it should look like appConfig.Config top level, so:
	// /metadata /drivers /automations /zones
	directory, err := readDir(dbRoot)
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
