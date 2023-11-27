package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/rs/zerolog/log"
)

type VMI struct {
	Accessor
	ctx context.Context
}

func (vmi *VMI) Init(ctx context.Context) {
	vmi.ctx = ctx
}

func (vmi *VMI) List(ctx context.Context) ([]Object, error) {

	ins, err := gcp.ListImages(ctx)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}

func (vmi *VMI) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (vms *VMI) Describe(imageId string) (string, error) {

	instance, err := gcp.GetImage(vms.ctx, imageId)

	if err != nil {
		log.Err(fmt.Errorf("Error while fetching image: %s", err.Error()))
		return "", err
	}
	res, err := json.MarshalIndent(instance, "", " ")

	if err != nil {
		log.Err(fmt.Errorf("Error while parsing json: %s", err.Error()))
		return "", err
	}
	return fmt.Sprintf("%v", string(res)), nil
}
