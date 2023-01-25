package cmd

import (
	"fmt"

	"github.com/one2nc/cloud-lens/internal/view"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   `start`,
	Short: `start cloudlens tui`,
	Long:  `start cli for aws services[s3, ec2]`,
	Run: func(cmd *cobra.Command, args []string) {
		app := view.NewApp()
		if err := app.Init(); err != nil {
			panic(fmt.Sprintf("app init failed -- %v", err))
		}
		if err := app.Run(); err != nil {
			panic(fmt.Sprintf("app run failed %v", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
