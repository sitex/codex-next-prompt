#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 3 ]]; then
  printf '%s\n' 'usage: scripts/verify-release-tag.sh PUBLIC_KEY TAG EXPECTED_PRIMARY_FINGERPRINT' >&2
  exit 2
fi

public_key=$1
tag=$2
expected_primary_fingerprint=$3
verification_home=$(mktemp -d "${RUNNER_TEMP:-${TMPDIR:-/tmp}}/release-gnupg.XXXXXX")
trap 'rm -rf "$verification_home"' EXIT HUP INT TERM
chmod 700 "$verification_home"
export GNUPGHOME=$verification_home

mapfile -t primary_fingerprints < <(
  gpg --batch --with-colons --show-keys --fingerprint "$public_key" |
    awk -F: '$1 == "pub" {primary = 1; next} primary && $1 == "fpr" {print $10; primary = 0}'
)
if [[ ${#primary_fingerprints[@]} -ne 1 ]]; then
  printf '%s\n' 'release signing key must contain exactly one primary key' >&2
  exit 1
fi
if [[ ${primary_fingerprints[0]} != "$expected_primary_fingerprint" ]]; then
  printf '%s\n' 'release signing key primary fingerprint does not match expected fingerprint' >&2
  exit 1
fi

gpg --batch --import "$public_key"
if ! verification_status=$(git verify-tag --raw "$tag" 2>&1); then
  printf '%s\n' "$verification_status" >&2
  exit 1
fi
mapfile -t valid_primary_fingerprints < <(
  awk '$1 == "[GNUPG:]" && $2 == "VALIDSIG" {print $NF}' <<<"$verification_status"
)
if [[ ${#valid_primary_fingerprints[@]} -ne 1 ]]; then
  printf '%s\n' 'release tag must contain exactly one valid signature' >&2
  exit 1
fi
if [[ ${valid_primary_fingerprints[0]} != "$expected_primary_fingerprint" ]]; then
  printf '%s\n' 'release tag signature primary fingerprint does not match expected fingerprint' >&2
  exit 1
fi
