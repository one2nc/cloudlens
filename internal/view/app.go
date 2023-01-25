package view

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/one2nc/cloud-lens/internal/config"

	"github.com/one2nc/cloud-lens/internal/ui"
	"github.com/rs/zerolog/log"
)

const (
	splashDelay = 1 * time.Second
)

type App struct {
	*ui.App
	cancelFn            context.CancelFunc
	showHeader          bool
	IsPageContentSorted bool
}

func NewApp() *App {
	a := App{App: ui.NewApp(), IsPageContentSorted: false}
	a.Views()["statusIndicator"] = ui.NewStatusIndicator(a.App)
	return &a
}

func (a *App) Init() error {
	ctx := context.WithValue(context.Background(), internal.KeyApp, a)

	a.App.Init()
	a.layout(ctx)
	return nil
}

func (a *App) layout(ctx context.Context) *tview.Flex {
	cfg, _ := config.Get()

	main := tview.NewFlex().SetDirection(tview.FlexRow)
	//----Flash Mesages----
	flash := ui.NewFlash(a.App)
	go flash.Watch(ctx, a.Flash().Channel())
	//------menu-----
	menuColFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	//menuColFlex.SetBorder(true)
	ddRowFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	//ddRowFlex.SetBorder(true)

	var ins []aws.EC2Resp
	var buckets []aws.BucketResp

	var currentRegion *string
	var currentProfile *string
	profiles := cfg.Profiles
	regions := aws.GetAllRegions()

	currentProfile = &profiles[0]
	currentRegion = &regions[0]
	sess, _ := config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)
	servicePage := tview.NewFlex().SetDirection(tview.FlexRow)
	servicePageContent := a.DisplayEc2Instances(ins, sess)

	profileDropdown := tview.NewDropDown().
		SetLabel("Profile ▼ ").
		SetOptions(profiles, func(text string, index int) {
			currentProfile = &text
			sess, _ = config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)
			ins, _ = aws.GetInstances(*sess)
			buckets, _ = aws.ListBuckets(*sess)
			if servicePage.ItemAt(0) != nil {
				servicePage.RemoveItemAtIndex(0)
			}
			servicePageContent = a.DisplayEc2Instances(ins, sess)
			servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
			servicePage.AddItem(servicePageContent, 0, 6, true)
		})
	profileDropdown.SetBorderFocusColor(tcell.ColorSpringGreen)
	profileDropdown.SetCurrentOption(0)
	profileDropdown.SetBorderPadding(2, 0, 0, 0)
	//profileDropdown.SetBorder(true)

	regionDropdown := tview.NewDropDown().
		SetLabel("Region ▼ ").
		SetOptions(regions, func(text string, index int) {
			currentRegion = &text
			sess, _ = config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)
			ins, _ = aws.GetInstances(*sess)
			buckets, _ = aws.ListBuckets(*sess)
			if servicePage.ItemAt(0) != nil {
				servicePage.RemoveItemAtIndex(0)
			}
			servicePageContent = a.DisplayEc2Instances(ins, sess)
			servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
			servicePage.AddItem(servicePageContent, 0, 6, true)
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

	servicePage = tview.NewFlex().SetDirection(tview.FlexRow)
	a.Flash().Info("Loading EC2 instacnes...")
	servicePageContent = a.DisplayEc2Instances(ins, sess)
	servicePageContent.SetBorderFocusColor(tcell.ColorSpringGreen)
	a.Application.SetFocus(servicePageContent)
	servicePageContent.SetSelectable(true, false)
	servicePageContent.Select(1, 1).SetFixed(1, 1)
	servicePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//sorting s3 Buckets
		//66 - Key B
		if event.Rune() == 66 {
			servicePage.RemoveItemAtIndex(0)
			if a.IsPageContentSorted {
				sort.Sort(sort.Reverse(aws.ByBucketName(buckets)))
				a.IsPageContentSorted = false
			} else {
				sort.Sort(aws.ByBucketName(buckets))
				a.IsPageContentSorted = true
			}
			servicePageContent = a.DisplayS3Buckets(sess, buckets)
			hc := servicePageContent.GetCell(0, 0)
			if a.IsPageContentSorted {
				hc.SetText(hc.Text + "↑")
			} else {
				hc.SetText(hc.Text + "↓")
			}
			servicePageContent.SetBorderFocusColor(tcell.ColorSpringGreen)
			servicePage.AddItem(servicePageContent, 0, 6, true)
		}
		//sorting ec2 instances
		//73 - Key I
		if event.Rune() == 73 {
			servicePage.RemoveItemAtIndex(0)
			if a.IsPageContentSorted {
				sort.Sort(sort.Reverse(aws.ByInstanceId(ins)))
				a.IsPageContentSorted = false
			} else {
				sort.Sort(aws.ByInstanceId(ins))
				a.IsPageContentSorted = true
			}
			servicePageContent = a.DisplayEc2Instances(ins, sess)
			hc := servicePageContent.GetCell(0, 0)
			if a.IsPageContentSorted {
				hc.SetText(hc.Text + "↑")
			} else {
				hc.SetText(hc.Text + "↓")
			}
			servicePageContent.SetBorderFocusColor(tcell.ColorSpringGreen)
			servicePage.AddItem(servicePageContent, 0, 6, true)
		}

		//84 - Key T
		if event.Rune() == 84 {
			servicePage.RemoveItemAtIndex(0)
			if a.IsPageContentSorted {
				sort.Sort(sort.Reverse(aws.ByInstanceType(ins)))
				a.IsPageContentSorted = false
			} else {
				sort.Sort(aws.ByInstanceType(ins))
				a.IsPageContentSorted = true
			}
			servicePageContent = a.DisplayEc2Instances(ins, sess)
			hc := servicePageContent.GetCell(0, 2)
			if a.IsPageContentSorted {
				hc.SetText(hc.Text + "↑")
			} else {
				hc.SetText(hc.Text + "↓")
			}
			servicePageContent.SetBorderFocusColor(tcell.ColorSpringGreen)
			servicePage.AddItem(servicePageContent, 0, 6, true)
		}

		//76 - Key L
		if event.Rune() == 76 {
			servicePage.RemoveItemAtIndex(0)
			if a.IsPageContentSorted {
				sort.Sort(sort.Reverse(aws.ByLaunchTime(ins)))
				a.IsPageContentSorted = false
			} else {
				sort.Sort(aws.ByLaunchTime(ins))
				a.IsPageContentSorted = true
			}
			servicePageContent = a.DisplayEc2Instances(ins, sess)
			hc := servicePageContent.GetCell(0, 7)
			if a.IsPageContentSorted {
				hc.SetText(hc.Text + "↑")
			} else {
				hc.SetText(hc.Text + "↓")
			}
			servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
			servicePage.AddItem(servicePageContent, 0, 6, true)
		}
		return event
	})
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
				a.Flash().Info("Loading S3 Buckets...")
				servicePage.RemoveItemAtIndex(0)
				servicePageContent = a.DisplayS3Buckets(sess, buckets)
				servicePage.AddItem(servicePageContent, 0, 6, true)
				a.Application.SetFocus(servicePageContent)
				inputField.SetText("")

			case "EC2", "ec2", "Ec2", "eC2":
				a.Flash().Info("Loading EC2 instacnes...")
				servicePage.RemoveItemAtIndex(0)
				servicePageContent = a.DisplayEc2Instances(ins, sess)
				// ec2Page.AddItem(menuColFlex, 0, 2, false)
				servicePage.AddItem(servicePageContent, 0, 6, true)
				a.Application.SetFocus(servicePageContent)
				inputField.SetText("")

			default:
				inputField.SetText("")
				a.Flash().Err(fmt.Errorf("NO SERVICE..."))
			}
		}
	})

	inputField.SetBorder(true)
	inputField.SetBorderFocusColor(tcell.ColorSpringGreen)

	a.Views()["pAndRMenu"] = menuColFlex
	a.Views()["cmd"] = inputField
	a.Views()["content"] = servicePage

	main.AddItem(menuColFlex, 0, 2, false)
	main.AddItem(inputField, 0, 1, false)
	main.AddItem(servicePage, 0, 8, true)
	main.AddItem(flash, 0, 1, false)
	a.Main.AddPage("main", main, true, false)
	a.Main.AddPage("splash", ui.NewSplash("0.0.1"), true, true)
	return main
}

func (a *App) DisplayEc2Instances(ins []aws.EC2Resp, sess *session.Session) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetBorderFocusColor(tcell.ColorSpringGreen)
	// flex.AddItem(table, 0, 1, true).SetDirection(tview.FlexRow)
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
	}).SetSelectionChangedFunc(func(row int, column int) {
		table.SetSelectable(true, false)
	})

	table.SetSelectedFunc(func(row, column int) {
		insId := ins[row].InstanceId
		a.DisplayEc2InstanceJson(sess, insId)
	})

	return table
}

func (a *App) DisplayEc2InstanceJson(sess *session.Session, instanceId string) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	tvForEc2Json := tview.NewTextView()
	tvForEc2Json.SetBorder(true)
	tvForEc2Json.SetBorderFocusColor(tcell.ColorSpringGreen)
	tvForEc2Json.SetTitle(fmt.Sprintf(" EC2/%v/::[json] ", instanceId))
	tvForEc2Json.SetTitleColor(tcell.ColorLightSkyBlue)
	tvForEc2Json.SetText(aws.GetSingleInstance(*sess, instanceId).GoString())
	flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
	flex.AddItem(a.Views()["cmd"], 0, 1, false)
	flex.AddItem(tvForEc2Json, 0, 9, true)
	a.Main.AddAndSwitchToPage("main:ece2-json", flex, true)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			a.Main.SwitchToPage("main")
			a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
		}
		return event
	})
}

func (a *App) DisplayS3Buckets(sess *session.Session, buckets []aws.BucketResp) *tview.Table {
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
	}
	//r := 0
	a.Application.SetFocus(table)
	table.SetBorderFocusColor(tcell.ColorSpringGreen)

	table.SetSelectable(true, false)
	table.Select(1, 1).SetFixed(1, 1)
	s3DataT := tview.NewTable()
	s3DataT.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetCell(0, 2, tview.NewTableCell("Last modified").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetCell(0, 3, tview.NewTableCell("Size").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetCell(0, 4, tview.NewTableCell("Storage class").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetBorder(true)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { //Bucket979-0r
		if event.Key() == tcell.KeyEnter {
			flex.Clear()
			flex.RemoveItem(table)
			r, _ := table.GetSelection()
			bucketName := buckets[r-1].BucketName
			bucketInfo := aws.GetInfoAboutBucket(*sess, bucketName, "/", "")
			folderArrayInfo, fileArrayInfo := getLevelInfo(bucketInfo)
			indx := 0
			if len(folderArrayInfo) == 0 && len(fileArrayInfo) == 0 {
				s3DataT.SetTitle(bucketName)
				s3DataT.SetCell((indx + 2), 0, tview.NewTableCell("No Data Found inside Bucket").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
				flex.AddItem(a.Views()["cmd"], 0, 1, false)
				flex.AddItem(s3DataT, 0, 9, true)
				s3DataT.SetBorderFocusColor(tcell.ColorSpringGreen)
				a.Main.AddAndSwitchToPage("s3data", flex, true)
				s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { // Empty
					if event.Key() == tcell.KeyESC {
						a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
						flex.Clear()
						s3DataT.Clear()
						a.Main.RemovePage("s3data")
						a.Main.SwitchToPage("main")
						a.Application.SetFocus(table)
					}
					return event
				})
			} else {
				for _, bi := range bucketInfo.CommonPrefixes {
					keyA := strings.Split(*bi.Prefix, "/")
					s3DataT.SetTitle(bucketName)
					s3DataT.SetTitleColor(tcell.ColorYellow)
					s3DataT.SetCell((indx + 2), 0, tview.NewTableCell(keyA[len(keyA)-2]+"/").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
					s3DataT.SetCell((indx + 2), 1, tview.NewTableCell("Folder").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
					s3DataT.SetCell((indx + 2), 2, tview.NewTableCell("-").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
					s3DataT.SetCell((indx + 2), 3, tview.NewTableCell("0").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
					s3DataT.SetCell((indx + 2), 4, tview.NewTableCell("-").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
					indx++
				}

				for _, fi := range bucketInfo.Contents {
					keyA := strings.Split(*fi.Key, "/")
					s3DataT.SetTitle(bucketName)
					s3DataT.SetTitleColor(tcell.ColorYellow)
					s3DataT.SetCell((indx + 2), 0, tview.NewTableCell(keyA[len(keyA)-1]).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
					s3DataT.SetCell((indx + 2), 1, tview.NewTableCell("File").SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
					s3DataT.SetCell((indx + 2), 2, tview.NewTableCell(fi.LastModified.String()).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
					s3DataT.SetCell((indx + 2), 3, tview.NewTableCell(strconv.Itoa(int(*fi.Size))).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
					s3DataT.SetCell((indx + 2), 4, tview.NewTableCell(*fi.StorageClass).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
					indx++
				}

				flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
				flex.AddItem(a.Views()["cmd"], 0, 1, false)
				flex.AddItem(s3DataT, 0, 9, true)
				s3DataT.SetBorderFocusColor(tcell.ColorSpringGreen)
				a.Main.AddAndSwitchToPage("s3data", flex, true)

				if len(bucketInfo.CommonPrefixes) != 0 || len(bucketInfo.Contents) != 0 {
					s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { // Empty
						if event.Key() == tcell.KeyEnter { //d
							a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
							flex.Clear()
							s3DataT.Clear()
							a.Main.RemovePage("s3data")
							r, _ := s3DataT.GetSelection()
							cell := s3DataT.GetCell(r, 0)
							a.DisplayS3Objects(s3DataT, flex, cell.Text, fileArrayInfo, *sess, bucketName)
						} else if event.Key() == tcell.KeyESC {
							if strings.Count(folderArrayInfo[0], "/") == 1 {
								a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
								flex.Clear()
								s3DataT.Clear()
								a.Main.RemovePage("s3data")
								a.Main.SwitchToPage("main")
								a.Application.SetFocus(table)
							}
						}
						return event
					})
					a.Application.SetFocus(s3DataT)
					s3DataT.SetSelectable(true, false)
					s3DataT.Select(1, 1).SetFixed(1, 1)
				}
			}
		}
		return event
	})

	return table
}

func (a *App) DisplayS3Objects(s3DataTable *tview.Table, flex *tview.Flex, folderName string, fileArrayInfo []string, sess session.Session, bucketName string) {

	s3DataT := tview.NewTable()
	bucketInfo := aws.GetInfoAboutBucket(sess, bucketName, "/", folderName)
	_, fileArrayInfoTemp := getLevelInfo(bucketInfo)

	if len(bucketInfo.CommonPrefixes) != 0 || len(bucketInfo.Contents) != 0 {
		a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
		flex.Clear()
		s3DataT.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
		s3DataT.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
		s3DataT.SetCell(0, 2, tview.NewTableCell("Last modified").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
		s3DataT.SetCell(0, 3, tview.NewTableCell("Size").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
		s3DataT.SetCell(0, 4, tview.NewTableCell("Storage class").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
		s3DataT.SetBorder(true)
		indx := 0
		for _, bi := range bucketInfo.CommonPrefixes {
			keyA := strings.Split(*bi.Prefix, "/")
			s3DataT.SetTitle(bucketName + "/" + folderName)
			s3DataT.SetTitleColor(tcell.ColorYellow)
			s3DataT.SetCell((indx + 2), 0, tview.NewTableCell(keyA[len(keyA)-2]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			s3DataT.SetCell((indx + 2), 1, tview.NewTableCell("Folder").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			s3DataT.SetCell((indx + 2), 2, tview.NewTableCell("_").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			s3DataT.SetCell((indx + 2), 3, tview.NewTableCell("0").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			s3DataT.SetCell((indx + 2), 4, tview.NewTableCell("_").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			indx++
		}

		for _, fi := range bucketInfo.Contents {
			keyA := strings.Split(*fi.Key, "/")
			s3DataT.SetTitle(bucketName + "/" + folderName)
			s3DataT.SetTitleColor(tcell.ColorYellow)
			s3DataT.SetCell((indx + 2), 0, tview.NewTableCell(keyA[len(keyA)-1]).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
			s3DataT.SetCell((indx + 2), 1, tview.NewTableCell("File").SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
			s3DataT.SetCell((indx + 2), 2, tview.NewTableCell(fi.LastModified.String()).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
			s3DataT.SetCell((indx + 2), 3, tview.NewTableCell(strconv.Itoa(int(*fi.Size))).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
			s3DataT.SetCell((indx + 2), 4, tview.NewTableCell(*fi.StorageClass).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
			indx++
		}
		flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
		flex.AddItem(a.Views()["cmd"], 0, 1, false)
		flex.AddItem(s3DataT, 0, 9, true)
		s3DataT.SetBorderFocusColor(tcell.ColorSpringGreen)

		s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { //Tiger
			if event.Key() == tcell.KeyEnter {
				r, _ := s3DataT.GetSelection()
				cell := s3DataT.GetCell(r, 0)
				foldN := cell.Text
				a.DisplayS3Objects(s3DataT, flex, folderName+foldN+"/", fileArrayInfoTemp, sess, bucketName)
			}

			if event.Key() == tcell.KeyESC {
				r, _ := s3DataT.GetSelection()
				cell := s3DataT.GetCell(r, 1)
				cellTxt := cell.Text
				passF := ""
				if strings.Count(folderName, "/") < 1 {
					a.Main.RemovePage("s3dataView")
					a.Main.SwitchToPage("main")
				} else {
					if cellTxt == "File" {
						slashed := strings.Split(folderName, "/")
						for i := 0; i < len(slashed)-2; i++ {
							passF = passF + slashed[i] + "/"
						}
					} else {
						slashed := strings.Split(folderName, "/")
						for i := 0; i < len(slashed)-2; i++ {
							passF = passF + slashed[i] + "/"
						}

					}
					a.DisplayS3Objects(s3DataT, flex, passF, fileArrayInfo, sess, bucketName)
				}
			}
			return event
		})
		a.Application.SetFocus(s3DataT)
		s3DataT.SetSelectable(true, false)
		s3DataT.Select(1, 1).SetFixed(1, 1)
		a.Main.AddAndSwitchToPage("s3dataView", flex, true)
	}

}

func getLevelInfo(bucketInfo *s3.ListObjectsV2Output) ([]string, []string) {
	var folderArrayInfo []string
	var fileArrayInfo []string
	for _, i := range bucketInfo.CommonPrefixes {
		folderArrayInfo = append(folderArrayInfo, *i.Prefix)
	}

	for i := 0; i < len(bucketInfo.Contents); i++ {
		fileArrayInfo = append(fileArrayInfo, *bucketInfo.Contents[i].Key)
	}
	return folderArrayInfo, fileArrayInfo
}

func (a *App) tempLayout(ctx context.Context) {
	flash := ui.NewFlash(a.App)
	go flash.Watch(ctx, a.Flash().Channel())

	main := tview.NewFlex().SetDirection(tview.FlexRow)
	main.AddItem(a.statusIndicator(), 1, 1, false)
	a.Main.AddPage("main", main, true, false)
	a.Main.AddPage("splash", ui.NewSplash("0.0.1"), true, true)
	main.AddItem(flash, 1, 1, false)
	a.toggleHeader(true)
}

// QueueUpdateDraw queues up a ui action and redraw the ui.
func (a *App) QueueUpdateDraw(f func()) {
	if a.Application == nil {
		return
	}
	go func() {
		a.Application.QueueUpdateDraw(f)
	}()
}

func (a *App) Run() error {
	//a.Resume()
	go func() {
		<-time.After(splashDelay)
		a.QueueUpdateDraw(func() {
			a.Main.SwitchToPage("main")
			a.Application.SetFocus(a.Main.CurrentPage().Item)
		})
	}()

	if err := a.Application.Run(); err != nil {
		return err
	}

	return nil
}

func (a *App) toggleHeader(header bool) {
	a.showHeader = header

	flex, ok := a.Main.GetPrimitive("main").(*tview.Flex)
	if !ok {
		log.Fatal().Msg("Expecting valid flex view")
	}
	if a.showHeader {
		flex.RemoveItemAtIndex(0)
		flex.AddItemAtIndex(0, a.buildHeader(), 7, 1, false)
	} else {
		flex.RemoveItemAtIndex(0)
		flex.AddItemAtIndex(0, a.statusIndicator(), 1, 1, false)
	}
}

func (a *App) buildHeader() tview.Primitive {
	header := tview.NewFlex()
	header.SetDirection(tview.FlexColumn)
	if !a.showHeader {
		return header
	}
	header.AddItem(a.Menu(), 0, 1, false)

	return header
}

func (a *App) bindKeys() {
	a.AddActions(ui.KeyActions{
		tcell.KeyCtrlE: ui.NewKeyAction("ToggleHeader", a.toggleHeaderCmd, false),
	})
}

func (a *App) toggleHeaderCmd(evt *tcell.EventKey) *tcell.EventKey {

	a.QueueUpdateDraw(func() {
		a.showHeader = !a.showHeader
		a.toggleHeader(a.showHeader)
	})

	return nil
}

func (a *App) statusIndicator() *ui.StatusIndicator {
	return a.Views()["statusIndicator"].(*ui.StatusIndicator)
}
