package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestSubnetRender(t *testing.T) {
	resp := aws.SubnetResp{SubnetId: "subnet-1", OwnerId: "000000000000", CidrBlock: "172.31.0.0/16", AvailabilityZone: "us-east-1", State: "disabled"}
	var subnet Subnet

	r := NewRow(5)
	err := subnet.Render(resp, "subnet", &r)

	assert.Nil(t, err)
	assert.Equal(t, "subnet", r.ID)

	e := Fields{"subnet-1", "000000000000", "172.31.0.0/16", "us-east-1", "disabled"}
	assert.Equal(t, e, r.Fields[0:])

	headers := subnet.Header()
	assert.Equal(t, 0, headers.IndexOf("Subnet-Id", false))
	assert.Equal(t, 1, headers.IndexOf("Owner-Id", false))
	assert.Equal(t, 2, headers.IndexOf("Cidr Block", false))
	assert.Equal(t, 3, headers.IndexOf("Availability Zone", false))
	assert.Equal(t, 4, headers.IndexOf("Subnet-State", false))
}
