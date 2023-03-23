package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type IamRole struct {
	Accessor
	ctx context.Context
}

func (iamu *IamRole) Init(ctx context.Context) {
	iamu.ctx = ctx
}

func (iamu *IamRole) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	usr := aws.GetIamRoles(cfg)
	objs := make([]Object, len(usr))
	for i, obj := range usr {
		objs[i] = obj
	}
	return objs, nil
}

func (iamu *IamRole) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
