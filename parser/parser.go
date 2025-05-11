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
	powerShl
	powerAdd
	powerBor
	powerMul
	powerPre
)

// @TokenKind
var tokenPowers = [token.COUNT]int{
	token.Add: powerAdd,
	token.Sub: powerAdd,

	token.Mul: powerMul,
	token.Div: powerMul,

	token.Shl:  powerShl,
	token.Shr:  powerShl,
	token.BOr:  powerBor,
	token.BAnd: powerBor,

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
	local bool
	Nodes []node.Node
}

// @TokenKind
func (p *Parser) parseType() node.Node {
	tok := p.lexer.Next()
	switch tok.Kind {
	case token.Ident:
		return &node.Atom{
			Token: tok,
		}

	case token.BAnd:
		return &node.Unary{
			Token:   tok,
			Operand: p.parseType(),
		}

	default:
		errorUnexpected(tok)
	}

	panic("unreachable")
}

// @TokenKind
func (p *Parser) parseExpr(mbp int) node.Node {
	var n node.Node

	tok := p.lexer.Next()
	switch tok.Kind {
	case token.Int, token.Bool, token.Ident:
		n = &node.Atom{
			Token: tok,
		}

	case token.Sub, token.Mul, token.BAnd, token.BNot, token.LNot:
		n = &node.Unary{
			Token:   tok,
			Operand: p.parseExpr(powerPre),
		}

	case token.LParen:
		n = p.parseExpr(powerSet)
		p.lexer.Expect(token.RParen)

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
			Rhs:   p.parseExpr(lbp),
		}
	}

	return n
}

func (p *Parser) localAssert(tok token.Token, local bool) {
	if p.local != local {
		scope := "global"
		if p.local {
			scope = "local"
		}

		fmt.Fprintf(
			os.Stderr,
			"%s: ERROR: Unexpected %s in %s scope\n",
			tok.Pos,
			token.Names[tok.Kind],
			scope,
		)
		os.Exit(1)
	}
}

// @TokenKind
func (p *Parser) parseStmt() node.Node {
	switch tok := p.lexer.Next(); tok.Kind {
	case token.Print:
		p.localAssert(tok, true)
		return &node.Print{
			Token:   tok,
			Operand: p.parseExpr(powerSet),
		}

	case token.If:
		p.localAssert(tok, true)
		condition := p.parseExpr(powerSet)

		p.lexer.Buffer(p.lexer.Expect(token.LBrace))
		consequent := p.parseStmt()

		antecedent := node.Node(&node.Block{})

		if p.lexer.Read(token.Else) {
			p.lexer.Buffer(p.lexer.Expect(token.LBrace, token.If))
			antecedent = p.parseStmt()
		}

		return &node.If{
			Token:      tok,
			Condition:  condition,
			Antecedent: antecedent,
			Consequent: consequent,
		}

	case token.While:
		p.localAssert(tok, true)
		condition := p.parseExpr(powerSet)

		p.lexer.Buffer(p.lexer.Expect(token.LBrace))
		body := p.parseStmt()

		return &node.While{
			Token:     tok,
			Condition: condition,
			Body:      body,
		}

	case token.Fn:
		p.localAssert(tok, false) // TODO: Nested functions
		name := p.lexer.Expect(token.Ident)
		fn := node.Fn{Token: name}

		p.local = true // TODO: Assuming functions can't be nested
		{
			p.lexer.Expect(token.LParen)
			p.lexer.Expect(token.RParen) // TODO: Function arguments

			p.lexer.Buffer(p.lexer.Expect(token.LBrace)) // TODO: Function return
			fn.Body = p.parseStmt()
		}
		p.local = false

		return &fn

	case token.Let:
		name := p.lexer.Expect(token.Ident)
		let := node.Let{Token: name}

		if p.lexer.Read(token.Set) {
			let.Assign = p.parseExpr(powerSet)
		} else {
			let.DefType = p.parseType()
			if p.lexer.Read(token.Set) {
				let.Assign = p.parseExpr(powerSet)
			}
		}

		return &let

	case token.LBrace:
		p.localAssert(tok, true)
		body := []node.Node{}
		for !p.lexer.Read(token.RBrace) {
			body = append(body, p.parseStmt())
		}

		return &node.Block{
			Token: tok,
			Body:  body,
		}

	default:
		p.localAssert(tok, true)
		p.lexer.Buffer(tok)
		return p.parseExpr(powerNil)
	}
}

func (p *Parser) File(lexer lexer.Lexer) {
	save := p.lexer
	if p.Nodes == nil {
		p.Nodes = []node.Node{}
	}

	p.lexer = lexer
	for !p.lexer.Read(token.Eof) {
		p.Nodes = append(p.Nodes, p.parseStmt())
	}
	p.lexer = save
}
