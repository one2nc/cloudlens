package dao

import (
	"context"
	"fmt"

	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/rs/zerolog/log"
)

type SG struct {
	Accessor
}

func (sg *SG) List(ctx context.Context) ([]Object, error) {
	log.Info().Msg(fmt.Sprintf("ctx type: %T", ctx.Value(internal.KeySession)))
	cfg, _ := config.Get()
	session, _ := config.GetSession(cfg.Profiles[0], "us-west-2", cfg.AwsConfig)
	//TODO: make dynamic
	//session := ctx.Value(internal.KeySession).(*session.Session)
	sgs := aws.GetSecGrps(*session)
	log.Info().Msg(fmt.Sprintf("ins: %d", len(sgs)))

	objs := make([]Object, len(sgs))
	for i, obj := range sgs {
		objs[i] = obj
	}
	return objs, nil
}

func (sg *SG) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
