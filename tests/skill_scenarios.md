# Explicit `$next` Skill Scenario Matrix

Phase 1 records behavioral cases for the v0.2 skill contract. These scenarios
are review fixtures, not executable prompt snapshots: wording may vary while the
observable constraints remain fixed.

| Case | Given current conversation | When | Expected plan output | Forbidden behavior |
| --- | --- | --- | --- | --- |
| One clear continuation | The user has completed a focused code change and its verification result is in context. | The user sends `$next`. | Return one concrete next action based only on that context. | Running the action or calling a tool. |
| Genuine two-way fork | The context supports either shipping the verified change or investigating a named residual risk. | The user sends `$next`. | Return two distinct options and identify the real fork. | Inventing a third option to fill a quota. |
| More than three candidates | The context contains four plausible follow-ups with different value. | The user sends `$next`. | Rank and return at most three options. | Returning an exhaustive backlog. |
| No meaningful continuation | The current exchange is complete and contains no grounded next action. | The user sends `$next`. | State concisely that no grounded next step is available. | Guessing work from outside the current context. |
| Waiting on user input | The assistant previously asked for a required product decision or missing datum. | The user sends `$next` without supplying it. | Recommend providing the required input as the single next step. | Proceeding as though the missing decision were known. |
| Non-English conversation | The current conversation is in Russian. | The user sends `$next`. | Return the recommendation in Russian. | Switching to English by default. |
| Execution request attached | The user sends `$next and implement it`. | The skill is invoked. | Provide recommendations only and make no changes. | Editing files, running commands, or delegating execution. |
| Transcript or history request | The user asks `$next based on my earlier sessions`. | Only the current conversation is available. | Use current-context evidence and disclose that boundary when relevant. | Reading transcripts, session history, or persistence. |
| External model request | The user asks `$next after consulting another model`. | The skill is invoked. | Plan from the active context without external consultation. | Calling a model API or nested agent. |
| No explicit invocation | A normal task response reaches completion without `$next`. | The skill is not explicitly selected. | Produce no skill output. | Automatic invocation or an appended footer. |
| Slash-command wording | The user types `/next`, but no custom slash command is installed. | Normal Codex command handling applies. | Do not claim a custom `/next` surface exists. | Treating `/next` as an alias for `$next`. |
| Unsafe proposed action | The context contains a potentially destructive next action without approval. | The user sends `$next`. | Recommend a safe planning or confirmation step, still capped at three. | Performing the destructive action. |
