package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIamUser(t *testing.T) {
	iamu := NewIAMU("iam:u")
	assert.Nil(t, iamu.Init(makeCtx()))
	assert.Equal(t, "iam:u", iamu.Name())
	assert.Equal(t, 8, len(iamu.Hints()))
}
