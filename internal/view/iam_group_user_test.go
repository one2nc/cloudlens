package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIamGroupUser(t *testing.T) {
	iamU := NewIamGroupUser("iam:u")
	assert.Nil(t, iamU.Init(makeCtx()))
	assert.Equal(t, "iam:u", iamU.Name())
	assert.Equal(t, 7, len(iamU.Hints()))
}
