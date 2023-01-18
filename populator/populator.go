package pop

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/brianvoe/gofakeit"
)

func CreateBuckets(sess *session.Session) error {
	s3Service := s3.New(sess)
	for i := 0; i < 10; i++ {
		gofakeit.Seed(0)
		rWord := gofakeit.Password(true, false, true, false, false, 5)
		_, err := s3Service.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(rWord + "-bucket" + strconv.Itoa(gofakeit.Number(0, 999999999999999999))), CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: aws.String("ap-south-1")}})
		if err != nil {
			log.Println("S3 err here: ", rWord, err)
			return err
		}
	}
	return nil
}

func CreateEC2Instances(sess []*session.Session) error {
	insType := []string{"t2.micro", "t2.nano", "t2.small", "t2.medium", "t2.large",
		"t3a.nano", "t3a.micro", "t3a.small", "t3a.medium", "t3a.large",
		"t3.nano", "t3.micro", "t3.small", "t3.medium", "t3.large"}

	for i := 0; i < 200; i++ {
		ec2Service := ec2.New(sess[gofakeit.Number(0, len(sess)-1)])
		gofakeit.Seed(0)

		// we can change the core count according to the instance type
		// Some instance types allow only one core cpu to use
		coreCount := int64(gofakeit.Number(1, 1))

		// Each and Every Instances use 1 or 2 threads only
		// The all instance types we are creating here, use 2 threads.
		threads := int64(gofakeit.Number(1, 2))
		cpuOptReq := &ec2.CpuOptionsRequest{
			CoreCount:      &coreCount,
			ThreadsPerCore: &threads,
		}
		tag := []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("Instance" + strconv.Itoa(gofakeit.Number(0, 100000))),
					},
					{
						Key:   aws.String("Owner"),
						Value: aws.String("One2N"),
					},
				},
			},
		}
		monitoring := &ec2.RunInstancesMonitoringEnabled{
			Enabled: aws.Bool(true),
		}
		params := &ec2.RunInstancesInput{
			ImageId:           aws.String("ami-" + strconv.Itoa(gofakeit.Number(0, 9999999))), // specify the ID of the image you want to use
			InstanceType:      aws.String(insType[gofakeit.Number(0, len(insType)-1)]),        // specify the instance type
			MinCount:          aws.Int64(1),
			MaxCount:          aws.Int64(2),
			CpuOptions:        cpuOptReq,
			TagSpecifications: tag,
			Monitoring:        monitoring,
		}

		_, err := ec2Service.RunInstances(params)
		if err != nil {
			log.Println("Error in creating instances:", err)
			return err
		}
	}
	return nil
}

// It will be used in upcoming times along with EC2 instances
func CreateSecGrp(name, desc, vpcId string, sess *session.Session) {
	if name == "" || desc == "" {
		log.Fatal("Group name and description require")
	}

	ec2Service := ec2.New(sess)

	if vpcId == "" {
		vpc, err := ec2Service.DescribeVpcs(nil)
		if err != nil {
			log.Fatalf("Unable to describe VPCs, %v", err)
		}
		if len(vpc.Vpcs) == 0 {
			log.Fatal("No VPCs found to associate security group with.")
		}
		vpcId = aws.StringValue(vpc.Vpcs[0].VpcId)
	}
	createRes, err := ec2Service.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
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
