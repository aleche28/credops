package cli

import (
	"context"
	"fmt"
	"io"
)

// Run executes the credops command line with args (without the program name), writing to
// stdout and stderr, and returns the process exit code. It never exits the process itself.
//
// Exit codes are provisional until #9 defines the contract: 0 success, 1 any error.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	root := newRootCmd()
	// Never pass nil: SetArgs(nil) makes cobra fall back to os.Args.
	if args == nil {
		args = []string{}
	}
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)

	if err := root.ExecuteContext(ctx); err != nil {
		// Nothing sensible to do if stderr itself can't be written.
		_, _ = fmt.Fprintf(stderr, "credops: %v\n", err)
		return 1
	}
	return 0
}
