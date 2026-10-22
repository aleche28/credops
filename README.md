# credops

Monitor expiring Entra ID app registration secrets and certificates across multiple tenants.

> **Status: design phase — nothing is implemented yet.**
> This repository currently contains the design only. There is no code, no release, and no
> installable binary. See [`docs/design.md`](docs/design.md).

---

## Why

Entra ID app registrations authenticate with client secrets and certificates, and those
credentials expire. When one expires unnoticed the failure is silent until it isn't — an
integration starts returning `401`, usually outside working hours, and the cause is rarely
obvious to whoever is on call.

Checking for this natively means going tenant by tenant, application by application, in the
portal. There is no single place to answer *"what expires in the next 30 days, across every
tenant we operate?"*

`credops` answers that question, keeps answering it on a schedule, and tells the right people
before the outage rather than after it.

## What it does

- Builds a **persistent inventory** of credentials on app registrations and service
  principals, across **many tenants**.
- **Alerts before expiry** through email, Teams, Slack and Jira, with thresholds, routing and
  suppressions that survive contact with thousands of applications.
- **Never fails silently** — "nothing is expiring" and "the monitor has been broken for a
  week" are different states, and the tool says which one you are in.
- Later: **manages** credentials — add, rotate and remove — as a staged lifecycle rather than
  a single destructive action.

### Planned — MVP

- Attach tenants; read-only, unattended collection
- Inventory of application and service principal credentials, with sync-status tracking
- List and threshold checks, machine-readable output, CI-friendly exit codes
- One notification channel end to end, with per-tier de-duplication
- Suppressions with expiry, and an unsuppressible alert on `credops`' own credentials
- Silent baselining on first sync, and a dry-run notification mode

### Planned — later

- Guided tenant onboarding
- Remaining channels, plus a proper Jira issue lifecycle
- Tags, ownership routing, layered configuration and `explain`
- Team-shared mode: one inventory, many operators
- Credential management: add / rotate / remove
- Key Vault credentials, and hygiene findings beyond expiry

## Key concepts

**One binary, three modes.** A one-shot sync for `cron` or a scheduled task, a long-running
agent that schedules itself, and an interactive CLI. The agent is a convenience, never a
requirement.

**The store is the contract.** The CLI and the agent never talk to each other directly; both
go through the state store. That is what lets a teammate on another machine drive a
collector they cannot reach, without `credops` becoming a network service.

**Read-only by default, writes run as you.** The unattended collector holds read-only
application permissions. Anything that modifies a credential runs with *your* delegated
identity, so your roles, Conditional Access and PIM apply, and the audit log names you.

**Silence is not "all clear".** Per-tenant sync status is first class. A failed
authentication is an alert, not a log line.

**Alert fatigue is a design constraint, not a tuning problem.** Per-tier de-duplication,
digests on chat channels, per-item issues in Jira, snooze-until suppressions that resurface
when they expire, and silent baselining when a new tenant is attached.

## Usage

> **Illustrative only — proposed, subject to change. None of this works yet.**

```console
$ credops tenant attach
$ credops sync
$ credops status
$ credops list --expiring-within 30d
$ credops check --expiring-within 30d        # exit code reflects findings, for CI
$ credops suppress <credential> --until 2026-12-01 --reason "migration in progress"
$ credops explain <credential>
$ credops channel test
$ credops agent run
```

## Documentation

- [`docs/design.md`](docs/design.md) — full design document: scope, goals and non-goals,
  architecture, notification and configuration models, alternatives considered, open
  questions, and milestones.

## Status and contributing

Design phase. The design document is the place to push back — particularly on the open
questions in §7.

## License

[MIT](LICENSE) © 2026 Alessio Chessa
