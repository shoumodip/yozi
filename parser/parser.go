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
func (p *Parser) expr(mbp int) node.Expr {
	var n node.Expr

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
func (p *Parser) stmt() node.Stmt {
	var n node.Stmt

	switch tok := p.lexer.Next(); tok.Kind {
	case token.Print:
		n = &node.Print{
			Token:   tok,
			Operand: p.expr(powerSet),
		}

	case token.If:
		n = p.ifBody(tok)

	case token.While:
		condition := p.expr(powerSet)
		body := p.block()

		n = &node.While{
			Token:     tok,
			Condition: condition,
			Body:      body,
		}

	case token.LBrace:
		n = p.blockBody(tok)

	default:
		p.lexer.Buffer(tok)
		n = &node.Block{}
	}

	return n
}

func (p *Parser) ifBody(tok token.Token) *node.If {
	condition := p.expr(powerSet)
	consequent := p.block()
	antecedent := node.Stmt(&node.Block{})

	if p.lexer.Read(token.Else) {
		switch tok := p.lexer.Next(); tok.Kind {
		case token.LBrace:
			antecedent = p.blockBody(tok)
		case token.If:
			antecedent = p.ifBody(tok)
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
	brace := p.lexer.Next()
	if brace.Kind != token.LBrace {
		errorUnexpected(p.lexer.Peek())
	}

	return p.blockBody(brace)
}

func (p *Parser) blockBody(tok token.Token) *node.Block {
	block := []node.Stmt{}

	for !p.lexer.Read(token.RBrace) {
		block = append(block, p.stmt())
	}

	return &node.Block{
		Token: tok,
		Stmts: block,
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
