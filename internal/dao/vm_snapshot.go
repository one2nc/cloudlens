package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/rs/zerolog/log"
)

type VMS struct {
	Accessor
	ctx context.Context
}

func (vms *VMS) Init(ctx context.Context) {
	vms.ctx = ctx
}

func (vms *VMS) List(ctx context.Context) ([]Object, error) {

	snapshots, err := gcp.ListSnapshots(ctx)
	objs := make([]Object, len(snapshots))
	for i, obj := range snapshots {
		objs[i] = obj
	}
	return objs, err
}

func (vms *VMS) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (vms *VMS) Describe(snapshotId string) (string, error) {

	instance, err := gcp.GetSnapshot(vms.ctx, snapshotId)

	if err != nil {
		log.Err(fmt.Errorf("Error while fetching snapshot: %s", err.Error()))
		return "", err
	}
	res, err := json.MarshalIndent(instance, "", " ")

	if err != nil {
		log.Err(fmt.Errorf("Error while parsing json: %s", err.Error()))
		return "", err
	}
	return fmt.Sprintf("%v", string(res)), nil
}
