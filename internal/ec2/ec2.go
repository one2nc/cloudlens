package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Resp struct {
	InstanceId       string
	InstanceState    string
	InstanceType     string
	AvailabilityZone string
	MonitoringState  string
	LaunchTime       string
}
type Ec2Service interface {
	CreateInstances()
	GetInstances() ([]EC2Resp, error)
}

type ec2Service struct {
	client ec2.EC2
}

func NewEc2Service(sess session.Session) Ec2Service {
	return ec2Service{client: *ec2.New(&sess)}
}

func (e ec2Service) GetInstances() ([]EC2Resp, error) {
	var ec2Info []EC2Resp
	result, err := e.client.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error fetching instances:", err)
		return nil, err
	}

	// Iterate through the instances and print their ID and state
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			ec2Resp := &EC2Resp{InstanceId: *instance.InstanceId, InstanceState: *instance.State.Name, InstanceType: *instance.InstanceType, AvailabilityZone: *instance.Placement.AvailabilityZone, MonitoringState: *instance.Monitoring.State, LaunchTime: instance.LaunchTime.Format("2006-01-02 15:04:05")}
			ec2Info = append(ec2Info, *ec2Resp)
		}
	}
	return ec2Info, nil
}
func (e ec2Service) CreateInstances() {
	params := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-12345678"), // specify the ID of the image you want to use
		InstanceType: aws.String("t2.micro"),     // specify the instance type
		MinCount:     aws.Int64(3),
		MaxCount:     aws.Int64(3),
	}

	_, err := e.client.RunInstances(params)
	if err != nil {
		fmt.Println("Error in creating instances:", err)
	}
}
