package aws

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/zerolog/log"
)

type BucketResp struct {
	BucketName   string
	CreationTime string
	Region       string
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

func GetInfoAboutBucket(sess session.Session, bucketName string, delimiter string, prefix string) *s3.ListObjectsV2Output {
	s3Serv := *s3.New(&sess)
	result, err := s3Serv.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(bucketName),
		Delimiter: aws.String(delimiter),
		Prefix:    aws.String(prefix)})
	if err != nil {
		fmt.Println("Error is:", err)
		return nil
	}
	return result
}

func GetPreSignedUrl(sess session.Session, bucketName, key string) string {
	s3Serv := *s3.New(&sess)
	req, _ := s3Serv.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	url, _ := req.Presign(15 * time.Minute)
	return url
}

func DownloadObject(sess session.Session, bucketName, key string) string {
	downloader := s3manager.NewDownloader(&sess)
	usr, err := user.Current()
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in getting the machine's user: %v", err))
	}
	path := usr.HomeDir + "/cloud-lens/s3objects/"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in creating s3 Object directory: %v", err))
	}
	files := strings.Split(key, "/")
	objectName := files[len(files)-1]
	p := fmt.Sprintf("%v%v", path, objectName)
	log.Info().Msg(fmt.Sprintf("path: %v", p))
	f, err := os.Create(p)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Failed to create file, err: %v", err))
		return ""
	}
	defer f.Close()
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("failed to download file, err: %v", err))
		return ""
	}
	clipboard.WriteAll(p)

	return fmt.Sprintf("%v with size %d bytes, downloaded and its path copied to the clipboard", objectName, n)
}

func PutObjects(sess session.Session) {
	body := strings.NewReader("Hello, I'm working on aws cli!")
	s3Serv := *s3.New(&sess)
	_, err := s3Serv.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("test-bucket12948611666145821"),
		Key:    aws.String(""),
		Body:   body,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("uploaded object")
}

func GetBuckEncryption(sess session.Session, bucketName string) *s3.ServerSideEncryptionConfiguration {
	s3Serv := *s3.New(&sess)
	sse, _ := s3Serv.GetBucketEncryption(&s3.GetBucketEncryptionInput{
		Bucket: &bucketName,
	})
	//fmt.Println("sse string is :", sse.GoString())
	return sse.ServerSideEncryptionConfiguration
}

func GetBuckLifecycle(sess session.Session, bucketName string) *s3.GetBucketLifecycleConfigurationOutput {
	s3Serv := *s3.New(&sess)
	blc, _ := s3Serv.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
	})
	return blc
}
