# AGENTS.md

Test: `just test` | Before committing: `just ready`

Read [CONTRIBUTING.md](CONTRIBUTING.md) first. It gives the prerequisites, the
setup, the code style, the test conventions and the pull request workflow. All
of it applies to agents the same as to people. This file has only the rules that
are specific to agents.

## Running tools

Run tools through `mise`, not from your path:

```bash
mise exec -- just test
```

`mise` is active in a person's shell. It supplies the versions that `.mise.toml`
declares. An agent's shell does not activate `mise`. Thus a bare `just` runs the
version that is installed globally, which is usually older.

The symptom is a check that fails here and passes in continuous integration, on
a file that nobody changed. If this occurs, find which version ran. Do this
before you accept the failure as real.

## Where the rules come from

[CONTRIBUTING.md](CONTRIBUTING.md#before-you-start) names the specification. If
a convention here and the specification do not agree, the specification is
correct. Say so. Do not follow the code.

## Commit trailer

When you commit with Claude Code, end the message with:

```
🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Task tracking

**Write the design in the design pages.** Write or change the page in
[osapi-io/specs](https://github.com/osapi-io/specs) before you build. If the
build shows that the page is wrong, correct the page. A design page in this
repository becomes incorrect when the code changes, and no check finds the
difference.
