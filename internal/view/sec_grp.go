package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type SG struct {
	ResourceViewer
}

// NewPod returns a new viewer.
func NewSG(resource string) ResourceViewer {
	var sg SG
	sg.ResourceViewer = NewBrowser(resource)
	sg.AddBindKeysFn(sg.bindKeys)
	return &sg
}

func (sg SG) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Group-Id", sg.GetTable().SortColCmd("Group-Id", true), true),
		ui.KeyShiftN:    ui.NewKeyAction("Sort Group-Name", sg.GetTable().SortColCmd("Group-Name", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", sg.App().PrevCmd, true),
	})
}
