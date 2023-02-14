package aws

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/rs/zerolog/log"
)

func GetAllLambdaFunctions(sess session.Session) ([]LambdaResp, error) {
	responseA := []LambdaResp{}
	lambdaServ := lambda.New(&sess)
	response, err := lambdaServ.ListFunctions(nil)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error getting Lambda functions : %v", err))
		return nil, err
	}
	for _, r := range response.Functions {
		lastModTime, err := time.Parse("2006-01-02T15:04:05.999-0700", *r.LastModified)
		if err != nil {
			log.Info().Msg(fmt.Sprintf("error in converting 8601 %v", err))
			return nil, err
		}
		loc, _ := time.LoadLocation("Asia/Kolkata")
		t := lastModTime.In(loc)
		IST := t.In(loc)
		log.Info().Msg(fmt.Sprintf("IST IS %v", IST))
		lr := LambdaResp{
			FunctionName: *r.FunctionName,
			Description:  *r.Description,
			Role:         *r.Role,
			FunctionArn:  *r.FunctionArn,
			CodeSize:     strconv.Itoa(int(*r.CodeSize)),
			LastModified: IST.Format("Mon Jan _2 15:04:05 2006"),
		}
		responseA = append(responseA, lr)
	}
	return responseA, nil
}
