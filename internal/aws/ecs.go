package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// --- ECS Clusters ---

func ListEcsClusters(cfg aws.Config) ([]EcsClusterResp, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	resultListClusters, err := ecsClient.ListClusters(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	var ecsClusterArns []string
	ecsClusterArns = append(ecsClusterArns, resultListClusters.ClusterArns...)
	describedClusters, err := DescribeEcsClusters(ecsClient, ecsClusterArns)
	if err != nil {
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
		return ecs.DescribeClustersOutput{}, err
	}
	return *detailedClusters, nil

}

func GetClusterJSONResponse(cfg aws.Config, clusterName string) (string, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	// Describe the specific cluster
	describeClustersInput := &ecs.DescribeClustersInput{
		Clusters: []string{clusterName},
	}
	result, err := ecsClient.DescribeClusters(context.TODO(), describeClustersInput)
	if err != nil {
		return "", err
	}
	// Check if the cluster was found
	if len(result.Clusters) == 0 {
		errMessage := fmt.Sprintf("Cluster %s not found", clusterName)
		return "", fmt.Errorf(errMessage)
	}
	// Marshal the cluster into a JSON string
	jsonResponse, err := json.MarshalIndent(result.Clusters[0], "", " ")
	if err != nil {
		return "", err
	}
	return string(jsonResponse), nil
}

// --- ECS Services ---

func ListEcsServices(cfg aws.Config, clusterName string) ([]EcsServiceResp, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	listServicesInput := &ecs.ListServicesInput{
		Cluster: &clusterName,
	}

	result, err := ecsClient.ListServices(context.TODO(), listServicesInput)
	if err != nil {
		return nil, err
	}
	var ecsServiceArns []string
	ecsServiceArns = append(ecsServiceArns, result.ServiceArns...)
	describedServices, err := DescribeEcsServices(ecsClient, clusterName, ecsServiceArns)
	if err != nil {
		return nil, err
	}

	var detailedServices []EcsServiceResp
	for _, service := range describedServices.Services {
		s := &EcsServiceResp{
			ServiceName:    *service.ServiceName,
			Status:         *service.Status,
			DesiredCount:   fmt.Sprint(service.DesiredCount),
			RunningCount:   fmt.Sprint(service.RunningCount),
			TaskDefinition: *service.TaskDefinition,
			ServiceArn:     *service.ServiceArn,
		}
		detailedServices = append(detailedServices, *s)
	}

	return detailedServices, nil
}

func DescribeEcsServices(ecsClient *ecs.Client, clusterName string, serviceArns []string) (ecs.DescribeServicesOutput, error) {
	describeServicesInput := &ecs.DescribeServicesInput{
		Cluster:  &clusterName,
		Services: serviceArns,
	}
	result, err := ecsClient.DescribeServices(context.TODO(), describeServicesInput)
	if err != nil {
		return ecs.DescribeServicesOutput{}, err
	}
	return *result, nil
}

func GetEcsServiceJSONResponse(cfg aws.Config, clusterName, serviceName string) (string, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	// Describe the specific service within the cluster
	describeServicesInput := &ecs.DescribeServicesInput{
		Cluster:  &clusterName,
		Services: []string{serviceName},
	}
	result, err := ecsClient.DescribeServices(context.TODO(), describeServicesInput)
	if err != nil {
		return "", err
	}
	// Check if the service was found
	if len(result.Services) == 0 {
		errMessage := fmt.Sprintf("Service %s not found in cluster %s", serviceName, clusterName)
		return "", fmt.Errorf(errMessage)
	}
	// Marshal the service into a JSON string
	jsonResponse, err := json.MarshalIndent(result.Services[0], "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonResponse), nil
}

// --- ECS Tasks ---

func ListEcsTasks(cfg aws.Config, clusterName, serviceName string) ([]EcsTaskResp, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	taskDetails, err := DescribeEcsTasksForService(ecsClient, clusterName, serviceName)
	if err != nil {
		return nil, err
	}
	tasks := make([]EcsTaskResp, len(taskDetails.Tasks))
	for i, task := range taskDetails.Tasks {
		t := &EcsTaskResp{
			TaskId: GetTaskIDFromArn(*task.TaskArn),
			Task:   &task,
		}
		tasks[i] = *t
	}
	return tasks, nil
}

func DescribeEcsTasksForService(ecsClient *ecs.Client, clusterName, serviceName string) (*ecs.DescribeTasksOutput, error) {
	listTasksInput := &ecs.ListTasksInput{
		Cluster:     &clusterName,
		ServiceName: &serviceName,
	}

	result, err := ecsClient.ListTasks(context.TODO(), listTasksInput)
	if err != nil {
		return nil, err
	}

	describeTasksInput := &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   result.TaskArns,
	}

	taskDetails, err := ecsClient.DescribeTasks(context.TODO(), describeTasksInput)
	if err != nil {
		return nil, err
	}

	return taskDetails, nil
}

func GetTaskJSONResponse(cfg aws.Config, clusterName, taskArn string) (string, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	describeTasksInput := &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   []string{GetTaskIDFromArn(taskArn)},
	}
	taskDetails, err := ecsClient.DescribeTasks(context.TODO(), describeTasksInput)
	if err != nil {
		return "", err
	}
	// Convert containers to JSON
	jsonResponse, err := json.MarshalIndent(taskDetails.Tasks[0], "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonResponse), nil
}

// --- ECS Containers ---

func ListContainersForTask(cfg aws.Config, clusterName, taskId string) ([]types.Container, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	describeTasksInput := &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   []string{taskId},
	}
	taskDetails, err := ecsClient.DescribeTasks(context.TODO(), describeTasksInput)
	if err != nil {
		return nil, err
	}
	if len(taskDetails.Tasks) == 0 {
		return nil, fmt.Errorf("task with ID %s not found", taskId)
	}
	return taskDetails.Tasks[0].Containers, nil
}

func GetECSContainerJsonResponse(cfg aws.Config, clusterName, taskId, runtimeId string) (string, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	describeTasksInput := &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   []string{taskId},
	}
	taskDetails, err := ecsClient.DescribeTasks(context.Background(), describeTasksInput)
	if err != nil {
		return "", err
	}
	for _, container := range taskDetails.Tasks[0].Containers {
		if *container.RuntimeId == runtimeId {
			jsonResponse, err := json.MarshalIndent(container, "", "  ")
			if err != nil {
				return "", err
			}
			return string(jsonResponse), nil
		}
	}
	return "", fmt.Errorf("container %s not found in task %s", runtimeId, taskId)
}

func GetTaskIDFromArn(taskArn string) string {
	parts := strings.Split(taskArn, "/")
	return parts[len(parts)-1]
}
