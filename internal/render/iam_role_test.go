package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestIamRoleRender(t *testing.T) {
	resp := aws.IamRoleResp{RoleId: "iam-role-1", RoleName: "role", ARN: "arn:aws:iam:00000000000:role/role-1", CreationTime: "9:00:00"}

	var iamRole IamRole
	r := NewRow(4)
	err := iamRole.Render(resp, "iamRole", &r)
	assert.Nil(t, err)

	assert.Equal(t, "iamRole", r.ID)
	e := Fields{"iam-role-1", "role", "arn:aws:iam:00000000000:role/role-1", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])
}
