package aws

import (
	"fmt"
	"time"

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
	GetInstances(sess session.Session) ([]EC2Resp, error)
	// CreateSecGrp(name, desc, vpcId string)
	// CreateKeyPair(keyName string) (*ec2.CreateKeyPairOutput, error)
	// ListAllKeyPairs() ([]*string, error)
	// CreateInstances()
	GetSingleInstance() *ec2.DescribeInstancesOutput
	// CreateEc2() (*ec2.Reservation, error)
}

// type ec2Service struct {
// 	client ec2.EC2
// }

// func NewEc2Service(sess session.Session) Ec2Service {
// 	return ec2Service{client: *ec2.New(&sess)}
// }

func GetAllRegions() []string {
	return []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "af-south-1", "ap-east-1", "ap-south-2", "ap-southeast-3", "ap-south-1", "ap-northeast-3", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-south-2", "eu-north-1", "eu-central-2", "me-south-1", "me-central-1", "sa-east-1", "us-gov-east-1", "us-gov-west-1"}
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
				InstanceId:       *instance.InstanceId,
				InstanceType:     *instance.InstanceType,
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				InstanceState:    *instance.State.Name,
				PublicDNS:        *instance.PublicDnsName,
				PublicIPv4:       *instance.PublicIpAddress,
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

// func (e ec2Service) CreateInstances() {
// 	params := &ec2.RunInstancesInput{
// 		ImageId:      aws.String("ami-12345678"), // specify the ID of the image you want to use
// 		InstanceType: aws.String("t2.micro"),     // specify the instance type
// 		MinCount:     aws.Int64(3),
// 		MaxCount:     aws.Int64(3),
// 	}

// 	_, err := e.client.RunInstances(params)
// 	if err != nil {
// 		fmt.Println("Error in creating instances:", err)
// 	}
// }

// func (e ec2Service) CreateEc2() (*ec2.Reservation, error) {
// 	res, err := e.client.RunInstances(&ec2.RunInstancesInput{
// 		ImageId:      aws.String("new-image-id-one2n"),
// 		MinCount:     aws.Int64(int64(2)),
// 		MaxCount:     aws.Int64(int64(3)),
// 		InstanceType: aws.String("t4g.xlarge"),
// 		KeyName:      aws.String(""),
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

// func CreateSecGrp(name, desc, vpcId string) {
// 	if name == "" || desc == "" {
// 		log.Fatal("Group name and description require")
// 	}

// 	if vpcId == "" {
// 		vpc, err := e.client.DescribeVpcs(nil)
// 		if err != nil {
// 			log.Fatalf("Unable to describe VPCs, %v", err)
// 		}
// 		if len(vpc.Vpcs) == 0 {
// 			log.Fatal("No VPCs found to associate security group with.")
// 		}
// 		vpcId = aws.StringValue(vpc.Vpcs[0].VpcId)
// 	}
// 	createRes, err := e.client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
// 		GroupName:   aws.String(name),
// 		Description: aws.String(desc),
// 		VpcId:       aws.String(vpcId),
// 	})
// 	if err != nil {
// 		if aerr, ok := err.(awserr.Error); ok {
// 			switch aerr.Code() {
// 			case "InvalidVpcID.NotFound":
// 				log.Fatalf("invalid VPCID")
// 			case "InvalidGroup.Duplicate":
// 				log.Fatalf("Security group already exists.")
// 			}
// 		}
// 		log.Fatalf("Unable to create security group %v", err)
// 	}

// 	fmt.Printf("Created security group %s with VPC %s.\n",
// 		aws.StringValue(createRes.GroupId), vpcId)
// }

// func CreateKeyPair(keyName string) (*ec2.CreateKeyPairOutput, error) {
// 	result, err := e.client.CreateKeyPair(&ec2.CreateKeyPairInput{
// 		KeyName: aws.String(keyName),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func ListAllKeyPairs() ([]*string, error) {
// 	allKeyPairs, err := e.client.DescribeKeyPairs(nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var keyPairsNames []*string
// 	for _, pair := range allKeyPairs.KeyPairs {
// 		keyPairsNames = append(keyPairsNames, pair.KeyName)
// 	}
// 	return keyPairsNames, err
// }
