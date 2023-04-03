package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestS3Render(t *testing.T) {
	pom := aws.BucketResp{
		BucketName: "test-bucket-1",
	}

	var s3 S3
	r := NewRow(1)
	err := s3.Render(pom, "s3", &r)
	assert.Nil(t, err)

	assert.Equal(t, "s3", r.ID)
	e := Fields{"test-bucket-1"}
	assert.Equal(t, e, r.Fields[:1])
}
