package cmd

import (
	"os"

	"github.com/one2nc/cloudlens/internal"
	"github.com/spf13/cobra"
)

func gcpCommand() *cobra.Command {

	command := cobra.Command{
		Use:   "gcp",
		Short: "Select gcp (set GOOGLE_APPLICATION_CREDENTIALS env variable or use \"gcp\" sub-command )",
		Long:  "Selects gcp as default cloud",
		Run: func(cmd *cobra.Command, args []string) {
			selectGCP()
		},
	}

	command.Flags().StringVarP(&gcpCredFilePath, "cf", "", "", "Read GCP credential file")
	command.MarkFlagRequired("cf")
	return &command
}

func selectGCP() {
	cloudConfig.SelectedCloud = internal.GCP
	os.Setenv(internal.GOOGLE_APPLICATION_CREDENTIALS, gcpCredFilePath)
	cloudConfig.CredFilePath = gcpCredFilePath

	initView()
}
