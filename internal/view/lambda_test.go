package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLambda(t *testing.T) {
	lambda := NewLambda("lambda")
	assert.Nil(t, lambda.Init(makeCtx()))
	assert.Equal(t, "lambda", lambda.Name())
	assert.Equal(t, 9, len(lambda.Hints()))
}
