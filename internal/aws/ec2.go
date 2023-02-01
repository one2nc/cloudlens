package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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

func GetAllRegions() []string {
	return []string{
		"us-east-1", "us-east-2", "us-west-1", "us-west-2", "af-south-1", "ap-east-1",
		"ap-southeast-3", "ap-south-1", "ap-northeast-3", "ap-northeast-2",
		"ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1",
		"eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-north-1",
		"me-south-1", "me-central-1", "sa-east-1", "us-gov-east-1", "us-gov-west-1"}
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

func GetSecGrps(sess session.Session) []*ec2.SecurityGroup {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		fmt.Println("Error in fetching Security Groups: ", " err: ", err)
		return nil
	}
	return result.SecurityGroups
}

func GetSingleSecGrp(sess session.Session, sgId string) *ec2.DescribeSecurityGroupsOutput {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{&sgId},
	})
	if err != nil {
		fmt.Println("Error in fetching Security Group: ", sgId, " err: ", err)
		return nil
	}
	return result
}

/*
Volumes(ebs) are region specific
Localstack doesn't have default volumes, so at some regions, there won't be any volumes.
*/
func GetVolumes(sess session.Session) []*ec2.Volume {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeVolumes(&ec2.DescribeVolumesInput{})
	if err != nil {
		fmt.Println("Error in fetching Volumes: ", " err: ", err)
		return nil
	}
	return result.Volumes
}

func GetSingleVolume(sess session.Session, vId string) *ec2.Volume {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeVolumes(&ec2.DescribeVolumesInput{
		VolumeIds: []*string{&vId},
	})
	if err != nil {
		fmt.Println("Error in fetching Volume: ", vId, " err: ", err)
		return nil
	}
	return result.Volumes[0]
}

/*
Snapshots are region specific
Localstack does have default snapshots, so we can see some of the snapshots that we never created
*/
func GetSnapshots(sess session.Session) []*ec2.Snapshot {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeSnapshots(&ec2.DescribeSnapshotsInput{})
	if err != nil {
		fmt.Println("Error in fetching Snapshots: ", " err: ", err)
		return nil
	}
	return result.Snapshots
}

func GetSingleSnapshot(sess session.Session, sId string) *ec2.Snapshot {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
		SnapshotIds: []*string{&sId},
	})
	if err != nil {
		fmt.Println("Error in fetching Snapshot: ", sId, " err: ", err)
		return nil
	}
	return result.Snapshots[0]
}

/*
	AMIs are region specific
	Localstack does have default some AMIs, so we can see some of the AMIs that we never created
*/

func GetAMIs(sess session.Session) []*ec2.Image {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeImages(&ec2.DescribeImagesInput{})
	if err != nil {
		fmt.Println("Error in fetching AMIs: ", " err: ", err)
		return nil
	}
	return result.Images
}

func GetSingleAMI(sess session.Session, amiId string) *ec2.Image {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeImages(&ec2.DescribeImagesInput{
		ImageIds: []*string{&amiId},
	})
	if err != nil {
		fmt.Println("Error in fetching AMI: ", amiId, " err: ", err)
		return nil
	}
	return result.Images[0]
}

func GetVPCs(sess session.Session) []*ec2.Vpc {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		fmt.Println("Error in fetching VPCs: ", " err: ", err)
		return nil
	}
	return result.Vpcs
}

func GetSingleVPC(sess session.Session, vpcId string) *ec2.Vpc {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{&vpcId},
	})
	if err != nil {
		fmt.Println("Error in fetching VPC: ", vpcId, " err: ", err)
		return nil
	}
	return result.Vpcs[0]
}

func GetSubnets(sess session.Session) []*ec2.Subnet {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		fmt.Println("Error in fetching Subnets: ", " err: ", err)
		return nil
	}
	return result.Subnets
}

func GetSingleSubnet(sess session.Session, sId string) *ec2.Subnet {
	ec2Serv := *ec2.New(&sess)
	result, err := ec2Serv.DescribeSubnets(&ec2.DescribeSubnetsInput{
		SubnetIds: []*string{&sId},
	})
	if err != nil {
		fmt.Println("Error in fetching Subnet: ", sId, " err: ", err)
		return nil
	}
	return result.Subnets[0]
}
