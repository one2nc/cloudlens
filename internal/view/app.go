package view

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/one2nc/cloud-lens/internal/config"

	"github.com/one2nc/cloud-lens/internal/ui"
)

type S3Response struct {
	Name         string
	Type         string
	LastModified string
	Size         int64
	StorageClass string
}

type App struct {
	*ui.App
	IsPageContentSorted bool
}

func NewApp() *App {
	a := App{App: ui.NewApp(), IsPageContentSorted: false}
	a.Views()["statusIndicator"] = ui.NewStatusIndicator(a.App)
	return &a
}

func (a *App) Init() error {
	_ = context.WithValue(context.Background(), internal.KeyApp, a)

	a.App.Init()
	a.layout()
	a.App.Run()
	return nil
}

func (a *App) layout() *tview.Flex {
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
	servicePage := tview.NewFlex().SetDirection(tview.FlexRow)
	servicePageContent := a.DisplayEc2Instances(ins, sess)

	profileDropdown := tview.NewDropDown().
		SetLabel("Profile ▼ ").
		SetOptions(profiles, func(text string, index int) {
			currentProfile = &text
			sess, _ = config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)
			ins, _ = aws.GetInstances(*sess)
			buckets, _ = aws.ListBuckets(*sess)
			textv.SetText("🌈🌧️ cloudlens profile..." + fmt.Sprintf("%v", a.Main.Pages.CurrentPage()))
			if servicePage.ItemAt(0) != nil {
				servicePage.RemoveItemAtIndex(0)
			}
			servicePageContent = a.DisplayEc2Instances(ins, sess)
			servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
			servicePage.AddItem(servicePageContent, 0, 6, true)
		})
	profileDropdown.SetBorderFocusColor(tcell.ColorSpringGreen)
	profileDropdown.SetCurrentOption(0)
	//profileDropdown.SetTextOptions(" ▲ ", "", " ▼ ", " ", "-")
	profileDropdown.SetBorderPadding(2, 0, 0, 0)
	//profileDropdown.SetBorder(true)

	regionDropdown := tview.NewDropDown().
		SetLabel("Region ▼ ").
		SetOptions(regions, func(text string, index int) {
			currentRegion = &text
			sess, _ = config.GetSession(*currentProfile, *currentRegion, cfg.AwsConfig)
			ins, _ = aws.GetInstances(*sess)
			buckets, _ = aws.ListBuckets(*sess)
			textv.SetText("🌈🌧️ cloudlens region... " + fmt.Sprintf("%v", ins))
			if servicePage.ItemAt(0) != nil {
				servicePage.RemoveItemAtIndex(0)
			}
			servicePageContent = a.DisplayEc2Instances(ins, sess)
			servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
			servicePage.AddItem(servicePageContent, 0, 6, true)
		})
	regionDropdown.SetBorderFocusColor(tcell.ColorSpringGreen)
	regionDropdown.SetCurrentOption(0)
	//regionDropdown.SetTextOptions(" ▲ ", "", " ▼ ", " ", "-")
	regionDropdown.SetBorderPadding(0, 0, 0, 0)
	//regionDropdown.SetBorder(true)

	ddRowFlex.AddItem(profileDropdown, 0, 1, false)
	ddRowFlex.AddItem(regionDropdown, 0, 1, false)

	menuColFlex.AddItem(tview.NewBox(), 0, 1, false)
	menuColFlex.AddItem(tview.NewBox(), 0, 1, false)
	menuColFlex.AddItem(tview.NewBox(), 0, 1, false)
	menuColFlex.AddItem(ddRowFlex, 0, 1, false)

	servicePage = tview.NewFlex().SetDirection(tview.FlexRow)
	servicePageContent = a.DisplayEc2Instances(ins, sess)
	servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
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
				hc.SetText(hc.Text+"↑")
			}else{
				hc.SetText(hc.Text+"↓")
			}
			servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
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
			servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
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
			servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
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
				servicePage.RemoveItemAtIndex(0)
				servicePageContent = a.DisplayS3Buckets(sess, buckets)
				servicePageContent.SetBorderFocusColor(tcell.ColorDarkSeaGreen)
				servicePage.AddItem(servicePageContent, 0, 6, true)
				inputField.SetText("")

			case "EC2", "ec2", "Ec2", "eC2":
				servicePage.RemoveItemAtIndex(0)
				servicePageContent = a.DisplayEc2Instances(ins, sess)
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
	return main
}

func (a *App) DisplayEc2Instances(ins []aws.EC2Resp, sess *session.Session) *tview.Table {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
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
	r := 0
	table.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectionChangedFunc(func(row int, column int) {
		table.SetSelectable(true, false)
		r = row - 1
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 100 {
			insId := ins[r].InstanceId
			newPage := tview.NewTextView()
			newPage.SetBorder(true)
			newPage.SetTitle(" JSON ")
			newPage.SetText(aws.GetSingleInstance(*sess, insId).GoString())
			desc := tview.NewTextView()
			desc.SetText("<esc> shift to previous page")
			flex.AddItem(desc, 0, 1, true)
			flex.AddItem(newPage, 0, 10, true)
			a.Main.AddAndSwitchToPage("json", flex, true)
			flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyESC {
					flex.RemoveItem(desc)
					flex.RemoveItem(newPage)
					a.Main.RemovePage("json")
					a.Main.ShowPage("main")
				}
				return event
			})
		}
		return event
	})
	return table
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
	table.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectionChangedFunc(func(row int, column int) {
		table.SetSelectable(true, false)
		//	r = row - 1
	})
	s3DataT := tview.NewTable()
	s3DataT.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetCell(0, 2, tview.NewTableCell("Last modified").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetCell(0, 3, tview.NewTableCell("Size").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetCell(0, 4, tview.NewTableCell("Storage class").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	s3DataT.SetBorder(true)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 100 {
			flex.RemoveItem(table)
			r, _ := table.GetSelection()
			bucketName := buckets[r-1].BucketName
			bucketInfo := aws.GetInfoAboutBucket(*sess, bucketName, "/", "")
			_, fileArrayInfo := getLevelInfo(bucketInfo)
			indx := 0
			for _, bi := range bucketInfo.CommonPrefixes {
				s3DataT.SetCell((indx + 2), 0, tview.NewTableCell(*bi.Prefix).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 1, tview.NewTableCell("Folder").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 2, tview.NewTableCell("_").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 3, tview.NewTableCell("0").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 4, tview.NewTableCell("_").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				indx++
			}
			for _, fi := range bucketInfo.Contents {
				s3DataT.SetCell((indx + 2), 0, tview.NewTableCell(*fi.Key).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 1, tview.NewTableCell("File").SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 2, tview.NewTableCell(fi.LastModified.String()).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 3, tview.NewTableCell(strconv.Itoa(int(*fi.Size))).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 4, tview.NewTableCell(*fi.StorageClass).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				indx++

			}
			s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Rune() == 100 { //d
					flex.RemoveItem(s3DataT)
					r, _ := s3DataT.GetSelection()
					cell := s3DataT.GetCell(r, 0)
					s3DataTR := tview.NewTable()
					s3DataTR.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetCell(0, 2, tview.NewTableCell("Last modified").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetCell(0, 3, tview.NewTableCell("Size").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetCell(0, 4, tview.NewTableCell("Storage class").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetBorder(true)
					flex.AddItem(s3DataTR, 0, 1, false)
					a.Main.AddAndSwitchToPage("s3dataView", flex, true)
					a.inputCaptureS3(s3DataTR, flex, cell.Text, fileArrayInfo, *sess, bucketName)
				}
				return event
			})
			//	}
			s3DataT.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					s3DataT.SetSelectable(true, false)
				}
			}).SetSelectionChangedFunc(func(row int, column int) {
				s3DataT.SetSelectable(true, false)
			})

			flex.AddItem(s3DataT, 0, 10, true)
			a.Main.AddAndSwitchToPage("s3data", flex, true)
		}
		return event
	})

	return table
}

func (a *App) inputCaptureS3(s3DataT *tview.Table, flex *tview.Flex, folderName string, fileArrayInfo []string, sess session.Session, bucketName string) {
	s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 100 {
			bucketInfo := aws.GetInfoAboutBucket(sess, bucketName, "/", folderName)
			_, fileArrayInfoTemp := getLevelInfo(bucketInfo)
			indx := 0
			for _, bi := range bucketInfo.CommonPrefixes {
				keyA := strings.Split(*bi.Prefix, "/")
				s3DataT.SetCell((indx + 2), 0, tview.NewTableCell(keyA[len(keyA)-2]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 1, tview.NewTableCell("Folder").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 2, tview.NewTableCell("-").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 3, tview.NewTableCell("0").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 4, tview.NewTableCell("-").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				indx++

			}
			for _, fi := range bucketInfo.Contents {
				keyA := strings.Split(*fi.Key, "/")
				s3DataT.SetCell((indx + 2), 0, tview.NewTableCell(keyA[len(keyA)-1]).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 1, tview.NewTableCell("File").SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 2, tview.NewTableCell(fi.LastModified.String()).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 3, tview.NewTableCell(strconv.Itoa(int(*fi.Size))).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				s3DataT.SetCell((indx + 2), 4, tview.NewTableCell(*fi.StorageClass).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignCenter))
				indx++
			}
			s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Rune() == 100 {
					flex.RemoveItem(s3DataT)
					r, _ := s3DataT.GetSelection()
					cell := s3DataT.GetCell(r, 0)
					s3DataTR := tview.NewTable()
					s3DataTR.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetCell(0, 2, tview.NewTableCell("Last modified").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetCell(0, 3, tview.NewTableCell("Size").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetCell(0, 4, tview.NewTableCell("Storage class").SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
					s3DataTR.SetBorder(true)
					flex.AddItem(s3DataTR, 0, 1, false)
					a.Main.AddAndSwitchToPage("s3dataView", flex, true)
					foldN := cell.Text
					a.inputCaptureS3(s3DataTR, flex, folderName+foldN+"/", fileArrayInfoTemp, sess, bucketName)
				}
				return event
			})
			s3DataT.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					s3DataT.SetSelectable(true, false)
				}
			}).SetSelectionChangedFunc(func(row int, column int) {
				s3DataT.SetSelectable(true, false)
			})
		}
		return event
	})
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

func (a *App) statusIndicator() *ui.StatusIndicator {
	return a.Views()["statusIndicator"].(*ui.StatusIndicator)
}
