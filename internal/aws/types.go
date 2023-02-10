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

type IAMUSerResp struct {
	UserId       string
	UserName     string
	ARN          string
	CreationTime string
}

type IAMUSerGroupResp struct {
	GroupId      string
	GroupName    string
	ARN          string
	CreationTime string
}

type IAMUSerPolicyResponse struct {
	PolicyArn  string
	PolicyName string
}

type EBSResp struct {
	VolumeId         string
	Size             string
	VolumeType       string
	State            string
	AvailabilityZone string
	Snapshot         string
	CreationTime     string
}

type IAMUSerGroupPolicyResponse struct {
	PolicyArn  string
	PolicyName string
}

type IamRoleResp struct {
	RoleId       string
	RoleName     string
	ARN          string
	CreationTime string
}

type IamRolePolicyResponse struct {
	PolicyArn  string
	PolicyName string
}
