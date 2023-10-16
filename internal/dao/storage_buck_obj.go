package dao

import (
	"context"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/rs/zerolog/log"
)

type SBObj struct {
	Accessor
}

func (bo SBObj) List(ctx context.Context) ([]Object, error) {

	storageObjects, err := gcp.GetInfoAboutBucket(ctx)
	objs := make([]Object, len(storageObjects))
	if err != nil {
		log.Print("Error while listing objects: ", err.Error())
		return objs, err
	}
	for i, obj := range storageObjects {
		objs[i] = obj
	}
	return objs, nil
}

func (bo SBObj) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
