package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type VPC struct {
	Accessor
	ctx context.Context
}

func (v *VPC) Init(ctx context.Context) {
	v.ctx = ctx
}

func (v *VPC) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	vpcs := aws.GetVPCs(cfg)
	objs := make([]Object, len(vpcs))
	for i, obj := range vpcs {
		objs[i] = obj
	}
	return objs, nil
}

func (v *VPC) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (v *VPC) Describe(vpcId string) (string, error) {
	cfg, ok := v.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	res := aws.GetSingleVPC(cfg, vpcId)
	return fmt.Sprintf("%v", res), nil
}
