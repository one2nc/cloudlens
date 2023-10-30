package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	s := NewStorage("storage")
	assert.Nil(t, s.Init(makeCtx()))
	assert.Equal(t, "storage", s.Name())
	assert.Equal(t, 7, len(s.Hints()))
}
