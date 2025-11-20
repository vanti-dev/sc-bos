// Command package-import-analyzer analyzes GitHub repositories in an organization
// to find which packages from this repository are imported by other repos.
package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	org         = flag.String("org", "", "GitHub organization name to scan")
	authToken   = flag.String("token", "", "GitHub personal access token (or set GITHUB_TOKEN env var)")
	baseModule  = flag.String("base-module", "github.com/vanti-dev/sc-bos", "Base module to look for imports")
	verbose     = flag.Bool("verbose", false, "Enable verbose output")
	veryVerbose = flag.Bool("vv", false, "Enable very verbose output (includes cache hits/misses and API calls)")
	output      = flag.String("output", "json", "Output format: json, csv, or text")
	outputFile  = flag.String("output-file", "", "Output file path (if not specified, outputs to stdout)")
	cacheDir    = flag.String("cache-dir", "", "Cache directory for API responses (default: ~/.package-import-analyzer-cache)")
)

// CacheEntry holds cached API response with metadata
type CacheEntry struct {
	Data      any       `json:"data"`
	CachedAt  time.Time `json:"cached_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// APICache handles persistent caching of API responses
type APICache struct {
	cacheDir  string
	apiCalls  int // track total API calls made
	cacheHits int // track cache hits
}

type GitHubRepo struct {
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	CloneURL      string    `json:"clone_url"`
	Language      string    `json:"language"`
	DefaultBranch string    `json:"default_branch"`
	PushedAt      time.Time `json:"pushed_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GitTreeItem struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type GitTree struct {
	SHA       string        `json:"sha"`
	URL       string        `json:"url"`
	Tree      []GitTreeItem `json:"tree"`
	Truncated bool          `json:"truncated"`
}

type ImportInfo struct {
	Repository string    `json:"repository"`
	Package    string    `json:"package"`
	Files      []string  `json:"files"`
	LastPush   time.Time `json:"last_push"`
	LastUpdate time.Time `json:"last_update"`
}

type AnalysisResult struct {
	BaseModule       string                  `json:"base_module"`
	Organization     string                  `json:"organization"`
	TotalRepos       int                     `json:"total_repos"`
	ReposWithGoFiles int                     `json:"repos_with_go_files"`
	ReposWithImports int                     `json:"repos_with_imports"`
	PackageUsage     map[string][]ImportInfo `json:"package_usage"`
	Timestamp        time.Time               `json:"timestamp"`
}

func main() {
	flag.Parse()

	if *org == "" {
		log.Fatal("Organization name is required. Use -org flag.")
	}

	githubToken := *authToken
	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
	}
	if githubToken == "" {
		log.Fatal("GitHub token is required. Use -token flag or set GITHUB_TOKEN environment variable.")
	}

	ctx := context.Background()

	if *verbose {
		log.Printf("Analyzing imports from %s in organization %s", *baseModule, *org)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	// Initialize persistent cache
	cache, err := newAPICache(*cacheDir)
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	if *verbose {
		log.Printf("Using cache directory: %s", cache.cacheDir)
	}

	repos, err := getOrgRepos(ctx, client, *org, githubToken, cache)
	if err != nil {
		log.Fatalf("Failed to get repositories: %v", err)
	}

	if *verbose {
		log.Printf("Found %d repositories (API calls: %d, cache hits: %d)",
			len(repos), cache.apiCalls, cache.cacheHits)
	}

	result := AnalysisResult{
		BaseModule:   *baseModule,
		Organization: *org,
		TotalRepos:   len(repos),
		PackageUsage: make(map[string][]ImportInfo),
		Timestamp:    time.Now(),
	}

	reposWithGoFiles := 0
	reposWithImports := 0

	// Extract the repository name from the base module to skip it
	baseRepoName := extractRepoNameFromModule(*baseModule)

	for _, repo := range repos {
		// Skip the repository that contains the base module
		if baseRepoName != "" && repo.FullName == baseRepoName {
			if *verbose {
				log.Printf("Skipping %s - this is the base module repository", repo.FullName)
			}
			continue
		}

		if *verbose {
			log.Printf("Analyzing repository: %s (language: %s, last push: %s)",
				repo.FullName, repo.Language, repo.PushedAt.Format("2006-01-02"))
		}

		imports, hasGoFiles, err := analyzeRepoImports(ctx, client, repo, githubToken, *baseModule, cache)
		if err != nil {
			if *verbose {
				log.Printf("Failed to analyze %s: %v (API calls: %d, cache hits: %d)",
					repo.FullName, err, cache.apiCalls, cache.cacheHits)
			}
			continue
		}

		if hasGoFiles {
			reposWithGoFiles++
		}

		if len(imports) > 0 {
			reposWithImports++
			for pkg, files := range imports {
				result.PackageUsage[pkg] = append(result.PackageUsage[pkg], ImportInfo{
					Repository: repo.FullName,
					Package:    pkg,
					Files:      files,
					LastPush:   repo.PushedAt,
					LastUpdate: repo.UpdatedAt,
				})
			}

			if *verbose {
				log.Printf("Found %d packages in %s (API calls: %d, cache hits: %d)",
					len(imports), repo.FullName, cache.apiCalls, cache.cacheHits)
			}
		}
	}

	result.ReposWithGoFiles = reposWithGoFiles
	result.ReposWithImports = reposWithImports

	if *verbose {
		log.Printf("Analysis complete. Total API calls: %d, cache hits: %d (%.1f%% cache hit rate)",
			cache.apiCalls, cache.cacheHits,
			float64(cache.cacheHits)/float64(cache.apiCalls+cache.cacheHits)*100)
	}

	if err := outputResult(result); err != nil {
		log.Fatalf("Failed to output result: %v", err)
	}
}

// newAPICache creates a new API cache with the specified cache directory
func newAPICache(cacheDir string) (*APICache, error) {
	if cacheDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %v", err)
		}
		cacheDir = filepath.Join(homeDir, ".package-import-analyzer-cache")
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %v", err)
	}

	return &APICache{
		cacheDir:  cacheDir,
		apiCalls:  0,
		cacheHits: 0,
	}, nil
}

// getCacheKey generates a cache key for a URL
func (c *APICache) getCacheKey(url string) string {
	hash := md5.Sum([]byte(url))
	return fmt.Sprintf("%x", hash)
}

// getCacheFilePath returns the full path to a cache file
// Files are organized into subdirectories named by the first two characters of the hash
func (c *APICache) getCacheFilePath(cacheKey string) string {
	// Use first two characters as subdirectory name
	subDir := cacheKey[:2]
	dirPath := filepath.Join(c.cacheDir, subDir)

	// Create subdirectory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		// Fall back to flat structure if subdirectory creation fails
		return filepath.Join(c.cacheDir, cacheKey+".json")
	}

	return filepath.Join(dirPath, cacheKey+".json")
}

// get retrieves data from cache if it exists and is not expired
func (c *APICache) get(url string, target any) bool {
	cacheKey := c.getCacheKey(url)
	filePath := c.getCacheFilePath(cacheKey)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return false // Cache miss
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return false // Invalid cache entry
	}

	// Check if cache entry is expired
	if time.Now().After(entry.ExpiresAt) {
		os.Remove(filePath) // Clean up expired cache
		return false
	}

	// Decode the cached data into target
	entryData, err := json.Marshal(entry.Data)
	if err != nil {
		return false
	}

	if err := json.Unmarshal(entryData, target); err != nil {
		return false
	}

	c.cacheHits++
	if *veryVerbose {
		log.Printf("Cache HIT for %s", url)
	}
	return true
}

// set stores data in cache with expiration time
func (c *APICache) set(url string, data any, ttl time.Duration) error {
	cacheKey := c.getCacheKey(url)
	filePath := c.getCacheFilePath(cacheKey)

	entry := CacheEntry{
		Data:      data,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}

	entryData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, entryData, 0644)
}

func getOrgRepos(ctx context.Context, client *http.Client, org, githubToken string, cache *APICache) ([]GitHubRepo, error) {
	var allRepos []GitHubRepo
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?page=%d&per_page=%d", org, page, perPage)

		var repos []GitHubRepo
		if err := httpGetJSONWithCache(ctx, client, url, githubToken, &repos, cache, 30*time.Minute); err != nil {
			return nil, err
		}

		if len(repos) == 0 {
			break
		}

		allRepos = append(allRepos, repos...)

		if len(repos) < perPage {
			break
		}
		page++
	}

	return allRepos, nil
}

// extractRepoNameFromModule extracts the repository owner/name from a Go module path
// For example: "github.com/vanti-dev/sc-bos" -> "vanti-dev/sc-bos"
func extractRepoNameFromModule(modulePath string) string {
	// Remove common prefixes and extract owner/repo
	parts := strings.Split(modulePath, "/")
	if len(parts) < 3 {
		return ""
	}

	// Handle github.com/owner/repo format
	if parts[0] == "github.com" && len(parts) >= 3 {
		return parts[1] + "/" + parts[2]
	}

	return ""
}

// httpGetJSONWithCache checks cache first, then makes API call if needed
func httpGetJSONWithCache(ctx context.Context, client *http.Client, url, githubToken string, target any, cache *APICache, ttl time.Duration) error {
	// Try cache first
	if cache.get(url, target) {
		return nil // Cache hit
	}

	// Cache miss - make API call
	if err := httpGetJSON(ctx, client, url, githubToken, target, cache); err != nil {
		return err
	}

	// Cache the response
	if err := cache.set(url, target, ttl); err != nil && *veryVerbose {
		log.Printf("Failed to cache response for %s: %v", url, err)
	}

	return nil
}

// httpGetJSON makes an authenticated HTTP GET request to the GitHub API and decodes the JSON response
func httpGetJSON(ctx context.Context, client *http.Client, url, githubToken string, target any, cache *APICache) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Track API call
	cache.apiCalls++
	if *veryVerbose {
		log.Printf("API CALL to %s (total: %d)", url, cache.apiCalls)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: %d - %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// repoHasGoFiles checks if a repository has Go files using the GitHub Languages API
func repoHasGoFiles(ctx context.Context, client *http.Client, repo GitHubRepo, githubToken string, cache *APICache) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/languages", repo.FullName)

	var languages map[string]int
	if err := httpGetJSONWithCache(ctx, client, url, githubToken, &languages, cache, 24*time.Hour); err != nil {
		if err.Error() == "not found" {
			return false, nil // Repository not accessible or empty
		}
		return false, err
	}

	// Check if "Go" is in the languages map
	_, hasGo := languages["Go"]
	return hasGo, nil
}

// repoHasModuleDependency checks if any go.mod file in the repository contains the base module as a dependency
// Returns the list of module directories (directories containing go.mod with the dependency)
// This helps avoid scanning repositories that don't use the module at all
func repoHasModuleDependency(ctx context.Context, client *http.Client, repo GitHubRepo, githubToken, baseModule string, cache *APICache, tree *GitTree) ([]string, error) {
	// Find all go.mod files in the tree
	var goModPaths []string
	for _, item := range tree.Tree {
		if item.Type == "blob" && strings.HasSuffix(item.Path, "go.mod") {
			goModPaths = append(goModPaths, item.Path)
		}
	}

	// If no go.mod files found, return nil to scan all files (conservative approach)
	if len(goModPaths) == 0 {
		if *verbose {
			log.Printf("No go.mod files found in %s, will scan all Go files", repo.FullName)
		}
		return nil, nil
	}

	// Check each go.mod file for the dependency and collect matching module directories
	var moduleDirs []string
	for _, goModPath := range goModPaths {
		hasDep, err := checkGoModFile(ctx, client, repo, githubToken, baseModule, goModPath, cache)
		if err != nil {
			if *verbose {
				log.Printf("Failed to check %s/%s: %v", repo.FullName, goModPath, err)
			}
			continue // Try other go.mod files
		}
		if hasDep {
			// Extract the directory containing this go.mod
			moduleDir := filepath.Dir(goModPath)
			if moduleDir == "." {
				moduleDir = "" // Root directory
			}
			moduleDirs = append(moduleDirs, moduleDir)
			if *verbose {
				log.Printf("Found dependency on %s in %s/%s", baseModule, repo.FullName, goModPath)
			}
		}
	}

	return moduleDirs, nil
}

// checkGoModFile checks if a specific go.mod file contains the base module as a dependency
func checkGoModFile(ctx context.Context, client *http.Client, repo GitHubRepo, githubToken, baseModule, goModPath string, cache *APICache) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo.FullName, goModPath)

	var fileData map[string]any
	if err := httpGetJSONWithCache(ctx, client, url, githubToken, &fileData, cache, 6*time.Hour); err != nil {
		return false, err
	}

	content, ok := fileData["content"].(string)
	if !ok {
		return false, fmt.Errorf("no content in go.mod response")
	}

	// Decode base64 content
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(content, "\n", ""))
	if err != nil {
		return false, err
	}

	// Simple string search for the base module in go.mod
	// This is faster than parsing the entire go.mod file
	goModContent := string(decoded)
	return strings.Contains(goModContent, baseModule), nil
}

func analyzeRepoImports(ctx context.Context, client *http.Client, repo GitHubRepo, githubToken, baseModule string, cache *APICache) (map[string][]string, bool, error) {
	// First check if the repository has Go files using the Languages API
	hasGoFiles, err := repoHasGoFiles(ctx, client, repo, githubToken, cache)
	if err != nil {
		return nil, false, err
	}

	// If no Go files, skip scanning
	if !hasGoFiles {
		return make(map[string][]string), false, nil
	}

	// Get the Git Tree for the repository (we'll use this for both go.mod check and finding Go files)
	tree, err := getRepoTree(ctx, client, repo, githubToken, cache)
	if err != nil {
		return nil, hasGoFiles, err
	}

	if tree == nil {
		// Empty repository
		return make(map[string][]string), hasGoFiles, nil
	}

	// Check go.mod files to see if the base module is a dependency
	// This returns the list of module directories that have the dependency
	moduleDirs, err := repoHasModuleDependency(ctx, client, repo, githubToken, baseModule, cache, tree)
	if err != nil && *verbose {
		log.Printf("Failed to check go.mod for %s: %v (will scan all files anyway)", repo.FullName, err)
		// Continue anyway - repo might not have go.mod or it might be an error
		moduleDirs = nil // Scan all files
	} else if len(moduleDirs) == 0 && moduleDirs != nil {
		// We found go.mod files but none have the dependency
		if *verbose {
			log.Printf("Skipping %s - no dependency on %s in any go.mod", repo.FullName, baseModule)
		}
		return make(map[string][]string), hasGoFiles, nil
	}

	// Get all Go files from the tree, filtered by module directories
	goFiles := extractGoFilesFromTree(tree, moduleDirs)

	if *verbose && moduleDirs != nil {
		log.Printf("Analyzing %d Go files in %d module(s) for %s", len(goFiles), len(moduleDirs), repo.FullName)
	}

	imports := make(map[string][]string)

	// Analyze each Go file
	for _, filePath := range goFiles {
		if err := analyzeGoFile(ctx, client, repo, githubToken, baseModule, filePath, imports, cache); err != nil && *verbose {
			log.Printf("Failed to analyze Go file %s/%s: %v", repo.FullName, filePath, err)
		}
	}

	return imports, hasGoFiles, nil
}

// getRepoTree uses the Git Trees API to get the full file tree of the repository
// This is much more efficient than recursively fetching directory contents
func getRepoTree(ctx context.Context, client *http.Client, repo GitHubRepo, githubToken string, cache *APICache) (*GitTree, error) {
	// Use the Git Trees API with recursive=1 to get all files in one request
	url := fmt.Sprintf("https://api.github.com/repos/%s/git/trees/%s?recursive=1", repo.FullName, repo.DefaultBranch)

	var tree GitTree
	if err := httpGetJSONWithCache(ctx, client, url, githubToken, &tree, cache, 1*time.Hour); err != nil {
		if err.Error() == "not found" {
			return nil, nil // Empty repository or branch doesn't exist
		}
		return nil, err
	}

	if tree.Truncated {
		log.Printf("Warning: tree for %s was truncated, some files may be missing", repo.FullName)
	}

	return &tree, nil
}

// extractGoFilesFromTree extracts all Go file paths from a Git Tree
// If moduleDirs is not nil, only includes files that are within those module directories
func extractGoFilesFromTree(tree *GitTree, moduleDirs []string) []string {
	var goFiles []string
	for _, item := range tree.Tree {
		// Only process files (blobs), not directories (trees)
		if item.Type != "blob" {
			continue
		}

		// Check if it's a Go file
		if !strings.HasSuffix(item.Path, ".go") {
			continue
		}

		// Skip files in directories we typically want to ignore
		if shouldSkipPath(item.Path) {
			continue
		}

		// If moduleDirs is specified, only include files within those directories
		if moduleDirs != nil {
			if !isInModuleDirs(item.Path, moduleDirs) {
				continue
			}
		}

		goFiles = append(goFiles, item.Path)
	}
	return goFiles
}

// isInModuleDirs checks if a file path is within any of the specified module directories
func isInModuleDirs(filePath string, moduleDirs []string) bool {
	for _, moduleDir := range moduleDirs {
		if moduleDir == "" {
			// Root directory - all files are included
			return true
		}
		// Check if file is in this module directory
		if strings.HasPrefix(filePath, moduleDir+"/") {
			return true
		}
	}
	return false
}

// shouldSkipPath checks if a file path should be skipped based on directory patterns
func shouldSkipPath(path string) bool {
	skipDirs := []string{
		".git/", ".github/", "vendor/", "node_modules/",
		".vscode/", ".idea/", "dist/", "build/", "target/",
		"docs/", "documentation/", "examples/", "test-data/",
		"ui/", "manifests/", "deploy/", "scripts/",
	}

	for _, skip := range skipDirs {
		if strings.HasPrefix(path, skip) || strings.Contains(path, "/"+skip) {
			return true
		}
	}
	return false
}

func analyzeGoFile(ctx context.Context, client *http.Client, repo GitHubRepo, githubToken, baseModule, filePath string, imports map[string][]string, cache *APICache) error {
	// Get file contents using cache
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo.FullName, filePath)

	var fileData map[string]any
	if err := httpGetJSONWithCache(ctx, client, url, githubToken, &fileData, cache, 6*time.Hour); err != nil {
		return err
	}

	content, ok := fileData["content"].(string)
	if !ok {
		return fmt.Errorf("no content in file response")
	}

	// Decode base64 content
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(content, "\n", ""))
	if err != nil {
		return err
	}

	// Parse Go file and extract imports
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, decoded, parser.ImportsOnly)
	if err != nil {
		return err
	}

	for _, importSpec := range node.Imports {
		importPath := strings.Trim(importSpec.Path.Value, "\"")

		if strings.HasPrefix(importPath, baseModule+"/") || importPath == baseModule {
			// This is an import from our base module
			imports[importPath] = append(imports[importPath], filePath)
		}
	}

	return nil
}

func outputResult(result AnalysisResult) error {
	var writer io.Writer = os.Stdout

	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer file.Close()
		writer = file

		if *verbose {
			log.Printf("Writing output to file: %s", *outputFile)
		}
	}

	switch *output {
	case "json":
		return outputJSON(result, writer)
	case "csv":
		return outputCSV(result, writer)
	case "text":
		return outputText(result, writer)
	default:
		return fmt.Errorf("unsupported output format: %s", *output)
	}
}

func outputJSON(result AnalysisResult, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func outputCSV(result AnalysisResult, writer io.Writer) error {
	fmt.Fprintln(writer, "Package,Repository,File,LastPush,LastUpdate")

	// Sort packages for consistent output
	var packages []string
	for pkg := range result.PackageUsage {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)

	for _, pkg := range packages {
		for _, importInfo := range result.PackageUsage[pkg] {
			for _, file := range importInfo.Files {
				fmt.Fprintf(writer, "%s,%s,%s,%s,%s\n", pkg, importInfo.Repository, file,
					importInfo.LastPush.Format("2006-01-02"), importInfo.LastUpdate.Format("2006-01-02"))
			}
		}
	}
	return nil
}

func outputText(result AnalysisResult, writer io.Writer) error {
	fmt.Fprintf(writer, "Package Import Analysis for %s\n", result.BaseModule)
	fmt.Fprintf(writer, "Organization: %s\n", result.Organization)
	fmt.Fprintf(writer, "Analysis Date: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(writer, "Total Repositories: %d\n", result.TotalRepos)
	fmt.Fprintf(writer, "Repositories with Go Files: %d\n", result.ReposWithGoFiles)
	fmt.Fprintf(writer, "Repositories with Imports: %d\n", result.ReposWithImports)
	fmt.Fprintln(writer)

	if len(result.PackageUsage) == 0 {
		fmt.Fprintln(writer, "No imports found from the base module.")
		return nil
	}

	// Sort packages for consistent output
	var packages []string
	for pkg := range result.PackageUsage {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)

	fmt.Fprintf(writer, "Packages and their usage:\n")
	fmt.Fprintf(writer, "========================\n\n")

	for _, pkg := range packages {
		importInfos := result.PackageUsage[pkg]

		// Count unique repositories
		repoMap := make(map[string]ImportInfo)
		totalFiles := 0
		for _, info := range importInfos {
			repoMap[info.Repository] = info
			totalFiles += len(info.Files)
		}

		fmt.Fprintf(writer, "Package: %s\n", pkg)
		fmt.Fprintf(writer, "  Used by %d repositories (%d files total)\n", len(repoMap), totalFiles)

		// Sort repositories for consistent output
		var repos []string
		for repo := range repoMap {
			repos = append(repos, repo)
		}
		sort.Strings(repos)

		for _, repo := range repos {
			info := repoMap[repo]
			fmt.Fprintf(writer, "  - %s (last push: %s, last update: %s)\n", repo,
				info.LastPush.Format("2006-01-02"), info.LastUpdate.Format("2006-01-02"))
			// Find files for this repo
			for _, file := range info.Files {
				fmt.Fprintf(writer, "    * %s\n", file)
			}
		}
		fmt.Fprintln(writer)
	}

	return nil
}
