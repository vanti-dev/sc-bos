package sysconf

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
)

// LoadAllFromJSON populates dst with json formatted config from all dst.ConfigDirs x dst.ConfigFiles.
// If any file path doesn't exist then it will be skipped.
func LoadAllFromJSON(dst *Config) error {
	var allErrs error
	for _, dir := range dst.ConfigDirs {
		for _, file := range dst.ConfigFiles {
			err := LoadFromJSONFile(dst, filepath.Join(dir, file))
			if !errors.Is(err, os.ErrNotExist) {
				allErrs = multierr.Append(allErrs, err)
			}
		}
	}
	return allErrs
}

// LoadFromConfigFilesJSON populates dst with json formatted config from all dst.ConfigFiles (based on current working directory).
func LoadFromConfigFilesJSON(dst *Config) error {
	var allErrs error
	for _, file := range dst.ConfigFiles {
		err := LoadFromJSONFile(dst, filepath.Join(".", file))
		if !errors.Is(err, os.ErrNotExist) {
			allErrs = multierr.Append(allErrs, err)
		}
	}
	return allErrs
}

// LoadFromJSONFile reads filePath and calls LoadFromJSON.
func LoadFromJSONFile(dst *Config, filePath string) error {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return LoadFromJSON(dst, jsonData)
}

// LoadFromJSON populates dst with the given json formatted bytes.
func LoadFromJSON(dst *Config, jsonData []byte) error {
	return json.Unmarshal(jsonData, dst)
}
