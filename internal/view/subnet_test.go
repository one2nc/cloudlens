package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSubnet(t *testing.T) {
	subnet := NewSubnet("subnet")
	assert.Nil(t, subnet.Init(makeCtx()))
	assert.Equal(t, "subnet", subnet.Name())
	assert.Equal(t, 7, len(subnet.Hints()))
}
