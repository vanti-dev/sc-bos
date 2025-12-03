package uiproto

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/smart-core-os/sc-bos/cmd/tools/genproto/internal/generator"
)

var Step = generator.Step{
	ID:   "uiproto",
	Desc: "UI protoc code generation",
	Run:  run,
}

func run(ctx *generator.Context) error {
	uiGenDir := filepath.Join(ctx.RootDir, "ui", "ui-gen")

	ctx.Verbose("Running yarn gen in %s", uiGenDir)

	if ctx.DryRun {
		ctx.Info("[DRY RUN] Would run: yarn run gen")
		return nil
	}

	cmd := exec.Command("yarn", "run", "gen")
	cmd.Dir = uiGenDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running yarn gen: %w", err)
	}

	return nil
}
