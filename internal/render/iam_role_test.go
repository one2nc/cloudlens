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

	headers := iamRole.Header()
	assert.Equal(t, 0, headers.IndexOf("Role-Id", false))
	assert.Equal(t, 1, headers.IndexOf("Role-Name", false))
	assert.Equal(t, 2, headers.IndexOf("ARN", false))
	assert.Equal(t, 3, headers.IndexOf("Created-Date", false))
}
