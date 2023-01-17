package pop

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/brianvoe/gofakeit"
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
