package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/gofrs/uuid"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/app/appconf"
	mockcfg "github.com/vanti-dev/sc-bos/pkg/driver/mock/config"
)

var (
	appconfFilepath string
	mockName        string
	lightsCount     int
)

func init() {
	flag.StringVar(&appconfFilepath, "appconfFilepath", "example/config/vanti-ugs/app.conf.json", "path to appconf file")
	flag.IntVar(&lightsCount, "lightsCount", 1_000, "number of mock lights to generate")
	flag.StringVar(&mockName, "mockName", "van/uk/brum/ugs/eg-ac-01/driver/ground", "name of the mock to generate lights into")
	flag.Parse()
}

func main() {
	cfg, err := appconf.LoadLocalConfig("", appconfFilepath)
	if err != nil {
		panic(err)
	}

	var svgs []string

	var elements []string

	for idx, driver := range cfg.Drivers {
		if driver.Name != mockName {
			continue
		}

		var devices mockcfg.Root
		if err := json.Unmarshal(driver.Raw, &devices); err != nil {
			panic(err)
		}

		for i := 0; i < lightsCount; i++ {
			id, err := uuid.NewV4()
			if err != nil {
				panic(err)
			}

			scName := fmt.Sprintf("%s%02d", "van/uk/brum/ugs/devices/LTF-L00-", i)
			svgName := fmt.Sprintf("LTF-L00-%02d", i)

			svgs = append(svgs, fmt.Sprintf(svgTemplate, svgName))

			elements = append(elements, fmt.Sprintf(`{"template": {"ref": "spotGroup", "el": "%s", "sc": "%s"}}`, svgName, scName))

			devices.Devices = append(devices.Devices, mockcfg.Device{
				Metadata: &traits.Metadata{
					Name: scName,
					Traits: []*traits.TraitMetadata{
						{
							Name: "smartcore.traits.Light",
						},
						{
							Name: "smartcore.bos.Status",
						},
					},
					Appearance: &traits.Metadata_Appearance{Title: id.String()[0:8]},
					Membership: &traits.Metadata_Membership{Subsystem: "lighting"},
					Location: &traits.Metadata_Location{
						Floor: "Ground Floor",
					},
				},
			})
		}

		devJson, err := json.MarshalIndent(devices, "", "  ")
		if err != nil {
			panic(err)
		}
		driver.Raw = devJson
		cfg.Drivers[idx] = driver
	}

	outFile, err := os.OpenFile(appconfFilepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := outFile.Close(); err != nil {
			panic(err)
		}
	}()

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		panic(err)
	}

	outSvgFile, err := os.OpenFile("lights-svg-snippets.svg", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := outSvgFile.Close(); err != nil {
			panic(err)
		}
	}()

	svgHeader := `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`
	svgFooter := `</svg>`

	if _, err := outSvgFile.WriteString(svgHeader); err != nil {
		panic(err)
	}

	for _, svg := range svgs {
		if _, err := outSvgFile.WriteString(svg); err != nil {
			panic(err)
		}
	}
	if _, err := outSvgFile.WriteString(svgFooter); err != nil {
		panic(err)
	}

	fmt.Println(elements)
}

var svgTemplate = `
<g id="%s">
<g id="Lights-Single-Circle-134:826">
<rect x="569.5" y="556.5" width="17" height="17" rx="8.5" fill="#2A3138"/>
<rect x="569.5" y="556.5" width="17" height="17" rx="8.5" stroke="#FBBF24"/>
<ellipse id="light-widget_26" cx="578" cy="565" rx="6" ry="6" transform="rotate(-90 578 565)" fill="#4A4F55" stroke="#4A4F55" stroke-width="0.559238" stroke-linejoin="round"/>
</g>
<g id="Lights-Single-Circle-134:827">
<rect x="569.5" y="625.5" width="17" height="17" rx="8.5" fill="#2A3138"/>
<rect x="569.5" y="625.5" width="17" height="17" rx="8.5" stroke="#FBBF24"/>
<ellipse id="light-widget_27" cx="578" cy="634" rx="6" ry="6" transform="rotate(-90 578 634)" fill="#4A4F55" stroke="#4A4F55" stroke-width="0.559238" stroke-linejoin="round"/>
</g>
<g id="Lights-Single-Circle-134:828">
<rect x="569.5" y="691.5" width="17" height="17" rx="8.5" fill="#2A3138"/>
<rect x="569.5" y="691.5" width="17" height="17" rx="8.5" stroke="#FBBF24"/>
<ellipse id="light-widget_28" cx="578" cy="700" rx="6" ry="6" transform="rotate(-90 578 700)" fill="#4A4F55" stroke="#4A4F55" stroke-width="0.559238" stroke-linejoin="round"/>
</g>
</g>`
