package ui

import (
	"github.com/derailed/tview"
)

type App struct {
	*tview.Application
	Main  *Pages
	views map[string]tview.Primitive
}

func NewApp() *App {
	a := App{
		Application: tview.NewApplication(),
		Main:        NewPages(),
	}
	return &a
}

func (a *App) Init() {
	a.SetRoot(a.Main, true).EnableMouse(true)
}
