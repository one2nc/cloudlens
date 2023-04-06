package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestLambdaRender(t *testing.T) {
	resp := aws.LambdaResp{FunctionName: "lambda-func-1", Description: "", Role: "arn:aws:iam:000000000000:role/Andre", FunctionArn: "arn:aws:lambda:us-east-1:000000000000:function:lambda-func-1", CodeSize: "861", LastModified: "9:00:00"}
	var lambda Lambda

	r := NewRow(6)
	err := lambda.Render(resp, "lambda", &r)

	assert.Nil(t, err)
	assert.Equal(t, "lambda", r.ID)

	e := Fields{"lambda-func-1", "", "arn:aws:iam:000000000000:role/Andre", "arn:aws:lambda:us-east-1:000000000000:function:lambda-func-1", "861", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])

	headers := lambda.Header()

	assert.Equal(t, 0, headers.IndexOf("Function-Name", false))
	assert.Equal(t, 1, headers.IndexOf("Description", false))
	assert.Equal(t, 2, headers.IndexOf("Role", false))
	assert.Equal(t, 3, headers.IndexOf("Function-Arn", false))
	assert.Equal(t, 4, headers.IndexOf("Code-Size", false))
	assert.Equal(t, 5, headers.IndexOf("Last-Modified", false))
}
