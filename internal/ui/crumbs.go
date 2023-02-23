package ui

import (
	"fmt"
	"strings"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/model"
)

// Crumbs represents user breadcrumbs.
type Crumbs struct {
	*tview.TextView

	stack *model.Stack
}

// NewCrumbs returns a new breadcrumb view.
func NewCrumbs() *Crumbs {
	c := Crumbs{
		stack:    model.NewStack(),
		TextView: tview.NewTextView(),
	}
	c.SetTextAlign(tview.AlignLeft)
	c.SetBorderPadding(0, 0, 1, 1)
	c.SetDynamicColors(true)

	return &c
}

// StackPushed indicates a new item was added.
func (c *Crumbs) StackPushed(comp model.Component) {
	c.stack.Push(comp)
	c.refresh(c.stack.Flatten())
}

// StackPopped indicates an item was deleted.
func (c *Crumbs) StackPopped(_, _ model.Component) {
	c.stack.Pop()
	c.refresh(c.stack.Flatten())
}

// StackTop indicates the top of the stack.
func (c *Crumbs) StackTop(top model.Component) {}

// Refresh updates view with new crumbs.
func (c *Crumbs) refresh(crumbs []string) {
	c.Clear()
	last, bgColor := len(crumbs)-1, "#FFE4E1"
	for i, crumb := range crumbs {
		if i == last {
			bgColor = "orange"
		}
		fmt.Fprintf(c, "[%s:%s:b] <%s> [-:-:-] ",
			"#000437",
			bgColor, strings.Replace(strings.ToLower(crumb), " ", "", -1))
	}
}
