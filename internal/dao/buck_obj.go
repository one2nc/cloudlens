package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/dustin/go-humanize"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type BObj struct {
	Accessor
}

func (bo *BObj) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	bucketName := fmt.Sprintf("%v", ctx.Value(internal.BucketName))
	fn := fmt.Sprintf("%v", ctx.Value(internal.FolderName))
	var s3Objects []aws.S3Object
	bucketInfo, err := aws.GetInfoAboutBucket(cfg, bucketName, "/", fn)
	if err != nil {
		s3Objects = append(s3Objects, aws.S3Object{
			Name: "No objects found",
		})
	} else {
		s3Objects = setFoldersAndFiles(bucketInfo.CommonPrefixes, bucketInfo.Contents)
	}
	objs := make([]Object, len(s3Objects))
	for i, obj := range s3Objects {
		objs[i] = obj
	}
	return objs, nil
}

func (bo *BObj) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func getBucLevelInfo(bucketInfo *s3.ListObjectsV2Output) ([]string, []string) {
	var folderArrayInfo []string
	var fileArrayInfo []string
	for _, i := range bucketInfo.CommonPrefixes {
		folderArrayInfo = append(folderArrayInfo, *i.Prefix)
	}
	log.Info().Msg(fmt.Sprintf("%v", folderArrayInfo))

	for i := 0; i < len(bucketInfo.Contents); i++ {
		fileArrayInfo = append(fileArrayInfo, *bucketInfo.Contents[i].Key)
	}
	log.Info().Msg(fmt.Sprintf("%v", fileArrayInfo))
	return folderArrayInfo, fileArrayInfo
}

func setFoldersAndFiles(folders []types.CommonPrefix, files []types.Object) []aws.S3Object {
	var s3Objects []aws.S3Object
	indx := 0

	if len(folders) != 0 {
		for _, bi := range folders {
			keyA := strings.Split(*bi.Prefix, "/")
			o := aws.S3Object{
				Name:         keyA[len(keyA)-2],
				ObjectType:   "Folder",
				LastModified: "-",
				Size:         "-",
				StorageClass: "-",
			}
			s3Objects = append(s3Objects, o)
			indx++
		}
	}

	if len(files) != 0 {
		for _, fi := range files {
			localZone, err := GetLocalTimeZone() // Empty string loads the local timezone
			if err != nil {
				fmt.Println("Error loading local timezone:", err)
				return nil
			}
			loc, _ := time.LoadLocation(localZone)
			launchTime := fi.LastModified
			IST := launchTime.In(loc)
			keyA := strings.Split(*fi.Key, "/")
			if keyA[len(keyA)-1] != "" {
				o := aws.S3Object{
					Name:         keyA[len(keyA)-1],
					ObjectType:   "File",
					LastModified: IST.Format("Mon Jan _2 15:04:05 2006"),
					Size:         humanize.Bytes(uint64(fi.Size)),
					StorageClass: string(fi.StorageClass),
				}
				s3Objects = append(s3Objects, o)
				indx++
			}
		}
	}
	return s3Objects
}

func GetLocalTimeZone() (string, error) {
	localZone, err := time.LoadLocation("") // Empty string loads the local timezone
	if err != nil {
		fmt.Println("Error loading local timezone:", err)
		return "", err
	}
	return localZone.String(), nil
}
