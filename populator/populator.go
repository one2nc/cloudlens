package pop

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/brianvoe/gofakeit"
)

func CreateBuckets(sess *session.Session) error {
	s3Service := s3.New(sess)
	for i := 0; i < 20; i++ {
		gofakeit.Seed(0)
		rWord := gofakeit.Password(true, false, true, false, false, 5)
		bName := aws.String(rWord + "-bucket" + strconv.Itoa(gofakeit.Number(0, 999999999999999999)))
		_, err := s3Service.CreateBucket(&s3.CreateBucketInput{
			Bucket:                    bName,
			CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: aws.String("ap-south-1")},
		})
		if err != nil {
			log.Println("S3 err here: ", rWord, err)
			return err
		}

		key := []string{"Delhi", "Asia/Japan", "Asia/China/Beijing", "Jakarta", "Africa/Ghana", "North-America/Canada/Toronto",
			"Africa", "Africa/Jamaice", "Europe/England/London", "Vietnam", "Asia/South-Korea", "Asia/India/Kolkata",
			"Australia/Sydney", "Paris", "India/Kerala/Kochi", "Asia/Sri-Lanka", "Asia/Indonesia", "Europe/France", "Europe/Sweden",
			"Africa/West-Indies/City1", "North-America/USA/New-York", "Asia/India/Bangalore", "Asia/Nepal", "Asia/Burma"}
		for i := 0; i < len(key); i++ {

			body := []byte(gofakeit.Name())
			s3Service.PutObject(&s3.PutObjectInput{
				Bucket: bName,
				Key:    aws.String(key[i]),
				Body:   bytes.NewReader(body),
			})
		}
	}
	return nil
}

func CreateEC2Instances(sess []*session.Session) error {

	insType := []string{"t2.micro", "t2.nano", "t2.small", "t2.medium", "t2.large",
		"t3a.nano", "t3a.micro", "t3a.small", "t3a.medium", "t3a.large",
		"t3.nano", "t3.micro", "t3.small", "t3.medium", "t3.large"}

	// Consider 3-5 sessions
	// Create VPC, Subnets, SG, Volumes
	// Create 5-10 EC Instances for each Sessions.
	for _, s := range sess {
		ec2Service := ec2.New(s)
		gofakeit.Seed(0)
		vpc, err := createVpc(ec2Service)
		if err != nil {
			fmt.Println(err)
		}

		sn, err := createSubnet(ec2Service, vpc.VpcId)
		if err != nil {
			fmt.Println(err)
		}

		sg, err := createSecGrp(ec2Service, vpc.VpcId)
		if err != nil {
			fmt.Println(err)
		}

		blockDeviceMapping := createEbsMapping()

		for i := 0; i < 10; i++ {

			ec2Tag := []*ec2.TagSpecification{
				{
					ResourceType: aws.String("instance"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("instance-" + strconv.Itoa(gofakeit.Number(0, 100000))),
						},
						{
							Key:   aws.String("Owner"),
							Value: aws.String("One2N"),
						},
					},
				},
			}

			_, err := ec2Service.RunInstances(&ec2.RunInstancesInput{
				BlockDeviceMappings: blockDeviceMapping,
				ImageId:             aws.String("ami-" + strconv.Itoa(gofakeit.Number(0, 9999999))),
				InstanceType:        aws.String(gofakeit.RandString(insType)),
				MaxCount:            aws.Int64(2),
				MinCount:            aws.Int64(1),
				SecurityGroupIds:    []*string{sg.GroupId},
				SubnetId:            aws.String(*sn.SubnetId),
				TagSpecifications:   ec2Tag,
			})
			if err != nil {
				fmt.Println(fmt.Errorf("error in creating EC2 instance: %v", err))
			}
		}
	}
	return nil
}

func CreateKeyPair(sess []*session.Session) error {
	for _, s := range sess {
		service := ec2.New(s)
		gofakeit.Seed(0)
		_, err := service.CreateKeyPair(&ec2.CreateKeyPairInput{
			KeyFormat: aws.String("PEM"),
			KeyName:   aws.String(gofakeit.FirstName()),
			KeyType:   aws.String("ssh-rsa"),
		})
		if err != nil {
			return fmt.Errorf("error in creating Key Pairs: %v", err)
		}
	}
	return nil
}

func createVpc(service *ec2.EC2) (*ec2.Vpc, error) {
	vpcTag := []*ec2.TagSpecification{
		{
			ResourceType: aws.String("instance"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("VPC" + strconv.Itoa(gofakeit.Number(0, 100000))),
				},
				{
					Key:   aws.String("Owner"),
					Value: aws.String("One2N"),
				},
			},
		},
	}
	vpc, err := service.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock:         aws.String("192.0.0.0/16"),
		TagSpecifications: vpcTag,
	})
	if err != nil {
		return nil, fmt.Errorf("error in creating VPC: %v", err)
	}
	return vpc.Vpc, nil
}

func createSubnet(service *ec2.EC2, vpcId *string) (*ec2.Subnet, error) {
	snTag := []*ec2.TagSpecification{
		{
			ResourceType: aws.String("instance"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("Subnet" + strconv.Itoa(gofakeit.Number(0, 100000))),
				},
				{
					Key:   aws.String("Owner"),
					Value: aws.String("One2N"),
				},
			},
		},
	}

	subnet, err := service.CreateSubnet(&ec2.CreateSubnetInput{
		CidrBlock:         aws.String("192.0.0.0/24"),
		TagSpecifications: snTag,
		VpcId:             vpcId,
	})
	if err != nil {
		return nil, fmt.Errorf("error in creating Subnet: %v", err)
	}

	return subnet.Subnet, nil
}

func createSecGrp(service *ec2.EC2, vpcId *string) (*ec2.CreateSecurityGroupOutput, error) {
	sgTag := []*ec2.TagSpecification{
		{
			ResourceType: aws.String("instance"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("Security-Group" + strconv.Itoa(gofakeit.Number(0, 100000))),
				},
				{
					Key:   aws.String("Owner"),
					Value: aws.String("One2N"),
				},
			},
		},
	}
	sg, err := service.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		Description:       aws.String("desc"),
		GroupName:         aws.String("name"),
		TagSpecifications: sgTag,
		VpcId:             vpcId,
	})
	if err != nil {
		return nil, fmt.Errorf("error in creating Security Group: %v", err)
	}

	return sg, nil
}

func createEbsMapping() []*ec2.BlockDeviceMapping {
	volumeSize := []int64{1, 2, 4, 8, 32, 64, 128}
	deviceName := []string{"dev/sdf", "dev/sdg", "dev/sdh", "dev/sdi", "dev/sdj", "dev/sdk", "dev/sdl", "dev/sdm", "dev/sdn", "dev/sdo", "dev/sdp"}
	blockDeviceMapping := []*ec2.BlockDeviceMapping{}

	for i := 0; i < gofakeit.Number(1, 3); i++ {
		ebs := ec2.EbsBlockDevice{
			VolumeSize: aws.Int64(volumeSize[gofakeit.Number(0, len(volumeSize)-1)]),
			VolumeType: aws.String("gp2"), // for t-type instances gp2 alone is supported
		}
		bdm := ec2.BlockDeviceMapping{
			DeviceName: aws.String(gofakeit.RandString(deviceName)),
			Ebs:        &ebs,
		}
		blockDeviceMapping = append(blockDeviceMapping, &bdm)
	}

	return blockDeviceMapping
}

func createVolume(service *ec2.EC2) (*ec2.Volume, error) {
	volumeSize := []int64{1, 2, 4, 8, 32, 64, 128}
	vTag := []*ec2.TagSpecification{
		{
			ResourceType: aws.String("instance"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("Volume" + strconv.Itoa(gofakeit.Number(0, 100000))),
				},
				{
					Key:   aws.String("Owner"),
					Value: aws.String("One2N"),
				},
			},
		},
	}

	v, err := service.CreateVolume(&ec2.CreateVolumeInput{
		AvailabilityZone:  aws.String("us-east-1a"),
		Size:              aws.Int64(volumeSize[gofakeit.Number(0, len(volumeSize)-1)]),
		TagSpecifications: vTag,
		VolumeType:        aws.String("gp2"),
	})

	if err != nil {
		return nil, fmt.Errorf("error in creating Volume: %v", err)
	}
	return v, nil
}

// working on it
func makeSubnetPublic(service *ec2.EC2, vpcId, snId, insId *string) {
	igw, _ := service.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})

	_, _ = service.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: igw.InternetGateway.InternetGatewayId,
		VpcId:             vpcId,
	})

	rt, _ := service.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: vpcId,
	})

	service.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock:   aws.String("192.0.0.1/24"),
		GatewayId:              igw.InternetGateway.InternetGatewayId,
		InstanceId:             insId,
		RouteTableId:           rt.RouteTable.RouteTableId,
		VpcEndpointId:          vpcId,
		VpcPeeringConnectionId: vpcId,
	})

	service.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		GatewayId:    igw.InternetGateway.InternetGatewayId,
		RouteTableId: rt.RouteTable.RouteTableId,
		SubnetId:     snId,
	})
}
