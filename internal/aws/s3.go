package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type BucketResp struct {
	BucketName   string
	CreationTime *time.Time
	Region       string
}

type S3Service interface {
	CreateBucket(bucketName string) error
	ListBuckets() ([]BucketResp, error)
}

type s3Service struct {
	client s3.S3
}

func NewS3Service(sess session.Session) S3Service {
	// Create an Amazon S3 service client
	return s3Service{client: *s3.New(&sess)}
}

func (s s3Service) CreateBucket(bucketName string) error {
	opcb, err := s.client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: aws.String("ap-south-2")}})
	if err != nil {
		fmt.Println("Error occured while creating a bucket:", err.Error())
		return err
	}
	fmt.Println("The location of bucket is:", *opcb.Location)
	return nil
}

func (s s3Service) ListBuckets() ([]BucketResp, error) {
	var bucketInfo []BucketResp
	lbop, err := s.client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Error in listing buckets")
		return nil, err
	}
	for _, buc := range lbop.Buckets {
		reg, err := s.client.GetBucketLocationWithContext(context.Background(), &s3.GetBucketLocationInput{Bucket: aws.String(*buc.Name)})
		if err != nil {
			fmt.Println("error getting bucket location")
			return nil, err
		}
		bucketresp := &BucketResp{BucketName: *buc.Name, CreationTime: buc.CreationDate, Region: aws.StringValue(reg.LocationConstraint)}
		bucketInfo = append(bucketInfo, *bucketresp)
	}
	return bucketInfo, nil
}
