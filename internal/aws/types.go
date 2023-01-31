package aws

import (
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
