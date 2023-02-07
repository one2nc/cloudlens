package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/rs/zerolog/log"
)

type EC2 struct {
	Accessor
	ctx context.Context
}

func (ec2 *EC2) Init(ctx context.Context) {
	ec2.ctx = ctx
}

func (e *EC2) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	ins, err := aws.GetInstances(*sess)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}

func (e *EC2) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
