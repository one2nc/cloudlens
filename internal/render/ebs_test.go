package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestEBSRender(t *testing.T) {
	pom := aws.EBSResp{VolumeId: "vol-ebs-1", Size: "32", VolumeType: "gp2", State: "in-use", AvailabilityZone: "us-east-1e", Snapshot: "snapshot", CreationTime: "9:00:00"}

	var ebs EBS
	r := NewRow(7)
	err := ebs.Render(pom, "ebs", &r)
	assert.Nil(t, err)

	assert.Equal(t, "ebs", r.ID)
	e := Fields{"vol-ebs-1", "32", "gp2", "in-use", "us-east-1e", "snapshot", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])
}
