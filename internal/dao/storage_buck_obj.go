package dao

import (
	"context"

	"github.com/one2nc/cloudlens/internal/gcp"
)

type SBObj struct {
	Accessor
}

func (bo SBObj) List(ctx context.Context) ([]Object, error) {

	storageObjects := gcp.GetInfoAboutBucket(ctx)

	objs := make([]Object, len(storageObjects))
	for i, obj := range storageObjects {
		objs[i] = obj
	}
	return objs, nil
}

func (bo SBObj) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
