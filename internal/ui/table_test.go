package ui

import (
	"testing"

	"github.com/one2nc/cloudlens/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewTableStartsWithFilterReady(t *testing.T) {
	tbl := NewTable("ec2")

	assert.NotNil(t, tbl.CmdBuff())
	assert.Equal(t, model.FilterBuffer, tbl.CmdBuff().GetKind())
}

func TestTypingDoesNothingWhenFilterNotActive(t *testing.T) {
	tbl := NewTable("ec2")

	got := tbl.FilterInput('w')

	assert.False(t, got)
	assert.Equal(t, "", tbl.CmdBuff().GetText())
}

func TestTypingUpdatesFilterTextWhenActive(t *testing.T) {
	tbl := NewTable("ec2")
	tbl.CmdBuff().SetActive(true)

	got := tbl.FilterInput('w')

	assert.True(t, got)
	assert.Equal(t, "w", tbl.CmdBuff().GetText())
}

func TestFilterAccessorReturnsSameFilter(t *testing.T) {
	tbl := NewTable("ec2")

	assert.Same(t, tbl.cmdBuff, tbl.CmdBuff())
}

func TestTitleShowsFilterTextWhenFiltering(t *testing.T) {
	tbl := NewTable("ec2")
	tbl.CmdBuff().SetActive(true)
	tbl.CmdBuff().Add('w')

	tbl.UpdateTitle()

	assert.Contains(t, tbl.GetTitle(), "(filter: w)")
}

func TestTitleOmitsFilterTextWhenNotFiltering(t *testing.T) {
	tbl := NewTable("ec2")

	tbl.UpdateTitle()

	assert.NotContains(t, tbl.GetTitle(), "(filter:")
}
