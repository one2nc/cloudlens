package ui

import (
	"sort"

	"github.com/derailed/tview"
	"github.com/one2nc/cloud-lens/internal/model"
)

const (
	menuIndexFmt = " [key:-:b]<%d> [fg:-:d]%s "
	maxRows      = 6
)

// Menu presents menu options.
type Menu struct {
	*tview.Table
}

// NewMenu returns a new menu.
func NewMenu() *Menu {
	m := Menu{
		Table: tview.NewTable(),
	}
	return &m
}

func (m *Menu) HydrateMenu(hh model.MenuHints) {
	m.Clear()
	sort.Sort(hh)

	table := make([]model.MenuHints, maxRows+1)
	colCount := (len(hh) / maxRows) + 1

	for row := 0; row < maxRows; row++ {
		table[row] = make(model.MenuHints, colCount)
	}

	//TODO: change logic
	for col := 0; col < 1; col++ {
		for row := 0; row < len(hh) && row < 6; row++ {
			c := tview.NewTableCell(hh[row].Mnemonic)
			m.SetCell(row, col, c)
		}
	}
}
