#!/usr/bin/env bash
set -e

DEPLOY_DIR="$(dirname "$(readlink -f $0)")"
ROOT_DIR="$(realpath "$DEPLOY_DIR/..")"

echo "Found project root at $ROOT_DIR"

VENDOR_DIR="$(mktemp -d)"
echo "Vendoring go modules into $VENDOR_DIR"
(cd "$ROOT_DIR" && go mod vendor -v -o "$VENDOR_DIR")

IMAGE_TAG="sc-bos:latest"
echo "\nBuilding container image $IMAGE_TAG"
podman build -f "$DEPLOY_DIR/Dockerfile" \
  -t "$IMAGE_TAG" \
  --secret id=npmrc,src=$HOME/.npmrc \
  --ssh default \
  -v "$ROOT_DIR/.git:/.git" \
  -v "$VENDOR_DIR:/work/vendor" \
  "$ROOT_DIR"

rm -r "$VENDOR_DIR"