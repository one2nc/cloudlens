package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
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
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	userName := fmt.Sprintf("%v", ctx.Value(internal.UserName))
	usrPolicy := aws.GetPoliciesOfUser(*sess, userName)
	objs := make([]Object, len(usrPolicy))
	for i, obj := range usrPolicy {
		objs[i] = obj
	}
	return objs, nil
}

func (iamup *IAMUP) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
