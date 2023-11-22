package view

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/ui"
)

type EcsServices struct {
	name string
	ResourceViewer
}

func NewEcsService(clusterName string) ResourceViewer {
	var ecs EcsServices
	ecs.name = clusterName
	ecs.ResourceViewer = NewBrowser(internal.LowercaseEcsServices)
	ecs.AddBindKeysFn(ecs.bindKeys)
	return &ecs
}

func (ecs *EcsServices) Name() string {
	return ecs.name
}

func (ecs *EcsServices) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyD:         ui.NewKeyAction("Describe", ecs.describeEcsService, true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ecs.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", ecs.enterCmd, false),
	})
}

func (ecs *EcsServices) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	serviceName := ecs.GetTable().GetSelectedItem()
	if serviceName == "" {
		return nil
	}
	ecsTaskScreen := NewEcsTask(serviceName)
	ctx := ecs.App().GetContext()
	clusterName := ctx.Value(internal.ECSClusterName).(string)
	ctx = context.WithValue(ctx, internal.ECSServiceName, serviceName)
	ecs.App().SetContext(ctx)
	ecs.App().inject(ecsTaskScreen)
	ecsTaskScreen.GetTable().SetTitle(fmt.Sprintf(" ecs://%s/%s ", clusterName, serviceName))
	ecsTaskScreen.App().Flash().Info(fmt.Sprintf("Viewing %s service...", serviceName))
	return nil
}

func (ecs *EcsServices) describeEcsService(evt *tcell.EventKey) *tcell.EventKey {
	serviceName := ecs.GetTable().GetSelectedItem()
	if serviceName == "" {
		return nil
	}
	f := describeResource
	if ecs.GetTable().enterFn != nil {
		f = ecs.GetTable().enterFn
	}
	f(ecs.App(), ecs.GetTable().GetModel(), ecs.Resource(), serviceName)
	ecs.App().Flash().Infof("Service %s", serviceName)
	return nil
}
