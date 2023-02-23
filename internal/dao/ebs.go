package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type EBS struct {
	Accessor
	ctx context.Context
}

func (ebs *EBS) Init(ctx context.Context) {
	ebs.ctx = ctx
}

func (ebs *EBS) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	vols, err := aws.GetVolumes(*sess)
	if err != nil {
		return nil, err
	}
	objs := make([]Object, len(vols))
	for i, obj := range vols {
		objs[i] = obj
	}
	return objs, nil
}

func (ebs *EBS) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (ebs *EBS) Describe(volId string) (string, error) {
	sess, ok := ebs.ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	res := aws.GetSingleVolume(*sess, volId)
	return res.GoString(), nil
}
