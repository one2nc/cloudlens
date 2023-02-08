package view

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

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
		tcell.KeyCtrlW: ui.NewKeyAction("Toggle Wide", t.toggleWideCmd, true),
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
	row := t.GetRowCount()
	col := t.GetColumnCount()

	for i := 0; i < row; i++ {
		var row []string
		for j := 0; j < col; j++ {
			row = append(row, t.GetCell(i, j).Text)
		}
		tableData = append(tableData, row)
	}

	csvName := strings.Split(t.GetTitle(), " ")
	err := os.MkdirAll("./resource/CSV", os.ModePerm)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in creating CSV directory: %v", err))
	}
	file, err := os.Create("./resource/CSV/" + csvName[len(csvName)-2] + ".csv")
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error in creating .csv file: %v", err))
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, record := range tableData {
		err := writer.Write(record)
		if err != nil {
			log.Info().Msg(fmt.Sprintf("error in writing records to csv file: %v", err))
		}
	}
	writer.Flush()
	t.app.Flash().Info("CSV Created.")
	return nil
}
