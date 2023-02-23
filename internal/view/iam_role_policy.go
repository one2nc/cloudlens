package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type IamRolePloicy struct {
	ResourceViewer
}

func NewIamRolePloicy(resource string) ResourceViewer {
	var irp IamRolePloicy
	irp.ResourceViewer = NewBrowser(resource)
	irp.AddBindKeysFn(irp.bindKeys)
	return &irp
}

func (up *IamRolePloicy) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		tcell.KeyEscape: ui.NewKeyAction("Back", up.App().PrevCmd, true),
		ui.KeyShiftA:    ui.NewKeyAction("Policy-ARN", up.GetTable().SortColCmd("Policy-ARN", true), true),
		ui.KeyShiftN:    ui.NewKeyAction("Policy-Name", up.GetTable().SortColCmd("Policy-Name", true), true),
	})
}
