package cmd

import (
	"context"
	"fmt"
	"os"

	cfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/mattn/go-colorable"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/one2nc/cloudlens/internal/color"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/one2nc/cloudlens/internal/view"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	profile, region string
	version         = "dev"
	commit          = "dev"
	date            = "today"
	rootCmd         = &cobra.Command{
		Use:   `cloudlens`,
		Short: `cli for aws services`,
		Long:  `cli for aws services[s3, ec2, security-groups, iam]`,
		Run:   run,
	}
	out = colorable.NewColorableStdout()
)

func init() {
	rootCmd.AddCommand(versionCmd())
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Read aws profile")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "Read aws region")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile("./cloudlens.log", mod, 0644)
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
	//TODO profiles and regions should under aws
	profiles := readAndValidateProfile()
	if profiles[0] == "default" && len(region) == 0 {
		region = getDefaultAWSRegion()
	} else if len(region) == 0 {
		region = "ap-south-1"
	}

	regions := readAndValidateRegion()
	sess, err := aws.GetSession(profiles[0], regions[0])
	if err != nil {
		panic(fmt.Sprintf("aws session init failed -- %v", err))
	}
	ctx := context.WithValue(context.Background(), internal.KeySession, sess)
	app := view.NewApp()

	// TODO pass the AWS session instead of profiles and regions
	if err := app.Init(ctx, profiles, regions, version); err != nil {
		panic(fmt.Sprintf("app init failed -- %v", err))
	}
	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("app run failed %v", err))
	}
}

func readAndValidateProfile() []string {
	profiles, err := aws.GetProfiles()
	if err != nil {
		panic(fmt.Sprintf("failed to read profiles -- %v", err))
	}
	profiles, isSwapped := config.SwapFirstIndexWithValue(profiles, profile)
	if !isSwapped {
		fmt.Printf("Profile '%v' not found, using profile '%v'... ", color.Colorize(profile, color.Red), color.Colorize(profiles[0], color.Green))
	}
	return profiles
}

func readAndValidateRegion() []string {
	regions := aws.GetAllRegions()
	regions, isSwapped := config.SwapFirstIndexWithValue(regions, region)
	if !isSwapped {
		fmt.Printf("Region '%v' not found, using %v..", color.Colorize(region, color.Red), color.Colorize(regions[0], color.Green))
	}
	return regions
}

func getDefaultAWSRegion() string {
	cfg, err := cfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load AWS SDK config: %v\n", err)
		os.Exit(1)
	}
	region := cfg.Region
	return region
}
