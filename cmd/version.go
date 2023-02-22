package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//var short bool

var command = &cobra.Command{
	Use:   "version",
	Short: "Print version/build info",
	Long:  "Print version/build information",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion("version")
	},
}

//command.PersistentFlags().BoolVarP(&short, "short", "s", false, "Prints K9s version info in short format")

func init() {
	rootCmd.AddCommand(command)
}

func printVersion(fmat string) {
	fmt.Println("v0.0.8")
}
