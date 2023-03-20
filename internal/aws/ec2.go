package aws

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	ecc "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rs/zerolog/log"
)

//	func GetInstances(sess session.Session) ([]EC2Resp, error) {
//		var ec2Info []EC2Resp
//		ec2Serv := *ec2.New(&sess)
//		result, err := ec2Serv.DescribeInstances(nil)
//		if err != nil {
//			log.Info().Msg(fmt.Sprintf("Error fetching instances: %v", err))
//			return nil, err
//		}
//		// Iterate through the instances and print their ID and state
//		for _, reservation := range result.Reservations {
//			for _, instance := range reservation.Instances {
//				launchTime := instance.LaunchTime
//				localZone, err := GetLocalTimeZone() // Empty string loads the local timezone
//				if err != nil {
//					fmt.Println("Error loading local timezone:", err)
//					return nil, err
//				}
//				loc, _ := time.LoadLocation(localZone)
//				IST := launchTime.In(loc)
//				ec2Resp := &EC2Resp{
//					Instance:         *instance,
//					InstanceId:       *instance.InstanceId,
//					InstanceType:     *instance.InstanceType,
//					AvailabilityZone: *instance.Placement.AvailabilityZone,
//					InstanceState:    *instance.State.Name,
//					PublicDNS:        *instance.PublicDnsName,
//					MonitoringState:  *instance.Monitoring.State,
//					LaunchTime:       IST.Format("Mon Jan _2 15:04:05 2006")}
//				ec2Info = append(ec2Info, *ec2Resp)
//			}
//		}
//		return ec2Info, nil
//	}
func GetInstances(cfg awsV2.Config) ([]EC2Resp, error) {
	var ec2Info []EC2Resp
	ec2Client := ecc.NewFromConfig(cfg)
	resultec2, err := ec2Client.DescribeInstances(context.TODO(), nil)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error fetching instances: %v", err))
		return nil, err
	}

	// Iterate through the instances and print their ID and state
	for _, reservation := range resultec2.Reservations {
		for _, instance := range reservation.Instances {
			launchTime := instance.LaunchTime
			localZone, err := GetLocalTimeZone() // Empty string loads the local timezone
			if err != nil {
				fmt.Println("Error loading local timezone:", err)
				return nil, err
			}
			loc, _ := time.LoadLocation(localZone)
			IST := launchTime.In(loc)
			ec2Resp := &EC2Resp{
				Instance:         ec2.Instance{},
				InstanceId:       *instance.InstanceId,
				InstanceType:     string(instance.InstanceType),
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				InstanceState:    string(instance.State.Name),
				PublicDNS:        *instance.PublicDnsName,
				MonitoringState:  string(instance.Monitoring.State),
				LaunchTime:       IST.Format("Mon Jan _2 15:04:05 2006")}
			ec2Info = append(ec2Info, *ec2Resp)
		}
	}
	return ec2Info, nil
}

func GetSingleInstance(cfg awsV2.Config, insId string) *ecc.DescribeInstancesOutput {
	ec2Client := ecc.NewFromConfig(cfg)
	result, err := ec2Client.DescribeInstances(context.Background(), &ecc.DescribeInstancesInput{
		InstanceIds: []string{insId},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error fetching instance with id: %s, err: %v", insId, err))
		return nil
	}
	return result
}

func GetSecGrps(cfg awsV2.Config) ([]types.SecurityGroup, error) {
	ec2Client := ecc.NewFromConfig(cfg)
	result, err := ec2Client.DescribeSecurityGroups(context.Background(), &ecc.DescribeSecurityGroupsInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Security Groups. err: %v ", err))
		return nil, err
	}
	for _, sg := range result.SecurityGroups {
		log.Info().Msgf("Security Group ID: %s, Name: %s\n", *sg.GroupId, *sg.GroupName)
	}
	return result.SecurityGroups, nil
}

func GetSingleSecGrp(cfg awsV2.Config, sgId string) *ecc.DescribeSecurityGroupsOutput {
	ec2Serv := *ecc.NewFromConfig(cfg)
	result, err := ec2Serv.DescribeSecurityGroups(context.Background(), &ecc.DescribeSecurityGroupsInput{
		GroupIds: []string{sgId},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Security Group: %s err: %v ", sgId, err))
		return nil
	}
	return result
}

/*
Volumes(ebs) are region specific
Localstack doesn't have default volumes, so at some regions, there won't be any volumes.
*/
func GetVolumes(cfg awsV2.Config) ([]EBSResp, error) {
	var volumes []EBSResp
	ec2Client := ecc.NewFromConfig(cfg)
	result, err := ec2Client.DescribeVolumes(context.Background(), &ecc.DescribeVolumesInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Volumes. err: %v", err))
		return nil, err
	}
	for _, v := range result.Volumes {
		launchTime := v.CreateTime
		localZone, err := GetLocalTimeZone() // Empty string loads the local timezone
		if err != nil {
			fmt.Println("Error loading local timezone:", err)
			return nil, err
		}
		loc, _ := time.LoadLocation(localZone)
		IST := launchTime.In(loc)
		IST.Format("Mon Jan _2 15:04:05 2006")
		volume := EBSResp{
			VolumeId:         *v.VolumeId,
			Size:             strconv.Itoa(int(*v.Size)) + " GB",
			VolumeType:       string(v.VolumeType),
			State:            string(v.State),
			AvailabilityZone: *v.AvailabilityZone,
			Snapshot:         *v.SnapshotId,
			CreationTime:     IST.String(),
		}
		volumes = append(volumes, volume)
	}
	return volumes, nil
}

func GetSingleVolume(cfg awsV2.Config, vId string) types.Volume {
	ec2Client := ecc.NewFromConfig(cfg)
	result, err := ec2Client.DescribeVolumes(context.Background(), &ecc.DescribeVolumesInput{
		VolumeIds: []string{vId},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Volume: %s err: %v", vId, err))
	}
	return result.Volumes[0]
}

/*
Snapshots are region specific
Localstack does have default snapshots, so we can see some of the snapshots that we never created
*/
func GetSnapshots(cfg awsV2.Config) []Snapshot {
	ec2Client := ecc.NewFromConfig(cfg)
	result, err := ec2Client.DescribeSnapshots(context.Background(), &ecc.DescribeSnapshotsInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Snapshots, err: %v", err))
		return nil
	}
	var snapshots []Snapshot
	for _, s := range result.Snapshots {
		launchTime := s.StartTime
		localZone, err := GetLocalTimeZone() // Empty string loads the local timezone
		if err != nil {
			fmt.Println("Error loading local timezone:", err)
			return nil
		}
		loc, _ := time.LoadLocation(localZone)
		IST := launchTime.In(loc)
		IST.Format("Mon Jan _2 15:04:05 2006")
		snapshot := Snapshot{
			SnapshotId: *s.SnapshotId,
			OwnerId:    *s.OwnerId,
			VolumeId:   *s.VolumeId,
			VolumeSize: strconv.Itoa(int(*s.VolumeSize)),
			StartTime:  IST.String(),
			State:      string(s.State),
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots
}

func GetSingleSnapshot(cfg awsV2.Config, sId string) types.Snapshot {
	ec2Serv := ecc.NewFromConfig(cfg)
	result, err := ec2Serv.DescribeSnapshots(context.Background(), &ecc.DescribeSnapshotsInput{
		SnapshotIds: []string{sId},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Snapshot: %s err: %v", sId, err))
	}
	return result.Snapshots[0]
}

/*
	AMIs are region specific
	Localstack does have default some AMIs, so we can see some of the AMIs that we never created
*/

func GetAMIs(cfg awsV2.Config) []ImageResp {
	ec2Serv := ecc.NewFromConfig(cfg)
	result, err := ec2Serv.DescribeImages(context.Background(), &ecc.DescribeImagesInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching AMIs, err: %v", err))
		return nil
	}
	var images []ImageResp
	for _, i := range result.Images {
		image := ImageResp{
			ImageId:       *i.ImageId,
			OwnerId:       *i.OwnerId,
			ImageLocation: *i.ImageLocation,
			Name:          *i.Name,
			ImageType:     string(i.ImageType),
		}
		images = append(images, image)
	}
	return images
}

func GetSingleAMI(cfg awsV2.Config, amiId string) types.Image {
	ec2Serv := ecc.NewFromConfig(cfg)
	result, err := ec2Serv.DescribeImages(context.Background(), &ecc.DescribeImagesInput{
		ImageIds: []string{amiId},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching AMI: %s err: %v ", amiId, err))
	}
	return result.Images[0]
}

func GetVPCs(cfg awsV2.Config) []VpcResp {
	ec2Serv := ecc.NewFromConfig(cfg)
	result, err := ec2Serv.DescribeVpcs(context.Background(), &ecc.DescribeVpcsInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching VPCs. err: %v ", err))
		return nil
	}
	var vpcs []VpcResp
	for _, v := range result.Vpcs {
		vpc := VpcResp{
			VpcId:           *v.VpcId,
			OwnerId:         *v.OwnerId,
			CidrBlock:       *v.CidrBlock,
			InstanceTenancy: string(v.InstanceTenancy),
			State:           string(v.State),
		}
		vpcs = append(vpcs, vpc)
	}
	return vpcs
}

func GetSingleVPC(cfg awsV2.Config, vpcId string) types.Vpc {
	ec2Serv := ecc.NewFromConfig(cfg)
	result, err := ec2Serv.DescribeVpcs(context.Background(), &ecc.DescribeVpcsInput{
		VpcIds: []string{vpcId},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching VPC: %s, err: %v", vpcId, err))
		return result.Vpcs[0]
	}
	return result.Vpcs[0]
}

func GetSubnets(cfg awsV2.Config, vpcId string) []SubnetResp {
	ec2Serv := ecc.NewFromConfig(cfg)
	result, err := ec2Serv.DescribeSubnets(context.Background(),
		&ecc.DescribeSubnetsInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{(vpcId)},
				},
			},
		})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Subnets. err: %v", err))
		return nil
	}
	var subnets []SubnetResp
	for _, s := range result.Subnets {
		subnet := SubnetResp{
			SubnetId:         *s.SubnetId,
			OwnerId:          *s.OwnerId,
			CidrBlock:        *s.CidrBlock,
			AvailabilityZone: *s.AvailabilityZone,
			State:            string(s.State),
		}
		subnets = append(subnets, subnet)
	}
	return subnets
}

func GetSingleSubnet(cfg awsV2.Config, sId string) types.Subnet {
	ec2Serv := ecc.NewFromConfig(cfg)
	result, err := ec2Serv.DescribeSubnets(context.Background(), &ecc.DescribeSubnetsInput{
		SubnetIds: []string{sId},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Subnet: %s, err: %v", sId, err))
		return result.Subnets[0]
	}
	return result.Subnets[0]
}

func GetLocalTimeZone() (string, error) {
	localZone, err := time.LoadLocation("") // Empty string loads the local timezone
	if err != nil {
		fmt.Println("Error loading local timezone:", err)
		return "", err
	}
	return localZone.String(), nil
}
