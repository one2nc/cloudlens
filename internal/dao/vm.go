package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/rs/zerolog/log"
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

func (vm *VM) Describe(instanceId string) (string, error) {

	instance, err := gcp.GetInstance(vm.ctx, instanceId)

	if err != nil {
		log.Err(fmt.Errorf("Error while fetching instance: %s", err.Error()))
		return "", err
	}
	res, err := json.MarshalIndent(instance, "", " ")

	if err != nil {
		log.Err(fmt.Errorf("Error while parsing json: %s", err.Error()))
		return "", err
	}
	return fmt.Sprintf("%v", string(res)), nil
}

