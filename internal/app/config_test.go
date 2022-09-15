package app

import (
	"encoding/json"
	"fmt"

	"github.com/vanti-dev/bsp-ew/internal/driver"
)

func ExampleControllerConfig() {
	type ExampleDriverConfig struct {
		driver.BaseConfig
		Property string `json:"property"`
	}

	// language=json
	buf := []byte(`{
	"drivers": [
		{
			"name": "foo",
			"type": "example",
			"property": "bar"
		}
    ]
}`)
	var config ControllerConfig
	err := json.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	for _, d := range config.Drivers {
		switch d.Type {
		case "example":
			var exampleConfig ExampleDriverConfig
			err := json.Unmarshal(d.Raw, &exampleConfig)
			if err != nil {
				panic(err)
			}
			fmt.Printf("ExampleDriver name=%q property=%q\n", d.Name, exampleConfig.Property)
		default:
			fmt.Printf("unknown driver type %q\n", d.Type)
		}
	}
	// Output: ExampleDriver name="foo" property="bar"
}
