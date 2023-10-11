package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-colorable"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/one2nc/cloudlens/internal/view"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	profile, region, gcpCredFilePath string
	isAWS, isGCP                     bool
	version                          = "dev"
	commit                           = "dev"
	date                             = "today"
	rootCmd                          = &cobra.Command{
		Use:   `cloudlens`,
		Short: `cli for aws services`,
		Long:  `cli for aws services[s3, ec2, security-groups, iam]`,
		Run:   run,
	}
	out = colorable.NewColorableStdout()
)

func init() {
	rootCmd.AddCommand(versionCmd(), updateCmd())
	rootCmd.PersistentFlags().BoolVar(&isAWS, "aws", false, "Select AWS")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Read aws profile")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "Read aws region")
	rootCmd.PersistentFlags().BoolVar(&isGCP, "gcp", false, "Select GCP")
	rootCmd.PersistentFlags().StringVar(&gcpCredFilePath, "cf", "", "Read GCP credential file")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	dirPath := "tmp"
	filename := "cloudlens.log"

	// Create the full path to the file
	filePath := filepath.Join(dirPath, filename)

	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(filePath, mod, 0644)
	if err != nil {
		log.Printf("Could not open cloudlens.log. Writing logs to stdout instead.")
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if err == nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})
	}

	cloudConfig := config.NewCloudConfig()
	if isAWS {
		cloudConfig.SelectedCloud = internal.AWS
		cloudConfig.AWSConfig.Profile = profile
		cloudConfig.AWSConfig.Region = region
	} else if isGCP {
		cloudConfig.SelectedCloud = internal.GCP
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gcpCredFilePath)
		cloudConfig.CredFilePath = gcpCredFilePath
	}
	//TODO profiles and regions should under aws
	//var sess *session.Session

	app := view.NewApp()
	app.Init(version, cloudConfig)

	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("app run failed %v", err))
	}

}
