package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type ECSContainers struct {
	Accessor
	ctx context.Context
}

func (ecscn *ECSContainers) Init(ctx context.Context) {
	ecscn.ctx = ctx
}

func (ecscn *ECSContainers) List(ctx context.Context) ([]Object, error) {
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
	taskId, ok := ctx.Value(internal.ECSTaskId).(string)
	if !ok || taskId == "" {
		errMsg = "failed to get ECS taskId from context"
		log.Err(fmt.Errorf(errMsg))
		return nil, fmt.Errorf(errMsg)
	}

	listContainers, err := aws.ListContainersForTask(cfg, clusterName, taskId)
	if err != nil {
		errMsg = fmt.Sprintf("failed to list ECS containers: %v", err)
		log.Err(fmt.Errorf(errMsg))
		return nil, fmt.Errorf(errMsg)
	}
	objs := make([]Object, len(listContainers))
	for i, obj := range listContainers {
		objs[i] = obj
	}
	return objs, err
}

func (ecscn *ECSContainers) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (ecscn *ECSContainers) Describe(runtimeId string) (string, error) {
	var errMsg string
	cfg, ok := ecscn.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		errMsg = fmt.Sprintf("conversion err: Expected awsV2.Config but got %v", cfg)
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	clusterName, ok := ecscn.ctx.Value(internal.ECSClusterName).(string)
	if !ok || clusterName == "" {
		errMsg = "failed to get ECS cluster name from context"
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	taskId, ok := ecscn.ctx.Value(internal.ECSTaskId).(string)
	if !ok || taskId == "" {
		errMsg = "failed to get ECS taskId from context"
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	res, err := aws.GetECSContainerJsonResponse(cfg, clusterName, taskId, runtimeId)
	if err != nil {
		errMsg = fmt.Sprintf("failed to get continers: %v", err)
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	return fmt.Sprintf("%v", res), nil
}
