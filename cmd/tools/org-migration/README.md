# Organization Migration Tool

A command-line tool to migrate code references from the `vanti-dev` organization to the `smart-core-os` organization.

## Overview

The `sc-bos` repository was transferred from `vanti-dev` to `smart-core-os` on GitHub. This tool helps update references
in your codebase to match the new organization.

## How It Works

The tool scans files in your project and applies specific replacements based on the preset you choose:

### What Gets Replaced

| Replacement Type         | Pattern                               | Example                                             |
|--------------------------|---------------------------------------|-----------------------------------------------------|
| **Go imports**           | `github.com/vanti-dev/sc-bos`         | `import "github.com/vanti-dev/sc-bos/pkg/node"`     |
| **npm packages**         | `@vanti-dev/sc-bos`                   | `"@vanti-dev/sc-bos": "^1.0.0"`                     |
| **Docker images**        | `ghcr.io/vanti-dev/sc-bos`            | `image: ghcr.io/vanti-dev/sc-bos:latest`            |
| **GitHub URLs**          | `https://github.com/vanti-dev/sc-bos` | `https://github.com/vanti-dev/sc-bos/blob/main/...` |
| **IDEA run configs**     | `github.com/vanti-dev/sc-bos`         | Package paths in `.run/*.xml` files                 |

### Presets

**Presets control which replacements are applied**, not which files are scanned. This ensures references in
documentation, config files, and code are all updated together when you migrate a specific component.

| Preset   | What It Does                                             | Use When                                                        |
|----------|----------------------------------------------------------|-----------------------------------------------------------------|
| `go`     | Updates Go import paths everywhere (code, docs, configs) | Your project imports Go packages from `sc-bos`                  |
| `js`     | Updates npm package references everywhere                | Your project depends on npm packages from `@vanti-dev/sc-bos`   |
| `docker` | Updates Docker image references everywhere               | Your project uses Docker images from `ghcr.io/vanti-dev/sc-bos` |
| `docs`   | Updates GitHub URLs everywhere                           | You want to update documentation links only                     |
| `all`    | Applies all replacements (default)                       | Your project uses multiple components from `sc-bos`             |

**By default, all file types are scanned** (`.go`, `.mod`, `.js`, `.ts`, `.md`, `.yml`, `.yaml`, `.json`, `.xml`, etc.), 
including **IntelliJ IDEA run configuration files** in `.run/` directories. Use `--types` to limit
which files are processed.

## Installation

### Install from GitHub

```bash
# Install the latest version directly from GitHub
go install github.com/smart-core-os/sc-bos/cmd/tools/org-migration@latest

# Then use it anywhere
org-migration --dry-run --preset go --path /path/to/your/project
```

### Run from source

If you have the repository cloned locally:

```bash
go run ./cmd/tools/org-migration --dry-run --preset go --path /path/to/your/project
```

## Quick Start

```bash
# Preview what would change (always do this first!)
org-migration --dry-run --preset go --path /path/to/your/project

# Apply the changes
org-migration --preset go --path /path/to/your/project
```

## Usage Examples

```bash
# Migrate Go imports if your project imports sc-bos packages
# (updates imports in .go files AND in READMEs, docs, etc.)
org-migration --preset go --path ~/my-project

# Migrate npm packages if your project uses @vanti-dev/sc-bos
org-migration --preset js --path ~/my-project

# Migrate Docker images if your project uses sc-bos containers
org-migration --preset docker --path ~/my-project

# Update only GitHub URLs in documentation
org-migration --preset docs --path ~/my-project

# Scan only specific file types
org-migration --types go,mod --preset go --path ~/my-project

# Verbose output to see exactly what changes
org-migration --dry-run --verbose --preset all --path .
```

**Note:** If you're running from source instead of using `go install`, replace `org-migration` with
`go run ./cmd/tools/org-migration` in the examples above.

## Command-Line Options

```
--dry-run       Preview changes without modifying files (recommended first step)
--verbose       Show detailed output including line-by-line changes
--path PATH     Directory to scan (default: current directory)
--preset NAME   Which replacements to apply: all, go, js, docker, docs (default: all)
--types EXTS    Comma-separated file extensions to scan (default: all supported types)
```

## Safety Features

- **Always use `--dry-run` first** to preview changes
- Preserves file permissions and line endings (CRLF vs LF)
- Skips most hidden directories (`.git`, `.github`) but **processes `.run/` for IDEA run configurations**
- Skips build artifacts (`node_modules`, `vendor`, `dist`)
- Won't modify files in the `org-migration` tool directory itself
- **Automatically skips generated code** (files with `// Code generated ... DO NOT EDIT.`)
  - If generated files contain references, you'll get a warning to regenerate them
  - Look for `//go:generate` directives or run `go generate ./...` after updating source files
- **Automatically renames files** with `vanti-dev` in their names (e.g., IDEA run configuration files)

## Post-Migration Steps

### After `--preset go`

1. Run `go mod tidy` to update dependencies
2. **If you see warnings about generated files**: Run `go generate ./...` to regenerate code
3. Build and test your Go code
4. Verify imports resolve correctly

### After `--preset js`

1. Run `npm install` or `yarn install` to fetch updated packages
2. Clear build caches if needed
3. Test your frontend builds

### After `--preset docker`

1. Update CI/CD pipelines that reference image names
2. Test Docker builds
3. Update deployment configs (Kubernetes, docker-compose, etc.)

### After `--preset all`

1. Test your complete application
2. Search for any remaining hardcoded references: `git grep vanti-dev`
3. Commit changes: `git commit -am "chore: migrate to smart-core-os organization"`

## IntelliJ IDEA Support

The tool automatically updates **IntelliJ IDEA run configuration files**:

- **Run configurations** in `.run/*.xml` - Package paths in Go Application run configs
- **Filename renaming** - Files with `vanti-dev` in their names are automatically renamed

### Examples of What Gets Updated

**In `.run/*.xml` files:**

```xml
<!-- Before -->
<package value="github.com/vanti-dev/sc-bos/cmd/bos"/>

    <!-- After -->
<package value="github.com/smart-core-os/sc-bos/cmd/bos"/>
```

**File renaming:**

- `go build github.com_vanti-dev_sc-bos_cmd_tools_export-alerts.run.xml`
- → `go build github.com_smart-core-os_sc-bos_cmd_tools_export-alerts.run.xml`

No special flags are needed - IDEA run configuration files are updated automatically when using the `go` or `all` presets.

## Migration Strategy

For a staged migration (recommended):

1. **First**: If your project imports Go packages from `sc-bos`
   ```bash
   org-migration --dry-run --preset go --path .
   org-migration --preset go --path .
   go mod tidy
   # If the tool warned about generated files:
   go generate ./...
   # Build and test
   go build ./...
   ```

2. **Then**: If your project uses npm packages from `@vanti-dev/sc-bos`
   ```bash
   org-migration --dry-run --preset js --path .
   org-migration --preset js --path .
   npm install && npm run build
   ```

3. **Then**: If your project uses Docker images from `ghcr.io/vanti-dev/sc-bos`
   ```bash
   org-migration --dry-run --preset docker --path .
   org-migration --preset docker --path .
   docker-compose build
   ```

4. **Finally**: Update documentation links
   ```bash
   org-migration --preset docs --path .
   ```

Or for a complete migration in one step:

```bash
org-migration --dry-run --preset all --path .
org-migration --preset all --path .
```

