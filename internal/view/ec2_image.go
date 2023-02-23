package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type EC2I struct {
	ResourceViewer
}

func NewEC2I(resource string) ResourceViewer {
	var es EC2I
	es.ResourceViewer = NewBrowser(resource)
	es.AddBindKeysFn(es.bindKeys)
	return &es
}

func (ei *EC2I) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Image-Id", ei.GetTable().SortColCmd("Image-Id", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ei.App().PrevCmd, true),
		tcell.KeyEnter:  ui.NewKeyAction("View", ei.enterCmd, true),
	})
}

func (ei *EC2I) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	imageId := ei.GetTable().GetSelectedItem()
	if imageId != "" {
		f := describeResource
		if ei.GetTable().enterFn != nil {
			f = ei.GetTable().enterFn
		}
		f(ei.App(), ei.GetTable().GetModel(), ei.Resource(), imageId)
		ei.App().Flash().Info("Image-Id: " + imageId)
	}
	return nil
}
