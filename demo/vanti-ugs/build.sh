#!/usr/bin/env bash

set -euo pipefail

# Builds the sc-bos docker image, which includes the sc-bos binary and ops-ui web app.
# Then extends that image with demo configuration and opinionated defaults.

containerCmd="podman"

if ! command -v $containerCmd &> /dev/null; then
  echo "Error: $containerCmd is not installed."
  exit 1
fi

# The sc-bos build requires secrets that are usually in the users home .npmrc file,
# so make sure it exists.
if [ ! -f "$HOME/.npmrc" ]; then
  echo "Error: $HOME/.npmrc file not found. Please create it with your npm credentials."
  exit 1
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
GIT_VERSION=$(git describe --tags --always --dirty)
BASE_IMAGE="localhost/vanti-dev-sc-bos:demo-ugs-base"
CR_TAG_PREFIX="ghcr.io/vanti-dev/sc-bos"
PLATFORMS="linux/amd64,linux/arm64"

cd $REPO_ROOT
echo "Preparing $containerCmd for building the image..."
$containerCmd manifest rm "$BASE_IMAGE" || true
$containerCmd manifest rm "$CR_TAG_PREFIX/demo-ugs-sc-bos:latest" || true
$containerCmd manifest rm "$CR_TAG_PREFIX/demo-ugs-seed-db:latest" || true

# Build using the shared Dockerfile in the project root directory
echo "Building the sc-bos base image with version: $GIT_VERSION"
$containerCmd build \
  --platform=$PLATFORMS \
  --jobs=2 \
  --build-arg GIT_VERSION=$GIT_VERSION \
  --secret=id=npmrc,src=$HOME/.npmrc \
  --manifest "$BASE_IMAGE" \
  -f $REPO_ROOT/Dockerfile \
  .

# Add config to the sc-bos image
echo "Building the sc-bos demo image with version: $GIT_VERSION"
$containerCmd build \
  --from="$BASE_IMAGE" \
  --pull=never \
  --all-platforms \
  --manifest "$CR_TAG_PREFIX/demo-ugs-sc-bos:latest" \
  -f demo/vanti-ugs/Dockerfile-Ugs .
# Build the db seeder image
echo "Building the sc-bos demo seed-db image"
$containerCmd build \
  --platform=$PLATFORMS \
  --manifest "$CR_TAG_PREFIX/demo-ugs-seed-db:latest" \
  -f demo/vanti-ugs/Dockerfile-SeedDb .




