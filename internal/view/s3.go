package view

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/ui"
	"github.com/rs/zerolog/log"
)

type S3 struct {
	ResourceViewer
}

func NewS3(resource string) ResourceViewer {
	var s3 S3
	s3.ResourceViewer = NewBrowser(resource)
	s3.AddBindKeysFn(s3.bindKeys)
	// s3.GetTable().SetEnterFn(s3.describeInstace)
	return &s3
}
func (s3 *S3) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftB:    ui.NewKeyAction("Sort Bucket-Name", s3.GetTable().SortColCmd("Bucket-Name", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Creation-Time", s3.GetTable().SortColCmd("Creation-Time", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", s3.App().PrevCmd, true),
		tcell.KeyEnter:  ui.NewKeyAction("View", s3.enterCmd, true),
	})
}

func (s3 *S3) describeInstace(app *App, model ui.Tabular, resource string) {
	log.Info().Msg(fmt.Sprintf("TODO: describe: %v", resource))
	// if err := app.inject(co); err != nil {
	// 	app.Flash().Err(err)
	// }
}
func (s3 *S3) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	bName := s3.GetTable().GetSelectedItem()
	o := NewS3FileViewer("OBJ")
	log.Info().Msg(fmt.Sprintf("Before Assigning Bucket Name: %v", bName))

	ctx := context.WithValue(s3.App().GetContext(), internal.BucketName, bName)
	s3.App().SetContext(ctx)
	ctx = context.WithValue(s3.App().GetContext(), internal.FolderName, "")
	s3.App().SetContext(ctx)

	s3.App().Flash().Info("After Bucket Name: " + bName)
	// println(bName)
	s3.App().inject(o)
	return nil
}
