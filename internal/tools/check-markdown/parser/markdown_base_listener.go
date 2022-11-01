// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // markdown

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BasemarkdownListener is a complete listener for a parse tree produced by markdownParser.
type BasemarkdownListener struct{}

var _ markdownListener = &BasemarkdownListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasemarkdownListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasemarkdownListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasemarkdownListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasemarkdownListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterFile_ is called when production file_ is entered.
func (s *BasemarkdownListener) EnterFile_(ctx *File_Context) {}

// ExitFile_ is called when production file_ is exited.
func (s *BasemarkdownListener) ExitFile_(ctx *File_Context) {}

// EnterHeader is called when production header is entered.
func (s *BasemarkdownListener) EnterHeader(ctx *HeaderContext) {}

// ExitHeader is called when production header is exited.
func (s *BasemarkdownListener) ExitHeader(ctx *HeaderContext) {}

// EnterList is called when production list is entered.
func (s *BasemarkdownListener) EnterList(ctx *ListContext) {}

// ExitList is called when production list is exited.
func (s *BasemarkdownListener) ExitList(ctx *ListContext) {}

// EnterLine is called when production line is entered.
func (s *BasemarkdownListener) EnterLine(ctx *LineContext) {}

// ExitLine is called when production line is exited.
func (s *BasemarkdownListener) ExitLine(ctx *LineContext) {}
