package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type Subnet struct {
	ResourceViewer
}

// NewPod returns a new viewer.
func NewSubnet(resource string) ResourceViewer {
	var sn Subnet
	sn.ResourceViewer = NewBrowser(resource)
	sn.AddBindKeysFn(sn.bindKeys)
	return &sn
}

func (sn *Subnet) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Subnet-Id", sn.GetTable().SortColCmd("Subnet-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Subnet-State", sn.GetTable().SortColCmd("Subnet-State", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", sn.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", sn.enterCmd, true),
	})
}

func (sn *Subnet) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	subnetId := sn.GetTable().GetSelectedItem()
	if subnetId != "" {
		f := describeResource
		if sn.GetTable().enterFn != nil {
			f = sn.GetTable().enterFn
		}
		f(sn.App(), sn.GetTable().GetModel(), sn.Resource(), subnetId)
		sn.App().Flash().Info("Subnet Id: " + subnetId)
	}

	return nil
}
