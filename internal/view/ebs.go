package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type EBS struct {
	ResourceViewer
}

func NewEBS(resource string) ResourceViewer {
	var ebs EBS
	ebs.ResourceViewer = NewBrowser(resource)
	ebs.AddBindKeysFn(ebs.bindKeys)
	// s3.GetTable().SetEnterFn(s3.describeInstace)
	return &ebs
}
func (ebs *EBS) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Volume-Id", ebs.GetTable().SortColCmd("Volume-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Size", ebs.GetTable().SortColCmd("Size", true), true),
		ui.KeyShiftV:    ui.NewKeyAction("Sort Volume-Type", ebs.GetTable().SortColCmd("Volume-Type", true), true),
		ui.KeyShiftZ:    ui.NewKeyAction("Sort Availability-Zone", ebs.GetTable().SortColCmd("Availability-Zone", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Creation-Time", ebs.GetTable().SortColCmd("Creation-Time", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ebs.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", ebs.enterCmd, false),
	})
}

func (ebs *EBS) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	volId := ebs.GetTable().GetSelectedItem()
	ebs.App().Flash().Info("volume id: " + volId)

	f := describeResource
	if ebs.GetTable().enterFn != nil {
		f = ebs.GetTable().enterFn
	}
	f(ebs.App(), ebs.GetTable().GetModel(), ebs.Resource(), volId)
	return nil
}
