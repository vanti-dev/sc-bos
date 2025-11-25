package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FileUpdate represents a single line change in a file
type FileUpdate struct {
	path         string
	lineNumber   int
	originalLine string
	updatedLine  string
	replacement  string
}

// processFile scans a file and returns all necessary updates
func processFile(filePath string, presets []string, replacements []Replacement) ([]FileUpdate, error) {
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
			// Skip this replacement if it doesn't apply to the current presets
			if !appliesToPreset(repl, presets) {
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

// applyUpdates writes all updates to their respective files
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

// updateFile applies updates to a single file
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

// isGeneratedFile checks if a file is auto-generated
func isGeneratedFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Only check first 10 lines
	lineCount := 0
	for scanner.Scan() && lineCount < 10 {
		lineCount++
		line := scanner.Text()
		// Standard Go generated code marker
		if strings.Contains(line, "Code generated") && strings.Contains(line, "DO NOT EDIT") {
			return true, nil
		}
	}
	return false, scanner.Err()
}

// checkGeneratedFilesForReferences checks if any generated files contain references
func checkGeneratedFilesForReferences(files []string, presets []string, replacements []Replacement) []string {
	var filesWithRefs []string
	for _, file := range files {
		hasRef, err := fileHasReferences(file, presets, replacements)
		if err != nil {
			continue
		}
		if hasRef {
			filesWithRefs = append(filesWithRefs, file)
		}
	}
	return filesWithRefs
}

// fileHasReferences checks if a file contains any matching patterns
func fileHasReferences(filePath string, presets []string, replacements []Replacement) (bool, error) {
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
			if !appliesToPreset(repl, presets) {
				continue
			}
			if repl.pattern.MatchString(line) {
				return true, nil
			}
		}
	}
	return false, scanner.Err()
}
