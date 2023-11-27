package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type VMS struct {
	ResourceViewer
}

func NewVMS(resource string) ResourceViewer {
	var vms VMS
	vms.ResourceViewer = NewBrowser(resource)
	vms.AddBindKeysFn(vms.bindKeys)
	return &vms
}

func (vms *VMS) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Snapshot-Id", vms.GetTable().SortColCmd("Snapshot-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Volume-Size", vms.GetTable().SortColCmd("Volume-Size", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Start-Time", vms.GetTable().SortColCmd("Start-Time", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", vms.App().PrevCmd, true),
		tcell.KeyEnter:  ui.NewKeyAction("View", vms.enterCmd, true),
	})
}

func (vms *VMS) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	snapshotId := vms.GetTable().GetSelectedItem()
	if snapshotId != "" {
		f := describeResource
		if vms.GetTable().enterFn != nil {
			f = vms.GetTable().enterFn
		}
		f(vms.App(), vms.GetTable().GetModel(), vms.Resource(), snapshotId)
		vms.App().Flash().Info("Snapshot-Id: " + snapshotId)
	}
	return nil
}
