package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVpc(t *testing.T) {
	vpc := NewVPC("vpc")
	assert.Nil(t, vpc.Init(makeCtx()))
	assert.Equal(t, "vpc", vpc.Name())
	assert.Equal(t, 8, len(vpc.Hints()))
}
