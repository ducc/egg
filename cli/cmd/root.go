package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "egg",
	Short: "egg - the simple error aggregator",
	Long: `egg - the simple error aggregator
egg ingests errors and aggregates them
documentation available at https://github.com/ducc/egg`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
