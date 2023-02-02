package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/one2nc/cloud-lens/internal/view"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   `cloudlens`,
	Short: `cli for aws services`,
	Long:  `cli for aws services[s3, ec2]`,
	Run:   run,
}

var (
	profile, region string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Read aws profile")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "ap-south-1", "Read aws region")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
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

	//get config
	cfg, err := config.Get()
	if err != nil {
		panic(fmt.Sprintf("app get config failed -- %v", err))
	}
	sess, err := config.GetSession(profile, region, cfg.AwsConfig)
	if err != nil {
		panic(fmt.Sprintf("aws session init failed -- %v", err))
	}

	ctx := context.WithValue(context.Background(), internal.KeySession, sess)

	app := view.NewApp()
	//init app
	if err := app.Init(ctx); err != nil {
		panic(fmt.Sprintf("app init failed -- %v", err))
	}
	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("app run failed %v", err))
	}
}
