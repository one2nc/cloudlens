package view

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/ui"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	ResourceViewer
}

func NewStorage(resource string) ResourceViewer {
	var s Storage
	s.ResourceViewer = NewBrowser(resource)
	s.AddBindKeysFn(s.bindKeys)
	return &s
}

func (s *Storage) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftB:    ui.NewKeyAction("Sort Bucket-Name", s.GetTable().SortColCmd("Bucket-Name", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Creation-Time", s.GetTable().SortColCmd("Creation-Time", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", s.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", s.enterCmd, false),
	})
}

func (s *Storage) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	bName := s.GetTable().GetSelectedItem()
	log.Print(bName)
	if bName != "" {
		o := NewStorageFileViewer()
		ctx := context.WithValue(s.App().GetContext(), internal.BucketName, bName)
		s.App().SetContext(ctx)
		ctx = context.WithValue(s.App().GetContext(), internal.FolderName, "")
		s.App().SetContext(ctx)
		s.App().Flash().Info("Bucket Name: " + bName)
		s.App().inject(o)
		o.GetTable().SetTitle(bName)
	}

	return nil
}
