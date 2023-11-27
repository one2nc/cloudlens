package gcp

import (
	"context"
	"fmt"

	"github.com/one2nc/cloudlens/internal"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/compute/v1"
)

func ListSnapshots(ctx context.Context) ([]SnapshotResp, error) {
	var snapshotsResp []SnapshotResp

	client, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("Failed to create Compute Engine client: %v\n", err)
		return snapshotsResp, err
	}
	project := ctx.Value(internal.KeyActiveProject).(string)
	snapshots, err := client.Snapshots.List(project).Context(ctx).Do()
	if err != nil {
		log.Printf("Failed to fetch snapshots: %v\n", err)
		return snapshotsResp, err
	}
	for _, snapshot := range snapshots.Items {
		createdAt, err := GetLocalTime(snapshot.CreationTimestamp)
		if err != nil {
			log.Print(err)
			return snapshotsResp, err
		}
		res := SnapshotResp{
			Name:      snapshot.Name,
			Size:      fmt.Sprintf("%v GB", snapshot.DiskSizeGb),
			CreatedAt: createdAt,
		}
		snapshotsResp = append(snapshotsResp, res)
	}

	return snapshotsResp, nil
}
func GetSnapshot(ctx context.Context, snapshotID string) (*compute.Snapshot, error) {

	client, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("Failed to create Compute Engine client: %v\n", err)
		return nil, err
	}
	project := ctx.Value(internal.KeyActiveProject).(string)
	snapshot, err := client.Snapshots.Get(project, snapshotID).Context(ctx).Do()

	if err != nil {
		log.Printf("Failed to fetch snapshots: %v\n", err)
		return nil, err
	}

	return snapshot, nil
}
