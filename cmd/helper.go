package cmd

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/one2nc/cloud-lens/internal/aws"
)

func CreateAndListS3Buckets(sess *session.Session) []aws.BucketResp {
	s3Service := aws.NewS3Service(*sess)
	for i := 0; i < 5; i++ {
		s3Service.CreateBucket("test" + "-bucket" + strconv.Itoa(gofakeit.Number(0, 999999999999999999)))
	}
	bucketList, err := s3Service.ListBuckets()
	if err != nil {
		fmt.Println("Error while listing the bucket")
	}
	return bucketList
}

func CreateAndListEC2Instances(sess *session.Session) []aws.EC2Resp {
	ec2Service := aws.NewEc2Service(*sess)
	ec2Service.CreateInstances()
	ec2InstanceInfoList, err := ec2Service.GetInstances()
	if err != nil {
		fmt.Println("Error while listing the bucket")
	}
	return ec2InstanceInfoList
}
