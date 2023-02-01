package view

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/derailed/tview"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/aws"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/one2nc/cloud-lens/internal/model"

	"github.com/one2nc/cloud-lens/internal/ui"
	"github.com/rs/zerolog/log"
)

const (
	splashDelay = 1 * time.Second
)

type App struct {
	*ui.App
	Content             *PageStack
	cancelFn            context.CancelFunc
	showHeader          bool
	IsPageContentSorted bool
}

type BucketInfoJson struct {
	LifeCycleJson  *s3.ServerSideEncryptionConfiguration
	EncryptionJson []*s3.LifecycleRule
}

func NewApp() *App {
	a := App{
		App:                 ui.NewApp(),
		Content:             NewPageStack(),
		IsPageContentSorted: false,
	}
	a.Views()["statusIndicator"] = ui.NewStatusIndicator(a.App)
	return &a
}

func (a *App) Init() error {
	ctx := context.WithValue(context.Background(), internal.KeyApp, a)
	if err := a.Content.Init(ctx); err != nil {
		return err
	}
	a.Content.Stack.AddListener(a.Menu())
	a.App.Init()
	a.SetInputCapture(a.keyboard)
	a.bindKeys()
	a.layout(ctx)
	//a.tempLayout(ctx)
	return nil
}

func (a *App) layout(ctx context.Context) *tview.Flex {
	cfg, _ := config.Get()

	main := tview.NewFlex().SetDirection(tview.FlexRow)
	//----Flash Mesages----
	flash := ui.NewFlash(a.App)
	go flash.Watch(ctx, a.Flash().Channel())

	a.Views()["flash"] = flash

	//------menu-----
	menuColFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	//menuColFlex.SetBorder(true)
	ddRowFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	//ddRowFlex.SetBorder(true)

	var ins []aws.EC2Resp
	var buckets []aws.BucketResp
	var secGrp []*ec2.SecurityGroup
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
			secGrp = aws.GetSecGrps(*sess)
			if servicePage.ItemAt(0) != nil {
				servicePage.RemoveItemAtIndex(0)
			}

			if servicePageContent.GetCell(0, 1).Text == "Creation-Time" {
				servicePageContent = a.DisplayS3Buckets(sess, buckets)
			} else if servicePageContent.GetCell(0, 0).Text == "Group-Id" {
				servicePageContent = a.DisplaySecurityGroup(sess, secGrp)
			} else {
				servicePageContent = a.DisplayEc2Instances(ins, sess)
			}
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

	main.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if servicePageContent.GetCell(0, 0).Text == "Instance-Id" {
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
				a.Application.SetFocus(servicePageContent)
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
				a.Application.SetFocus(servicePageContent)
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
				a.Application.SetFocus(servicePageContent)
			}
		}

		if servicePageContent.GetCell(0, 1).Text == "Creation-Time" {
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
				a.Application.SetFocus(servicePageContent)
			}
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
	servicePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			a.Application.SetFocus(inputField)
		}
		return event
	})
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
				servicePage.AddItem(servicePageContent, 0, 6, true)
				a.Application.SetFocus(servicePageContent)
				inputField.SetText("")

			case "Security Group", "sg":
				a.Flash().Info("Loading Security Groups...")
				servicePage.RemoveItemAtIndex(0)
				secGrp = aws.GetSecGrps(*sess)
				servicePageContent = a.DisplaySecurityGroup(sess, secGrp)
				// ec2Page.AddItem(menuColFlex, 0, 2, false)
				servicePage.AddItem(servicePageContent, 0, 6, true)
				a.Application.SetFocus(servicePageContent)
				inputField.SetText("")

			default:
				inputField.SetText("")
				a.Flash().Err(fmt.Errorf("NO SERVICE..."))
			}

			servicePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyTab {
					a.Application.SetFocus(inputField)
				}
				return event
			})

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

func (a *App) DisplaySecurityGroup(sess *session.Session, secGrp []*ec2.SecurityGroup) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetBorderFocusColor(tcell.ColorSpringGreen)
	a.setTableHeaderForSecGrp(table)
	a.setTableContentForSecGrp(table, secGrp)
	table.SetSelectable(true, false)
	a.Application.SetFocus(table)
	table.Select(1, 1).SetFixed(1, 1)
	table.SetSelectedFunc(func(row, column int) {
		grpId := secGrp[row-2].GroupId
		a.DisplaySecGrpJson(sess, *grpId)
	})
	return table
}

func (a *App) DisplaySecGrpJson(sess *session.Session, grpId string) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	tvForSecGrpJson := tview.NewTextView()
	tvForSecGrpJson.SetBorder(true)
	tvForSecGrpJson.SetBorderFocusColor(tcell.ColorSpringGreen)
	tvForSecGrpJson.SetTitle(fmt.Sprintf(" SecurityGroup/%v/::[json] ", grpId))
	tvForSecGrpJson.SetTitleColor(tcell.ColorLightSkyBlue)
	tvForSecGrpJson.SetText(aws.GetSingleSecGrp(*sess, grpId).GoString())
	flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
	inputPrompt := tview.NewInputField().
		SetLabel("🐶>").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return true // accept any input
		})
	inputPrompt.SetFieldBackgroundColor(tcell.ColorBlack)
	inputPrompt.SetBorder(true)

	flex.AddItem(inputPrompt, 0, 1, false)
	buckets, _ := aws.ListBuckets(*sess)
	ins, _ := aws.GetInstances(*sess)
	a.SearchUtility(inputPrompt, sess, buckets, flex, nil, ins)
	flex.AddItem(tvForSecGrpJson, 0, 9, true)
	flex.AddItem(a.FlashView(), 0, 1, false)
	a.Main.AddAndSwitchToPage("main:secGrpjson", flex, true)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			flex.Clear()
			a.Main.SwitchToPage("main")
			a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
		}
		if event.Key() == tcell.KeyTab {
			a.Application.SetFocus(inputPrompt)
		}

		return event
	})
	tvForSecGrpJson.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 67 {
			text := tvForSecGrpJson.GetText(true)
			clipboard.WriteAll(text)
			a.Flash().Info("Text Copied.")
		}
		return event
	})
}

func (a *App) setTableHeaderForSecGrp(secGrpTable *tview.Table) *tview.Table {
	secGrpTable.SetTitleColor(tcell.ColorYellow)
	secGrpTable.SetCell(0, 0, tview.NewTableCell("Group-Id").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	secGrpTable.SetCell(0, 1, tview.NewTableCell("Group-Name").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	secGrpTable.SetCell(0, 2, tview.NewTableCell("Description").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	secGrpTable.SetCell(0, 3, tview.NewTableCell("Owner-Id").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	secGrpTable.SetCell(0, 4, tview.NewTableCell("Vpc-Id").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))

	return secGrpTable
}

func (a *App) setTableContentForSecGrp(table *tview.Table, secGrp []*ec2.SecurityGroup) *tview.Table {
	indx := 0
	for _, grp := range secGrp {
		table.SetCell((indx + 2), 0, tview.NewTableCell(*grp.GroupId).SetAlign(tview.AlignCenter))
		table.SetCell((indx + 2), 1, tview.NewTableCell(*grp.GroupName).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((indx + 2), 2, tview.NewTableCell(*grp.Description).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((indx + 2), 3, tview.NewTableCell(*grp.OwnerId).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((indx + 2), 4, tview.NewTableCell(*grp.VpcId).SetExpansion(1).SetAlign(tview.AlignCenter))
		indx++
	}
	table.SetBorderFocusColor(tcell.ColorSpringGreen)
	return table
}

func (a *App) DisplayEc2Instances(ins []aws.EC2Resp, sess *session.Session) *tview.Table {

	table := tview.NewTable()
	table.SetBorder(true)
	table.SetBorderFocusColor(tcell.ColorSpringGreen)
	// flex.AddItem(table, 0, 1, true).SetDirection(tview.FlexRow)
	//table data
	table.SetCell(0, 0, tview.NewTableCell("Instance-Id").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Instance-State").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Instance-Type").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Availability-zone").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("Public-DNS").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 5, tview.NewTableCell("Public-IPV4").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 6, tview.NewTableCell("Monitoring-State").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))
	table.SetCell(0, 7, tview.NewTableCell("Launch-Time").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignCenter))

	for i, in := range ins {
		table.SetCell((i + 1), 0, tview.NewTableCell(in.InstanceId).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 1, tview.NewTableCell(in.InstanceState).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 2, tview.NewTableCell(in.InstanceType).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 3, tview.NewTableCell(in.AvailabilityZone).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 4, tview.NewTableCell(in.PublicDNS).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 5, tview.NewTableCell(in.PublicIPv4).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 6, tview.NewTableCell(in.MonitoringState).SetExpansion(1).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 7, tview.NewTableCell(in.LaunchTime).SetExpansion(1).SetAlign(tview.AlignCenter))
	}
	table.SetSelectable(true, false)
	a.Application.SetFocus(table)
	if table.GetCell(1, 1).Text != "" {
		table.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				table.SetSelectable(true, false)
			}
		}).SetSelectionChangedFunc(func(row int, column int) {
			table.SetSelectable(true, false)
		})

		table.SetSelectedFunc(func(row, column int) {
			insId := ins[row-1].InstanceId
			a.DisplayEc2InstanceJson(sess, insId)
		})
	}
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
	inputPrompt := tview.NewInputField().
		SetLabel("🐶>").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return true // accept any input
		})
	inputPrompt.SetFieldBackgroundColor(tcell.ColorBlack)
	inputPrompt.SetBorder(true)

	flex.AddItem(inputPrompt, 0, 1, false)
	buckets, _ := aws.ListBuckets(*sess)
	ins, _ := aws.GetInstances(*sess)
	a.SearchUtility(inputPrompt, sess, buckets, flex, nil, ins)
	flex.AddItem(tvForEc2Json, 0, 9, true)
	a.Main.AddAndSwitchToPage("main:ece2-json", flex, true)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			a.Main.SwitchToPage("main")
			a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
		}
		if event.Key() == tcell.KeyTab {
			a.Application.SetFocus(inputPrompt)
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
	table.SetCell(0, 0, tview.NewTableCell("Bucket-Name").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignLeft))
	table.SetCell(0, 1, tview.NewTableCell("Creation-Time").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignLeft))
	for i, b := range buckets {
		table.SetCell((i + 1), 0, tview.NewTableCell(b.BucketName).SetAlign(tview.AlignLeft))
		table.SetCell((i + 1), 1, tview.NewTableCell(fmt.Sprintf("%v", b.CreationTime)).SetAlign(tview.AlignLeft))
	}
	//r := 0
	table.SetSelectable(true, false)
	table.Select(1, 1).SetFixed(1, 1)
	a.Application.SetFocus(table)
	table.SetBorderFocusColor(tcell.ColorSpringGreen)
	s3DataT := tview.NewTable()
	s3DataT.SetBorder(true)

	if table.GetCell(1, 1).Text != "" {
		table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { //Bucket979-0r
			if event.Key() == tcell.KeyEnter {
				flex.Clear()
				flex.RemoveItem(table)
				r, _ := table.GetSelection()
				bucketName := buckets[r-1].BucketName
				bucketInfo := aws.GetInfoAboutBucket(*sess, bucketName, "/", "")
				folderArrayInfo, fileArrayInfo := getBucLevelInfo(bucketInfo)
				if len(folderArrayInfo) == 0 && len(fileArrayInfo) == 0 {
					a.DisplayS3ObjectForEmptyBuc(s3DataT, flex, bucketName, *sess)
				} else {
					a.setTableHeaderForS3(s3DataT, bucketName)
					a.setTableContentForS3(s3DataT, bucketInfo.CommonPrefixes, bucketInfo.Contents)

					flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)

					//extract to method
					inputPrompt := tview.NewInputField().
						SetLabel("🐶>").
						SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
							return true // accept any input
						})
					inputPrompt.SetFieldBackgroundColor(tcell.ColorBlack)
					inputPrompt.SetBorder(true)

					flex.AddItem(inputPrompt, 0, 1, false)
					flex.AddItem(s3DataT, 0, 9, true)
					a.Main.AddAndSwitchToPage("s3data", flex, true)

					//extract to method
					flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						if event.Key() == tcell.KeyTab {
							a.Application.SetFocus(inputPrompt)
						}
						return event
					})

					ins, _ := aws.GetInstances(*sess)
					a.SearchUtility(inputPrompt, sess, buckets, flex, table, ins)

					if len(bucketInfo.CommonPrefixes) != 0 || len(bucketInfo.Contents) != 0 {
						s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { // Empty
							if event.Key() == tcell.KeyEnter { //d
								r, _ := s3DataT.GetSelection()
								cell := s3DataT.GetCell(r, 0)
								flex.Clear()
								s3DataT.Clear()
								a.Main.RemovePage("s3data")
								a.DisplayS3Objects(s3DataT, flex, cell.Text+"/", fileArrayInfo, *sess, bucketName)
							} else if event.Key() == tcell.KeyESC {
								if strings.Count(folderArrayInfo[0], "/") == 1 {
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
			} else if event.Rune() == 68 { // press D
				r, _ := table.GetSelection()
				cell := table.GetCell(r, 0)
				flex.Clear()
				s3DataT.Clear()
				a.Main.RemovePage("s3data")
				a.DisplayS3Json(sess, cell.Text)
			}
			return event
		})
	}

	return table
}

func (a *App) DisplayS3Json(sess *session.Session, bucketName string) {
	json1 := aws.GetBuckEncryption(*sess, bucketName)
	json2 := aws.GetBuckLifecycle(*sess, bucketName)
	res := concatJson(json1, json2.Rules)
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	tvForS3Json := tview.NewTextView()
	tvForS3Json.SetBorder(true)
	tvForS3Json.SetBorderFocusColor(tcell.ColorSpringGreen)
	tvForS3Json.SetTitle(bucketName)
	tvForS3Json.SetTitleColor(tcell.ColorLightSkyBlue)
	tvForS3Json.SetText(res)
	flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
	inputPrompt := tview.NewInputField().
		SetLabel("🐶>").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return true
		})
	inputPrompt.SetFieldBackgroundColor(tcell.ColorBlack)
	inputPrompt.SetBorder(true)

	flex.AddItem(inputPrompt, 0, 1, false)
	buckets, _ := aws.ListBuckets(*sess)
	ins, _ := aws.GetInstances(*sess)
	a.SearchUtility(inputPrompt, sess, buckets, flex, nil, ins)
	flex.AddItem(tvForS3Json, 0, 8, true)
	flex.AddItem(a.FlashView(), 0, 1, false)
	a.Main.AddAndSwitchToPage("s3Json", flex, true)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			tvForS3Json.SetTextColor(tcell.ColorWhite)
			a.Main.SwitchToPage("main")
			a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
		}
		if event.Key() == tcell.KeyTab {
			a.Application.SetFocus(inputPrompt)
		}
		return event
	})
	tvForS3Json.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 67 {
			text := tvForS3Json.GetText(true)
			clipboard.WriteAll(text)
			a.Flash().Info("Text Copied.")
		}
		return event
	})

}

func (a *App) DisplayS3Objects(s3DataTable *tview.Table, flex *tview.Flex, folderName string, fileArrayInfo []string, sess session.Session, bucketName string) {
	s3DataT := tview.NewTable()
	s3DataT.SetBorder(true)
	bucketInfo := aws.GetInfoAboutBucket(sess, bucketName, "/", folderName)
	_, fileArrayInfoTemp := getBucLevelInfo(bucketInfo)

	if len(bucketInfo.CommonPrefixes) != 0 || len(bucketInfo.Contents) != 0 {
		flex.Clear()
		a.setTableHeaderForS3(s3DataT, bucketName+"/"+folderName)
		a.setTableContentForS3(s3DataT, bucketInfo.CommonPrefixes, bucketInfo.Contents)

		flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
		inputPrompt := tview.NewInputField().
			SetLabel("🐶>").
			SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
				return true // accept any input
			})
		inputPrompt.SetFieldBackgroundColor(tcell.ColorBlack)
		inputPrompt.SetBorder(true)

		flex.AddItem(inputPrompt, 0, 1, false)
		flex.AddItem(s3DataT, 0, 9, true)
		buckets, _ := aws.ListBuckets(sess)
		a.SearchUtility(inputPrompt, &sess, buckets, flex, s3DataT, nil)
		flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				a.Application.SetFocus(inputPrompt)
			}
			return event
		})
		s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

func (a *App) DisplayS3ObjectForEmptyBuc(s3DataT *tview.Table, flex *tview.Flex, bucketName string, sess session.Session) {
	s3DataT.SetTitle(bucketName)
	s3DataT.SetTitleColor(tcell.ColorYellow)
	s3DataT.SetCell(1, 0, tview.NewTableCell("No Data Found inside the  Bucket").SetTextColor(tcell.ColorPeachPuff).SetAlign(tview.AlignCenter))
	flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
	flex.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
	inputPrompt := tview.NewInputField().
		SetLabel("🐶>").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return true // accept any input
		})
	inputPrompt.SetFieldBackgroundColor(tcell.ColorBlack)
	inputPrompt.SetBorder(true)

	flex.AddItem(inputPrompt, 0, 1, false)
	flex.AddItem(s3DataT, 0, 9, true)
	s3DataT.SetBorderFocusColor(tcell.ColorSpringGreen)
	buckets, _ := aws.ListBuckets(sess)
	ins, _ := aws.GetInstances(sess)
	a.SearchUtility(inputPrompt, &sess, buckets, flex, s3DataT, ins)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			a.Application.SetFocus(inputPrompt)
		}
		return event
	})
	a.Main.AddAndSwitchToPage("s3data", flex, true)
	s3DataT.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { // Empty
		if event.Key() == tcell.KeyESC {
			a.Application.SetFocus(a.Views()["content"].(*tview.Flex).ItemAt(0))
			flex.Clear()
			s3DataT.Clear()
			a.Main.RemovePage("s3data")
			a.Main.SwitchToPage("main")
		}
		return event
	})
}

func (a *App) SearchUtility(inputField *tview.InputField, sess *session.Session, buckets []aws.BucketResp, servicePage *tview.Flex, servicePageContent *tview.Table, ins []aws.EC2Resp) {
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			serviceName := inputField.GetText()
			switch serviceName {
			case "S3", "s3":
				a.Flash().Info("Loading S3 Buckets...")
				servicePage.Clear()
				servicePageContent = a.DisplayS3Buckets(sess, buckets)
				servicePage.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
				servicePage.AddItem(inputField, 0, 1, false)
				servicePage.AddItem(servicePageContent, 0, 8, true)
				a.Application.SetFocus(servicePageContent)
				inputField.SetText("")

			case "EC2", "ec2", "Ec2", "eC2":
				a.Flash().Info("Loading EC2 instacnes...")
				servicePage.Clear()
				servicePageContent = a.DisplayEc2Instances(ins, sess)
				servicePage.AddItem(a.Views()["pAndRMenu"], 0, 2, false)
				servicePage.AddItem(inputField, 0, 1, false)
				servicePage.AddItem(servicePageContent, 0, 8, true)
				a.Application.SetFocus(servicePageContent)
				inputField.SetText("")

			default:
				inputField.SetText("")
				a.Flash().Err(fmt.Errorf("NO SERVICE..."))
			}
		}
		servicePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				a.Application.SetFocus(inputField)
			}
			return event
		})
	})
}

func (a *App) setTableHeaderForS3(s3DataT *tview.Table, tableTitle string) *tview.Table {
	s3DataT.SetTitle(tableTitle)
	s3DataT.SetTitleColor(tcell.ColorYellow)
	s3DataT.SetCell(0, 0, tview.NewTableCell("Name").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignLeft))
	s3DataT.SetCell(0, 1, tview.NewTableCell("Type").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignLeft))
	s3DataT.SetCell(0, 2, tview.NewTableCell("Last modified").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignLeft))
	s3DataT.SetCell(0, 3, tview.NewTableCell("Size").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignLeft))
	s3DataT.SetCell(0, 4, tview.NewTableCell("Storage class").SetExpansion(1).SetSelectable(false).SetTextColor(tcell.ColorOrangeRed).SetAlign(tview.AlignLeft))

	return s3DataT
}

func (a *App) setTableContentForS3(table *tview.Table, Folder []*s3.CommonPrefix, File []*s3.Object) *tview.Table {
	indx := 0
	for _, bi := range Folder {
		keyA := strings.Split(*bi.Prefix, "/")
		table.SetCell((indx + 2), 0, tview.NewTableCell(keyA[len(keyA)-2]).SetExpansion(1).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
		table.SetCell((indx + 2), 1, tview.NewTableCell("Folder").SetExpansion(1).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
		table.SetCell((indx + 2), 2, tview.NewTableCell("_").SetExpansion(1).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
		table.SetCell((indx + 2), 3, tview.NewTableCell("0").SetExpansion(1).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
		table.SetCell((indx + 2), 4, tview.NewTableCell("_").SetExpansion(1).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
		indx++
	}

	for _, fi := range File {
		keyA := strings.Split(*fi.Key, "/")
		table.SetCell((indx + 2), 0, tview.NewTableCell(keyA[len(keyA)-1]).SetExpansion(1).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignLeft))
		table.SetCell((indx + 2), 1, tview.NewTableCell("File").SetExpansion(1).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignLeft))
		IST := getIST(fi.LastModified)
		table.SetCell((indx + 2), 2, tview.NewTableCell(IST).SetExpansion(1).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignLeft))
		size := humanize.Bytes(uint64(*fi.Size))
		table.SetCell((indx + 2), 3, tview.NewTableCell(size).SetExpansion(1).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignLeft))
		table.SetCell((indx + 2), 4, tview.NewTableCell(*fi.StorageClass).SetExpansion(1).SetTextColor(tcell.ColorAntiqueWhite).SetAlign(tview.AlignLeft))
		indx++
	}
	table.SetBorderFocusColor(tcell.ColorSpringGreen)

	return table
}

func getBucLevelInfo(bucketInfo *s3.ListObjectsV2Output) ([]string, []string) {
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
	main.AddItem(a.Content, 0, 10, true)
	main.AddItem(flash, 1, 1, false)

	a.Main.AddPage("main", main, true, false)
	a.Main.AddPage("splash", ui.NewSplash("0.0.1"), true, true)
	a.toggleHeader(true)

	//Testing only
	a.inject(NewHelp(a))
	a.inject(NewEC2("ec2"))
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
	header.AddItem(tview.NewBox(), 0, 1, false)
	header.AddItem(a.Menu(), 0, 2, false)
	header.AddItem(tview.NewBox(), 0, 1, false)

	return header
}

func (a *App) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	if k, ok := a.HasAction(ui.AsKey(evt)); ok && !a.Content.IsTopDialog() {
		return k.Action(evt)
	}

	return evt
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

func (a *App) helpCmd(evt *tcell.EventKey) *tcell.EventKey {
	top := a.Content.Top()

	if top != nil && top.Name() == "help" {
		a.Content.Pop()
		return nil
	}

	if err := a.inject(NewHelp(a)); err != nil {
		a.Flash().Err(err)
	}

	return nil
}

func (a *App) inject(c model.Component) error {
	ctx := context.WithValue(context.Background(), internal.KeyApp, a)
	if err := c.Init(ctx); err != nil {
		log.Error().Err(err).Msgf("component init failed for %q", c.Name())
		//dialog.ShowError(a.Styles.Dialog(), a.Content.Pages, err.Error())
	}
	a.Content.Push(c)

	return nil
}

// PrevCmd pops the command stack.
func (a *App) PrevCmd(evt *tcell.EventKey) *tcell.EventKey {
	if !a.Content.IsLast() {
		a.Content.Pop()
	}

	return nil
}

func (a *App) statusIndicator() *ui.StatusIndicator {
	return a.Views()["statusIndicator"].(*ui.StatusIndicator)
}

func getIST(launchTime *time.Time) string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	IST := launchTime.In(loc)
	return IST.Format("Mon Jan _2 15:04:05 2006")
}

func concatJson(json1 *s3.ServerSideEncryptionConfiguration, json2 []*s3.LifecycleRule) string {
	res := BucketInfoJson{LifeCycleJson: json1, EncryptionJson: json2}
	tempRes2, _ := json.MarshalIndent(res, "", " ")
	return string(tempRes2)
}

