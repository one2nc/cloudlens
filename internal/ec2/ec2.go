package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2Service interface {
	GetInstances()
}

type ec2Service struct {
	client ec2.EC2
}

func NewEc2Service(sess session.Session) Ec2Service {
	return ec2Service{client: *ec2.New(&sess)}
}

func (e ec2Service) GetInstances() {
	result, err := e.client.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error fetching instances:", err)
		return
	}

	// Iterate through the instances and print their ID and state
	for r, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			//fmt.Printf("%+v\n", instance)
			fmt.Println("****", r+1, "****")
			fmt.Println("Instance ID:", *instance.InstanceId)
			fmt.Println("Instance State:", *instance.State.Name)
			fmt.Println("Instance Type:", *instance.InstanceType)
			fmt.Println("Availability zone:", *instance.Placement.AvailabilityZone)
			fmt.Println("Monitoring: ", *instance.Monitoring.State)
			fmt.Println("Launch time: ", instance.LaunchTime.Format("2006-01-02 15:04:05"))
		}
	}
}
