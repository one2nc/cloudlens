package gcp

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

func ListInstances(ctx context.Context) ([]VMResp, error) {
	var vmResp = []VMResp{}
	projectID := "projectID"
	zone := "zone"

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

		vm := VMResp{
			InstanceId:       *instance.Name,
			InstanceType:     *instance.MachineType,
			AvailabilityZone: *instance.Zone,
			InstanceState:    *instance.Status,
			LaunchTime:       *instance.CreationTimestamp,
			
		}
		vmResp = append(vmResp, vm)
	}
	return vmResp, nil
}
