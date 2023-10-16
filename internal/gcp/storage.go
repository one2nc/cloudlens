package gcp

import (
	"context"
	"fmt"
	"strings"
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
			log.Print("Error loading local timezone:", err)
			return nil, err
		}
		loc, _ := time.LoadLocation(localZone)
		IST := launchTime.In(loc)
		storageResp := &StorageResp{BucketName: bucket.Name, CreationTime: IST.Format("Mon Jan _2 15:04:05 2006")}
		bucketInfo = append(bucketInfo, *storageResp)
	}

	return bucketInfo, nil
}

func GetInfoAboutBucket(ctx context.Context) ([]StorageObjResp, error) {
	objs := []StorageObjResp{}
	client, err := storage.NewClient(ctx)
	if err != nil {
		
		return objs, err
	}
	defer client.Close()
	bucketName := fmt.Sprintf("%v", ctx.Value(internal.BucketName))
	fn := fmt.Sprintf("%v", ctx.Value(internal.FolderName))
	query := &storage.Query{Delimiter: "/", Prefix: fn} // List all objects in the bucket.
	it := client.Bucket(bucketName).Objects(ctx, query)

	// Iterate through the objects in the bucket.
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return objs,err
		}

		obj := StorageObjResp{}

		if attrs.Prefix != "" {
			obj.Name = attrs.Prefix
			obj.ObjectType = "Folder"
			obj.Size = "-"
			obj.LastModified = "-"
			obj.StorageClass = "-"
		} else {
			if attrs.Name == fn {
				continue
			}
			// remove folder name from file name
			splitName := strings.Split(attrs.Name, "/")
			obj.Name = splitName[len(splitName)-1]
			obj.ObjectType = "File"
			obj.Size = humanize.Bytes(uint64(attrs.Size))
			obj.StorageClass = attrs.StorageClass
			obj.SizeInBytes = attrs.Size

			launchTime := attrs.Updated
			localZone, err := config.GetLocalTimeZone() // Empty string loads the local timezone
			if err != nil {
				log.Print("Error loading local timezone:", err)
				continue
			}
			loc, _ := time.LoadLocation(localZone)
			IST := launchTime.In(loc)
			obj.LastModified = IST.Format("Mon Jan _2 15:04:05 2006")
		}

		// remove parent folder name from child folder name
		if strings.Contains(obj.Name, "/") {
			pathList := strings.Split(obj.Name, "/")
			if len(pathList) > 2 {
				obj.Name = pathList[len(pathList)-2] + "/"
			}
		}
		objs = append(objs, obj)
	}

	return objs,nil
}
