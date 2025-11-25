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

	// Check if any skipped generated files contain references
	if len(skippedGenerated) > 0 {
		generatedWithRefs := checkGeneratedFilesForReferences(skippedGenerated, presets, replacements)
		if len(generatedWithRefs) > 0 {
			fmt.Fprintf(os.Stderr, "\n⚠️  Warning: The following generated files contain references but were skipped:\n")
			for _, file := range generatedWithRefs {
				fmt.Fprintf(os.Stderr, "  - %s\n", file)
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

// displayResults shows the changes that were found
func displayResults(updates map[string][]FileUpdate) {
	if len(updates) == 0 {
		fmt.Println("No references to vanti-dev found.")
		return
	}

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
