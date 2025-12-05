package goproto

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/smart-core-os/sc-bos/cmd/tools/genproto/internal/generator"
)

func fixGeneratedFiles(ctx *generator.Context, genDir string) error {
	pkgDir := filepath.Join(ctx.RootDir, "pkg")
	traitDir := filepath.Join(pkgDir, "trait")

	// Check for old version fixes
	oldGenDir := filepath.Join(traitDir, "gen")
	if info, err := os.Stat(oldGenDir); err == nil && info.IsDir() {
		ctx.Verbose("Moving files from %q to %q", relPath(ctx.RootDir, oldGenDir), relPath(ctx.RootDir, genDir))
		if err := processOldVersion(ctx, oldGenDir, genDir); err != nil {
			return err
		}
	}

	// Check for new version fixes
	newGenDir := filepath.Join(traitDir, "genpb")
	if info, err := os.Stat(newGenDir); err == nil && info.IsDir() {
		ctx.Verbose("Moving files from %q to %q", relPath(ctx.RootDir, newGenDir), relPath(ctx.RootDir, genDir))
		if err := processNewVersion(ctx, newGenDir, genDir); err != nil {
			return err
		}
	}

	// Clean up trait directory
	if _, err := os.Stat(traitDir); err == nil {
		ctx.Verbose("Removing trait directory %q", relPath(ctx.RootDir, traitDir))
		if !ctx.DryRun {
			if err := os.RemoveAll(traitDir); err != nil {
				return fmt.Errorf("removing trait directory: %w", err)
			}
		}
	}

	return nil
}

func relPath(from, to string) string {
	rel, err := filepath.Rel(from, to)
	if err != nil {
		return to
	}
	return rel
}

func processOldVersion(ctx *generator.Context, srcDir, dstDir string) error {
	files, err := filepath.Glob(filepath.Join(srcDir, "*.pb.go"))
	if err != nil {
		return fmt.Errorf("finding pb.go files: %w", err)
	}

	for _, file := range files {
		if err := fixAndMoveFile(ctx, file, dstDir, false); err != nil {
			return err
		}
	}

	return nil
}

func processNewVersion(ctx *generator.Context, srcDir, dstDir string) error {
	files, err := filepath.Glob(filepath.Join(srcDir, "*.pb.go"))
	if err != nil {
		return fmt.Errorf("finding pb.go files: %w", err)
	}

	for _, file := range files {
		if err := fixAndMoveFile(ctx, file, dstDir, true); err != nil {
			return err
		}
	}

	return nil
}

func fixAndMoveFile(ctx *generator.Context, srcPath, dstDir string, isNewVersion bool) error {
	ctx.Debug("Processing %s", srcPath)

	if ctx.DryRun {
		ctx.Verbose("[DRY RUN] Would process and move file")
		return nil
	}

	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	// Apply all transformations in a single pass
	opts := RemoveGenImports | RemoveGenQualifiers
	if isNewVersion {
		opts |= RenamePackageToGen
	}
	fixed := TransformGoFile(content, opts)

	dstPath := filepath.Join(dstDir, filepath.Base(srcPath))
	if err := os.WriteFile(dstPath, fixed, 0644); err != nil {
		return fmt.Errorf("writing fixed file: %w", err)
	}

	return nil
}
