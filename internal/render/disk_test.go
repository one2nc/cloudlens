package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/stretchr/testify/assert"
)

func TestDiskRender(t *testing.T) {
	resp := gcp.DiskResp{Name: "disk-1", Size: "32", Type: "pd-balanced", Status: "READY", Zone: "asia-south1-c", CreationTime: "9:00:00"}
	var disk Disk

	r := NewRow(6)
	err := disk.Render(resp, "disk", &r)

	assert.Nil(t, err)
	assert.Equal(t, "disk", r.ID)

	e := Fields{"disk-1", "32", "pd-balanced", "READY", "asia-south1-c", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])

	headers := disk.Header()
	assert.Equal(t, 0, headers.IndexOf("Volume-Id", false))
	assert.Equal(t, 1, headers.IndexOf("Size", false))
	assert.Equal(t, 2, headers.IndexOf("Volume-Type", false))
	assert.Equal(t, 3, headers.IndexOf("Status", false))
	assert.Equal(t, 4, headers.IndexOf("Availability-Zone", false))
	assert.Equal(t, 5, headers.IndexOf("Creation-Time", false))
}
