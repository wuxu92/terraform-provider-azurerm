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
		"", "'\\n'", "'#'", "'*'", "' '", "'_'", "'['", "']'", "'('", "')'",
		"'>'", "'`'", "'\\r'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "WS",
	}
	staticData.ruleNames = []string{
		"file_", "elem", "header", "para", "paraContent", "bold", "astericks",
		"underscore", "italics", "link", "quote", "quoteElem", "list", "listElem",
		"text", "nl",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 13, 158, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		1, 0, 4, 0, 34, 8, 0, 11, 0, 12, 0, 35, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3,
		1, 43, 8, 1, 1, 2, 4, 2, 46, 8, 2, 11, 2, 12, 2, 47, 1, 2, 5, 2, 51, 8,
		2, 10, 2, 12, 2, 54, 9, 2, 1, 2, 1, 2, 1, 3, 5, 3, 59, 8, 3, 10, 3, 12,
		3, 62, 9, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 68, 8, 3, 1, 4, 1, 4, 1, 4,
		1, 4, 1, 4, 1, 4, 1, 4, 4, 4, 77, 8, 4, 11, 4, 12, 4, 78, 1, 5, 1, 5, 1,
		5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1,
		8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 5, 9, 104, 8, 9, 10,
		9, 12, 9, 107, 9, 9, 1, 9, 1, 9, 1, 10, 4, 10, 112, 8, 10, 11, 10, 12,
		10, 113, 1, 10, 1, 10, 1, 11, 1, 11, 5, 11, 120, 8, 11, 10, 11, 12, 11,
		123, 9, 11, 1, 11, 1, 11, 1, 12, 4, 12, 128, 8, 12, 11, 12, 12, 12, 129,
		1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 3, 13, 138, 8, 13, 3, 13, 140,
		8, 13, 3, 13, 142, 8, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 14, 4, 14, 149,
		8, 14, 11, 14, 12, 14, 150, 1, 15, 3, 15, 154, 8, 15, 1, 15, 1, 15, 1,
		15, 0, 0, 16, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
		0, 4, 1, 0, 1, 1, 2, 0, 1, 1, 4, 4, 1, 0, 9, 9, 3, 0, 1, 3, 5, 7, 10, 11,
		166, 0, 33, 1, 0, 0, 0, 2, 42, 1, 0, 0, 0, 4, 45, 1, 0, 0, 0, 6, 60, 1,
		0, 0, 0, 8, 76, 1, 0, 0, 0, 10, 80, 1, 0, 0, 0, 12, 85, 1, 0, 0, 0, 14,
		89, 1, 0, 0, 0, 16, 93, 1, 0, 0, 0, 18, 98, 1, 0, 0, 0, 20, 111, 1, 0,
		0, 0, 22, 117, 1, 0, 0, 0, 24, 127, 1, 0, 0, 0, 26, 141, 1, 0, 0, 0, 28,
		148, 1, 0, 0, 0, 30, 153, 1, 0, 0, 0, 32, 34, 3, 2, 1, 0, 33, 32, 1, 0,
		0, 0, 34, 35, 1, 0, 0, 0, 35, 33, 1, 0, 0, 0, 35, 36, 1, 0, 0, 0, 36, 1,
		1, 0, 0, 0, 37, 43, 3, 4, 2, 0, 38, 43, 3, 6, 3, 0, 39, 43, 3, 20, 10,
		0, 40, 43, 3, 24, 12, 0, 41, 43, 5, 1, 0, 0, 42, 37, 1, 0, 0, 0, 42, 38,
		1, 0, 0, 0, 42, 39, 1, 0, 0, 0, 42, 40, 1, 0, 0, 0, 42, 41, 1, 0, 0, 0,
		43, 3, 1, 0, 0, 0, 44, 46, 5, 2, 0, 0, 45, 44, 1, 0, 0, 0, 46, 47, 1, 0,
		0, 0, 47, 45, 1, 0, 0, 0, 47, 48, 1, 0, 0, 0, 48, 52, 1, 0, 0, 0, 49, 51,
		8, 0, 0, 0, 50, 49, 1, 0, 0, 0, 51, 54, 1, 0, 0, 0, 52, 50, 1, 0, 0, 0,
		52, 53, 1, 0, 0, 0, 53, 55, 1, 0, 0, 0, 54, 52, 1, 0, 0, 0, 55, 56, 5,
		1, 0, 0, 56, 5, 1, 0, 0, 0, 57, 59, 5, 1, 0, 0, 58, 57, 1, 0, 0, 0, 59,
		62, 1, 0, 0, 0, 60, 58, 1, 0, 0, 0, 60, 61, 1, 0, 0, 0, 61, 63, 1, 0, 0,
		0, 62, 60, 1, 0, 0, 0, 63, 64, 3, 8, 4, 0, 64, 67, 5, 1, 0, 0, 65, 68,
		3, 30, 15, 0, 66, 68, 5, 0, 0, 1, 67, 65, 1, 0, 0, 0, 67, 66, 1, 0, 0,
		0, 68, 7, 1, 0, 0, 0, 69, 77, 3, 28, 14, 0, 70, 77, 3, 10, 5, 0, 71, 77,
		3, 16, 8, 0, 72, 77, 3, 18, 9, 0, 73, 77, 3, 12, 6, 0, 74, 77, 3, 14, 7,
		0, 75, 77, 5, 1, 0, 0, 76, 69, 1, 0, 0, 0, 76, 70, 1, 0, 0, 0, 76, 71,
		1, 0, 0, 0, 76, 72, 1, 0, 0, 0, 76, 73, 1, 0, 0, 0, 76, 74, 1, 0, 0, 0,
		76, 75, 1, 0, 0, 0, 77, 78, 1, 0, 0, 0, 78, 76, 1, 0, 0, 0, 78, 79, 1,
		0, 0, 0, 79, 9, 1, 0, 0, 0, 80, 81, 5, 3, 0, 0, 81, 82, 8, 1, 0, 0, 82,
		83, 3, 28, 14, 0, 83, 84, 5, 3, 0, 0, 84, 11, 1, 0, 0, 0, 85, 86, 5, 13,
		0, 0, 86, 87, 5, 3, 0, 0, 87, 88, 5, 13, 0, 0, 88, 13, 1, 0, 0, 0, 89,
		90, 5, 13, 0, 0, 90, 91, 5, 5, 0, 0, 91, 92, 5, 13, 0, 0, 92, 15, 1, 0,
		0, 0, 93, 94, 5, 5, 0, 0, 94, 95, 8, 1, 0, 0, 95, 96, 3, 28, 14, 0, 96,
		97, 5, 5, 0, 0, 97, 17, 1, 0, 0, 0, 98, 99, 5, 6, 0, 0, 99, 100, 3, 28,
		14, 0, 100, 101, 5, 7, 0, 0, 101, 105, 5, 8, 0, 0, 102, 104, 8, 2, 0, 0,
		103, 102, 1, 0, 0, 0, 104, 107, 1, 0, 0, 0, 105, 103, 1, 0, 0, 0, 105,
		106, 1, 0, 0, 0, 106, 108, 1, 0, 0, 0, 107, 105, 1, 0, 0, 0, 108, 109,
		5, 9, 0, 0, 109, 19, 1, 0, 0, 0, 110, 112, 3, 22, 11, 0, 111, 110, 1, 0,
		0, 0, 112, 113, 1, 0, 0, 0, 113, 111, 1, 0, 0, 0, 113, 114, 1, 0, 0, 0,
		114, 115, 1, 0, 0, 0, 115, 116, 3, 30, 15, 0, 116, 21, 1, 0, 0, 0, 117,
		121, 5, 10, 0, 0, 118, 120, 8, 0, 0, 0, 119, 118, 1, 0, 0, 0, 120, 123,
		1, 0, 0, 0, 121, 119, 1, 0, 0, 0, 121, 122, 1, 0, 0, 0, 122, 124, 1, 0,
		0, 0, 123, 121, 1, 0, 0, 0, 124, 125, 5, 1, 0, 0, 125, 23, 1, 0, 0, 0,
		126, 128, 3, 26, 13, 0, 127, 126, 1, 0, 0, 0, 128, 129, 1, 0, 0, 0, 129,
		127, 1, 0, 0, 0, 129, 130, 1, 0, 0, 0, 130, 131, 1, 0, 0, 0, 131, 132,
		3, 30, 15, 0, 132, 133, 3, 30, 15, 0, 133, 25, 1, 0, 0, 0, 134, 139, 5,
		4, 0, 0, 135, 137, 5, 4, 0, 0, 136, 138, 5, 4, 0, 0, 137, 136, 1, 0, 0,
		0, 137, 138, 1, 0, 0, 0, 138, 140, 1, 0, 0, 0, 139, 135, 1, 0, 0, 0, 139,
		140, 1, 0, 0, 0, 140, 142, 1, 0, 0, 0, 141, 134, 1, 0, 0, 0, 141, 142,
		1, 0, 0, 0, 142, 143, 1, 0, 0, 0, 143, 144, 5, 3, 0, 0, 144, 145, 5, 13,
		0, 0, 145, 146, 3, 8, 4, 0, 146, 27, 1, 0, 0, 0, 147, 149, 8, 3, 0, 0,
		148, 147, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150, 148, 1, 0, 0, 0, 150,
		151, 1, 0, 0, 0, 151, 29, 1, 0, 0, 0, 152, 154, 5, 12, 0, 0, 153, 152,
		1, 0, 0, 0, 153, 154, 1, 0, 0, 0, 154, 155, 1, 0, 0, 0, 155, 156, 5, 1,
		0, 0, 156, 31, 1, 0, 0, 0, 17, 35, 42, 47, 52, 60, 67, 76, 78, 105, 113,
		121, 129, 137, 139, 141, 150, 153,
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
	markdownParserEOF   = antlr.TokenEOF
	markdownParserT__0  = 1
	markdownParserT__1  = 2
	markdownParserT__2  = 3
	markdownParserT__3  = 4
	markdownParserT__4  = 5
	markdownParserT__5  = 6
	markdownParserT__6  = 7
	markdownParserT__7  = 8
	markdownParserT__8  = 9
	markdownParserT__9  = 10
	markdownParserT__10 = 11
	markdownParserT__11 = 12
	markdownParserWS    = 13
)

// markdownParser rules.
const (
	markdownParserRULE_file_       = 0
	markdownParserRULE_elem        = 1
	markdownParserRULE_header      = 2
	markdownParserRULE_para        = 3
	markdownParserRULE_paraContent = 4
	markdownParserRULE_bold        = 5
	markdownParserRULE_astericks   = 6
	markdownParserRULE_underscore  = 7
	markdownParserRULE_italics     = 8
	markdownParserRULE_link        = 9
	markdownParserRULE_quote       = 10
	markdownParserRULE_quoteElem   = 11
	markdownParserRULE_list        = 12
	markdownParserRULE_listElem    = 13
	markdownParserRULE_text        = 14
	markdownParserRULE_nl          = 15
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

func (s *File_Context) AllElem() []IElemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IElemContext); ok {
			len++
		}
	}

	tst := make([]IElemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IElemContext); ok {
			tst[i] = t.(IElemContext)
			i++
		}
	}

	return tst
}

func (s *File_Context) Elem(i int) IElemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElemContext); ok {
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

	return t.(IElemContext)
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
	p.SetState(33)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&14206) != 0 {
		{
			p.SetState(32)
			p.Elem()
		}

		p.SetState(35)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IElemContext is an interface to support dynamic dispatch.
type IElemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsElemContext differentiates from other interfaces.
	IsElemContext()
}

type ElemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyElemContext() *ElemContext {
	var p = new(ElemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_elem
	return p
}

func (*ElemContext) IsElemContext() {}

func NewElemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElemContext {
	var p = new(ElemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_elem

	return p
}

func (s *ElemContext) GetParser() antlr.Parser { return s.parser }

func (s *ElemContext) Header() IHeaderContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHeaderContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHeaderContext)
}

func (s *ElemContext) Para() IParaContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParaContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParaContext)
}

func (s *ElemContext) Quote() IQuoteContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQuoteContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQuoteContext)
}

func (s *ElemContext) List() IListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListContext)
}

func (s *ElemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ElemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ElemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterElem(s)
	}
}

func (s *ElemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitElem(s)
	}
}

func (p *markdownParser) Elem() (localctx IElemContext) {
	this := p
	_ = this

	localctx = NewElemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, markdownParserRULE_elem)

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

	p.SetState(42)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(37)
			p.Header()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(38)
			p.Para()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(39)
			p.Quote()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(40)
			p.List()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(41)
			p.Match(markdownParserT__0)
		}

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
	p.EnterRule(localctx, 4, markdownParserRULE_header)
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

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(45)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(44)
				p.Match(markdownParserT__1)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(47)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())
	}
	p.SetState(52)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&16380) != 0 {
		{
			p.SetState(49)
			_la = p.GetTokenStream().LA(1)

			if _la <= 0 || _la == markdownParserT__0 {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(54)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(55)
		p.Match(markdownParserT__0)
	}

	return localctx
}

// IParaContext is an interface to support dynamic dispatch.
type IParaContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParaContext differentiates from other interfaces.
	IsParaContext()
}

type ParaContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParaContext() *ParaContext {
	var p = new(ParaContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_para
	return p
}

func (*ParaContext) IsParaContext() {}

func NewParaContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParaContext {
	var p = new(ParaContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_para

	return p
}

func (s *ParaContext) GetParser() antlr.Parser { return s.parser }

func (s *ParaContext) ParaContent() IParaContentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParaContentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParaContentContext)
}

func (s *ParaContext) Nl() INlContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INlContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INlContext)
}

func (s *ParaContext) EOF() antlr.TerminalNode {
	return s.GetToken(markdownParserEOF, 0)
}

func (s *ParaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParaContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterPara(s)
	}
}

func (s *ParaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitPara(s)
	}
}

func (p *markdownParser) Para() (localctx IParaContext) {
	this := p
	_ = this

	localctx = NewParaContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, markdownParserRULE_para)

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

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(60)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(57)
				p.Match(markdownParserT__0)
			}

		}
		p.SetState(62)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext())
	}
	{
		p.SetState(63)
		p.ParaContent()
	}
	{
		p.SetState(64)
		p.Match(markdownParserT__0)
	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case markdownParserT__0, markdownParserT__11:
		{
			p.SetState(65)
			p.Nl()
		}

	case markdownParserEOF:
		{
			p.SetState(66)
			p.Match(markdownParserEOF)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IParaContentContext is an interface to support dynamic dispatch.
type IParaContentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParaContentContext differentiates from other interfaces.
	IsParaContentContext()
}

type ParaContentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParaContentContext() *ParaContentContext {
	var p = new(ParaContentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_paraContent
	return p
}

func (*ParaContentContext) IsParaContentContext() {}

func NewParaContentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParaContentContext {
	var p = new(ParaContentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_paraContent

	return p
}

func (s *ParaContentContext) GetParser() antlr.Parser { return s.parser }

func (s *ParaContentContext) AllText() []ITextContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITextContext); ok {
			len++
		}
	}

	tst := make([]ITextContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITextContext); ok {
			tst[i] = t.(ITextContext)
			i++
		}
	}

	return tst
}

func (s *ParaContentContext) Text(i int) ITextContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITextContext); ok {
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

	return t.(ITextContext)
}

func (s *ParaContentContext) AllBold() []IBoldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBoldContext); ok {
			len++
		}
	}

	tst := make([]IBoldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBoldContext); ok {
			tst[i] = t.(IBoldContext)
			i++
		}
	}

	return tst
}

func (s *ParaContentContext) Bold(i int) IBoldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBoldContext); ok {
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

	return t.(IBoldContext)
}

func (s *ParaContentContext) AllItalics() []IItalicsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IItalicsContext); ok {
			len++
		}
	}

	tst := make([]IItalicsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IItalicsContext); ok {
			tst[i] = t.(IItalicsContext)
			i++
		}
	}

	return tst
}

func (s *ParaContentContext) Italics(i int) IItalicsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IItalicsContext); ok {
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

	return t.(IItalicsContext)
}

func (s *ParaContentContext) AllLink() []ILinkContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILinkContext); ok {
			len++
		}
	}

	tst := make([]ILinkContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILinkContext); ok {
			tst[i] = t.(ILinkContext)
			i++
		}
	}

	return tst
}

func (s *ParaContentContext) Link(i int) ILinkContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILinkContext); ok {
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

	return t.(ILinkContext)
}

func (s *ParaContentContext) AllAstericks() []IAstericksContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAstericksContext); ok {
			len++
		}
	}

	tst := make([]IAstericksContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAstericksContext); ok {
			tst[i] = t.(IAstericksContext)
			i++
		}
	}

	return tst
}

func (s *ParaContentContext) Astericks(i int) IAstericksContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAstericksContext); ok {
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

	return t.(IAstericksContext)
}

func (s *ParaContentContext) AllUnderscore() []IUnderscoreContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IUnderscoreContext); ok {
			len++
		}
	}

	tst := make([]IUnderscoreContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IUnderscoreContext); ok {
			tst[i] = t.(IUnderscoreContext)
			i++
		}
	}

	return tst
}

func (s *ParaContentContext) Underscore(i int) IUnderscoreContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnderscoreContext); ok {
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

	return t.(IUnderscoreContext)
}

func (s *ParaContentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParaContentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParaContentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterParaContent(s)
	}
}

func (s *ParaContentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitParaContent(s)
	}
}

func (p *markdownParser) ParaContent() (localctx IParaContentContext) {
	this := p
	_ = this

	localctx = NewParaContentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, markdownParserRULE_paraContent)

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

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(76)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(76)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) {
			case 1:
				{
					p.SetState(69)
					p.Text()
				}

			case 2:
				{
					p.SetState(70)
					p.Bold()
				}

			case 3:
				{
					p.SetState(71)
					p.Italics()
				}

			case 4:
				{
					p.SetState(72)
					p.Link()
				}

			case 5:
				{
					p.SetState(73)
					p.Astericks()
				}

			case 6:
				{
					p.SetState(74)
					p.Underscore()
				}

			case 7:
				{
					p.SetState(75)
					p.Match(markdownParserT__0)
				}

			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(78)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())
	}

	return localctx
}

// IBoldContext is an interface to support dynamic dispatch.
type IBoldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoldContext differentiates from other interfaces.
	IsBoldContext()
}

type BoldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoldContext() *BoldContext {
	var p = new(BoldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_bold
	return p
}

func (*BoldContext) IsBoldContext() {}

func NewBoldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoldContext {
	var p = new(BoldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_bold

	return p
}

func (s *BoldContext) GetParser() antlr.Parser { return s.parser }

func (s *BoldContext) Text() ITextContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITextContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITextContext)
}

func (s *BoldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BoldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterBold(s)
	}
}

func (s *BoldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitBold(s)
	}
}

func (p *markdownParser) Bold() (localctx IBoldContext) {
	this := p
	_ = this

	localctx = NewBoldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, markdownParserRULE_bold)
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
		p.SetState(80)
		p.Match(markdownParserT__2)
	}
	{
		p.SetState(81)
		_la = p.GetTokenStream().LA(1)

		if _la <= 0 || _la == markdownParserT__0 || _la == markdownParserT__3 {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(82)
		p.Text()
	}
	{
		p.SetState(83)
		p.Match(markdownParserT__2)
	}

	return localctx
}

// IAstericksContext is an interface to support dynamic dispatch.
type IAstericksContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAstericksContext differentiates from other interfaces.
	IsAstericksContext()
}

type AstericksContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAstericksContext() *AstericksContext {
	var p = new(AstericksContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_astericks
	return p
}

func (*AstericksContext) IsAstericksContext() {}

func NewAstericksContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AstericksContext {
	var p = new(AstericksContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_astericks

	return p
}

func (s *AstericksContext) GetParser() antlr.Parser { return s.parser }

func (s *AstericksContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(markdownParserWS)
}

func (s *AstericksContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(markdownParserWS, i)
}

func (s *AstericksContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AstericksContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AstericksContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterAstericks(s)
	}
}

func (s *AstericksContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitAstericks(s)
	}
}

func (p *markdownParser) Astericks() (localctx IAstericksContext) {
	this := p
	_ = this

	localctx = NewAstericksContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, markdownParserRULE_astericks)

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
		p.SetState(85)
		p.Match(markdownParserWS)
	}
	{
		p.SetState(86)
		p.Match(markdownParserT__2)
	}
	{
		p.SetState(87)
		p.Match(markdownParserWS)
	}

	return localctx
}

// IUnderscoreContext is an interface to support dynamic dispatch.
type IUnderscoreContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUnderscoreContext differentiates from other interfaces.
	IsUnderscoreContext()
}

type UnderscoreContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnderscoreContext() *UnderscoreContext {
	var p = new(UnderscoreContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_underscore
	return p
}

func (*UnderscoreContext) IsUnderscoreContext() {}

func NewUnderscoreContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnderscoreContext {
	var p = new(UnderscoreContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_underscore

	return p
}

func (s *UnderscoreContext) GetParser() antlr.Parser { return s.parser }

func (s *UnderscoreContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(markdownParserWS)
}

func (s *UnderscoreContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(markdownParserWS, i)
}

func (s *UnderscoreContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnderscoreContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnderscoreContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterUnderscore(s)
	}
}

func (s *UnderscoreContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitUnderscore(s)
	}
}

func (p *markdownParser) Underscore() (localctx IUnderscoreContext) {
	this := p
	_ = this

	localctx = NewUnderscoreContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, markdownParserRULE_underscore)

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
		p.SetState(89)
		p.Match(markdownParserWS)
	}
	{
		p.SetState(90)
		p.Match(markdownParserT__4)
	}
	{
		p.SetState(91)
		p.Match(markdownParserWS)
	}

	return localctx
}

// IItalicsContext is an interface to support dynamic dispatch.
type IItalicsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsItalicsContext differentiates from other interfaces.
	IsItalicsContext()
}

type ItalicsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyItalicsContext() *ItalicsContext {
	var p = new(ItalicsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_italics
	return p
}

func (*ItalicsContext) IsItalicsContext() {}

func NewItalicsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ItalicsContext {
	var p = new(ItalicsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_italics

	return p
}

func (s *ItalicsContext) GetParser() antlr.Parser { return s.parser }

func (s *ItalicsContext) Text() ITextContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITextContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITextContext)
}

func (s *ItalicsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ItalicsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ItalicsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterItalics(s)
	}
}

func (s *ItalicsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitItalics(s)
	}
}

func (p *markdownParser) Italics() (localctx IItalicsContext) {
	this := p
	_ = this

	localctx = NewItalicsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, markdownParserRULE_italics)
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
		p.SetState(93)
		p.Match(markdownParserT__4)
	}
	{
		p.SetState(94)
		_la = p.GetTokenStream().LA(1)

		if _la <= 0 || _la == markdownParserT__0 || _la == markdownParserT__3 {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(95)
		p.Text()
	}
	{
		p.SetState(96)
		p.Match(markdownParserT__4)
	}

	return localctx
}

// ILinkContext is an interface to support dynamic dispatch.
type ILinkContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLinkContext differentiates from other interfaces.
	IsLinkContext()
}

type LinkContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLinkContext() *LinkContext {
	var p = new(LinkContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_link
	return p
}

func (*LinkContext) IsLinkContext() {}

func NewLinkContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LinkContext {
	var p = new(LinkContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_link

	return p
}

func (s *LinkContext) GetParser() antlr.Parser { return s.parser }

func (s *LinkContext) Text() ITextContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITextContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITextContext)
}

func (s *LinkContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LinkContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LinkContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterLink(s)
	}
}

func (s *LinkContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitLink(s)
	}
}

func (p *markdownParser) Link() (localctx ILinkContext) {
	this := p
	_ = this

	localctx = NewLinkContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, markdownParserRULE_link)
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
		p.SetState(98)
		p.Match(markdownParserT__5)
	}
	{
		p.SetState(99)
		p.Text()
	}
	{
		p.SetState(100)
		p.Match(markdownParserT__6)
	}
	{
		p.SetState(101)
		p.Match(markdownParserT__7)
	}
	p.SetState(105)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&15870) != 0 {
		{
			p.SetState(102)
			_la = p.GetTokenStream().LA(1)

			if _la <= 0 || _la == markdownParserT__8 {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(108)
		p.Match(markdownParserT__8)
	}

	return localctx
}

// IQuoteContext is an interface to support dynamic dispatch.
type IQuoteContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQuoteContext differentiates from other interfaces.
	IsQuoteContext()
}

type QuoteContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuoteContext() *QuoteContext {
	var p = new(QuoteContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_quote
	return p
}

func (*QuoteContext) IsQuoteContext() {}

func NewQuoteContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QuoteContext {
	var p = new(QuoteContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_quote

	return p
}

func (s *QuoteContext) GetParser() antlr.Parser { return s.parser }

func (s *QuoteContext) Nl() INlContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INlContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INlContext)
}

func (s *QuoteContext) AllQuoteElem() []IQuoteElemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IQuoteElemContext); ok {
			len++
		}
	}

	tst := make([]IQuoteElemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IQuoteElemContext); ok {
			tst[i] = t.(IQuoteElemContext)
			i++
		}
	}

	return tst
}

func (s *QuoteContext) QuoteElem(i int) IQuoteElemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQuoteElemContext); ok {
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

	return t.(IQuoteElemContext)
}

func (s *QuoteContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QuoteContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QuoteContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterQuote(s)
	}
}

func (s *QuoteContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitQuote(s)
	}
}

func (p *markdownParser) Quote() (localctx IQuoteContext) {
	this := p
	_ = this

	localctx = NewQuoteContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, markdownParserRULE_quote)
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
	p.SetState(111)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == markdownParserT__9 {
		{
			p.SetState(110)
			p.QuoteElem()
		}

		p.SetState(113)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(115)
		p.Nl()
	}

	return localctx
}

// IQuoteElemContext is an interface to support dynamic dispatch.
type IQuoteElemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQuoteElemContext differentiates from other interfaces.
	IsQuoteElemContext()
}

type QuoteElemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuoteElemContext() *QuoteElemContext {
	var p = new(QuoteElemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_quoteElem
	return p
}

func (*QuoteElemContext) IsQuoteElemContext() {}

func NewQuoteElemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QuoteElemContext {
	var p = new(QuoteElemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_quoteElem

	return p
}

func (s *QuoteElemContext) GetParser() antlr.Parser { return s.parser }
func (s *QuoteElemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QuoteElemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QuoteElemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterQuoteElem(s)
	}
}

func (s *QuoteElemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitQuoteElem(s)
	}
}

func (p *markdownParser) QuoteElem() (localctx IQuoteElemContext) {
	this := p
	_ = this

	localctx = NewQuoteElemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, markdownParserRULE_quoteElem)
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
		p.SetState(117)
		p.Match(markdownParserT__9)
	}
	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&16380) != 0 {
		{
			p.SetState(118)
			_la = p.GetTokenStream().LA(1)

			if _la <= 0 || _la == markdownParserT__0 {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(123)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(124)
		p.Match(markdownParserT__0)
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

func (s *ListContext) AllNl() []INlContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INlContext); ok {
			len++
		}
	}

	tst := make([]INlContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INlContext); ok {
			tst[i] = t.(INlContext)
			i++
		}
	}

	return tst
}

func (s *ListContext) Nl(i int) INlContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INlContext); ok {
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

	return t.(INlContext)
}

func (s *ListContext) AllListElem() []IListElemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IListElemContext); ok {
			len++
		}
	}

	tst := make([]IListElemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IListElemContext); ok {
			tst[i] = t.(IListElemContext)
			i++
		}
	}

	return tst
}

func (s *ListContext) ListElem(i int) IListElemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListElemContext); ok {
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

	return t.(IListElemContext)
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
	p.EnterRule(localctx, 24, markdownParserRULE_list)
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
	p.SetState(127)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == markdownParserT__2 || _la == markdownParserT__3 {
		{
			p.SetState(126)
			p.ListElem()
		}

		p.SetState(129)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(131)
		p.Nl()
	}
	{
		p.SetState(132)
		p.Nl()
	}

	return localctx
}

// IListElemContext is an interface to support dynamic dispatch.
type IListElemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsListElemContext differentiates from other interfaces.
	IsListElemContext()
}

type ListElemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyListElemContext() *ListElemContext {
	var p = new(ListElemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_listElem
	return p
}

func (*ListElemContext) IsListElemContext() {}

func NewListElemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListElemContext {
	var p = new(ListElemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_listElem

	return p
}

func (s *ListElemContext) GetParser() antlr.Parser { return s.parser }

func (s *ListElemContext) WS() antlr.TerminalNode {
	return s.GetToken(markdownParserWS, 0)
}

func (s *ListElemContext) ParaContent() IParaContentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParaContentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParaContentContext)
}

func (s *ListElemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListElemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListElemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterListElem(s)
	}
}

func (s *ListElemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitListElem(s)
	}
}

func (p *markdownParser) ListElem() (localctx IListElemContext) {
	this := p
	_ = this

	localctx = NewListElemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, markdownParserRULE_listElem)
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
	p.SetState(141)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == markdownParserT__3 {
		{
			p.SetState(134)
			p.Match(markdownParserT__3)
		}
		p.SetState(139)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == markdownParserT__3 {
			{
				p.SetState(135)
				p.Match(markdownParserT__3)
			}
			p.SetState(137)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == markdownParserT__3 {
				{
					p.SetState(136)
					p.Match(markdownParserT__3)
				}

			}

		}

	}
	{
		p.SetState(143)
		p.Match(markdownParserT__2)
	}
	{
		p.SetState(144)
		p.Match(markdownParserWS)
	}
	{
		p.SetState(145)
		p.ParaContent()
	}

	return localctx
}

// ITextContext is an interface to support dynamic dispatch.
type ITextContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTextContext differentiates from other interfaces.
	IsTextContext()
}

type TextContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTextContext() *TextContext {
	var p = new(TextContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_text
	return p
}

func (*TextContext) IsTextContext() {}

func NewTextContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TextContext {
	var p = new(TextContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_text

	return p
}

func (s *TextContext) GetParser() antlr.Parser { return s.parser }
func (s *TextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TextContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterText(s)
	}
}

func (s *TextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitText(s)
	}
}

func (p *markdownParser) Text() (localctx ITextContext) {
	this := p
	_ = this

	localctx = NewTextContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, markdownParserRULE_text)
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

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(148)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(147)
				_la = p.GetTokenStream().LA(1)

				if _la <= 0 || (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3310) != 0 {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(150)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext())
	}

	return localctx
}

// INlContext is an interface to support dynamic dispatch.
type INlContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNlContext differentiates from other interfaces.
	IsNlContext()
}

type NlContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNlContext() *NlContext {
	var p = new(NlContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = markdownParserRULE_nl
	return p
}

func (*NlContext) IsNlContext() {}

func NewNlContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NlContext {
	var p = new(NlContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = markdownParserRULE_nl

	return p
}

func (s *NlContext) GetParser() antlr.Parser { return s.parser }
func (s *NlContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NlContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NlContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.EnterNl(s)
	}
}

func (s *NlContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(markdownListener); ok {
		listenerT.ExitNl(s)
	}
}

func (p *markdownParser) Nl() (localctx INlContext) {
	this := p
	_ = this

	localctx = NewNlContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, markdownParserRULE_nl)
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
	p.SetState(153)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == markdownParserT__11 {
		{
			p.SetState(152)
			p.Match(markdownParserT__11)
		}

	}
	{
		p.SetState(155)
		p.Match(markdownParserT__0)
	}

	return localctx
}
