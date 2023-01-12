package view

import (
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/config"

	"github.com/one2nc/cloud-lens/internal/ui"
)

type App struct {
	*ui.App
}

func NewApp() App {
	app := App{App: ui.NewApp()}
	return app
}

func (a *App) Init() {
	a.App.Init()
	a.layout()
	a.App.Run()
}

func (a *App) layout() {
	main := tview.NewFlex().SetDirection(tview.FlexRow)
	
	textv := tview.NewTextView()
	textv.SetBackgroundColor(tcell.ColorBlack)
	textv.SetBorder(true)

	content := tview.NewFlex().SetDirection(tview.FlexRow)
	content.SetBorder(true)

	textView := tview.NewTextView().SetText("Services")
	textView.SetTextAlign(tview.AlignCenter)
	content.AddItem(textView, 1, 1, false)

	cfg, _ := config.Get()
	services := []string{"S3", "EC2"}

	for _, s := range services {
		textView := tview.NewTextView().SetText(s)
		textView.SetTextAlign(tview.AlignCenter)
		content.AddItem(textView, 1, 1, false)
	}

	ddflex := tview.NewFlex().SetDirection(tview.FlexRow)
	ddflex.SetBorder(true)
	ddflex.SetBackgroundColor(tcell.ColorDarkRed)
	ddflex.SetTitleAlign(tview.AlignRight)

	profileDropdown := tview.NewDropDown().
		SetLabel("Profile:").
		SetOptions(cfg.Profiles, func(text string, index int) {
			textv.SetText("🌈🌧️ cloudlens starting up..." + text)
		})
	profileDropdown.SetCurrentOption(0)
	profileDropdown.SetBorder(true)

	regions := []string{"us-east-2", "us-east-1", "us-west-1", "us-west-2", "af-south-1", "ap-east-1", "ap-south-2", "ap-southeast-3", "ap-south-1", "ap-northeast-3", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-south-2", "eu-north-1", "eu-central-2", "me-south-1", "me-central-1", "sa-east-1", "us-gov-east-1", "us-gov-west-1"}
	regionDropdown := tview.NewDropDown().
		SetLabel("Region:").
		SetOptions(regions, func(text string, index int) {
			textv.SetText("🌈🌧️ cloudlens starting up..." + text)
		})
	regionDropdown.SetCurrentOption(0)
	regionDropdown.SetBorder(true)

	ddflex.AddItem(profileDropdown, 5, 2, false)
	ddflex.AddItem(regionDropdown, 5, 2, false)

	main.AddItem(ddflex, 0, 4, false)
	main.AddItem(content, 0, 4, false)
	main.AddItem(textv, 0, 2, true)
	a.Main.AddPage("main", main, true, false)
	a.Main.ShowPage("main")
}
