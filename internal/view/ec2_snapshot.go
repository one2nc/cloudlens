package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type EC2S struct {
	ResourceViewer
}

func NewEC2S(resource string) ResourceViewer {
	var es EC2S
	es.ResourceViewer = NewBrowser(resource)
	es.AddBindKeysFn(es.bindKeys)
	return &es
}

func (es *EC2S) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Snapshot-Id", es.GetTable().SortColCmd("Snapshot-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Volume-Size", es.GetTable().SortColCmd("Volume-Size", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Start-Time", es.GetTable().SortColCmd("Start-Time", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", es.App().PrevCmd, true),
		tcell.KeyEnter:  ui.NewKeyAction("View", es.enterCmd, true),
	})
}

func (es *EC2S) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	snapshotId := es.GetTable().GetSelectedItem()
	if snapshotId != "" {
		f := describeResource
		if es.GetTable().enterFn != nil {
			f = es.GetTable().enterFn
		}
		f(es.App(), es.GetTable().GetModel(), es.Resource(), snapshotId)
		es.App().Flash().Info("Snapshot-Id: " + snapshotId)
	}
	return nil
}
