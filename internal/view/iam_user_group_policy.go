package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type iamUserGroupPloicy struct {
	ResourceViewer
}

func NewIamUserGroupPloicy(resource string) ResourceViewer {
	var ugp iamUserGroupPloicy
	ugp.ResourceViewer = NewBrowser(resource)
	ugp.AddBindKeysFn(ugp.bindKeys)
	return &ugp
}

func (ugp *iamUserGroupPloicy) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		tcell.KeyEscape: ui.NewKeyAction("Back", ugp.App().PrevCmd, true),
		ui.KeyShiftA:    ui.NewKeyAction("Policy-ARN", ugp.GetTable().SortColCmd("Policy-ARN", true), true),
		ui.KeyShiftN:    ui.NewKeyAction("Policy-Name", ugp.GetTable().SortColCmd("Policy-Name", true), true),
	})
}
