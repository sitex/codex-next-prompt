# `$next` Examples

These examples describe the Codex Next Prompt 0.2.0 standalone skill. Invoke `$next` explicitly. Each
result is ready-to-send prompt text, and none of the results are executed by the
skill.

## One clear action

Context: the implementation is complete, but focused tests haven't run.

Invoke:

```text
$next
```

Output:

```text
Run the focused tests for the changed module, fix any regression caused by the change, then report the commands and results.
```

## A real fork

Context: a new API can use either the existing synchronous client or a new
asynchronous boundary, and the choice changes the public interface.

Output:

1. **Keep the synchronous interface**

```text
Keep the public API synchronous, implement the new call through the existing client, and add focused contract tests for the new behavior.
```

2. **Introduce an asynchronous boundary**

```text
Change the public API to an asynchronous boundary, document the compatibility impact, and add focused tests for cancellation and errors.
```

These are separate prompts because they lead to materially different designs.

## Cosmetic alternatives stay together

Context: a heading could be shorter or more descriptive without changing the
document's structure.

Output:

```text
Rewrite the heading for clarity, compare one concise and one descriptive wording, choose the better fit for the surrounding document, and apply only that choice.
```

Minor wording options don't justify separate prompts.

## Missing context

Context: the user asks for a deployment prompt, but no target environment is
known.

Output:

```text
Ask which environment and deployment method should be used, then wait for the answer before drafting deployment steps.
```

The result asks for what is missing. It doesn't invent a host, insert a
placeholder, or assume a provider.

## Russian conversation

Context: the conversation is in Russian and a bug fix is ready for verification.

Output:

```text
Запусти точечный тест для исправленного сценария, затем полный набор проверок и кратко перечисли команды и результаты.
```

## Never execute

Context: the best next action is to update a file and run tests.

Output:

```text
Обнови документацию под текущий контракт, затем запусти точечные проверки и полный make check.
```

`$next` returns that text only. It doesn't edit the documentation or run the
checks.
