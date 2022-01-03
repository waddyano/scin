package lexers

type TokenType int

const (
	END         TokenType = -1
	KEYWORD     TokenType = 1
	IDENTIFIER            = 2
	PUNCTUATION           = 3
	LITERAL               = 4
	CHARLITERAL           = 5
	STRINGWORD            = 6
	COMMENTWORD           = 7
)

func TypeString(l TokenType) string {
	switch l {
	case KEYWORD:
		return "KEYWORD"
	case IDENTIFIER:
		return "IDENTIFIER"
	case PUNCTUATION:
		return "PUNCTUATION"
	case LITERAL:
		return "LITERAL"
	case CHARLITERAL:
		return "CHARLITERAL"
	case STRINGWORD:
		return "STRINGWORD"
	case COMMENTWORD:
		return "COMMENTWORD"
	default:
		return "???"
	}

}
