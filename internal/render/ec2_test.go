package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestEc2Render(t *testing.T) {
	resp := aws.EC2Resp{InstanceId: "ec2-instance-1", InstanceState: "running", InstanceType: "t2.micro", MonitoringState: "disabled", PublicDNS: "public-dns", LaunchTime: "9:00:00", AvailabilityZone: "us-east-1e"}

	var ec2 EC2
	r := NewRow(7)
	err := ec2.Render(resp, "ec2", &r)
	assert.Nil(t, err)

	assert.Equal(t, "ec2", r.ID)
	e := Fields{"ec2-instance-1", "running", "t2.micro", "disabled", "9:00:00", "public-dns", "us-east-1e"}
	assert.Equal(t, e, r.Fields[0:])
}
