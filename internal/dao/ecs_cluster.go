package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type ECSClusters struct {
	Accessor
	ctx context.Context
}

func (ecsClusters *ECSClusters) Init(ctx context.Context) {
	ecsClusters.ctx = ctx
}

func (ecsClusters *ECSClusters) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	listClustersResp, err := aws.ListEcsClusters(cfg)
	if err != nil {
		log.Err(fmt.Errorf("failed to list ECS clusters: %v", err))
		return nil, err
	}
	objs := make([]Object, len(listClustersResp))
	for i, obj := range listClustersResp {
		objs[i] = obj
	}
	return objs, err
}

func (ecsClusters *ECSClusters) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (ecsClusters *ECSClusters) Describe(clusterName string) (string, error) {
	var errMsg string
	cfg, ok := ecsClusters.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		errMsg = fmt.Sprintf("conversion err: Expected awsV2.Config but got %v", cfg)
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	res, err := aws.GetClusterJSONResponse(cfg, clusterName)
	if err != nil {
		errMsg = fmt.Sprintf("failed to get ECS cluster: %v", err)
		log.Err(fmt.Errorf(errMsg))
		return "", fmt.Errorf(errMsg)
	}
	return fmt.Sprintf("%v", res), nil
}
