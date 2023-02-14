package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func GetAllLambdaFunctions(sess session.Session) {
	lambdaServ := lambda.New(&sess)
	response, err := lambdaServ.ListFunctions(nil)
	if err != nil {
		fmt.Println("Error getting list of functions: ", err)
	}
	fmt.Println("response is:", response)
}
