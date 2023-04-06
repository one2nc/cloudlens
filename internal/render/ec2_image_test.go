package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestEc2ImageRender(t *testing.T) {
	resp := aws.ImageResp{ImageId: "image-1", OwnerId: "8011", ImageLocation: "amazon/getting-started", Name: "Windows_server_2016", ImageType: "machine"}
	var ec2i EC2I

	r := NewRow(5)
	err := ec2i.Render(resp, "ec2i", &r)

	assert.Nil(t, err)
	assert.Equal(t, "ec2i", r.ID)

	e := Fields{"image-1", "8011", "amazon/getting-started", "Windows_server_2016", "machine"}
	assert.Equal(t, e, r.Fields[0:])

	headers := ec2i.Header()
	assert.Equal(t, 0, headers.IndexOf("Image-Id", false))
	assert.Equal(t, 1, headers.IndexOf("Owner-Id", false))
	assert.Equal(t, 2, headers.IndexOf("Image-Location", false))
	assert.Equal(t, 3, headers.IndexOf("Name", false))
	assert.Equal(t, 4, headers.IndexOf("Image-Type", false))
}
