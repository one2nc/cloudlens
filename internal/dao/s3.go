package dao

import (
	"context"
	"fmt"

	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/rs/zerolog/log"
)

type S3 struct {
	Accessor
}

func (s3 *S3) List(context context.Context) ([]Object, error) {
	cfg, _ := config.Get()
	sess, _ := config.GetSession(cfg.Profiles[0], "ap-south-1", cfg.AwsConfig)
	buckResp, err := aws.ListBuckets(*sess)
	log.Info().Msg(fmt.Sprintf("s3 inssss: %d", len(buckResp)))
	objs := make([]Object, len(buckResp))
	for i, obj := range buckResp {
		objs[i] = obj
	}
	return objs, err
}

func (s3 *S3) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
