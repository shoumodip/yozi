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

	case node.TypeI64:
		sb.WriteString("i64")

	case node.TypeUnit:
		panic("unreachable")

	case node.TypeFn:
		// TODO: Function return
		sb.WriteString("void (")
		for i, arg := range t.Spec.(*node.Fn).Args {
			if i != 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(llvmFormatType(arg.GetType()))
		}
		sb.WriteString(")*")

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

// @NodeKind
func (c *Compiler) compileExpr(n node.Node, ref bool) string {
	switch n := n.(type) {
	case *node.Atom:
		// @TokenKind
		switch n.Token.Kind {
		case token.Int, token.Bool:
			return fmt.Sprintf("%d", n.Token.I64)

		case token.Ident:
			if _, ok := n.Defined.(*node.Fn); ok {
				ref = true
			}
			normalizedName := n.Defined.Literal().Str

			if ref || normalizedName[1] == 'a' {
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

		// TODO: Function return
		fmt.Fprint(c.out, "    call void(")
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
		return ""

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
			return c.binaryOp(n, "add")

		case token.Sub:
			return c.binaryOp(n, "sub")

		case token.Mul:
			return c.binaryOp(n, "mul")

		case token.Div:
			return c.binaryOp(n, "sdiv")

		case token.Shl:
			return c.binaryOp(n, "shl")

		case token.Shr:
			return c.binaryOp(n, "ashr")

		case token.BOr:
			return c.binaryOp(n, "or")

		case token.BAnd:
			return c.binaryOp(n, "and")

		case token.Set:
			lhs := c.compileExpr(n.Lhs, true)
			rhs := c.compileExpr(n.Rhs, false)

			llvmType := llvmFormatType(n.Lhs.GetType())
			fmt.Fprintf(c.out, "    store %s %s, %s* %s\n", llvmType, rhs, llvmType, lhs)
			return ""

		case token.Gt:
			return c.binaryOp(n, "icmp sgt")

		case token.Ge:
			return c.binaryOp(n, "icmp sge")

		case token.Lt:
			return c.binaryOp(n, "icmp slt")

		case token.Le:
			return c.binaryOp(n, "icmp sle")

		case token.Eq:
			return c.binaryOp(n, "icmp eq")

		case token.Ne:
			return c.binaryOp(n, "icmp ne")

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
	case *node.Print:
		operand := c.compileExpr(n.Operand, false)

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

	case *node.Block:
		for _, stmt := range n.Body {
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

	case *node.Let:
		llvmType := llvmFormatType(n.Type)
		if n.Assign != nil {
			assign := c.compileExpr(n.Assign, false)
			fmt.Fprintf(c.out, "    store %s %s, %s* %s\n", llvmType, assign, llvmType, n.Token.Str)
		} else {
			if n.Type.Ref != 0 || n.Type.Kind == node.TypeFn {
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

		if len(mainType.Spec.(*node.Fn).Args) != 0 {
			fmt.Fprintf(
				os.Stderr,
				"%s: ERROR: The entry function 'main' cannot take any arguments\n",
				mainTok.Pos,
			)
			os.Exit(1)
		}

		// TODO: Function return
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
			name := g.Token.Str
			if name == "@main" {
				// Yozi expects  fn main()
				// Clang expects fn main() i32
				name = "@.main"
			}

			fmt.Fprintf(compiler.out, "define void %s(", name)

			for i, arg := range g.Args {
				if i != 0 {
					fmt.Fprint(compiler.out, ", ")
				}

				arg.Token.Str = fmt.Sprintf("%%a%d", i)
				fmt.Fprintf(compiler.out, "%s %s", llvmFormatType(arg.Type), arg.Token.Str)
			}

			fmt.Fprintln(compiler.out, ") {")
			fmt.Fprintln(compiler.out, "$0:")

			for i, l := range g.Locals {
				// TODO: Assuming functions can't be nested
				switch l := l.(type) {
				case *node.Let:
					l.Token.Str = fmt.Sprintf("%%v%d", i)
					fmt.Fprintf(
						compiler.out,
						"    %s = alloca %s\n",
						l.Token.Str,
						llvmFormatType(l.Type),
					)
				}
			}

			compiler.compileStmt(g.Body)

			fmt.Fprintln(compiler.out, "    ret void")
			fmt.Fprintln(compiler.out, "}")

		case *node.Let:
			fmt.Fprintf(
				compiler.out,
				"%s = global %s ",
				g.Literal().Str,
				llvmFormatType(globalType),
			)

			if globalType.Ref != 0 || globalType.Kind == node.TypeFn {
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
