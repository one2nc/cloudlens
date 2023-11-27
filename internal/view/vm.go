package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type VM struct {
	ResourceViewer
}

// NewPod returns a new viewer.
func NewVM(resource string) ResourceViewer {
	var e VM
	e.ResourceViewer = NewBrowser(resource)
	e.AddBindKeysFn(e.bindKeys)
	// e.GetTable().SetEnterFn(e.describeInstace)
	return &e
}

func (e *VM) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI: ui.NewKeyAction("Sort Instance-Id", e.GetTable().SortColCmd("Instance-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Instance-State", e.GetTable().SortColCmd("Instance-State", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Instance-Type", e.GetTable().SortColCmd("Instance-Type", true), true),
		ui.KeyShiftL:    ui.NewKeyAction("Sort Launch-Time", e.GetTable().SortColCmd("Launch-Time", true), true),
		// ui.KeyShiftM:    ui.NewKeyAction("Sort Monitoring-State", e.GetTable().SortColCmd("Monitoring-State", true), true),
		// ui.KeyShiftP:    ui.NewKeyAction("Sort Public-DNS", e.GetTable().SortColCmd("Public-DNS", true), false),
		tcell.KeyEscape: ui.NewKeyAction("Back", e.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", e.enterCmd, false),
	})
}

func (e *VM) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	instanceId := e.GetTable().GetSelectedItem()
	if instanceId != "" {
		f := describeResource
		if e.GetTable().enterFn != nil {
			f = e.GetTable().enterFn
		}
		f(e.App(), e.GetTable().GetModel(), e.Resource(), instanceId)
		e.App().Flash().Info("Instance Id: " + instanceId)
	}

	return nil
}
