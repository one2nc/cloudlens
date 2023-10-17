package gcp

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/atotto/clipboard"
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
			return objs, err
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

	return objs, nil
}

func DownloadObject(ctx context.Context, bucketName string, path string, fileName string) string {

	usr, err := user.Current()
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in getting the machine's user: %v", err))
	}
	dirPath := usr.HomeDir + "/cloudlens/gcp_storage_objects/"
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in creating GCP storage Object directory: %v", err))
		return ""
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Failed to create client: %v", err))
		return ""
	}
	defer client.Close()

	// Open the Google Cloud Storage object for reading.
	reader, err := client.Bucket(bucketName).Object(fmt.Sprintf("%v%v", path, fileName)).NewReader(ctx)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Failed to open object for reading: %v", err))
		return ""
	}
	defer reader.Close()

	locaFileName := fileName
	// Create a local file to save the downloaded object.
	localFilePath := fmt.Sprintf("%v%v", dirPath, locaFileName)
	localFile, err := os.Create(localFilePath)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Failed to create local file: %v", err))
		return ""
	}
	defer localFile.Close()

	// Copy the object content to the local file.
	if _, err := fmt.Fprint(localFile, reader); err != nil {
		log.Info().Msg(fmt.Sprintf("Failed to copy object content to local file: %v", err))
		return ""
	}
	clipboard.WriteAll(localFilePath)

	return fmt.Sprintf("%v with size %d bytes, downloaded and its path copied to the clipboard", fileName, reader.Attrs.Size)
}
