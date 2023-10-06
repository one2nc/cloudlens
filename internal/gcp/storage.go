package gcp

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

type StorageResp struct {
	BucketName   string
	CreationTime string
	Region       string
}

func ListBuckets() ([]StorageResp, error) {
	var bucketInfo []StorageResp

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: handle error.
	}
	it := client.Buckets(ctx, "centering-aegis-400910") //TODO: Pass project id through gcpconfig

	for {
		bucket, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			log.Printf("Error fetching bucket: %v", err)
			continue
		}

		storageResp := &StorageResp{BucketName: bucket.Name, CreationTime: bucket.Created.Local().String()}
		bucketInfo = append(bucketInfo, *storageResp)
	}

	return bucketInfo, nil
}
