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
	return info.IsDir(), nil
}

// normaliseDeviceName replaces all instances of /, :, and spaces with -
func normaliseDeviceName(s string) string {
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ":", "-")
	return s
}

// mergeWithReflection does the same thing as mergeRawStruct but uses reflection to merge the changes.
// used only for Metadata at the moment, probably unnecessary as we can just use mergeRawStruct
func mergeWithReflection(v reflect.Value, path string) error {
	if v.Kind() != reflect.Ptr || v.IsNil() {
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
						err = mergeWithReflection(value.Field(i), path)
					} else {
						err = mergeWithReflection(value.Field(i).Addr(), path)
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
		file, err := readFile(path)
		if err != nil {
			return err
		}
		var pageFile page
		err = json.Unmarshal(file, &pageFile)
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
	case reflect.Float32, reflect.Float64:
		f, err := pageFile.Value.(json.Number).Float64()
		if err != nil {
			return err
		}
		value.SetFloat(f)
	case reflect.Struct, reflect.Slice, reflect.Map:
		bytes, err := json.Marshal(pageFile.Value)
		err = json.Unmarshal(bytes, value.Interface())
		if err != nil {
			return err
		}
	default:
	}
	return nil
}

// mergeRawStruct recursively merges the file structure located at the root path with the raw struct s
// looks at path, if it is a directory we try to match the name of the directory with the name of a field in the struct
// if the dir name matches a field, then recursively call mergeRawStruct(field, dirPath) again until the dirPath is a file not a dir
// when path is a file then we are at the lowest level and we apply the change located in file to the corresponding field in struct
// when we encounter a map or a slice we look at the dir name and try to match it with the key of the map or the `name` field of the slice element
// if the dir at slice level contains an alternate_key.json file then we use the key specified in that file instead of the default `name`
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

// mergeDbWithExtConfig merges changes found in the root with the app config
// this is done by reading the root and then recursively merging the changes
// ie. if the root contains a file at root/drivers/floor-01-bms/localInterface
// the value in this file will overwrite the config at appConfig.Drivers["floor-01/bms"].LocalInterface
func mergeDbWithExtConfig(appConfig *Config, root string) error {

	directory, err := readDir(root)
	if err != nil {
		return err
	}

	var errs error
	for _, d := range directory {
		dirName := strings.ToLower(d.Name())
		switch dirName {
		case "metadata":
			path := filepath.Join(root, d.Name())
			// todo, probably unnecessary now as this uses a different method to perform the merge as the exact structure is known. can just use mergeRawStruct
			err := mergeWithReflection(reflect.ValueOf(appConfig.Metadata), path)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
		case "drivers":
			path := filepath.Join(root, d.Name())
			autoDirectory, err := readDir(path)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
			for i := 0; i < len(appConfig.Drivers); i++ {
				auto := &appConfig.Drivers[i]
				for _, nextDir := range autoDirectory {
					if strings.EqualFold(normaliseDeviceName(auto.Name), nextDir.Name()) {
						s := make(map[string]interface{})

						err = json.Unmarshal(auto.Raw, &s)
						if err != nil {
							errs = multierr.Append(errs, err)
						}

						path := filepath.Join(root, strings.ToLower(d.Name()), nextDir.Name())
						err := mergeRawStruct(s, path)
						if err != nil {
							errs = multierr.Append(errs, err)
						}

						auto.Raw, err = json.Marshal(s)
						if err != nil {
							errs = multierr.Append(errs, err)
						}
					}
				}
			}
		case "automation":
			// todo this is a c&p of case "drivers", if drivers, autos et. al implemented GetName/GetRaw it could be done using generics
			path := filepath.Join(root, d.Name())
			autoDirectory, err := readDir(path)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
			for i := 0; i < len(appConfig.Automation); i++ {
				auto := &appConfig.Automation[i]
				for _, nextDir := range autoDirectory {
					if strings.EqualFold(normaliseDeviceName(auto.Name), nextDir.Name()) {
						s := make(map[string]interface{})

						err = json.Unmarshal(auto.Raw, &s)
						if err != nil {
							errs = multierr.Append(errs, err)
						}

						path := filepath.Join(root, strings.ToLower(d.Name()), nextDir.Name())
						err := mergeRawStruct(s, path)
						if err != nil {
							errs = multierr.Append(errs, err)
						}

						auto.Raw, err = json.Marshal(s)
						if err != nil {
							errs = multierr.Append(errs, err)
						}
					}
				}
			}
		case "zones":
			// todo this is a c&p of case "drivers", if drivers, autos et. al implemented GetName/GetRaw it could be done using generics
			path := filepath.Join(root, d.Name())
			zonesDirectory, err := readDir(path)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
			for i := 0; i < len(appConfig.Zones); i++ {
				auto := &appConfig.Zones[i]
				for _, nextDir := range zonesDirectory {
					if strings.EqualFold(normaliseDeviceName(auto.Name), nextDir.Name()) {
						s := make(map[string]interface{})

						err = json.Unmarshal(auto.Raw, &s)
						if err != nil {
							errs = multierr.Append(errs, err)
						}

						path := filepath.Join(root, strings.ToLower(d.Name()), nextDir.Name())
						err := mergeRawStruct(s, path)
						if err != nil {
							errs = multierr.Append(errs, err)
						}

						auto.Raw, err = json.Marshal(s)
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
