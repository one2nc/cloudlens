package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type ECSServices struct {
	Accessor
	ctx context.Context
}

func (ecsServices *ECSServices) Init(ctx context.Context) {
	ecsServices.ctx = ctx
}

func (ecsServices *ECSServices) List(ctx context.Context) ([]Object, error) {
	var errMsg string
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		errMsg = fmt.Sprintf("conversion err: Expected awsV2.Config but got %v", cfg)
		log.Err(fmt.Errorf(errMsg))
	}
	clusterName, ok := ctx.Value(internal.ECSClusterName).(string)
	if !ok || clusterName == "" {
		errMsg = "failed to get ECS cluster name from context"
		log.Err(fmt.Errorf(errMsg))
		return nil, fmt.Errorf(errMsg)
	}
	listEcsServiceResp, err := aws.ListEcsServices(cfg, clusterName)
	if err != nil {
		errMsg = fmt.Sprintf("failed to list ECS services: %v", err)
		log.Err(fmt.Errorf(errMsg))
		return nil, fmt.Errorf(errMsg)
	}
	objs := make([]Object, len(listEcsServiceResp))
	for i, obj := range listEcsServiceResp {
		objs[i] = obj
	}
	return objs, err

}

func (ecsServices *ECSServices) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (ecsServices *ECSServices) Describe(serviceName string) (string, error) {
	var errMsg string
	cfg, ok := ecsServices.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		errMsg = fmt.Sprintf("conversion err: Expected awsV2.Config but got %v", cfg)
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	clusterName, ok := ecsServices.ctx.Value(internal.ECSClusterName).(string)
	if !ok || clusterName == "" {
		errMsg = "failed to get ECS cluster name from context"
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	res, err := aws.GetEcsServiceJSONResponse(cfg, clusterName, serviceName)
	if err != nil {
		errMsg = fmt.Sprintf("failed to get ECS service: %v", err)
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	return fmt.Sprintf("%v", res), nil
}
