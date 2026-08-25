# Project Instructions

- Keep this repository skill-only. Do not add hooks, executable runtime code,
  native launchers, MCP servers, apps, or automatic invocation surfaces.
- Keep the `$next` skill local-only, offline, and standard-library-only.
- Use only the current conversation context. Do not read transcripts or history,
  call model APIs, persist data, or execute a recommended follow-up.
- Preserve explicit `$next` invocation and return only ready-to-send prompt text.
- Follow TDD: write a failing behavioral test first, make it pass minimally,
  then refactor without weakening the contract.
- Run the narrowest relevant test first, then `make test`. Run `make check` only
  when the portable packager exists. Do not hide failures or silently skip gates.
- Keep credentials, private signing keys, user data, generated packages, release
  archives, and local planning artifacts out of the repository.
- Stage explicit paths only when preparing a commit; do not use broad staging.
