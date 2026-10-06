#!/bin/sh
# Installs goremi: finds the OS and architecture, downloads the release binary to ~/.local/bin.
# The whole script is one function called on the last line: when `curl | sh` is cut off, sh gets an unfinished function and runs nothing.
set -eu

release_url="https://github.com/nhhthong/goremi/releases/latest/download"

# START: dependency hint

# hint prints how to install a missing tool: brew on macOS, apt and the release binary on Linux (Debian or Ubuntu first).
hint() {
	case "$os:$1" in
	darwin:*) printf ' (install it with: brew install %s)' "$1" ;;
	linux:mpv) printf ' (Debian or Ubuntu: sudo apt install mpv)' ;;
	linux:yt-dlp) printf ' (install the release binary: curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o ~/.local/bin/yt-dlp && chmod a+rx ~/.local/bin/yt-dlp)' ;;
	esac
}

# END: dependency hint

main() {
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
			printf '%s not found%s\n' "$tool" "$(hint "$tool")"
		fi
	done
	# END: check dependencies

	# START: check PATH
	case ":$PATH:" in
	*":$dir:"*) ;;
	*) printf '~/.local/bin is not in your PATH; add it to your shell profile: export PATH="$HOME/.local/bin:$PATH"\n' ;;
	esac
	# END: check PATH
}

main "$@"
