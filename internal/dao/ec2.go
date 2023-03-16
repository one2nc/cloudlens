package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type EC2 struct {
	Accessor
	ctx context.Context
}

func (ec2 *EC2) Init(ctx context.Context) {
	ec2.ctx = ctx
}

//	func (e *EC2) List(ctx context.Context) ([]Object, error) {
//		sess, ok := ctx.Value(internal.KeySession).(*session.Session)
//		if !ok {
//			log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
//		}
//		ins, err := aws.GetInstances(*sess)
//		objs := make([]Object, len(ins))
//		for i, obj := range ins {
//			objs[i] = obj
//		}
//		return objs, err
//	}

func (e *EC2) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	ins, err := aws.GetInstances(cfg)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}

func (e *EC2) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (e *EC2) Describe(instanceId string) (string, error) {
	sess, ok := e.ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	res := aws.GetSingleInstance(*sess, instanceId).GoString()
	return res, nil
}
