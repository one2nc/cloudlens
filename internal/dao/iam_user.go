package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/rs/zerolog/log"
)

type IAMU struct {
	Accessor
	ctx context.Context
}

func (iamu *IAMU) Init(ctx context.Context) {
	iamu.ctx = ctx
}

func (iamu *IAMU) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	usr := aws.GetUsers(*sess)
	objs := make([]Object, len(usr))
	for i, obj := range usr {
		objs[i] = obj
	}
	return objs, nil
}

func (iamu *IAMU) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
