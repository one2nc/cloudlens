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
		views: make(map[string]tview.Primitive),
	}
	a.views["menu"] = NewMenu()
	return &a
}

func (a *App) Init() {
	a.SetRoot(a.Main, true).EnableMouse(true)
}

// Views return the application root views.
func (a *App) Views() map[string]tview.Primitive {
	return a.views
}

// View Accessors...
func (a *App) Menu() *Menu {
	return a.views["menu"].(*Menu)
}
