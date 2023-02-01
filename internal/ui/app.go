package ui

import (
	"os"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/model"
)

type App struct {
	*tview.Application
	Main    *Pages
	flash   *model.Flash
	actions KeyActions
	views   map[string]tview.Primitive
}

func NewApp() *App {
	a := App{
		Application: tview.NewApplication(),
		Main:        NewPages(),
		flash:       model.NewFlash(model.DefaultFlashDelay),
		actions:     make(KeyActions),
		views:       make(map[string]tview.Primitive),
	}
	a.views["menu"] = NewMenu()
	return &a
}

func (a *App) Init() {
	a.bindKeys()
	a.SetRoot(a.Main, true).EnableMouse(true)
}

// HasAction checks if key matches a registered binding.
func (a *App) HasAction(key tcell.Key) (KeyAction, bool) {
	act, ok := a.actions[key]
	return act, ok
}

// GetActions returns a collection of actions.
func (a *App) GetActions() KeyActions {
	return a.actions
}

// AddActions returns the application actions.
func (a *App) AddActions(aa KeyActions) {
	for k, v := range aa {
		a.actions[k] = v
	}
}

func (a *App) bindKeys() {
	a.actions = KeyActions{
		tcell.KeyCtrlR: NewKeyAction("Redraw", a.redrawCmd, false),
		tcell.KeyCtrlC: NewKeyAction("Quit", a.quitCmd, false),
	}
}

// RedrawCmd forces a redraw.
func (a *App) redrawCmd(evt *tcell.EventKey) *tcell.EventKey {
	a.QueueUpdateDraw(func() {})
	return evt
}

func (a *App) quitCmd(evt *tcell.EventKey) *tcell.EventKey {
	a.BailOut()
	// overwrite the default ctrl-c behavior of tview
	return nil
}

// BailOut exits the application.
func (a *App) BailOut() {
	a.Stop()
	os.Exit(0)
}

// Views return the application root views.
func (a *App) Views() map[string]tview.Primitive {
	return a.views
}

// View Accessors...
func (a *App) Menu() *Menu {
	return a.views["menu"].(*Menu)
}

func (a *App) FlashView() *Flash {
	return a.views["flash"].(*Flash)
}

// Flash returns a flash model.
func (a *App) Flash() *model.Flash {
	return a.flash
}

// AsKey converts rune to keyboard key.,.
func AsKey(evt *tcell.EventKey) tcell.Key {
	if evt.Key() != tcell.KeyRune {
		return evt.Key()
	}
	key := tcell.Key(evt.Rune())
	if evt.Modifiers() == tcell.ModAlt {
		key = tcell.Key(int16(evt.Rune()) * int16(evt.Modifiers()))
	}
	return key
}
