package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Resp struct {
	Instance         ec2.Instance
	InstanceId       string
	InstanceType     string
	AvailabilityZone string
	InstanceState    string
	PublicDNS        string
	MonitoringState  string
	LaunchTime       string
}

type S3Object struct {
	SizeInBytes int64
	Name, ObjectType, LastModified, Size, StorageClass string
}

type BucketInfo struct {
	EncryptionConfiguration types.ServerSideEncryptionConfiguration
	LifeCycleRules          []types.LifecycleRule
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

type SQSResp struct {
	Name              string
	URL               string
	Type              string
	Created           string
	MessagesAvailable string
	Encryption        string
	MaxMessageSize    string
}

type Snapshot struct {
	SnapshotId string
	OwnerId    string
	VolumeId   string
	VolumeSize string
	StartTime  string
	State      string
}

type ImageResp struct {
	ImageId       string
	OwnerId       string
	ImageLocation string
	Name          string
	ImageType     string
}

type VpcResp struct {
	VpcId           string
	OwnerId         string
	CidrBlock       string
	InstanceTenancy string
	State           string
}

type LambdaResp struct {
	FunctionName string
	Description  string
	Role         string
	FunctionArn  string
	CodeSize     string
	LastModified string
}

type SubnetResp struct {
	SubnetId         string
	OwnerId          string
	CidrBlock        string
	AvailabilityZone string
	State            string
}

type SGResp struct {
	GroupId     string
	GroupName   string
	Description string
	OwnerId     string
	VpcId       string
}

type EcsClusterResp struct {
	ClusterName       string
	Status            string
	ClusterArn        string
	RunningTasksCount string
}
