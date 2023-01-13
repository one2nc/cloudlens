package view

import (
	"github.com/derailed/tview"
	"github.com/one2nc/cloud-lens/internal/config"

)


type Menu struct {
	profileDropdown *tview.DropDown
	regionDropdown *tview.DropDown
}

func NewMenu(cfg config.Config) *Menu {
	return nil
}

func (m *Menu) GetCurrentRegion(){

}