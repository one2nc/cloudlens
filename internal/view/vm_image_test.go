package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVMI(t *testing.T) {
	vmi := NewVMI("vmi")
	assert.Nil(t, vmi.Init(makeCtx()))
	assert.Equal(t, "vmi", vmi.Name())
	assert.Equal(t, 6, len(vmi.Hints()))
}
