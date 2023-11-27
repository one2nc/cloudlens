package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type VMI struct {
	ResourceViewer
}

func NewVMI(resource string) ResourceViewer {
	var vmi VMI
	vmi.ResourceViewer = NewBrowser(resource)
	vmi.AddBindKeysFn(vmi.bindKeys)
	return &vmi
}

func (vmi *VMI) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftT:    ui.NewKeyAction("Sort Creation-Time", vmi.GetTable().SortColCmd("Creation-Time", true), true),
		ui.KeyShiftI:    ui.NewKeyAction("Sort Image-Id", vmi.GetTable().SortColCmd("Image-Id", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", vmi.App().PrevCmd, true),
		tcell.KeyEnter:  ui.NewKeyAction("View", vmi.enterCmd, true),
	})
}

func (vmi *VMI) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	imageId := vmi.GetTable().GetSelectedItem()
	if imageId != "" {
		f := describeResource
		if vmi.GetTable().enterFn != nil {
			f = vmi.GetTable().enterFn
		}
		f(vmi.App(), vmi.GetTable().GetModel(), vmi.Resource(), imageId)
		vmi.App().Flash().Info("Image-Id: " + imageId)
	}
	return nil
}
