package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/rs/zerolog/log"
)

type EC2I struct {
	Accessor
	ctx context.Context
}

func (ei *EC2I) Init(ctx context.Context) {
	ei.ctx = ctx
}

func (ei *EC2I) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	ins := aws.GetAMIs(*sess)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, nil
}

func (ei *EC2I) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (ei *EC2I) Describe(imageId string) (string, error) {
	sess, ok := ei.ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	res := aws.GetSingleAMI(*sess, imageId).GoString()
	return res, nil
}