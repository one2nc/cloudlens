package view

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
)

const (
	DefaultRefreshRate = time.Second * 20
)

// Table represents a table viewer.
type Table struct {
	*ui.Table

	app        *App
	enterFn    EnterFunc
	bindKeysFn []BindKeysFunc
}

// NewTable returns a new viewer.
func NewTable(res string) *Table {
	t := Table{
		Table: ui.NewTable(res),
	}
	return &t
}

// Init initializes the component.
func (t *Table) Init(ctx context.Context) (err error) {
	if t.app, err = extractApp(ctx); err != nil {
		return err
	}

	t.Table.Init(ctx)
	t.SetInputCapture(t.keyboard)
	t.bindKeys()
	t.GetModel().SetRefreshRate(DefaultRefreshRate)

	return nil
}

// App returns the current app handle.
func (t *Table) App() *App {
	return t.app
}

// Start runs the component.
func (t *Table) Start() {
}

// Stop terminates the component.
func (t *Table) Stop() {
}

// SetEnterFn specifies the default enter behavior.
func (t *Table) SetEnterFn(f EnterFunc) {
	t.enterFn = f
}

func (t *Table) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	key := evt.Key()
	if key == tcell.KeyUp || key == tcell.KeyDown {
		return evt
	}

	if a, ok := t.Actions()[ui.AsKey(evt)]; ok {
		return a.Action(evt)
	}

	return evt
}

func (t *Table) bindKeys() {
	t.Actions().Add(ui.KeyActions{
		tcell.KeyCtrlW: ui.NewKeyAction("Toggle Wide", t.toggleWideCmd, true),
		ui.KeyHelp:     ui.NewKeyAction("Help", t.App().helpCmd, true),
	})
}

// Name returns the table name.
func (t *Table) Name() string { return t.Table.Resource() }

// AddBindKeysFn adds additional key bindings.
func (t *Table) AddBindKeysFn(f BindKeysFunc) {
	t.bindKeysFn = append(t.bindKeysFn, f)
}

func (t *Table) toggleWideCmd(evt *tcell.EventKey) *tcell.EventKey {
	t.ToggleWide()
	return nil
}
