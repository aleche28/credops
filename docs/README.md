# credops docs

| Where                  | What it holds                                                                          |
| ---------------------- | -------------------------------------------------------------------------------------- |
| [design.md](design.md) | The current design: architecture, notification and configuration models, milestones. Kept up to date. |
| [adr/](adr/README.md)  | Architecture decision records: why the design is the way it is. Append-only.           |

Folders are added when they have content. The planned ones are:

- `guides/`: task-oriented how-tos, such as onboarding a tenant or setting up a channel.
- `reference/`: config schema, CLI, output and permissions reference, generated from code where possible.

**Keeping docs in order**

- A change to a public contract (CLI, machine-readable output, exit codes, config, required Graph permissions), a core dependency or milestone scope needs an ADR, and `design.md` updated to match.
- Every doc is linked from this page or from an index it links to.
- `python3 scripts/check_docs.py` checks ADR numbering and metadata, the ADR index, supersede links and relative links; `--write` regenerates the ADR index.
- With Claude Code, the `adr` and `docs-maintenance` skills in `.claude/skills/` follow these rules; `roadmap-next` picks the next task in the current milestone.
