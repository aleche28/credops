# 0003. Use SQLite via modernc.org/sqlite for the embedded store, with goose migrations

- Status: Accepted
- Date: 2026-11-27

## Context

Single-machine mode needs an embedded, file-based store with zero infrastructure, and it's the default first-run experience ([design.md §4.7](../design.md#47-state-store-and-deployment-modes)). This ADR answers the embedded half of open question [§7.5](../design.md#7-open-questions) (#6). The networked store belongs to M5.

Constraints:

- Releases are cross-platform builds for Linux, macOS and Windows (#29), and §6 asks for a minimal dependency surface ([§6](../design.md#6-cross-cutting-concerns)). A store that needs cgo would need a C cross-toolchain for every target.
- The inventory schema (#12) is relational: tenants, applications, service principals, credentials, owners, sync runs, suppressions and the notification log. Queries such as `list --within` and `check --within` (#19, #20) filter and join across those tables. #12 also requires versioned migrations.
- The embedded store must never be pointed at a network share, and the tool must refuse that misconfiguration (§4.7). SQLite's own documentation warns that its file locking is unreliable on network filesystems and that WAL mode doesn't work over them.

## Considered options

Engine:

- **modernc.org/sqlite**: SQLite transpiled to pure Go, used as a standard `database/sql` driver, so builds work with `CGO_ENABLED=0`. It's the most widely used cgo-free SQLite for Go. It brings about 11 modules (the `modernc.org/libc` family and a few small utilities) and makes the binary larger.
- **github.com/ncruces/go-sqlite3**: also pure Go with no cgo, shipped through its own `go-sqlite3-wasm` module, with fewer extra modules. It is pre-1.0 (v0.38), less widely used, and offers its own API alongside the `database/sql` driver.
- **github.com/mattn/go-sqlite3**: cgo bindings to the SQLite C library. It needs a C toolchain per release target and doesn't work with `CGO_ENABLED=0`.
- **go.etcd.io/bbolt**: a pure-Go key/value store. With no SQL, the filters, joins, indexes and migrations would all be written by hand, and it shares nothing with a networked SQL store.

Migrations:

- **pressly/goose v3**: an established library with a version table, SQL or Go migrations, and embedded migration files (`fs.FS`).
- **An in-house runner**: embedded `.sql` files applied in order, with the version kept in `PRAGMA user_version`. No dependency, but we would write and maintain it ourselves.

Network-share refusal:

- **Detect the filesystem type and refuse, with no override.**
- **Detect and refuse, with an override flag** for filesystems wrongly detected as remote.

## Decision

We will back single-machine mode with SQLite through `modernc.org/sqlite`, build with `CGO_ENABLED=0`, and apply schema migrations with `github.com/pressly/goose/v3` from SQL files embedded in the binary.

- Migrations are versioned SQL files embedded with `go:embed` and applied by goose when the store is opened. The store refuses to open a database whose schema version is newer than the binary knows.
- Before opening the database, credops checks the filesystem holding the store's directory (after resolving symlinks). It checks the filesystem type with `statfs` on Linux and macOS (for example NFS, SMB/CIFS, AFP), and on Windows it checks for a UNC path or a remote drive type. A remote store is refused with an error that points to shared mode. There is no override flag.

## Consequences

- Release builds (#29) are plain `CGO_ENABLED=0` cross-compiles, with no C toolchain anywhere in the pipeline.
- The schema is SQL, so a networked SQL store in M5 can probably reuse most of it and the migration tooling. M5 still chooses its own engine.
- The dependency surface grows by the modernc SQLite modules. Dependabot (#3) and govulncheck (#2) cover them.
- goose's own `go.mod` lists drivers for every database it supports (MySQL, Postgres, MSSQL, ClickHouse and others), but its core package imports none of them. We measured this with modernc.org/sqlite v1.63.0 and goose v3.30.2 on a minimal `CGO_ENABLED=0` linux/amd64 build. goose added four modules (`pressly/goose/v3`, `mfridman/interpolate`, `golang.org/x/sync`, `go.uber.org/multierr`) and took the binary from 9.5 MB to 11 MB. Importing anything outside goose's core package (such as `cmd/goose` or its database helpers) would pull drivers in, so we don't. If goose's core package ever starts bringing drivers with it, switch to an in-house runner; the migration SQL files would stay the same.
- Detecting remote filesystems is platform-specific code with a test per platform. Without an override, a filesystem we wrongly flag as remote has to be fixed in the detection list. If that comes up in practice, an override flag is the obvious follow-up.
- Revisit the engine if modernc.org/sqlite falls behind upstream SQLite security fixes, or if binary size becomes a real adoption problem.
