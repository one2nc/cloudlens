package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type Subnet struct {
	Accessor
	ctx context.Context
}

func (sn *Subnet) Init(ctx context.Context) {
	sn.ctx = ctx
}

func (sn *Subnet) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	vpcId := fmt.Sprintf("%v", ctx.Value(internal.VpcId))
	subnets := aws.GetSubnets(cfg, vpcId)
	objs := make([]Object, len(subnets))
	for i, obj := range subnets {
		objs[i] = obj
	}
	return objs, nil
}

func (sn *Subnet) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (sn *Subnet) Describe(vpcId string) (string, error) {
	cfg, ok := sn.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	res := aws.GetSingleSubnet(cfg, vpcId)
	return fmt.Sprintf("%v", res), nil
}
