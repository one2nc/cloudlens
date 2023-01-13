package view

import (
	"fmt"
	"time"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/config"

	"github.com/one2nc/cloud-lens/internal/ec2"
	"github.com/one2nc/cloud-lens/internal/s3"
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

	var ins []ec2.EC2Resp
	var buckets []s3.BucketResp

	var currentRegion *string
	var currentProfile *string
	prfoiles := cfg.Profiles
	regions := []string{"us-east-2", "us-east-1", "us-west-1", "us-west-2", "af-south-1", "ap-east-1", "ap-south-2", "ap-southeast-3", "ap-south-1", "ap-northeast-3", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-south-2", "eu-north-1", "eu-central-2", "me-south-1", "me-central-1", "sa-east-1", "us-gov-east-1", "us-gov-west-1"}

	currentProfile = &prfoiles[0]
	currentRegion = &regions[0]
	sess, _ := config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)

	profileDropdown := tview.NewDropDown().
		SetLabel("Profile  ").
		SetOptions(prfoiles, func(text string, index int) {
			currentProfile = &text
			sess, _ = config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)
			ins, _ = ec2.NewEc2Service(*sess).GetInstances()
			buckets, _ = s3.NewS3Service(*sess).ListBuckets()
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
			ins, _ = ec2.NewEc2Service(*sess).GetInstances()
			buckets, _ = s3.NewS3Service(*sess).ListBuckets()
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

	//------body-----
	content := tview.NewFlex().SetDirection(tview.FlexColumn)
	content.SetTitle("<SERVICES>").SetTitleColor(tcell.NewRGBColor(0, 252, 255))
	content.SetBorderColor(tcell.ColorOrangeRed)
	content.SetBorderFocusColor(tcell.ColorSpringGreen)
	content.SetBorder(true)

	tableFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	tableV := tview.NewTable()
	tableV.SetBorder(true)
	tableV.SetBorderFocusColor(tcell.ColorSpringGreen)
	tableV.SetBorderPadding(1, 1, 1, 1)
	tableV.SetSelectable(true, false)

	services := []string{"S3", "EC2"}
	for i, s := range services {
		tableV.SetCell(i, 25, tview.NewTableCell(s).SetAlign(tview.AlignRight).SetSelectable(true))
	}

	ec2Page := tview.NewFlex().SetDirection(tview.FlexRow)
	s3Page := tview.NewFlex().SetDirection(tview.FlexRow)

	tableV.SetSelectedFunc(func(row, column int) {
		textv.SetText("🌈🌧️ cloudlens service..." + fmt.Sprintf("%v %v %v", currentProfile, currentRegion, row))
		switch services[row] {
		case "S3":
			a.Main.Pages.RemovePage("s3Page")
			s3PageContent := DisplayS3Buckets(buckets)
			s3Page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				s3Page.RemoveItem(s3Page.RemoveItemAtIndex(1))
				s3Page.AddItemAtIndex(1, DisplayS3Buckets(buckets), 0, 6, true)
				textv.SetText("🌈🌧️ cloudlens event..." + fmt.Sprintf("%T", s3PageContent))
				if event.Rune() == 98 {
					a.Main.Pages.SwitchToPage("main")
				}
				return event
			})
			s3Page.AddItem(menuColFlex, 0, 2, false)
			s3Page.AddItem(s3PageContent, 0, 6, true)
			s3Page.AddItem(textv, 0, 2, false)
			a.Main.Pages.AddAndSwitchToPage("s3Page", s3Page, true)

		case "EC2":
			a.Main.Pages.RemovePage("ec2Page")
		
			ec2PageContent := DisplayEc2Instances(ins)
			ec2Page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEnter {
					textv.SetText("🌈🌧️ cloudlens event..." + time.Now().String())	
				}
				ec2Page.RemoveItem(ec2Page.RemoveItemAtIndex(1))
				ec2Page.AddItemAtIndex(1, DisplayEc2Instances(ins), 0, 6, true)
				if event.Rune() == 98 {
					a.Main.Pages.SwitchToPage("main")
				}
				return event
			})
			ec2PageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
			ec2Page.AddItem(menuColFlex, 0, 2, false)
			ec2Page.AddItem(ec2PageContent, 0, 6, true)
			ec2Page.AddItem(textv, 0, 2, false)
			a.Main.Pages.AddAndSwitchToPage("ec2Page", ec2Page, true)

		default:
			textv.SetText("🌈🌧️ No service...")
		}
	})

	tableFlex.AddItem(tableV, 0, 1, true)

	content.AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
	content.AddItem(tableFlex, 0, 1, true)
	content.AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	main.AddItem(menuColFlex, 0, 2, false)
	main.AddItem(content, 0, 6, true)
	main.AddItem(textv, 0, 2, false)
	a.Main.AddPage("main", main, true, false)
	a.Main.ShowPage("main")
}

func DisplayEc2Instances(ins []ec2.EC2Resp) *tview.Flex {
	flex := tview.NewFlex()
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetBorderFocusColor(tcell.ColorSpringGreen)
	//table data
	table.SetCell(0, 0, tview.NewTableCell("Instance-Id").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Instance-State").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Instance-Type").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Availability-zone").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("Monitoring-State").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 5, tview.NewTableCell("Launch-Time").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))

	for i, in := range ins {
		table.SetCell((i + 1), 0, tview.NewTableCell(in.InstanceId).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 1, tview.NewTableCell(in.InstanceState).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 2, tview.NewTableCell(in.InstanceType).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 3, tview.NewTableCell(in.AvailabilityZone).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 4, tview.NewTableCell(in.MonitoringState).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 5, tview.NewTableCell(in.LaunchTime).SetAlign(tview.AlignCenter))
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

func DisplayS3Buckets(buckets []s3.BucketResp) *tview.Flex {
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
