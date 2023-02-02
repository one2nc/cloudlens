package view

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type SG struct {
	ResourceViewer
}

// NewPod returns a new viewer.
func NewSG(resource string) ResourceViewer {
	cfg, _ := config.Get()
	session, _ := config.GetSession(cfg.Profiles[0], "ap-east-1", cfg.AwsConfig)
	ctx := context.WithValue(context.Background(), internal.KeySession, session)

	var sg SG
	sg.ResourceViewer = NewBrowser(resource, ctx)
	sg.AddBindKeysFn(sg.bindKeys)
	return &sg
}

func (sg SG) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Group-Id", sg.GetTable().SortColCmd("Group-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Group-Name", sg.GetTable().SortColCmd("Group-Name", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", sg.App().PrevCmd, true),
	})
}
