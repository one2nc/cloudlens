package ui

import (
	"strings"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

// const selectCloudCalvinSBanner = `

// ╔═╗┬  ┌─┐┬ ┬┌┬┐  ┌─┐
// ║  │  │ ││ │ ││   ┌┘
// ╚═╝┴─┘└─┘└─┘─┴┘   o

// `

const selectCloudANSIShadowBanner = `

██████╗██╗      ██████╗ ██╗   ██╗██████╗     ██████╗ 
██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗    ╚════██╗
██║     ██║     ██║   ██║██║   ██║██║  ██║      ▄███╔╝
██║     ██║     ██║   ██║██║   ██║██║  ██║      ▀▀══╝ 
╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝      ██╗   
 ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝       ╚═╝   
                                                       
`

type OptionWithAction map[string]func()

type CloudSelectionScreen struct {
	*tview.Flex
	focusItem *tview.List
}

func NewCloudSelectionScreen(optionsWithAction OptionWithAction) *CloudSelectionScreen {
	cloudSelectScreen := &CloudSelectionScreen{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
	}

	// Add selection title
	selectionTitle := tview.NewTextView().SetText(selectCloudANSIShadowBanner).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreen)
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

	for option, action := range optionsWithAction {
		cs.AddItem(strings.Repeat("  ", 4)+option+strings.Repeat("  ", 4), "", 0, action)
	}

	// design selection screen part
	cs.SetTitleColor(tcell.ColorGreen)
	cs.SetTitleColor(tcell.ColorGreen)
	cs.SetTitleAlign(tview.AlignCenter)
	cs.SetBorder(true)
	cs.SetBorderColor(cs.GetBackgroundColor())
	cs.ShowSecondaryText(false)
	cs.SetHighlightFullLine(false)
	cs.SetSelectedFocusOnly(true)
	cs.SetWrapAround(true)
	cs.Blur()

	cloudSelectScreen.AddItem(cscFlexCol, 0, 1, true)
	cloudSelectScreen.AddItem(tview.NewBox(), 0, 1, false)
	return cloudSelectScreen
}

func (c *CloudSelectionScreen) GetFocusItem() *tview.List {
	return c.focusItem
}
