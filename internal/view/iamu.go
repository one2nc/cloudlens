package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type IAMU struct {
	ResourceViewer
}

// NewSG returns a new viewer.
func NewIAMU(resource string) ResourceViewer {
	var iamu IAMU
	iamu.ResourceViewer = NewBrowser(resource)
	iamu.AddBindKeysFn(iamu.bindKeys)
	return &iamu
}

func (iamu IAMU) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		// ui.KeyShiftI:    ui.NewKeyAction("Sort ", iamu.GetTable().SortColCmd("Group-Id", true), true),
		// ui.KeyShiftN:    ui.NewKeyAction("Sort Group-Name", iamu.GetTable().SortColCmd("Group-Name", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", iamu.App().PrevCmd, true),
	})
}

