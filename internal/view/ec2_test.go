package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEc2(t *testing.T) {
	ec2 := NewEC2("ec2")
	assert.Nil(t, ec2.Init(makeCtx()))
	assert.Equal(t, "ec2", ec2.Name())
	assert.Equal(t, 11, len(ec2.Hints()))
}
