package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
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

type S3Object struct {
	Name, ObjectType, LastModified, Size, StorageClass string
}

type BucketInfo struct {
	EncryptionConfiguration *s3.ServerSideEncryptionConfiguration
	LifeCycleRules          []*s3.LifecycleRule
}
