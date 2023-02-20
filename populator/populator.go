package pop

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/brianvoe/gofakeit"
)

/*
	- While creating s3 buckets and ec2 instances, Create policy there itself and add policy to a common slice
	- Use Interface so that we can create Iam policies for any service as needed.
	- In IAM service create users then user-group
	- Add random users to the group
	- Attach policies to the users and user group
	- Create roles

*/

type IamPolicy interface {
	createIamPolicy(iamService *iam.IAM) (*iam.Policy, error)
}

type s3IamPolicy struct {
	bName *string
}

type insIamPolicy struct {
	region, iId *string
}

var iamPolicies []*iam.Policy

func CreateBuckets(sess *session.Session) error {
	s3Service := s3.New(sess)
	iamService := iam.New(sess)
	for i := 0; i < 10; i++ {
		gofakeit.Seed(0)
		rWord := gofakeit.Password(true, false, true, false, false, 5)
		bName := aws.String("test-" + rWord + "-bucket" + strconv.Itoa(gofakeit.Number(0, 99999999)))
		_, err := s3Service.CreateBucket(&s3.CreateBucketInput{
			Bucket:                     bName,
			ObjectLockEnabledForBucket: aws.Bool(true),
		})
		if err != nil {
			log.Println("s3 create bucket error: ", rWord, err)
			return err
		}
		sip := s3IamPolicy{
			bName: bName,
		}
		iamPolicy, err := sip.createIamPolicy(iamService)
		if err != nil {
			return err
		}
		iamPolicies = append(iamPolicies, iamPolicy)

		key := []string{"Delhi", "Asia/Japan", "Asia/China/Beijing", "Jakarta", "Africa/Ghana", "North-America/Canada/Toronto",
			"Africa", "Africa/Jamaica", "Europe/England/London", "Vietnam", "Asia/South-Korea", "Asia/India/Kolkata",
			"Australia/Sydney", "Paris", "India/Kerala/Kochi", "Asia/Sri-Lanka", "Asia/Indonesia", "Europe/France", "Europe/Sweden",
			"Africa/West-Indies/City1", "North-America/USA/New-York", "Asia/India/Bangalore", "Asia/Nepal", "Asia/Burma"}
		for j := 0; j < len(key); j++ {
			if i == 4 {
				continue
			}
			body := []byte(gofakeit.Name())
			s3Service.PutObject(&s3.PutObjectInput{
				Bucket: bName,
				Key:    aws.String(key[j]),
				Body:   bytes.NewReader(body),
			})
		}

		addBucketEncryption(s3Service, bName)
		addBucketLifecycle(s3Service, bName)
		addBucketTag(s3Service, bName)
		addObjectLockConf(s3Service, bName)
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

	for i, s := range sess {
		n := 0
		if i == 0 {
			n = 50
		} else {
			n = 10
		}
		ec2Service := ec2.New(s)
		iamService := iam.New(s)
		gofakeit.Seed(0)
		vl, err := createVolume(ec2Service)
		if err != nil {
			fmt.Println(err)
		}
		ec2Service.CreateSnapshot(&ec2.CreateSnapshotInput{
			Description: aws.String("backup"),
			VolumeId:    vl.VolumeId,
		})

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
		// var lastEc2Policy *iam.Policy
		for i := 0; i < n; i++ {
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
			ec2, err := ec2Service.RunInstances(&ec2.RunInstancesInput{
				BlockDeviceMappings: blockDeviceMapping,
				// ImageId:             aws.String("ami-" + strconv.Itoa(gofakeit.Number(0, 9999999))),
				InstanceType:      aws.String(gofakeit.RandString(insType)),
				MaxCount:          aws.Int64(2),
				MinCount:          aws.Int64(1),
				SecurityGroupIds:  []*string{sg.GroupId},
				SubnetId:          aws.String(*sn.SubnetId),
				TagSpecifications: ec2Tag,
			})
			if err != nil {
				fmt.Println(fmt.Errorf("error in creating EC2 instance: %v", err))
			}

			for _, in := range ec2.Instances {
				log.Println("Ec2 instanced created: ", *in.InstanceId)
				iip := insIamPolicy{
					region: s.Config.Region,
					iId:    in.InstanceId,
				}
				iamPolicy, err := iip.createIamPolicy(iamService)
				if err != nil {
					fmt.Println(fmt.Errorf("error in creating Iam policy for instance:%v , err:%v", *in.InstanceId, err))
				}
				iamPolicies = append(iamPolicies, iamPolicy)
				// lastEc2Policy = iamPolicy.Policy
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
			fmt.Println(fmt.Errorf("error in creating Key Pairs: %v", err))
		}
	}
	return nil
}

func addBucketEncryption(service *s3.S3, bName *string) {
	service.PutBucketEncryption(&s3.PutBucketEncryptionInput{
		Bucket: bName,
		ServerSideEncryptionConfiguration: &s3.ServerSideEncryptionConfiguration{

			Rules: []*s3.ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: &s3.ServerSideEncryptionByDefault{
						SSEAlgorithm: aws.String("AES256"), // another algo: "aws:kms"
					},
				},
			},
		},
	})
	log.Println("Encryption added to the bucket: ", *bName)
}

func addBucketLifecycle(service *s3.S3, bName *string) {
	service.PutBucketLifecycleConfiguration(&s3.PutBucketLifecycleConfigurationInput{
		Bucket: bName,
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: []*s3.LifecycleRule{{
				Expiration: &s3.LifecycleExpiration{
					Days: aws.Int64(3650),
				},
				ID:     aws.String("Lifecycle Rule"),
				Prefix: aws.String("test"),
				Status: aws.String("Enabled"),
				Transitions: []*s3.Transition{
					{
						Days:         aws.Int64(365),
						StorageClass: aws.String("S3 Glacier Flexible Retrieval"),
					},
				},
			}},
		},
	})

	log.Println("Lifecycle added to the bucket: ", *bName)
}

func addBucketTag(service *s3.S3, bName *string) {
	gofakeit.Seed(0)
	service.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: bName,
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				{
					Key:   aws.String("Company"),
					Value: aws.String(gofakeit.Company()),
				},
			},
		},
	})

	log.Println("Tags added to the bucket: ", *bName)
}

func addObjectLockConf(service *s3.S3, bName *string) {
	service.PutObjectLockConfiguration(&s3.PutObjectLockConfigurationInput{
		Bucket: bName,
		ObjectLockConfiguration: &s3.ObjectLockConfiguration{
			ObjectLockEnabled: aws.String("Enable"),
			Rule: &s3.ObjectLockRule{
				DefaultRetention: &s3.DefaultRetention{
					Mode: aws.String("COMPLIANCE"),
					Days: aws.Int64(30),
				},
			},
		},
	})
	log.Println("Object Lock added to the bucket: ", *bName)
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
	log.Println("VPC created: ", *vpc.Vpc.VpcId)
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

	log.Println("Subnet created: ", *subnet.Subnet.SubnetId)
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
		GroupName:         aws.String("name"),
		Description:       aws.String("desc"),
		TagSpecifications: sgTag,
		VpcId:             vpcId,
	})
	if err != nil {
		return nil, fmt.Errorf("error in creating Security Group: %v", err)
	}

	log.Println("Security Group created: ", *sg.GroupId)

	addSecGrpRules(service, sg.GroupId)
	return sg, nil
}

func addSecGrpRules(service *ec2.EC2, sgId *string) {
	var sgRule = &ec2.IpPermission{
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
		IpProtocol: aws.String("tcp"),
		IpRanges: []*ec2.IpRange{
			{
				CidrIp:      aws.String("192.1.0.0/24"),
				Description: aws.String("Allow login (SSH) port"),
			},
		},
	}

	service.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:        aws.String("192.1.0.0/24"),
		FromPort:      aws.Int64(22),
		GroupId:       sgId,
		IpPermissions: []*ec2.IpPermission{sgRule},
		IpProtocol:    aws.String("tcp"),
		ToPort:        aws.Int64(22),
	})

	service.AuthorizeSecurityGroupEgress(&ec2.AuthorizeSecurityGroupEgressInput{
		CidrIp:        aws.String("192.1.0.0/24"),
		FromPort:      aws.Int64(22),
		GroupId:       sgId,
		IpPermissions: []*ec2.IpPermission{sgRule},
		IpProtocol:    aws.String("tcp"),
		ToPort:        aws.Int64(22),
	})
	log.Println("Security Group Rules added to ", *sgId)
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
	log.Println("Volume Created: ", *v.VolumeId)
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

func (sip s3IamPolicy) createIamPolicy(iamService *iam.IAM) (*iam.Policy, error) {
	pd := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": "s3:*",
				"Resource": "arn:aws:s3:::` + *sip.bName + `/*"
			}
		]
	}`
	iamPolicy, err := iamService.CreatePolicy(&iam.CreatePolicyInput{
		Description:    aws.String("Grant access to all objects of this bucket" + *sip.bName),
		PolicyDocument: aws.String(pd),
		PolicyName:     aws.String(gofakeit.FirstName() + "-s3-policy"),
	})
	if err != nil {
		return nil, fmt.Errorf("create iam policy for s3 bucket error: %v", err)
	}
	log.Println("Iam policy created for the S3 bucket: ", *sip.bName)
	return iamPolicy.Policy, nil
}

func (iip insIamPolicy) createIamPolicy(iamService *iam.IAM) (*iam.Policy, error) {
	pd := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"ec2:*"
				],
				"Resource": "arn:aws:ec2:` + *iip.region + `:000000000000:instance/` + *iip.iId + `"
			}
		]
	}`
	iamPolicy, err := iamService.CreatePolicy(&iam.CreatePolicyInput{
		Description:    aws.String("Grant all types of access to this instance: " + *iip.iId),
		PolicyDocument: aws.String(pd),
		PolicyName:     aws.String(gofakeit.FirstName() + "-ec2-policy"),
	})
	if err != nil {
		return nil, fmt.Errorf("create iam policy for instance error: %v", err)
	}
	log.Println("Iam policy created for the EC2 Instance: ", *iip.iId)
	return iamPolicy.Policy, nil
}

func IamAwsSrv(sess *session.Session) error {
	srv := iam.New(sess)
	gofakeit.Seed(0)

	var users []*iam.User

	//Creating 10 users
	for i := 0; i < 10; i++ {
		gofakeit.Seed(0)
		cuo, err := srv.CreateUser(&iam.CreateUserInput{
			UserName: aws.String(gofakeit.LastName()),
		})
		if err != nil {
			return fmt.Errorf("error creating user: %v", err)
		}
		log.Println("user created: ", *cuo.User.UserName)
		users = append(users, cuo.User)
	}

	//Creating User Group
	cgo, _ := srv.CreateGroup(&iam.CreateGroupInput{
		GroupName: aws.String(fmt.Sprintf("one2n-engineers-%d", gofakeit.Number(1, 2000))),
	})
	log.Println("user-group created: ", *cgo.Group.GroupName)
	for i := 0; i < 5; i++ {
		srv.AddUserToGroup(&iam.AddUserToGroupInput{
			GroupName: cgo.Group.GroupName,
			UserName:  users[i].UserName,
		})
		log.Println("user: ", *users[i].UserName, "added to the group: ", *cgo.Group.GroupName)

		srv.AttachUserPolicy(&iam.AttachUserPolicyInput{
			PolicyArn: iamPolicies[gofakeit.Number(0, len(iamPolicies)-1)].Arn,
			UserName:  users[i+5].UserName,
		})

		srv.AttachGroupPolicy(&iam.AttachGroupPolicyInput{
			GroupName: cgo.Group.GroupName,
			PolicyArn: aws.String(*iamPolicies[i].Arn),
		})
	}
	for i := 0; i < 4; i++ {
		err := createIamRole(srv)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating role: %v", err))
		}
	}

	return nil
}

func createIamRole(srv *iam.IAM) error {
	// Have to add more actions in future
	rpd := `{
		"Version":"2012-10-17",
		"Statement":[
			{
				"Effect":"Allow",
				"Principal":{"Service":["ec2.amazonaws.com"]},
				"Action":["s3:*", "ec2:*", "lambda:*"]
			}
		]
	}`

	cr, err := srv.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(rpd),
		RoleName:                 aws.String(gofakeit.FirstName()),
	})
	log.Println("Iam Role Created: ", *cr.Role.RoleId)
	if err != nil {
		return fmt.Errorf("error creating role: %v", err)
	}
	srv.AttachRolePolicy(&iam.AttachRolePolicyInput{
		PolicyArn: iamPolicies[gofakeit.Number(0, len(iamPolicies)-1)].Arn,
		RoleName:  cr.Role.RoleName,
	})
	// if err != nil {
	// 	return fmt.Errorf("error creating role policy: %v", err)
	// }
	log.Println("Iam Role Policy Created for: ", *cr.Role.RoleId)
	return nil
}

func CreateQueueAndSetMessages(sessions []*session.Session) error {
	for _, sess := range sessions {
		sqsServ := *sqs.New(sess)
		for i := 0; i < 5; i++ {
			gofakeit.Seed(0)
			qName := gofakeit.FirstName()
			res, err := sqsServ.CreateQueue(&sqs.CreateQueueInput{
				QueueName: aws.String(qName + "-queue-" + strconv.Itoa(i)),
			})
			if err != nil {
				fmt.Println("Error in creating queue: ", qName, " err: ", err)
				return err
			}
			log.Println("Queue created..."+qName, " and url is:", res.QueueUrl)
			queueUrl := res.QueueUrl
			messages := []*sqs.SendMessageBatchRequestEntry{
				{
					Id:          aws.String(strconv.Itoa(i)),
					MessageBody: aws.String("Hello-world-" + strconv.Itoa(i) + "!"),
				},
				{
					Id:          aws.String(strconv.Itoa(gofakeit.Number(i, 99999999))),
					MessageBody: aws.String("Hello-world-" + strconv.Itoa(i+gofakeit.Number(i, 99999999)) + "!"),
				},
			}
			batchRequest := &sqs.SendMessageBatchInput{
				QueueUrl: aws.String(*queueUrl),
				Entries:  messages,
			}
			_, err = sqsServ.SendMessageBatch(batchRequest)
			if err != nil {
				fmt.Println("Error sending messages:", err)
				return err
			}
			log.Println("Messages added to queue:", res.QueueUrl)
		}
	}
	return nil
}

func CreateLambdaFunction(sessions []*session.Session) error {
	for _, sess := range sessions {
		lambdaServ := lambda.New(sess)
		codeBytes, err := ioutil.ReadFile("code.zip")
		if err != nil {
			log.Fatalf("Failed to read function code: %v", err)
			return err
		}
		params := &lambda.CreateFunctionInput{
			FunctionName: aws.String("lambdaaa-func-" + strconv.Itoa(gofakeit.Number(0, 999999))),
			//Runtime:      aws.String(*aws.String("python3.7")),
			Role: aws.String("arn:aws:iam::000000000000:role/Andre"),
			Code: &lambda.FunctionCode{
				ZipFile: codeBytes,
			},
		}
		resp, err := lambdaServ.CreateFunction(params)
		if err != nil {
			log.Fatalf("Failed to create function: %v", err)
			return err
		}
		log.Printf("Function ARN: %s\n", aws.StringValue(resp.FunctionArn))
	}
	return nil
}

func GetDefaultAWSRegion() string {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load AWS SDK config: %v\n", err)
		os.Exit(1)
	}
	region := cfg.Region
	return region
}
