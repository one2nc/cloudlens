package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type SG struct {
	Accessor
	ctx context.Context
}

func (sg *SG) Init(ctx context.Context) {
	sg.ctx = ctx
}

func (sg *SG) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	sgs := aws.GetSecGrps(*sess)
	objs := make([]Object, len(sgs))
	for i, obj := range sgs {
		objs[i] = obj
	}
	return objs, nil
}

func (sg *SG) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (sg *SG) Describe(path string) (string, error) {
	sess, ok := sg.ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	sgInfo := aws.GetSingleSecGrp(*sess, path)
	return sgInfo.GoString(), nil
}
