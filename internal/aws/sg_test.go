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

type SgAPI interface {
	GetSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
}

func GetSecurityGroupsTest(ctx context.Context, api SgAPI) (*ec2.DescribeSecurityGroupsOutput, error) {
	object, err := api.GetSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetSecurityGroups func(ctx context.Context, params *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)

func (m mockGetSecurityGroups) GetSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return m(ctx, params)
}

func TestGetSecurityGroups(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) SgAPI
		expect ec2.DescribeSecurityGroupsOutput
	}{
		{
			client: func(t *testing.T) SgAPI {
				return mockGetSecurityGroups(func(ctx context.Context, params *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					t.Helper()
					return &ec2.DescribeSecurityGroupsOutput{SecurityGroups: []types.SecurityGroup{{GroupId: aws.String("sec-group-1")}}}, nil
				})
			},
			expect: ec2.DescribeSecurityGroupsOutput{SecurityGroups: []types.SecurityGroup{{GroupId: aws.String("sec-group-1")}}},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetSecurityGroupsTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.SecurityGroups); i++ {
				fmt.Println("got:", *got.SecurityGroups[i].GroupId)
				fmt.Println("expect:", *tt.expect.SecurityGroups[i].GroupId)

				if *got.SecurityGroups[i].GroupId != *tt.expect.SecurityGroups[i].GroupId {
					t.Errorf("expect %v, got %v", *tt.expect.SecurityGroups[i].GroupId, *got.SecurityGroups[i].GroupId)
				}
			}
		})
	}
}
