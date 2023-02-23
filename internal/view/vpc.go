package view

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/ui"
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
		tcell.KeyEnter:  ui.NewKeyAction("View", v.enterCmd, true),
		ui.KeyS:         ui.NewKeyAction("View Subnets", v.subnetCmd, true),
	})
}

func (v *VPC) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	vpcId := v.GetTable().GetSelectedItem()
	if vpcId != "" {
		f := describeResource
		if v.GetTable().enterFn != nil {
			f = v.GetTable().enterFn
		}
		f(v.App(), v.GetTable().GetModel(), v.Resource(), vpcId)
		v.App().Flash().Info("VPC Id: " + vpcId)
	}

	return nil
}

func (iamug *VPC) subnetCmd(evt *tcell.EventKey) *tcell.EventKey {
	vpcId := iamug.GetTable().GetSelectedItem()
	if vpcId != "" {
		sn := NewSubnet("subnet")
		ctx := context.WithValue(iamug.App().GetContext(), internal.VpcId, vpcId)
		iamug.App().SetContext(ctx)
		iamug.App().Flash().Info("VPC ID: " + vpcId)
		iamug.App().inject(sn)
	}
	return nil
}
