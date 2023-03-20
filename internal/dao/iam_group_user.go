package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type IamGroupUser struct {
	Accessor
	ctx context.Context
}

func (igu *IamGroupUser) Init(ctx context.Context) {
	igu.ctx = ctx
}

func (igu *IamGroupUser) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	gp := fmt.Sprintf("%v", ctx.Value(internal.GroupName))
	groupUsers := aws.GetGroupUsers(cfg, gp)
	objs := make([]Object, len(groupUsers))
	for i, obj := range groupUsers {
		objs[i] = obj
	}
	return objs, nil
}

func (iamu *IamGroupUser) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
