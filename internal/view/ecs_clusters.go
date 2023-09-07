package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type EcsClusters struct {
	ResourceViewer
}

func NewEcs(resource string) ResourceViewer {
	var ecs EcsClusters
	ecs.ResourceViewer = NewBrowser(resource)
	ecs.AddBindKeysFn(ecs.bindKeys)
	return &ecs
}
func (ecs *EcsClusters) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftB:    ui.NewKeyAction("Sort Cluster-Arn", ecs.GetTable().SortColCmd("Cluster-Arn", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ecs.App().PrevCmd, false),
	})
}
