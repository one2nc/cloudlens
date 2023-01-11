package cmd

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/one2nc/cloud-lens/internal/ec2"
	"github.com/one2nc/cloud-lens/internal/s3"
)

func CreateAndListS3Buckets(sess *session.Session) []s3.BucketResp {
	s3Service := s3.NewS3Service(*sess)
	for i := 0; i < 5; i++ {
		s3Service.CreateBucket("test" + "-bucket" + strconv.Itoa(gofakeit.Number(0, 999999999999999999)))
	}
	bucketList, err := s3Service.ListBuckets()
	if err != nil {
		fmt.Println("Error while listing the bucket")
	}
	return bucketList
}

func CreateAndListEC2Instances(sess *session.Session) []ec2.EC2Resp {

	ec2Service := ec2.NewEc2Service(*sess)
	ec2Service.CreateInstances()

	ec2InstanceInfoList, err := ec2Service.GetInstances()
	if err != nil {
		fmt.Println("Error while listing the bucket")
	}
	return ec2InstanceInfoList
}
