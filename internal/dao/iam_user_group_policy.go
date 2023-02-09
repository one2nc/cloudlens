package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
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
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	grpName := fmt.Sprintf("%v", ctx.Value(internal.GroupName))
	grpPolicy := aws.GetPoliciesOfGrp(*sess, grpName)
	objs := make([]Object, len(grpPolicy))
	for i, obj := range grpPolicy {
		objs[i] = obj
	}
	return objs, nil
}

func (iamugp *IAMUGP) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
