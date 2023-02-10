package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type IamGroupUser struct {
	ResourceViewer
}

func NewIamGroupUser(resource string) ResourceViewer {
	var igu IamGroupUser
	igu.ResourceViewer = NewBrowser(resource)
	igu.AddBindKeysFn(igu.bindKeys)
	return &igu
}

func (igu *IamGroupUser) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort User-Id ", igu.GetTable().SortColCmd("User-Id", true), true),
		ui.KeyShiftN:    ui.NewKeyAction("Sort User-Name", igu.GetTable().SortColCmd("User-Name", true), true),
		ui.KeyShiftD:    ui.NewKeyAction("Sort Created-Date", igu.GetTable().SortColCmd("Created-Date", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", igu.App().PrevCmd, true),
	})
}
