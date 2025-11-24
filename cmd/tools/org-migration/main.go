// Command org-migration helps migrate references from vanti-dev to smart-core-os organization.
// This tool scans local files and updates references to the old organization name with the new one.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Flags
	dryRun    = flag.Bool("dry-run", false, "Show what would be changed without making changes")
	verbose   = flag.Bool("verbose", false, "Show detailed output")
	path      = flag.String("path", ".", "Path to scan for files")
	fileTypes = flag.String("types", "", "Comma-separated list of file extensions to process (default: all supported types)")
	preset    = flag.String("preset", "all", "Preset category: all, go, js, docker, docs")
)

// Default file extensions to scan when --types is not specified
const defaultFileTypes = "go,mod,sum,sh,yml,yaml,json,md,js,ts,jsx,tsx,mjs,cjs,vue,html,xml,Dockerfile"

// Preset categories define which replacements to apply
var presetInfo = map[string]string{
	"all":    "All replacements (default)",
	"go":     "Go import path replacements only",
	"js":     "JavaScript/npm package replacements only",
	"docker": "Docker image replacements only",
	"docs":   "GitHub URL replacements only",
}

// Replacement represents a pattern replacement with preset restrictions
type Replacement struct {
	name    string
	pattern *regexp.Regexp
	repl    string
	presets []string // Which presets this replacement applies to (empty means all)
}

// Replacement patterns for different contexts
// Order matters: more specific patterns should come first
var replacements = []Replacement{
	{
		name:    "GHCR image path",
		pattern: regexp.MustCompile(`ghcr\.io/vanti-dev/sc-bos`),
		repl:    "ghcr.io/smart-core-os/sc-bos",
		presets: []string{"docker", "all"},
	},
	{
		name:    "npm package",
		pattern: regexp.MustCompile(`@vanti-dev/sc-bos`),
		repl:    "@smart-core-os/sc-bos",
		presets: []string{"js", "all"},
	},
	{
		name:    "GitHub raw URL",
		pattern: regexp.MustCompile(`https?://raw\.githubusercontent\.com/vanti-dev/sc-bos`),
		repl:    "https://raw.githubusercontent.com/smart-core-os/sc-bos",
		presets: []string{"docs", "all"},
	},
	{
		name:    "GitHub URL",
		pattern: regexp.MustCompile(`https?://github\.com/vanti-dev/sc-bos`),
		repl:    "https://github.com/smart-core-os/sc-bos",
		presets: []string{"docs", "all"},
	},
	{
		name:    "Go import path",
		pattern: regexp.MustCompile(`(^|[^:/])github\.com/vanti-dev/sc-bos`), // special handling to avoid matching URLs
		repl:    "${1}github.com/smart-core-os/sc-bos",
		presets: []string{"go", "all"},
	},
}

type FileUpdate struct {
	path         string
	lineNumber   int
	originalLine string
	updatedLine  string
	replacement  string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Migrates references from vanti-dev to smart-core-os organization.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nPresets (control which replacements to apply):\n")
		// Generate preset list dynamically from presetInfo
		presetOrder := []string{"all", "go", "js", "docker", "docs"}
		for _, name := range presetOrder {
			if desc, ok := presetInfo[name]; ok {
				fmt.Fprintf(os.Stderr, "  %-6s - %s\n", name, desc)
			}
		}
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --dry-run --preset go --path ~/my-project\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --preset js --path ~/my-project\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --preset docker --verbose --path ~/my-project\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --types go,mod,yaml --path .\n", os.Args[0])
	}
	flag.Parse()

	// Validate preset
	if _, ok := presetInfo[*preset]; !ok {
		fmt.Fprintf(os.Stderr, "Unknown preset: %s\n", *preset)
		fmt.Fprintf(os.Stderr, "Available presets: all, go, js, docker, docs\n")
		os.Exit(1)
	}

	// Determine which file types to process
	var typesStr string
	if *fileTypes != "" {
		// User-specified types
		typesStr = *fileTypes
	} else {
		// Default to all supported file types
		typesStr = defaultFileTypes
	}

	// Parse file types
	extensions := parseExtensions(typesStr)

	// Collect files to process
	files, err := collectFiles(*path, extensions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error collecting files: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Scanning %d files in %s\n", len(files), *path)
	}

	// Process files
	updates := make(map[string][]FileUpdate)
	var skippedGenerated []string
	for _, file := range files {
		// Skip generated files
		isGenerated, err := isGeneratedFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking if %s is generated: %v\n", file, err)
			continue
		}
		if isGenerated {
			skippedGenerated = append(skippedGenerated, file)
			continue
		}

		fileUpdates, err := processFile(file, *preset)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file, err)
			continue
		}
		if len(fileUpdates) > 0 {
			updates[file] = fileUpdates
		}
	}

	// Display results
	if len(updates) == 0 {
		fmt.Println("No references to vanti-dev found.")
	} else {
		fmt.Printf("\nFound references in %d file(s):\n\n", len(updates))

		totalUpdates := 0
		for file, fileUpdates := range updates {
			totalUpdates += len(fileUpdates)
			fmt.Printf("%s: (%d change(s))\n", file, len(fileUpdates))
			if *verbose {
				for _, update := range fileUpdates {
					fmt.Printf("  Line %d [%s]:\n", update.lineNumber, update.replacement)
					fmt.Printf("    - %s\n", strings.TrimSpace(update.originalLine))
					fmt.Printf("    + %s\n", strings.TrimSpace(update.updatedLine))
				}
			}
		}

		fmt.Printf("\nTotal: %d change(s) in %d file(s)\n", totalUpdates, len(updates))
	}

	// Check if any skipped generated files contain references that would need regeneration
	if len(skippedGenerated) > 0 {
		generatedWithRefs := checkGeneratedFilesForReferences(skippedGenerated, *preset)
		if len(generatedWithRefs) > 0 {
			fmt.Printf("\n⚠️  Warning: %d generated file(s) contain vanti-dev references:\n", len(generatedWithRefs))

			// Group files by directory for cleaner output
			if len(generatedWithRefs) > 5 {
				dirs := groupFilesByDirectory(generatedWithRefs)
				for dir, count := range dirs {
					if count == 1 {
						fmt.Printf("  - %s (1 file)\n", dir)
					} else {
						fmt.Printf("  - %s (%d files)\n", dir, count)
					}
				}
			} else {
				// Show individual files if there aren't many
				for _, file := range generatedWithRefs {
					fmt.Printf("  - %s\n", file)
				}
			}

			fmt.Println("\nThese files were not modified because they are generated code.")
			fmt.Println("You will need to regenerate them after updating the source files.")
			fmt.Println("Look for '//go:generate' directives or run 'go generate ./...' in the repository.")
		}
	}

	if len(updates) == 0 {
		return
	}

	if *dryRun {
		fmt.Println("\n[DRY RUN] No changes were made. Run without --dry-run to apply changes.")
		return
	}

	// Apply updates
	if err := applyUpdates(updates); err != nil {
		fmt.Fprintf(os.Stderr, "Error applying updates: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ Successfully updated all files.")

	// Rename files that have vanti-dev in their names
	renamedFiles, err := renameFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming files: %v\n", err)
		os.Exit(1)
	}
	if len(renamedFiles) > 0 {
		fmt.Printf("\n✓ Renamed %d file(s) with vanti-dev in their names:\n", len(renamedFiles))
		for oldName, newName := range renamedFiles {
			fmt.Printf("  %s -> %s\n", oldName, newName)
		}
	}
}

func parseExtensions(typesStr string) map[string]bool {
	extensions := make(map[string]bool)
	for _, ext := range strings.Split(typesStr, ",") {
		ext = strings.TrimSpace(ext)
		if ext != "" {
			extensions[ext] = true
		}
	}
	return extensions
}

func appliesToPreset(repl Replacement, preset string) bool {
	// Empty presets list means applies to all
	if len(repl.presets) == 0 {
		return true
	}
	// Check if the current preset is in the replacement's preset list
	for _, p := range repl.presets {
		if p == preset {
			return true
		}
	}
	return false
}

func isGeneratedFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Check first few lines for generated code marker
	lineCount := 0
	for scanner.Scan() && lineCount < 5 {
		lineCount++
		line := scanner.Text()
		// Standard Go generated code marker
		if strings.Contains(line, "Code generated") && strings.Contains(line, "DO NOT EDIT") {
			return true, nil
		}
	}
	return false, scanner.Err()
}

func checkGeneratedFilesForReferences(files []string, preset string) []string {
	var filesWithRefs []string
	for _, file := range files {
		hasRef, err := fileHasReferences(file, preset)
		if err != nil {
			continue
		}
		if hasRef {
			filesWithRefs = append(filesWithRefs, file)
		}
	}
	return filesWithRefs
}

func fileHasReferences(filePath string, preset string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if any replacement pattern matches
		for _, repl := range replacements {
			if !appliesToPreset(repl, preset) {
				continue
			}
			if repl.pattern.MatchString(line) {
				return true, nil
			}
		}
	}
	return false, scanner.Err()
}

func groupFilesByDirectory(files []string) map[string]int {
	dirs := make(map[string]int)
	for _, file := range files {
		dir := filepath.Dir(file)
		dirs[dir]++
	}
	return dirs
}

func collectFiles(rootPath string, extensions map[string]bool) ([]string, error) {
	var files []string

	// Directories to skip during traversal
	// Note: .run is NOT skipped to allow updating IDEA run configurations
	ignoredDirs := []string{"node_modules", "vendor", "dist", "org-migration"}

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and common ignored directories
		if d.IsDir() {
			name := d.Name()
			// Don't skip the root directory even if it's "."
			if path != rootPath {
				// Allow .run directory (IDEA run configurations)
				if name == ".run" {
					return nil
				}
				// Skip other hidden directories (starting with .)
				if strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
				// Skip explicitly ignored directories
				for _, ignored := range ignoredDirs {
					if name == ignored {
						return filepath.SkipDir
					}
				}
			}
			return nil
		}

		// Check if file extension matches
		name := d.Name()
		// Handle Dockerfile specially (no extension)
		if name == "Dockerfile" || strings.HasPrefix(name, "Dockerfile-") {
			files = append(files, path)
			return nil
		}

		// Check regular extensions
		ext := strings.TrimPrefix(filepath.Ext(name), ".")
		if ext != "" && extensions[ext] {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func processFile(filePath string, preset string) ([]FileUpdate, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var updates []FileUpdate
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		updatedLine := line

		// Try each replacement pattern
		for _, repl := range replacements {
			// Skip this replacement if it doesn't apply to the current preset
			if !appliesToPreset(repl, preset) {
				continue
			}

			if repl.pattern.MatchString(updatedLine) {
				newLine := repl.pattern.ReplaceAllString(updatedLine, repl.repl)
				if newLine != updatedLine {
					updates = append(updates, FileUpdate{
						path:         filePath,
						lineNumber:   lineNumber,
						originalLine: line,
						updatedLine:  newLine,
						replacement:  repl.name,
					})
					updatedLine = newLine
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return updates, nil
}

func applyUpdates(updates map[string][]FileUpdate) error {
	for filePath, fileUpdates := range updates {
		if err := updateFile(filePath, fileUpdates); err != nil {
			return fmt.Errorf("failed to update %s: %w", filePath, err)
		}
		if *verbose {
			fmt.Printf("Updated %s\n", filePath)
		}
	}
	return nil
}

func updateFile(filePath string, updates []FileUpdate) error {
	// Read the entire file once
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Detect line ending style
	lineEnding := "\n"
	if strings.Contains(string(content), "\r\n") {
		lineEnding = "\r\n"
	}

	// Split into lines
	lines := strings.Split(string(content), lineEnding)

	// Apply updates (updates are 1-indexed)
	updateMap := make(map[int]string)
	for _, update := range updates {
		updateMap[update.lineNumber] = update.updatedLine
	}

	for i := range lines {
		if newLine, ok := updateMap[i+1]; ok {
			lines[i] = newLine
		}
	}

	// Write back to file
	output := strings.Join(lines, lineEnding)

	return os.WriteFile(filePath, []byte(output), 0644)
}

func renameFiles(files []string) (map[string]string, error) {
	renamedFiles := make(map[string]string)

	for _, filePath := range files {
		// Check if filename contains vanti-dev
		fileName := filepath.Base(filePath)
		if !strings.Contains(fileName, "vanti-dev") {
			continue
		}

		// Create new filename by replacing vanti-dev with smart-core-os
		newFileName := strings.ReplaceAll(fileName, "vanti-dev", "smart-core-os")
		newFilePath := filepath.Join(filepath.Dir(filePath), newFileName)

		// Rename the file
		if err := os.Rename(filePath, newFilePath); err != nil {
			return renamedFiles, fmt.Errorf("failed to rename %s to %s: %w", filePath, newFilePath, err)
		}

		renamedFiles[filePath] = newFilePath
		if *verbose {
			fmt.Printf("Renamed %s to %s\n", filePath, newFilePath)
		}
	}

	return renamedFiles, nil
}
