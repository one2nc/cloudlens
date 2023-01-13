package cmd

import (
	"github.com/one2nc/cloud-lens/internal/view"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   `start`,
	Short: `start cloudlens tui`,
	Long:  `start cli for aws services[s3, ec2]`,
	Run: func(cmd *cobra.Command, args []string) {
		app := view.NewApp()
		app.Init()
	},
}


func init(){
	rootCmd.AddCommand(startCmd)
}