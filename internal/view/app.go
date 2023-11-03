package view

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	cfg "github.com/aws/aws-sdk-go-v2/config"
	awsS "github.com/aws/aws-sdk-go/aws"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/one2nc/cloudlens/internal/model"
	"github.com/one2nc/cloudlens/internal/ui"
	"github.com/one2nc/cloudlens/internal/ui/dialog"
	"github.com/rs/zerolog/log"
)

const (
	splashDelay = 1 * time.Second
)

var (
	availableCloud  = []string{"AWS", "GCP"}
	profile, region string
)

type App struct {
	*ui.App
	Content             *PageStack
	command             *Command
	context             context.Context
	cancelFn            context.CancelFunc
	showHeader          bool
	IsPageContentSorted bool
	version             string
	cloudConfig         config.CloudConfig
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

// TODO keep context param at first place always
func (a *App) Init(version string, cloudConfig config.CloudConfig) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, internal.KeyApp, a)
	a.SetContext(ctx)

	a.version = model.NormalizeVersion(version)
	if err := a.Content.Init(ctx); err != nil {
		return err
	}
	a.Content.Stack.AddListener(a.Menu())
	a.Content.Stack.AddListener(a.Crumbs())

	a.App.Init()
	a.SetInputCapture(a.keyboard)
	a.bindKeys()

	a.CmdBuff().SetSuggestionFn(a.suggestCommand())

	a.layout(ctx)

	if cloudConfig.SelectedCloud != "" {
		a.cloudConfig = cloudConfig
		err := a.handleCloudSelection(cloudConfig.SelectedCloud)
		if err != nil {
			return err
		}
	} else {
		a.showCloudSelectionScreen()
	}
	return nil
}

func (a *App) handleAWS() {

	region = a.cloudConfig.Region
	profile = a.cloudConfig.Profile
	awsConfigInput := aws.AWSConfigInput{
		UseLocalStack: a.cloudConfig.UseLocalStack,
	}
	var regions, profiles []string
	if a.cloudConfig.UseLocalStack {
		profile = "localstack"
		profiles = []string{profile}
		regions = readAndValidateRegion()
	} else {
		var err error
		profiles, err = readAndValidateProfile()
		if err != nil {
			panic(err)
		}
		if len(profiles) > 0 {
			if profiles[0] == "default" && len(region) == 0 {
				region = getDefaultAWSRegion()
			} else if len(region) == 0 {
				region = "ap-south-1"
			}

			regions = readAndValidateRegion()

		} else {
			profile := awsS.String(os.Getenv(internal.AWS_PROFILE))
			profiles = []string{*profile}
			region := awsS.String(os.Getenv(internal.AWS_DEFAULT_REGION))
			regions = []string{*region}
			awsConfigInput.UseEnvVariables = true
		}
	}

	awsConfigInput.Profile = profiles[0]
	awsConfigInput.Region = regions[0]
	cfg, err := aws.GetCfg(awsConfigInput)
	if err != nil {
		panic(fmt.Sprintf("aws session init failed -- %v", err))
	}
	ctx := context.WithValue(a.context, internal.KeySession, cfg)
	a.SetContext(ctx)

	ctx = context.WithValue(ctx, internal.KeyActiveProfile, profiles[0])
	ctx = context.WithValue(ctx, internal.KeyActiveRegion, regions[0])
	ctx = context.WithValue(ctx, internal.KeySelectedCloud, internal.AWS)
	a.SetContext(ctx)
	a.App.UpdateContext(ctx)

	p := ui.NewDropDown("Profile:", profiles)
	p.SetSelectedFunc(a.profileChanged)
	a.Views()["profile"] = p

	r := ui.NewDropDown("Region:", regions)
	r.SetSelectedFunc(a.regionChanged)
	a.Views()["region"] = r

	infoData := map[string]tview.Primitive{
		"profile": a.profile(),
		"region":  a.region(),
	}
	a.Views()["info"] = ui.NewInfo(infoData)
	a.toggleHeader(true)

}

func (a *App) handleGCP() error {
	ctx := a.GetContext()
	ctx = context.WithValue(ctx, internal.KeySelectedCloud, internal.GCP)
	credFilePath := os.Getenv(internal.GOOGLE_APPLICATION_CREDENTIALS)
	if credFilePath == "" {
		go func() {
			<-time.After(splashDelay)
			dialog.ShowError(a.Content.Pages, "Invalid path to google credentials")

		}()
	}
	serviceAccount, err := gcp.FetchProjectID(credFilePath)
	if err != nil {
		go func() {
			<-time.After(splashDelay)
			dialog.ShowError(a.Content.Pages, "Invalid path to google credentials")

		}()
	}

	ctx = context.WithValue(ctx, internal.KeyActiveProject, serviceAccount.ProjectID)

	p := ui.NewDropDown("Projects:", []string{serviceAccount.ProjectID})
	p.SetSelectedFunc(a.projectchanged)
	a.Views()["project"] = p

	infoData := map[string]tview.Primitive{
		"project": a.project(),
	}
	a.Views()["info"] = ui.NewInfo(infoData)
	a.SetContext(ctx)
	a.App.UpdateContext(ctx)
	a.toggleHeader(true)
	return nil
}

func (a *App) showCloudSelectionScreen() {
	cloudSelectScreen := ui.NewCloudSelectionScreen(ui.OptionWithAction{
		"AWS": func() {
			a.handleCloudSelection(internal.AWS)
		},
		"GCP": func() {
			a.handleCloudSelection(internal.GCP)
		},
	}, a.version)
	a.App.SetFocus(cloudSelectScreen.GetFocusItem())
	a.Main.AddPage(internal.MAIN_SCREEN, cloudSelectScreen, true, true)

}

func (a *App) handleCloudSelection(seletedCloud string) error {
	switch seletedCloud {
	case internal.AWS:
		a.handleAWS()
		a.Main.SwitchToPage(internal.AWS_SCREEN)
	case internal.GCP:
		err := a.handleGCP()
		if err != nil {
			return err
		}
		a.Main.SwitchToPage(internal.GCP_SCREEN)
	}
	a.command = NewCommand(a)
	if err := a.command.Init(); err != nil {
		log.Print(err)
		return err
	}
	if err := a.command.defaultCmd(); err != nil {
		return err
	}
	return nil
}

func (a *App) layout(ctx context.Context) {
	flash := ui.NewFlash(a.App)
	go flash.Watch(ctx, a.Flash().Channel())

	aws := tview.NewFlex().SetDirection(tview.FlexRow)
	aws.AddItem(a.statusIndicator(), 1, 1, false)
	aws.AddItem(a.Content, 0, 10, true)
	aws.AddItem(a.Crumbs(), 1, 1, false)
	aws.AddItem(flash, 1, 1, false)
	a.Main.AddPage(internal.AWS_SCREEN, aws, true, false)

	gcp := tview.NewFlex().SetDirection(tview.FlexRow)
	gcp.AddItem(a.statusIndicator(), 1, 1, false)
	gcp.AddItem(a.Content, 0, 10, true)
	gcp.AddItem(a.Crumbs(), 1, 1, false)
	gcp.AddItem(flash, 1, 1, false)
	a.Main.AddPage(internal.GCP_SCREEN, gcp, true, false)

	a.Main.AddPage(internal.SPLASH_SCREEN, ui.NewSplash("0.1.3"), true, true)
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

			switch a.cloudConfig.SelectedCloud {
			case "":
				a.Main.SwitchToPage(internal.MAIN_SCREEN)

			case internal.AWS:
				a.Main.SwitchToPage(internal.AWS_SCREEN)

			case internal.GCP:
				a.Main.SwitchToPage(internal.GCP_SCREEN)
			}
		})
	}()
	a.SetRunning(true)

	if err := a.Application.Run(); err != nil {
		return err
	}

	return nil
}

func (a *App) GetContext() context.Context {
	return a.context
}

func (a *App) SetContext(ctx context.Context) {
	a.context = ctx
}

func (a *App) toggleHeader(header bool) {
	a.showHeader = header
	cloud := a.context.Value(internal.KeySelectedCloud).(string)

	flex, ok := a.Main.GetPrimitive(cloud).(*tview.Flex)
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
	header.AddItem(a.info(), 50, 1, false)
	header.AddItem(a.Menu(), 0, 1, false)
	header.AddItem(ui.NewLogo(), 26, 1, false)
	return header
}

func (a *App) suggestCommand() model.SuggestionFunc {
	return func(s string) (entries sort.StringSlice) {
		// if s == "" {
		// 	if a.cmdHistory.Empty() {
		// 		return
		// 	}
		// 	return a.cmdHistory.List()
		// }

		s = strings.ToLower(s)
		for _, k := range a.command.alias.Keys() {
			if k == s {
				continue
			}
			if strings.HasPrefix(k, s) {
				entries = append(entries, strings.Replace(k, s, "", 1))
			}
		}
		if len(entries) == 0 {
			return nil
		}
		entries.Sort()
		return
	}
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
		tcell.KeyEnter: ui.NewKeyAction("Goto", a.gotoCmd, false),
		tcell.KeyTAB:   ui.NewKeyAction("switch", NewTab(a).tabAction, false),
	})
}

func (a *App) toggleHeaderCmd(evt *tcell.EventKey) *tcell.EventKey {

	a.QueueUpdateDraw(func() {
		a.showHeader = !a.showHeader
		a.toggleHeader(a.showHeader)
	})

	return nil
}

func (a *App) gotoCmd(evt *tcell.EventKey) *tcell.EventKey {
	if a.CmdBuff().IsActive() && !a.CmdBuff().Empty() {
		a.gotoResource(a.GetCmd(), "", true)
		a.ResetCmd()
		return nil
	}

	return evt
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

func (a *App) profileChanged(profile string, index int) {
	region := a.GetContext().Value(internal.KeyActiveRegion).(string)
	a.refreshSession(profile, region)
}

func (a *App) projectchanged(project string, index int) {
	a.refreshProject(project)
}

func (a *App) regionChanged(region string, index int) {
	profile := a.GetContext().Value(internal.KeyActiveProfile).(string)
	a.refreshSession(profile, region)
}

func (a *App) refreshSession(profile string, region string) {

	awsConfigInput := aws.AWSConfigInput{
		UseLocalStack: a.cloudConfig.UseLocalStack,
		Profile:       profile,
		Region:        region,
	}
	cfg, err := aws.GetCfg(awsConfigInput)
	//sess, err := aws.GetSession(profile, region)
	if err != nil {
		a.App.Flash().Err(err)
		return
	}
	ctx := context.WithValue(a.GetContext(), internal.KeySession, cfg)
	a.SetContext(ctx)
	stackedViews := a.Content.Pages.Stack.Flatten()
	a.gotoResource(stackedViews[0], "", true)
	a.App.Flash().Infof("Refreshing %v...", stackedViews[0])
}

func (a *App) refreshProject(project string) {
	ctx := context.WithValue(a.GetContext(), internal.KeyActiveProject, project)
	a.SetContext(ctx)
	stackedViews := a.Content.Pages.Stack.Flatten()
	a.gotoResource(stackedViews[0], "", true)
	a.App.Flash().Infof("Refreshing %v...", stackedViews[0])
}

func (a *App) gotoResource(cmd, path string, clearStack bool) {
	err := a.command.run(cmd, path, clearStack)
	if err != nil {
		dialog.ShowError(a.Content.Pages, err.Error())
	}
}

func (a *App) inject(c model.Component) error {
	if err := c.Init(a.context); err != nil {
		log.Error().Err(err).Msgf("component init failed for %q", c.Name())
		dialog.ShowError(a.Content.Pages, err.Error())
	}
	a.Content.Push(c)
	return nil
}

// PrevCmd pops the command stack.
func (a *App) PrevCmd(evt *tcell.EventKey) *tcell.EventKey {
	if !a.Content.IsLast() {
		a.Content.Pop()
		fn := fmt.Sprintf("%v", a.context.Value(internal.FolderName))
		newFn := ""
		if strings.Count(fn, "/") > 0 {
			fl := strings.Split(fn, "/")
			for i := 0; i < len(fl)-2; i++ {
				newFn = newFn + fl[i] + "/"
			}
			ctx := context.WithValue(a.GetContext(), internal.FolderName, newFn)
			a.SetContext(ctx)
			log.Info().Msg(fmt.Sprintf("inside prv cmd: %v", a.context.Value(internal.FolderName)))
		}
	} else {

		if a.cloudConfig.SelectedCloud == "" {
			a.Main.SwitchToPage(internal.MAIN_SCREEN)
		}
	}

	return nil
}

func (a *App) statusIndicator() *ui.StatusIndicator {
	return a.Views()["statusIndicator"].(*ui.StatusIndicator)
}

func (a *App) info() *ui.Info {
	return a.Views()["info"].(*ui.Info)
}

func (a *App) profile() *ui.DropDown {
	return a.Views()["profile"].(*ui.DropDown)
}

func (a *App) project() *ui.DropDown {
	return a.Views()["project"].(*ui.DropDown)
}

func (a *App) region() *ui.DropDown {
	return a.Views()["region"].(*ui.DropDown)
}

func readAndValidateProfile() ([]string, error) {
	profiles, err := aws.GetProfiles()
	if err != nil {
		log.Printf("failed to read profiles -- %v", err)
		return nil, err
	}
	profiles, isSwapped := config.SwapFirstIndexWithValue(profiles, profile)
	if !isSwapped {
		if profile != "" {
			fmt.Fprintf(os.Stderr, "Could not load profile: %v\n", profile)
			os.Exit(1)
		}
	}
	return profiles, nil
}

func readAndValidateRegion() []string {
	regions := aws.GetAllRegions()
	regions, isSwapped := config.SwapFirstIndexWithValue(regions, region)
	if !isSwapped {
		if region != "" {
			fmt.Fprintf(os.Stderr, "Could not load region: %v\n", region)
			os.Exit(1)
		}
	}
	return regions
}

func getDefaultAWSRegion() string {
	cfg, err := cfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load AWS SDK config: %v\n", err)
		os.Exit(1)
	}
	region := cfg.Region
	return region
}
