package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIamRole(t *testing.T) {
	iamr := NewIamRole("iam:r")
	assert.Nil(t, iamr.Init(makeCtx()))
	assert.Equal(t, "iam:r", iamr.Name())
	assert.Equal(t, 8, len(iamr.Hints()))
}
