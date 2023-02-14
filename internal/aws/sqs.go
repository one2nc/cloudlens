package aws

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
)

func GetAllQueues(sess session.Session) ([]SQSResp, error) {
	queueResp := []SQSResp{}
	sqsServ := *sqs.New(&sess)
	res, err := sqsServ.ListQueues(nil)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching all queues, err: %v", err))
		return nil, err
	}

	for _, qUrl := range res.QueueUrls {
		qA := strings.Split(*qUrl, "/")
		qName := qA[len(qA)-1]
		qAttributes, err := sqsServ.GetQueueAttributes(&sqs.GetQueueAttributesInput{
			AttributeNames: aws.StringSlice([]string{"All"}),
			QueueUrl:       qUrl,
		})
		if err != nil {
			log.Info().Msg(fmt.Sprintf("Error in fetching queue attributes: %v", err))
			return nil, err
		}
		mp := qAttributes.Attributes
		launchTime, _ := strconv.Atoi(*mp["CreatedTimestamp"])
		loc, _ := time.LoadLocation("Asia/Kolkata")
		IST := time.Unix(int64(launchTime), 0).In(loc)
		qR := SQSResp{
			Name:              qName,
			URL:               *qUrl,
			Type:              *mp["QueueArn"],
			Created:           IST.Format("Mon Jan _2 15:04:05 2006"),
			MessagesAvailable: *mp["ApproximateNumberOfMessages"],
			Encryption:        *mp["SqsManagedSseEnabled"],
			MaxMessageSize:    *mp["MaximumMessageSize"],
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
