# 0002. Use Go for the implementation and cobra for the CLI

- Status: Accepted
- Date: 2026-11-26

## Context

credops is one binary invoked three ways: one-shot sync, agent and interactive CLI ([design.md §3](../design.md#3-design-overview)). It ships as signed, cross-platform releases for Linux, macOS and Windows (#29), and §6 asks for a minimal dependency surface because the tool requests tenant-wide read access to an identity system ([§6](../design.md#6-cross-cutting-concerns)).

The M1 issues (#1, #2) and the naming row in [§5](../design.md#5-alternatives-considered) already assume Go, but the language was never recorded. #1 also needs a CLI library. M1 alone brings around ten commands (`tenant attach`, `sync`, `status`, `list`, `check`, `suppress` with sub-actions, `version`), with global flags such as the configuration path and output format, and exit codes that become a public contract (#9).

## Considered options

Language: no alternatives were weighed. Go was the working assumption from the start, and this ADR records it.

CLI library:

- **cobra**: nested subcommands, persistent flags, generated help, and shell completion for bash, zsh, fish and PowerShell. It's widely used in Go CLIs, so contributors are likely to know it. It compiles in `spf13/pflag`, plus `inconshreveable/mousetrap` on Windows only (checked against v1.10.2). `SilenceUsage` and `SilenceErrors` have to be set so credops owns its error output and exit codes.
- **kong**: commands declared as tagged structs, with typed flags and no runtime dependencies. Less widely known, and shell completion needs an extra library.
- **urfave/cli v3**: no dependencies, built-in completion, `context` passed to every action. Smaller ecosystem, and v3 is recent.
- **Standard library `flag`**: no dependencies. Flags stop at the first positional argument, and subcommand dispatch, help and completion would all be written by hand.

## Decision

We will implement credops in Go and build the CLI on `github.com/spf13/cobra`.

- The `go` directive in `go.mod` tracks the older of the two supported Go releases, currently `go 1.26.0`. A `toolchain` line pins the release we develop and build with, currently `go1.27.x`.
- Commands live in `internal/cli`. `cmd/credops/main.go` only wires dependencies and calls `os.Exit` with the code returned by the CLI, so exit codes are decided in one place.

## Consequences

- Shell completion and consistent help come with cobra rather than being written by hand.
- The dependency surface grows by two modules (one of them Windows-only). Dependabot (#3) and govulncheck (#2) cover them.
- Commands must be built by constructor functions rather than package-level `cobra.Command` variables, so tests can run the CLI with fresh state and captured output.
- Durations like `30d` (#19, #20) aren't understood by `time.ParseDuration`, so they need a custom `pflag.Value`.
- Raising the `go` directive is a deliberate step, taken when a feature requires it or when the older release leaves support.
- Revisit if cobra goes unmaintained, or if its dependency set grows beyond what §6 tolerates.
