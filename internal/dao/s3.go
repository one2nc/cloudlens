package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/rs/zerolog/log"
)

type S3 struct {
	Accessor
	ctx context.Context
}

func (s3 *S3) Init(ctx context.Context) {
	s3.ctx = ctx
}

func (s3 *S3) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	buckResp, err := aws.ListBuckets(*sess)
	objs := make([]Object, len(buckResp))
	for i, obj := range buckResp {
		objs[i] = obj
	}
	return objs, err
}

func (s3 *S3) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
