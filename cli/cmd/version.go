package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.0.1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of egg",
	Long:  `all software has versions. this is egg's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("egg %s\n", version)
	},
}
