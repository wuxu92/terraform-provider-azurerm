package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"io"
	"os"
)

func NewMarkdownDoc(conf string) *MD {
	is := antlr.NewInputStream(conf)
	lex := NewmarkdownLexer(is)
	ts := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
	p := NewmarkdownParser(ts)

	// parse content
	l := &MD{}

	antlr.ParseTreeWalkerDefault.Walk(l, p.File_())
	return l
}

func NewMarkdownDocByFilename(f string) (*MD, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	md := NewMarkdownDoc(string(content))
	return md, nil
}
