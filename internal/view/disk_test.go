package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDisk(t *testing.T) {
	disk := NewDisk("disk")
	assert.Nil(t, disk.Init(makeCtx()))
	assert.Equal(t, "disk", disk.Name())
	assert.Equal(t, 10, len(disk.Hints()))
}
