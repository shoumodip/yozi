package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"yozi/node"
	"yozi/token"
)

type Compiler struct {
	out     *os.File
	valueId int
}

func (c *Compiler) valueNew() string {
	c.valueId++
	return fmt.Sprintf("%%%d", c.valueId-1)
}

func (c *Compiler) binaryOp(n *node.Binary, op string) string {
	lhs := c.expr(n.Lhs)
	rhs := c.expr(n.Rhs)
	result := c.valueNew()

	fmt.Fprintf(c.out, "    %s = %s i64 %s, %s\n", result, op, lhs, rhs)
	return result
}

// @NodeKind
func (c *Compiler) expr(n node.Node) string {
	switch n := n.(type) {
	case *node.Atom:
		literal := n.Literal()

		// @TokenKind
		switch literal.Kind {
		case token.Int, token.Bool:
			return fmt.Sprintf("%d", literal.I64)

		default:
			panic("unreachable")
		}

	case *node.Unary:
		// @TokenKind
		switch n.Literal().Kind {
		case token.Sub:
			operand := c.expr(n.Operand)
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
		operand := c.expr(n.Operand)
		c.valueNew()
		fmt.Fprintf(c.out, "    call i32 (ptr, ...) @printf(ptr @.print, i64 %s)\n", operand)

	case *node.Block:
		for _, stmt := range n.Stmts {
			c.stmt(stmt)
		}

	case *node.If:
		conditionName := c.valueNew()
		consequentName := "L" + c.valueNew()[1:]
		antecendentName := "L" + c.valueNew()[1:]
		confluenceName := "L" + c.valueNew()[1:]
		condition := c.expr(n.Condition)
		fmt.Fprintf(c.out, "    %s = icmp ne i32 %s, 0\n", conditionName, condition)
		fmt.Fprintf(c.out, "    br i1 %s, label %%%s, label %%%s\n", conditionName, consequentName, antecendentName)
		fmt.Fprintf(c.out, "%s:\n", consequentName)
		c.stmt(n.Consequent)
		fmt.Fprintf(c.out, "    br label %%%s\n", confluenceName)
		fmt.Fprintf(c.out, "%s:\n", antecendentName)
		c.stmt(n.Antecedent)
		fmt.Fprintf(c.out, "    br label %%%s\n", confluenceName)
		fmt.Fprintf(c.out, "%s:\n", confluenceName)

	default:
		c.expr(n)
	}
}

func Program(nodes []node.Node, exePath string) {
	asmPath := exePath + ".ll"

	var err error
	var compiler Compiler

	compiler.out, err = os.Create(asmPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	// @Temporary
	fmt.Fprintln(compiler.out, `@.print = private unnamed_addr constant [5 x i8] c"%ld\0A\00"`)
	fmt.Fprintln(compiler.out, `declare i32 @printf(ptr, ...)`)
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
