package lexers

import (
	"bytes"
    "fmt"
)

const (
	STATE_NORMAL = iota
	STATE_CHARLITERAL
	STATE_STRINGLITERAL
	STATE_RAWSTRINGLITERAL
	STATE_EOLCOMMENT
	STATE_MLCOMMENT
)

func cpp_lex_str(in *Input) TokenType {
	for {
		in.token = in.cursor
    /*!re2c
		re2c:eof = 0;
		re2c:define:YYCTYPE    = byte;
		re2c:define:YYPEEK     = "peek(in)";
		re2c:define:YYSKIP     = "in.cursor += 1";
		re2c:define:YYBACKUP   = "in.marker = in.cursor";
		re2c:define:YYRESTORE  = "in.cursor = in.marker";
		re2c:define:YYLESSTHAN = "in.limit <= in.cursor + @@{len}";
		re2c:define:YYFILL     = "fill(in) == 0";
        *                    { continue }
        $                    { return -1 }
        "\""                 { if in.state != STATE_STRINGLITERAL { continue }; in.state = STATE_NORMAL; return END }
        "\'"                 { if in.state != STATE_CHARLITERAL { continue }; in.state = STATE_NORMAL; return END }
        [a-zA-Z_0-9]+        { return STRINGWORD }
        "\\a"                { continue }
        "\\b"                { continue }
        "\\f"                { continue }
        "\\n"                { continue }
        "\\r"                { continue }
        "\\t"                { continue }
        "\\v"                { continue }
        "\\\\"               { continue }
        "\\'"                { continue }
        "\\\""               { continue }
        "\\?"                { continue }
        //"\\" [0-7]{1,3}      { lex_oct(in.tok, in.cur, u); continue; }
        //"\\u" [0-9a-fA-F]{4} { lex_hex(in.tok, in.cur, u); continue; }
        //"\\U" [0-9a-fA-F]{8} { lex_hex(in.tok, in.cur, u); continue; }
        //"\\x" [0-9a-fA-F]+   { if (!lex_hex(in.tok, in.cur, u)) return false; continue; }	
	*/
	}
}

func cpp_lex_raw_str(in *Input, start bool) TokenType {
	for {
		if start {
			in.raw_str_delim = in.data[in.token + 2:in.cursor - 1]
			start = false
		}
		in.token = in.cursor
    /*!re2c
		re2c:eof = 0;
		re2c:define:YYCTYPE    = byte;
		re2c:define:YYPEEK     = "peek(in)";
		re2c:define:YYSKIP     = "in.cursor += 1";
		re2c:define:YYBACKUP   = "in.marker = in.cursor";
		re2c:define:YYRESTORE  = "in.cursor = in.marker";
		re2c:define:YYLESSTHAN = "in.limit <= in.cursor + @@{len}";
		re2c:define:YYFILL     = "fill(in) == 0";
        *                    { continue }
        $                    { return -1 }
        "\n"                 { in.line += 1; continue }

		dchar = [a-zA-Z0-9_{}[\]#<>%:;.?*+-/^&|~!=,"’];
        ")" dchar* "\""      { if bytes.Compare(in.raw_str_delim, in.data[in.token+1:in.cursor-1]) != 0 { 
									//fmt.Printf("%s:%d: not end %s\n", in.filename, in.line, in.data[in.token+1:in.cursor-1]); 
									continue
								}
								in.state = STATE_NORMAL; in.raw_str_delim = nil; return END }
        [a-zA-Z_0-9]+        { return STRINGWORD }
	*/
	}
}

func cpp_lex_eol_comment(in *Input) TokenType {
	for {
		in.token = in.cursor
    /*!re2c
		re2c:eof = 0;
		re2c:define:YYCTYPE    = byte;
		re2c:define:YYPEEK     = "peek(in)";
		re2c:define:YYSKIP     = "in.cursor += 1";
		re2c:define:YYBACKUP   = "in.marker = in.cursor";
		re2c:define:YYRESTORE  = "in.cursor = in.marker";
		re2c:define:YYLESSTHAN = "in.limit <= in.cursor + @@{len}";
		re2c:define:YYFILL     = "fill(in) == 0";
        *                    { continue }
        "\n"                 { in.state = STATE_NORMAL; in.cursor -= 1; return END }
        $                    { return END }
        [a-zA-Z_0-9]+        { return COMMENTWORD }
	*/
	}
}

func cpp_lex_ml_comment(in *Input) TokenType {
	for {
		in.token = in.cursor
    /*!re2c
		re2c:eof = 0;
		re2c:define:YYCTYPE    = byte;
		re2c:define:YYPEEK     = "peek(in)";
		re2c:define:YYSKIP     = "in.cursor += 1";
		re2c:define:YYBACKUP   = "in.marker = in.cursor";
		re2c:define:YYRESTORE  = "in.cursor = in.marker";
		re2c:define:YYLESSTHAN = "in.limit <= in.cursor + @@{len}";
		re2c:define:YYFILL     = "fill(in) == 0";
        *                    { continue }
        "\n"                 { in.line += 1; continue }
        "*/"                 { in.state = STATE_NORMAL; return END }
        $                    { return END }
        [a-zA-Z_0-9]+        { return COMMENTWORD }
	*/
	}
}

func cpp_lex(in *Input) TokenType {
	for {
		in.token = in.cursor
        //fmt.Printf("start at %d\n", in.token)
		if in.state == STATE_STRINGLITERAL || in.state == STATE_CHARLITERAL {
			t := cpp_lex_str(in)
			if t >= 0 {
				return t
			}
		} else if (in.state == STATE_RAWSTRINGLITERAL) {
			t := cpp_lex_raw_str(in, false)
			if t >= 0 {
				return t
			}
		} else if (in.state == STATE_EOLCOMMENT) {
			t := cpp_lex_eol_comment(in)
			if t >= 0 {
				return t
			}
		} else if (in.state == STATE_MLCOMMENT) {
			t := cpp_lex_ml_comment(in)
			if t >= 0 {
				return t
			}
		}

	    was_bol := in.bol
		in.bol = false
    /*!re2c
		re2c:eof = 0;
		re2c:define:YYCTYPE    = byte;
		re2c:define:YYPEEK     = "peek(in)";
		re2c:define:YYSKIP     = "in.cursor += 1";
		re2c:define:YYBACKUP   = "in.marker = in.cursor";
		re2c:define:YYRESTORE  = "in.cursor = in.marker";
		re2c:define:YYLESSTHAN = "in.limit <= in.cursor + @@{len}";
		re2c:define:YYFILL     = "fill(in) == 0";

        scm = "//" [^\n]*;
        wsp = [ \t\v\r]+;
		newline = [\n];
		"\\" { continue }
        wsp { in.bol = was_bol; continue }
		newline { in.bol = true; in.line += 1; continue }

        * { fmt.Printf("%s: %d: match %2x\n", in.filename, in.line, peek(in)); continue }
        $ { return END }

        "@" { continue } // Objective-c

		decimal = [1-9][0-9]*;
		hex = "0x" [0-9a-fA-F]+;
		octal = "0" [0-7]*;

		decimal { return LITERAL }
		hex { return LITERAL }
		octal { return LITERAL }

		"#" [^\n]* { continue }
		"\"" { in.state = STATE_STRINGLITERAL; t := cpp_lex_str(in); if t >= 0 { return t }; continue }
		"\'" { in.state = STATE_CHARLITERAL; t := cpp_lex_str(in); if t >= 0 { return t }; continue }
		//dchar = [a-zA-Z0-9_{}[\]#<>%:;.?*+-/^&|~!=,"’];
		"R\"" dchar * "(" { in.state = STATE_RAWSTRINGLITERAL; t := cpp_lex_raw_str(in, true); if t >= 0 { return t }; continue }
		"//" { in.state = STATE_EOLCOMMENT; t := cpp_lex_eol_comment(in); if t >= 0 { return t }; continue }
		"/*" { in.state = STATE_MLCOMMENT; t := cpp_lex_ml_comment(in); if t >= 0 { return t }; continue }

        "auto" { return KEYWORD }
        "bool" { return KEYWORD }
        "break" { return KEYWORD }
        "char" { return KEYWORD }
        "const" { return KEYWORD }
        "continue" { return KEYWORD }
        "double" { return KEYWORD }
        "for" { return KEYWORD }
        "if" { return KEYWORD }
        "int" { return KEYWORD }
        "long" { return KEYWORD }
        "short" { return KEYWORD }
        "sizeof" { return KEYWORD }
        "struct" { return KEYWORD }
        "this" { return KEYWORD }
        "typedef" { return KEYWORD }
        "unsigned" { return KEYWORD }
        "return" { return KEYWORD }
        "void" { return KEYWORD }

		"+" { return PUNCTUATION }
		"-" { return PUNCTUATION }
		"/" { return PUNCTUATION }
		"*" { return PUNCTUATION }
		"%" { return PUNCTUATION }
		"&" { return PUNCTUATION }
		"|" { return PUNCTUATION }
		"&&" { return PUNCTUATION }
		"||" { return PUNCTUATION }
		"!" { return PUNCTUATION }
		"^" { return PUNCTUATION }
		"~" { return PUNCTUATION }
		";" { return PUNCTUATION }
		"." { return PUNCTUATION }
		"," { return PUNCTUATION }
		"(" { return PUNCTUATION }
		")" { return PUNCTUATION }
		"{" { return PUNCTUATION }
		"}" { return PUNCTUATION }
		"[" { return PUNCTUATION }
		"]" { return PUNCTUATION }
		"=" { return PUNCTUATION }
		"==" { return PUNCTUATION }
		"!=" { return PUNCTUATION }
		"<" { return PUNCTUATION }
		">" { return PUNCTUATION }
		"<=" { return PUNCTUATION }
		">=" { return PUNCTUATION }
		"->" { return PUNCTUATION }
		"?" { return PUNCTUATION }
		":" { return PUNCTUATION }
		"::" { return PUNCTUATION }

        word = [a-zA-Z_$][a-zA-Z_0-9$]*;
        word { return IDENTIFIER }
    */
    }
}
