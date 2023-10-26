package dao

import (
	"context"

	"github.com/one2nc/cloudlens/internal/gcp"
)

type VM struct {
	Accessor
	ctx context.Context
}

func (vm *VM) Init(ctx context.Context) {
	vm.ctx = ctx
}

func (vm *VM) List(ctx context.Context) ([]Object, error) {

	ins, err := gcp.ListInstances(ctx)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}

// func (vm *VM) Get(ctx context.Context, path string) (Object, error) {
// 	return nil, nil
// }

// func (vm *VM) Describe(instanceId string) (string, error) {
// 	cfg, ok := vm.ctx.Value(internal.KeySession).(awsV2.Config)
// 	if !ok {
// 		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
// 	}
// 	res := aws.GetSingleInstance(cfg, instanceId)
// 	return fmt.Sprintf("%v", res), nil
// }
