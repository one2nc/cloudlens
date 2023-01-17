package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/spf13/cobra"
)

func CreateBuckets(sess *session.Session) error {
	s3Service := s3.New(sess)
	for i := 0; i < 5; i++ {
		_, err := s3Service.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("test" + "-bucket" + strconv.Itoa(gofakeit.Number(0, 999999999999999999))), CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: aws.String("ap-south-1")}})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func CreateEC2Instances(sess *session.Session) error {
	ec2Service := ec2.New(sess)

	params := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-12345678"), // specify the ID of the image you want to use
		InstanceType: aws.String("t2.micro"),     // specify the instance type
		MinCount:     aws.Int64(3),
		MaxCount:     aws.Int64(3),
	}

	_, err := ec2Service.RunInstances(params)
	if err != nil {
		fmt.Println("Error in creating instances:", err)
		return err
	}
	return nil
}

var lspop = &cobra.Command{
	Use:   `lspop`,
	Short: ``,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal("err: ", err)
		}
		sess, err := config.GetSession("test", "ap-south-1", cfg.AwsConfig)
		if err != nil {
			log.Fatal("err: ", err)
		}

		errCB := CreateBuckets(sess)
		if errCB != nil {
			log.Fatal("err: ", errCB)
		}
		errCEI := CreateEC2Instances(sess)
		if errCEI != nil {
			log.Fatal("err: ", errCEI)
		}
	},
}

func init() {
	rootCmd.AddCommand(lspop)
}
