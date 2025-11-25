// Command org-migration helps migrate references from vanti-dev to smart-core-os organization.
// This tool scans local files and updates references to the old organization name with the new one.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	var skippedGenerated []string
	var skippedGoSum []string

	for _, file := range files {
		// Skip go.sum files - these should be regenerated with `go mod tidy`
		if strings.HasSuffix(file, "go.sum") {
			skippedGoSum = append(skippedGoSum, file)
			continue
		}

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
	displayResults(updates)

	// Inform about skipped go.sum files
	if len(skippedGoSum) > 0 {
		sort.Strings(skippedGoSum)
		fmt.Printf("\nSkipped %d go.sum file(s):\n", len(skippedGoSum))
		for _, file := range skippedGoSum {
			fmt.Printf("  - %s\n", file)
		}
		fmt.Println("\nRun the following command to regenerate go.sum files:")
		fmt.Println("  $ go mod tidy")
	}

	// Check if any skipped generated files contain references
	if len(skippedGenerated) > 0 {
		generatedWithRefs := checkGeneratedFilesForReferences(skippedGenerated, presets, replacements)
		if len(generatedWithRefs) > 0 {
			// If there are more than 5 files, just show unique directories
			if len(generatedWithRefs) > 5 {
				fmt.Fprintf(os.Stderr, "\nWarning: %d generated files in the following directories contain references but were skipped:\n", len(generatedWithRefs))
				displayUniqueDirectories(generatedWithRefs, os.Stderr)
			} else {
				sort.Strings(generatedWithRefs)
				fmt.Fprintf(os.Stderr, "\nWarning: %d generated file(s) contain references but were skipped:\n", len(generatedWithRefs))
				for _, file := range generatedWithRefs {
					fmt.Fprintf(os.Stderr, "  - %s\n", file)
				}
			}

			fmt.Fprintf(os.Stderr, "\nYou may need to regenerate these files after updating their source files.\n")
			fmt.Fprintf(os.Stderr, "Look for //go:generate directives or run: go generate ./...\n")
		}
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

	fmt.Println("\n✓ Successfully updated all files.")

	// Rename files that have vanti-dev in their names for the selected projects
	renamedFiles, err := renameFiles(files, projects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming files: %v\n", err)
		os.Exit(1)
	}
	if len(renamedFiles) > 0 {
		fmt.Printf("\n✓ Renamed %d file(s) with vanti-dev in their names:\n", len(renamedFiles))
		// Sort old names for consistent output
		var oldNames []string
		for oldName := range renamedFiles {
			oldNames = append(oldNames, oldName)
		}
		sort.Strings(oldNames)
		for _, oldName := range oldNames {
			fmt.Printf("  %s -> %s\n", oldName, renamedFiles[oldName])
		}
	}
}

// displayResults shows the changes that were found
func displayResults(updates map[string][]FileUpdate) {
	if len(updates) == 0 {
		fmt.Println("No references to vanti-dev found.")
		return
	}

	fmt.Printf("\nFound references in %d file(s):\n\n", len(updates))

	// Sort files for consistent output
	var files []string
	for file := range updates {
		files = append(files, file)
	}
	sort.Strings(files)

	totalUpdates := 0
	for _, file := range files {
		fileUpdates := updates[file]
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

// displayUniqueDirectories extracts and displays unique directories containing the files
func displayUniqueDirectories(files []string, output *os.File) {
	// Group files by directory to count them
	dirCounts := make(map[string]int)
	for _, file := range files {
		dir := filepath.Dir(file)
		dirCounts[dir]++
	}

	// Sort directories for consistent output
	var dirs []string
	for dir := range dirCounts {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	// Display directories with counts
	for _, dir := range dirs {
		count := dirCounts[dir]
		fmt.Fprintf(output, "  - %s/ (%d file(s))\n", dir, count)
	}

	if *verbose {
		fmt.Fprintf(output, "\nIndividual files:\n")
		// Sort files for consistent output
		sortedFiles := make([]string, len(files))
		copy(sortedFiles, files)
		sort.Strings(sortedFiles)
		for _, file := range sortedFiles {
			fmt.Fprintf(output, "  - %s\n", file)
		}
	} else {
		fmt.Fprintf(output, "  (run with --verbose to see individual files)\n")
	}
}
