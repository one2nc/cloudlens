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

type SnapshotAPI interface {
	GetSnapshots(ctx context.Context, params *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error)
}

func GetSnapshotsTest(ctx context.Context, api SnapshotAPI) (*ec2.DescribeSnapshotsOutput, error) {
	object, err := api.GetSnapshots(ctx, &ec2.DescribeSnapshotsInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetSnapshotsAPI func(ctx context.Context, params *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error)

func (m mockGetSnapshotsAPI) GetSnapshots(ctx context.Context, params *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
	return m(ctx, params)
}

func TestGetSnapshots(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) SnapshotAPI
		expect ec2.DescribeSnapshotsOutput
	}{
		{
			client: func(t *testing.T) SnapshotAPI {
				return mockGetSnapshotsAPI(func(ctx context.Context, params *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
					t.Helper()
					return &ec2.DescribeSnapshotsOutput{Snapshots: []types.Snapshot{{SnapshotId: aws.String("snap-1"), OwnerId: aws.String("owner-1"), VolumeId: aws.String("vol-1"), VolumeSize: aws.Int32(32), State: types.SnapshotStateCompleted}}}, nil
				})
			},
			expect: ec2.DescribeSnapshotsOutput{Snapshots: []types.Snapshot{{SnapshotId: aws.String("snap-1"), OwnerId: aws.String("owner-1"), VolumeId: aws.String("vol-1"), VolumeSize: aws.Int32(32), State: types.SnapshotStateCompleted}}},
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetSnapshotsTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.Snapshots); i++ {
				fmt.Println("got:", *got.Snapshots[i].SnapshotId)
				fmt.Println("expect:", *tt.expect.Snapshots[i].SnapshotId)

				if *got.Snapshots[i].SnapshotId != *tt.expect.Snapshots[i].SnapshotId {
					t.Errorf("expect %v, got %v", *tt.expect.Snapshots[i].SnapshotId, *got.Snapshots[i].SnapshotId)
				}

				fmt.Println("got:", *got.Snapshots[i].OwnerId)
				fmt.Println("expect:", *tt.expect.Snapshots[i].OwnerId)

				if *got.Snapshots[i].OwnerId != *tt.expect.Snapshots[i].OwnerId {
					t.Errorf("expect %v, got %v", *tt.expect.Snapshots[i].OwnerId, *got.Snapshots[i].OwnerId)
				}

				fmt.Println("got:", *got.Snapshots[i].VolumeId)
				fmt.Println("expect:", *tt.expect.Snapshots[i].VolumeId)

				if *got.Snapshots[i].VolumeId != *tt.expect.Snapshots[i].VolumeId {
					t.Errorf("expect %v, got %v", *tt.expect.Snapshots[i].VolumeId, *got.Snapshots[i].VolumeId)
				}

				fmt.Println("got:", *got.Snapshots[i].VolumeSize)
				fmt.Println("expect:", *tt.expect.Snapshots[i].VolumeSize)
				if *got.Snapshots[i].VolumeSize != *tt.expect.Snapshots[i].VolumeSize {
					t.Errorf("expect %v, got %v", *tt.expect.Snapshots[i].VolumeSize, *got.Snapshots[i].VolumeSize)
				}

				fmt.Println("got:", *&got.Snapshots[i].State)
				fmt.Println("expect:", *&tt.expect.Snapshots[i].State)
				if *&got.Snapshots[i].State != *&tt.expect.Snapshots[i].State {
					t.Errorf("expect %v, got %v", *&tt.expect.Snapshots[i].State, *&got.Snapshots[i].State)
				}
			}
		})
	}
}
