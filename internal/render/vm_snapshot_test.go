package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/stretchr/testify/assert"
)

func TestVMSnapshotRender(t *testing.T) {
	resp := gcp.SnapshotResp{Name: "snap-1", Size: "10 GB", CreatedAt: "9:00:00"}
	var vms VMS

	r := NewRow(6)
	err := vms.Render(resp, "vms", &r)

	assert.Nil(t, err)
	assert.Equal(t, "vms", r.ID)

	e := Fields{"snap-1", "10 GB", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])

	headers := vms.Header()
	assert.Equal(t, 0, headers.IndexOf("Snapshot-Id", false))

	assert.Equal(t, 1, headers.IndexOf("Volume-Size", false))

	assert.Equal(t, 2, headers.IndexOf("Start-Time", false))
}
