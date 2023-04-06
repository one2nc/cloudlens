package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIamUserGroup(t *testing.T) {
	iamug := NewIAMUG("iam:g")
	assert.Nil(t, iamug.Init(makeCtx()))
	assert.Equal(t, "iam:g", iamug.Name())
	assert.Equal(t, 9, len(iamug.Hints()))
}
