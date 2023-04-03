package aws

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LabmdaFunctionsAPI interface {
	GetLambdaFunctions(ctx context.Context, params *lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error)
}

func GetAllLambdaFunctionsTest(ctx context.Context, api LabmdaFunctionsAPI) (*lambda.ListFunctionsOutput, error) {
	object, err := api.GetLambdaFunctions(ctx, &lambda.ListFunctionsInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetAllLambdaFunctionsAPI func(ctx context.Context, params *lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error)

func (m mockGetAllLambdaFunctionsAPI) GetLambdaFunctions(ctx context.Context, params *lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error) {
	return m(ctx, params)
}

func TestGetAllLambdaFunctions(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) LabmdaFunctionsAPI
		expect lambda.ListFunctionsOutput
		region string
	}{
		{
			client: func(t *testing.T) LabmdaFunctionsAPI {
				return mockGetAllLambdaFunctionsAPI(func(ctx context.Context, params *lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error) {
					t.Helper()
					funcArr := []types.FunctionConfiguration{{FunctionName: aws.String("lambda-func-1"), Role: aws.String("aen:aws:iam:000000000000:role/Role"), FunctionArn: aws.String("arn:aws:lambda:us-east-1:0000000000000:function:lambda-func-1")}}
					return &lambda.ListFunctionsOutput{Functions: funcArr}, nil
				})
			},
			expect: lambda.ListFunctionsOutput{Functions: []types.FunctionConfiguration{{FunctionName: aws.String("lambda-func-1"), Role: aws.String("aen:aws:iam:000000000000:role/Role"), FunctionArn: aws.String("arn:aws:lambda:us-east-1:0000000000000:function:lambda-func-1")}}},
			region: "us-east-1",
		},
		{
			client: func(t *testing.T) LabmdaFunctionsAPI {
				return mockGetAllLambdaFunctionsAPI(func(ctx context.Context, params *lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error) {
					t.Helper()
					funcArr := []types.FunctionConfiguration{{FunctionName: aws.String("lambda-func-2"), Role: aws.String("aen:aws:iam:000000000000:role/Role"), FunctionArn: aws.String("arn:aws:lambda:us-east-2:0000000000000:function:lambda-func-2")}, {FunctionName: aws.String("lambda-func-3"), Role: aws.String("aen:aws:iam:000000000000:role/Role"), FunctionArn: aws.String("arn:aws:lambda:us-east-2:0000000000000:function:lambda-func-3")}}
					return &lambda.ListFunctionsOutput{Functions: funcArr}, nil
				})
			},
			expect: lambda.ListFunctionsOutput{Functions: []types.FunctionConfiguration{{FunctionName: aws.String("lambda-func-2"), Role: aws.String("aen:aws:iam:000000000000:role/Role"), FunctionArn: aws.String("arn:aws:lambda:us-east-2:0000000000000:function:lambda-func-2")}, {FunctionName: aws.String("lambda-func-3"), Role: aws.String("aen:aws:iam:000000000000:role/Role"), FunctionArn: aws.String("arn:aws:lambda:us-east-2:0000000000000:function:lambda-func-3")}}},
			region: "us-east-2",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetAllLambdaFunctionsTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if tt.region == "us-east-1" || tt.region == "us-east-2" {
				for i := 0; i < len(got.Functions); i++ {
					fmt.Println("got:", *got.Functions[i].FunctionName)
					fmt.Println("expect:", *tt.expect.Functions[i].FunctionName)

					if *got.Functions[i].FunctionName != *tt.expect.Functions[i].FunctionName {
						t.Errorf("expect %v, got %v", tt.expect.Functions[i].FunctionName, got.Functions[i].FunctionName)
					}
				}
			}
		})
	}
}

func TestGetLambdaFunctions(t *testing.T) {
	lambdaFunctions := []string{"lambdaaa-func-0", "lambdaaa-func-1", "lambdaaa-func-2"}
	localStackConf, err := GetLocalstackCfg("us-east-1")
	if err != nil {
		fmt.Println("Error in getting config")
	}

	lambdaInfo, err := GetAllLambdaFunctions(localStackConf)
	for i := 0; i < len(lambdaFunctions); i++ {
		got := lambdaInfo[i].FunctionName
		want := lambdaFunctions[i]
		if got != want {
			t.Errorf("Got %v, want %v", got, want)
		}
	}
}
