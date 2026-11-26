// Package version reports which build of credops is running.
package version

import (
	"fmt"
	"runtime/debug"
)

// Build metadata, set at link time by `make build` with
// -ldflags "-X github.com/aleche28/credops/internal/version.Version=..." (ADR-0002).
// They must stay uninitialized: -X cannot override a value set in code.
// An empty value means "not injected", and Get falls back to the build info Go embeds.
var (
	Version string
	Commit  string
	Date    string
)

// Defaults when neither -ldflags nor the embedded build info provide a value.
const (
	devVersion = "dev"
	unknown    = "unknown"
)

// shortCommitLen is how much of the commit hash the human-readable line shows.
const shortCommitLen = 7

// Info describes the running build.
type Info struct {
	Version string // release tag or git describe output, e.g. v0.1.0 or f147819-dirty
	Commit  string // full commit hash
	Date    string // commit time in UTC, RFC 3339
}

// Get returns the running build's metadata: values injected with -ldflags win, then the
// build info Go embeds (go build in a checkout, go install), then the defaults.
func Get() Info {
	bi, _ := debug.ReadBuildInfo() // nil when unavailable; resolve handles that
	return resolve(Version, Commit, Date, bi)
}

// resolve picks each field independently. It is separate from Get so tests can drive it
// with made-up inputs: go test binaries carry no VCS information.
func resolve(ldVersion, ldCommit, ldDate string, bi *debug.BuildInfo) Info {
	var biVersion, biCommit, biDate string
	if bi != nil {
		// go build in a checkout records a pseudo-version (with +dirty); "(devel)" means
		// no version was available, e.g. outside a repository or with -buildvcs=false.
		if bi.Main.Version != "(devel)" {
			biVersion = bi.Main.Version
		}
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				biCommit = s.Value
			case "vcs.time":
				biDate = s.Value
			}
		}
	}

	return Info{
		Version: firstNonEmpty(ldVersion, biVersion, devVersion),
		Commit:  firstNonEmpty(ldCommit, biCommit, unknown),
		Date:    firstNonEmpty(ldDate, biDate, unknown),
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// String renders the build on one line, e.g.
// "credops v0.1.0 (commit f147819, built 2026-11-26T18:38:26Z)".
func (i Info) String() string {
	commit := i.Commit
	if len(commit) > shortCommitLen && commit != unknown {
		commit = commit[:shortCommitLen]
	}
	return fmt.Sprintf("credops %s (commit %s, built %s)", i.Version, commit, i.Date)
}
