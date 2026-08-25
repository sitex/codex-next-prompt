# Project Instructions for 0.2.0

- Keep this repository skill-only. Don't add executable runtime code, native
  launchers, MCP servers, apps, or automatic invocation surfaces.
- Keep the `$next` skill local-only, offline, and standard-library-only.
- Use only the current conversation context. Don't read conversation records or
  history, call model APIs, persist data, or execute a recommended follow-up.
- Preserve explicit `$next` invocation and return only ready-to-send prompt text.
- Return one default prompt. Return two or three separate prompts only for a
  real fork. Keep cosmetic alternatives in one prompt.
- Ask for missing context. Don't invent values or emit placeholders.
- Follow TDD: write a focused failing contract test first, make it pass minimally,
  then refactor without pinning prose.
- Run the narrowest relevant test first, then `make test`, `make clean`, and
  `make check`. Don't hide failures or silently skip gates.
- Keep credentials, private signing keys, user data, generated packages, release
  ZIP files, and local planning artifacts out of the repository.
- Stage explicit paths only when preparing a commit. Don't use broad staging.
