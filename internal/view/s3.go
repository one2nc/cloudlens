package view

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/ui"
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
		tcell.KeyEscape: ui.NewKeyAction("Back", s3.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", s3.enterCmd, false),
		ui.KeyD:         ui.NewKeyAction("Describe", s3.describeBucket, true),
	})
}

func (s3 *S3) describeBucket(evt *tcell.EventKey) *tcell.EventKey {
	bName := s3.GetTable().GetSelectedItem()
	f := describeResource
	if s3.GetTable().enterFn != nil {
		f = s3.GetTable().enterFn
	}
	if bName != "" {
		s3.App().Flash().Info("Bucket-Name: " + bName)
		f(s3.App(), s3.GetTable().GetModel(), s3.Resource(), bName)
	}

	return nil
}

func (s3 *S3) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	bName := s3.GetTable().GetSelectedItem()
	if bName != "" {
		o := NewS3FileViewer("s3://", bName)
		ctx := context.WithValue(s3.App().GetContext(), internal.BucketName, bName)
		s3.App().SetContext(ctx)
		ctx = context.WithValue(s3.App().GetContext(), internal.FolderName, "")
		s3.App().SetContext(ctx)
		s3.App().Flash().Info("Bucket Name: " + bName)
		s3.App().inject(o)
		o.GetTable().SetTitle(o.path)
	}

	return nil
}
