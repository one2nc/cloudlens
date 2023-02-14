package aws

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func GetAllQueues(sess session.Session) ([]QueueResp, error) {
	queueResp := []QueueResp{}
	sqsServ := *sqs.New(&sess)
	res, err := sqsServ.ListQueues(nil)
	if err != nil {
		fmt.Println("Error in fetching all queues,", " err: ", err)
		return nil, err
	}

	for _, qUrl := range res.QueueUrls {
		fmt.Println("queue url is:", *qUrl)
		qA := strings.Split(*qUrl, "/")
		qName := qA[len(qA)-1]
		qAttributes, err := sqsServ.GetQueueAttributes(&sqs.GetQueueAttributesInput{
			AttributeNames: aws.StringSlice([]string{"All"}),
			QueueUrl:       qUrl,
		})
		if err != nil {
			fmt.Println("Error in fetching queue attributes", " err: ", err)
			return nil, err
		}
		mp := qAttributes.Attributes
		qR := QueueResp{
			Name:              qName,
			URL:               *qUrl,
			Type:              *mp["QueueArn"],
			Created:           *mp["CreatedTimestamp"],
			MessagesAvailable: *mp["ApproximateNumberOfMessages"],
			Encryption:        *mp["SqsManagedSseEnabled"],
			MaxMessageSize:    *mp["MaximumMessageSize"],
		}
		queueResp = append(queueResp, qR)
	}
	fmt.Println("qRespp is:", queueResp)
	return queueResp, nil
}

func GetMessageFromQueue(sess session.Session, queueUrl string) (*sqs.ReceiveMessageOutput, error) {
	sqsServ := *sqs.New(&sess)
	result, err := sqsServ.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(100),
	})
	if err != nil {
		fmt.Println("Error in fetching queue attributes", " err: ", err)
		return nil, err
	}
	return result, nil
}

func CreateLambdaFunction(sess session.Session) {
	lambdaServ := lambda.New(&sess)
	codeBytes, err := ioutil.ReadFile("code.zip")
	if err != nil {
		log.Fatalf("Failed to read function code: %v", err)
	}
	fmt.Println("codebytes are:", codeBytes)
	params := &lambda.CreateFunctionInput{
		FunctionName: aws.String("lambdaaaTemp1"),
		//Runtime:      aws.String(*aws.String("python3.7")),
		//Runtime: aws.String("go1.x"),
		Role: aws.String("arn:aws:iam::000000000000:role/Andre"),
		Code: &lambda.FunctionCode{
			ZipFile: codeBytes,
		},
	}
	resp, err := lambdaServ.CreateFunction(params)
	if err != nil {
		log.Fatalf("Failed to create function: %v", err)
	}

	fmt.Printf("Function ARN: %s\n", aws.StringValue(resp.FunctionArn))
}

func GetAllLambdaFunctions(sess session.Session) {
	lambdaServ := lambda.New(&sess)
	response, err := lambdaServ.ListFunctions(nil)
	if err != nil {
		fmt.Println("Error getting list of functions: ", err)
	}
	fmt.Println("response is:", response)
}
