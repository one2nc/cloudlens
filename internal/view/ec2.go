package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type EC2 struct {
	ResourceViewer
}

// NewPod returns a new viewer.
func NewEC2(resource string) ResourceViewer {
	var e EC2
	e.ResourceViewer = NewBrowser(resource)
	e.AddBindKeysFn(e.bindKeys)
	// e.GetTable().SetEnterFn(e.describeInstace)
	return &e
}

func (e *EC2) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Instance-Id", e.GetTable().SortColCmd("Instance-Id", true), false),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Instance-State", e.GetTable().SortColCmd("Instance-State", true), false),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Instance-Type", e.GetTable().SortColCmd("Instance-Type", true), false),
		ui.KeyShiftL:    ui.NewKeyAction("Sort Launch-Time", e.GetTable().SortColCmd("Launch-Time", true), false),
		ui.KeyShiftM:    ui.NewKeyAction("Sort Monitoring-State", e.GetTable().SortColCmd("Monitoring-State", true), false),
		ui.KeyShiftP:    ui.NewKeyAction("Sort Public-DNS", e.GetTable().SortColCmd("Public-DNS", true), false),
		tcell.KeyEscape: ui.NewKeyAction("Back", e.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", e.enterCmd, false),
	})
}

func (e *EC2) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	instanceId := e.GetTable().GetSelectedItem()
	f := describeResource
	if e.GetTable().enterFn != nil {
		f = e.GetTable().enterFn
	}
	f(e.App(), e.GetTable().GetModel(), e.Resource(), instanceId)
	return nil
}
