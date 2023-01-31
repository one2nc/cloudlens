package view

import "github.com/one2nc/cloud-lens/internal/ui"

type (
	// EnterFunc represents an enter key action.
	EnterFunc func(app *App, model ui.Tabular, resource string)

	// BindKeysFunc adds new menu actions.
	BindKeysFunc func(ui.KeyActions)
)
