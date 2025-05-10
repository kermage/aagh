#!/usr/bin/env sh

set -eu

VERSION="0.5.0"
BIN_NAME="aagh"
BIN_DIR="/usr/local/bin"
BASE_URL="https://github.com/kermage/aagh/releases"
DESCRIPTION="A cross-platform executable for handling Git hooks"


BOLD="$(tput bold 2>/dev/null || printf '')"
UNDERLINE="$(tput smul 2>/dev/null || printf '')"
RESET="$(tput sgr0 2>/dev/null || printf '')"
GREY="$(tput setaf 0 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
GREEN="$(tput setaf 2 2>/dev/null || printf '')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '')"
BLUE="$(tput setaf 4 2>/dev/null || printf '')"
MAGENTA="$(tput setaf 5 2>/dev/null || printf '')"


info() {
	printf '%s\n' "${BOLD}${GREY}➤ ${RESET} $*"
}

warn() {
	printf '%s\n' "${YELLOW}ℹ︎ $*${RESET}"
}

error() {
	printf '%s\n' "${RED}✕ $*${RESET}" >&2
}

success() {
	printf '%s\n' "${GREEN}✓${RESET} $*"
}


has() {
	command -v "$1" 1>/dev/null 2>&1
}

writeable() {
	path="${1}/${BIN_NAME}.tmp"

	if touch "$path" 2>/dev/null; then
		rm "$path"
		return 0
	else
		return 1
	fi
}

maybe_elevate() {
	if writeable "$BIN_DIR"; then
		return
	fi

	if ! has sudo; then
		if [ "$1" = "Windows" ]; then
			error "Please run the shell as administrator."
		else
			error "Please run as root or go install sudo."
		fi

		exit 1
	fi

	warn "Escalated permission required, trying…"

	if ! sudo -v; then
		error "Substitute user not granted, aborting…"
		exit 1
	fi
}

download() {
	file="$1"
	url="$2"

	if has curl; then
		cmd="curl -fsLo $file $url"
	elif has wget; then
		cmd="wget -qO $file $url"
	else
		error "No HTTP download program found ${MAGENTA}(curl or wget)${RESET}"
		return 1
	fi

	$cmd && return 0 || rc=$?

	error "$(printf "%s\n  %s" "Command failed with exit code ($rc):" "${MAGENTA}${cmd}${RESET}")"
	return $rc
}

unpack() {
	archive="$1"
	path="$2"

	case "$archive" in
		*.tar.gz) cmd="tar -xzof $archive -C $path" ;;
		*.zip) cmd="unzip -qqo $archive -d $path" ;;
		*)
			error "Unknown archive format for ${MAGENTA}${archive}${RESET}"
			return 1
			;;
	esac

	$cmd && return 0 || rc=$?

	error "$(printf "%s\n  %s" "Command failed with exit code ($rc):" "${MAGENTA}${cmd}${RESET}")"
	return $rc
}


get_goos() {
	platform="$(uname -s)"

	case "${platform}" in
		Win* | MYSYS* | MINGW* | CYGWIN*) platform="Windows" ;;
	esac

	printf '%s' "$platform"
}

get_goarch() {
	arch="$(uname -m)"

	case "$arch" in
		amd64 | i86pc | x64 | x86-64) arch="x86_64" ;;
		386 | x86 | i386 | i686) arch="i386" ;;
		aarch64) arch="arm64" ;;
	esac

	printf '%s' "$arch"
}

get_target() {
	arch="$1"
	platform="$2"
	target="UNKNOWN"

	case "$platform" in
		Darwin | Linux | Windows) target="$platform" ;;
	esac

	case "$arch" in
		x86_64 | i386 | arm64) target="${target}_${arch}" ;;
	esac

	printf '%s' "$target"
}

usage() {
	printf '%s\n' "${BOLD}Usage${RESET}:    ${BLUE}$0${RESET} ${MAGENTA}[options]${RESET}"
	printf '%s\n' "${BOLD}Options${RESET}:"
	printf '%s\n' "  -b, --bin-dir ${GREEN}<path>${RESET}  Set the destination directory"
	printf '%s\n' "  -h, --help            Show this help message"

	if [ "${1:-0}" != "1" ]; then
		printf '\n'
	fi
}


PLATFORM="$(get_goos)"
ARCH="$(get_goarch)"
CURRENT=""
CURRENT_BIN_DIR=""

if has $BIN_NAME; then
	BIN_DIR=$(dirname $(command -v $BIN_NAME))
	CURRENT="$($BIN_NAME --version 2>/dev/null | cut -d ' ' -f3)"
	CURRENT_BIN_DIR="$BIN_DIR"
fi

while [ "$#" -gt 0 ]; do
	case "$1" in
	-h | --help) usage 1 && exit 0 ;;
	-b | --bin-dir) BIN_DIR="$2" && shift 2 ;;
	-b=* | --bin-dir=*) BIN_DIR="${1#*=}" && shift 1 ;;
	*) usage && error "Unknown option: $1" && exit 1 ;;
	esac
done

printf "\n  %s\n\n" "${UNDERLINE}${BLUE}${DESCRIPTION}${RESET}"
info "${BOLD}Version${RESET}:      ${GREEN}${VERSION}${RESET}"
info "${BOLD}Destination${RESET}:  ${GREEN}${BIN_DIR}${RESET}"
info "${BOLD}Platform${RESET}:     ${GREEN}${PLATFORM}${RESET}"
info "${BOLD}Arch${RESET}:         ${GREEN}${ARCH}${RESET}"
printf '\n'

if [ ! -d "$BIN_DIR" ]; then
	error "Destination does not appear to be a directory."
	warn "Make sure ${BIN_DIR} exists and is writeable."
	exit 1
fi

if [ "$CURRENT" = "$VERSION" ]; then
	success "Already has latest ${UNDERLINE}${BLUE}${BIN_NAME}${RESET}!"
	exit 0
fi

TARGET="$(get_target "${ARCH}" "${PLATFORM}")"

if [ "$TARGET" = "UNKNOWN" ]; then
	error "Current machine is not supported."
	exit 1
fi

EXT="tar.gz"

if [ "$PLATFORM" = "Windows" ]; then
	EXT="zip"
fi

FILE="${BIN_NAME}_${VERSION}_${TARGET}.${EXT}"
URL="${BASE_URL}/download/v${VERSION}/${FILE}"

if [ -n "$CURRENT" ] && [ "$CURRENT_BIN_DIR" = "$BIN_DIR" ]; then
	warn "Updating current ${MAGENTA}v${CURRENT}${YELLOW}, please wait…"
else
	warn "Installation in progress, please wait…"
fi

maybe_elevate "$PLATFORM"
download "$FILE" "$URL"
unpack "$FILE" "$BIN_DIR"
rm -f "$FILE"
success "Latest ${UNDERLINE}${BLUE}${BIN_NAME}${RESET} is now ready!"
