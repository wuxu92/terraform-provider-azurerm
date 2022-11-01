// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // markdown

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type markdownParser struct {
	*antlr.BaseParser
}

var markdownParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func markdownParserInit() {
	staticData := &markdownParserStaticData
	staticData.literalNames = []string{
		"", "'-'", "'*'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "LIST", "HEAD", "LINE", "WS", "LN", "END",
	}
	staticData.ruleNames = []string{
		"file_", "header", "list", "line",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 8, 26, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 1, 0, 1, 0,
		1, 0, 4, 0, 12, 8, 0, 11, 0, 12, 0, 13, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1,
		2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 0, 0, 4, 0, 2, 4, 6, 0, 1, 1, 0, 1, 2,
		24, 0, 11, 1, 0, 0, 0, 2, 17, 1, 0, 0, 0, 4, 20, 1, 0, 0, 0, 6, 23, 1,
		0, 0, 0, 8, 12, 3, 2, 1, 0, 9, 12, 3, 4, 2, 0, 10, 12, 3, 6, 3, 0, 11,
		8, 1, 0, 0, 0, 11, 9, 1, 0, 0, 0, 11, 10, 1, 0, 0, 0, 12, 13, 1, 0, 0,
		0, 13, 11, 1, 0, 0, 0, 13, 14, 1, 0, 0, 0, 14, 15, 1, 0, 0, 0, 15, 16,
		5, 0, 0, 1, 16, 1, 1, 0, 0, 0, 17, 18, 5, 4, 0, 0, 18, 19, 5, 5, 0, 0,
		19, 3, 1, 0, 0, 0, 20, 21, 7, 0, 0, 0, 21, 22, 5, 5, 0, 0, 22, 5, 1, 0,
		0, 0, 23, 24, 5, 5, 0, 0, 24, 7, 1, 0, 0, 0, 2, 11, 13,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// markdownParserInit initializes any static state used to implement markdownParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewmarkdownParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func MarkdownParserInit() {
	staticData := &markdownParserStaticData
	staticData.once.Do(markdownParserInit)
}

// NewmarkdownParser produces a new parser instance for the optional input antlr.TokenStream.
func NewmarkdownParser(input antlr.TokenStream) *markdownParser {
	MarkdownParserInit()
	this := new(markdownParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &markdownParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "java-escape"

	return this
}

// markdownParser tokens.
const (
	markdownParserEOF  = antlr.TokenEOF
	markdownParserT__0 = 1
	markdownParserT__1 = 2
	markdownParserLIST = 3
	markdownParserHEAD = 4
	markdownParserLINE = 5
	markdownParserWS   = 6
	markdownParserLN   = 7
	markdownParserEND  = 8
)

// markdownParser rules.
const (
	markdownParserRULE_file_  = 0
	markdownParserRULE_header = 1
	markdownParserRULE_list   = 2
	markdownParserRULE_line   = 3
)

// IFile_Context is an interface to support dynamic dispatch.
type IFile_Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFile_Context differentiates from other interfaces.
	IsFile_Context()
}

type File_Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFile_Context() *File_Context {
	var p = new(File_Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_file_
	return p
}

func (*File_Context) IsFile_Context() {}

func NewFile_Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *File_Context {
	var p = new(File_Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_file_

	return p
}

func (s *File_Context) GetParser() antlr.Parser { return s.parser }

func (s *File_Context) EOF() antlr.TerminalNode {
	return s.GetToken(markdownParserEOF, 0)
}

func (s *File_Context) AllHeader() []IHeaderContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHeaderContext); ok {
			len++
		}
	}

	tst := make([]IHeaderContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHeaderContext); ok {
			tst[i] = t.(IHeaderContext)
			i++
		}
	}

	return tst
}

func (s *File_Context) Header(i int) IHeaderContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHeaderContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHeaderContext)
}

func (s *File_Context) AllList() []IListContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IListContext); ok {
			len++
		}
	}

	tst := make([]IListContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IListContext); ok {
			tst[i] = t.(IListContext)
			i++
		}
	}

	return tst
}

func (s *File_Context) List(i int) IListContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListContext)
}

func (s *File_Context) AllLine() []ILineContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILineContext); ok {
			len++
		}
	}

	tst := make([]ILineContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILineContext); ok {
			tst[i] = t.(ILineContext)
			i++
		}
	}

	return tst
}

func (s *File_Context) Line(i int) ILineContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILineContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILineContext)
}

func (s *File_Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *File_Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *File_Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterFile_(s)
	}
}

func (s *File_Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitFile_(s)
	}
}

func (p *markdownParser) File_() (localctx IFile_Context) {
	this := p
	_ = this

	localctx = NewFile_Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, markdownParserRULE_file_)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(11)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&54) != 0 {
		p.SetState(11)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case markdownParserHEAD:
			{
				p.SetState(8)
				p.Header()
			}

		case markdownParserT__0, markdownParserT__1:
			{
				p.SetState(9)
				p.List()
			}

		case markdownParserLINE:
			{
				p.SetState(10)
				p.Line()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(13)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(15)
		p.Match(markdownParserEOF)
	}

	return localctx
}

// IHeaderContext is an interface to support dynamic dispatch.
type IHeaderContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHeaderContext differentiates from other interfaces.
	IsHeaderContext()
}

type HeaderContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHeaderContext() *HeaderContext {
	var p = new(HeaderContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_header
	return p
}

func (*HeaderContext) IsHeaderContext() {}

func NewHeaderContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HeaderContext {
	var p = new(HeaderContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_header

	return p
}

func (s *HeaderContext) GetParser() antlr.Parser { return s.parser }

func (s *HeaderContext) HEAD() antlr.TerminalNode {
	return s.GetToken(markdownParserHEAD, 0)
}

func (s *HeaderContext) LINE() antlr.TerminalNode {
	return s.GetToken(markdownParserLINE, 0)
}

func (s *HeaderContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HeaderContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HeaderContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterHeader(s)
	}
}

func (s *HeaderContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitHeader(s)
	}
}

func (p *markdownParser) Header() (localctx IHeaderContext) {
	this := p
	_ = this

	localctx = NewHeaderContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, markdownParserRULE_header)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(17)
		p.Match(markdownParserHEAD)
	}
	{
		p.SetState(18)
		p.Match(markdownParserLINE)
	}

	return localctx
}

// IListContext is an interface to support dynamic dispatch.
type IListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsListContext differentiates from other interfaces.
	IsListContext()
}

type ListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyListContext() *ListContext {
	var p = new(ListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_list
	return p
}

func (*ListContext) IsListContext() {}

func NewListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListContext {
	var p = new(ListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_list

	return p
}

func (s *ListContext) GetParser() antlr.Parser { return s.parser }

func (s *ListContext) LINE() antlr.TerminalNode {
	return s.GetToken(markdownParserLINE, 0)
}

func (s *ListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterList(s)
	}
}

func (s *ListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitList(s)
	}
}

func (p *markdownParser) List() (localctx IListContext) {
	this := p
	_ = this

	localctx = NewListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, markdownParserRULE_list)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(20)
		_la = p.GetTokenStream().LA(1)

		if !(_la == markdownParserT__0 || _la == markdownParserT__1) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(21)
		p.Match(markdownParserLINE)
	}

	return localctx
}

// ILineContext is an interface to support dynamic dispatch.
type ILineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLineContext differentiates from other interfaces.
	IsLineContext()
}

type LineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLineContext() *LineContext {
	var p = new(LineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_line
	return p
}

func (*LineContext) IsLineContext() {}

func NewLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineContext {
	var p = new(LineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_line

	return p
}

func (s *LineContext) GetParser() antlr.Parser { return s.parser }

func (s *LineContext) LINE() antlr.TerminalNode {
	return s.GetToken(markdownParserLINE, 0)
}

func (s *LineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterLine(s)
	}
}

func (s *LineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitLine(s)
	}
}

func (p *markdownParser) Line() (localctx ILineContext) {
	this := p
	_ = this

	localctx = NewLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, markdownParserRULE_line)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(23)
		p.Match(markdownParserLINE)
	}

	return localctx
}
