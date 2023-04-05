package aws

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IamAPI interface {
	GetIamUsers(ctx context.Context, params iam.ListUsersInput) (*iam.ListUsersOutput, error)
	GetIamGroups(ctx context.Context, params iam.ListGroupsInput) (*iam.ListGroupsOutput, error)
}

func GetAllIamUsersTest(ctx context.Context, api IamAPI) (*iam.ListUsersOutput, error) {
	object, err := api.GetIamUsers(ctx, iam.ListUsersInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func GetAllIamGroupsTest(ctx context.Context, api IamAPI) (*iam.ListGroupsOutput, error) {
	object, err := api.GetIamGroups(ctx, iam.ListGroupsInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetIamUsersAPI func(ctx context.Context, params iam.ListUsersInput) (*iam.ListUsersOutput, error)
type mockGetIamGroupsAPI func(ctx context.Context, params iam.ListGroupsInput) (*iam.ListGroupsOutput, error)

func (m mockGetIamUsersAPI) GetIamUsers(ctx context.Context, params iam.ListUsersInput) (*iam.ListUsersOutput, error) {
	return m(ctx, params)
}
func (m mockGetIamUsersAPI) GetIamGroups(ctx context.Context, params iam.ListGroupsInput) (*iam.ListGroupsOutput, error) {
	return nil, nil
}
func (m mockGetIamGroupsAPI) GetIamUsers(ctx context.Context, params iam.ListUsersInput) (*iam.ListUsersOutput, error) {
	return nil, nil
}
func (m mockGetIamGroupsAPI) GetIamGroups(ctx context.Context, params iam.ListGroupsInput) (*iam.ListGroupsOutput, error) {
	return m(ctx, params)
}

func TestIamUsers(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) IamAPI
		expect iam.ListUsersOutput
	}{
		{
			client: func(t *testing.T) IamAPI {
				return mockGetIamUsersAPI(func(ctx context.Context, params iam.ListUsersInput) (*iam.ListUsersOutput, error) {
					t.Helper()
					var usrArr []types.User
					usr := types.User{Arn: aws.String("arn:aws:iam:000000000000:user/Erdman"), UserId: aws.String("vyt1qsgh"), UserName: aws.String("Erdman"), CreateDate: aws.Time(time.Now())}
					usrArr = append(usrArr, usr)
					return &iam.ListUsersOutput{Users: usrArr}, nil
				})
			},
			expect: iam.ListUsersOutput{Users: []types.User{{Arn: aws.String("arn:aws:iam:000000000000:user/Erdman"), UserId: aws.String("vyt1qsgh"), UserName: aws.String("Erdman"), CreateDate: aws.Time(time.Now())}}},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetAllIamUsersTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.Users); i++ {
				fmt.Println("got:", *got.Users[i].UserId)
				fmt.Println("expect:", *tt.expect.Users[i].UserId)

				if *got.Users[i].UserId != *tt.expect.Users[i].UserId {
					t.Errorf("expect %v, got %v", *tt.expect.Users[i].UserId, *got.Users[i].UserId)
				}
			}
		})
	}
}

func TestIamGroups(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) IamAPI
		expect iam.ListGroupsOutput
	}{
		{
			client: func(t *testing.T) IamAPI {
				return mockGetIamGroupsAPI(func(ctx context.Context, params iam.ListGroupsInput) (*iam.ListGroupsOutput, error) {
					t.Helper()
					var grpArr []types.Group
					grp := types.Group{Arn: aws.String("arn:aws:iam:000000000000:group/Erdman"), GroupId: aws.String("ibaciunsoinonioucqnoiu"), GroupName: aws.String("Erdman-Group"), CreateDate: aws.Time(time.Now())}
					grpArr = append(grpArr, grp)
					return &iam.ListGroupsOutput{Groups: grpArr}, nil
				})
			},
			expect: iam.ListGroupsOutput{Groups: []types.Group{{Arn: aws.String("arn:aws:iam:000000000000:group/Erdman"), GroupId: aws.String("ibaciunsoinonioucqnoiu"), GroupName: aws.String("Erdman-Group"), CreateDate: aws.Time(time.Now())}}},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetAllIamGroupsTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.Groups); i++ {
				fmt.Println("got:", *got.Groups[i].GroupId)
				fmt.Println("expect:", *tt.expect.Groups[i].GroupId)

				if *got.Groups[i].GroupId != *tt.expect.Groups[i].GroupId {
					t.Errorf("expect %v, got %v", *tt.expect.Groups[i].GroupId, *got.Groups[i].GroupId)
				}
			}
		})
	}
}
