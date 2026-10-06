#!/bin/sh
# Builds the release binaries into dist/: goremi_<os>_<arch>, with .exe on Windows.
set -eu
cd "$(dirname "$0")/.."
rm -rf dist
mkdir dist
for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do
	os=${target%/*}
	arch=${target#*/}
	ext=""
	[ "$os" = windows ] && ext=".exe"
	GOOS=$os GOARCH=$arch go build -o "dist/goremi_${os}_${arch}${ext}" ./cmd/goremi
done
