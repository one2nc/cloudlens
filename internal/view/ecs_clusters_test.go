package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEcsCluster(t *testing.T) {
	ecs := NewEcs("ecs:c")
	assert.Nil(t, ecs.Init(makeCtx()))
	assert.Equal(t, "ecs:c", ecs.Name())
	assert.Equal(t, 5, len(ecs.Hints()))
}
