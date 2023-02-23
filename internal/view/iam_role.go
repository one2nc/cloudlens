package view

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/ui"
)

type IamRole struct {
	ResourceViewer
}

// NewSG returns a new viewer.
func NewIamRole(resource string) ResourceViewer {
	var iamu IamRole
	iamu.ResourceViewer = NewBrowser(resource)
	iamu.AddBindKeysFn(iamu.bindKeys)
	return &iamu
}

func (ir IamRole) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Role-Id ", ir.GetTable().SortColCmd("Role-Id", true), true),
		ui.KeyShiftN:    ui.NewKeyAction("Sort Role-Name", ir.GetTable().SortColCmd("Role-Name", true), true),
		ui.KeyShiftD:    ui.NewKeyAction("Sort Created-Date", ir.GetTable().SortColCmd("Created-Date", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ir.App().PrevCmd, true),
		ui.KeyShiftP:    ui.NewKeyAction("View", ir.enterCmd, true),
	})
}

func (ir *IamRole) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	roleName := ir.GetTable().GetSecondColumn()
	if roleName != "" {
		irp := NewIamRolePloicy("Role Policy")
		ctx := context.WithValue(ir.App().GetContext(), internal.RoleName, roleName)
		ir.App().SetContext(ctx)
		ir.App().Flash().Info("Role Name: " + roleName)
		ir.App().inject(irp)
	}
	return nil
}
