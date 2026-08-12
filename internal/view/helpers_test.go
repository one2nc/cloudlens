package view

import (
	"context"
	"testing"

	"github.com/one2nc/cloudlens/internal"
	"github.com/stretchr/testify/assert"
)

func TestDescribeResourcePushesLiveViewOntoContentStack(t *testing.T) {
	app := NewApp()
	ctx := context.WithValue(context.Background(), internal.KeyApp, app)
	app.SetContext(ctx)
	assert.Nil(t, app.Content.Init(ctx))

	describeResource(app, nil, "ec2", "i-123")

	top := app.Content.Current()
	liveView, ok := top.(*LiveView)
	assert.True(t, ok)
	assert.Equal(t, "Describe", liveView.Name())
}
