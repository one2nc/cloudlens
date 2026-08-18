package view

import (
	"context"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/ui"
	"github.com/stretchr/testify/assert"
)

func newTestBrowser(t *testing.T) *Browser {
	app := NewApp()
	ctx := context.WithValue(context.Background(), internal.KeyApp, app)
	ctx = context.WithValue(ctx, internal.KeySelectedCloud, internal.AWS)
	app.SetContext(ctx)
	app.UpdateContext(ctx)

	b := NewBrowser("ec2").(*Browser)
	assert.Nil(t, b.Init(ctx))
	return b
}

func TestFilterKeysAreActuallyBound(t *testing.T) {
	b := newTestBrowser(t)

	actions := b.Actions()

	assert.Equal(t, "Filter Mode", actions[ui.KeySlash].Description)
	assert.Equal(t, "Filter Reset", actions[tcell.KeyEscape].Description)
}

func TestSlashKeyOpensFilterWhenNotAlreadyFiltering(t *testing.T) {
	b := newTestBrowser(t)

	b.activateCmd(nil)

	assert.True(t, b.CmdBuff().IsActive())
}

func TestEscapeClearsFilterTextAndTurnsOffFilterBox(t *testing.T) {
	b := newTestBrowser(t)
	b.CmdBuff().SetActive(true)
	b.CmdBuff().Add('w')

	b.resetCmd(nil)

	assert.Equal(t, "", b.CmdBuff().GetText())
	assert.False(t, b.CmdBuff().IsActive())
}

func TestEscapeDoesNothingHarmfulWhenNoFilterIsActive(t *testing.T) {
	b := newTestBrowser(t)

	assert.NotPanics(t, func() {
		b.resetCmd(nil)
	})
	assert.Equal(t, "", b.CmdBuff().GetText())
}
