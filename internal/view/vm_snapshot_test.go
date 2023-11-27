package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVMS(t *testing.T) {
	vms := NewVMS("vms")
	assert.Nil(t, vms.Init(makeCtx()))
	assert.Equal(t, "vms", vms.Name())
	assert.Equal(t, 8, len(vms.Hints()))
}