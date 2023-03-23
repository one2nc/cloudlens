package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type SG struct {
	ResourceViewer
}

func NewSG(resource string) ResourceViewer {
	var sg SG
	sg.ResourceViewer = NewBrowser(resource)
	sg.AddBindKeysFn(sg.bindKeys)
	return &sg
}
func (sg *SG) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	groupId := sg.GetTable().GetSelectedItem()
	if groupId != "" {
		f := describeResource
		if sg.GetTable().enterFn != nil {
			f = sg.GetTable().enterFn
		}
		f(sg.App(), sg.GetTable().GetModel(), sg.Resource(), groupId)
		sg.App().Flash().Info("groupId: " + groupId)

	}
	return nil
}

func (sg *SG) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Group-Id", sg.GetTable().SortColCmd("Group-Id", true), true),
		ui.KeyShiftN:    ui.NewKeyAction("Sort Group-Name", sg.GetTable().SortColCmd("Group-Name", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", sg.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", sg.enterCmd, false),
	})
}
