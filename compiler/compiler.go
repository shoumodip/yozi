package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"yozi/checker"
	"yozi/node"
	"yozi/token"
)

type Compiler struct {
	context *checker.Context

	out     *os.File
	labelId int
	valueId int
}

// @Temporary
// @TypeKind
func llvmFormatType(t node.Type) string {
	sb := strings.Builder{}
	switch t.Kind {
	case node.TypeBool:
		sb.WriteString("i1")

	case node.TypeI8, node.TypeU8:
		sb.WriteString("i8")

	case node.TypeI16, node.TypeU16:
		sb.WriteString("i16")

	case node.TypeI32, node.TypeU32:
		sb.WriteString("i32")

	case node.TypeI64, node.TypeU64:
		sb.WriteString("i64")

	case node.TypeUnit:
		sb.WriteString("void")

	case node.TypeFn:
		fn := t.Spec.(*node.Fn)
		sb.WriteString(llvmFormatType(fn.ReturnType()))
		sb.WriteString(" (")
		for i, arg := range fn.Args {
			if i != 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(llvmFormatType(arg.GetType()))
		}
		sb.WriteString(")*")

	case node.TypeRawptr:
		sb.WriteString("i8*")

	default:
		panic("unreachable")
	}

	for range t.Ref {
		sb.WriteByte('*')
	}

	return sb.String()
}

func (c *Compiler) valueNew() string {
	c.valueId++
	return fmt.Sprintf("%%%d", c.valueId-1)
}

func (c *Compiler) labelNew() string {
	c.labelId++
	return fmt.Sprintf("L%d", c.labelId-1)
}

func (c *Compiler) binaryOp(n *node.Binary, op string) string {
	lhs := c.compileExpr(n.Lhs, false)
	rhs := c.compileExpr(n.Rhs, false)
	result := c.valueNew()

	fmt.Fprintf(c.out, "    %s = %s %s %s, %s\n", result, op, llvmFormatType(n.Lhs.GetType()), lhs, rhs)
	return result
}

// Like binaryOp, but pointers are first cast to i64
func (c *Compiler) binaryArithOp(n *node.Binary, op string) string {
	lhs := c.compileExpr(n.Lhs, false)
	rhs := c.compileExpr(n.Rhs, false)

	exprType := n.Lhs.GetType()
	tempType := "i64"

	llvmType := llvmFormatType(exprType)
	if exprType.Ref != 0 {
		lhsTemp := c.valueNew()
		rhsTemp := c.valueNew()
		fmt.Fprintf(c.out, "    %s = ptrtoint %s %s to %s\n", lhsTemp, llvmType, lhs, tempType)
		fmt.Fprintf(c.out, "    %s = ptrtoint %s %s to %s\n", rhsTemp, llvmType, rhs, tempType)

		lhs = lhsTemp
		rhs = rhsTemp
	}

	if exprType.Ref != 0 {
		tempResult := c.valueNew()
		fmt.Fprintf(c.out, "    %s = %s %s %s, %s\n", tempResult, op, tempType, lhs, rhs)

		result := c.valueNew()
		fmt.Fprintf(
			c.out,
			"    %s = inttoptr %s %s to %s\n",
			result,
			tempType,
			tempResult,
			llvmType,
		)

		return result
	} else {
		result := c.valueNew()
		fmt.Fprintf(c.out, "    %s = %s %s %s, %s\n", result, op, llvmType, lhs, rhs)
		return result
	}
}

func (c *Compiler) binaryLogicalOp(n *node.Binary) string {
	var shortCircuitValue int

	evalRhs := c.labelNew()
	shortCircuit := c.labelNew()
	merge := c.labelNew()

	lhs := c.compileExpr(n.Lhs, false)
	switch n.Token.Kind {
	case token.LOr:
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", lhs, shortCircuit, evalRhs)
		shortCircuitValue = 1

	case token.LAnd:
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", lhs, evalRhs, shortCircuit)
		shortCircuitValue = 0

	default:
		panic("unreachable")
	}
	fmt.Fprintf(c.out, "%s:\n", evalRhs)

	rhs := c.compileExpr(n.Rhs, false)
	fmt.Fprintf(c.out, "    br label %%%s\n", merge)
	fmt.Fprintf(c.out, "%s:\n", shortCircuit)
	fmt.Fprintf(c.out, "    br label %%%s\n", merge)
	fmt.Fprintf(c.out, "%s:\n", merge)

	result := c.valueNew()
	fmt.Fprintf(
		c.out,
		"    %s = phi i1 [ %d, %%%s ], [ %s, %%%s ]\n",
		result,
		shortCircuitValue,
		shortCircuit,
		rhs,
		evalRhs,
	)
	return result
}

// @TypeKind
func (c *Compiler) castOp(from node.Node, to node.Node) string {
	fromExpr := c.compileExpr(from, false)

	toType := to.GetType()
	fromType := from.GetType()
	if fromType.Equal(toType) {
		return fromExpr
	}

	command := ""

	llvmTo := llvmFormatType(toType)
	llvmFrom := llvmFormatType(fromType)

	intSize := func(k node.TypeKind) int {
		switch k {
		case node.TypeI8, node.TypeU8:
			return 8

		case node.TypeI16, node.TypeU16:
			return 16

		case node.TypeI32, node.TypeU32:
			return 32

		case node.TypeI64, node.TypeU64:
			return 64

		default:
			return -1
		}
	}

	toIntSize := intSize(toType.Kind)
	fromIntSize := intSize(fromType.Kind)

	if fromType.Ref != 0 || fromType.Kind == node.TypeRawptr {
		if toType.Ref != 0 || toType.Kind == node.TypeRawptr {
			// Pointer -> Pointer
			command = "bitcast"
		} else {
			// Pointer -> Integer
			command = "ptrtoint"
		}
	} else if fromType.Kind == node.TypeBool {
		// Boolean -> Integer
		command = "zext"
	} else if fromIntSize != -1 {
		if toType.Ref != 0 {
			// Integer -> Pointer
			command = "inttoptr"
		} else if toType.Kind == node.TypeBool {
			// Integer -> Boolean
			result := c.valueNew()
			fmt.Fprintf(c.out, "    %s = icmp ne %s %s, 0\n", result, llvmFrom, fromExpr)
			return result
		} else if toIntSize != -1 {
			// Integer -> Integer
			if fromIntSize == toIntSize {
				return fromExpr
			}

			if fromIntSize < toIntSize {
				// Extend
				if fromType.IsSignedInt() {
					command = "sext"
				} else {
					command = "zext"
				}
			} else {
				// Truncate
				command = "trunc"
			}
		} else {
			panic("unreachable")
		}
	} else {
		panic("unreachable")
	}

	result := c.valueNew()
	fmt.Fprintf(c.out, "    %s = %s %s %s to %s\n", result, command, llvmFrom, fromExpr, llvmTo)
	return result
}

// @NodeKind
func (c *Compiler) compileExpr(n node.Node, ref bool) string {
	switch n := n.(type) {
	case *node.Atom:
		if n.Token.IsInteger() {
			return fmt.Sprintf("%d", n.Token.Int)
		}

		// @TokenKind
		switch n.Token.Kind {
		case token.Bool:
			return fmt.Sprintf("%d", n.Token.Int)

		case token.Ident:
			switch def := n.Defined.(type) {
			case *node.Fn:
				ref = true

			case *node.Let:
				if def.Kind == node.LetArg {
					ref = true
				}
			}

			normalizedName := n.Defined.Literal().Str
			if ref {
				return fmt.Sprintf("%s", normalizedName)
			}

			llvmType := llvmFormatType(n.GetType())
			result := c.valueNew()

			fmt.Fprintf(c.out, "    %s = load %s, %s* %s\n", result, llvmType, llvmType, normalizedName)
			return result

		default:
			panic("unreachable")
		}

	case *node.Call:
		fn := c.compileExpr(n.Fn, false)
		args := []string{}

		for _, arg := range n.Args {
			args = append(args, c.compileExpr(arg, false))
		}

		result := ""

		fmt.Fprint(c.out, "    ")
		if !n.Type.Equal(node.Type{Kind: node.TypeUnit}) {
			result = c.valueNew()
			fmt.Fprintf(c.out, "%s = ", result)
		}

		fmt.Fprintf(c.out, "call %s(", llvmFormatType(n.Type))
		for i, arg := range n.Args {
			if i != 0 {
				fmt.Fprint(c.out, ", ")
			}

			fmt.Fprint(c.out, llvmFormatType(arg.GetType()))
		}
		fmt.Fprintf(c.out, ") %s(", fn)
		for i, arg := range args {
			if i != 0 {
				fmt.Fprint(c.out, ", ")
			}

			fmt.Fprintf(c.out, "%s %s", llvmFormatType(n.Args[i].GetType()), arg)
		}

		fmt.Fprintln(c.out, ")")
		return result

	case *node.Unary:
		// @TokenKind
		switch n.Token.Kind {
		case token.Sub:
			operand := c.compileExpr(n.Operand, false)
			result := c.valueNew()
			fmt.Fprintf(c.out, "    %s = sub %s 0, %s\n", result, llvmFormatType(n.Operand.GetType()), operand)
			return result

		case token.Mul:
			operand := c.compileExpr(n.Operand, false)
			if ref {
				return operand
			}

			llvmType := llvmFormatType(n.GetType())

			result := c.valueNew()
			fmt.Fprintf(c.out, "    %s = load %s, %s* %s\n", result, llvmType, llvmType, operand)

			return result

		case token.BAnd:
			return c.compileExpr(n.Operand, true)

		case token.BNot:
			operand := c.compileExpr(n.Operand, false)
			result := c.valueNew()
			fmt.Fprintf(c.out, "    %s = xor %s %s, -1\n", result, llvmFormatType(n.Operand.GetType()), operand)
			return result

		case token.LNot:
			operand := c.compileExpr(n.Operand, false)
			result := c.valueNew()
			fmt.Fprintf(c.out, "    %s = xor i1 %s, true\n", result, operand)
			return result

		default:
			panic("unreachable")
		}

	case *node.Binary:
		// @TokenKind
		switch n.Token.Kind {
		case token.Add:
			return c.binaryArithOp(n, "add")

		case token.Sub:
			return c.binaryArithOp(n, "sub")

		case token.Mul:
			return c.binaryArithOp(n, "mul")

		case token.Div:
			if n.Type.IsSignedInt() {
				return c.binaryArithOp(n, "sdiv")
			} else {
				return c.binaryArithOp(n, "udiv")
			}

		case token.Shl:
			return c.binaryArithOp(n, "shl")

		case token.Shr:
			if n.Type.IsSignedInt() {
				return c.binaryArithOp(n, "ashr")
			} else {
				return c.binaryArithOp(n, "lshr")
			}

		case token.BOr:
			return c.binaryArithOp(n, "or")

		case token.BAnd:
			return c.binaryArithOp(n, "and")

		case token.LOr:
			return c.binaryLogicalOp(n)

		case token.LAnd:
			return c.binaryLogicalOp(n)

		case token.Set:
			lhs := c.compileExpr(n.Lhs, true)
			rhs := c.compileExpr(n.Rhs, false)

			llvmType := llvmFormatType(n.Lhs.GetType())
			fmt.Fprintf(c.out, "    store %s %s, %s* %s\n", llvmType, rhs, llvmType, lhs)
			return ""

		case token.Gt:
			if n.Lhs.GetType().IsSignedInt() {
				return c.binaryOp(n, "icmp sgt")
			} else {
				return c.binaryOp(n, "icmp ugt")
			}

		case token.Ge:
			if n.Lhs.GetType().IsSignedInt() {
				return c.binaryOp(n, "icmp sge")
			} else {
				return c.binaryOp(n, "icmp uge")
			}

		case token.Lt:
			if n.Lhs.GetType().IsSignedInt() {
				return c.binaryOp(n, "icmp slt")
			} else {
				return c.binaryOp(n, "icmp ult")
			}

		case token.Le:
			if n.Lhs.GetType().IsSignedInt() {
				return c.binaryOp(n, "icmp sle")
			} else {
				return c.binaryOp(n, "icmp ule")
			}

		case token.Eq:
			return c.binaryOp(n, "icmp eq")

		case token.Ne:
			return c.binaryOp(n, "icmp ne")

		case token.As:
			return c.castOp(n.Lhs, n.Rhs)

		default:
			panic("unreachable")
		}

	case *node.Debug:
		operand := c.compileExpr(n.Operand, false)

		switch n.Token.Kind {
		case token.DebugAlloc:
			result := c.valueNew()
			fmt.Fprintf(
				c.out,
				"    %s = call i8* (i64) @malloc(i64 %s)\n",
				result,
				operand,
			)
			return result

		case token.DebugPrint:
			fmtPointer := c.valueNew()
			fmt.Fprintf(
				c.out,
				"    %s = getelementptr [5 x i8], [5 x i8]* @.print, i64 0, i64 0\n",
				fmtPointer,
			)

			c.valueNew() // For the call
			fmt.Fprintf(
				c.out,
				"    call i32 (i8*, ...) @printf(i8* %s, %s %s)\n",
				fmtPointer,
				llvmFormatType(n.Operand.GetType()),
				operand,
			)
			return ""

		default:
			panic("unreachable")
		}

	default:
		panic("unreachable")
	}
}

// @NodeKind
func (c *Compiler) compileStmt(n node.Node) {
	switch n := n.(type) {
	case *node.Block:
		for _, stmt := range n.Nodes {
			c.compileStmt(stmt)
		}

	case *node.If:
		consequent := c.labelNew()
		antecendent := c.labelNew()
		confluence := c.labelNew()

		condition := c.compileExpr(n.Condition, false)
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", condition, consequent, antecendent)

		fmt.Fprintf(c.out, "%s:\n", consequent)
		c.compileStmt(n.Consequent)
		fmt.Fprintf(c.out, "    br label %%%s\n", confluence)

		fmt.Fprintf(c.out, "%s:\n", antecendent)
		c.compileStmt(n.Antecedent)
		fmt.Fprintf(c.out, "    br label %%%s\n", confluence)

		fmt.Fprintf(c.out, "%s:\n", confluence)

	case *node.While:
		start := c.labelNew()
		body := c.labelNew()
		finally := c.labelNew()

		fmt.Fprintf(c.out, "    br label %%%s\n", start)
		fmt.Fprintf(c.out, "%s:\n", start)

		condition := c.compileExpr(n.Condition, false)
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", condition, body, finally)

		fmt.Fprintf(c.out, "%s:\n", body)
		c.compileStmt(n.Body)
		fmt.Fprintf(c.out, "    br label %%%s\n", start)

		fmt.Fprintf(c.out, "%s:\n", finally)

	case *node.Return:
		if n.Operand != nil {
			expr := c.compileExpr(n.Operand, false)
			fmt.Fprintf(c.out, "    ret %s %s\n", llvmFormatType(n.Type), expr)
		} else {
			fmt.Fprintln(c.out, "    ret void")
		}
		c.valueNew() // Apparently this is needed??

	case *node.Let:
		llvmType := llvmFormatType(n.Type)
		if n.Assign != nil {
			assign := c.compileExpr(n.Assign, false)
			fmt.Fprintf(c.out, "    store %s %s, %s* %s\n", llvmType, assign, llvmType, n.Token.Str)
		} else {
			if n.Type.Ref != 0 || n.Type.Kind == node.TypeFn || n.Type.Kind == node.TypeRawptr {
				fmt.Fprintf(c.out, "    store %s null, %s* %s\n", llvmType, llvmType, n.Token.Str)
			} else {
				fmt.Fprintf(c.out, "    store %s 0, %s* %s\n", llvmType, llvmType, n.Token.Str)
			}
		}

	default:
		c.compileExpr(n, false)
	}
}

// TODO: Test this
func ensureMainFunction(context *checker.Context) {
	if main, ok := context.Globals["main"]; ok {
		mainTok := main.Literal()
		mainType := main.GetType()

		if mainType.Kind != node.TypeFn {
			fmt.Fprintf(
				os.Stderr,
				"%s: ERROR: The identifier 'main' must be a function\n",
				mainTok.Pos,
			)
			os.Exit(1)
		}

		if mainType.Ref != 0 {
			fmt.Fprintf(
				os.Stderr,
				"%s: ERROR: The entry function 'main' cannot be a pointer\n",
				mainTok.Pos,
			)
			os.Exit(1)
		}

		mainFn := mainType.Spec.(*node.Fn)
		if len(mainFn.Args) != 0 {
			fmt.Fprintf(
				os.Stderr,
				"%s: ERROR: The entry function 'main' cannot take any arguments\n",
				mainTok.Pos,
			)
			os.Exit(1)
		}

		if mainFn.Return != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s: ERROR: The entry function 'main' cannot return anything\n",
				mainTok.Pos,
			)
			os.Exit(1)
		}

		// Yozi expects  fn main()
		// Clang expects fn main() i32
		mainFn.Token.Str = ".main"
		return
	}

	fmt.Fprintln(os.Stderr, "ERROR: The entry function 'main' has not been defined")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "+ fn main() {")
	fmt.Fprintln(os.Stderr, "+     // This function MUST be defined")
	fmt.Fprintln(os.Stderr, "+ }")
	os.Exit(1)
}

func normalizeGlobalNames(context *checker.Context) {
	for _, g := range context.Globals {
		switch g := g.(type) {
		case *node.Fn:
			g.Token.Str = "@" + g.Token.Str

		case *node.Let:
			g.Token.Str = "@" + g.Token.Str
		}
	}
}

func Program(context *checker.Context, exePath string) {
	ensureMainFunction(context)
	normalizeGlobalNames(context)

	asmPath := exePath + ".ll"

	var err error
	var compiler Compiler

	compiler.context = context
	compiler.out, err = os.Create(asmPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	// Compile the globals
	for _, g := range compiler.context.Globals {
		compiler.valueId = 0
		compiler.labelId = 0

		globalType := g.GetType()
		switch g := g.(type) {
		case *node.Fn:
			returnType := g.ReturnType()
			fmt.Fprintf(compiler.out, "define %s %s(", llvmFormatType(returnType), g.Token.Str)

			for i, arg := range g.Args {
				if i != 0 {
					fmt.Fprint(compiler.out, ", ")
				}

				arg.Token.Str = fmt.Sprintf("%%a%d", i)
				fmt.Fprintf(compiler.out, "%s %s", llvmFormatType(arg.Type), arg.Token.Str)
			}

			fmt.Fprintln(compiler.out, ") {")
			fmt.Fprintln(compiler.out, "$0:")

			for i, arg := range g.Args {
				if arg.Kind == node.LetLocalArg {
					arg.Token.Str = fmt.Sprintf("%%v%d", i)
					fmt.Fprintf(
						compiler.out,
						"    %s = alloca %s\n",
						arg.Token.Str,
						llvmFormatType(arg.Type),
					)

					llvmType := llvmFormatType(arg.Type)
					fmt.Fprintf(
						compiler.out,
						"    store %s %%a%d, %s* %s\n",
						llvmType,
						i,
						llvmType,
						arg.Token.Str,
					)
				}
			}

			for i, l := range g.Locals {
				// TODO: Assuming functions can't be nested
				switch l := l.(type) {
				case *node.Let:
					l.Token.Str = fmt.Sprintf("%%v%d", len(g.Args)+i)
					fmt.Fprintf(
						compiler.out,
						"    %s = alloca %s\n",
						l.Token.Str,
						llvmFormatType(l.Type),
					)
				}
			}

			compiler.compileStmt(g.Body)

			if returnType.Equal(node.Type{Kind: node.TypeUnit}) {
				fmt.Fprintln(compiler.out, "    ret void")
			}

			fmt.Fprintln(compiler.out, "}")

		case *node.Let:
			fmt.Fprintf(
				compiler.out,
				"%s = global %s ",
				g.Literal().Str,
				llvmFormatType(globalType),
			)

			if globalType.Ref != 0 || globalType.Kind == node.TypeFn || globalType.Kind == node.TypeRawptr {
				fmt.Fprintf(compiler.out, "null\n")
			} else {
				fmt.Fprintf(compiler.out, "0\n")
			}

		default:
			panic("unreachable")
		}
	}

	// @Temporary
	fmt.Fprintln(compiler.out, `@.print = private unnamed_addr constant [5 x i8] c"%ld\0A\00"`)
	fmt.Fprintln(compiler.out, "declare i32 @printf(i8*, ...)")
	fmt.Fprintln(compiler.out, "declare i8* @malloc(i64)")

	fmt.Fprintln(compiler.out, "define i32 @main() {")
	fmt.Fprintln(compiler.out, "$0:")

	// Assign the global variables
	for _, g := range compiler.context.Globals {
		if _, ok := g.(*node.Let); ok {
			compiler.compileStmt(g)
		}
	}

	fmt.Fprintln(compiler.out, "    call void @.main()")
	fmt.Fprintln(compiler.out, "    ret i32 0")
	fmt.Fprintln(compiler.out, "}")

	compiler.out.Close()

	cmd := exec.Command("clang", "-Wno-override-module", "-o", exePath, asmPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	// @Temporary: Turn this on later
	// os.Remove(asmPath)
}
