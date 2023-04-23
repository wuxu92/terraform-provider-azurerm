package util

import "github.com/fatih/color"

var (
	Bold       = color.New(color.Bold).Sprint
	ItalicCode = color.New(color.Italic, color.FgCyan).Sprint
	FormatCode = color.New(color.BgMagenta).Sprint
	Blue       = color.New(color.BgBlue).Sprint
	IssueLine  = color.New(color.BgYellow).Sprint
	FixedCode  = color.New(color.BgGreen).Sprint
)
