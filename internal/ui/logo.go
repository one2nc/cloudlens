package ui

import (
	"strings"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

// LogoSmall cls small log.
var LogoSmall = []string{
	`_________ .__       `,
	`\_   ___ \|  |   ______`,
	`/    \  \/|  |  /  ___/`,
	`\     \___|  |__\___ \ `,
	` \______  /____/____  >`,
	`        \/          \/ `,
}

// LogoBig cls big logo for splash page.
var LogoBig = []string{
	`█▀▀ █   █▀▀█ █  █ █▀▀▄ █   █▀▀ █▀▀▄ █▀▀`,
	`█   █   █  █ █  █ █  █ █   █▀▀ █  █ ▀▀█`,
	`▀▀▀ ▀▀▀ ▀▀▀▀  ▀▀▀ ▀▀▀  ▀▀▀ ▀▀▀ ▀  ▀ ▀▀▀`,
}

type Logo struct {
	*tview.Flex
	logo *tview.TextView
}

func NewLogo() *Logo {
	l := Logo{
		Flex: tview.NewFlex(),
		logo: tview.NewTextView(),
	}
	l.SetDirection(tview.FlexRow)
	l.buildLogo()
	l.AddItem(l.logo, 6, 1, false)
	return &l
}

func (l *Logo) buildLogo() {
	l.logo.SetText(strings.Join(LogoSmall, "\n"))
	l.logo.SetTextColor(tcell.ColorOrange)
}
