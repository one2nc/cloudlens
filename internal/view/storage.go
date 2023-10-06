package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type Storage struct {
	ResourceViewer
}

func NewStorage(resource string) ResourceViewer {
	var l Storage
	l.ResourceViewer = NewBrowser(resource)
	l.AddBindKeysFn(l.bindKeys)
	return &l
}

func (s *Storage) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftB:    ui.NewKeyAction("Sort Bucket-Name", s.GetTable().SortColCmd("Bucket-Name", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Creation-Time", s.GetTable().SortColCmd("Creation-Time", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", s.App().PrevCmd, false),
	})
}
