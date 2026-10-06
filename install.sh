#!/bin/sh
# Installs goremi: finds the OS and architecture, downloads the release binary to ~/.local/bin.
set -eu

release_url="https://github.com/nhhthong/goremi/releases/latest/download"

# START: detect system
raw_os=$(uname -s)
raw_arch=$(uname -m)
case "$raw_os" in
Linux) os=linux ;;
Darwin) os=darwin ;;
*) os="" ;;
esac
case "$raw_arch" in
x86_64) arch=amd64 ;;
aarch64 | arm64) arch=arm64 ;;
*) arch="" ;;
esac
if [ -z "$os" ] || [ -z "$arch" ]; then
	printf 'goremi install: unsupported system: %s %s\n' "$raw_os" "$raw_arch" >&2
	exit 1
fi
# END: detect system

# START: download and install
dir="$HOME/.local/bin"
mkdir -p "$dir"
tmp=$(mktemp "$dir/.goremi.XXXXXX")
trap 'rm -f "$tmp"' EXIT
if ! curl -fsSL -o "$tmp" "$release_url/goremi_${os}_${arch}"; then
	printf 'goremi install: download failed: %s/goremi_%s_%s\n' "$release_url" "$os" "$arch" >&2
	exit 1
fi
chmod 755 "$tmp"
mv -f "$tmp" "$dir/goremi"
# END: download and install

# START: check dependencies
for tool in mpv yt-dlp; do
	if ! command -v "$tool" >/dev/null 2>&1; then
		hint=""
		[ "$os" = darwin ] && hint=" (install it with: brew install $tool)"
		printf '%s not found%s\n' "$tool" "$hint"
	fi
done
# END: check dependencies
