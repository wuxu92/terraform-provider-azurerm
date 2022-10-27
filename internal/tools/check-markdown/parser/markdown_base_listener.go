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

// EnterElem is called when production elem is entered.
func (s *BasemarkdownListener) EnterElem(ctx *ElemContext) {}

// ExitElem is called when production elem is exited.
func (s *BasemarkdownListener) ExitElem(ctx *ElemContext) {}

// EnterHeader is called when production header is entered.
func (s *BasemarkdownListener) EnterHeader(ctx *HeaderContext) {}

// ExitHeader is called when production header is exited.
func (s *BasemarkdownListener) ExitHeader(ctx *HeaderContext) {}

// EnterPara is called when production para is entered.
func (s *BasemarkdownListener) EnterPara(ctx *ParaContext) {}

// ExitPara is called when production para is exited.
func (s *BasemarkdownListener) ExitPara(ctx *ParaContext) {}

// EnterParaContent is called when production paraContent is entered.
func (s *BasemarkdownListener) EnterParaContent(ctx *ParaContentContext) {}

// ExitParaContent is called when production paraContent is exited.
func (s *BasemarkdownListener) ExitParaContent(ctx *ParaContentContext) {}

// EnterBold is called when production bold is entered.
func (s *BasemarkdownListener) EnterBold(ctx *BoldContext) {}

// ExitBold is called when production bold is exited.
func (s *BasemarkdownListener) ExitBold(ctx *BoldContext) {}

// EnterAstericks is called when production astericks is entered.
func (s *BasemarkdownListener) EnterAstericks(ctx *AstericksContext) {}

// ExitAstericks is called when production astericks is exited.
func (s *BasemarkdownListener) ExitAstericks(ctx *AstericksContext) {}

// EnterUnderscore is called when production underscore is entered.
func (s *BasemarkdownListener) EnterUnderscore(ctx *UnderscoreContext) {}

// ExitUnderscore is called when production underscore is exited.
func (s *BasemarkdownListener) ExitUnderscore(ctx *UnderscoreContext) {}

// EnterItalics is called when production italics is entered.
func (s *BasemarkdownListener) EnterItalics(ctx *ItalicsContext) {}

// ExitItalics is called when production italics is exited.
func (s *BasemarkdownListener) ExitItalics(ctx *ItalicsContext) {}

// EnterLink is called when production link is entered.
func (s *BasemarkdownListener) EnterLink(ctx *LinkContext) {}

// ExitLink is called when production link is exited.
func (s *BasemarkdownListener) ExitLink(ctx *LinkContext) {}

// EnterQuote is called when production quote is entered.
func (s *BasemarkdownListener) EnterQuote(ctx *QuoteContext) {}

// ExitQuote is called when production quote is exited.
func (s *BasemarkdownListener) ExitQuote(ctx *QuoteContext) {}

// EnterQuoteElem is called when production quoteElem is entered.
func (s *BasemarkdownListener) EnterQuoteElem(ctx *QuoteElemContext) {}

// ExitQuoteElem is called when production quoteElem is exited.
func (s *BasemarkdownListener) ExitQuoteElem(ctx *QuoteElemContext) {}

// EnterList is called when production list is entered.
func (s *BasemarkdownListener) EnterList(ctx *ListContext) {}

// ExitList is called when production list is exited.
func (s *BasemarkdownListener) ExitList(ctx *ListContext) {}

// EnterListElem is called when production listElem is entered.
func (s *BasemarkdownListener) EnterListElem(ctx *ListElemContext) {}

// ExitListElem is called when production listElem is exited.
func (s *BasemarkdownListener) ExitListElem(ctx *ListElemContext) {}

// EnterText is called when production text is entered.
func (s *BasemarkdownListener) EnterText(ctx *TextContext) {}

// ExitText is called when production text is exited.
func (s *BasemarkdownListener) ExitText(ctx *TextContext) {}

// EnterNl is called when production nl is entered.
func (s *BasemarkdownListener) EnterNl(ctx *NlContext) {}

// ExitNl is called when production nl is exited.
func (s *BasemarkdownListener) ExitNl(ctx *NlContext) {}
