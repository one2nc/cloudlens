package view

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/one2nc/cloud-lens/internal/ui"
	"github.com/rs/zerolog/log"
)

type EC2 struct {
	ResourceViewer
}

// NewPod returns a new viewer.
func NewEC2(resource string) ResourceViewer {
	cfg, _ := config.Get()
	session, _ := config.GetSession(cfg.Profiles[0], "ap-south-1", cfg.AwsConfig)
	ctx := context.WithValue(context.Background(), internal.KeySession, session)

	var e EC2
	e.ResourceViewer = NewBrowser(resource, ctx)
	e.AddBindKeysFn(e.bindKeys)
	e.GetTable().SetEnterFn(e.describeInstace)
	return &e
}

func (e *EC2) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		//ui.KeyShiftT: ui.NewKeyAction("Sort Restart", e.GetTable().SortColCmd("RESTARTS", false), false),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Type", nil, false),
		ui.KeyShiftL:    ui.NewKeyAction("Sort Launch-Time", nil, false),
		tcell.KeyEscape: ui.NewKeyAction("Back", e.App().PrevCmd, true),
	})
}

func (e *EC2) describeInstace(app *App, model ui.Tabular, resource string) {
	log.Info().Msg(fmt.Sprintf("TODO: describe: %v", resource))
	// if err := app.inject(co); err != nil {
	// 	app.Flash().Err(err)
	// }
}
