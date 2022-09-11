package ast

import (
	"gwine/token"
	"testing"
)

func TestString(t *testing.T) {

	pg := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET,Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT,Literal: "onevar"},
					Value: "onevar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT,Literal: "twovar"},
					Value: "twovar",
				},
			},
		},
	}

	if pg.String() != "let onevar = twovar;" {
		t.Fatalf("string fatal")
	}

}