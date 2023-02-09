package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/brianvoe/gofakeit"
	"github.com/one2nc/cloud-lens/internal/config"
	pop "github.com/one2nc/cloud-lens/populator"
	"github.com/spf13/cobra"
)

var lspop = &cobra.Command{
	Use:   `lspop`,
	Short: ``,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := config.GetSession(profile, region)
		if err != nil {
			log.Fatal("err: ", err)
		}

		errCB := pop.CreateBuckets(sess)
		if errCB != nil {
			log.Fatal("err: ", errCB)
		}

		regions := []string{
			"us-east-1", "us-east-2", "us-west-1", "us-west-2", "af-south-1", "ap-east-1",
			"ap-southeast-3", "ap-south-1", "ap-northeast-3", "ap-northeast-2",
			"ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1",
			"eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-north-1",
			"me-south-1", "me-central-1", "sa-east-1", "us-gov-east-1", "us-gov-west-1"}

		var sessions []*session.Session

		for i := 0; i < 4; i++ {
			gofakeit.Seed(0)
			sess, err := config.GetSession(profile, regions[gofakeit.Number(0, len(regions)-1)])
			if err != nil {
				log.Fatal("err: ", err)
			}
			sessions = append(sessions, sess)
		}

		errCEI := pop.CreateEC2Instances(sessions)
		if errCEI != nil {
			log.Fatal("err: ", errCEI)
		}

		errCKP := pop.CreateKeyPair(sessions)
		if errCKP != nil {
			log.Fatal("err: ", errCKP)
		}

		errIAM := pop.IamAwsSrv(sess)
		if errIAM != nil {
			log.Fatal("err: ", errIAM)
		}
	},
}

func init() {
	rootCmd.AddCommand(lspop)
}
