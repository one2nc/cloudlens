package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "version",
	Short: "Print version/build info",
	Long:  "Print version/build information",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion("version")
	},
}

func init() {
	rootCmd.AddCommand(command)
}

func printVersion(fmat string) {
	fmt.Println(version)
}
