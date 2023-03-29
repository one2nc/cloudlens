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

type VolumesAPI interface {
	GetVolumes(ctx context.Context, params *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error)
}

func GetVolumesTest(ctx context.Context, api VolumesAPI) (*ec2.DescribeVolumesOutput, error) {
	object, err := api.GetVolumes(ctx, &ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetVolumesAPI func(ctx context.Context, params *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error)

func (m mockGetVolumesAPI) GetVolumes(ctx context.Context, params *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	return m(ctx, params)
}

func TestGetVolumes(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) VolumesAPI
		expect ec2.DescribeVolumesOutput
	}{
		{
			client: func(t *testing.T) VolumesAPI {
				return mockGetVolumesAPI(func(ctx context.Context, params *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
					t.Helper()
					return &ec2.DescribeVolumesOutput{Volumes: []types.Volume{{VolumeId: aws.String("vol-1"), Size: aws.Int32(32), VolumeType: types.VolumeTypeGp2}}}, nil
				})
			},
			expect: ec2.DescribeVolumesOutput{Volumes: []types.Volume{{VolumeId: aws.String("vol-1"), Size: aws.Int32(32), VolumeType: types.VolumeTypeGp2}}},
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetVolumesTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.Volumes); i++ {
				fmt.Println("got:", *got.Volumes[i].VolumeId)
				fmt.Println("expect:", *tt.expect.Volumes[i].VolumeId)

				if *got.Volumes[i].VolumeId != *tt.expect.Volumes[i].VolumeId {
					t.Errorf("expect %v, got %v", *tt.expect.Volumes[i].VolumeId, *got.Volumes[i].VolumeId)
				}

				fmt.Println("got size:", *got.Volumes[i].Size)
				fmt.Println("expect size:", *tt.expect.Volumes[i].Size)

				if *got.Volumes[i].Size != *tt.expect.Volumes[i].Size {
					t.Errorf("expect %v, got %v", *tt.expect.Volumes[i].Size, *got.Volumes[i].Size)
				}

				fmt.Println("got vol type:", *&got.Volumes[i].VolumeType)
				fmt.Println("expect vol type:", *&tt.expect.Volumes[i].VolumeType)
				if *&got.Volumes[i].VolumeType != *&tt.expect.Volumes[i].VolumeType {
					t.Errorf("expect %v, got %v", *&tt.expect.Volumes[i].VolumeType, *&got.Volumes[i].VolumeType)
				}
			}
		})
	}
}
