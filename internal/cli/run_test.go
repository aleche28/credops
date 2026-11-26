package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aleche28/credops/internal/version"
)

func run(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = Run(context.Background(), args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestVersionCommandAndFlagMatch(t *testing.T) {
	want := version.Get().String() + "\n"

	for _, args := range [][]string{{"version"}, {"--version"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			code, stdout, stderr := run(t, args...)
			if code != 0 {
				t.Errorf("exit code = %d, want 0 (stderr %q)", code, stderr)
			}
			if stdout != want {
				t.Errorf("stdout = %q, want %q", stdout, want)
			}
			if stderr != "" {
				t.Errorf("stderr = %q, want empty", stderr)
			}
		})
	}
}

func TestErrorsExitOneOnStderr(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"unknown command", []string{"foo"}},
		{"version takes no arguments", []string{"version", "extra"}},
		{"unknown flag", []string{"--no-such-flag"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := run(t, tt.args...)
			if code != 1 {
				t.Errorf("exit code = %d, want 1", code)
			}
			if !strings.HasPrefix(stderr, "credops: ") {
				t.Errorf("stderr = %q, want prefix %q", stderr, "credops: ")
			}
			if stdout != "" {
				t.Errorf("stdout = %q, want empty", stdout)
			}
		})
	}
}

func TestNoArgumentsPrintsHelp(t *testing.T) {
	for _, args := range [][]string{{}, nil} {
		code, stdout, stderr := run(t, args...)
		if code != 0 {
			t.Errorf("args %#v: exit code = %d, want 0 (stderr %q)", args, code, stderr)
		}
		if !strings.Contains(stdout, "Usage:") {
			t.Errorf("args %#v: stdout = %q, want help text", args, stdout)
		}
	}
}
