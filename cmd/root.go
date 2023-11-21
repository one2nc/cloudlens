package cmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/one2nc/cloudlens/internal/view"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	profile, region, gcpCredFilePath, localStackPort string
	useLocalStack                                    bool
	version                                          = "dev"
	commit                                           = "dev"
	date                                             = "today"
	rootCmd                                          = &cobra.Command{
		Use:   `cloudlens`,
		Short: `cli for cloud services`,
		Long:  `cli for cloud services[aws,gcp]`,
		Run:   run,
	}
	cloudConfig config.CloudConfig
	out         = colorable.NewColorableStdout()
)

func init() {
	rootCmd.AddCommand(versionCmd(), updateCmd(), awsCommand(), gcpCommand())

}

func Execute() {
	file := ensureLogging()
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	initView()
}
func initView() {
	app := view.NewApp()
	app.Init(version, cloudConfig)

	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("app run failed %v", err))
	}
}

func ensureLogging() *os.File {
	filename := "cloudlens.log"
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(filename, mod, 0644)
	if err != nil {
		log.Printf("Could not open cloudlens.log. Writing logs to stdout instead.")
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	if err == nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})
	}

	return file
}
