package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/rs/zerolog/log"
)

type Disk struct {
	Accessor
	ctx context.Context
}

func (disk *Disk) Init(ctx context.Context) {
	disk.ctx = ctx
}

func (disk *Disk) List(ctx context.Context) ([]Object, error) {

	disks, err := gcp.ListDisks(ctx)
	if err != nil {
		return nil, err
	}
	objs := make([]Object, len(disks))
	for i, obj := range disks {
		objs[i] = obj
	}
	return objs, nil
}

func (disk *Disk) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (disk *Disk) Describe(instanceId string) (string, error) {

	instance, err := gcp.GetDisk(disk.ctx, instanceId)

	if err != nil {
		log.Err(fmt.Errorf("Error while fetching instance: %s", err.Error()))
		return "", err
	}
	res, err := json.MarshalIndent(instance, "", " ")

	if err != nil {
		log.Err(fmt.Errorf("Error while parsing json: %s", err.Error()))
		return "", err
	}
	return fmt.Sprintf("%v", string(res)), nil
}
