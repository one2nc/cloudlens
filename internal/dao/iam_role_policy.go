package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type IamRolePloicy struct {
	Accessor
	ctx context.Context
}

func (irp *IamRolePloicy) Init(ctx context.Context) {
	irp.ctx = ctx
}

func (irp *IamRolePloicy) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	rn := fmt.Sprintf("%v", ctx.Value(internal.RoleName))
	rp := aws.GetPoliciesOfRoles(cfg, rn)
	objs := make([]Object, len(rp))
	for i, obj := range rp {
		objs[i] = obj
	}
	return objs, nil
}

func (iamup *IamRolePloicy) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
