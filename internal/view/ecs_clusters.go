package view

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
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
		ui.KeyD:         ui.NewKeyAction("Describe", ecs.describeCluster, true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ecs.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", ecs.enterCmd, false),
	})
}

func (ecs *EcsClusters) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	clusterName := ecs.GetTable().GetSelectedItem()
	if clusterName == "" {
		return nil
	}
	ecsServiceScreen := NewEcsService(clusterName)
	ctx := context.WithValue(ecs.App().GetContext(), internal.ECSClusterName, clusterName)
	ecs.App().SetContext(ctx)
	ecs.App().inject(ecsServiceScreen)
	ecsServiceScreen.GetTable().SetTitle(fmt.Sprintf(" ecs://%s ", clusterName))
	ecsServiceScreen.App().Flash().Info(fmt.Sprintf("Viewing %s cluster...", clusterName))
	return nil
}

func (ecs *EcsClusters) describeCluster(evt *tcell.EventKey) *tcell.EventKey {
	clusterName := ecs.GetTable().GetSelectedItem()
	if clusterName == "" {
		return nil
	}
	f := describeResource
	if ecs.GetTable().enterFn != nil {
		f = ecs.GetTable().enterFn
	}
	f(ecs.App(), ecs.GetTable().GetModel(), ecs.Resource(), clusterName)
	return nil
}
