# Architecture decision records

An ADR records one decision: the context, the options we considered, what we chose and what follows from it. [design.md](../design.md) describes the current design; the ADRs explain why it looks that way.

**Rules**

- Copy [template.md](template.md) to `NNNN-kebab-title.md` with the next free number.
- ADRs are append-only. To change a decision, write a new ADR with `Supersedes: ADR-NNNN`, and only update the old one's status to `Superseded` with a `Superseded by: ADR-NNNN` line.
- Statuses: `Proposed`, `Accepted`, `Rejected`, `Deprecated`, `Superseded`.
- Don't edit the table below by hand: run `python3 scripts/check_docs.py --write`.

<!-- adr-index:start -->
| ADR | Title | Status | Date |
| --- | ----- | ------ | ---- |
| [0001](0001-record-architecture-decisions.md) | Record architecture decisions as ADRs | Accepted | 2026-11-26 |
<!-- adr-index:end -->
