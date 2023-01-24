package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Resp struct {
	Instance         ec2.Instance
	InstanceId       string
	InstanceType     string
	AvailabilityZone string
	InstanceState    string
	PublicDNS        string
	PublicIPv4       string
	MonitoringState  string
	LaunchTime       string
}

type ByInstanceId []EC2Resp

func (e ByInstanceId) Len() int           { return len(e) }
func (e ByInstanceId) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ByInstanceId) Less(i, j int) bool { return e[i].InstanceId < e[j].InstanceId }

type ByInstanceType []EC2Resp

func (e ByInstanceType) Len() int           { return len(e) }
func (e ByInstanceType) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ByInstanceType) Less(i, j int) bool { return e[i].InstanceType < e[j].InstanceType }

type ByLaunchTime []EC2Resp

func (e ByLaunchTime) Len() int      { return len(e) }
func (e ByLaunchTime) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e ByLaunchTime) Less(i, j int) bool {
	return e[i].Instance.LaunchTime.Before(*e[j].Instance.LaunchTime)
}

type Ec2Service interface {
	GetInstances(sess session.Session) ([]EC2Resp, error)
	GetSingleInstance() *ec2.DescribeInstancesOutput
}

func GetAllRegions() []string {
	return []string{"ap-south-1", "us-east-2", "us-west-1", "us-west-2", "af-south-1", "ap-east-1", "ap-south-2", "ap-southeast-3", "us-east-1", "ap-northeast-3", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-south-2", "eu-north-1", "eu-central-2", "me-south-1", "me-central-1", "sa-east-1", "us-gov-east-1", "us-gov-west-1"}
}

func GetInstances(sess session.Session) ([]EC2Resp, error) {
	var ec2Info []EC2Resp
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error fetching instances:", err)
		return nil, err
	}
	// Iterate through the instances and print their ID and state
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			launchTime := instance.LaunchTime
			loc, _ := time.LoadLocation("Asia/Kolkata")
			IST := launchTime.In(loc)
			ec2Resp := &EC2Resp{
				Instance:         *instance,
				InstanceId:       *instance.InstanceId,
				InstanceType:     *instance.InstanceType,
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				InstanceState:    *instance.State.Name,
				PublicDNS:        *instance.PrivateDnsName,
				PublicIPv4:       *instance.PrivateIpAddress,
				MonitoringState:  *instance.Monitoring.State,
				LaunchTime:       IST.Format("Mon Jan _2 15:04:05 2006")}
			ec2Info = append(ec2Info, *ec2Resp)
		}
	}
	return ec2Info, nil
}

func GetSingleInstance(sess session.Session, insId string) *ec2.DescribeInstancesOutput {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{&insId},
	})
	if err != nil {
		fmt.Println("Error fetching instance with id: ", insId, " err: ", err)
		return nil
	}
	return result
}
