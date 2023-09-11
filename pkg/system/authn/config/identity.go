package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Identities represents configured accounts and secrets for users and tenants.
// Identities can be unmarshalled from json in different ways that affect how Load behaves:
//  1. true indicates that Load should read identities from the defaults provided
//  2. false indicates that Load will always return no identities or error
//  3. A string represents a path, rooted in baseDir, to a file containing identities to load
//  4. A JSON array that contains a list of identities explicitly defined and returned by Load
type Identities struct {
	path    string
	content []Identity
}

func (p *Identities) Load(baseDir, defaultFilename string) ([]Identity, error) {
	if p == nil {
		return nil, nil // disabled
	}

	if len(p.content) > 0 {
		return p.content, nil
	}

	var usingDefault bool
	path := p.path
	if path == "" {
		path = defaultFilename
		usingDefault = true
	}
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(filepath.Join(baseDir, path))
	if err != nil {
		if usingDefault && errors.Is(err, os.ErrNotExist) {
			// special case for when
			return nil, nil
		}
		return nil, err
	}
	var c []Identity
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (p *Identities) UnmarshalJSON(data []byte) error {
	switch {
	case bytes.Equal(data, []byte("true")):
		*p = Identities{} // uses defaults during Load
		return nil
	case bytes.Equal(data, []byte("false")):
		return nil // Load always returns nothing
	case data[0] == '"':
		var path string
		err := json.Unmarshal(data, &path)
		if err != nil {
			return err
		}
		*p = Identities{path: path}
		return nil
	default:
		var c []Identity
		err := json.Unmarshal(data, &c)
		if err != nil {
			return err
		}
		*p = Identities{content: c}
		return nil
	}
}

type Identity struct {
	Title   string   `json:"title,omitempty"`
	ID      string   `json:"id,omitempty"`
	Secrets []Secret `json:"secrets,omitempty"`
	Zones   []string `json:"zones,omitempty"`
	Roles   []string `json:"roles,omitempty"`
}

type Secret struct {
	Note string `json:"note,omitempty"`
	Hash string `json:"hash,omitempty"`
}
