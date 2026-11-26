package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aleche28/credops/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print current version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), version.Get().String())
			return err
		},
	}
}
