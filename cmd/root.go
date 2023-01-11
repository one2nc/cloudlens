package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   `cloudlens`,
	Short: `cli for aws services`,
	Long:  `cli for aws services[s3, ec2]`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
