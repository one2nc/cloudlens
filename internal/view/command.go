package view

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/one2nc/cloud-lens/internal/model"
	"github.com/rs/zerolog/log"
)

var (
	customViewers MetaViewers

	canRX = regexp.MustCompile(`\Acan\s([u|g|s]):([\w-:]+)\b`)
)

// Command represents a user command.
type Command struct {
	app *App

	mx sync.Mutex
}

// NewCommand returns a new command.
func NewCommand(app *App) *Command {
	return &Command{
		app: app,
	}
}

// Init initializes the command.
func (c *Command) Init() error {
	customViewers = loadCustomViewers()
	return nil
}

// Reset resets Command and reload aliases.
func (c *Command) Reset(clear bool) error {

	return nil
}

// Exec the Command by showing associated display.
func (c *Command) run(cmd, path string, clearStack bool) error {
	if c.specialCmd(cmd, path) {
		return nil
	}
	cmds := strings.Split(cmd, " ")
	res, v, err := c.viewMetaFor(cmds[0])
	if err != nil {
		return err
	}

	switch cmds[0] {
	default:
		return c.exec(cmd, res, c.componentFor(res, path, v), clearStack)
	}
}

func (c *Command) defaultCmd() error {
	return c.run("ec2", "", true)
}

func (c *Command) specialCmd(cmd, path string) bool {
	cmds := strings.Split(cmd, " ")
	switch cmds[0] {
	// case "cow":
	// 	c.app.cowCmd(path)
	// 	return true
	case "q", "q!", "Q", "quit":
		c.app.BailOut()
		return true
	case "?", "h", "help":
		c.app.helpCmd(nil)
		return true
		// case "a", "alias":
		// 	c.app.aliasCmd(nil)
		return true
	default:
		if !canRX.MatchString(cmd) {
			return false
		}
	}
	return false
}

func (c *Command) viewMetaFor(cmd string) (string, *MetaViewer, error) {
	// gvr, ok := c.alias.AsGVR(cmd)
	// if !ok {
	// 	return "", nil, fmt.Errorf("`%s` command not found", cmd)
	// }

	v, ok := customViewers[cmd]
	log.Info().Msg(fmt.Sprintf("Is ok: %v", ok))
	if !ok {
		return cmd, &MetaViewer{viewerFn: NewBrowser}, nil
	}

	return cmd, &v, nil
}

func (c *Command) componentFor(res, path string, v *MetaViewer) ResourceViewer {
	var view ResourceViewer
	if v.viewerFn != nil {
		log.Info().Msg(fmt.Sprintf("If res: %v, Path: %v", res, path))
		view = v.viewerFn(res)
	} else {
		log.Info().Msg(fmt.Sprintf("else: res: %v, Path: %v", res, path))
		view = NewBrowser(res)
	}

	if v.enterFn != nil {
		view.GetTable().SetEnterFn(v.enterFn)
	}

	return view
}

func (c *Command) exec(cmd, gvr string, comp model.Component, clearStack bool) (err error) {
	defer func() {
		// if e := recover(); e != nil {
		// 	log.Error().Msgf("Something bad happened! %#v", e)
		// 	c.app.Content.Dump()
		// 	//log.Debug().Msgf("History %v", c.app.cmdHistory.List())

		// 	//hh := c.app.cmdHistory.List()
		// 	hh := []string{}
		// 	if len(hh) == 0 {
		// 		_ = c.run("pod", "", true)
		// 	} else {
		// 		_ = c.run(hh[0], "", true)
		// 	}
		// 	err = fmt.Errorf("Invalid command %q", cmd)
		// }
	}()
	log.Info().Msg(fmt.Sprintf("cmd: %v, res: %v, comp: %T", cmd, gvr, comp))

	if comp == nil {
		return fmt.Errorf("No component found for %s", gvr)
	}
	c.app.Flash().Infof("Viewing %s...", cmd)

	if clearStack {
		c.app.Content.Stack.Clear()
	}

	if err := c.app.inject(comp); err != nil {
		return err
	}
	//c.app.cmdHistory.Push(cmd)

	return
}
