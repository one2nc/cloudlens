package view

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/ui"
)

type EcsTask struct {
	name string
	ResourceViewer
}

func NewEcsTask(serviceName string) ResourceViewer {
	var ecs EcsTask
	ecs.name = serviceName
	ecs.ResourceViewer = NewBrowser(internal.LowercaseEcsTasks)
	//ecs.Actions().Clear()
	ecs.AddBindKeysFn(ecs.bindKeys)
	return &ecs
}

func (ecsTask *EcsTask) Name() string {
	return ecsTask.name
}

func (ecsTask *EcsTask) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyD:         ui.NewKeyAction("Describe", ecsTask.describeEcsTask, true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ecsTask.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", ecsTask.enterCmd, true),
	})
}

func (ecsTask *EcsTask) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	taskId := ecsTask.GetTable().GetSelectedItem()
	if taskId == "" {
		return nil
	}
	ecsContainerScreen := NewEcsContainer(taskId)
	ctx := ecsTask.App().GetContext()
	clusterName := ctx.Value(internal.ECSClusterName).(string)
	serviceName := ctx.Value(internal.ECSServiceName).(string)
	ctx = context.WithValue(ctx, internal.ECSTaskId, taskId)
	ecsTask.App().SetContext(ctx)
	ecsTask.App().inject(ecsContainerScreen)
	ecsContainerScreen.GetTable().SetTitle(fmt.Sprintf(" ecs://%s/%s/%s ", clusterName, serviceName, taskId))
	ecsContainerScreen.App().Flash().Info(fmt.Sprintf("Viewing %s containers...", taskId))
	return nil
}

func (ecsTask *EcsTask) describeEcsTask(evt *tcell.EventKey) *tcell.EventKey {
	taskId := ecsTask.GetTable().GetSelectedItem()
	if taskId == "" {
		return nil
	}
	f := describeResource
	if ecsTask.GetTable().enterFn != nil {
		f = ecsTask.GetTable().enterFn
	}
	f(ecsTask.App(), ecsTask.GetTable().GetModel(), ecsTask.Resource(), taskId)
	ecsTask.App().Flash().Infof("Task %s", taskId)
	return nil
}
