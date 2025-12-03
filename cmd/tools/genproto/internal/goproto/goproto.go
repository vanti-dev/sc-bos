package goproto

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/smart-core-os/sc-bos/cmd/tools/genproto/internal/generator"
)

var Step = generator.Step{
	ID:   "goproto",
	Desc: "Go protoc code generation",
	Run:  run,
}

func run(ctx *generator.Context) error {
	genDir := filepath.Join(ctx.RootDir, "pkg", "gen")

	ctx.Verbose("Running go generate in %s", genDir)

	if ctx.DryRun {
		ctx.Info("[DRY RUN] Would run: go generate")
	} else {
		cmd := exec.Command("go", "generate")
		cmd.Dir = genDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("running go generate: %w", err)
		}
	}

	// Fix generated files
	if err := fixGeneratedFiles(ctx, genDir); err != nil {
		return err
	}

	return nil
}
