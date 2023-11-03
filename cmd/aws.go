package cmd

import (
	"github.com/one2nc/cloudlens/internal"
	"github.com/spf13/cobra"
)

func awsCommand() *cobra.Command {

	command := cobra.Command{
		Use:   "aws",
		Short: "Select aws",
		Long:  "Selects aws as default cloud",
		Run: func(cmd *cobra.Command, args []string) {
			selectAWS()
		},
	}

	command.Flags().StringVarP(&profile, "profile", "p", "default", "Read aws profile")

	command.Flags().StringVarP(&region, "region", "r", "", "Read aws region")
	command.Flags().BoolVarP(&useLocalStack, "localstack", "l", false, "Use localsatck instead of AWS")

	return &command
}

func selectAWS() {
	cloudConfig.SelectedCloud = internal.AWS
	cloudConfig.AWSConfig.Profile = profile
	cloudConfig.AWSConfig.Region = region
	cloudConfig.AWSConfig.UseLocalStack = useLocalStack

	initView()
}
