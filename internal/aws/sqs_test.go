package aws

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSGetAllQueuesAPI interface {
	GetQueues(ctx context.Context, params *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error)
}

func GetAllQueuesTest(ctx context.Context, api SQSGetAllQueuesAPI) (*sqs.ListQueuesOutput, error) {
	object, err := api.GetQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetAllQueuesAPI func(ctx context.Context, params *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error)

func (m mockGetAllQueuesAPI) GetQueues(ctx context.Context, params *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	return m(ctx, params)
}

func TestGetAllQueues(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) SQSGetAllQueuesAPI
		expect sqs.ListQueuesOutput
	}{
		{
			client: func(t *testing.T) SQSGetAllQueuesAPI {
				return mockGetAllQueuesAPI(func(ctx context.Context, params *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
					t.Helper()
					return &sqs.ListQueuesOutput{
						QueueUrls: aws.ToStringSlice(aws.StringSlice([]string{"http://localhost:4566/000000000000/queue-0", "http://localhost:4566/000000000000/queue-1", "http://localhost:4566/000000000000/queue-2"})),
					}, nil
				})
			},
			expect: sqs.ListQueuesOutput{
				QueueUrls: aws.ToStringSlice(aws.StringSlice([]string{"http://localhost:4566/000000000000/queue-0", "http://localhost:4566/000000000000/queue-1", "http://localhost:4566/000000000000/queue-2"})),
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetAllQueuesTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.QueueUrls); i++ {
				fmt.Println("got:", got.QueueUrls[i])
				fmt.Println("expect:", tt.expect.QueueUrls[i])

				if got.QueueUrls[i] != tt.expect.QueueUrls[i] {
					t.Errorf("expect %v, got %v", tt.expect.QueueUrls[i], got.QueueUrls[i])
				}
			}
		})
	}
}
