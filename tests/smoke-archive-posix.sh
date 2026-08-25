#!/bin/sh
set -eu

if test "$#" -ne 1; then
	printf '%s\n' 'usage: tests/smoke-archive-posix.sh ARCHIVE' >&2
	exit 2
fi

archive_path=$1
temp_root=$(mktemp -d "${TMPDIR:-/tmp}/codex-next-prompt-archive-smoke.XXXXXX")
trap 'rm -rf "$temp_root"' EXIT HUP INT TERM

tar -xzf "$archive_path" -C "$temp_root"
plugin_root=
for candidate in "$temp_root"/*; do
	if test -d "$candidate"; then
		if test -n "$plugin_root"; then
			printf '%s\n' 'archive smoke: expected one release root' >&2
			exit 1
		fi
		plugin_root=$candidate
	fi
done
if test -z "$plugin_root"; then
	printf '%s\n' 'archive smoke: expected one release root' >&2
	exit 1
fi

session_output=$(printf '%s\n' '{"hook_event_name":"SessionStart","source":"startup"}' | "$plugin_root/hooks/run" session-start)
case "$session_output" in
*'"hookEventName":"SessionStart"'*'next_prompt_rules no_tools=\"true\" no_composer=\"true\"'*'Suggested next prompt:'*) ;;
*)
	printf '%s\n' "archive smoke: invalid SessionStart output: $session_output" >&2
	exit 1
	;;
esac

stop_output=$(printf '%s\n' '{"hook_event_name":"Stop","last_assistant_message":"Suggested next prompt:"}' | "$plugin_root/hooks/run" stop)
if test "$stop_output" != '{"systemMessage":"codex-next-prompt: invalid Suggested next prompt line"}'; then
	printf '%s\n' "archive smoke: invalid Stop output: $stop_output" >&2
	exit 1
fi

printf '%s\n' 'smoke: extracted POSIX archive launcher passed'
