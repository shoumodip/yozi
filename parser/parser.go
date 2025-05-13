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
	powerLor
	powerCmp
	powerShl
	powerAdd
	powerBor
	powerMul
	powerAs
	powerPre
	powerDot
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

	token.LOr:  powerLor,
	token.LAnd: powerLor,

	token.Set: powerSet,

	token.Gt: powerCmp,
	token.Ge: powerCmp,
	token.Lt: powerCmp,
	token.Le: powerCmp,
	token.Eq: powerCmp,
	token.Ne: powerCmp,

	token.LParen: powerDot,

	token.As: powerAs,
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

func tokenKindIsStartOfType(k token.Kind) bool {
	return k == token.Ident || k == token.LAnd || k == token.BAnd || k == token.Fn
}

// @TokenKind
func (p *Parser) parseType() node.Node {
	tok := p.lexer.Next()
	switch tok.Kind {
	case token.Ident:
		return &node.Atom{
			Token: tok,
		}

	case token.LAnd:
		tok = p.lexer.SplitAndBufferRight(tok)
		fallthrough

	case token.BAnd:
		return &node.Unary{
			Token:   tok,
			Operand: p.parseType(),
		}

	case token.Fn:
		p.lexer.Expect(token.LParen)
		fn := node.Fn{
			Token: tok,
			Args:  []*node.Let{},
		}

		for !p.lexer.Read(token.RParen) {
			arg := node.Let{}

			argToken := p.lexer.Next()
			if argToken.Kind == token.Ident {
				peek := p.lexer.Peek()
				if peek.Kind == token.Comma || peek.Kind == token.RParen {
					// argToken is type
					arg.Token = tok
					arg.DefType = &node.Atom{
						Token: argToken,
					}
				} else {
					// argToken is name
					arg.Token = argToken
					arg.DefType = p.parseType()
				}
			} else {
				// argToken is type
				p.lexer.Buffer(argToken)
				arg.Token = tok
				arg.DefType = p.parseType()
			}

			arg.Kind = node.LetArg
			fn.Args = append(fn.Args, &arg)

			if p.lexer.Expect(token.Comma, token.RParen).Kind == token.RParen {
				break
			}
		}

		tok = p.lexer.Peek()
		if !tok.OnNewline && tokenKindIsStartOfType(tok.Kind) {
			fn.Return = p.parseType()
		}

		return &fn

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

		switch tok.Kind {
		case token.LParen:
			call := node.Call{
				Token: tok,
				Fn:    n,
				Args:  []node.Node{},
			}

			for !p.lexer.Read(token.RParen) {
				call.Args = append(call.Args, p.parseExpr(powerSet))
				if p.lexer.Expect(token.Comma, token.RParen).Kind == token.RParen {
					break
				}
			}

			n = &call

		case token.As:
			n = &node.Binary{
				Token: tok,
				Lhs:   n,
				Rhs:   p.parseType(),
			}

		default:
			n = &node.Binary{
				Token: tok,
				Lhs:   n,
				Rhs:   p.parseExpr(lbp),
			}
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

	case token.Return:
		p.localAssert(tok, true)

		ret := node.Return{Token: tok}
		if peek := p.lexer.Peek(); !peek.OnNewline {
			ret.Operand = p.parseExpr(powerSet)
		}

		return &ret

	case token.Fn:
		p.localAssert(tok, false) // TODO: Nested functions
		fn := node.Fn{
			Token:  p.lexer.Expect(token.Ident),
			Args:   []*node.Let{},
			Locals: []node.Node{},
		}

		p.local = true
		{
			p.lexer.Expect(token.LParen)
			for !p.lexer.Read(token.RParen) {
				arg := node.Let{}
				arg.Token = p.lexer.Expect(token.Ident)
				arg.Kind = node.LetArg
				arg.DefType = p.parseType()
				fn.Args = append(fn.Args, &arg)

				if p.lexer.Expect(token.Comma, token.RParen).Kind == token.RParen {
					break
				}
			}

			if peek := p.lexer.Peek(); peek.Kind != token.LBrace {
				fn.Return = p.parseType()
			}

			p.lexer.Buffer(p.lexer.Expect(token.LBrace))
			fn.Body = p.parseStmt().(*node.Block)
		}
		p.local = false

		return &fn

	case token.Let:
		let := node.Let{
			Token: p.lexer.Expect(token.Ident),
		}

		if tok := p.lexer.Peek(); tok.Kind != token.Set {
			let.DefType = p.parseType()
		}

		if tok := p.lexer.Peek(); tok.Kind == token.Set {
			p.lexer.Unbuffer()
			let.Assign = p.parseExpr(powerSet)
		}

		if p.local {
			let.Kind = node.LetLocal
		} else {
			let.Kind = node.LetGlobal
		}
		return &let

	case token.LBrace:
		p.localAssert(tok, true)
		body := []node.Node{}
		for {
			// Propagate the '}' as the node token
			if tok = p.lexer.Peek(); tok.Kind == token.RBrace {
				p.lexer.Unbuffer()
				break
			}

			body = append(body, p.parseStmt())
		}

		return &node.Block{
			Token: tok,
			Nodes: body,
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
