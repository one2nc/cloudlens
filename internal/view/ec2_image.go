package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type EC2I struct {
	ResourceViewer
}

func NewEC2I(resource string) ResourceViewer {
	var es EC2S
	es.ResourceViewer = NewBrowser(resource)
	es.AddBindKeysFn(es.bindKeys)
	return &es
}

func (ei *EC2I) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Snapshot-Id", ei.GetTable().SortColCmd("Snapshot-Id", true), true),
		ui.KeyShiftV:    ui.NewKeyAction("Sort Volume-Size", ei.GetTable().SortColCmd("Volume-Size", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Start-Time", ei.GetTable().SortColCmd("Start-Time", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ei.App().PrevCmd, true),
	})
}