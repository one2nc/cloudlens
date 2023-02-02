package ui

import (
	"fmt"

	"github.com/derailed/tview"
)

// Crumbs represents user breadcrumbs.
type Info struct {
	*tview.TextView
	data map[string]string
}

// NewCrumbs returns a new breadcrumb view.
func NewInfo(data map[string]string) *Info {
	i := Info{
		TextView: tview.NewTextView(),
		data:     data,
	}
	i.SetTextAlign(tview.AlignLeft)
	i.SetBorderPadding(0, 0, 1, 1)
	i.SetDynamicColors(true)
	i.build()
	return &i
}

// Refresh updates view with new crumbs.
func (i *Info) build() {
	i.Clear()
	for k, v := range i.data {
		fmt.Fprintf(i, "[%s::b]%s: [%s::b]%s\n", "orange", k, "#ffffff", v)
	}
}
