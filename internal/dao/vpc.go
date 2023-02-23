package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
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
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	vpcs := aws.GetVPCs(*sess)
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
	sess, ok := v.ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	res := aws.GetSingleVPC(*sess, vpcId).GoString()
	return res, nil
}
