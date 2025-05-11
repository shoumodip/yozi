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

	case node.TypeNil:
		panic("unreachable")

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
	lhs := c.expr(n.Lhs, false)
	rhs := c.expr(n.Rhs, false)
	result := c.valueNew()

	fmt.Fprintf(c.out, "    %s = %s %s %s, %s\n", result, op, llvmFormatType(n.Lhs.GetType()), lhs, rhs)
	return result
}

// @NodeKind
func (c *Compiler) expr(n node.Node, ref bool) string {
	switch n := n.(type) {
	case *node.Atom:
		literal := n.Literal()

		// @TokenKind
		switch literal.Kind {
		case token.Int, token.Bool:
			return fmt.Sprintf("%d", literal.I64)

		case token.Ident:
			if ref {
				return fmt.Sprintf("@%s", n.Literal().Str)
			}

			llvmType := llvmFormatType(n.GetType())
			result := c.valueNew()

			fmt.Fprintf(c.out, "    %s = load %s, %s* @%s\n", result, llvmType, llvmType, n.Literal().Str)
			return result

		default:
			panic("unreachable")
		}

	case *node.Unary:
		// @TokenKind
		switch n.Literal().Kind {
		case token.Sub:
			operand := c.expr(n.Operand, false)
			result := c.valueNew()
			fmt.Fprintf(c.out, "    %s = sub i64 0, %s\n", result, operand)
			return result

		case token.Mul:
			operand := c.expr(n.Operand, false)
			if ref {
				return operand
			}

			llvmType := llvmFormatType(n.GetType())

			result := c.valueNew()
			fmt.Fprintf(c.out, "    %s = load %s, %s* %s\n", result, llvmType, llvmType, operand)

			return result

		case token.BAnd:
			return c.expr(n.Operand, true)

		default:
			panic("unreachable")
		}

	case *node.Binary:
		// @TokenKind
		switch n.Literal().Kind {
		case token.Add:
			return c.binaryOp(n, "add")

		case token.Sub:
			return c.binaryOp(n, "sub")

		case token.Mul:
			return c.binaryOp(n, "mul")

		case token.Div:
			return c.binaryOp(n, "sdiv")

		case token.Set:
			lhs := c.expr(n.Lhs, true)
			rhs := c.expr(n.Rhs, false)

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
func (c *Compiler) stmt(n node.Node) {
	switch n := n.(type) {
	case *node.Print:
		operand := c.expr(n.Operand, false)
		c.valueNew()
		fmt.Fprintf(
			c.out,
			"    call i32 (ptr, ...) @printf(ptr @.print, %s %s)\n",
			llvmFormatType(n.Operand.GetType()),
			operand,
		)

	case *node.Block:
		for _, stmt := range n.Body {
			c.stmt(stmt)
		}

	case *node.If:
		consequent := c.labelNew()
		antecendent := c.labelNew()
		confluence := c.labelNew()

		condition := c.expr(n.Condition, false)
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", condition, consequent, antecendent)

		fmt.Fprintf(c.out, "%s:\n", consequent)
		c.stmt(n.Consequent)
		fmt.Fprintf(c.out, "    br label %%%s\n", confluence)

		fmt.Fprintf(c.out, "%s:\n", antecendent)
		c.stmt(n.Antecedent)
		fmt.Fprintf(c.out, "    br label %%%s\n", confluence)

		fmt.Fprintf(c.out, "%s:\n", confluence)

	case *node.While:
		start := c.labelNew()
		body := c.labelNew()
		finally := c.labelNew()

		fmt.Fprintf(c.out, "    br label %%%s\n", start)
		fmt.Fprintf(c.out, "%s:\n", start)

		condition := c.expr(n.Condition, false)
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", condition, body, finally)

		fmt.Fprintf(c.out, "%s:\n", body)
		c.stmt(n.Body)
		fmt.Fprintf(c.out, "    br label %%%s\n", start)

		fmt.Fprintf(c.out, "%s:\n", finally)

	case *node.Let:
		assign := c.expr(n.Assign, false)
		llvmType := llvmFormatType(n.GetType())
		fmt.Fprintf(c.out, "    store %s %s, %s* @%s\n", llvmType, assign, llvmType, n.Literal().Str)

	default:
		c.expr(n, false)
	}
}

func Program(context *checker.Context, nodes []node.Node, exePath string) {
	asmPath := exePath + ".ll"

	var err error
	var compiler Compiler

	compiler.context = context
	compiler.out, err = os.Create(asmPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	// @Temporary
	fmt.Fprintln(compiler.out, `@.print = private unnamed_addr constant [5 x i8] c"%ld\0A\00"`)
	fmt.Fprintln(compiler.out, `declare i32 @printf(ptr, ...)`)

	for _, variable := range compiler.context.Globals {
		variableType := variable.GetType()
		fmt.Fprintf(
			compiler.out,
			"@%s = global %s ",
			variable.Literal().Str,
			llvmFormatType(variableType),
		)

		if variableType.Ref != 0 {
			fmt.Fprintf(compiler.out, "null\n")
		} else {
			fmt.Fprintf(compiler.out, "0\n")
		}
	}

	// @Temporary
	fmt.Fprintln(compiler.out, `define i32 @main() {`)
	fmt.Fprintln(compiler.out, `$0:`)

	for _, n := range nodes {
		compiler.stmt(n)
	}

	// @Temporary
	fmt.Fprintln(compiler.out, `    ret i32 0`)
	fmt.Fprintln(compiler.out, `}`)

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
