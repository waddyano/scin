//go:generate re2c --lang go cpp.re -o cpp.go

package lexers

import "path/filepath"

type Lexer struct {
	input *Input
}

func CanLex(filename string) bool {
	ext := filepath.Ext(filename)
	if ext != ".c" && ext != ".h" && ext != ".cpp" && ext != ".cc" && ext != ".cxx" {
		return false
	}
	return true
}

func NewLexer(filename string, input []byte) *Lexer {
	if !CanLex(filename) {
		return nil
	}

	in := &Input{
		filename: filename,
		file:     nil,
		data:     input,
		cursor:   0,
		marker:   0,
		token:    -1,
		limit:    len(input),
		line:     1,
		state:    STATE_NORMAL,
		eof:      false,
		bol:      true,
	}

	return &Lexer{input: in}
}

func (lexer *Lexer) Lex() TokenType {
	return cpp_lex(lexer.input)
}

func (lexer *Lexer) Line() int {
	return lexer.input.line
}

func (lexer *Lexer) TokenPos() (int, int) {
	return lexer.input.token, lexer.input.cursor
}

func (lexer *Lexer) Token() []byte {
	return lexer.input.data[lexer.input.token:lexer.input.cursor]
}

/*
	for {
		l := lex(in)
		//fmt.Printf("lex returns %d\n", l)
		if l < 0 {
			break
		}
		fmt.Printf("%d: next token %s %d \"%s\"\n", in.line, typeString(l), in.token, in.data[in.token:in.cursor])
	}
}
*/
