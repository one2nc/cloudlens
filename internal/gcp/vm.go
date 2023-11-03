package gcp

import (
	"context"
	"strings"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/config"
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

		splittedZoneURL := strings.Split(*instance.Zone, "/")
		zone := splittedZoneURL[len(splittedZoneURL)-1]

		splittedMachineURL := strings.Split(*instance.MachineType, "/")
		machineType := splittedMachineURL[len(splittedMachineURL)-1]

		createdTime := *instance.CreationTimestamp
		launchTime, err := time.Parse("2006-01-02T15:04:05.999-07:00", createdTime)
		if err != nil {
			log.Print("Error parsing timestamp :", err)
			return vmResp, err
		}
		localZone, err := config.GetLocalTimeZone()
		if err != nil {
			log.Print("Error loading local timezone:", err)
			return vmResp, err
		}
		loc, _ := time.LoadLocation(localZone)
		IST := launchTime.In(loc)

		vm := VMResp{
			InstanceId:       *instance.Name,
			InstanceType:     machineType,
			AvailabilityZone: zone,
			InstanceState:    *instance.Status,
			LaunchTime:       IST.Format("Mon Jan _2 15:04:05 2006"),
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
