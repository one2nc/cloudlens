package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
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
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	sgs, err := aws.GetSecGrps(cfg)
	if err != nil {
		log.Info().Msg("Error in getting security groups: " + err.Error())
	}
	objs := make([]Object, len(sgs))
	for i, obj := range sgs {
		log.Info().Msgf("SEC GRP INFO: %v", obj)
		objs[i] = obj
	}
	return objs, nil
}

func (sg *SG) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (sg *SG) Describe(path string) (string, error) {
	cfg, ok := sg.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	sgInfo := aws.GetSingleSecGrp(cfg, path)
	return fmt.Sprintf("%v", sgInfo), nil
}
