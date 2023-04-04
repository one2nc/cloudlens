package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestSQSRender(t *testing.T) {
	resp := aws.SQSResp{Name: "Queue-1", URL: "http://localhost:4566/000000000000", Type: "arn:aws:sqs:us-east-1:000000000000:P", Created: "9:00:00", MessagesAvailable: "10", Encryption: "ASE", MaxMessageSize: "128"}

	var sqs SQS
	r := NewRow(7)
	err := sqs.Render(resp, "sqs", &r)
	assert.Nil(t, err)

	assert.Equal(t, "sqs", r.ID)
	e := Fields{"http://localhost:4566/000000000000", "Queue-1", "arn:aws:sqs:us-east-1:000000000000:P", "9:00:00", "10", "ASE", "128"}
	assert.Equal(t, e, r.Fields[0:])
}
