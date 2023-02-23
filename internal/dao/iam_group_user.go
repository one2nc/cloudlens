package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
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
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	gp := fmt.Sprintf("%v", ctx.Value(internal.GroupName))
	groupUsers := aws.GetGroupUsers(*sess, gp)
	objs := make([]Object, len(groupUsers))
	for i, obj := range groupUsers {
		objs[i] = obj
	}
	return objs, nil
}

func (iamu *IamGroupUser) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
