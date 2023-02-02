package view

import (
	"context"
	"fmt"

	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/one2nc/cloud-lens/internal/ui"
	"github.com/rs/zerolog/log"
)

type S3 struct {
	ResourceViewer
}

func NewS3(resource string) ResourceViewer {
	cfg, _ := config.Get()
	session, _ := config.GetSession(cfg.Profiles[0], "ap-east-1", cfg.AwsConfig)
	ctx := context.WithValue(context.Background(), internal.KeySession, session)

	var s3 S3
	s3.ResourceViewer = NewBrowser(resource, ctx)
	s3.AddBindKeysFn(s3.bindKeys)
	//s3.GetTable().SetEnterFn(s3.describeInstace)
	return &s3
}
func (s3 *S3) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftB: ui.NewKeyAction("Sort Bucket-Name", s3.GetTable().SortColCmd("Bucket-Name", true), true),
		ui.KeyShiftT: ui.NewKeyAction("Sort Creation-Time", s3.GetTable().SortColCmd("Creation-Time", true), true),
	})
}

func (s3 *S3) describeInstace(app *App, model ui.Tabular, resource string) {
	log.Info().Msg(fmt.Sprintf("TODO: describe: %v", resource))
	// if err := app.inject(co); err != nil {
	// 	app.Flash().Err(err)
	// }
}
