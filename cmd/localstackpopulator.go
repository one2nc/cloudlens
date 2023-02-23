package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/brianvoe/gofakeit"
	"github.com/one2nc/cloudlens/internal/config"
	pop "github.com/one2nc/cloudlens/populator"
	"github.com/spf13/cobra"
)

var lspop = &cobra.Command{
	Use:   `lspop`,
	Short: ``,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := config.GetProfiles()
		if err != nil {
			panic(fmt.Sprintf("failed to read profiles -- %v", err))
		}
		if !config.LookupForValue(profiles, profile) {
			profile = profiles[0]
		}
		if len(region) == 0 {
			region = "ap-south-1"
		}
		sess, err := getSession(profile, region)
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
		sessDef, err := getSession(profile, pop.GetDefaultAWSRegion())
		sessions = append(sessions, sessDef)
		for i := 1; i < 4; i++ {
			gofakeit.Seed(0)
			sess, err := getSession(profile, regions[gofakeit.Number(0, len(regions)-1)])
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

		errSQS := pop.CreateQueueAndSetMessages(sessions)
		if errSQS != nil {
			log.Fatal("err: ", errSQS)
		}

		errLambda := pop.CreateLambdaFunction(sessions)
		if errLambda != nil {
			log.Fatal("err: ", errLambda)
		}
	},
}

func init() {
	rootCmd.AddCommand(lspop)
}

func getSession(profile, region string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{Config: aws.Config{
		//TODO: remove hardcoded enpoint
		Endpoint:         aws.String(localstackEndpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	},
		Profile: profile})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return nil, err
	}
	return sess, nil
}
