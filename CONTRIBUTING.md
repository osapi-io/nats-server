# Contributing

We welcome contributions to nats-server. Read this document before you open a
pull request. It gives the prerequisites, the setup, the conventions for code,
and the pull request workflow.

## Before you start

- Read the [Code of Conduct](CODE_OF_CONDUCT.md). It applies to all
  communication in this repository.
- **Design pages.** The design of this repository is in
  [osapi-io/specs](https://github.com/osapi-io/specs), under
  `components/nats-server/`. Read it before you change behaviour. When you
  change behaviour, change the design page too. A design page in this repository
  becomes incorrect when the code changes, and no check finds the difference.
- **Examine the work that exists.** Is there a pull request for it? Do issues
  discuss the change that you want to make? Use what these discussions say in
  your work.
- **Backwards compatibility.** Does your change break the programs that use
  nats-server now? A change that is backwards compatible is much more likely to
  merge. Find a method that keeps this compatibility, if you can. If you cannot,
  open an issue first. Then we can discuss the API change before you write the
  pull request.

## Prerequisites

Install the tools with [mise]:

```bash
mise install
```

- **[Go].** nats-server is written in Go. We support the latest two major Go
  versions. Make sure that your version is one of them.
- **[uv].** A Python package runner. `just md-fmt` formats Markdown with
  [mdformat] through `uvx`. It installs nothing into the repository.
- **[just].** A task runner. It builds, tests, formats and runs other
  development tasks. Install it with `brew install just`.

### Claude Code

If you use [Claude Code] for development, install this plugin from the default
marketplace:

```
/plugin install commit-commands@claude-plugins-official
```

- **commit-commands.** It gives the `/commit` and `/commit-push-pr` slash
  commands. They follow the commit conventions of the project.

**Write the design in the design pages.** Write or change the page in
[osapi-io/specs](https://github.com/osapi-io/specs) before you build. If the
build shows that the page is wrong, correct the page. A design page in this
repository becomes incorrect when the code changes, and no check finds the
difference.

## Setup

Get the shared justfiles and install all dependencies:

```bash
just fetch
just deps
```

## Project structure

- **`pkg/server/`.** The public SDK. It is the only package that consumers
  import.
  - `types.go`. `Server`, `Options`. It holds only type declarations.
  - `server.go`. `New()`, `Start()`, `Stop()`.
  - `server_wrapper.go`. `NATSServerInstance`, an interface over the embedded
    NATS server. Tests use it to replace the server.
  - `logger.go`. `SlogWrapper`. It connects `slog` to the logger interface of
    the embedded server.
  - `mocks/`. Generated mocks. See [Test doubles](#test-doubles).
- **`examples/`.** Programs that you can run, one for each auth mode.

## Code style

[gofumpt] formats the Go code and [golangci-lint] lints it. CI makes these
checks.

```bash
just go-fmt-check   # Check formatting
just go-fmt         # Auto-fix formatting
just go-vet         # Run linter
```

`.golangci.yml` declares the linters that run. Read them there. This page does
not have a list, because a copied list becomes incorrect when the configuration
changes. The formatter does not format generated files (`*.gen.go`, `*.pb.go`).

### Documentation

[mdformat] formats the Markdown files through `uvx`. CI checks this format.

```bash
just md-fmt-check   # Check formatting
just md-fmt         # Auto-fix formatting
```

## Code standards

### Function signatures

A function with parameters has one parameter on each line. The closing
parenthesis and the return types are on a line of their own:

```go
func FunctionName(
    param1 type1,
    param2 type2,
) (returnType, error) {
}
```

A function with no parameters stays on one line:

```go
func Name() string {
}
```

Thus a new parameter shows in the diff as one added line. The diff does not show
a rewritten signature.

### File naming

Name a file for what it holds. Do not use `helpers.go`, `utils.go` or names of
that kind. These names tell where code is, not what it is. Also, they collect
all code that has no other file.

`types.go` holds only type declarations: structs, interfaces, constants, and
aliases. Put a function in a file named for what the function does.

Name a test file for the production file that it tests. If the tests become too
large to read, split the production file first. Then each test file keeps its
production file. Do not split tests away from the file that they test.

### Go patterns

- Error wrapping: `fmt.Errorf("context: %w", err)`. Then the chain names each
  layer that the error went through. `errors.Is` and `errors.As` can still
  examine it.
- Return early. Do not put the successful path inside conditionals.
- Unused parameters: change the name to `_`.
- Import order: standard library, third party, then local. Put a blank line
  between each group.

### Test doubles

Use `mockgen` to generate a double for an interface that this organization
defines. Commit the double. Do not write a struct by hand to satisfy one.

Generated mocks are in a `mocks` package next to the code that they mock. A
`generate.go` file holds the directive that makes them:

```go
package mocks

//go:generate go tool go.uber.org/mock/mockgen -source=../types.go -destination=types.gen.go -package=mocks
```

Go finds the generator through the tool dependencies of the module. Thus every
checkout runs the version that `go.mod` records. Destination files end in
`.gen.go`. Commit them. Do not use `gen/` for mocks. API code generation uses
that name.

A sibling package cannot work for an **unexported** interface. The mock must
import the package to name the types in the interface. The tests of the package
must import the mock. Thus, generate the mock into the package. Use a
destination file that only tests compile. Then the mocking library is not in the
dependency graph of the code that imports the package:

```go
// generate.go, in the package that declares the interface
package thispackage

//go:generate go tool go.uber.org/mock/mockgen -source=thing.go -destination=thing.gen_test.go -package=thispackage
```

In the two cases, the directives are in a `generate.go` that holds no code. The
generated file has `.gen` in its name. This tells a reader not to edit it.

If many call sites repeat the same expectations, write a constructor that
returns a configured mock. Do not add a hand-written type. The generated mock
still satisfies the interface.

Write three types of double by hand. To generate them gives no benefit:

- A double for a standard library interface: `net.Conn`, `fs.File`, `io.Writer`,
  `slog.Handler`. These do not change when our code changes.
- A double that has a real implementation of the behavior under test. An example
  is a signature with a key pair that the test really generates.
- A recorder for a dependency that a goroutine calls, if the test cannot join
  that goroutine. A generated mock would assert a call count at a time that the
  test cannot know. Write this reason where you define the recorder.

### File headers

Every `.go` file MUST start with the MIT license header. For the exact format,
look at a Go file in the repository. A file with a build tag has `//go:build` on
line 1, then a blank line, then the header.

## Testing

```bash
just test           # Run all tests (lint + unit + coverage)
just go-unit       # Run unit tests only
just go-unit-cov   # Generate coverage report
go test -run TestName -v ./pkg/server/...  # Run a single test
```

The coverage gate is 100%. `just test` fails if the total coverage is less than
that. Thus a change that adds untested code fails locally and in CI:

```bash
just go-unit-cov-check   # Report coverage and fail below the target
```

`.github/codecov.yml` and the shared `go` justfile module declare the target.
Change the two at the same time.

### Test file conventions

- Public tests: `*_public_test.go` in the `_test` package of the package. They
  test the exported API. This is the default.
- Internal tests: `*_test.go` in the same package. They test the code that the
  exported API cannot reach.
- Suite naming: `*_public_test.go` → `{Name}PublicTestSuite`, `*_test.go` →
  `{Name}TestSuite`.
- `testify/suite` with table-driven cases.
- One suite method for each function under test. Success, errors, and edge cases
  are rows in one table. They are not separate methods.
- `export_test.go` gives external tests access to unexported symbols, with an
  alias or a setter. Do not use an alias to test again behavior that the test of
  the caller already reaches. Use the pattern for a helper with its own
  contract.

External tests in this repository are in `package server_test`. Their tables
have `validateFunc` callbacks.

## Before committing

Before you commit, run `just ready`. It makes sure that the generated code, the
package docs, the format and lint are all current:

```bash
just ready
```

## Branching

Make all changes on a feature branch. Make the branch from `main`. Name it
`type/short-description`, where `type` is the [Conventional Commits] type:

- `feat/add-retry-logic`
- `fix/null-pointer-crash`
- `docs/update-api-reference`
- `refactor/simplify-handler`
- `chore/update-dependencies`

If you are on `main`, the `/commit` command of Claude Code makes a branch for
you.

## Commit messages

Follow [Conventional Commits] with the 50/72 rule:

- **Subject line**: max 50 characters, imperative mood, capitalized, no period
- **Body**: wrap at 72 characters, separated from subject by a blank line
- **Format**: `type(scope): description`
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- Say "what" and "why", not "how"

Write commit messages that tell what changed. Do not put too many commits on a
pull request. Most pull requests have one commit. A large pull request can have
a small number. Use git squash and rebase to combine commits.

## Submitting a PR

- **Describe your changes.** Say what changed and why. A reviewer must not have
  to read the diff to learn the reason for it.
- **Issue/PR links.** Link to related issues or pull requests. Tell how your
  changes are different from that work, or how they add to it.
- **Examples.** Add examples or screenshots that show the effect of your
  changes.
- **Draft PRs.** If your changes are not complete and you want to discuss them,
  open the pull request as a draft. Add a comment to start a discussion. Use
  comments, not the description. Then you can update the description later, and
  the discussion stays.

## AI usage

We write this repository with help from AI. The [AI Usage Policy](AI_POLICY.md)
applies to all contributions. Tell which tool you used. Make sure that you can
explain what your change does without help from AI tools.

## FAQ

> I want to contribute, where do I start?

We welcome all types of contribution, a typo fix or a new feature. You can also
vote on issues, comment on them, or help to answer questions.

> I'm stuck, where can I get help?

If you have questions, open a [Discussion] on GitHub.

[claude code]: https://claude.ai/code
[conventional commits]: https://www.conventionalcommits.org
[discussion]: https://github.com/osapi-io/nats-server/discussions
[go]: https://go.dev
[gofumpt]: https://github.com/mvdan/gofumpt
[golangci-lint]: https://golangci-lint.run
[just]: https://just.systems
[mdformat]: https://pypi.org/project/mdformat/
[mise]: https://mise.jdx.dev
[uv]: https://docs.astral.sh/uv/
