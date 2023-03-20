package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type EC2I struct {
	Accessor
	ctx context.Context
}

func (ei *EC2I) Init(ctx context.Context) {
	ei.ctx = ctx
}

func (ei *EC2I) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	ins := aws.GetAMIs(cfg)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, nil
}

func (ei *EC2I) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (ei *EC2I) Describe(imageId string) (string, error) {
	cfg, ok := ei.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	res := aws.GetSingleAMI(cfg, imageId)
	return fmt.Sprintf("%v", res), nil
}
