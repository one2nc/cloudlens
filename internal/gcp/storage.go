package gcp

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

type StorageResp struct {
	BucketName   string
	CreationTime string
	Region       string
}

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
