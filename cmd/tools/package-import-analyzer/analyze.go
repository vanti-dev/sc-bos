package main

import (
	"flag"
	"log"
	"sort"
)

// runAnalyze executes the analyze command
func runAnalyze() {
	if *analyzeInput == "" {
		log.Fatal("Input file is required. Use -input flag.")
	}

	result, err := loadScanResult(*analyzeInput)
	if err != nil {
		log.Fatalf("Failed to load scan result: %v", err)
	}

	if *verbose {
		log.Printf("Loaded scan results from %s", *analyzeInput)
		log.Printf("Base module: %s, Organization: %s", result.BaseModule, result.Organization)
		log.Printf("Scanned on: %s", result.Timestamp.Format("2006-01-02 15:04:05"))
	}

	switch *analyzeType {
	case "all":
		runAllAnalysis(result)
	case "package-types":
		runPackageTypeAnalysis(result)
	case "repo-dependencies":
		runRepoDependenciesAnalysis(result)
	default:
		log.Fatalf("Unknown analysis type: %s", *analyzeType)
	}
}

func runAllAnalysis(result *ScanResult) {
	runPackageTypeAnalysis(result)
	println()
	runRepoDependenciesAnalysis(result)
	println()
	runBasicSummary(result)
}

func runPackageTypeAnalysis(result *ScanResult) {
	analysis := analyzePackageTypes(result)

	if err := outputPackageTypeAnalysis(analysis); err != nil {
		log.Fatalf("Failed to output analysis: %v", err)
	}
}

func runRepoDependenciesAnalysis(result *ScanResult) {
	analysis := analyzeRepoDependencies(result)

	if err := outputRepoDependencyAnalysis(analysis); err != nil {
		log.Fatalf("Failed to output analysis: %v", err)
	}
}

func runBasicSummary(result *ScanResult) {
	if err := outputBasicSummary(result); err != nil {
		log.Fatalf("Failed to output summary: %v", err)
	}
}

func analyzePackageTypes(result *ScanResult) []PackageTypeAnalysis {
	// Map of type -> list of import locations (repo:file pairs)
	typeImports := make(map[string]map[string]bool)    // type -> unique "repo:file" strings
	typeDependents := make(map[string]map[string]bool) // type -> dependent repos
	typePackages := make(map[string]map[string]bool)   // type -> packages being imported

	// Iterate through all package usage
	for pkg, importInfos := range result.PackageUsage {
		for _, info := range importInfos {
			// Classify each file based on where it's located in the dependent repo
			for _, filePath := range info.Files {
				fileType := classifyImportByFilePath(filePath)

				// Initialize maps if needed
				if typeImports[fileType] == nil {
					typeImports[fileType] = make(map[string]bool)
					typeDependents[fileType] = make(map[string]bool)
					typePackages[fileType] = make(map[string]bool)
				}

				// Track the import location
				importLocation := info.Repository + ":" + filePath
				typeImports[fileType][importLocation] = true

				// Track the dependent repository
				typeDependents[fileType][info.Repository] = true

				// Track the package being imported
				typePackages[fileType][pkg] = true
			}
		}
	}

	// Build the analysis results
	var analysis []PackageTypeAnalysis
	for fileType, imports := range typeImports {
		// Get unique packages for this type
		pkgList := make([]string, 0, len(typePackages[fileType]))
		for pkg := range typePackages[fileType] {
			pkgList = append(pkgList, pkg)
		}
		sort.Strings(pkgList)

		analysis = append(analysis, PackageTypeAnalysis{
			PackageType:    fileType,
			Count:          len(imports), // Count of unique file locations
			Packages:       pkgList,
			DependentRepos: len(typeDependents[fileType]),
		})
	}

	// Sort by count (descending)
	sort.Slice(analysis, func(i, j int) bool {
		return analysis[i].Count > analysis[j].Count
	})

	return analysis
}

func analyzeRepoDependencies(result *ScanResult) []RepoDependencyAnalysis {
	// Group by importing repository
	repoPackages := make(map[string]map[string]bool) // repo -> packages

	for pkg, importInfos := range result.PackageUsage {
		for _, info := range importInfos {
			if repoPackages[info.Repository] == nil {
				repoPackages[info.Repository] = make(map[string]bool)
			}
			repoPackages[info.Repository][pkg] = true
		}
	}

	var analysis []RepoDependencyAnalysis
	for repo, packages := range repoPackages {
		analysis = append(analysis, RepoDependencyAnalysis{
			Repository:   repo,
			PackageCount: len(packages),
		})
	}

	sort.Slice(analysis, func(i, j int) bool {
		return analysis[i].PackageCount > analysis[j].PackageCount
	})

	return analysis
}

var (
	// Analyze subcommand flags
	analyzeCmd    = flag.NewFlagSet("analyze", flag.ExitOnError)
	analyzeInput  = analyzeCmd.String("input", "", "Input file path with scan results (required)")
	analyzeOutput = analyzeCmd.String("output", "text", "Output format: json, csv, or text")
	analyzeFile   = analyzeCmd.String("output-file", "", "Output file path (if not specified, outputs to stdout)")
	analyzeType   = analyzeCmd.String("type", "all", "Analysis type: all, package-types, repo-dependencies")
)
