package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gaiaz-iusipov/rpi-violet/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommit: %s\nBuilt: %s\n", version.Version(), version.Commit(), version.Date())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
