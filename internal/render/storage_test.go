package render

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/stretchr/testify/assert"
)

func TestStorageRender(t *testing.T) {
	resp := gcp.StorageResp{BucketName: "test-bucket-1"}
	var s Storage

	r := NewRow(1)
	err := s.Render(resp, "storage", &r)

	assert.Nil(t, err)
	assert.Equal(t, "storage", r.ID)

	e := Fields{"test-bucket-1"}
	assert.Equal(t, e, r.Fields[:1])

	headers := s.Header()

	assert.Equal(t, 0, headers.IndexOf("Bucket-Name", false))
	assert.Equal(t, 1, headers.IndexOf("Creation-Time", false))
}
