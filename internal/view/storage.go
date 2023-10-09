package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
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
	})
}
