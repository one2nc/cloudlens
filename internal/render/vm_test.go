package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/stretchr/testify/assert"
)

func TestVMRender(t *testing.T) {
	resp := gcp.VMResp{InstanceId: "instance-1", InstanceState: "running", InstanceType: "e2-micro", LaunchTime: "9:00:00", AvailabilityZone: "asia-south1-c"}
	var vm VM

	r := NewRow(5)
	err := vm.Render(resp, "vm", &r)

	assert.Nil(t, err)
	assert.Equal(t, "vm", r.ID)

	e := Fields{"instance-1", "running", "e2-micro", "9:00:00", "asia-south1-c"}
	assert.Equal(t, e, r.Fields[0:])

	headers := vm.Header()
	assert.Equal(t, 0, headers.IndexOf("Instance-Id", false))
	assert.Equal(t, 1, headers.IndexOf("Instance-State", false))
	assert.Equal(t, 2, headers.IndexOf("Instance-Type", false))
	assert.Equal(t, 3, headers.IndexOf("Launch-Time", false))
	assert.Equal(t, 4, headers.IndexOf("Availability-Zone", false))
}
