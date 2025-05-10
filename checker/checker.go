package checker

import (
	"fmt"
	"os"
	"yozi/node"
	"yozi/token"
)

func typeAssert(n node.Node, expected node.Type) node.Type {
	actual := n.GetType()
	if !actual.Equal(expected) {
		fmt.Fprintf(os.Stderr, "%s: ERROR: Expected type %s, got %s\n", n.Literal().Pos, expected, actual)
		os.Exit(1)
	}

	return actual
}

func typeAssertArith(n node.Node) node.Type {
	actual := n.GetType()
	if actual.Kind != node.TypeI64 {
		fmt.Fprintf(os.Stderr, "%s: ERROR: Expected arithmetic type, got %s\n", n.Literal().Pos, actual)
		os.Exit(1)
	}

	return actual
}

// @NodeKind
func Check(n node.Node) {
	switch n := n.(type) {
	case *node.Atom:
		// @TokenKind
		switch n.Literal().Kind {
		case token.Int:
			n.SetType(node.Type{Kind: node.TypeI64})

		case token.Bool:
			n.SetType((node.Type{Kind: node.TypeBool}))

		default:
			panic("unreachable")
		}

	case *node.Unary:
		// @TokenKind
		switch n.Literal().Kind {
		case token.Sub:
			Check(n.Operand)
			n.SetType(typeAssertArith(n.Operand))

		default:
			panic("unreachable")
		}

	case *node.Binary:
		// @TokenKind
		switch n.Literal().Kind {
		case token.Add, token.Sub, token.Mul, token.Div:
			Check(n.Lhs)
			Check(n.Rhs)
			n.SetType(typeAssert(n.Rhs, typeAssertArith(n.Lhs)))

		default:
			panic("unreachable")
		}

	case *node.Print:
		Check(n.Operand)
		typeAssertArith(n.Operand)

	case *node.Block:
		for _, it := range n.Stmts {
			Check(it)
		}

	case *node.If:
		Check(n.Condition)
		typeAssert(n.Condition, node.Type{Kind: node.TypeBool})
		Check(n.Consequent)
		Check(n.Antecedent)

	case *node.While:
		Check(n.Condition)
		typeAssert(n.Condition, node.Type{Kind: node.TypeBool})
		Check(n.Body)

	default:
		panic("unreachable")
	}
}
