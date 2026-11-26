// Package cli builds the credops command tree and runs it.
package cli

import (
	"github.com/spf13/cobra"

	"github.com/aleche28/credops/internal/version"
)

func newRootCmd() *cobra.Command {
	// No Run: a bare `credops` prints help, and an unknown subcommand is an error.
	root := &cobra.Command{
		Use:           "credops",
		Short:         "monitor expiring credentials",
		Long:          `Monitor expiring Entra ID app registration secrets and certificates across multiple tenants.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version.Get().String(),
	}

	// not shorthand 'v' cause maybe I will use that for 'verbose'
	root.Flags().Bool("version", false, "Print current version")
	root.AddCommand(newVersionCmd())
	root.SetVersionTemplate("{{.Version}}\n")
	return root
}
