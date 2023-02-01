package ui

import (
	"fmt"
	"strings"

	"github.com/one2nc/cloud-lens/internal/color"
	"github.com/one2nc/cloud-lens/internal/render"
	"github.com/rs/zerolog/log"
)

const (
	descIndicator = "↓"
	ascIndicator  = "↑"
)

// TrimCell removes superfluous padding.
func TrimCell(tv *SelectTable, row, col int) string {
	c := tv.GetCell(row, col)
	if c == nil {
		log.Error().Err(fmt.Errorf("No cell at location [%d:%d]", row, col)).Msg("Trim cell failed!")
		return ""
	}
	return strings.TrimSpace(c.Text)
}

func sortIndicator(sort, asc bool, hc render.HeaderColumn) string {
	if !sort {
		return hc.Name
	}

	order := descIndicator
	if asc {
		order = ascIndicator
	}
	return fmt.Sprintf("%s%s", color.ColorizeAt(hc.Name, hc.SortIndicatorIdx, "red"), color.ColorizeAt(order, 0, "green"))
}
