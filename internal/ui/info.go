package ui

import (
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type Info struct {
	*tview.Flex
	dropdowns []DropDown
	items     map[string]tview.Primitive
}

func NewInfo(items map[string]tview.Primitive) *Info {
	i := Info{
		Flex:      tview.NewFlex(),
		items:     items,
		dropdowns: []DropDown{},
	}
	i.padDropDownLabels()
	i.build()
	return &i
}

func (i *Info) padDropDownLabels() {
	for _, p := range i.items {
		d, ok := p.(DropDown)
		if ok {
			i.dropdowns = append(i.dropdowns, d)
		}
	}
}

func (i *Info) build() {
	i.Clear()
	i.SetDirection(tview.FlexRow)
	i.SetBorderColor(tcell.ColorBlack.TrueColor())
	i.SetBorderPadding(1, 1, 1, 1)
	for _, p := range i.items {
		i.AddItem(p, 0, 1, false)
	}
	// for k, v := range i.data {
	// 	fmt.Fprintf(i, "[%s::b]%s: [%s::b]%s\n", "orange", k, "#ffffff", v)
	// }
}
