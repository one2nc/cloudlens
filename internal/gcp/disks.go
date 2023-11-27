package gcp

import (
	"context"
	"fmt"

	"github.com/one2nc/cloudlens/internal"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/compute/v1"
)

func ListDisks(ctx context.Context) ([]DiskResp, error) {
	var diskResp = []DiskResp{}

	projectID := ctx.Value(internal.KeyActiveProject).(string)
	zone := ctx.Value(internal.KeyActiveZone).(string)
	client, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("Failed to create Compute Engine client: %v\n", err)
		return diskResp, err
	}

	disks, err := client.Disks.List(projectID, zone).Context(ctx).Do()
	if err != nil {
		log.Printf("Failed to list disks: %v\n", err)
		return diskResp, err
	}

	for _, zoneDisks := range disks.Items {

		createdAt, err := GetLocalTime(zoneDisks.CreationTimestamp)
		if err != nil {
			log.Print(err)
			return diskResp, err
		}
		diskType := GetResourceFromURL(zoneDisks.Type)
		disk := DiskResp{
			Name:         zoneDisks.Name,
			Status:       zoneDisks.Status,
			Type:         diskType,
			Size:         fmt.Sprintf("%v GB", zoneDisks.SizeGb),
			Zone:         zone,
			CreationTime: createdAt,
		}
		diskResp = append(diskResp, disk)

	}
	return diskResp, err
}

func GetDisk(ctx context.Context, diskId string) (*compute.Disk, error) {
	projectID := ctx.Value(internal.KeyActiveProject).(string)
	zone := ctx.Value(internal.KeyActiveZone).(string)
	client, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("Failed to create Compute Engine client: %v\n", err)
		return nil, err
	}

	disks, err := client.Disks.Get(projectID, zone, diskId).Context(ctx).Do()
	if err != nil {
		log.Printf("Failed to list disks: %v\n", err)
		return nil, err
	}

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return disks, nil
}
