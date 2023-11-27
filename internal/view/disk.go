package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type Disk struct {
	ResourceViewer
}

func NewDisk(resource string) ResourceViewer {
	var disk Disk
	disk.ResourceViewer = NewBrowser(resource)
	disk.AddBindKeysFn(disk.bindKeys)
	// s3.GetTable().SetEnterFn(s3.describeInstace)
	return &disk
}
func (disk *Disk) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Volume-Id", disk.GetTable().SortColCmd("Volume-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Size", disk.GetTable().SortColCmd("Size", true), true),
		ui.KeyShiftV:    ui.NewKeyAction("Sort Volume-Type", disk.GetTable().SortColCmd("Volume-Type", true), true),
		ui.KeyShiftZ:    ui.NewKeyAction("Sort Availability-Zone", disk.GetTable().SortColCmd("Availability-Zone", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Creation-Time", disk.GetTable().SortColCmd("Creation-Time", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", disk.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", disk.enterCmd, false),
	})
}

func (disk *Disk) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	volId := disk.GetTable().GetSelectedItem()
	disk.App().Flash().Info("volume id: " + volId)

	f := describeResource
	if disk.GetTable().enterFn != nil {
		f = disk.GetTable().enterFn
	}
	f(disk.App(), disk.GetTable().GetModel(), disk.Resource(), volId)
	return nil
}
