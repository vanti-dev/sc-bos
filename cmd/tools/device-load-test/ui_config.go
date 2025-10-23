package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/smart-core-os/sc-golang/pkg/trait"
)

type uiConfig struct {
	Config struct {
		Ops opsConfig `json:"ops"`
	} `json:"config"`
}

type opsConfig struct {
	Pages []pageConfig `json:"pages,omitempty"`
}

type pageConfig struct {
	Title    string            `json:"title,omitempty"`
	Path     string            `json:"path,omitempty"`
	Main     []componentConfig `json:"main,omitempty"`
	After    []componentConfig `json:"after,omitempty"`
	Children []pageConfig      `json:"children,omitempty"`
}

type componentConfig struct {
	Component string       `json:"component,omitempty"`
	Props     graphicProps `json:"props,omitempty"`
}

type graphicProps struct {
	Layers []graphicLayer `json:"layers,omitempty"`
}

type graphicLayer struct {
	Title      string `json:"title,omitempty"`
	ConfigPath string `json:"configPath,omitempty"`
}

type layerConfig struct {
	Templates map[string]layerTemplate `json:"templates,omitempty"`
	Elements  []layerElement           `json:"elements,omitempty"`
}

type layerTemplate struct {
	Sources map[string]layerSource `json:"sources,omitempty"`
}

type layerSource struct {
	Trait   trait.Name    `json:"trait"`
	Request sourceRequest `json:"request"`
}

type sourceRequest struct {
	Name string `json:"name"`
}

type layerElement struct {
	Template templateRef `json:"template"`
}

type templateRef struct {
	Ref  string            `json:"ref"`
	Vars map[string]string `json:"-"`
}

func (tr *templateRef) UnmarshalJSON(data []byte) error {
	props := make(map[string]string)
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}
	tr.Ref = props["ref"]
	delete(props, "ref")
	tr.Vars = props
	return nil
}

func (tr templateRef) MarshalJSON() ([]byte, error) {
	vars := maps.Clone(tr.Vars)
	vars["ref"] = tr.Ref
	return json.Marshal(vars)
}

// loadUIConfigLayers loads UI config from disk, loads the config for all referenced graphic layers, and returns
// a map of page paths to the graphic layer configs used on that page.
func loadUIConfigLayers(uiConfigPath string) (map[string][]layerConfig, error) {
	uiConfigDir := filepath.Dir(uiConfigPath)
	var uiCfg uiConfig
	if err := loadJSON(uiConfigPath, &uiCfg); err != nil {
		return nil, err
	}

	// cache of loaded layer configs by path (only load each once)
	layerConfigs := make(map[string]layerConfig)
	loadLayerConfig := func(layerPath string) (layerConfig, error) {
		if cfg, ok := layerConfigs[layerPath]; ok {
			return cfg, nil
		}

		var layerCfg layerConfig
		fullPath := filepath.Join(uiConfigDir, layerPath)
		if err := loadJSON(fullPath, &layerCfg); err != nil {
			return layerConfig{}, err
		}
		layerConfigs[layerPath] = layerCfg
		return layerCfg, nil
	}

	// visit all pages and components to find graphic layers
	pageLayers := make(map[string][]layerConfig)
	visitComponent := func(c componentConfig, path string) {
		if c.Component == "builtin:graphic/LayeredGraphic" {
			for _, layer := range c.Props.Layers {
				layerCfg, err := loadLayerConfig(layer.ConfigPath)
				if err != nil {
					log.Printf("Error loading layer config %s: %v", layer.ConfigPath, err)
					continue
				}
				pageLayers[path] = append(pageLayers[path], layerCfg)
			}
		}
	}
	var visitPage func(p pageConfig, basePath string)
	visitPage = func(p pageConfig, basePath string) {
		subPath := p.Path
		if subPath == "" {
			subPath = p.Title
		}
		fullPath := path.Join(basePath, subPath)

		for _, c := range p.Main {
			visitComponent(c, fullPath)
		}
		for _, c := range p.After {
			visitComponent(c, fullPath)
		}
		for _, child := range p.Children {
			visitPage(child, fullPath)
		}
	}
	for _, page := range uiCfg.Config.Ops.Pages {
		visitPage(page, "")
	}

	return pageLayers, nil
}

func loadJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

type deviceTrait struct {
	Name  string
	Trait trait.Name
}

func compareDeviceTrait(a, b deviceTrait) int {
	if a.Name < b.Name {
		return -1
	} else if a.Name > b.Name {
		return 1
	} else {
		return strings.Compare(a.Trait.String(), b.Trait.String())
	}
}

func layerDeviceTraits(layer layerConfig) ([]deviceTrait, error) {
	var traits []deviceTrait
	for _, elem := range layer.Elements {
		tpl, ok := layer.Templates[elem.Template.Ref]
		if !ok {
			return nil, fmt.Errorf("unknown template %q", elem.Template.Ref)
		}

		for srcName, src := range tpl.Sources {
			substitutedName, err := subTemplate(src.Request.Name, elem.Template.Vars)
			if err != nil {
				return nil, fmt.Errorf("failed to substitute template for source %q: %w", srcName, err)
			}
			traits = append(traits, deviceTrait{
				Name:  substitutedName,
				Trait: src.Trait,
			})
		}
	}

	return traits, nil
}

// subTemplate substitutes variables in the template reference and returns the resulting string.
// Variables are denoted by {{varName}} in the template string.
//
// For example, given the template string "device/{{deviceID}}/status" and vars map
// {"deviceID": "12345"}, the result would be "device/12345/status".
func subTemplate(s string, vars map[string]string) (string, error) {
	var result strings.Builder
	for len(s) > 0 {
		start := strings.Index(s, "{{")
		if start == -1 {
			// no more substitutions in the input, copy the rest
			result.WriteString(s)
			break
		}
		result.WriteString(s[:start])
		s = s[start+2:]
		end := strings.Index(s, "}}")
		if end == -1 {
			return "", fmt.Errorf("unclosed variable in template: %q", s)
		}
		varName := strings.TrimSpace(s[:end])
		varValue, ok := vars[varName]
		if !ok {
			return "", fmt.Errorf("unknown variable %q in template", varName)
		}
		result.WriteString(varValue)
		s = s[end+2:]
	}
	return result.String(), nil
}
