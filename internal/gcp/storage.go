package gcp

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/one2nc/cloudlens/internal"
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

		storageResp := &StorageResp{BucketName: bucket.Name, CreationTime: bucket.Created.Local().String()}
		bucketInfo = append(bucketInfo, *storageResp)
	}

	return bucketInfo, nil
}
