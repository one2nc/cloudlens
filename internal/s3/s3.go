package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service interface {
	CreateBucket(bucketName string) error
	ListBuckets() error
}

type s3Service struct {
	client s3.S3
}

func NewS3Service(sess session.Session) S3Service {
	// Create an Amazon S3 service client
	return s3Service{client: *s3.New(&sess)}
}

func (s s3Service) CreateBucket(bucketName string) error {
	opcb, err := s.client.CreateBucketWithContext(context.Background(), &s3.CreateBucketInput{Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: aws.String("ap-south-2")}})
	if err != nil {
		fmt.Println("Error occured while creating a bucket:", err.Error())
		return err
	}
	fmt.Println("The location of bucket is:", *opcb.Location)
	return nil
}

func (s s3Service) ListBuckets() error {
	lbop, err := s.client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Error in listing buckets")
		return nil
	}
	for i, buc := range lbop.Buckets {
		reg, err := s.client.GetBucketLocationWithContext(context.Background(), &s3.GetBucketLocationInput{Bucket: aws.String(*buc.Name)})
		if err != nil {
			fmt.Println("error getting bucket location")
			return err
		}
		fmt.Println("******", i+1, "******")
		fmt.Println("Bucket Name is:", *buc.Name)
		fmt.Println("Bucket Creation Date is:", buc.CreationDate)
		fmt.Println("Region of bucket is:", aws.StringValue(reg.LocationConstraint))
	}
	return nil
}
