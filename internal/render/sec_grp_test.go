package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestSecGrpRender(t *testing.T) {
	pom := aws.SGResp{GroupId: "sg-1", GroupName: "default", Description: "DefaultGroup", OwnerId: "000000000000", VpcId: "vpc-1"}

	var s3 SG
	r := NewRow(5)
	err := s3.Render(pom, "sg", &r)
	assert.Nil(t, err)

	assert.Equal(t, "sg", r.ID)
	e := Fields{"sg-1", "default", "DefaultGroup", "000000000000", "vpc-1"}
	assert.Equal(t, e, r.Fields[0:])
}
