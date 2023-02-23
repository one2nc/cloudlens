package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type Lambda struct {
	Accessor
	ctx context.Context
}

func (l *Lambda) Init(ctx context.Context) {
	l.ctx = ctx
}

func (l *Lambda) List(ctx context.Context) ([]Object, error) {
	sess, ok := ctx.Value(internal.KeySession).(*session.Session)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected session.session but got %v", sess))
	}
	ins, err := aws.GetAllLambdaFunctions(*sess)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}
