package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSQS(t *testing.T) {
	sqs := NewSQS("sqs")
	assert.Nil(t, sqs.Init(makeCtx()))
	assert.Equal(t, "sqs", sqs.Name())
	assert.Equal(t, 9, len(sqs.Hints()))
}
