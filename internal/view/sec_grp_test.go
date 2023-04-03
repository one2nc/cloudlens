package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecGrp(t *testing.T) {
	sg := NewSG("sg")
	assert.Nil(t, sg.Init(makeCtx()))
	assert.Equal(t, "sg", sg.Name())
	assert.Equal(t, 7, len(sg.Hints()))
}
