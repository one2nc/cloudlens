package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEc2Snapshot(t *testing.T) {
	ec2s := NewEC2S("ec2:s")
	assert.Nil(t, ec2s.Init(makeCtx()))
	assert.Equal(t, "ec2:s", ec2s.Name())
	assert.Equal(t, 8, len(ec2s.Hints()))
}
