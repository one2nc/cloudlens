package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEc2Image(t *testing.T) {
	ec2i := NewEC2I("ec2:i")
	assert.Nil(t, ec2i.Init(makeCtx()))
	assert.Equal(t, "ec2:i", ec2i.Name())
	assert.Equal(t, 6, len(ec2i.Hints()))
}
