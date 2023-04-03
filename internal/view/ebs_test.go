package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEbs(t *testing.T) {
	ebs := NewEBS("ebs")
	assert.Nil(t, ebs.Init(makeCtx()))
	assert.Equal(t, "ebs", ebs.Name())
	assert.Equal(t, 10, len(ebs.Hints()))
}
