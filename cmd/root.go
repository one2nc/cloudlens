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

	//TODO profiles and regions should under aws
	profiles := readAndValidateProfile()
	if profiles[0] == "default" && len(region) == 0 {
		region = getDefaultAWSRegion()
	} else {
		region = "ap-south-1"
	}

	regions := readAndValidateRegion()
	//TODO Move this in the AWS folder
	sess, err := config.GetSession(profiles[0], getDefaultAWSRegion())
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
	profiles, err := config.GetProfiles()
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
	loop:
		for {
			var input string
			fmt.Printf("Region '%v' not found, would you like to pick one from regions[%v,..] ["+color.Colorize("y", color.Cyan)+"/"+color.Colorize("n", color.Red)+"]: ", color.Colorize(region, color.Red), regions[0])
			fmt.Scanln(&input)
			switch input {
			case internal.LowercaseY, internal.UppercaseY, internal.LowercaseYes, internal.UppercaseYes:
				break loop
			case internal.LowercaseN, internal.UppercaseN, internal.LowercaseNo, internal.UppercaseNo:
				fmt.Printf("Region '%v' not found, exiting..", region)
				os.Exit(0)
			}
		}
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
