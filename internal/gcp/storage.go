package gcp

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/dustin/go-humanize"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

func ListBuckets(ctx context.Context) ([]StorageResp, error) {
	var bucketInfo []StorageResp

	client, err := storage.NewClient(ctx)
	if err != nil {
		return bucketInfo, err
	}
	project := ctx.Value(internal.KeyActiveProject).(string)
	it := client.Buckets(ctx, project)

	for {
		bucket, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			log.Printf("Error fetching bucket: %v", err)
			break
		}

		launchTime := bucket.Created
		localZone, err := config.GetLocalTimeZone() // Empty string loads the local timezone
		if err != nil {
			fmt.Println("Error loading local timezone:", err)
			return nil, err
		}
		loc, _ := time.LoadLocation(localZone)
		IST := launchTime.In(loc)
		storageResp := &StorageResp{BucketName: bucket.Name, CreationTime: IST.Format("Mon Jan _2 15:04:05 2006")}
		bucketInfo = append(bucketInfo, *storageResp)
	}

	return bucketInfo, nil
}

func GetInfoAboutBucket(ctx context.Context) []StorageObjResp {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}
	defer client.Close()
	bucketName := fmt.Sprintf("%v", ctx.Value(internal.BucketName))
	
	query := &storage.Query{Prefix: ""} // List all objects in the bucket.
	it := client.Bucket(bucketName).Objects(ctx, query)

	// Iterate through the objects in the bucket.
	objs := []StorageObjResp{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Print("Failed to iterate through objects: %v", err)
			break
		}

		obj := StorageObjResp{
			Name: attrs.Name,
		}
		log.Print(attrs)
		// Check if the object is a folder by examining its ContentType.
		if attrs.ContentType == "application/x-directory" {

			obj.ObjectType = "Folder"
			obj.Size = "-"
			obj.LastModified = attrs.Updated.String()
			obj.StorageClass = "-"
			fmt.Printf("Folder: %s\n", attrs.Prefix)
		} else {
			obj.ObjectType = "File"
			obj.Size = humanize.Bytes(uint64(attrs.Size))
			obj.LastModified = attrs.Updated.String()
			obj.StorageClass = "-"
			obj.SizeInBytes = attrs.Size
		}
		objs = append(objs, obj)
	}

	return objs
}
