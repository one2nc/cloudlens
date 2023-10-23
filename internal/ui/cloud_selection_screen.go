package ui

import (
	"fmt"
	"strings"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

const selectCloudSlantBanner = `
   ____      __           __         __               __
  / __/___  / /___  ____ / /_  ____ / /___  __ __ ___/ /
 _\ \ / -_)/ // -_)/ __// __/ / __// // _ \/ // // _  / 
/___/ \__//_/ \__/ \__/ \__/  \__//_/ \___/\_,_/ \_,_/  
													  
`

type OptionWithAction map[string]func()

type CloudSelectionScreen struct {
	*tview.Flex
	focusItem *tview.List
}

func NewCloudSelectionScreen(optionsWithAction OptionWithAction, version string) *CloudSelectionScreen {
	cloudSelectScreen := &CloudSelectionScreen{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
	}

	// Add logo
	logoString := fmt.Sprintf("%s\n%s%s", strings.Join(LogoBig, "\n"), strings.Repeat("  ", 23), version)
	logo := tview.NewTextView().SetText(logoString).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorWheat)
	logo.SetBorderPadding(2, 0, 0, 0)
	cloudSelectScreen.AddItem(logo, 0, 2, false)

	// Add selection title
	selectionTitle := tview.NewTextView().SetText(selectCloudSlantBanner).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreen)
	cloudSelectScreen.AddItem(selectionTitle, 0, 1, false)

	cscFlexCol := tview.NewFlex().SetDirection(tview.FlexColumn)
	lf1 := tview.NewBox()
	cscFlexCol.AddItem(lf1, 0, 1, false)

	cs := tview.NewList()
	cloudSelectScreen.focusItem = cs
	cscFlexCol.AddItem(cs, 0, 2, true)

	lf3 := tview.NewBox()
	cscFlexCol.AddItem(lf3, 0, 1, false)

	isCsFirstDraw := true
	cs.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		if isCsFirstDraw {
			_, _, aw, _ := lf1.GetRect()
			_, _, bw, _ := lf3.GetRect()
			x = (aw + width + bw) / 2
			x = x - 11
			//tview.Print(screen, fmt.Sprintf("ax: %v aw:%v x:%v xw:%v bx:%v bw:%v", ax, aw, x, width, bx, bw), x, y, width-2, tview.AlignCenter, tcell.ColorYellow)
		}
		return x + 2, y + 2, width, height
	})
	// Add selection options
	for option, action := range optionsWithAction {
		cs.AddItem(strings.Repeat("  ", 4)+option+strings.Repeat("  ", 4), "", 0, action)
	}

	// design selection screen part
	cs.SetBorder(true)
	cs.SetBorderColor(cs.GetBackgroundColor())
	cs.ShowSecondaryText(false)
	cs.SetHighlightFullLine(false)
	cs.SetSelectedFocusOnly(true)
	cs.SetWrapAround(true)
	cs.Blur()

	cloudSelectScreen.AddItem(cscFlexCol, 0, 1, true)
	cloudSelectScreen.AddItem(tview.NewBox(), 0, 1, false)
	cloudSelectScreen.AddItem(tview.NewBox(), 0, 1, false)
	cloudSelectScreen.AddItem(tview.NewBox(), 0, 1, false)
	return cloudSelectScreen
}

func (c *CloudSelectionScreen) GetFocusItem() *tview.List {
	return c.focusItem
}
