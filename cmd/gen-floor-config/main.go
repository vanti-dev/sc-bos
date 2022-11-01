package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	flagData   string
	flagOutDir string
)

func init() {
	flag.StringVar(&flagData, "data-file", "config/data.json", "path to lighting counts data file")
	flag.StringVar(&flagOutDir, "out-dir", "config", "path to output directory")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	err := run(ctx)
	cancel()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	flag.Parse()

	// read and parse data file
	dataRaw, err := os.ReadFile(flagData)
	if err != nil {
		return fmt.Errorf("read %q: %w", flagData, err)
	}
	var dataParsed data
	err = json.Unmarshal(dataRaw, &dataParsed)
	if err != nil {
		return fmt.Errorf("parse %q: %w", flagData, err)
	}

	floorNo := promptInt("Which floor? ")
	floorDir := filepath.Join(flagOutDir, fmt.Sprintf("floor-%02d", floorNo))
	floorData, ok := findFloorData(dataParsed, floorNo)
	if !ok {
		return fmt.Errorf("floor %d not present in data file", floorNo)
	}
	err = os.MkdirAll(floorDir, 0777)
	if err != nil {
		return fmt.Errorf("create dir %q: %w", floorDir, err)
	}

	controllerOut := genController(floorData)
	daylightOut := controllerOut.Drivers[0]

	controllerPath := filepath.Join(floorDir, "area-controller.local.json")
	daylightPath := filepath.Join(floorDir, "daylight-power.json")

	fmt.Printf("Writing output to %q\n", floorDir)
	err = writeJSON(daylightPath, daylightOut)
	if err != nil {
		return fmt.Errorf("write %q: %w", daylightPath, err)
	}
	fmt.Printf("Wrote daylight config to %q\n", daylightPath)
	err = writeJSON(controllerPath, controllerOut)
	if err != nil {
		return fmt.Errorf("write %q, %w", controllerPath, err)
	}
	fmt.Printf("Wrote area controller local config to %q\n", controllerPath)

	return nil
}

func genFloor(floorData floor, name string) (o daliOutput) {
	o.Name = "dali"
	o.Type = "tc3dali"
	o.ADS = daliAds{
		NetID: floorData.NetID,
		Port:  851,
	}
	o.Buses = append(o.Buses, genTenant(1, floorData.T1Count, name)...)
	o.Buses = append(o.Buses, genTenant(2, floorData.T2Count, name)...)
	o.Buses = append(o.Buses, genTenant(3, floorData.T3Count, name)...)
	o.Buses = append(o.Buses, genTenant(4, floorData.T4Count, name)...)
	o.Buses = append(o.Buses, genLandlord(floorData.LandlordCount, name)...)
	if floorData.HasLifeSafety {
		o.Buses = append(o.Buses, daliBus{
			Name:         fmt.Sprintf("%s/dali/bus/LS", name),
			BridgePrefix: "GVL_Bridges.bus_LS",
		})
	}
	return
}

func genController(floorData floor) (o controllerOutput) {
	name := fmt.Sprintf("ns/bsp/sites/enterprise-wharf/floors/%d", floorData.Floor)
	driver := genFloor(floorData, name)

	o.Name = name
	o.Drivers = []daliOutput{driver}
	return
}

func genTenant(tenant int, count int, baseName string) (buses []daliBus) {
	for i := 1; i <= count; i++ {
		buses = append(buses, daliBus{
			Name:         fmt.Sprintf("%s/dali/bus/T%d_%d", baseName, tenant, i),
			BridgePrefix: fmt.Sprintf("GVL_Bridges.bus_T%d_%d", tenant, i),
		})
	}
	return
}

func genLandlord(count int, baseName string) (buses []daliBus) {
	for i := 1; i <= count; i++ {
		buses = append(buses, daliBus{
			Name:         fmt.Sprintf("%s/dali/bus/LL_%d", baseName, i),
			BridgePrefix: fmt.Sprintf("GVL_Bridges.bus_LL%d", i),
		})
	}
	return
}

func promptInt(prompt string) (result int) {
	for {
		fmt.Print(prompt)
		_, err := fmt.Scanf("%d\n", &result)
		if err == nil {
			return
		}
		fmt.Printf("INVALID INPUT, TRY AGAIN: %s\n", err.Error())
	}
}

func findFloorData(input data, floor int) (out floor, ok bool) {
	for _, f := range input.Floors {
		if f.Floor == floor {
			out = f
			ok = true
			return
		}
	}
	return
}

func writeJSON(path string, data any) error {
	encoded, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, encoded, 0777)
	return err
}

type data struct {
	Floors []floor `json:"floors"`
}

type floor struct {
	Floor         int    `json:"floor"`
	NetID         string `json:"netID"`
	HasLifeSafety bool   `json:"hasLifeSafety"`
	T1Count       int    `json:"t1Count"`
	T2Count       int    `json:"t2Count"`
	T3Count       int    `json:"t3Count"`
	T4Count       int    `json:"t4Count"`
	LandlordCount int    `json:"landlordCount"`
}

type daliOutput struct {
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	ADS   daliAds   `json:"ads"`
	Buses []daliBus `json:"buses"`
}

type daliAds struct {
	NetID string `json:"netID"`
	Port  int    `json:"port"`
}

type daliBus struct {
	Name         string `json:"name"`
	BridgePrefix string `json:"bridgePrefix"`
}

type controllerOutput struct {
	Name    string       `json:"name"`
	Drivers []daliOutput `json:"drivers"`
}
