#!/usr/bin/env bash

set -e

PROJECT_ROOT="$(realpath "$(dirname "$0")/..")"
BUILD_ROOT="/tmp/bsp-ew-build"
mkdir -p "$BUILD_ROOT"

echo "Found project in $PROJECT_ROOT"
echo "Building in $BUILD_ROOT"

echo "Archiving Go code..."
# copy the sources
AREA_CONTROLLER_COPY="$BUILD_ROOT/area-controller"
rsync -a "$PROJECT_ROOT/"{cmd,internal,pkg,go.mod,go.sum} "$AREA_CONTROLLER_COPY"

# vendor all dependencies
cd "$AREA_CONTROLLER_COPY"
go mod vendor
# pack into an archive
cd "$BUILD_ROOT"
tar -cf "$BUILD_ROOT/area-controller.tar" area-controller
# remove the copy
rm -rf "$AREA_CONTROLLER_COPY"
echo "Archive done"

echo "Archiving Conductor UI..."
cd "$PROJECT_ROOT/ui/conductor/dist"
tar -cf "$BUILD_ROOT/conductor.tar" .