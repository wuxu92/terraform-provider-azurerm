// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // markdown

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// markdownListener is a complete listener for a parse tree produced by markdownParser.
type markdownListener interface {
	antlr.ParseTreeListener

	// EnterFile_ is called when entering the file_ production.
	EnterFile_(c *File_Context)

	// EnterHeader is called when entering the header production.
	EnterHeader(c *HeaderContext)

	// EnterList is called when entering the list production.
	EnterList(c *ListContext)

	// EnterLine is called when entering the line production.
	EnterLine(c *LineContext)

	// ExitFile_ is called when exiting the file_ production.
	ExitFile_(c *File_Context)

	// ExitHeader is called when exiting the header production.
	ExitHeader(c *HeaderContext)

	// ExitList is called when exiting the list production.
	ExitList(c *ListContext)

	// ExitLine is called when exiting the line production.
	ExitLine(c *LineContext)
}
