// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type markdownLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var markdownlexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func markdownlexerLexerInit() {
	staticData := &markdownlexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.literalNames = []string{
		"", "'\\n'", "'#'", "'*'", "' '", "'_'", "'['", "']'", "'('", "')'",
		"'>'", "'`'", "'\\r'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "WS",
	}
	staticData.ruleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
		"T__9", "T__10", "T__11", "WS",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 13, 56, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1,
		2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1,
		8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1, 12, 4, 12, 53, 8, 12, 11,
		12, 12, 12, 54, 0, 0, 13, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15,
		8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 1, 0, 1, 2, 0, 9, 9, 32, 32,
		56, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0,
		0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0,
		0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1,
		0, 0, 0, 0, 25, 1, 0, 0, 0, 1, 27, 1, 0, 0, 0, 3, 29, 1, 0, 0, 0, 5, 31,
		1, 0, 0, 0, 7, 33, 1, 0, 0, 0, 9, 35, 1, 0, 0, 0, 11, 37, 1, 0, 0, 0, 13,
		39, 1, 0, 0, 0, 15, 41, 1, 0, 0, 0, 17, 43, 1, 0, 0, 0, 19, 45, 1, 0, 0,
		0, 21, 47, 1, 0, 0, 0, 23, 49, 1, 0, 0, 0, 25, 52, 1, 0, 0, 0, 27, 28,
		5, 10, 0, 0, 28, 2, 1, 0, 0, 0, 29, 30, 5, 35, 0, 0, 30, 4, 1, 0, 0, 0,
		31, 32, 5, 42, 0, 0, 32, 6, 1, 0, 0, 0, 33, 34, 5, 32, 0, 0, 34, 8, 1,
		0, 0, 0, 35, 36, 5, 95, 0, 0, 36, 10, 1, 0, 0, 0, 37, 38, 5, 91, 0, 0,
		38, 12, 1, 0, 0, 0, 39, 40, 5, 93, 0, 0, 40, 14, 1, 0, 0, 0, 41, 42, 5,
		40, 0, 0, 42, 16, 1, 0, 0, 0, 43, 44, 5, 41, 0, 0, 44, 18, 1, 0, 0, 0,
		45, 46, 5, 62, 0, 0, 46, 20, 1, 0, 0, 0, 47, 48, 5, 96, 0, 0, 48, 22, 1,
		0, 0, 0, 49, 50, 5, 13, 0, 0, 50, 24, 1, 0, 0, 0, 51, 53, 7, 0, 0, 0, 52,
		51, 1, 0, 0, 0, 53, 54, 1, 0, 0, 0, 54, 52, 1, 0, 0, 0, 54, 55, 1, 0, 0,
		0, 55, 26, 1, 0, 0, 0, 2, 0, 54, 0,
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

// markdownLexerInit initializes any static state used to implement markdownLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewmarkdownLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func MarkdownLexerInit() {
	staticData := &markdownlexerLexerStaticData
	staticData.once.Do(markdownlexerLexerInit)
}

// NewmarkdownLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewmarkdownLexer(input antlr.CharStream) *markdownLexer {
	MarkdownLexerInit()
	l := new(markdownLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &markdownlexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "markdown.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// markdownLexer tokens.
const (
	markdownLexerT__0  = 1
	markdownLexerT__1  = 2
	markdownLexerT__2  = 3
	markdownLexerT__3  = 4
	markdownLexerT__4  = 5
	markdownLexerT__5  = 6
	markdownLexerT__6  = 7
	markdownLexerT__7  = 8
	markdownLexerT__8  = 9
	markdownLexerT__9  = 10
	markdownLexerT__10 = 11
	markdownLexerT__11 = 12
	markdownLexerWS    = 13
)
