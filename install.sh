#!/bin/sh
# Installs goremi: finds the OS and architecture, downloads the release binary to ~/.local/bin.
# The whole script is one function called on the last line: when `curl | sh` is cut off, sh gets an unfinished function and runs nothing.
set -eu

release_url="https://github.com/nhhthong/goremi/releases/latest/download"

# START: dependency offer

# find_pm names the package manager of this system: brew on macOS, then apt, dnf, pacman and brew on Linux; empty when none is installed.
find_pm() {
	if [ "$os" = darwin ]; then
		command -v brew >/dev/null 2>&1 && printf 'brew'
		return 0
	fi
	for candidate in apt dnf pacman brew; do
		if command -v "$candidate" >/dev/null 2>&1; then
			printf '%s' "$candidate"
			return 0
		fi
	done
}

# mpv_command is the install command of mpv for a package manager, without sudo.
mpv_command() {
	case "$1" in
	apt) printf 'apt install -y mpv' ;;
	dnf) printf 'dnf install -y mpv' ;;
	pacman) printf 'pacman -S --noconfirm mpv' ;;
	*) printf 'brew install mpv' ;;
	esac
}

# ytdlp_asset is the yt-dlp release file of this system: a standalone binary, no Python needed.
ytdlp_asset() {
	case "$os:$arch" in
	darwin:*) printf 'yt-dlp_macos' ;;
	linux:arm64) printf 'yt-dlp_linux_aarch64' ;;
	*) printf 'yt-dlp_linux' ;;
	esac
}

# file_hash prints the SHA-256 of a file.
file_hash() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | cut -d' ' -f1
	else
		shasum -a 256 "$1" | cut -d' ' -f1
	fi
}

# install_ytdlp downloads the yt-dlp binary into the install directory and checks it against SHA2-256SUMS of the same release; a mismatch installs nothing.
install_ytdlp() {
	asset=$(ytdlp_asset)
	base="https://github.com/yt-dlp/yt-dlp/releases/latest/download"
	bin=$(mktemp "$dir/.yt-dlp.XXXXXX") || return 1
	sums=$(mktemp "$dir/.yt-dlp-sums.XXXXXX") || { rm -f "$bin"; return 1; }
	if ! curl -fsSL -o "$bin" "$base/$asset" || ! curl -fsSL -o "$sums" "$base/SHA2-256SUMS"; then
		rm -f "$bin" "$sums"
		printf 'goremi install: could not download yt-dlp\n' >&2
		return 1
	fi
	want=$(grep "  $asset\$" "$sums" | cut -d' ' -f1)
	got=$(file_hash "$bin")
	rm -f "$sums"
	if [ -z "$want" ] || [ "$want" != "$got" ]; then
		rm -f "$bin"
		printf 'goremi install: the checksum of yt-dlp does not match SHA2-256SUMS; yt-dlp not installed\n' >&2
		return 1
	fi
	chmod 755 "$bin"
	mv -f "$bin" "$dir/yt-dlp"
}

# offer_dependencies lists the commands that install what is missing, asks once, and runs them on a yes. A failed step is reported and never undoes the goremi install.
offer_dependencies() {
	pm=$(find_pm)
	run_as=""
	if [ -n "$pm" ] && [ "$pm" != brew ] && [ "$(id -u 2>/dev/null || printf 1)" != 0 ]; then
		run_as="sudo "
	fi
	shown=${pm:-apt}
	[ "$os" = darwin ] && shown=${pm:-brew}
	shown_as="$run_as"
	if [ -z "$pm" ]; then # no manager: show the Debian or Ubuntu command (the Homebrew one on macOS) as the example
		shown_as="sudo "
		[ "$os" = darwin ] && shown_as=""
	fi
	names=""
	printf 'These commands install what is missing:\n'
	for tool in $missing; do
		names="${names:+$names and }$tool"
		case "$tool" in
		mpv)
			printf '  %s%s\n' "$shown_as" "$(mpv_command "$shown")"
			[ -z "$pm" ] && printf '  (no known package manager was found: install mpv with yours)\n'
			;;
		yt-dlp) printf '  curl -fsSL -o ~/.local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/%s\n' "$(ytdlp_asset)" ;;
		esac
	done
	if [ "$assume_yes" = 1 ]; then
		answer=y
	else
		printf 'Install %s? [y/N] ' "$names"
		if ! IFS= read -r answer 2>/dev/null <"${GOREMI_TTY:-/dev/tty}"; then
			answer=n
			printf '\nno answer: run the commands above yourself, or run the installer again with -y\n'
		fi
	fi
	case "$answer" in y | Y | yes) ;; *) return 0 ;; esac
	for tool in $missing; do
		case "$tool" in
		mpv)
			if [ -n "$pm" ]; then
				# shellcheck disable=SC2086
				$run_as$(mpv_command "$pm") || printf 'goremi install: could not install mpv\n' >&2
			fi
			;;
		yt-dlp) install_ytdlp || true ;;
		esac
	done
}

# END: dependency offer

main() {
	assume_yes=0
	for arg in "$@"; do
		[ "$arg" = -y ] && assume_yes=1
	done
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
	missing=""
	for tool in mpv yt-dlp; do
		if ! command -v "$tool" >/dev/null 2>&1; then
			printf '%s not found\n' "$tool"
			missing="$missing $tool"
		fi
	done
	[ -n "$missing" ] && offer_dependencies
	# END: check dependencies

	# START: check PATH
	case ":$PATH:" in
	*":$dir:"*) ;;
	*) printf '~/.local/bin is not in your PATH; add it to your shell profile: export PATH="$HOME/.local/bin:$PATH"\n' ;;
	esac
	# END: check PATH
}

main "$@"
