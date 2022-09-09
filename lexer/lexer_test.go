package lexer

import (
	"gwine/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		token := l.NextToken()
		token.Print()
		if token.Type != tt.expectedType {
			t.Fatalf("test %v :token expected type wrong,expected %v,got %v", i, tt.expectedType, token.Type)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("test %v :token expected literal wrong,expected %v,got %v", i, tt.expectedLiteral, token.Literal)
		}
	}
}

func TestRawToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = fn(x,y) {
		return x + y;
	}
	let result = add(five,ten);
	!-/*5;
	5<10>5;
	if (5 < 10) {
		return false;
	}else {
		return true;
	}
	10 != 1
	1 == 1
	`
	l := New(input)
	for {
		token := l.NextToken()
		token.Print()
		if token.Literal == "" {
			break
		}
	}
}
