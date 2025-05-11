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
	powerCmp
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

	token.Set: powerSet,

	token.Gt: powerCmp,
	token.Ge: powerCmp,
	token.Lt: powerCmp,
	token.Le: powerCmp,
	token.Eq: powerCmp,
	token.Ne: powerCmp,
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
	case token.Int, token.Bool, token.Ident:
		n = &node.Atom{
			Token: tok,
		}

	case token.Sub, token.Mul, token.BAnd:
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
	switch tok := p.lexer.Next(); tok.Kind {
	case token.Print:
		return &node.Print{
			Token:   tok,
			Operand: p.expr(powerSet),
		}

	case token.If:
		return p.ifBody(tok)

	case token.While:
		condition := p.expr(powerSet)
		body := p.block()

		return &node.While{
			Token:     tok,
			Condition: condition,
			Body:      body,
		}

	case token.Let:
		name := p.lexer.Expect(token.Ident)

		p.lexer.Expect(token.Set) // TODO: Implement definition by type
		value := p.expr(powerSet)

		return &node.Let{
			Token:  name,
			Assign: value,
		}

	case token.LBrace:
		return p.blockBody(tok)

	default:
		p.lexer.Buffer(tok)
		return p.expr(powerNil)
	}
}

func (p *Parser) ifBody(tok token.Token) *node.If {
	condition := p.expr(powerSet)
	consequent := p.block()
	antecedent := node.Node(&node.Block{})

	if p.lexer.Read(token.Else) {
		switch tok := p.lexer.Expect(token.LBrace, token.If); tok.Kind {
		case token.LBrace:
			antecedent = p.blockBody(tok)

		case token.If:
			antecedent = p.ifBody(tok)

		default:
			panic("unreachable")
		}
	}

	return &node.If{
		Token:      tok,
		Condition:  condition,
		Antecedent: antecedent,
		Consequent: consequent,
	}
}

func (p *Parser) block() *node.Block {
	brace := p.lexer.Expect(token.LBrace)
	return p.blockBody(brace)
}

func (p *Parser) blockBody(tok token.Token) *node.Block {
	block := []node.Node{}

	for !p.lexer.Read(token.RBrace) {
		block = append(block, p.stmt())
	}

	return &node.Block{
		Token: tok,
		Body:  block,
	}
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
