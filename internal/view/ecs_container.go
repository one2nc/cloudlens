package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/ui"
)

type EcsContainer struct {
	name string
	ResourceViewer
}

func NewEcsContainer(taskId string) ResourceViewer {
	var ecs EcsContainer
	ecs.name = taskId
	ecs.ResourceViewer = NewBrowser(internal.LowercaseEcsContainer)
	ecs.AddBindKeysFn(ecs.bindKeys)
	return &ecs
}

func (ecs *EcsContainer) Name() string {
	return ecs.name
}

func (ecs *EcsContainer) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyD:         ui.NewKeyAction("Describe", ecs.describeEcsContainer, true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ecs.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", ecs.enterCmd, false),
	})
}

func (ecs *EcsContainer) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	return ecs.describeEcsContainer(evt)
}

func (ecs *EcsContainer) describeEcsContainer(evt *tcell.EventKey) *tcell.EventKey {
	containerId := ecs.GetTable().GetSelectedCell(1)
	if containerId == "" {
		return nil
	}
	f := describeResource
	if ecs.GetTable().enterFn != nil {
		f = ecs.GetTable().enterFn
	}
	f(ecs.App(), ecs.GetTable().GetModel(), ecs.Resource(), containerId)
	ecs.App().Flash().Infof("Container %s", containerId)
	return nil
}
