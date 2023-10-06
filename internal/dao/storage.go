package dao

import (
	"context"

	"github.com/one2nc/cloudlens/internal/gcp"
)

type Storage struct {
	Accessor
	ctx context.Context
}

func (s *Storage) Init(ctx context.Context) {
	s.ctx = ctx
}

func (s *Storage) List(ctx context.Context) ([]Object, error) {
	// cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	// if !ok {
	// 	log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	// }
	ins, err := gcp.ListBuckets()
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}

// func (e *EC2) Get(ctx context.Context, path string) (Object, error) {
// 	return nil, nil
// }

// func (e *EC2) Describe(instanceId string) (string, error) {
// 	cfg, ok := e.ctx.Value(internal.KeySession).(awsV2.Config)
// 	if !ok {
// 		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
// 	}
// 	res := aws.GetSingleInstance(cfg, instanceId)
// 	return fmt.Sprintf("%v", res), nil
// }


