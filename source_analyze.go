package main

import (
	"lexers"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
)

type SourceAnalyzer struct {
	filename string
	comments bool
}

func (a SourceAnalyzer) Analyze(input []byte) analysis.TokenStream {
	ts := sourceAnalyze(a.filename, input, a.comments)
	ts = token.NewLowerCaseFilter().Filter(ts)
	return ts
}

func sourceAnalyze(filename string, input []byte, comments bool) analysis.TokenStream {
	lex := lexers.NewLexer(filename, input)
	res := make([]*analysis.Token, 0)
	//fmt.Printf("do analyze\n")

	for {
		tk := lex.Lex()
		//fmt.Printf("lex returns %d\n", l)
		if tk == lexers.END {
			break
		}

		if tk == lexers.PUNCTUATION {
			continue
		}

		isComment := tk == lexers.COMMENTWORD
		if comments != isComment {
			continue
		}
		pos, end := lex.TokenPos()
		//fmt.Printf("%d: next token %s %d \"%s\"\n", lex.Line(), lexers.TypeString(tk), pos, lex.Token())
		atok := &analysis.Token{Start: pos, End: end, Term: lex.Token()}
		res = append(res, atok)
	}

	return res
}
