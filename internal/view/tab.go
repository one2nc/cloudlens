package view

import (
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
)

type Tab struct {
	*App
	items []tview.Primitive
}

func NewTab(app *App) *Tab {
	t := Tab{App: app, items: []tview.Primitive{}}
	t.Add()
	return &t
}

func (t *Tab) Add() {

	ctx := t.App.GetContext()

	cloud := ctx.Value(internal.KeySelectedCloud)

	switch cloud {
	case internal.AWS:
		t.items = append(t.items, t.App.profile())
		t.items = append(t.items, t.App.region())
	case internal.GCP:
		t.items = append(t.items, t.App.project())
	}
}

func (t *Tab) tabAction(event *tcell.EventKey) *tcell.EventKey {
	if t.InCmdMode() {
		return event
	}

	focusIdx := t.currentFocusIdx()

	if event.Key() == tcell.KeyTAB  {
		if focusIdx + 1 == len(t.items) {
			t.App.Application.SetFocus(t.Content.Pages.Current())
			return event
		}
		focusIdx = focusIdx + 1
	}
	if focusIdx < 0 {
		focusIdx = 0
	}
	t.App.Application.SetFocus(t.items[focusIdx])
	return event
}

func (t *Tab) currentFocusIdx() int {
	focusIdx := -1
	for i, p := range t.items {
		if p.HasFocus() {
			focusIdx = i
			break
		}
	}
	return focusIdx
}
