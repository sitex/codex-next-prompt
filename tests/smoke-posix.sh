#!/bin/sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/codex-next-prompt-smoke.XXXXXX")
trap 'rm -rf "$temp_root"' EXIT HUP INT TERM

case "$(uname -s)" in
Linux)
	os=linux
	;;
Darwin)
	os=darwin
	;;
*)
	printf '%s\n' "smoke: unsupported host operating system: $(uname -s)" >&2
	exit 1
	;;
esac

case "$(uname -m)" in
x86_64 | amd64)
	arch=amd64
	;;
arm64 | aarch64)
	arch=arm64
	;;
*)
	printf '%s\n' "smoke: unsupported host architecture: $(uname -m)" >&2
	exit 1
	;;
esac

plugin_root="$temp_root/plugin"
binary_dir="$plugin_root/bin/$os-$arch"
mkdir -p "$binary_dir" "$plugin_root/hooks" "$temp_root/work/subdirectory"
cp "$repo_root/hooks/run" "$plugin_root/hooks/run"
chmod +x "$plugin_root/hooks/run"

GOTOOLCHAIN=auto go build -o "$binary_dir/codex-next-prompt" "$repo_root/cmd/codex-next-prompt"

session_output=$(
	cd "$temp_root/work/subdirectory"
	printf '%s\n' '{"hook_event_name":"SessionStart","source":"startup"}' | "$plugin_root/hooks/run" session-start
)
case "$session_output" in
*'"hookEventName":"SessionStart"'*'"additionalContext":"'*'Suggested next prompt:'*)
	;;
*)
	printf '%s\n' "smoke: invalid SessionStart output: $session_output" >&2
	exit 1
	;;
esac

stop_output=$(
	cd "$temp_root/work/subdirectory"
	printf '%s\n' '{"hook_event_name":"Stop","last_assistant_message":"Suggested next prompt:"}' | "$plugin_root/hooks/run" stop
)
case "$stop_output" in
'{"systemMessage":"codex-next-prompt: invalid Suggested next prompt line"}')
	;;
*)
	printf '%s\n' "smoke: invalid Stop output: $stop_output" >&2
	exit 1
	;;
esac

mkdir -p "$temp_root/unsupported-bin"
cat >"$temp_root/unsupported-bin/uname" <<'EOF'
#!/bin/sh
printf '%s\n' Plan9
EOF
chmod +x "$temp_root/unsupported-bin/uname"
unsupported_stdout="$temp_root/unsupported.stdout"
unsupported_stderr="$temp_root/unsupported.stderr"
if ! PATH="$temp_root/unsupported-bin:$PATH" "$plugin_root/hooks/run" stop >"$unsupported_stdout" 2>"$unsupported_stderr"; then
	printf '%s\n' 'smoke: unsupported platform must fail open with exit 0' >&2
	exit 1
fi
if test -s "$unsupported_stdout"; then
	printf '%s\n' 'smoke: unsupported platform wrote to stdout' >&2
	exit 1
fi
unsupported_message=$(cat "$unsupported_stderr")
if test "$unsupported_message" != 'codex-next-prompt: unsupported operating system: Plan9'; then
	printf '%s\n' "smoke: invalid unsupported-platform diagnostic: $unsupported_message" >&2
	exit 1
fi

printf '%s\n' 'smoke: POSIX launcher passed'
