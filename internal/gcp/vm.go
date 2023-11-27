package gcp

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/one2nc/cloudlens/internal"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

func ListInstances(ctx context.Context) ([]VMResp, error) {
	var vmResp = []VMResp{}

	projectID := ctx.Value(internal.KeyActiveProject).(string)
	zone := ctx.Value(internal.KeyActiveZone).(string)
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return vmResp, err
	}
	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
	}

	it := instancesClient.List(ctx, req)

	for {
		instance, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return vmResp, err
		}

		zone := GetResourceFromURL(*instance.Zone)
		machineType := GetResourceFromURL(*instance.MachineType)
		createdAt, err := GetLocalTime(*instance.CreationTimestamp)
		if err != nil {
			log.Print(err)
			return vmResp, err
		}

		vm := VMResp{
			InstanceId:       *instance.Name,
			InstanceType:     machineType,
			AvailabilityZone: zone,
			InstanceState:    *instance.Status,
			LaunchTime:       createdAt,
		}
		vmResp = append(vmResp, vm)
	}
	return vmResp, nil
}

func GetInstance(ctx context.Context, instaceId string) (*computepb.Instance, error) {
	projectID := ctx.Value(internal.KeyActiveProject).(string)
	zone := ctx.Value(internal.KeyActiveZone).(string)
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer instancesClient.Close()

	instance, err := instancesClient.Get(ctx, &computepb.GetInstanceRequest{
		Zone:     zone,
		Instance: instaceId,
		Project:  projectID,
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}
