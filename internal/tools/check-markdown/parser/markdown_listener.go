// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // markdown

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// markdownListener is a complete listener for a parse tree produced by markdownParser.
type markdownListener interface {
	antlr.ParseTreeListener

	// EnterFile_ is called when entering the file_ production.
	EnterFile_(c *File_Context)

	// EnterElem is called when entering the elem production.
	EnterElem(c *ElemContext)

	// EnterHeader is called when entering the header production.
	EnterHeader(c *HeaderContext)

	// EnterPara is called when entering the para production.
	EnterPara(c *ParaContext)

	// EnterParaContent is called when entering the paraContent production.
	EnterParaContent(c *ParaContentContext)

	// EnterBold is called when entering the bold production.
	EnterBold(c *BoldContext)

	// EnterAstericks is called when entering the astericks production.
	EnterAstericks(c *AstericksContext)

	// EnterUnderscore is called when entering the underscore production.
	EnterUnderscore(c *UnderscoreContext)

	// EnterItalics is called when entering the italics production.
	EnterItalics(c *ItalicsContext)

	// EnterLink is called when entering the link production.
	EnterLink(c *LinkContext)

	// EnterQuote is called when entering the quote production.
	EnterQuote(c *QuoteContext)

	// EnterQuoteElem is called when entering the quoteElem production.
	EnterQuoteElem(c *QuoteElemContext)

	// EnterList is called when entering the list production.
	EnterList(c *ListContext)

	// EnterListElem is called when entering the listElem production.
	EnterListElem(c *ListElemContext)

	// EnterText is called when entering the text production.
	EnterText(c *TextContext)

	// EnterNl is called when entering the nl production.
	EnterNl(c *NlContext)

	// ExitFile_ is called when exiting the file_ production.
	ExitFile_(c *File_Context)

	// ExitElem is called when exiting the elem production.
	ExitElem(c *ElemContext)

	// ExitHeader is called when exiting the header production.
	ExitHeader(c *HeaderContext)

	// ExitPara is called when exiting the para production.
	ExitPara(c *ParaContext)

	// ExitParaContent is called when exiting the paraContent production.
	ExitParaContent(c *ParaContentContext)

	// ExitBold is called when exiting the bold production.
	ExitBold(c *BoldContext)

	// ExitAstericks is called when exiting the astericks production.
	ExitAstericks(c *AstericksContext)

	// ExitUnderscore is called when exiting the underscore production.
	ExitUnderscore(c *UnderscoreContext)

	// ExitItalics is called when exiting the italics production.
	ExitItalics(c *ItalicsContext)

	// ExitLink is called when exiting the link production.
	ExitLink(c *LinkContext)

	// ExitQuote is called when exiting the quote production.
	ExitQuote(c *QuoteContext)

	// ExitQuoteElem is called when exiting the quoteElem production.
	ExitQuoteElem(c *QuoteElemContext)

	// ExitList is called when exiting the list production.
	ExitList(c *ListContext)

	// ExitListElem is called when exiting the listElem production.
	ExitListElem(c *ListElemContext)

	// ExitText is called when exiting the text production.
	ExitText(c *TextContext)

	// ExitNl is called when exiting the nl production.
	ExitNl(c *NlContext)
}
