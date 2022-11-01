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
		"", "'-'", "'*'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "LIST", "HEAD", "LINE", "WS", "LN", "END",
	}
	staticData.ruleNames = []string{
		"T__0", "T__1", "LIST", "HEAD", "LINE", "WS", "LN", "END",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 8, 54, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 1, 0, 1, 0, 1, 1, 1, 1, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 28, 8, 2, 1, 3, 4, 3, 31, 8, 3,
		11, 3, 12, 3, 32, 1, 4, 1, 4, 5, 4, 37, 8, 4, 10, 4, 12, 4, 40, 9, 4, 1,
		5, 1, 5, 1, 6, 4, 6, 45, 8, 6, 11, 6, 12, 6, 46, 1, 6, 1, 6, 1, 7, 1, 7,
		1, 7, 1, 7, 0, 0, 8, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8,
		1, 0, 4, 2, 0, 42, 42, 45, 45, 5, 0, 10, 10, 35, 35, 42, 42, 45, 45, 91,
		93, 2, 0, 10, 10, 13, 13, 2, 0, 9, 9, 32, 32, 58, 0, 1, 1, 0, 0, 0, 0,
		3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0,
		11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 1, 17, 1, 0, 0, 0,
		3, 19, 1, 0, 0, 0, 5, 27, 1, 0, 0, 0, 7, 30, 1, 0, 0, 0, 9, 34, 1, 0, 0,
		0, 11, 41, 1, 0, 0, 0, 13, 44, 1, 0, 0, 0, 15, 50, 1, 0, 0, 0, 17, 18,
		5, 45, 0, 0, 18, 2, 1, 0, 0, 0, 19, 20, 5, 42, 0, 0, 20, 4, 1, 0, 0, 0,
		21, 28, 7, 0, 0, 0, 22, 23, 5, 91, 0, 0, 23, 28, 5, 93, 0, 0, 24, 25, 5,
		91, 0, 0, 25, 26, 5, 120, 0, 0, 26, 28, 5, 93, 0, 0, 27, 21, 1, 0, 0, 0,
		27, 22, 1, 0, 0, 0, 27, 24, 1, 0, 0, 0, 28, 6, 1, 0, 0, 0, 29, 31, 5, 35,
		0, 0, 30, 29, 1, 0, 0, 0, 31, 32, 1, 0, 0, 0, 32, 30, 1, 0, 0, 0, 32, 33,
		1, 0, 0, 0, 33, 8, 1, 0, 0, 0, 34, 38, 8, 1, 0, 0, 35, 37, 8, 2, 0, 0,
		36, 35, 1, 0, 0, 0, 37, 40, 1, 0, 0, 0, 38, 36, 1, 0, 0, 0, 38, 39, 1,
		0, 0, 0, 39, 10, 1, 0, 0, 0, 40, 38, 1, 0, 0, 0, 41, 42, 7, 3, 0, 0, 42,
		12, 1, 0, 0, 0, 43, 45, 7, 2, 0, 0, 44, 43, 1, 0, 0, 0, 45, 46, 1, 0, 0,
		0, 46, 44, 1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47, 48, 1, 0, 0, 0, 48, 49,
		6, 6, 0, 0, 49, 14, 1, 0, 0, 0, 50, 51, 5, 0, 0, 1, 51, 52, 1, 0, 0, 0,
		52, 53, 6, 7, 0, 0, 53, 16, 1, 0, 0, 0, 5, 0, 27, 32, 38, 46, 1, 6, 0,
		0,
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
	markdownLexerT__0 = 1
	markdownLexerT__1 = 2
	markdownLexerLIST = 3
	markdownLexerHEAD = 4
	markdownLexerLINE = 5
	markdownLexerWS   = 6
	markdownLexerLN   = 7
	markdownLexerEND  = 8
)
