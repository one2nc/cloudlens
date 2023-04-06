package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestSecGrpRender(t *testing.T) {
	resp := aws.SGResp{GroupId: "sg-1", GroupName: "default", Description: "DefaultGroup", OwnerId: "000000000000", VpcId: "vpc-1"}
	var sg SG

	r := NewRow(5)
	err := sg.Render(resp, "sg", &r)

	assert.Nil(t, err)
	assert.Equal(t, "sg", r.ID)

	e := Fields{"sg-1", "default", "DefaultGroup", "000000000000", "vpc-1"}
	assert.Equal(t, e, r.Fields[0:])

	headers := sg.Header()

	assert.Equal(t, 0, headers.IndexOf("Group-Id", false))
	assert.Equal(t, 1, headers.IndexOf("Group-Name", false))
	assert.Equal(t, 2, headers.IndexOf("Description", false))
	assert.Equal(t, 3, headers.IndexOf("Owner-Id", false))
	assert.Equal(t, 4, headers.IndexOf("VPC-Id", false))
}
