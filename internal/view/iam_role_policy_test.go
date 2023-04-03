package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIamRolePolicy(t *testing.T) {
	iamr := NewIamRolePloicy("iam:r")
	assert.Nil(t, iamr.Init(makeCtx()))
	assert.Equal(t, "iam:r", iamr.Name())
	assert.Equal(t, 6, len(iamr.Hints()))
}
