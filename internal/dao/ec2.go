package dao

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
)

type EC2 struct {
	Accessor
}

func (e *EC2) List(ctx context.Context) ([]Object, error) {
	session := ctx.Value(internal.KeySession).(session.Session)
	ins, err := aws.GetInstances(session)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		ins[i] = obj
	}
	return objs, err
}

func (e *EC2) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
