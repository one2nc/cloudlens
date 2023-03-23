package dao

import (
	"context"
	"encoding/json"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
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
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	buckResp, err := aws.ListBuckets(cfg)
	objs := make([]Object, len(buckResp))
	for i, obj := range buckResp {
		objs[i] = obj
	}
	return objs, err
}

func (s3 *S3) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (s3 *S3) Describe(BName string) (string, error) {
	cfg, ok := s3.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	be := aws.GetBuckEncryption(cfg, BName)
	blc := aws.GetBuckLifecycle(cfg, BName)
	log.Info().Msgf("be is: %v", be)
	log.Info().Msgf("blc is: %v", blc)
	return merge(*be, blc.Rules), nil
}

func merge(sse types.ServerSideEncryptionConfiguration, lcr []types.LifecycleRule) string {
	bi := aws.BucketInfo{EncryptionConfiguration: sse, LifeCycleRules: lcr}
	bij, _ := json.MarshalIndent(bi, "", " ")
	return string(bij)
}
