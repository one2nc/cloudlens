package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/rs/zerolog/log"
)

func ListEcsClusters(cfg aws.Config) ([]EcsClusterResp, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	resultListClusters, err := ecsClient.ListClusters(context.TODO(), nil)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error fetching ECS Clusters: %v", err))
		return nil, err
	}
	var ecsClusterArns []string
	for _, cluster := range resultListClusters.ClusterArns {
		ecsClusterArns = append(ecsClusterArns, cluster)
	}
	describedClusters, err := DescribeEcsClusters(ecsClient, ecsClusterArns)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error describing ECS Clusters"))
		return nil, err
	}
	var detailedClusters []EcsClusterResp
	for _, cluster := range describedClusters.Clusters {
		c := &EcsClusterResp{ClusterName: *cluster.ClusterName, Status: *cluster.Status, RunningTasksCount: fmt.Sprint(cluster.RunningTasksCount), ClusterArn: *cluster.ClusterArn}
		detailedClusters = append(detailedClusters, *c)
	}
	return detailedClusters, nil

}

func DescribeEcsClusters(ecsClient *ecs.Client, clusters []string) (ecs.DescribeClustersOutput, error) {
	detailedClusters, err := ecsClient.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{Clusters: clusters})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error describing ECS Clusters"))
		return ecs.DescribeClustersOutput{}, err
	}
	return *detailedClusters, nil

}
