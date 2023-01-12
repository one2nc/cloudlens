package ui

import (
	"github.com/derailed/tview"
)

type Pages struct {
	*tview.Pages
}

func NewPages() *Pages {
	pages := Pages{Pages: tview.NewPages()}
	return &pages
}
