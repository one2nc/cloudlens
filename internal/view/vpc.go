package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type VPC struct {
	ResourceViewer
}

// NewPod returns a new viewer.
func NewVPC(resource string) ResourceViewer {
	var v VPC
	v.ResourceViewer = NewBrowser(resource)
	v.AddBindKeysFn(v.bindKeys)
	return &v
}

func (v *VPC) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort VPC-Id", v.GetTable().SortColCmd("VPC-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort VPC-State", v.GetTable().SortColCmd("VPC-State", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", v.App().PrevCmd, false),
	})
}
