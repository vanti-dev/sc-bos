// Command org-migration helps migrate references from vanti-dev to smart-core-os organization.
// This tool scans local files and updates references to the old organization name with the new one.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Setup and parse flags
	setupFlags()
	if err := parseAndValidateFlags(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Determine which file types to process
	var typesStr string
	if len(fileTypes) > 0 {
		typesStr = strings.Join(fileTypes, ",")
	} else {
		typesStr = defaultFileTypes
	}
	extensions := parseExtensions(typesStr)

	// Build replacements for the specified projects
	replacements := buildReplacements(projects)

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
	var skippedGeneratedWithRefs []string
	var skippedGoSumWithRefs []string

	for _, file := range files {
		// Skip go.sum files - these should be regenerated with `go mod tidy`
		isGoSum, hasReferences, err := isGoSumFileWithReferences(file, presets, replacements)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking if %s is go.sum: %v\n", file, err)
			continue
		}
		if isGoSum {
			// Only track go.sum files that have references
			if hasReferences {
				skippedGoSumWithRefs = append(skippedGoSumWithRefs, file)
			}
			continue
		}

		// Skip generated files
		isGenerated, hasReferences, err := isGeneratedFileWithReferences(file, presets, replacements)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking if %s is generated: %v\n", file, err)
			continue
		}
		if isGenerated {
			// Only track generated files that have references
			if hasReferences {
				skippedGeneratedWithRefs = append(skippedGeneratedWithRefs, file)
			}
			continue
		}

		fileUpdates, err := processFile(file, presets, replacements)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file, err)
			continue
		}
		if len(fileUpdates) > 0 {
			updates[file] = fileUpdates
		}
	}

	// Display results
	logResults(updates)

	// Warn about skipped files that contain references
	if len(skippedGoSumWithRefs) > 0 {
		logGoSumWarning(skippedGoSumWithRefs)
	}
	if len(skippedGeneratedWithRefs) > 0 {
		logGeneratedFilesWarning(skippedGeneratedWithRefs)
	}

	// Exit early if no changes or dry-run
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

	// Rename files that have vanti-dev in their names for the selected projects
	renamedFiles, err := renameFiles(files, projects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming files: %v\n", err)
		os.Exit(1)
	}
	logRenamedFiles(renamedFiles)

	fmt.Println("\n✓ Successfully updated all files.")
}
