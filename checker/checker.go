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
	if actual.Kind != node.TypeI64 && actual.Ref == 0 {
		fmt.Fprintf(os.Stderr, "%s: ERROR: Expected arithmetic type, got %s\n", n.Literal().Pos, actual)
		os.Exit(1)
	}

	return actual
}

func typeAssertScalar(n node.Node) node.Type {
	actual := n.GetType()
	if actual.Kind != node.TypeBool && actual.Kind != node.TypeI64 && actual.Ref == 0 {
		fmt.Fprintf(os.Stderr, "%s: ERROR: Expected scalar type, got %s\n", n.Literal().Pos, actual)
		os.Exit(1)
	}

	return actual
}

func errorUndefined(n node.Node, label string) {
	literal := n.Literal()
	fmt.Fprintf(os.Stderr, "%s: ERROR: Undefined %s '%s'\n", literal.Pos, label, literal.Str)
	os.Exit(1)
}

func errorRedefinition(n node.Node, prev node.Node, label string) {
	literal := n.Literal()
	fmt.Fprintf(os.Stderr, "%s: ERROR: Redefinition of %s '%s'\n", literal.Pos, label, literal.Str)
	fmt.Fprintf(os.Stderr, "%s: NOTE: Defined here\n", prev.Literal().Pos)
	os.Exit(1)
}

type Context struct {
	Globals map[string]node.Node
}

func NewContext() Context {
	return Context{
		Globals: make(map[string]node.Node),
	}
}

// @NodeKind
func (c *Context) Check(n node.Node) {
	switch n := n.(type) {
	case *node.Atom:
		literal := n.Literal()

		// @TokenKind
		switch literal.Kind {
		case token.Int:
			n.Type = node.Type{Kind: node.TypeI64}

		case token.Bool:
			n.Type = (node.Type{Kind: node.TypeBool})

		case token.Ident:
			if defined, ok := c.Globals[literal.Str]; ok {
				n.Defined = defined
				n.Type = defined.GetType()
				n.Memory = true
			} else {
				errorUndefined(n, "identifier")
			}

		default:
			panic("unreachable")
		}

	case *node.Unary:
		// @TokenKind
		switch n.Literal().Kind {
		case token.Sub:
			c.Check(n.Operand)
			n.Type = typeAssertArith(n.Operand)

		case token.Mul:
			c.Check(n.Operand)

			operandType := n.Operand.GetType()
			if operandType.Ref == 0 {
				fmt.Fprintf(
					os.Stderr,
					"%s: ERROR: Expected pointer, got %s\n",
					n.Operand.Literal().Pos,
					operandType,
				)
				os.Exit(1)
			}

			n.Type = operandType
			n.Type.Ref--
			n.Memory = true

		case token.BAnd:
			c.Check(n.Operand)
			if !n.Operand.IsMemory() {
				fmt.Fprintf(
					os.Stderr,
					"%s: ERROR: Cannot take reference of value not in memory\n",
					n.Operand.Literal().Pos,
				)
				os.Exit(1)
			}

			n.Type = n.Operand.GetType()
			n.Type.Ref++

		default:
			panic("unreachable")
		}

	case *node.Binary:
		// @TokenKind
		switch n.Literal().Kind {
		case token.Add, token.Sub, token.Mul, token.Div:
			c.Check(n.Lhs)
			c.Check(n.Rhs)
			n.Type = typeAssert(n.Rhs, typeAssertArith(n.Lhs))

		case token.Set:
			c.Check(n.Lhs)
			if !n.Lhs.IsMemory() {
				fmt.Fprintf(os.Stderr, "%s: ERROR: Cannot assign to value not in memory\n", n.Lhs.Literal().Pos)
				os.Exit(1)
			}

			c.Check(n.Rhs)
			typeAssert(n.Rhs, n.Lhs.GetType())
			n.Type = node.Type{Kind: node.TypeNil}

		case token.Gt, token.Ge, token.Lt, token.Le, token.Eq, token.Ne:
			c.Check(n.Lhs)
			c.Check(n.Rhs)
			typeAssert(n.Rhs, typeAssertArith(n.Lhs))
			n.Type = node.Type{Kind: node.TypeBool}

		default:
			panic("unreachable")
		}

	case *node.Print:
		c.Check(n.Operand)
		typeAssertScalar(n.Operand)

	case *node.Block:
		for _, it := range n.Body {
			c.Check(it)
		}

	case *node.If:
		c.Check(n.Condition)
		typeAssert(n.Condition, node.Type{Kind: node.TypeBool})
		c.Check(n.Consequent)
		c.Check(n.Antecedent)

	case *node.While:
		c.Check(n.Condition)
		typeAssert(n.Condition, node.Type{Kind: node.TypeBool})
		c.Check(n.Body)

	case *node.Let:
		if previous, ok := c.Globals[n.Token.Str]; ok {
			errorRedefinition(n, previous, "global identifier")
		}

		c.Check(n.Assign)
		n.Type = n.Assign.GetType()
		c.Globals[n.Token.Str] = n

	default:
		panic("unreachable")
	}
}
