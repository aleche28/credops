// Command credops monitors expiring Entra ID application credentials across tenants.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aleche28/credops/internal/cli"
)

func main() {
	// Ctrl-C or SIGTERM cancels ctx, so long-running commands (sync, agent) can stop cleanly.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop() // explicit, not deferred: os.Exit skips deferred calls
	os.Exit(code)
}
