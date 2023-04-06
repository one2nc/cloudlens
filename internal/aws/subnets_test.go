package aws

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SubnetAPI interface {
	GetSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
}

func GetSubnetsTest(ctx context.Context, api SubnetAPI) (*ec2.DescribeSubnetsOutput, error) {
	object, err := api.GetSubnets(ctx, &ec2.DescribeSubnetsInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetSubnetAPI func(ctx context.Context, params *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)

func (m mockGetSubnetAPI) GetSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return m(ctx, params)
}

func TestGetSubnets(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) SubnetAPI
		expect ec2.DescribeSubnetsOutput
	}{
		{
			client: func(t *testing.T) SubnetAPI {
				return mockGetSubnetAPI(func(ctx context.Context, params *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					t.Helper()
					return &ec2.DescribeSubnetsOutput{Subnets: []types.Subnet{{SubnetId: aws.String("subnet-1"), AvailabilityZone: aws.String("ap-south-1"), OwnerId: aws.String("000000000000")}}}, nil
				})
			},
			expect: ec2.DescribeSubnetsOutput{Subnets: []types.Subnet{{SubnetId: aws.String("subnet-1"), AvailabilityZone: aws.String("ap-south-1"), OwnerId: aws.String("000000000000")}}},
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetSubnetsTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.Subnets); i++ {
				fmt.Println("got:", *got.Subnets[i].SubnetId)
				fmt.Println("expect:", *tt.expect.Subnets[i].SubnetId)

				if *got.Subnets[i].SubnetId != *tt.expect.Subnets[i].SubnetId {
					t.Errorf("expect %v, got %v", *tt.expect.Subnets[i].SubnetId, *got.Subnets[i].SubnetId)
				}

				fmt.Println("got:", *got.Subnets[i].AvailabilityZone)
				fmt.Println("expect:", *tt.expect.Subnets[i].AvailabilityZone)

				if *got.Subnets[i].AvailabilityZone != *tt.expect.Subnets[i].AvailabilityZone {
					t.Errorf("expect %v, got %v", *tt.expect.Subnets[i].AvailabilityZone, *got.Subnets[i].AvailabilityZone)
				}

				fmt.Println("got:", *got.Subnets[i].OwnerId)
				fmt.Println("expect:", *tt.expect.Subnets[i].OwnerId)

				if *got.Subnets[i].OwnerId != *tt.expect.Subnets[i].OwnerId {
					t.Errorf("expect %v, got %v", *tt.expect.Subnets[i].OwnerId, *got.Subnets[i].OwnerId)
				}
			}
		})
	}
}
