package dao

import (
	"context"
	"fmt"

	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/rs/zerolog/log"
)

type EC2 struct {
	Accessor
}

func (e *EC2) List(ctx context.Context) ([]Object, error) {
	log.Info().Msg(fmt.Sprintf("ctx type: %T", ctx.Value(internal.KeySession)))
	cfg, _ := config.Get()
	session, _ := config.GetSession(cfg.Profiles[0], "ap-south-1", cfg.AwsConfig)
	//TODO: make dynamic
	//session := ctx.Value(internal.KeySession).(*session.Session)
	ins, err := aws.GetInstances(*session)
	log.Info().Msg(fmt.Sprintf("ins: %d", len(ins)))

	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}

func (e *EC2) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
