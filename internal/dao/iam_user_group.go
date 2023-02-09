package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/rs/zerolog/log"
)

type IAMUG struct {
	Accessor
	ctx context.Context
}

func (iamug *IAMUG) Init(ctx context.Context) {
	iamug.ctx = ctx
}

func (iamug *IAMUG) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	usrGroup := aws.GetUserGroups(*sess)
	objs := make([]Object, len(usrGroup))
	for i, obj := range usrGroup {
		objs[i] = obj
	}
	return objs, nil
}

func (iamug *IAMUG) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
