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
	CreationTime string
	Region       string
}

type S3Service interface {
	ListBuckets(sess session.Session) ([]BucketResp, error)
}

func ListBuckets(sess session.Session) ([]BucketResp, error) {
	var bucketInfo []BucketResp
	s3Serv := *s3.New(&sess)
	lbop, err := s3Serv.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Error in listing buckets")
		return nil, err
	}
	for _, buc := range lbop.Buckets {
		reg, err := s3Serv.GetBucketLocationWithContext(context.Background(), &s3.GetBucketLocationInput{Bucket: aws.String(*buc.Name)})
		if err != nil {
			fmt.Println("error getting bucket location")
			return nil, err
		}
		launchTime := buc.CreationDate
		loc, _ := time.LoadLocation("Asia/Kolkata")
		IST := launchTime.In(loc)
		bucketresp := &BucketResp{BucketName: *buc.Name, CreationTime: IST.Format("Mon Jan _2 15:04:05 2006"), Region: aws.StringValue(reg.LocationConstraint)}
		bucketInfo = append(bucketInfo, *bucketresp)
	}
	return bucketInfo, nil
}
