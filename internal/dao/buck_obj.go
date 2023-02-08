package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dustin/go-humanize"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/rs/zerolog/log"
)

type BObj struct {
	Accessor
}

func (bo *BObj) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	bucketName := fmt.Sprintf("%v", ctx.Value(internal.BucketName))
	fn := fmt.Sprintf("%v", ctx.Value(internal.FolderName))
	log.Info().Msg(fmt.Sprintf("In Dao Bucket Name: %v", bucketName))
	log.Info().Msg(fmt.Sprintf("In Dao Folder Name: %v", fn))
	bucketInfo := aws.GetInfoAboutBucket(*sess, bucketName, "/", fn)
	folderArrayInfo, fileArrayInfo := getBucLevelInfo(bucketInfo)
	var s3Objects []aws.S3Object
	if len(folderArrayInfo) != 0 || len(fileArrayInfo) != 0 {
		s3Objects = setFoldersAndFiles(bucketInfo.CommonPrefixes, bucketInfo.Contents)
	} else {
		s3Objects = append(s3Objects, aws.S3Object{
			Name: "No objects found",
		})
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

func setFoldersAndFiles(Folder []*s3.CommonPrefix, File []*s3.Object) []aws.S3Object {
	var s3Objects []aws.S3Object
	indx := 0
	for _, bi := range Folder {
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

	for _, fi := range File {
		keyA := strings.Split(*fi.Key, "/")
		o := aws.S3Object{
			Name:         keyA[len(keyA)-1],
			ObjectType:   "File",
			LastModified: fi.LastModified.String(),
			Size:         humanize.Bytes(uint64(*fi.Size)),
			StorageClass: *fi.StorageClass,
		}
		s3Objects = append(s3Objects, o)
		indx++
	}

	return s3Objects
}
