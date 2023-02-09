package view

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
	"github.com/rs/zerolog/log"
)

const (
	DefaultRefreshRate = time.Second * 20
)

// Table represents a table viewer.
type Table struct {
	*ui.Table

	app        *App
	enterFn    EnterFunc
	bindKeysFn []BindKeysFunc
}

// NewTable returns a new viewer.
func NewTable(res string) *Table {
	t := Table{
		Table: ui.NewTable(res),
	}
	return &t
}

// Init initializes the component.
func (t *Table) Init(ctx context.Context) (err error) {
	if t.app, err = extractApp(ctx); err != nil {
		return err
	}

	t.Table.Init(ctx)
	t.SetInputCapture(t.keyboard)
	t.bindKeys()
	t.GetModel().SetRefreshRate(DefaultRefreshRate)

	return nil
}

// App returns the current app handle.
func (t *Table) App() *App {
	return t.app
}

// Start runs the component.
func (t *Table) Start() {
}

// Stop terminates the component.
func (t *Table) Stop() {
}

// SetEnterFn specifies the default enter behavior.
func (t *Table) SetEnterFn(f EnterFunc) {
	t.enterFn = f
}

func (t *Table) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	key := evt.Key()
	if key == tcell.KeyUp || key == tcell.KeyDown {
		return evt
	}

	if a, ok := t.Actions()[ui.AsKey(evt)]; ok {
		return a.Action(evt)
	}

	return evt
}

func (t *Table) bindKeys() {
	t.Actions().Add(ui.KeyActions{
		tcell.KeyCtrlW: ui.NewKeyAction("Toggle Wide", t.toggleWideCmd, false),
		ui.KeyHelp:     ui.NewKeyAction("Help", t.App().helpCmd, true),
		ui.KeyZ:        ui.NewKeyAction("CSV", t.importAsCSV, true),
	})
}

// Name returns the table name.
func (t *Table) Name() string { return t.Table.Resource() }

// AddBindKeysFn adds additional key bindings.
func (t *Table) AddBindKeysFn(f BindKeysFunc) {
	t.bindKeysFn = append(t.bindKeysFn, f)
}

func (t *Table) toggleWideCmd(evt *tcell.EventKey) *tcell.EventKey {
	t.ToggleWide()
	return nil
}
func (t *Table) importAsCSV(evt *tcell.EventKey) *tcell.EventKey {
	var tableData [][]string
	rowCount := t.GetRowCount()
	colCount := t.GetColumnCount()
	for i := 0; i < rowCount; i++ {
		var row []string
		for j := 0; j < colCount; j++ {
			text := t.GetCell(i, j).Text
			if strings.Contains(text, "↑") || strings.Contains(text, "↓") {
				text = decolorize(text)
				text = text[:len(text)-3]
			}
			row = append(row, text)
		}
		tableData = append(tableData, row)
	}
	csvFileName := strings.Split(t.GetTitle(), " ")
	usr, err := user.Current()
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in getting the machine's user: %v", err))
	}
	path := usr.HomeDir + "/cloud-lens/csv/"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in creating csv directory: %v", err))
	}
	path = filepath.Join(path + "/" + csvFileName[len(csvFileName)-2] + ".csv")
	file, err := os.Create(path)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in creating csv file: %v", err))
	}
	writer := csv.NewWriter(file)
	for _, record := range tableData {
		err := writer.Write(record)
		if err != nil {
			log.Info().Msg(fmt.Sprintf("error in writing records to csv file: %v", err))
		}
	}
	writer.Flush()
	t.app.Flash().Info("CSV file created and CSV file path copied to clipboard.")
	clipboard.WriteAll(path)
	return nil
}

func decolorize(input string) string {
	re := regexp.MustCompile("\\[.*?\\]")
	return re.ReplaceAllString(input, "")
}
