package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

type Pos int

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT   = "IDENT"   // let, if, for, etc
	BUILTIN = "BUILTIN" // print, len, etc
	INT     = "INT"     // 1, 2, 3
	STRING  = "STRING"

	// Operators
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	ASSIGN   = "="
	BANG     = "!"
	GT       = ">"
	LT       = "<"
	EQ       = "=="
	NEQ      = "!="
	DEC      = "--"
	INC      = "++"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	SINGLE_COMMENT = "//"

	// Keywords
	FUNCTION    = "FUNCTION"
	EXTEND      = "EXTEND"
	STRUCT      = "STRUCT"
	STRUCT_TYPE = "STRUCT_TYPE"
	LET         = "LET"
	IF          = "IF"
	ELSE        = "ELSE"
	RETURN      = "RETURN"
	TRUE        = "TRUE"
	FALSE       = "FALSE"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"struct": STRUCT,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
