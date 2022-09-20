package lexer

import (
	"gwine/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}

}
func (l *Lexer) NextToken() token.Token {

	l.skipWhitespace()
	t := token.Token{Literal: string(l.ch)}
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			t.Type = token.EQ
			t.Literal = "=="
			l.readChar()
		} else {
			t.Type = token.ASSIGN
		}
	case '+':
		t.Type = token.PLUS
	case '-':
		t.Type = token.MINUS
	case '*':
		t.Type = token.ASTERISK
	case '/':
		t.Type = token.SLASH
	case '!':
		if l.peekChar() == '=' {
			t.Type = token.NEQ
			t.Literal = "!="
			l.readChar()
		} else {
			t.Type = token.BANG
		}
	case '<':
		t.Type = token.LT
	case '>':
		t.Type = token.GT
	case ',':
		t.Type = token.COMMA
	case ';':
		t.Type = token.SEMICOLON
	case '(':
		t.Type = token.LPAREN
	case ')':
		t.Type = token.RPAREN
	case '{':
		t.Type = token.LBRACE
	case '}':
		t.Type = token.RBRACE
	case '[':
		t.Type = token.LBRACKET
	case ']':
		t.Type = token.RBRACKET
	case '"':
		t.Type = token.STRING
		t.Literal = l.readString()
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.ch) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdent(t.Literal)
			return t
		} else if isDigital(l.ch) {
			t.Literal = l.readNumber()
			t.Type = token.INT
			return t
		} else {
			t.Type = token.ILLEGAL
		}
	}
	l.readChar()
	return t

}
func (l *Lexer) readIdentifier() string {
	p := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[p:l.position]
}
func (l *Lexer) readNumber() string {
	p := l.position
	for isDigital(l.ch) {
		l.readChar()
	}
	return l.input[p:l.position]
}
func (l *Lexer) readString() string{
	p := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0{
			break
		}
	}
	return l.input[p:l.position]
}
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
func isDigital(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func newToken(tt token.TokenType, ch byte) token.Token {
	return token.Token{Type: tt, Literal: string(ch)}
}
