package parser

import (
	"fmt"
	"strings"
)

const (
	PosDefault = iota
	PosExample
	PosArgument
	PosAttribute
	PosTimeout
	PosImport
)

type MD struct {
	BasemarkdownListener

	Args     []string // get all args
	Attrs    []string // get all attrs
	Timeouts []string
	curPos   int
}

// EnterHeader is called when production header is entered.
func (m *MD) EnterHeader(ctx *HeaderContext) {
	text := ctx.GetText()
	if strings.Contains(text, "Argument") {
		m.curPos = PosArgument
	} else if strings.Contains(text, "Attribute") {
		m.curPos = PosAttribute
	} else if strings.Contains(text, "Example") {
		m.curPos = PosExample
	} else if strings.Contains(text, "Import") {
		m.curPos = PosImport
	} else if strings.Contains(text, "Timeout") {
		m.curPos = PosTimeout
	} else {
		m.curPos = PosDefault
	}
}

// ExitHeader is called when production header is exited.
func (m *MD) ExitHeader(ctx *HeaderContext) {}

// EnterListElem is called when production listElem is entered.
func (m *MD) EnterList(ctx *ListContext) {
	txt := ctx.GetText()
	switch m.curPos {
	case PosArgument:
		m.Args = append(m.Args, txt)
	case PosAttribute:
		m.Attrs = append(m.Attrs, txt)
	case PosTimeout:
		m.Timeouts = append(m.Timeouts, txt)
	default:
	}
}

func (m *MD) EnterLine(ctx *LineContext) {
	txt := ctx.GetText()
	switch m.curPos {
	case PosArgument:
		m.Args = append(m.Args, txt)
	case PosAttribute:
		m.Attrs = append(m.Attrs, txt)
	case PosTimeout:
		m.Timeouts = append(m.Timeouts, txt)
	default:
	}
}

// ExitListElem is called when production listElem is exited.
func (m *MD) ExitList(ctx *ListContext) {
	txt := ctx.GetText()
	fmt.Println(txt)
}

// EnterCode is called when production code is entered.
//func (m *MD) EnterCode(ctx *CodeContext) {}
//
//// ExitCode is called when production code is exited.
//func (m *MD) ExitCode(ctx *CodeContext) {}
