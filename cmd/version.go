package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is set via ldflags during build
var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gitstuff",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gitstuff version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}