package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type ECSTasks struct {
	Accessor
	ctx context.Context
}

func (ecsTasks *ECSTasks) Init(ctx context.Context) {
	ecsTasks.ctx = ctx
}

func (ecsTasks *ECSTasks) List(ctx context.Context) ([]Object, error) {
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
	serviceName, ok := ctx.Value(internal.ECSServiceName).(string)
	if !ok || serviceName == "" {
		errMsg = "failed to get ECS service name from context"
		log.Err(fmt.Errorf(errMsg))
		return nil, fmt.Errorf(errMsg)
	}

	listEcsTasks, err := aws.ListEcsTasks(cfg, clusterName, serviceName)
	if err != nil {
		errMsg = fmt.Sprintf("failed to list ECS tasks: %v", err)
		log.Err(fmt.Errorf(errMsg))
		return nil, fmt.Errorf(errMsg)
	}
	objs := make([]Object, len(listEcsTasks))
	for i, obj := range listEcsTasks {
		objs[i] = obj
	}
	return objs, err
}

func (ecsTasks *ECSTasks) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (ecsTasks *ECSTasks) Describe(taskArn string) (string, error) {
	var errMsg string
	cfg, ok := ecsTasks.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		errMsg = fmt.Sprintf("conversion err: Expected awsV2.Config but got %v", cfg)
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	clusterName, ok := ecsTasks.ctx.Value(internal.ECSClusterName).(string)
	if !ok || clusterName == "" {
		errMsg = "failed to get ECS cluster name from context"
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	res, err := aws.GetTaskJSONResponse(cfg, clusterName, taskArn)
	if err != nil {
		errMsg = fmt.Sprintf("failed to get ECS service: %v", err)
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	return fmt.Sprintf("%v", res), nil
}
