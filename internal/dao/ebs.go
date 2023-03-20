package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type EBS struct {
	Accessor
	ctx context.Context
}

func (ebs *EBS) Init(ctx context.Context) {
	ebs.ctx = ctx
}

func (ebs *EBS) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	vols, err := aws.GetVolumes(cfg)
	if err != nil {
		return nil, err
	}
	objs := make([]Object, len(vols))
	for i, obj := range vols {
		objs[i] = obj
	}
	return objs, nil
}

func (ebs *EBS) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (ebs *EBS) Describe(volId string) (string, error) {
	cfg, ok := ebs.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	res := aws.GetSingleVolume(cfg, volId)
	return fmt.Sprintf("%v", res), nil
}
