# Package Import Analyzer

This tool analyzes GitHub repositories in an organization to find which packages from the `sc-bos` repository are imported by other repositories. It helps identify package dependencies and understand which packages are being used externally, which is useful for determining breaking changes and API stability requirements.

## Features

- Scans all repositories in a GitHub organization for Go files (not just those marked as "Go" language)
- Identifies imports from the specified base module (defaults to `github.com/vanti-dev/sc-bos`)
- Supports multiple output formats: JSON, CSV, and human-readable text
- Supports writing output to files while keeping logs on console
- **Persistent API response caching** with organized file structure to avoid rate limits and improve performance across runs
- **Efficient Git Trees API** usage to fetch all repository files in a single API call instead of recursive directory traversal
- Provides detailed information about which files use each package
- Includes repository liveliness information (last push and update dates)
- Handles repository pagination for organizations with many repositories
- Smart directory filtering to skip common non-source directories

## Prerequisites

- Go 1.25 or later
- A GitHub Personal Access Token with repository read permissions

## Installation

```bash
cd cmd/tools/package-import-analyzer
go build -o package-import-analyzer
```

## Usage

### Basic Usage

```bash
# Analyze the "vanti-dev" organization
./package-import-analyzer -org vanti-dev -token YOUR_GITHUB_TOKEN

# Or set the token as an environment variable
export GITHUB_TOKEN=your_github_token
./package-import-analyzer -org vanti-dev
```

### Command Line Options

- `-org`: GitHub organization name to scan (required)
- `-token`: GitHub personal access token (or set `GITHUB_TOKEN` environment variable)
- `-base-module`: Base module to look for imports (default: `github.com/vanti-dev/sc-bos`)
- `-output`: Output format: `json`, `csv`, or `text` (default: `json`)
- `-output-file`: Output file path (if not specified, outputs to stdout)
- `-cache-dir`: Cache directory for API responses (default: `~/.package-import-analyzer-cache`)
- `-verbose`: Enable verbose logging (repository analysis progress and summary statistics)
- `-vv`: Enable very verbose logging (includes individual cache hits/misses and API calls)

### Output Formats

#### JSON Output (default)
```bash
./package-import-analyzer -org vanti-dev -output json
```

Returns structured JSON with detailed analysis results including timestamps, repository counts, and detailed package usage information.

#### CSV Output
```bash
./package-import-analyzer -org vanti-dev -output csv
```

Returns CSV format with columns: Package, Repository, File, LastPush, LastUpdate

#### Text Output
```bash
./package-import-analyzer -org vanti-dev -output text
```

Returns human-readable text format with summary statistics and detailed package usage breakdown.

### Examples

#### Analyze with verbose output
```bash
./package-import-analyzer -org vanti-dev -verbose -output text
```

#### Analyze different base module
```bash
./package-import-analyzer -org your-org -base-module github.com/your-org/your-repo -output csv
```

#### Save results to file
```bash
./package-import-analyzer -org vanti-dev -output json -output-file analysis-results.json
./package-import-analyzer -org vanti-dev -output csv -output-file analysis-results.csv
./package-import-analyzer -org vanti-dev -output text -output-file analysis-results.txt

# Logs will still appear on console when using -verbose
./package-import-analyzer -org vanti-dev -verbose -output json -output-file results.json
```

#### Verbose logging levels
```bash
# No verbose logging (only errors and final summary)
./package-import-analyzer -org vanti-dev

# Verbose logging (-verbose): Repository analysis progress, summaries, and cache statistics
./package-import-analyzer -org vanti-dev -verbose

# Very verbose logging (-vv): All of the above plus individual cache hits/misses and API calls
./package-import-analyzer -org vanti-dev -vv
```

#### Using custom cache directory
```bash
# Use a custom cache directory
./package-import-analyzer -org vanti-dev -cache-dir ./my-cache

# Clear cache by removing the directory
rm -rf ~/.package-import-analyzer-cache
./package-import-analyzer -org vanti-dev -verbose
```

### Caching

The tool implements persistent file-based caching to dramatically improve performance and reduce API usage:

- **Repository Lists**: Cached for 30 minutes (repositories don't change frequently)
- **Language Information**: Cached for 24 hours (language statistics are stable)
- **Git Trees (File Lists)**: Cached for 1 hour (repository structure changes occasionally)
- **File Contents**: Cached for 6 hours (import statements change less frequently)

#### Cache Structure:
- Cache files are organized into subdirectories based on the first two characters of their hash (e.g., `.cache/0e/`, `.cache/a3/`)
- This organization improves file system performance when handling many cached files
- Up to 256 subdirectories (00-ff) distribute cache files efficiently

#### Cache Benefits:
- **Subsequent runs**: Much faster execution when analyzing the same organization
- **Rate limit protection**: Significantly reduced GitHub API calls
- **Interrupted runs**: Restart from where you left off without re-fetching data
- **Development workflow**: Quick iterations when testing different base modules

#### Cache Management:
- Cache files are stored as JSON with expiration timestamps
- Expired cache entries are automatically cleaned up
- Use `-verbose` to see cache statistics summary and `-vv` for individual cache hits/misses
- Clear cache manually by removing the cache directory when needed

### Verbose Logging Levels

The tool supports multiple verbosity levels to control log output:

- **Default (no flags)**: Only errors and final results
- **`-verbose`**: Shows:
  - Repository analysis progress
  - Module dependency checks
  - Package import discoveries
  - Cache statistics summary (total API calls, cache hits, hit rate)
  - Skipped repositories and reasons
- **`-vv` (very verbose)**: Shows all of the above plus:
  - Individual cache hits for each API call
  - Individual cache misses and subsequent API calls
  - Running totals of API calls
  - Cache write failures

Use `-verbose` for normal operation monitoring, and `-vv` for debugging cache behavior or API usage patterns.

### Output and Logging

When using the `-output-file` flag:
- **Analysis results** are written to the specified file
- **Log messages** (including verbose output) continue to appear on the console
- **Cache information** is logged when using `-verbose` flag
- **Error messages** are displayed on the console, and no output file is created on error
- This separation allows you to save clean results while still monitoring the analysis progress

## Output Example

### Text Format
```
Package Import Analysis for github.com/vanti-dev/sc-bos
Organization: vanti-dev
Analysis Date: 2025-11-19 10:30:45
Total Repositories: 25
Repositories with Go Files: 18
Repositories with Imports: 3

Packages and their usage:
========================

Package: github.com/vanti-dev/sc-bos/pkg/gen
  Used by 2 repositories (5 files total)
  - vanti-dev/control-system (last push: 2025-11-15, last update: 2025-11-15)
    * cmd/main.go
    * internal/client/bos.go
  - vanti-dev/dashboard (last push: 2025-11-10, last update: 2025-11-12)
    * server/api/devices.go
    * server/api/alerts.go
    * ui/src/services/api.go

Package: github.com/vanti-dev/sc-bos/pkg/node
  Used by 1 repositories (2 files total)
  - vanti-dev/edge-gateway (last push: 2025-11-08, last update: 2025-11-08)
    * internal/node/manager.go
    * cmd/gateway/main.go
```

## Implementation Details

The tool:

1. **Repository Discovery**: Uses the GitHub API to fetch all repositories in the specified organization
2. **Efficient Go Detection**: Uses the GitHub Languages API to quickly determine if a repository contains Go code before scanning files
3. **File Tree Retrieval**: Uses the Git Trees API with `recursive=1` to fetch all files in a repository with a single API call
4. **Smart Filtering**: Filters the file tree for `.go` files while skipping common non-source directories
5. **Go File Parsing**: Downloads and parses Go source files to extract import statements
6. **Import Matching**: Identifies imports that start with the specified base module path
7. **Liveliness Tracking**: Captures last push and update timestamps for dependency assessment
8. **Result Aggregation**: Collects and organizes results by package and repository

## Limitations

- Only analyzes public repositories (or those accessible with the provided token)
- Large organizations with many repositories may take some time to analyze on first run (subsequent runs use cache)
- GitHub API rate limits apply (typically 5000 requests/hour for authenticated requests)
- Only analyzes the default branch of each repository
- Very large repositories may have truncated file trees (GitHub limits recursive tree results to 7MB or 100,000 entries)
- Skips common non-source directories (vendor, node_modules, .git, docs, ui, etc.)

## Performance Optimizations

The tool includes several optimizations to minimize API usage and improve speed:

- **Git Trees API**: Uses the GitHub Git Trees API with `recursive=1` to fetch all repository files in a single API call instead of making recursive calls for each directory
- **Languages API First**: Uses the GitHub Languages API to quickly identify repositories with Go code before scanning files
- **Skip Non-Go Repositories**: Completely skips repositories that don't contain Go files, saving significant API calls
- **Directory Filtering**: Skips common directories that typically don't contain relevant source code (vendor, node_modules, docs, ui, etc.)
- **Efficient Pagination**: Handles repository pagination efficiently for large organizations
- **Persistent Caching**: File-based caching system with intelligent TTL policies:
  - Repository data cached for extended periods (repositories change infrequently)
  - File trees cached to avoid re-scanning repository structure
  - File contents cached to avoid re-downloading during analysis
  - Smart cache expiration prevents stale data while maximizing reuse
- **Organized Cache Storage**: Cache files distributed across subdirectories (based on first two characters of hash) for improved file system performance
- **Cache Hit Logging**: Verbose mode shows cache hit rate and API call statistics

## Error Handling

The tool gracefully handles:
- Repositories without Go files
- Empty repositories
- Network timeouts
- GitHub API errors
- Malformed Go files

Errors are logged in verbose mode but don't stop the overall analysis.

## GitHub Token Requirements

The GitHub token needs at least the following scopes:
- `repo` (for private repositories) or `public_repo` (for public repositories only)

To generate a token:
1. Go to GitHub Settings > Developer settings > Personal access tokens
2. Click "Generate new token"
3. Select appropriate scopes
4. Copy the generated token
