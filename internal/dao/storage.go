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
	ins, err := gcp.ListBuckets(ctx)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}
