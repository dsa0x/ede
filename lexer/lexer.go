package lexer

import (
	"ede/token"
	"unicode"
)

type Lexer struct {
	input []byte
	// readers func
	prevPos int
	currPos int
	readPos int
	char    byte
	charStr string

	line     int
	column   int
	startCol int
}

func New(input string) *Lexer {
	l := &Lexer{input: []byte(input), line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.char = byte(0)
	} else {
		l.char = l.input[l.readPos]
	}

	l.charStr = string(l.char)
	l.prevPos = l.currPos
	l.currPos = l.readPos
	l.readPos += 1

	if l.char == byte(0) {
		return
	}

	// if the previous line has ended, we increment the
	// line count, and restart the column count
	if l.input[l.prevPos] == byte('\n') {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

func (l *Lexer) readNChars(n int) {
	for i := 0; i < n; i++ {
		l.readChar()
	}
}

func (l *Lexer) readIdent() []byte {
	start := l.currPos
	for l.isIdentifier(l.peekChar()) {
		// if l.char == 34
		l.readChar()
	}
	return l.input[start:l.readPos]
}

func (l *Lexer) readDigit() []byte {
	start := l.currPos
	for unicode.IsDigit(rune(l.peekChar())) && !l.peekCharIs(byte(0)) {
		l.readChar()
	}
	return l.input[start:l.readPos]
}

func (l *Lexer) readString() []byte {
	l.readChar() // read the beginner
	start := l.currPos
	for l.char != '"' && l.readPos < len(l.input) && l.char != '\n' {
		l.readChar()
	}
	return l.input[start:l.currPos]
}

func (l *Lexer) readStringLiteral() []byte {
	l.readChar() // read the beginner
	start := l.currPos
	for l.char != '`' && l.readPos < len(l.input) {
		l.readChar()
	}
	return l.input[start:l.currPos]
}

func (l *Lexer) readReturn() []byte {
	start := l.currPos
	if l.peekCharIs('-') {
		l.readChar()
	}
	return l.input[start:l.readPos]
}

func (l *Lexer) readSingleComment() []byte {
	l.readNChars(2) // read '//'
	start := l.currPos
	for !(l.peekCharIs(';') || l.peekCharIs('\n')) {
		l.readChar()
	}
	l.readChar() // read closing char (\n or ;)
	return l.input[start:l.readPos]
}

func (l *Lexer) readStruct() []byte {
	start := l.currPos
	for unicode.IsDigit(rune(l.char)) {
		l.readChar()
	}
	return l.input[start:l.currPos]
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) currCharIs(char byte) bool {
	return l.char == char
}

func (l *Lexer) peekCharIs(char byte) bool {
	return l.peekChar() == char
}

func (l *Lexer) prevCharIs(char byte) bool {
	if l.prevPos >= len(l.input) {
		return false
	}
	return l.input[l.prevPos] == char
}

func (l *Lexer) eatWhitespace() {
	for l.isWhitespace(l.char) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()
	l.startCol = l.column // start of the token in the column

	switch l.char {
	case '=':
		// handling ==
		if l.peekCharIs('=') {
			l.readChar()
			tok = newToken(token.EQ, []byte("==")...)
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case ';', '{', '}', '(', ')', ',', '[', ']', ':':
		tok = charTokens[l.char]
	case '\n':
		tok = charTokens[l.char]
	case '+':
		if l.peekCharIs('+') {
			l.readChar()
			tok = newToken(token.INC, []byte("++")...)
		} else if l.peekCharIs('=') {
			l.readChar()
			tok = newToken(token.PLUS_EQUAL, []byte("+=")...)
		} else {
			tok = newToken(token.PLUS, l.char)
		}
	case '-':
		if l.peekCharIs('-') {
			l.readChar()
			tok = newToken(token.DEC, []byte("--")...)
		} else if l.peekCharIs('=') {
			l.readChar()
			tok = newToken(token.MINUS_EQUAL, []byte("-=")...)
		} else {
			tok = newToken(token.MINUS, l.char)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '/':
		if l.peekCharIs('/') {
			byt := l.readSingleComment()
			tok = newToken(token.SINGLE_COMMENT, byt...)
		} else {
			tok = newToken(token.SLASH, l.char)
		}
	case '&':
		if l.peekCharIs('&') {
			l.readChar()
			tok = newToken(token.AND_AND, []byte("&&")...)
		}
	case '|':
		if l.peekCharIs('|') {
			l.readChar()
			tok = newToken(token.OR_OR, []byte("||")...)
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	case '!':
		if l.peekCharIs('=') {
			l.readChar()
			tok = newToken(token.NEQ, []byte("!=")...)
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case '>':
		if l.peekCharIs('=') {
			l.readChar()
			tok = newToken(token.GTE, []byte(">=")...)
		} else {
			tok = newToken(token.GT, l.char)
		}
	case '<':
		if l.peekCharIs('-') {
			tok = newToken(token.RETURN, l.readReturn()...)
		} else if l.peekCharIs('=') {
			l.readChar()
			tok = newToken(token.LTE, []byte("<=")...)
		} else {
			tok = newToken(token.LT, l.char)
		}
	case '"':
		str := l.readString()
		tok = newToken(token.STRING, str...)
		if l.char != '"' {
			tok = newToken(token.ILLEGAL, str...)
		}
	case '`':
		str := l.readStringLiteral()
		tok = newToken(token.STRING, str...)
		if l.char != '`' {
			tok = newToken(token.ILLEGAL, str...)
		}
	case '%':
		tok = newToken(token.MODULO, l.char)
	case 0:
		tok = newToken(token.EOF)
	case '.':
		if l.peekCharIs('.') {
			l.readChar()
			tok = newToken(token.RANGE_ARRAY, '.', '.')
		} else {
			tok = newToken(token.DOT, '.')
		}
	default:
		if l.isIdentifier(l.char) {
			ident := l.readIdent()
			tokenType := token.LookupIdent(string(ident))
			tok = newToken(tokenType, ident...)
			break
		} else if unicode.IsDigit(rune(l.char)) {
			digit := l.readDigit()
			pos := l.CurrPos()
			l.readChar() // reads the next char after the digit
			if l.currCharIs('.') {
				// if it is not float or list, then it's invalid
				if !(l.peekCharIs('.') || unicode.IsDigit(rune(l.peekChar()))) {
					tok = newToken(token.ILLEGAL, l.peekChar())
					tok.Pos = pos
					return tok
				}
				if unicode.IsDigit(rune(l.peekChar())) {
					l.readChar()
					fraction := l.readDigit()
					tok = newToken(token.FLOAT, append(append(digit, '.'), fraction...)...)
					break
				}
			}
			tok = newToken(token.INT, digit...)
			tok.Pos = pos
			return tok // because we read the next char above, we need to return so we dont read again
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	tok.Pos = l.CurrPos()
	l.readChar()
	return tok
}

func newToken(typ token.TokenType, literal ...byte) token.Token {
	return token.Token{Type: typ, Literal: string(literal)}
}

func (l *Lexer) isIdentifier(char byte) bool {
	return unicode.IsLetter(rune(char)) || char == '_'
}
func (l *Lexer) isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' // || char == '\n'
}

var charTokens = map[byte]token.Token{
	';':  newToken(token.SEMICOLON, ';'),
	'{':  newToken(token.LBRACE, '{'),
	'}':  newToken(token.RBRACE, '}'),
	'(':  newToken(token.LPAREN, '('),
	')':  newToken(token.RPAREN, ')'),
	'[':  newToken(token.LBRACKET, '['),
	']':  newToken(token.RBRACKET, ']'),
	',':  newToken(token.COMMA, ','),
	':':  newToken(token.COLON, ':'),
	'\n': newToken(token.NEWLINE, '\n'),
}

func (l *Lexer) CurrPos() token.Pos { return token.Pos{Line: l.line, Column: l.Column()} }

// Column is the current column of the lexer
func (l *Lexer) Column() int { return l.startCol }

// Line is the current line of the lexer
func (l *Lexer) Line() int { return l.line }
