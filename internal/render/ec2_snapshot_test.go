package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestEc2SnapshotRender(t *testing.T) {
	pom := aws.Snapshot{SnapshotId: "snap-1", OwnerId: "8011", VolumeId: "vol-1", VolumeSize: "15", StartTime: "9:00:00", State: "completed"}

	var ec2s EC2S
	r := NewRow(6)
	err := ec2s.Render(pom, "ec2s", &r)
	assert.Nil(t, err)

	assert.Equal(t, "ec2s", r.ID)
	e := Fields{"snap-1", "8011", "vol-1", "15", "completed", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])
}
