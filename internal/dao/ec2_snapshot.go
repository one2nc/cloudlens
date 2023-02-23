package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type EC2S struct {
	Accessor
	ctx context.Context
}

func (es *EC2S) Init(ctx context.Context) {
	es.ctx = ctx
}

func (es *EC2S) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	ins := aws.GetSnapshots(*sess)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, nil
}

func (es *EC2S) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (es *EC2S) Describe(snapshotId string) (string, error) {
	sess, ok := es.ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	res := aws.GetSingleSnapshot(*sess, snapshotId).GoString()
	return res, nil
}
