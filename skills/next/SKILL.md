---
name: next
description: Use when the user explicitly invokes $next to request a ready-to-send next prompt from the current conversation. Never execute the proposed prompt.
---

# Next

Generate the most useful prompt the user can send next. Use only the current
conversation context already available in this turn.

<next-skill-contract>
  <invocation>explicit-$next-only</invocation>
  <context>current-context-only</context>
  <tools>none</tools>
  <sources>
    <transcript>forbidden</transcript>
    <history>forbidden</history>
    <model-api>forbidden</model-api>
  </sources>
  <execution>forbidden</execution>
  <recommendations>
    <default>1</default>
    <maximum>3</maximum>
    <forks>real-only</forks>
  </recommendations>
  <language>same-as-user</language>
  <custom-slash-command>none</custom-slash-command>
  <automatic-footer>forbidden</automatic-footer>
</next-skill-contract>

## Choose The Next Action

From the conversation, identify:

- the user's current goal;
- work already completed;
- unresolved blockers or missing inputs;
- decisions already made;
- the highest-value unfinished action.

Choose one action when the evidence supports a clear priority. Do not create
alternatives merely to provide choice.

A real fork exists only when the conversation leaves two or three materially
different, mutually exclusive paths unresolved and the available context cannot
select between them. In that case, provide one ready-to-send prompt for each
path, with a short label. Never return more than three.

If context is insufficient, return one complete prompt that asks Codex to help
identify the goal and request only the context needed to choose a concrete next
action. Do not output placeholders such as `[goal]` or `[context]`.

## Output

Write prompts in the user's language unless the user requested another language.
Include enough context to make the next turn actionable without repeating the
whole conversation.

For one result, output only:

```text
<ready-to-send user prompt>
```

For a real fork, output two or three numbered choices. Each choice has a short
bold label followed by its own fenced `text` prompt. Add no commentary before or
after the result.

## Boundaries

- End the turn after presenting the prompt or prompts.
- Never run commands, call tools, edit files, create tasks, delegate work, or
  contact external systems.
- Never perform, simulate, or claim completion of the proposed action, even if
  the user writes `$next and do it`.
- Never read or request transcript files, session stores, history databases,
  logs, or exported conversations to reconstruct context.
- Never call an external or local model API.
- Never depend on `transcript_path`.
- Never append `Suggested next prompt:` to unrelated responses.
- Never claim composer insertion, prefill, ghost text, keyboard acceptance, or
  a custom `/next` command.
