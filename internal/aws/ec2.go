package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Resp struct {
	InstanceId       string
	InstanceType     string
	AvailabilityZone string
	InstanceState    string
	PublicDNS        string
	PublicIPv4       string
	MonitoringState  string
	LaunchTime       string
}
type Ec2Service interface {
	CreateInstances()
	GetInstances() ([]EC2Resp, error)
	CreateEc2() (*ec2.Reservation, error)
	CreateSecGrp(name, desc, vpcId string)
	RegionsAndAvailZones()
	CreateKeyPair(keyName string) (*ec2.CreateKeyPairOutput, error)
	ListAllKeyPairs() ([]*string, error)
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
			ec2Resp := &EC2Resp{
				InstanceId:       *instance.InstanceId,
				InstanceType:     *instance.InstanceType,
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				InstanceState:    *instance.State.Name,
				PublicDNS:        *instance.PublicDnsName,
				PublicIPv4:       *instance.PublicIpAddress,
				MonitoringState:  *instance.Monitoring.State,
				LaunchTime:       instance.LaunchTime.Format("2006-01-02 15:04:05")}
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

func (e ec2Service) CreateEc2() (*ec2.Reservation, error) {
	res, err := e.client.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String("new-image-id-one2n"),
		MinCount:     aws.Int64(int64(2)),
		MaxCount:     aws.Int64(int64(3)),
		InstanceType: aws.String("t4g.xlarge"),
		KeyName:      aws.String(""),
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e ec2Service) CreateSecGrp(name, desc, vpcId string) {
	if name == "" || desc == "" {
		log.Fatal("Group name and description require")
	}

	if vpcId == "" {
		vpc, err := e.client.DescribeVpcs(nil)
		if err != nil {
			log.Fatalf("Unable to describe VPCs, %v", err)
		}
		if len(vpc.Vpcs) == 0 {
			log.Fatal("No VPCs found to associate security group with.")
		}
		vpcId = aws.StringValue(vpc.Vpcs[0].VpcId)
	}
	createRes, err := e.client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(name),
		Description: aws.String(desc),
		VpcId:       aws.String(vpcId),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidVpcID.NotFound":
				log.Fatalf("invalid VPCID")
			case "InvalidGroup.Duplicate":
				log.Fatalf("Security group already exists.")
			}
		}
		log.Fatalf("Unable to create security group %v", err)
	}

	fmt.Printf("Created security group %s with VPC %s.\n",
		aws.StringValue(createRes.GroupId), vpcId)
}

func (e ec2Service) RegionsAndAvailZones() {
	regions, err := e.client.DescribeRegions(nil)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	availZones, err := e.client.DescribeAvailabilityZones(nil)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Regions: ", regions, "\n Available Zones: ", availZones.AvailabilityZones)
}

func (e ec2Service) CreateKeyPair(keyName string) (*ec2.CreateKeyPairOutput, error) {
	result, err := e.client.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(keyName),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e ec2Service) ListAllKeyPairs() ([]*string, error) {
	allKeyPairs, err := e.client.DescribeKeyPairs(nil)
	if err != nil {
		return nil, err
	}
	var keyPairsNames []*string
	for _, pair := range allKeyPairs.KeyPairs {
		keyPairsNames = append(keyPairsNames, pair.KeyName)
	}
	return keyPairsNames, err
}
