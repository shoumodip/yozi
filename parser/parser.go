package parser

import (
	"fmt"
	"os"
	"yozi/lexer"
	"yozi/node"
	"yozi/token"
)

const (
	powerNil = iota
	powerSet
	powerAdd
	powerMul
	powerPre
)

// @TokenKind
var tokenPowers = [token.COUNT]int{
	token.Add: powerAdd,
	token.Sub: powerAdd,

	token.Mul: powerMul,
	token.Div: powerMul,
}

func errorUnexpected(tok token.Token) {
	fmt.Fprintf(os.Stderr, "%s: ERROR: Unexpected %s\n", tok.Pos, token.Names[tok.Kind])
	os.Exit(1)
}

type Parser struct {
	lexer lexer.Lexer
	Nodes []node.Node
}

// @TokenKind
func (p *Parser) expr(mbp int) node.Node {
	var n node.Node

	tok := p.lexer.Next()
	switch tok.Kind {
	case token.Int, token.Bool:
		n = &node.Atom{
			Token: tok,
		}

	case token.Sub:
		n = &node.Unary{
			Token:   tok,
			Operand: p.expr(powerPre),
		}

	default:
		errorUnexpected(tok)
	}

	for true {
		tok = p.lexer.Peek()
		if tok.OnNewline {
			break
		}

		lbp := tokenPowers[tok.Kind]
		if lbp <= mbp {
			break
		}
		p.lexer.Unbuffer()

		n = &node.Binary{
			Token: tok,
			Lhs:   n,
			Rhs:   p.expr(lbp),
		}
	}

	return n
}

// @TokenKind
func (p *Parser) stmt() node.Node {
	var n node.Node

	switch tok := p.lexer.Next(); tok.Kind {
	case token.Print:
		n = &node.Print{
			Token:   tok,
			Operand: p.expr(powerSet),
		}

	default:
		p.lexer.Buffer(tok)
		n = p.expr(powerNil)
	}

	return n
}

func (p *Parser) File(lexer lexer.Lexer) {
	save := p.lexer
	if p.Nodes == nil {
		p.Nodes = []node.Node{}
	}

	p.lexer = lexer
	for !p.lexer.Read(token.Eof) {
		p.Nodes = append(p.Nodes, p.stmt())
	}
	p.lexer = save
}
