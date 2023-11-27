package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/stretchr/testify/assert"
)

func TestVMImageRender(t *testing.T) {
	resp := gcp.ImageResp{Name: "image-1", Location: "asia", Status: "READY", CreatedAt: "9:00:00"}
	var vmi VMI

	r := NewRow(4)
	err := vmi.Render(resp, "vmi", &r)

	assert.Nil(t, err)
	assert.Equal(t, "vmi", r.ID)

	e := Fields{"image-1", "asia", "READY", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])

	headers := vmi.Header()
	assert.Equal(t, 0, headers.IndexOf("Image-Id", false))

	assert.Equal(t, 1, headers.IndexOf("Image-Location", false))

	assert.Equal(t, 2, headers.IndexOf("Status", false))
	assert.Equal(t, 3, headers.IndexOf("Creation-Time", false))
}
