package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type IAMUGP struct {
	Accessor
	ctx context.Context
}

func (iamugp *IAMUGP) Init(ctx context.Context) {
	iamugp.ctx = ctx
}

func (iamugp *IAMUGP) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	grpName := fmt.Sprintf("%v", ctx.Value(internal.GroupName))
	grpPolicy := aws.GetPoliciesOfGrp(cfg, grpName)
	objs := make([]Object, len(grpPolicy))
	for i, obj := range grpPolicy {
		objs[i] = obj
	}
	return objs, nil
}

func (iamugp *IAMUGP) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
