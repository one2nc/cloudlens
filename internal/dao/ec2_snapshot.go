package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type EC2S struct {
	Accessor
	ctx context.Context
}

func (es *EC2S) Init(ctx context.Context) {
	es.ctx = ctx
}

func (es *EC2S) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	ins := aws.GetSnapshots(cfg)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, nil
}

func (es *EC2S) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (es *EC2S) Describe(snapshotId string) (string, error) {
	cfg, ok := es.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	res := aws.GetSingleSnapshot(cfg, snapshotId)
	return fmt.Sprintf("%v", res), nil
}
