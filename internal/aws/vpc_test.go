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

type VpcAPI interface {
	GetVPC(ctx context.Context, params *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
}

func GetVPCsTest(ctx context.Context, api VpcAPI) (*ec2.DescribeVpcsOutput, error) {
	object, err := api.GetVPC(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetVpcAPI func(ctx context.Context, params *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)

func (m mockGetVpcAPI) GetVPC(ctx context.Context, params *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return m(ctx, params)
}

func TestGetVPC(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) VpcAPI
		expect ec2.DescribeVpcsOutput
	}{
		{
			client: func(t *testing.T) VpcAPI {
				return mockGetVpcAPI(func(ctx context.Context, params *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					t.Helper()
					return &ec2.DescribeVpcsOutput{Vpcs: []types.Vpc{{VpcId: aws.String("vpc-1"), OwnerId: aws.String("000000000000"), CidrBlock: aws.String("172.31.0.0/16"), InstanceTenancy: types.TenancyDefault, State: types.VpcStateAvailable}}}, nil
				})
			},
			expect: ec2.DescribeVpcsOutput{Vpcs: []types.Vpc{{VpcId: aws.String("vpc-1"), OwnerId: aws.String("000000000000"), CidrBlock: aws.String("172.31.0.0/16"), InstanceTenancy: types.TenancyDefault, State: types.VpcStateAvailable}}},
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetVPCsTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.Vpcs); i++ {
				fmt.Println("got:", *got.Vpcs[i].VpcId)
				fmt.Println("expect:", *tt.expect.Vpcs[i].VpcId)

				if *got.Vpcs[i].VpcId != *tt.expect.Vpcs[i].VpcId {
					t.Errorf("expect %v, got %v", *tt.expect.Vpcs[i].VpcId, *got.Vpcs[i].VpcId)
				}

				fmt.Println("got:", *got.Vpcs[i].OwnerId)
				fmt.Println("expect:", *tt.expect.Vpcs[i].OwnerId)

				if *got.Vpcs[i].OwnerId != *tt.expect.Vpcs[i].OwnerId {
					t.Errorf("expect %v, got %v", *tt.expect.Vpcs[i].OwnerId, *got.Vpcs[i].OwnerId)
				}

				fmt.Println("got:", *got.Vpcs[i].OwnerId)
				fmt.Println("expect:", *tt.expect.Vpcs[i].OwnerId)

				if *got.Vpcs[i].OwnerId != *tt.expect.Vpcs[i].OwnerId {
					t.Errorf("expect %v, got %v", *tt.expect.Vpcs[i].OwnerId, *got.Vpcs[i].OwnerId)
				}

				fmt.Println("got:", *got.Vpcs[i].CidrBlock)
				fmt.Println("expect:", *tt.expect.Vpcs[i].CidrBlock)
				if *got.Vpcs[i].CidrBlock != *tt.expect.Vpcs[i].CidrBlock {
					t.Errorf("expect %v, got %v", *tt.expect.Vpcs[i].CidrBlock, *got.Vpcs[i].CidrBlock)
				}

				fmt.Println("got:", *&got.Vpcs[i].State)
				fmt.Println("expect:", *&tt.expect.Vpcs[i].State)
				if *&got.Vpcs[i].State != *&tt.expect.Vpcs[i].State {
					t.Errorf("expect %v, got %v", *&tt.expect.Vpcs[i].State, *&got.Vpcs[i].State)
				}
			}
		})
	}
}
