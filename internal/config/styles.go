package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	// DefaultColor represents  a default color.
	DefaultColor Color = "default"

	// TransparentColor represents the terminal bg color.
	TransparentColor Color = "-"
)

type Color string

func (c Color) String() string {
	if c.isHex() {
		return string(c)
	}
	if c == DefaultColor {
		return "-"
	}
	col := c.Color().TrueColor().Hex()
	if col < 0 {
		return "-"
	}

	return fmt.Sprintf("#%06x", col)
}

func (c Color) isHex() bool {
	return len(c) == 7 && c[0] == '#'
}

// Color returns a view color.
func (c Color) Color() tcell.Color {
	if c == DefaultColor {
		return tcell.ColorDefault
	}

	return tcell.GetColor(string(c)).TrueColor()
}

type Frame struct {
	Title  Title  `yaml:"title"`
	Border Border `yaml:"border"`
	Menu   Menu   `yaml:"menu"`
	Crumb  Crumb  `yaml:"crumbs"`
	Status Status `yaml:"status"`
}

type Title struct {
	FgColor        Color `yaml:"fgColor"`
	BgColor        Color `yaml:"bgColor"`
	HighlightColor Color `yaml:"highlightColor"`
	CounterColor   Color `yaml:"counterColor"`
	FilterColor    Color `yaml:"filterColor"`
}

// Border tracks border styles.
type Border struct {
	FgColor    Color `yaml:"fgColor"`
	FocusColor Color `yaml:"focusColor"`
}

// Crumb tracks crumbs styles.
type Crumb struct {
	FgColor     Color `yaml:"fgColor"`
	BgColor     Color `yaml:"bgColor"`
	ActiveColor Color `yaml:"activeColor"`
}
type Menu struct {
	FgColor     Color `yaml:"fgColor"`
	KeyColor    Color `yaml:"keyColor"`
	NumKeyColor Color `yaml:"numKeyColor"`
}

type Status struct {
	NewColor       Color `yaml:"newColor"`
	ModifyColor    Color `yaml:"modifyColor"`
	AddColor       Color `yaml:"addColor"`
	PendingColor   Color `yaml:"pendingColor"`
	ErrorColor     Color `yaml:"errorColor"`
	HighlightColor Color `yaml:"highlightColor"`
	KillColor      Color `yaml:"killColor"`
	CompletedColor Color `yaml:"completedColor"`
}
