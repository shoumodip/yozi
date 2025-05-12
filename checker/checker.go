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

	locals    []node.Node
	currentFn *node.Fn
}

func NewContext() Context {
	return Context{
		Globals: make(map[string]node.Node),
	}
}

func (c *Context) checkType(n node.Node) {
	switch n := n.(type) {
	case *node.Atom:
		switch n.Token.Str {
		case "i64":
			n.Type = node.Type{Kind: node.TypeI64}

		case "bool":
			n.Type = node.Type{Kind: node.TypeBool}

		default:
			errorUndefined(n, "type")
		}

	case *node.Unary:
		c.checkType(n.Operand)
		n.Type = n.Operand.GetType()
		n.Type.Ref++

	case *node.Fn:
		for _, arg := range n.Args {
			c.checkType(arg.DefType)
			arg.Type = arg.DefType.GetType()
		}

		if n.Return != nil {
			c.checkType(n.Return)
		}

		n.Type = node.Type{
			Kind: node.TypeFn,
			Spec: n,
		}

	default:
		panic("unreachable")
	}
}

func (c *Context) argumentFind(name string, until int) (node.Node, bool) {
	for i, arg := range c.currentFn.Args {
		if i == until {
			break
		}

		if arg.Token.Str == name {
			return arg, true
		}
	}

	return nil, false
}

func (c *Context) variableFind(name string) (node.Node, bool) {
	if c.currentFn != nil {
		for i := len(c.locals) - 1; i >= 0; i-- {
			l := c.locals[i]
			if l.Literal().Str == name {
				return l, true
			}
		}

		if arg, ok := c.argumentFind(name, len(c.currentFn.Args)); ok {
			return arg, true
		}
	}

	// Bruh
	global, ok := c.Globals[name]
	return global, ok
}

func checkIfMemory(n node.Node, message string) {
	if !n.IsMemory() {
		fmt.Fprintf(os.Stderr, "%s: ERROR: %s\n", n.Literal().Pos, message)
		os.Exit(1)
	}

	var checkIfMemoryImpl func(n node.Node)
	checkIfMemoryImpl = func(n node.Node) {
		switch n := n.(type) {
		case *node.Atom:
			if let, ok := n.Defined.(*node.Let); ok {
				if let.Kind == node.LetArg {
					let.Kind = node.LetLocalArg
				}
			}

		case *node.Unary:
			checkIfMemoryImpl(n.Operand)

		case *node.Binary:
			checkIfMemoryImpl(n.Lhs)

		default:
			panic("unreachable")
		}
	}
	checkIfMemoryImpl(n)
}

// @NodeKind
func (c *Context) Check(n node.Node) {
	switch n := n.(type) {
	case *node.Atom:
		// @TokenKind
		switch n.Token.Kind {
		case token.Int:
			n.Type = node.Type{Kind: node.TypeI64}

		case token.Bool:
			n.Type = (node.Type{Kind: node.TypeBool})

		case token.Ident:
			if defined, ok := c.variableFind(n.Token.Str); ok {
				n.Defined = defined
				n.Type = defined.GetType()
				_, n.Memory = defined.(*node.Let)
			} else {
				errorUndefined(n, "identifier")
			}

		default:
			panic("unreachable")
		}

	case *node.Call:
		c.Check(n.Fn)

		fnTok := n.Fn.Literal()
		fnType := n.Fn.GetType()

		if fnType.Kind != node.TypeFn {
			fmt.Fprintf(
				os.Stderr,
				"%s: ERROR: Expected function, got %s\n",
				fnTok.Pos,
				fnType,
			)
			os.Exit(1)
		}

		if fnType.Ref != 0 {
			fmt.Fprintf(
				os.Stderr,
				"%s: ERROR: Cannot call pointer to function. Dereference it first\n",
				fnTok.Pos,
			)
			os.Exit(1)
		}

		fnSig := fnType.Spec.(*node.Fn)
		if len(n.Args) != len(fnSig.Args) {
			fmt.Fprintf(
				os.Stderr,
				"%s: ERROR: Expected %d arguments, got %d\n",
				n.Token.Pos,
				len(fnSig.Args),
				len(n.Args),
			)
			os.Exit(1)
		}

		for i, aArg := range n.Args {
			c.Check(aArg)
			typeAssert(aArg, fnSig.Args[i].Type)
		}

		n.Type = fnSig.ReturnType()

	case *node.Unary:
		// @TokenKind
		switch n.Token.Kind {
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
			checkIfMemory(n.Operand, "Cannot take reference of value not in memory")

			n.Type = n.Operand.GetType()
			n.Type.Ref++

		case token.BNot:
			c.Check(n.Operand)
			n.Type = typeAssertArith(n.Operand)

		case token.LNot:
			c.Check(n.Operand)
			n.Type = typeAssert(n.Operand, node.Type{Kind: node.TypeBool})

		default:
			panic("unreachable")
		}

	case *node.Binary:
		// @TokenKind
		switch n.Token.Kind {
		case token.Add, token.Sub, token.Mul, token.Div:
			c.Check(n.Lhs)
			c.Check(n.Rhs)
			n.Type = typeAssert(n.Rhs, typeAssertArith(n.Lhs))

		// These only work on integers, whereas the standard arithmetics branch can also work on floats
		case token.Shl, token.Shr, token.BOr, token.BAnd:
			c.Check(n.Lhs)
			c.Check(n.Rhs)
			n.Type = typeAssert(n.Rhs, typeAssertArith(n.Lhs))

		case token.Set:
			c.Check(n.Lhs)
			checkIfMemory(n.Lhs, "Cannot assign to value not in memory")

			c.Check(n.Rhs)
			typeAssert(n.Rhs, n.Lhs.GetType())
			n.Type = node.Type{Kind: node.TypeUnit}

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
		scopeStart := len(c.locals)
		for _, it := range n.Nodes {
			c.Check(it)
		}
		c.locals = c.locals[0:scopeStart]

	case *node.If:
		c.Check(n.Condition)
		typeAssert(n.Condition, node.Type{Kind: node.TypeBool})
		c.Check(n.Consequent)
		c.Check(n.Antecedent)

	case *node.While:
		c.Check(n.Condition)
		typeAssert(n.Condition, node.Type{Kind: node.TypeBool})
		c.Check(n.Body)

	case *node.Return:
		if n.Operand != nil {
			c.Check(n.Operand)
			n.Type = n.Operand.GetType()
		} else if c.currentFn.Return != nil {
			n.Type = node.Type{Kind: node.TypeUnit}
		}
		typeAssert(n, c.currentFn.ReturnType())

	case *node.Fn:
		if previous, ok := c.Globals[n.Token.Str]; ok {
			errorRedefinition(n, previous, "global identifier")
		}

		n.Type = node.Type{Kind: node.TypeFn, Spec: n}

		c.currentFn = n // TODO: Assuming functions can't be nested
		{
			scopeStart := len(c.locals)
			for i, arg := range n.Args {
				if previous, ok := c.argumentFind(arg.Token.Str, i); ok {
					errorRedefinition(arg, previous, "argument")
				}

				c.checkType(arg.DefType)
				arg.Type = arg.DefType.GetType()
			}

			if n.Return != nil {
				c.checkType(n.Return)
			}

			c.Globals[n.Token.Str] = n
			c.Check(n.Body)

			if n.Return != nil {
				// TODO: Implement proper return analysis
				endsWithReturn := len(n.Body.Nodes) > 0
				if endsWithReturn {
					_, endsWithReturn = n.Body.Nodes[len(n.Body.Nodes)-1].(*node.Return)
				}

				if !endsWithReturn {
					fmt.Fprintf(
						os.Stderr,
						"%s: ERROR: Expected last statement to be 'return'\n",
						n.Body.Token.Pos,
					)
					os.Exit(1)
				}
			}

			c.locals = c.locals[0:scopeStart]
		}
		c.currentFn = nil

	case *node.Let:
		if n.Kind == node.LetGlobal {
			if previous, ok := c.Globals[n.Token.Str]; ok {
				errorRedefinition(n, previous, "global identifier")
			}
		}

		if n.DefType != nil {
			c.checkType(n.DefType)
			n.Type = n.DefType.GetType()
		}

		if n.Assign != nil {
			c.Check(n.Assign)

			assignType := n.Assign.GetType()
			if assignType.Equal(node.Type{Kind: node.TypeUnit}) {
				fmt.Fprintf(
					os.Stderr,
					"%s: ERROR: Cannot define variable with type %s\n",
					n.Token.Pos,
					assignType,
				)
				os.Exit(1)
			}

			if n.DefType != nil {
				typeAssert(n.Assign, n.Type)
			} else {
				n.Type = assignType
			}
		}

		if n.Kind == node.LetGlobal {
			c.Globals[n.Token.Str] = n
		} else {
			c.locals = append(c.locals, n)
			c.currentFn.Locals = append(c.currentFn.Locals, n)
		}

	default:
		panic("unreachable")
	}
}
