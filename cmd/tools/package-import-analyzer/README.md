# Package Import Analyzer

A tool to analyze which packages from the `sc-bos` repository are imported by other repositories in the organization.

## Overview

The tool consists of two subcommands:

1. **`scan`** - Scans a GitHub organization and outputs a JSON file with import information
2. **`analyze`** - Analyzes the scan results and provides different types of analysis

## Usage

### Scan Command

Scan a GitHub organization for package imports:

```bash
package-import-analyzer scan \
  -org vanti-dev \
  -token $GITHUB_TOKEN \
  -output scan-results.json
```

#### Scan Flags

- `-org` (required) - GitHub organization name to scan
- `-output` (required) - Output file path for scan results
- `-token` - GitHub personal access token (or set `GITHUB_TOKEN` env var)
- `-base-module` - Base module to look for imports (default: `github.com/vanti-dev/sc-bos`)
- `-cache-dir` - Cache directory for API responses (default: `~/.package-import-analyzer-cache`)

#### Scan Output

The scan command produces a JSON file containing:
- Organization and base module information
- Total repository statistics
- Package usage details including:
  - Package import paths
  - Importing repositories
  - Files containing the imports
  - Last push/update timestamps

### Analyze Command

Analyze scan results with various analysis types:

```bash
# Run all analyses
package-import-analyzer analyze -input scan-results.json

# Analyze package types only
package-import-analyzer analyze -input scan-results.json -type package-types

# Analyze repository dependencies only
package-import-analyzer analyze -input scan-results.json -type repo-dependencies
```

#### Analyze Flags

- `-input` (required) - Input file path with scan results
- `-type` - Analysis type: `all`, `package-types`, `repo-dependencies` (default: `all`)
- `-output` - Output format: `text`, `json`, `csv` (default: `text`)
- `-output-file` - Output file path (if not specified, outputs to stdout)

#### Analysis Types

**package-types**: Categorizes imports by **where they are used** in dependent repositories:
- Import location types (tools, drivers, internal, commands, etc.)
- Number of import occurrences per location type
- Number of dependent repositories per type
- List of packages imported at each location type

**Note**: This analysis classifies imports based on the file path in the **dependent repository** where the import occurs, not based on what type of package is being imported from sc-bos. For example, if a tool in another repository imports a driver package from sc-bos, it's counted as a "tools" import because it's located in the `cmd/tools/` directory of the dependent repo.

**repo-dependencies**: Shows which repositories depend on sc-bos packages:
- List of repositories sorted by number of imported packages
- Package count for each repository

**all**: Runs both analyses above plus a basic summary

## Examples

### Complete Workflow

```bash
# 1. Scan the organization
package-import-analyzer scan \
  -org vanti-dev \
  -token $GITHUB_TOKEN \
  -output vanti-dev-imports.json

# 2. Analyze package types
package-import-analyzer analyze \
  -input vanti-dev-imports.json \
  -type package-types

# 3. Export repository dependencies to CSV
package-import-analyzer analyze \
  -input vanti-dev-imports.json \
  -type repo-dependencies \
  -output csv \
  -output-file repo-deps.csv
```

### Output Formats

#### Text Output (default)
Human-readable formatted output with tables and summaries.

#### JSON Output
Structured data suitable for further processing:

```bash
package-import-analyzer analyze \
  -input scan-results.json \
  -type package-types \
  -output json \
  -output-file analysis.json
```

#### CSV Output
Tabular data for spreadsheet analysis:

```bash
package-import-analyzer analyze \
  -input scan-results.json \
  -type repo-dependencies \
  -output csv \
  -output-file repos.csv
```

## Import Location Classification

For package-types analysis, the tool classifies imports based on the file path in the **dependent repository** where the import occurs:

- **drivers** - Files under `internal/driver/` or `pkg/driver/`
- **tools** - Files under `cmd/tools/`
- **automations** - Files under `internal/auto/`, `pkg/auto/`, or similar
- **commands** - Files under `cmd/` (excluding tools)
- **traits** - Files under `pkg/gentrait/`, `pkg/trait/`, or `pkg/traits/`
- **app** - Files under `pkg/app/`
- **node** - Files under `pkg/node/`
- **generated** - Files under `pkg/gen/`
- **proto** - Files under `proto/`
- **tests** - Files under `test/` or `tests/`
- **examples** - Files under `example/` or `examples/`
- **other-internal** - Other files under `internal/`
- **other-pkg** - Other files under `pkg/`
- **other** - Files at the root level or other locations

### Analysis Output

The package-types analysis shows:
- **Import Count** - Number of unique file locations where imports occur
- **Percentage** - What percentage of total imports this type represents (e.g., "43.4% of imports are for drivers")
- **Dependent Repositories** - How many repositories have this type of import
- **Unique Packages Imported** - How many different sc-bos packages are imported at this location type

## Caching

The scan command caches GitHub API responses to reduce API calls and improve performance:

- Cache directory: `~/.package-import-analyzer-cache` (configurable)
- Cache TTL varies by endpoint (30 minutes to 24 hours)
- Cache survives between runs to speed up subsequent scans
- Expired cache entries are automatically cleaned up

## Global Flags

These flags work with any subcommand:

- `-verbose` - Enable verbose output
- `-vv` - Enable very verbose output (includes cache hits/misses and API calls)

