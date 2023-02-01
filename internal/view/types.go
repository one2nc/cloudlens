package view

import (
	"context"

	"github.com/one2nc/cloud-lens/internal/model"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type (
	// EnterFunc represents an enter key action.
	EnterFunc func(app *App, model ui.Tabular, resource string)

	// ContextFunc enhances a given context.
	ContextFunc func(context.Context) context.Context

	// BindKeysFunc adds new menu actions.
	BindKeysFunc func(ui.KeyActions)
)

// Viewer represents a component viewer.
type Viewer interface {
	model.Component

	// Actions returns active menu bindings.
	Actions() ui.KeyActions

	// App returns an app handle.
	App() *App

	// Refresh updates the viewer
	Refresh()
}

// TableViewer represents a tabular viewer.
type TableViewer interface {
	Viewer

	// Table returns a table component.
	GetTable() *Table
}

// ResourceViewer represents a generic resource viewer.
type ResourceViewer interface {
	TableViewer

	Resource() string

	// SetContextFn provision a custom context.
	SetContext(context.Context)

	// AddBindKeys provision additional key bindings.
	AddBindKeysFn(BindKeysFunc)
}
