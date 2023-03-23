package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type SQS struct {
	Accessor
	ctx context.Context
}

func (sqs *SQS) Init(ctx context.Context) {
	sqs.ctx = ctx
}

func (sqs *SQS) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	ins, err := aws.GetAllQueues(cfg)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}

func (sqs *SQS) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (sqs *SQS) Describe(queueUrl string) (string, error) {
	cfg, ok := sqs.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	res, _ := aws.GetMessageFromQueue(cfg, queueUrl)
	return fmt.Sprintf("%v", res), nil
}
