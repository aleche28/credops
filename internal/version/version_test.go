package version

import (
	"runtime/debug"
	"testing"
)

func buildInfo(mainVersion string, settings map[string]string) *debug.BuildInfo {
	bi := &debug.BuildInfo{Main: debug.Module{Path: "github.com/aleche28/credops", Version: mainVersion}}
	for k, v := range settings {
		bi.Settings = append(bi.Settings, debug.BuildSetting{Key: k, Value: v})
	}
	return bi
}

func TestResolve(t *testing.T) {
	const (
		rev  = "f1478190b81dffb27ae6301b13cb6f9e356a0b49"
		when = "2026-11-26T18:38:26Z"
		pseu = "v0.0.0-20261126183826-f1478190b81d+dirty"
	)
	vcs := map[string]string{"vcs.revision": rev, "vcs.time": when, "vcs.modified": "true"}

	tests := []struct {
		name                        string
		ldVersion, ldCommit, ldDate string
		bi                          *debug.BuildInfo
		want                        Info
	}{
		{
			name:      "ldflags win over build info",
			ldVersion: "v0.1.0", ldCommit: "abc1234def", ldDate: "2026-12-01T00:00:00Z",
			bi:   buildInfo(pseu, vcs),
			want: Info{Version: "v0.1.0", Commit: "abc1234def", Date: "2026-12-01T00:00:00Z"},
		},
		{
			name: "empty ldflags fall back to build info",
			bi:   buildInfo(pseu, vcs),
			want: Info{Version: pseu, Commit: rev, Date: when},
		},
		{
			name: "(devel) is not a version",
			bi:   buildInfo("(devel)", nil),
			want: Info{Version: devVersion, Commit: unknown, Date: unknown},
		},
		{
			name: "no build info gives defaults",
			bi:   nil,
			want: Info{Version: devVersion, Commit: unknown, Date: unknown},
		},
		{
			name:      "fields resolve independently",
			ldVersion: "v0.1.0",
			bi:        buildInfo("(devel)", vcs),
			want:      Info{Version: "v0.1.0", Commit: rev, Date: when},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolve(tt.ldVersion, tt.ldCommit, tt.ldDate, tt.bi)
			if got != tt.want {
				t.Errorf("resolve() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestInfoString(t *testing.T) {
	tests := []struct {
		name string
		info Info
		want string
	}{
		{
			name: "full commit is shortened",
			info: Info{Version: "v0.1.0", Commit: "f1478190b81dffb27ae6301b13cb6f9e356a0b49", Date: "2026-11-26T18:38:26Z"},
			want: "credops v0.1.0 (commit f147819, built 2026-11-26T18:38:26Z)",
		},
		{
			name: "defaults are printed as is",
			info: Info{Version: devVersion, Commit: unknown, Date: unknown},
			want: "credops dev (commit unknown, built unknown)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
