package view

import (
	"context"
	"testing"

	"github.com/one2nc/cloudlens/internal"
	"github.com/stretchr/testify/assert"
)

func TestNewS3(t *testing.T) {
	s3 := NewS3("s3")
	assert.Nil(t, s3.Init(makeCtx()))
	assert.Equal(t, "s3", s3.Name())
	assert.Equal(t, 8, len(s3.Hints()))
}

func makeCtx() context.Context {
	return context.WithValue(context.Background(), internal.KeyApp, NewApp())
}
