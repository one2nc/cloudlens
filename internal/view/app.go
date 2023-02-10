package view

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/one2nc/cloud-lens/internal/model"
	"github.com/one2nc/cloud-lens/internal/ui"
	"github.com/one2nc/cloud-lens/internal/ui/dialog"
	"github.com/rs/zerolog/log"
)

const (
	splashDelay = 1 * time.Second
)

type App struct {
	*ui.App
	Content             *PageStack
	command             *Command
	context             context.Context
	cancelFn            context.CancelFunc
	showHeader          bool
	IsPageContentSorted bool
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

func (a *App) Init(profiles, regions []string, ctx context.Context) error {
	ctx = context.WithValue(ctx, internal.KeyActiveProfile, profiles[0])
	ctx = context.WithValue(ctx, internal.KeyActiveRegion, regions[0])
	ctx = context.WithValue(ctx, internal.KeyApp, a)
	a.SetContext(ctx)

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

	if err := a.Content.Init(ctx); err != nil {
		return err
	}
	a.Content.Stack.AddListener(a.Menu())
	a.Content.Stack.AddListener(a.Crumbs())
	a.App.Init()
	a.SetInputCapture(a.keyboard)
	a.bindKeys()

	a.command = NewCommand(a)
	if err := a.command.Init(); err != nil {
		return err
	}
	a.CmdBuff().SetSuggestionFn(a.suggestCommand())
	a.layout(ctx)
	return nil
}

func (a *App) layout(ctx context.Context) {
	flash := ui.NewFlash(a.App)
	go flash.Watch(ctx, a.Flash().Channel())
	main := tview.NewFlex().SetDirection(tview.FlexRow)

	main.AddItem(a.statusIndicator(), 1, 1, false)
	main.AddItem(a.Content, 0, 10, true)
	main.AddItem(a.Crumbs(), 1, 1, false)
	main.AddItem(flash, 1, 1, false)

	a.Main.AddPage("main", main, true, false)
	a.Main.AddPage("splash", ui.NewSplash("0.0.1"), true, true)
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
		})
	}()

	if err := a.command.defaultCmd(); err != nil {
		return err
	}
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

func (a *App) regionChanged(region string, index int) {
	profile := a.GetContext().Value(internal.KeyActiveProfile).(string)
	a.refreshSession(profile, region)
}

func (a *App) refreshSession(profile string, region string) {
	sess, err := config.GetSession(profile, region)
	if err != nil {
		a.App.Flash().Err(err)
		return
	}
	ctx := context.WithValue(a.GetContext(), internal.KeySession, sess)
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

func (a *App) region() *ui.DropDown {
	return a.Views()["region"].(*ui.DropDown)
}
