package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"go.uber.org/zap"
)

var (
	daliConfigFile = flag.String("dali", "dali.json", "Dali bus configuration")
	level          = flag.Float64("level", 0, "The level to change the lights to")
)

func main() {
	fmt.Printf("Stopped: %v", run(context.Background()))
}

func run(ctx context.Context) error {
	flag.Parse()
	daliConfig, err := os.ReadFile(*daliConfigFile)
	if err != nil {
		return err
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	rootNode := node.New("lighting-power")
	lightRouter := light.NewApiRouter()
	parentRouter := parent.NewApiRouter()
	rootNode.Support(
		node.Routing(
			lightRouter,
			occupancysensor.NewApiRouter(),
			parentRouter,
		),
		node.Clients(
			light.WrapApi(lightRouter),
			parent.WrapApi(parentRouter),
		),
	)
	services := driver.Services{
		Logger: logger,
		Tasks:  &task.Group{},
		Node:   rootNode,
	}
	daliDriver := tc3dali.NewDriver(services)
	if err := daliDriver.Start(ctx); err != nil {
		return err
	}
	if err := daliDriver.Configure(daliConfig); err != nil {
		return err
	}

	var lightClient traits.LightApiClient
	if err := rootNode.Client(&lightClient); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	targetLevel := float32(*level)
	return forEachChild(ctx, rootNode, func(child *traits.Child) error {
		if !hasTrait(trait.Light, child.Traits...) {
			return nil
		}

		_, err := lightClient.UpdateBrightness(ctx, &traits.UpdateBrightnessRequest{
			Name: child.Name,
			Brightness: &traits.Brightness{
				LevelPercent: targetLevel,
			},
		})
		log.Printf("Unable to set %v to level %v: %v", child.Name, targetLevel, err)
		return nil
	})
}

func forEachChild(ctx context.Context, n *node.Node, do func(child *traits.Child) error) error {
	var parentClient traits.ParentApiClient
	if err := n.Client(&parentClient); err != nil {
		return err
	}

	req := &traits.ListChildrenRequest{
		Name: n.Name(),
	}
	// returns nextPageToken or "" if there are no more pages
	forChildrenOnPage := func(pageToken string) (string, error) {
		req.PageToken = pageToken
		children, err := parentClient.ListChildren(ctx, req)
		if err != nil {
			return "", err
		}
		for _, child := range children.Children {
			if err := do(child); err != nil {
				return "", err
			}
		}

		return children.NextPageToken, nil
	}

	var pageToken string
	for {
		var err error
		pageToken, err = forChildrenOnPage(pageToken)
		if err != nil {
			return err
		}
		if pageToken == "" {
			return nil
		}
	}
}

func hasTrait(t trait.Name, haystack ...*traits.Trait) bool {
	for _, item := range haystack {
		if item.Name == string(t) {
			return true
		}
	}
	return false
}
