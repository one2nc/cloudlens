package aws

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/rs/zerolog/log"
)

func GetAllLambdaFunctions(cfg awsV2.Config) ([]LambdaResp, error) {
	responseA := []LambdaResp{}
	lambdaServ := lambda.NewFromConfig(cfg)
	response, err := lambdaServ.ListFunctions(context.Background(), nil)
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
		localZone, err := config.GetLocalTimeZone() // Empty string loads the local timezone
		if err != nil {
			fmt.Println("Error loading local timezone:", err)
			return nil, err
		}
		loc, _ := time.LoadLocation(localZone)
		t := lastModTime.In(loc)
		IST := t.In(loc)
		log.Info().Msg(fmt.Sprintf("IST IS %v", IST))
		lr := LambdaResp{
			FunctionName: *r.FunctionName,
			Description:  *r.Description,
			Role:         *r.Role,
			FunctionArn:  *r.FunctionArn,
			CodeSize:     strconv.Itoa(int(r.CodeSize)),
			LastModified: IST.Format("Mon Jan _2 15:04:05 2006"),
		}
		responseA = append(responseA, lr)
	}
	return responseA, nil
}
