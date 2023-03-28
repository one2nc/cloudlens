package aws

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3GetObjectAPI interface {
	GetBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	GetPresignedUrl(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) *v4.PresignedHTTPRequest
}

func GetBucketsFromS3(ctx context.Context, api S3GetObjectAPI) (*s3.ListBucketsOutput, error) {
	object, err := api.GetBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func GetS3PresignedUrl(ctx context.Context, api S3GetObjectAPI, bucketName string, key string) string {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	psUrl := api.GetPresignedUrl(ctx, input)
	return psUrl.URL
}

type mockGetBucketsAPI func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
type mockGetPresignedUrlAPI func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) *v4.PresignedHTTPRequest

func (m mockGetBucketsAPI) GetBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m(ctx, params, optFns...)
}

func (m mockGetBucketsAPI) GetPresignedUrl(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) *v4.PresignedHTTPRequest {
	return nil
}

func (m mockGetPresignedUrlAPI) GetBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return nil, nil
}

func (m mockGetPresignedUrlAPI) GetPresignedUrl(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) *v4.PresignedHTTPRequest {
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

func TestGetPresignedUrl(t *testing.T) {
	cases := []struct {
		client     func(t *testing.T) S3GetObjectAPI
		expect     *v4.PresignedHTTPRequest
		key        string
		bucketName string
	}{
		{
			client: func(t *testing.T) S3GetObjectAPI {
				return mockGetPresignedUrlAPI(func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) *v4.PresignedHTTPRequest {
					t.Helper()
					return &v4.PresignedHTTPRequest{URL: *aws.String(fmt.Sprintf("https://%v.s3.amazonaws.com/%v", *params.Bucket, *params.Key))}
				})
			},
			expect:     &v4.PresignedHTTPRequest{URL: *aws.String("https://foo.s3.amazonaws.com/bar")},
			bucketName: "foo",
			key:        "bar",
		},
		{
			client: func(t *testing.T) S3GetObjectAPI {
				return mockGetPresignedUrlAPI(func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) *v4.PresignedHTTPRequest {
					t.Helper()
					return &v4.PresignedHTTPRequest{URL: *aws.String(fmt.Sprintf("https://%v.s3.amazonaws.com/%v", *params.Bucket, *params.Key))}
				})
			},
			expect:     &v4.PresignedHTTPRequest{URL: *aws.String("https://bar.s3.amazonaws.com/foo")},
			bucketName: "bar",
			key:        "foo",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got := GetS3PresignedUrl(ctx, tt.client(t), tt.bucketName, tt.key)
			expect := tt.expect.URL

			if got != expect {
				t.Errorf("expect %v, got %v", expect, got)
			}
		})
	}
}
