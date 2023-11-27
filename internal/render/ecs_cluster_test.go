package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestEcsClusterRender(t *testing.T) {
	resp := aws.EcsClusterResp{ClusterName: "ClusterOne", Status: "ACTIVE", ClusterArn: "arn:aws:ecs:eu-central-1:012345678901:cluster/ClusterOne", RunningTasksCount: "1"}
	var ecsCluster EcsClusters

	r := NewRow(2)
	err := ecsCluster.Render(resp, "ecsCluster", &r)

	assert.Nil(t, err)
	assert.Equal(t, "ecsCluster", r.ID)

	e := Fields{"ClusterOne", "ACTIVE", "1", "arn:aws:ecs:eu-central-1:012345678901:cluster/ClusterOne"}
	assert.Equal(t, e, r.Fields[0:])

	headers := ecsCluster.Header()
	assert.Equal(t, 0, headers.IndexOf("Name", false))
	assert.Equal(t, 1, headers.IndexOf("Status", false))
	assert.Equal(t, 2, headers.IndexOf("RunningTasksCount", false))
	assert.Equal(t, 3, headers.IndexOf("Arn", false))
}
