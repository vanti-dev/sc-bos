package app

import (
	"encoding/json"
	"fmt"
)

func ExampleBaseDriverConfig() {
	type ExampleDriverConfig struct {
		BaseDriverConfig
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

	for _, driver := range config.Drivers {
		switch driver.Type {
		case "example":
			var exampleConfig ExampleDriverConfig
			err := json.Unmarshal(driver.Raw, &exampleConfig)
			if err != nil {
				panic(err)
			}
			fmt.Printf("ExampleDriver name=%q property=%q\n", driver.Name, exampleConfig.Property)
		default:
			fmt.Printf("unknown driver type %q\n", driver.Type)
		}
	}
	// Output: ExampleDriver name="foo" property="bar"
}
