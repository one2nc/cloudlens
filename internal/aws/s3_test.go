package aws

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3GetObjectAPI interface {
	GetBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

func GetBucketsFromS3(ctx context.Context, api S3GetObjectAPI) (*s3.ListBucketsOutput, error) {
	object, err := api.GetBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetBucketsAPI func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)

func (m mockGetBucketsAPI) GetBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetObjectFromS3(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) S3GetObjectAPI
		expect s3.ListBucketsOutput
	}{
		{
			client: func(t *testing.T) S3GetObjectAPI {
				return mockGetBucketsAPI(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
					t.Helper()
					return &s3.ListBucketsOutput{
						Buckets: []types.Bucket{{CreationDate: aws.Time(time.Now()), Name: aws.String("bucket-1")}},
					}, nil
				})
			},
			expect: s3.ListBucketsOutput{
				Buckets: []types.Bucket{{CreationDate: aws.Time(time.Now()), Name: aws.String("bucket-1")}},
			},
		},
		{
			client: func(t *testing.T) S3GetObjectAPI {
				return mockGetBucketsAPI(func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
					t.Helper()
					return &s3.ListBucketsOutput{
						Buckets: []types.Bucket{{CreationDate: aws.Time(time.Now()), Name: aws.String("bucket-1")}, {CreationDate: aws.Time(time.Now()), Name: aws.String("bucket-2")}},
					}, nil
				})
			},
			expect: s3.ListBucketsOutput{
				Buckets: []types.Bucket{{CreationDate: aws.Time(time.Now()), Name: aws.String("bucket-1")}, {CreationDate: aws.Time(time.Now()), Name: aws.String("bucket-2")}},
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetBucketsFromS3(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.Buckets); i++ {
				fmt.Println("got:", *got.Buckets[i].Name)
				fmt.Println("expect:", *tt.expect.Buckets[i].Name)

				if *got.Buckets[i].Name != *tt.expect.Buckets[i].Name {
					t.Errorf("expect %v, got %v", tt.expect.Buckets[i], got.Buckets[i])
				}
			}
		})
	}
}
