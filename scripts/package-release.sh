#!/bin/sh
set -eu

usage() {
	printf '%s\n' 'usage: scripts/package-release.sh VERSION GOOS GOARCH' >&2
}

if test "$#" -ne 3; then
	usage
	exit 2
fi

version=$1
goos=$2
goarch=$3

case "$version" in
'' | *[!0-9A-Za-z._+-]* | .* | *..* | *--* | *++*)
	printf '%s\n' "package-release: invalid version: $version" >&2
	exit 2
	;;
esac

case "$goos" in
linux | darwin | windows)
	;;
*)
	printf '%s\n' "package-release: unsupported GOOS: $goos" >&2
	exit 2
	;;
esac

case "$goarch" in
amd64 | arm64)
	;;
*)
	printf '%s\n' "package-release: unsupported GOARCH: $goarch" >&2
	exit 2
	;;
esac

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
dist_dir=${DIST_DIR:-"$repo_root/dist"}
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/codex-next-prompt-package.XXXXXX")
trap 'rm -rf "$temp_root"' EXIT HUP INT TERM

release_name="codex-next-prompt-$version"
stage_root="$temp_root/$release_name"
target_dir="$stage_root/bin/$goos-$goarch"
mkdir -p "$stage_root/.codex-plugin" "$stage_root/hooks" "$target_dir" "$dist_dir"

sed 's/"version": "[^"]*"/"version": "'"$version"'"/' \
	"$repo_root/.codex-plugin/plugin.json" >"$stage_root/.codex-plugin/plugin.json"
cp "$repo_root/hooks/hooks.json" "$stage_root/hooks/hooks.json"
cp "$repo_root/hooks/run" "$stage_root/hooks/run"
cp "$repo_root/hooks/run.cmd" "$stage_root/hooks/run.cmd"

binary_name=codex-next-prompt
if test "$goos" = windows; then
	binary_name=codex-next-prompt.exe
fi

(
	cd "$repo_root"
	GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
		go build -trimpath -o "$target_dir/$binary_name" ./cmd/codex-next-prompt
)

chmod +x "$stage_root/hooks/run"
if test "$goos" != windows; then
	chmod +x "$target_dir/$binary_name"
fi

expected_files=$(cat <<EOF
$release_name/.codex-plugin/plugin.json
$release_name/bin/$goos-$goarch/$binary_name
$release_name/hooks/hooks.json
$release_name/hooks/run
$release_name/hooks/run.cmd
EOF
)
actual_files=$(cd "$temp_root" && find "$release_name" -type f -print | LC_ALL=C sort)
if test "$actual_files" != "$expected_files"; then
	printf '%s\n' 'package-release: staged payload does not match release contract' >&2
	printf '%s\n' "$actual_files" >&2
	exit 1
fi

if test "$goos" = windows; then
	archive_path="$dist_dir/$release_name-$goos-$goarch.zip"
	rm -f "$archive_path" "$archive_path.sha256"
	(cd "$temp_root" && zip -q -r "$archive_path" "$release_name")
	archive_files=$(unzip -Z1 "$archive_path" | grep -v '/$' | LC_ALL=C sort)
else
	archive_path="$dist_dir/$release_name-$goos-$goarch.tar.gz"
	rm -f "$archive_path" "$archive_path.sha256"
	(cd "$temp_root" && tar -czf "$archive_path" "$release_name")
	archive_files=$(tar -tzf "$archive_path" | grep -v '/$' | LC_ALL=C sort)
fi

if test "$archive_files" != "$expected_files"; then
	printf '%s\n' 'package-release: archive payload does not match release contract' >&2
	printf '%s\n' "$archive_files" >&2
	exit 1
fi

archive_base=$(basename "$archive_path")
if command -v sha256sum >/dev/null 2>&1; then
	checksum=$(sha256sum "$archive_path" | cut -d ' ' -f 1)
elif command -v shasum >/dev/null 2>&1; then
	checksum=$(shasum -a 256 "$archive_path" | cut -d ' ' -f 1)
else
	printf '%s\n' 'package-release: sha256sum or shasum is required' >&2
	exit 1
fi
printf '%s  %s\n' "$checksum" "$archive_base" >"$archive_path.sha256"

printf '%s\n' "$archive_path"
printf '%s\n' "$archive_path.sha256"
