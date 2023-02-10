package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
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
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	rn := fmt.Sprintf("%v", ctx.Value(internal.RoleName))
	rp := aws.GetPoliciesOfRoles(*sess, rn)
	objs := make([]Object, len(rp))
	for i, obj := range rp {
		objs[i] = obj
	}
	return objs, nil
}

func (iamup *IamRolePloicy) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
