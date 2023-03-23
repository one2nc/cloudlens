package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type IAMUP struct {
	Accessor
	ctx context.Context
}

func (iamup *IAMUP) Init(ctx context.Context) {
	iamup.ctx = ctx
}

func (iamup *IAMUP) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	userName := fmt.Sprintf("%v", ctx.Value(internal.UserName))
	usrPolicy := aws.GetPoliciesOfUser(cfg, userName)
	objs := make([]Object, len(usrPolicy))
	for i, obj := range usrPolicy {
		objs[i] = obj
	}
	return objs, nil
}

func (iamup *IAMUP) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
