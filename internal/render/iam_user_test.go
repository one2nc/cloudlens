package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestIamUserRender(t *testing.T) {
	pom := aws.IAMUSerResp{UserId: "iam-user-1", UserName: "Erdman", ARN: "arn:aws:iam:00000000000:user/Erdman", CreationTime: "9:00:00"}

	var iamU IAMU
	r := NewRow(4)
	err := iamU.Render(pom, "iamU", &r)
	assert.Nil(t, err)

	assert.Equal(t, "iamU", r.ID)
	e := Fields{"iam-user-1", "Erdman", "arn:aws:iam:00000000000:user/Erdman", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])
}
