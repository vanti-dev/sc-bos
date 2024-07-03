package appconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"go.uber.org/multierr"
)

const AlternateKey = "alternate_key.json"

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

func normaliseDeviceName(s string) string {
	return strings.ReplaceAll(s, "/", "-")
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
		var pageFile page
		err := json.Unmarshal(file, &pageFile)
		if err != nil {
			return err
		}
		switch kind {
		case reflect.Int:
			i, err := pageFile.Value.(json.Number).Int64()
			if err != nil {
				return err
			}
			value.SetInt(i)
		case reflect.String:
			value.SetString(pageFile.Value.(string))
		case reflect.Bool:
			value.SetBool((pageFile.Value).(bool))
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			f, err := pageFile.Value.(json.Number).Float64()
			if err != nil {
				return err
			}
			value.SetFloat(f)
		case reflect.Struct:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Map:
			bytes, err := json.Marshal(pageFile.Value)
			err = json.Unmarshal(bytes, v.Interface())
			if err != nil {
				return err
			}
		default:
		}
	}

	return nil
}

func setValue(s *any, path string) error {
	value := reflect.ValueOf(*s)
	kind := value.Kind()
	file, _ := readFile(path)
	var pageFile page
	err := json.Unmarshal(file, &pageFile)
	if err != nil {
		return err
	}
	switch kind {
	case reflect.Int:
		i, err := pageFile.Value.(json.Number).Int64()
		if err != nil {
			return err
		}
		value.SetInt(i)
	case reflect.String:
		*s = pageFile.Value.(string)
	case reflect.Bool:
		value.SetBool((pageFile.Value).(bool))
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		f, err := pageFile.Value.(json.Number).Float64()
		if err != nil {
			return err
		}
		value.SetFloat(f)
	case reflect.Struct:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		bytes, err := json.Marshal(pageFile.Value)
		err = json.Unmarshal(bytes, value.Interface())
		if err != nil {
			return err
		}
	default:
	}
	return nil
}

func mergeRawStruct(s any, path string) error {

	isDrcty, err := isDir(path)
	if err != nil {
		return err
	}
	if isDrcty {
		// if we are at a directory the corresponding part w
		directory, err := readDir(path)
		if err != nil {
			return err
		}
		for _, d := range directory {
			key := d.Name()

			value := reflect.ValueOf(s)
			kind := value.Kind()

			if kind == reflect.Map {
				if _, ok := s.(map[string]interface{})[key]; ok {
					nextPath := filepath.Join(path, key)

					// when we get to a file, then the corresponding value in the map could be a primitive.
					// we can't pass a pointer to the value of a map so primitives will be passed by value and we can't
					// do this to modify the value. so check if nextPath is a file and then modify the element in the map
					isDrcty, err := isDir(nextPath)
					if err != nil {
						return err
					}

					if isDrcty {
						err := mergeRawStruct(s.(map[string]interface{})[key], nextPath)
						if err != nil {
							return err
						}
					} else {
						file, _ := readFile(nextPath)
						var pageFile page
						err := json.Unmarshal(file, &pageFile)
						if err != nil {
							return err
						}
						// we are assuming that Value in the page file holds the correct type
						s.(map[string]interface{})[key] = pageFile.Value
					}
				}
			} else if kind == reflect.Slice {
				// when we get to a slice, we need to know which element of the slice we are looking at
				// default to 'name' but use getAlternateKey as we can define an alternate key
				for i := 0; i < len(s.([]interface{})); i++ {
					elem := s.([]interface{})[i]
					elementKey := "name"
					alternateKey, err := getAlternateKey(path)

					if err == nil {
						elementKey = alternateKey
					}

					if v, ok := elem.(map[string]interface{})[elementKey]; ok {
						normalisedName := normaliseDeviceName(v.(string))
						if normalisedName == key {
							nextPath := filepath.Join(path, key)
							err := mergeRawStruct(elem, nextPath)
							if err != nil {
								return err
							}
						}
					} else {
						return errors.New(fmt.Sprintf("no %s field in slice element", elementKey))
					}
				}
			}
		}
	} else {
		// we are at the file level so we just try to set the value in the map
		file, _ := readFile(path)
		var pageFile page
		err := json.Unmarshal(file, &pageFile)
		if err != nil {
			return err
		}

		err = setValue(&s, path)
		if err != nil {
			return err
		}

	}
	return nil
}

func getAlternateKey(path string) (string, error) {

	path = filepath.Join(path, AlternateKey)
	f, err := readFile(path)
	if err != nil {
		return "", err
	}

	var keyFile alternate
	err = json.Unmarshal(f, &keyFile)
	if err != nil {
		return "", err
	}
	return keyFile.Key, nil
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

	var errs error
	for _, d := range directory {
		switch strings.ToLower(d.Name()) {
		case "metadata":
			path := filepath.Join(dbRoot, d.Name())
			err := mergeField(reflect.ValueOf(appConfig.Metadata), path)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
		case "drivers":
			path := filepath.Join(dbRoot, d.Name())
			driversDirectory, err := readDir(path)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
			for i := 0; i < len(appConfig.Drivers); i++ {
				driver := &appConfig.Drivers[i]
				for _, d := range driversDirectory {
					if strings.EqualFold(driver.Type, d.Name()) {
						s := make(map[string]interface{})

						err = json.Unmarshal(driver.Raw, &s)
						if err != nil {
							errs = multierr.Append(errs, err)
						}

						path := filepath.Join(dbRoot, "drivers", d.Name())
						err := mergeRawStruct(s, path)
						if err != nil {
							errs = multierr.Append(errs, err)
						}

						driver.Raw, err = json.Marshal(s)
						if err != nil {
							errs = multierr.Append(errs, err)
						}
					}
				}
			}
		}
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
	Value any `json:"value,omitempty"`
}

type alternate struct {
	Key string `json:"key,omitempty"`
}

// writePageFile this writes a page file to at the given path
// a page file defines the value that is going to replace whatever config is located at the path
// and also includes an optional key which specifies the item in a collection that we are editing
// if no key is included then the entire config item located at path is replaced
func writePageFile(path string, value any) error {

	var pageFile page
	pageFile = page{
		Value: value,
	}

	pageFileJson, err := json.Marshal(pageFile)
	if err != nil {
		return err
	}

	return writeFile(path, pageFileJson, 0664)
}

func writeAlternateKey(path string, key string) error {
	keyFile := alternate{
		Key: key,
	}

	keyFileJson, err := json.Marshal(keyFile)
	if err != nil {
		return err
	}

	path = filepath.Join(path, AlternateKey)
	return writeFile(path, keyFileJson, 0664)
}
