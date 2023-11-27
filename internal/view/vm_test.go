package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVM(t *testing.T) {
	vm := NewVM("vm")
	assert.Nil(t, vm.Init(makeCtx()))
	assert.Equal(t, "vm", vm.Name())
	assert.Equal(t, 9, len(vm.Hints()))
}
