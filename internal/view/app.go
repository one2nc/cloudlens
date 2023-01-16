package view

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/aws"
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
	cfg, _ := config.Get()

	main := tview.NewFlex().SetDirection(tview.FlexRow)

	//------menu-----
	menuColFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	//menuColFlex.SetBorder(true)
	ddRowFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	//ddRowFlex.SetBorder(true)

	textv := tview.NewTextView()
	textv.SetBackgroundColor(tcell.ColorBlack)
	//textv.SetBorder(true)

	var ins []aws.EC2Resp
	var buckets []aws.BucketResp

	var currentRegion *string
	var currentProfile *string
	profiles := cfg.Profiles
	regions := aws.GetAllRegions()

	currentProfile = &profiles[0]
	currentRegion = &regions[0]
	sess, _ := config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)

	profileDropdown := tview.NewDropDown().
		SetLabel("Profile  ").
		SetOptions(profiles, func(text string, index int) {
			currentProfile = &text
			sess, _ = config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)
			ins, _ = aws.NewEc2Service(*sess).GetInstances()
			buckets, _ = aws.NewS3Service(*sess).ListBuckets()
			textv.SetText("🌈🌧️ cloudlens profile..." + fmt.Sprintf("%v", a.Main.Pages.CurrentPage()))
		})
	profileDropdown.SetBorderFocusColor(tcell.ColorSpringGreen)
	profileDropdown.SetCurrentOption(0)
	profileDropdown.SetBorderPadding(2, 0, 0, 0)
	//profileDropdown.SetBorder(true)

	regionDropdown := tview.NewDropDown().
		SetLabel("Region   ").
		SetOptions(regions, func(text string, index int) {
			currentRegion = &text
			sess, _ = config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)
			ins, _ = aws.NewEc2Service(*sess).GetInstances()
			buckets, _ = aws.NewS3Service(*sess).ListBuckets()
			textv.SetText("🌈🌧️ cloudlens region..." + fmt.Sprintf("%v", ins))
		})
	regionDropdown.SetBorderFocusColor(tcell.ColorSpringGreen)
	regionDropdown.SetCurrentOption(0)
	regionDropdown.SetBorderPadding(0, 0, 0, 0)
	//regionDropdown.SetBorder(true)

	ddRowFlex.AddItem(profileDropdown, 0, 1, false)
	ddRowFlex.AddItem(regionDropdown, 0, 1, false)

	menuColFlex.AddItem(tview.NewBox(), 0, 1, false)
	menuColFlex.AddItem(tview.NewBox(), 0, 1, false)
	menuColFlex.AddItem(tview.NewBox(), 0, 1, false)
	menuColFlex.AddItem(ddRowFlex, 0, 1, false)

	servicePage := tview.NewFlex().SetDirection(tview.FlexRow)
	servicePageContent := DisplayEc2Instances(ins)
	servicePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		servicePage.RemoveItem(servicePage.RemoveItemAtIndex(0))
		servicePage.AddItemAtIndex(0, DisplayEc2Instances(ins), 0, 6, true)
		return event
	})
	servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
	servicePage.AddItem(servicePageContent, 0, 6, true)

	inputField := tview.NewInputField().
		SetLabel("🐶>").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return true // accept any input
		})
	inputField.SetFieldBackgroundColor(tcell.ColorBlack)
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			serviceName := inputField.GetText()
			switch serviceName {
			case "S3", "s3":
				servicePage.RemoveItemAtIndex(0)
				servicePageContent = DisplayS3Buckets(buckets)
				servicePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					servicePage.RemoveItem(servicePage.RemoveItemAtIndex(0))
					servicePage.AddItemAtIndex(0, DisplayS3Buckets(buckets), 0, 6, true)
					return event
				})
				servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
				servicePage.AddItem(servicePageContent, 0, 6, true)
				inputField.SetText("")

			case "EC2", "ec2", "Ec2", "eC2":
				servicePage.RemoveItemAtIndex(0)
				servicePageContent = DisplayEc2Instances(ins)
				servicePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					servicePage.RemoveItem(servicePage.RemoveItemAtIndex(0))
					servicePage.AddItemAtIndex(0, DisplayEc2Instances(ins), 0, 6, true)
					return event
				})
				servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
				// ec2Page.AddItem(menuColFlex, 0, 2, false)
				servicePage.AddItem(servicePageContent, 0, 6, true)
				inputField.SetText("")

			default:
				textv.SetText("🌈🌧️ No service...")
				inputField.SetText("")
			}
		}
	})

	inputField.SetBorder(true)

	main.AddItem(menuColFlex, 0, 2, true)
	main.AddItem(inputField, 0, 1, false)
	// main.AddItem(content, 0, 6, true)
	main.AddItem(servicePage, 0, 6, false)
	main.AddItem(textv, 0, 2, false)
	a.Main.AddPage("main", main, true, false)
	a.Main.ShowPage("main")
}

func DisplayEc2Instances(ins []aws.EC2Resp) *tview.Flex {
	flex := tview.NewFlex()
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetBorderFocusColor(tcell.ColorSpringGreen)
	//table data
	table.SetCell(0, 0, tview.NewTableCell("Instance-Id").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Instance-State").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Instance-Type").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Availability-zone").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("Public-DNS").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 5, tview.NewTableCell("Public-IPV4").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 6, tview.NewTableCell("Monitoring-State").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 7, tview.NewTableCell("Launch-Time").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))

	for i, in := range ins {
		table.SetCell((i + 1), 0, tview.NewTableCell(in.InstanceId).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 1, tview.NewTableCell(in.InstanceState).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 2, tview.NewTableCell(in.InstanceType).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 3, tview.NewTableCell(in.AvailabilityZone).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 4, tview.NewTableCell(in.PublicDNS).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 5, tview.NewTableCell(in.PublicIPv4).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 6, tview.NewTableCell(in.MonitoringState).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 7, tview.NewTableCell(in.LaunchTime).SetAlign(tview.AlignCenter))
	}
	table.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.SetSelectable(true, false)
	})

	flex.AddItem(table, 0, 1, true).SetDirection(tview.FlexRow)
	return flex
}

func DisplayS3Buckets(buckets []aws.BucketResp) *tview.Flex {
	flex := tview.NewFlex()

	table := tview.NewTable()
	table.SetBorderFocusColor(tcell.ColorSpringGreen)
	table.SetBorder(true)

	//layout
	flex.AddItem(table, 0, 1, true).SetDirection(tview.FlexRow)

	//table data
	table.SetCell(0, 0, tview.NewTableCell("Bucket-Name").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Creation-Time").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))

	for i, b := range buckets {
		table.SetCell((i + 1), 0, tview.NewTableCell(b.BucketName).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 1, tview.NewTableCell(fmt.Sprintf("%v", b.CreationTime)).SetAlign(tview.AlignCenter))
		table.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				table.SetSelectable(true, false)
			}
		}).SetSelectedFunc(func(row int, column int) {
			table.SetSelectable(true, false)
		})
	}
	return flex
}
