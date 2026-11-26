---
name: roadmap-next
description: Determine the next task to tackle in the current credops milestone. Use when the user asks "what's next", "what should I do next", "what's left in this milestone", "where do I start", or otherwise asks for the next roadmap task. Grounds the answer in docs/design.md §8, the milestone's GitHub issues and their dependencies, open decision and tech-debt issues, the ADRs, git history and the actual source tree instead of inventing work.
model: sonnet
---

# Roadmap next

Answer the recurring "what's next?" question with a ranked, repo-grounded recommendation. This is a
**read-only planning** skill: do not edit code, docs, issues or milestones.

## Sources of truth (read these, in this order)

1. `docs/design.md` §8 Milestones: the milestone plan. The **current milestone** is the first one
   with unfinished work. Each milestone is a GitHub milestone named like its heading (for example
   "M1 — MVP"); its tasks are issues in that milestone:
   `gh issue list --milestone "<name>" --state all --json number,title,state,labels,body`.
2. Issue bodies: each task issue says `Part of <milestone>`, then `Goal:`, `Steps:`, `Done when:`
   and `Depends on #NN, ...`. The `Depends on` line is the dependency graph; `Done when` is the
   acceptance criteria. Use them rather than inventing either.
3. `decision` issues: open questions from design.md §7, each closed by an ADR. A task that depends
   on an open decision is not startable until that decision is recorded.
4. `docs/adr/` (index in `docs/adr/README.md`): decisions already made. An Accepted ADR overrides
   the design doc where they disagree.
5. Open `tech-debt` issues: `gh issue list --label tech-debt --state open`. Note which ones the
   current milestone's work touches or would naturally close.
6. The rest of `docs/design.md` (§4 detailed design, §6 cross-cutting concerns) for the target
   behaviour of any task, and §2 goals when a task's *why* is unclear.

The [credops project board](https://github.com/users/aleche28/projects/3) mirrors the milestones;
`gh project item-list 3 --owner aleche28 --format json` shows its status columns if they are in use.

## Steps

### 1. Establish where we actually are

- Read the current milestone from design.md §8 (first with open work) and list its issues, open
  and closed.
- Run `git log --oneline -15` and `git status --short` to separate **landed** work from
  **uncommitted** work from **promised-but-absent** work. Uncommitted changes count as in-flight,
  not done. Check open PRs too: `gh pr list --state open`.
- For each open issue, spot-check the source tree (Go module, `cmd/credops`, packages, tests,
  CI workflows) to confirm whether it already exists. A closed issue with no code, or code with no
  closed issue, is a mismatch to report. Do not trust issue state alone.

### 2. Build the dependency picture

- From the `Depends on` lines, work out which open issues are **startable now** (every dependency
  closed or landed) and which are blocked, and by what.
- An open `decision` issue that blocks startable-looking work is itself a candidate: resolving it
  means writing an ADR with the `adr` skill, which the user must agree to.
- List open `tech-debt` issues the milestone's work touches. These are not optional extras.

### 3. Rank candidate tasks

For each startable candidate, weigh:

- **Blocking power**: how many open issues list it, directly or transitively, in `Depends on`?
  Foundations (module skeleton, CI, store schema, the fake Graph seam) and blocking decisions
  usually win.
- **Dependencies met?**: is it actually startable now, or does it need a decision, package or
  schema that doesn't exist yet?
- **Fit with the milestone goal**: favour work that completes the current milestone over pulling
  forward a later one. Respect design.md's notes on what is not optional (for example M1 item 6).
- **Existing partial progress**: finish in-flight work (uncommitted changes, open PRs) before
  opening a new front.

### 4. Output

Produce a concise recommendation with these sections:

- **Current milestone**: name and a one-line status (what's done, what's left).
- **Recommended next task**: the single next thing, in imperative form, with its issue number.
- **Why now**: what it unblocks (by issue number), which decision or tech-debt issue it closes,
  and why not something else.
- **Files to touch**: concrete paths (new or existing), plus any docs the change will need under
  the `docs-maintenance` skill (design.md, README, an ADR).
- **Acceptance criteria**: from the issue's `Done when`, sharpened where the design doc is more
  specific (commands, output, exit codes, tests against the fake Graph).
- **Also considered / deferred**: 1–3 alternates with a one-line reason each.

Cite `file:line`, doc and section, or `#NN` for every claim. If design.md, the issues and the code
disagree, say so explicitly and recommend which to trust.

## Rules

- **Read-only.** Never write code, edit docs, or change issues, milestones or the project board
  here. Recording a decision is the `adr` skill's job; syncing docs after work lands is the
  `docs-maintenance` skill's.
- **Never invent tasks** not grounded in the current milestone, its issues, open decisions or
  design.md. If something needed is missing from the issues, say so as a gap; don't fold it into
  the recommendation as if it were planned.
- **Never trust issue state or a status line blindly**: verify against git and the tree.
- If the current milestone looks complete, say so and recommend its close-out (close the GitHub
  milestone, update design.md's `Status` line and the README status banner) or the first task of
  the next milestone.
- Keep it short: one recommendation, not a menu. Alternates go in the deferred section.
- If the user names a milestone or issue explicitly, honour that scope instead of re-deriving the
  current milestone.
