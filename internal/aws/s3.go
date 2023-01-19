package aws

import (
	"context"
	"fmt"
	"strings"
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

type ByBucketName []BucketResp

func (b ByBucketName) Len() int           { return len(b) }
func (b ByBucketName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByBucketName) Less(i, j int) bool { return b[i].BucketName < b[j].BucketName }

type S3Service interface {
	ListBuckets(sess session.Session) ([]BucketResp, error)
	GetInfoAboutBucket(sess session.Session)
	PutObjects(sess session.Session)
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

func GetInfoAboutBucket(sess session.Session, bucketName string) *s3.ListObjectsOutput {
	s3Serv := *s3.New(&sess)
	result, err := s3Serv.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucketName), Delimiter: aws.String("/"), Prefix: aws.String("")})
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}
	// fmt.Println(result)
	return result
}

func PutObjects(sess session.Session) {
	body := strings.NewReader("Hello, I'm working on aws cli!")
	s3Serv := *s3.New(&sess)
	_, err := s3Serv.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("test-bucket577006791947779410"),
		Key:    aws.String("folder1/folder2/key1.png"),
		Body:   body,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("uploaded object")
}
