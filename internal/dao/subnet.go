package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
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
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	vpcId := fmt.Sprintf("%v", ctx.Value(internal.VpcId))
	subnets := aws.GetSubnets(*sess, vpcId)
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
	sess, ok := sn.ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	res := aws.GetSingleSubnet(*sess, vpcId).GoString()
	return res, nil
}
