package cmd

import (
	"fmt"
	"os"

	"github.com/one2nc/cloud-lens/internal/view"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/rs/zerolog/log"
)

var startCmd = &cobra.Command{
	Use:   `start`,
	Short: `start cloudlens tui`,
	Long:  `start cli for aws services[s3, ec2]`,
	Run: func(cmd *cobra.Command, args []string) {
		mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
		file, err := os.OpenFile("./log.txt", mod, 0644)
		if err != nil {
			panic(err)
		}
		defer func() {
			if file != nil {
				_ = file.Close()
			}
		}()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})

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
