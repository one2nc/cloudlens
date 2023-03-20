package aws

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	sqss "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
)

func GetAllQueues(cfg awsV2.Config) ([]SQSResp, error) {
	queueResp := []SQSResp{}
	sqsServ := *sqss.NewFromConfig(cfg)
	res, err := sqsServ.ListQueues(context.Background(), nil)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching all queues, err: %v", err))
		return nil, err
	}

	for _, qUrl := range res.QueueUrls {
		qA := strings.Split(qUrl, "/")
		qName := qA[len(qA)-1]
		qAttributes, err := sqsServ.GetQueueAttributes(context.Background(), &sqss.GetQueueAttributesInput{
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
			QueueUrl:       &qUrl,
		})
		if err != nil {
			log.Info().Msg(fmt.Sprintf("Error in fetching queue attributes: %v", err))
			return nil, err
		}
		mp := qAttributes.Attributes
		launchTime, _ := strconv.Atoi(mp["CreatedTimestamp"])
		localZone, err := GetLocalTimeZone() // Empty string loads the local timezone
		if err != nil {
			fmt.Println("Error loading local timezone:", err)
			return nil, err
		}
		loc, _ := time.LoadLocation(localZone)
		IST := time.Unix(int64(launchTime), 0).In(loc)
		qR := SQSResp{
			Name:              qName,
			URL:               qUrl,
			Type:              mp["QueueArn"],
			Created:           IST.Format("Mon Jan _2 15:04:05 2006"),
			MessagesAvailable: mp["ApproximateNumberOfMessages"],
			Encryption:        mp["SqsManagedSseEnabled"],
			MaxMessageSize:    mp["MaximumMessageSize"],
		}
		queueResp = append(queueResp, qR)
	}
	return queueResp, nil
}

func GetMessageFromQueue(sess session.Session, queueUrl string) (*sqs.ReceiveMessageOutput, error) {
	sqsServ := *sqs.New(&sess)
	result, err := sqsServ.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(100),
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching queue attributes : %v", err))
		return nil, err
	}
	return result, nil
}
