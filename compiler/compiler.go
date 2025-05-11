package compiler

import (
	"fmt"
	"os"
	"os/exec"
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
func typeToLLVM(t node.Type) string {
	switch t.Kind {
	case node.TypeNil:
		panic("unreachable")

	case node.TypeBool:
		return "i1"

	case node.TypeI64:
		return "i64"

	default:
		panic("unreachable")
	}
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

	fmt.Fprintf(c.out, "    %s = %s i64 %s, %s\n", result, op, lhs, rhs)
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

			llvmType := typeToLLVM(n.GetType())
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

			llvmType := typeToLLVM(n.Lhs.GetType())
			fmt.Fprintf(c.out, "    store %s %s, %s* %s\n", llvmType, rhs, llvmType, lhs)
			return ""

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
		fmt.Fprintf(c.out, "    call i32 (ptr, ...) @printf(ptr @.print, i64 %s)\n", operand)

	case *node.Block:
		for _, stmt := range n.Body {
			c.stmt(stmt)
		}

	case *node.If:
		conditionName := c.valueNew()
		consequentName := c.labelNew()
		antecendentName := c.labelNew()
		confluenceName := c.labelNew()

		// TODO: condition is already a boolean
		condition := c.expr(n.Condition, false)
		fmt.Fprintf(c.out, "    %s = icmp ne i32 %s, 0\n", conditionName, condition)
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", conditionName, consequentName, antecendentName)
		fmt.Fprintf(c.out, "%s:\n", consequentName)
		c.stmt(n.Consequent)
		fmt.Fprintf(c.out, "    br label %%%s\n", confluenceName)
		fmt.Fprintf(c.out, "%s:\n", antecendentName)
		c.stmt(n.Antecedent)
		fmt.Fprintf(c.out, "    br label %%%s\n", confluenceName)
		fmt.Fprintf(c.out, "%s:\n", confluenceName)

	case *node.While:
		conditionName := c.valueNew()
		startName := c.labelNew()
		bodyName := c.labelNew()
		finallyName := c.labelNew()

		fmt.Fprintf(c.out, "    br label %%%s\n", startName)
		fmt.Fprintf(c.out, "%s:\n", startName)

		// TODO: condition is already a boolean
		condition := c.expr(n.Condition, false)
		fmt.Fprintf(c.out, "    %s = icmp ne i32 %s, 0\n", conditionName, condition)
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", conditionName, bodyName, finallyName)
		fmt.Fprintf(c.out, "%s:\n", bodyName)
		c.stmt(n.Body)
		fmt.Fprintf(c.out, "    br label %%%s\n", startName)
		fmt.Fprintf(c.out, "%s:\n", finallyName)

	case *node.Let:
		assign := c.expr(n.Assign, false)
		llvmType := typeToLLVM(n.GetType())
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

	for _, it := range compiler.context.Globals {
		fmt.Fprintf(compiler.out, "@%s = global %s 0\n", it.Literal().Str, typeToLLVM(it.GetType()))
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
