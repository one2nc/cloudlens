package aws

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
)

type BucketResp struct {
	BucketName   string
	CreationTime string
	Region       string
}

type Presigner struct {
	PresignClient *s3.PresignClient
}

func ListBuckets(cfg aws.Config) ([]BucketResp, error) {
	var bucketInfo []BucketResp
	s3Client := s3.NewFromConfig(cfg)
	lbop, err := s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in listing buckets. err: %v", err))
		return nil, err
	}
	for _, buc := range lbop.Buckets {
		reg, err := s3Client.GetBucketLocation(context.Background(), &s3.GetBucketLocationInput{Bucket: buc.Name})
		if err != nil {
			log.Info().Msg(fmt.Sprintf("error getting bucket location. err: %v", err))
			return nil, err
		}
		launchTime := buc.CreationDate
		localZone, err := GetLocalTimeZone() // Empty string loads the local timezone
		if err != nil {
			fmt.Println("Error loading local timezone:", err)
			return nil, err
		}
		loc, _ := time.LoadLocation(localZone)
		IST := launchTime.In(loc)
		bucketresp := &BucketResp{BucketName: *buc.Name, CreationTime: IST.Format("Mon Jan _2 15:04:05 2006"), Region: string(reg.LocationConstraint)}
		bucketInfo = append(bucketInfo, *bucketresp)
	}
	return bucketInfo, nil
}

func GetInfoAboutBucket(cfg aws.Config, bucketName string, delimiter string, prefix string) (*s3.ListObjectsV2Output, error) {
	s3Serv := *s3.NewFromConfig(cfg)
	result, err := s3Serv.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket:    &bucketName,
		Delimiter: &delimiter,
		Prefix:    &prefix})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error is here: %v", err))
		return nil, err
	}
	return result, nil
}

func GetPreSignedUrl(cfg aws.Config, bucketName, key string) string {
	// log.Info().Msg("Bucket name izzzz:" + bucketName)
	// log.Info().Msg("Key name izzzz:" + key)
	// s3Serv := ss3.NewFromConfig(cfg)
	// req, _ := s3Serv.GetObjectRetention(context.Background(), &ss3.GetObjectAttributesInput{
	// 	Bucket: &bucketName,
	// 	Key:    &key,
	// })
	var presigner Presigner
	request, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(15 * time.Minute)
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, key, err)
	}
	log.Info().Msgf("Presigned URL is: %v", request.URL)
	return request.URL

}

func DownloadObject(cfg aws.Config, bucketName, key string) string {
	abc := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(abc)
	usr, err := user.Current()
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in getting the machine's user: %v", err))
	}
	path := usr.HomeDir + "/cloudlens/s3objects/"
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
	n, err := downloader.Download(context.Background(), f, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("failed to download file, err: %v", err))
		return ""
	}
	clipboard.WriteAll(p)

	return fmt.Sprintf("%v with size %d bytes, downloaded and its path copied to the clipboard", objectName, n)
}

func GetBuckEncryption(cfg aws.Config, bucketName string) *types.ServerSideEncryptionConfiguration {
	s3Serv := *s3.NewFromConfig(cfg)
	sse, _ := s3Serv.GetBucketEncryption(context.Background(), &s3.GetBucketEncryptionInput{
		Bucket: &bucketName,
	})
	//fmt.Println("sse string is :", sse.GoString())
	return sse.ServerSideEncryptionConfiguration
}

func GetBuckLifecycle(cfg aws.Config, bucketName string) *s3.GetBucketLifecycleConfigurationOutput {
	s3Serv := *s3.NewFromConfig(cfg)
	blc, _ := s3Serv.GetBucketLifecycleConfiguration(context.Background(), &s3.GetBucketLifecycleConfigurationInput{
		Bucket: &bucketName,
	})
	return blc
}
